package main

import (
	"github.com/sachez/chose_git_config/cli"
	"github.com/sachez/chose_git_config/config"
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
			config.GetPathsByFileName(filePathsChan, usr.HomeDir, "config")
		}
		close(filePathsChan)
	}()
	go func() {
		// TODO: Add logging
		config.GetUserNamesAndEmail(filePathsChan, userEmailChan, userNamesChan)
		close(userEmailChan)
		close(userNamesChan)
	}()

	cli.NewConsoleUI(userNamesChan, userEmailChan)
}
