package main

import (
	"github.com/sachez/chose_git_config/cli"

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
	go func() {
		names := [][]byte{
			[]byte("Sasha"),
			[]byte("Petya"),
		}

		emails := [][]byte{
			[]byte("Sasha@Sasha.com"),
			[]byte("Petya@Petya.com"),
		}

		for _, name := range names {
			userNamesChan <- name
		}
		close(userNamesChan)

		for _, email := range emails {
			userEmailChan <- email
		}
		close(userEmailChan)
	}()

	cli.NewConsoleUI(userNamesChan, userEmailChan)
	// go GetUserNamesAndEmail(filePathsChan, userEmailChan, userNamesChan, finishChan)
}
