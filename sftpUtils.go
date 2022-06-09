package main

import (
	"fmt"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func downloadDB(ip string) {
	stashInst := "NFS2"
	fmt.Println("Importing DB from", stashInst)
	dbSrc := "/mnt/cache/appdata/" + stashInst + "/config/stash-go.sqlite"
	dbSrcShm := "/mnt/cache/appdata/" + stashInst + "/config/stash-go.sqlite-shm"
	dbSrcWal := "/mnt/cache/appdata/" + stashInst + "/config/stash-go.sqlite-wal"

	dbDest := "./stash-go.sqlite"
	dbDestShm := "./stash-go.sqlite-shm"
	dbDestWal := "./stash-go.sqlite-wal"

	deleteExisting(dbDest)
	deleteExisting(dbDestShm)
	deleteExisting(dbDestWal)

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("Nickel427"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", ip, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	sftp, err := sftp.NewClient(client)
	if err != nil {
		panic("Session failed: " + err.Error())
	}
	defer sftp.Close()

	sftpDownload(sftp, dbSrc, dbDest)

	if sftpLstat(sftp, dbSrcShm) {
		sftpDownload(sftp, dbSrcShm, dbDestShm)
	}

	if sftpLstat(sftp, dbSrcWal) {
		sftpDownload(sftp, dbSrcWal, dbDestWal)
	}
}

func deleteExisting(file string) {
	_, err := os.Lstat(file)
	if err == nil {
		e := os.Remove(file)
		checkErr(e)
		fmt.Println("Removed", file)
	}
}

func sftpLstat(sftp *sftp.Client, src string) (present bool) {
	_, err := sftp.Lstat(src)
	if err != nil {
		fmt.Println("No", src, "found")
		return false
	}
	return true
}

func sftpDownload(sftp *sftp.Client, src string, dest string) {
	srcFile, err := sftp.Open(src)
	checkErrMsg(err, "Could not open file "+src)
	defer srcFile.Close()

	dstFile, err := os.Create(dest)
	checkErrMsg(err, "Could not open file "+dest)
	defer dstFile.Close()

	fmt.Println("Attempting download")
	_, err = dstFile.ReadFrom(srcFile)
	checkErr(err)

	_, err = os.Lstat(dest)
	checkErr(err)
	fmt.Println("Download successful")
}
