package main

import (
	"fmt"
)

func main() {
	ip := "10.10.20.101:22"
	getDbNames(ip)
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
