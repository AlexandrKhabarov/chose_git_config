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
	userEmailChan := make(chan []byte, 1)
	userNamesChan := make(chan []byte, 1)
	filePathsChan := make(chan string)

	go func() {
		for _, path := range defaultPaths {
			filePathsChan <- path
		}
	}()
	go func() {
		usr, err := user.Current()
		if err == nil {
			// todo: Add logging
			GetPathsByFileName(filePathsChan, usr.HomeDir, "config")
		}
		close(filePathsChan)
	}()
	go func() {
		// todo: Add logging
		GetUserNamesAndEmail(filePathsChan, userEmailChan, userNamesChan)
		close(userEmailChan)
		close(userNamesChan)
	}() 

	NewConsoleUI(userNamesChan, userEmailChan)
}
