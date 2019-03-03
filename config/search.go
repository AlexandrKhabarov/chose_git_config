package config

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
)

func GetPathsByFileName(path, filename string) (chan string, chan error) {
	resultsPaths := make(chan string)
	errors := make(chan error)

	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					errors <- err
				}
				if info.Name() == filename {
					resultsPaths <- path
				}
				return nil
			})

			if err != nil {
				errors <- err
			}
		}(wg)

		wg.Wait()
		close(resultsPaths)
		close(errors)
	}()

	return resultsPaths, errors
}

func GetUserNamesAndEmail(paths chan string) (chan []byte, chan []byte, chan error) {
	userEmails := make(chan []byte)
	userNames := make(chan []byte)
	errors := make(chan error)

	go func() {
		wg := &sync.WaitGroup{}

		for path := range paths {
			absPath, err := preparePath(path)
			if err != nil {
				errors <- err
				continue
			}
			wg.Add(1)
			go func(path string, wg *sync.WaitGroup) {
				defer wg.Done()
				name, err := getUserName(path)
				if err != nil {
					errors <- err
				}
				if name != nil {
					userNames <- name
				}
				email, err := getUserEmail(path)
				if err != nil {
					errors <- err
				}
				if userEmails != nil {
					userEmails <- email
				}
			}(absPath, wg)
		}
		wg.Wait()
		close(userNames)
		close(userEmails)
		close(errors)
	}()

	return userEmails, userNames, errors
}

func getUserName(path string) ([]byte, error) {
	cmd := exec.Command("git", "config", "-f", path, "user.name")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	name, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%v: %v", err.Error(), stderr.String())
	}
	return name, nil
}

func getUserEmail(path string) ([]byte, error) {
	cmd := exec.Command("git", "config", "-f", path, "user.email")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	email, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%v: %v", err.Error(), stderr.String())
	}
	return email, nil
}

func preparePath(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	dir := usr.HomeDir
	switch {
	case path == "~":
		path = dir
	case strings.HasPrefix(path, "~/"):
		path = filepath.Join(dir, path[2:])
	default:
		path, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}
	return path, err
}
