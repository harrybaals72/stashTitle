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

func getDbNames(ip string, stashInst string) ([]int, []string, []string) {
	downloadDB(ip, stashInst)
	db, err := sql.Open("sqlite3", "./stash-go.sqlite")
	checkErr(err)

	scenes, err := db.Query("SELECT scenes.id, scenes.path, scenes.title FROM scenes")
	checkErr(err)

	// fmt.Println(scenes)
	var ids []int
	var titles []string
	var paths []string

	var id int
	var title string
	var path string
	for scenes.Next() {
		err = scenes.Scan(&id, &path, &title)
		checkErr(err)

		if filepath.Base(path) != title {
			ids = append(ids, id)
			titles = append(titles, title)
			paths = append(paths, path)
			// fmt.Println("Mismatch on", title, "and", filepath.Base(path))
		}
		// fmt.Println(path, "--", title)
	}

	return ids, titles, paths
}

func correctTitles(ip string, stashInst string, ids []int, titles []string, paths []string) {
	db, err := sql.Open("sqlite3", "./stash-go.sqlite")
	checkErr(err)

	fmt.Println("\n-------------------------------------------------------")
	fmt.Println("Beginning DB Titles Update")

	stmt, err := db.Prepare("UPDATE scenes SET title=? WHERE id=?")
	checkErr(err)

	for i := 0; i < len(ids); i++ {
		correctTitle := filepath.Base(paths[i])
		_, err := stmt.Exec(correctTitle, ids[i])
		checkErr(err)

		fmt.Println("Changed id:", ids[i], "\t", titles[i], "\t to \t", correctTitle)
	}

	err = db.Close()
	checkErr(err)
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

	_ = os.Mkdir("./backups/", os.ModePerm)
	destination, err := os.Create("./backups/" + filename)
	checkErr(err)
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func backupDbRemote(ip string, stashInst string) {
	fmt.Println("Backing up DB on server")
	dt := time.Now()
	fn := "stash-go.sqlite." + dt.Format("2006-01-02") + "T" + dt.Format("15-04-05") + ".sqlite"
	mvPrefix := "touch /mnt/cache/appdata/" + stashInst + "/config/stash-go.sqlite && mkdir -p /mnt/cache/appdata/" + stashInst + "/config/backups/ && mv /mnt/cache/appdata/" + stashInst + "/config/stash-go.sqlite /mnt/cache/appdata/" + stashInst + "/config/backups/"
	cmd := mvPrefix + fn
	sendCmd(ip, cmd)
	backupDbLocal(fn)
}
