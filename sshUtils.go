package main

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

func sendCmd(ip string, cmd string) []string {
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("Nickel427"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	fmt.Println("Attempting to ssh into server...")
	client, err := ssh.Dial("tcp", ip, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	fmt.Println("Opening session...")
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	fmt.Println("Sending command:", cmd)
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		panic("Failed to run: " + err.Error())
	}

	return strings.Split(b.String(), "\n")
}
