package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func getDbNames(ip string) ([]string, []string) {
	downloadDB(ip)
	db, err := sql.Open("sqlite3", "./stash-go.sqlite")
	checkErr(err)

	scenes, err := db.Query("SELECT scenes.path, scenes.title FROM scenes")
	checkErr(err)

	// fmt.Println(scenes)
	var titles []string
	var paths []string
	var title string
	var path string
	for scenes.Next() {
		err = scenes.Scan(&path, &title)
		checkErr(err)

		if filepath.Base(path) != title {
			titles = append(titles, title)
			paths = append(paths, path)
			// fmt.Println("Mismatch on", title, "and", filepath.Base(path))
		}
		// fmt.Println(path, "--", title)
	}

	return titles, paths
	// var dbPerfNames []string
	// var dbStudioNames []string
	// var dbCategoryNames []string

	// downloadDB(ip)

	// db, err := sql.Open("sqlite3", "./stash-go.sqlite")
	// checkErr(err)

	// fmt.Println("Querying DB performers...")
	// perfRows, err := db.Query("SELECT performers.name FROM performers")
	// checkErr(err)
	// dbPerfNames = getFromDB(perfRows)

	// fmt.Println("Querying DB studios...")
	// studioRows, err := db.Query("SELECT studios.name FROM studios")
	// checkErr(err)
	// dbStudioNames = getFromDB(studioRows)

	// fmt.Println("Querying DB categories...")
	// categoryRows, err := db.Query("SELECT tags.name FROM tags")
	// checkErr(err)
	// dbCategoryNames = getFromDB(categoryRows)

	// return dbPerfNames, dbStudioNames, dbCategoryNames
}

func backupDbLocal(filename string) (int64, error) {
	fmt.Println("Backing up DB locally")
	sourceFileStat, err := os.Stat("./stash-go.sqlite")
	checkErr(err)

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("./stash-go.sqlite is not a regular file")
	}

	source, err := os.Open("./stash-go.sqlite")
	checkErr(err)
	defer source.Close()

	destination, err := os.Create("./backups/" + filename)
	checkErr(err)
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func backupDbRemote(ip string) {
	fmt.Println("Backing up DB on server")
	dt := time.Now()
	fn := "stash-go.sqlite." + dt.Format("2006-01-02") + "T" + dt.Format("15-04-05")
	mvPrefix := "touch /mnt/cache/appdata/NFS2/config/stash-go.sqlite && mv /mnt/cache/appdata/NFS2/config/stash-go.sqlite /mnt/cache/appdata/NFS2/config/backups/"
	cmd := mvPrefix + fn
	sendCmd(ip, cmd)
	backupDbLocal(fn)
}
