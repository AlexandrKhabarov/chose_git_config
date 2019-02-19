package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
)

func ParseConfigFile(path string) ([][]byte, [][]byte, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, nil, err
	}

	return parseConfigFile(file)
}

func parseConfigFile(file io.Reader) ([][]byte, [][]byte, error) {
	fileScanner := bufio.NewScanner(file)

	userNames := make([][]byte, 0, 0)
	emails := make([][]byte, 0, 0)

LOOP:
	for fileScanner.Scan() {
		row := fileScanner.Bytes()
		if bytes.Contains(row, []byte("[user]")) {
			for fileScanner.Scan() {
				row = fileScanner.Bytes()
				switch {
				case bytes.Contains(row, []byte("[")):
					break LOOP
				case bytes.Contains(row, []byte("email")):
					index := bytes.Index(row, []byte("="))
					if index != -1 {
						if index+1 < len(row) {
							emails = append(emails, bytes.TrimSpace(row[index+1:]))
						}
					}
				case bytes.Contains(row, []byte("name")):
					index := bytes.Index(row, []byte("="))
					if index != -1 {
						if index+1 < len(row) {
							userNames = append(userNames, bytes.TrimSpace(row[index+1:]))
						}
					}
				}
			}
		}
	}
	return userNames, emails, nil
}

func GetUserNamesAndEmail(filePaths <-chan string, userEmail, userNames chan<- []byte, finish chan <- struct{}) {
	usersNamesChan := make(chan [][]byte)
	usersEmailChan := make(chan [][]byte)
	errChan := make(chan error)
	finishChan := make(chan struct{})

	go func() {
		for {
			select {
			case users := <-usersNamesChan:
				for _, u := range users {
					userNames <- u
				}
			case emails := <-usersEmailChan:
				for _, e := range emails {
					userEmail <- e
				}
			case err := <-errChan:
				if err != nil {
					fmt.Printf("%q\n", err)
				}
			case <-finishChan:
				close(finishChan)
				close(userNames)
				close(userEmail)
				finish <- struct{}{}
				return
			default:
			}
		}
	}()

	wg := &sync.WaitGroup{}

	for path := range filePaths {
		absPath, err := preparePath(path)
		if err != nil {
			errChan <- err
		}
		wg.Add(1)
		go func(path string, wg *sync.WaitGroup) {
			defer wg.Done()
			users, emails, err := ParseConfigFile(path)
			errChan <- err
			usersEmailChan <- emails
			usersNamesChan <- users

		}(absPath, wg)
	}

	wg.Wait()
	close(usersNamesChan)
	close(usersEmailChan)
	close(errChan)
	finishChan <- struct{}{}
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
