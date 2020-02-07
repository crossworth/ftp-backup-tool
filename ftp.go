package ftpbackuptool

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTPBackupTool struct {
	conn          *ftp.ServerConn
	path          []string
	totalFiles    int
	skippedFiles  int
	downloadFiles int
}

func New(user string, password string, server string, port uint, timeout uint) (*FTPBackupTool, error) {
	var err error
	f := &FTPBackupTool{}
	f.conn, err = ftp.Dial(fmt.Sprintf("%s:%d", server, port), ftp.DialWithTimeout(time.Duration(timeout)*time.Second))
	if err != nil {
		return f, err
	}

	err = f.conn.Login(user, password)
	if err != nil {
		return f, err
	}

	return f, err
}

func (f *FTPBackupTool) Start(startPath string, downloadPath string) error {
	start := time.Now()

	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		err = os.Mkdir(downloadPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create download dir: %s", err)
		}
	}

	err := os.Chdir(downloadPath)
	if err != nil {
		return fmt.Errorf("could not change dir: %s", err)
	}

	f.path = append(f.path, downloadPath)

	err = f.downloadDir(startPath)
	if err != nil {
		return err
	}

	if err := f.conn.Quit(); err != nil {
		return fmt.Errorf("could not close connection ftp: %s", err)
	}

	elapsed := time.Since(start)

	fmt.Printf("Backup completed, took: %f minutes\n", elapsed.Minutes())
	fmt.Printf("Files downloaded/skipped/total: %d/%d/%d\n", f.downloadFiles, f.skippedFiles, f.totalFiles)

	return nil
}

func (f *FTPBackupTool) downloadDir(path string) error {
	if filepath.Base(path) != "/" {
		f.path = append(f.path, filepath.Base(path))
	}

	err := f.conn.ChangeDir(path)
	if err != nil {
		return fmt.Errorf("could not change dir: %s", err)
	}

	if _, err := os.Stat("./" + path); os.IsNotExist(err) {
		err = os.Mkdir("./" + path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create download dir: %s", err)
		}
	}

	err = os.Chdir("./" + path)
	if err != nil {
		return fmt.Errorf("could not change dir: %s", err)
	}

	files, err := f.conn.List(".")
	if err != nil {
		return fmt.Errorf("could not list ftp files: %s", err)
	}

	for _, file := range files {
		if file.Type == ftp.EntryTypeFolder && (file.Name == "." || file.Name == "..") {
			continue
		}

		if file.Type == ftp.EntryTypeFolder {
			err := f.downloadDir(file.Name)
			if err != nil {
				log.Println(err)
			}
		}

		if file.Type == ftp.EntryTypeFile {
			err := f.downloadFile(file)
			if err != nil {
				log.Println(err)
			}
		}
	}

	err = f.conn.ChangeDirToParent()
	if err != nil {
		return fmt.Errorf("could not change to parent dir on ftp: %s", err)
	}

	err = os.Chdir("../")
	if err != nil {
		return fmt.Errorf("could not change to parent dir: %s", err)
	}

	if len(f.path) > 0 {
		f.path = f.path[:len(f.path)-1]
	}

	return nil
}

func (f *FTPBackupTool) downloadFile(file *ftp.Entry) error {
	f.totalFiles++

	basePath := strings.Join(f.path, "/")

	fileStat, err := os.Stat(file.Name)
	exists := !os.IsNotExist(err)

	if exists && fileStat.Size() == int64(file.Size) && fileStat.ModTime().Local() == file.Time.Local() {
		f.skippedFiles++
		log.Printf("File skip: %s/%s !\n", basePath, file.Name)
		return nil
	}

	r, err := f.conn.Retr(file.Name)
	if err != nil {
		return fmt.Errorf("could not download file %s/%s: %s", basePath, file.Name, err)
	}
	defer r.Close()

	out, err := os.Create(file.Name)
	if err != nil {
		return fmt.Errorf("could not create file %s/%s: %s", basePath, file.Name, err)
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	if err != nil {
		return fmt.Errorf("could not save file %s/%s: %s", basePath, file.Name, err)
	}

	err = os.Chtimes(file.Name, file.Time, file.Time)
	if err != nil {
		return fmt.Errorf("could not change file time %s/%s: %s", basePath, file.Name, err)
	}

	log.Printf("File saved: %s/%s!\n", basePath, file.Name)
	f.downloadFiles++
	return nil
}
