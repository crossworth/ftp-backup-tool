package main

import (
	"log"

	"github.com/integrii/flaggy"

	ftpbackuptool "github.com/crossworth/ftp-backup-tool"
)

var (
	user         string
	password     string
	server       string
	port         uint   = 21
	timeout      uint   = 30
	startPath    string = "/"
	downloadPath string = "./backup"
)

func init() {
	flaggy.SetName("FTP Backup tool")
	flaggy.SetDescription("Make FTP backup the slow way, one file at time!")

	flaggy.DefaultParser.ShowVersionWithVersionFlag = false
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/crossworth/ftp-backup-tool"

	flaggy.AddPositionalValue(&user, "user", 1, true, "The ftp user")
	flaggy.AddPositionalValue(&password, "password", 2, true, "The ftp password")
	flaggy.AddPositionalValue(&server, "server", 3, true, "The ftp server")

	flaggy.UInt(&port, "p", "port", "The ftp port")
	flaggy.UInt(&timeout, "t", "timeout", "The ftp connection timeout in seconds")
	flaggy.String(&startPath, "s", "start", "The start path")
	flaggy.String(&downloadPath, "d", "download", "The download path")

	flaggy.Parse()
}

func main() {
	ftp, err := ftpbackuptool.New(user, password, server, port, timeout)
	if err != nil {
		log.Fatal(err)
	}

	err = ftp.Start(startPath, downloadPath)
	log.Fatal(err)
}
