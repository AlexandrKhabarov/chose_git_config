package main

import (
	"github.com/sachez/chose_git_config/cli"
	"github.com/sachez/chose_git_config/config"
	"os/exec"
	"os/user"
)

var defaultPaths = []string{
	"/etc/gitconfig",
	"~/.gitconfig",
	"~/.config/git/config",
	".git/config",
}

func main() {
	_, err := exec.LookPath("git")
	if err != nil {
		panic(err)
	}
	usr, err := user.Current()
	errHandler := config.NewErrorHandler()
	errHandler.Run()
	if err == nil {
		paths, errors := config.GetPathsByFileName(usr.HomeDir, "config")
		go func() {
			for _, path := range defaultPaths {
				paths <- path
			}
		}()
		errHandler.Handle(errors)
		userEmailChan, userNamesChan, errors := config.GetUserNamesAndEmail(paths)
		errHandler.Handle(errors)
		ui := cli.NewConsoleUI()
		errHandler.Handle(ui.Errors)
		ui.RunUI(userNamesChan, userEmailChan)
		errHandler.Quit()
	} else {
		panic(err)
	}
}
