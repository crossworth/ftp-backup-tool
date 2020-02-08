# FTP Backup tool

A tool to make a complete backup of ftp accounts. It downloads **one file at time**, to avoid transfer errors. Mostly on shared servers where the server cannot handle to many connections.

It can be used to download **huge** ftp accounts, it supports file skipping when an identical (same size and modification date) is found.

## How to use

Its a command line program.

Example run:

```bash
chmod +x ftp-backup-tool
./ftp-backup-tool user@example.org mypassword ftp.example.org
```
By default it will try download the entire ftp account, but you can specify a few options.

* -p the port (default 21)
* -t connection timeout in seconds (default 30 seconds)
* -s the start path to download (default /)
* -d the download path (default to ./backup)
