package main

import (
	"os/user"
)

var defaultPaths = []string{
	"/etc/gitconfig",
	"~/.gitconfig",
	"~/.config/git/config",
	".git/config",
}

func main() {
	userEmailChan := make(chan []byte)
	userNamesChan := make(chan []byte)
	filePathsChan := make(chan string)

	go func() {
		for _, path := range defaultPaths {
			filePathsChan <- path
		}
	}()
	go func() {
		usr, err := user.Current()
		if err == nil {
			// TODO: Add logging
			GetPathsByFileName(filePathsChan, usr.HomeDir, "config")
		}
		close(filePathsChan)
	}()
	go func() {
		// TODO: Add logging
		GetUserNamesAndEmail(filePathsChan, userEmailChan, userNamesChan)
		close(userEmailChan)
		close(userNamesChan)
	}()

	NewConsoleUI(userNamesChan, userEmailChan)
}
