package main

import (
	"fmt"
	"strings"
)

func main() {
	ip := "10.10.20.101:22"
	stashInst := "NFS"

	ids, titles, paths := getDbNames(ip, stashInst)
	if len(ids) == len(titles) && len(ids) == len(paths) {
		fmt.Println("")
		for i := 0; i < len(ids); i++ {
			fmt.Println("Mismatch on id:", ids[i], "\t", titles[i], "\t", paths[i])
		}

		if len(ids) > 0 {
			var proceed string
			fmt.Println("\nProceed? y/n")
			fmt.Scanln(&proceed)

			if strings.ToLower(proceed) == "y" {
				stopContainer(ip, stashInst)
				backupDbRemote(ip, stashInst)
				correctTitles(ip, stashInst, ids, titles, paths)
				uploadDB(ip, stashInst)
				startContainer(ip, stashInst)
			}
		} else {
			fmt.Println("No mismatches found")
		}
	} else {
		fmt.Println("Slice length mismatch")
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func checkErrMsg(err error, msg string) {
	if err != nil {
		fmt.Println("Error:", msg)
		panic(err)
	}
}
