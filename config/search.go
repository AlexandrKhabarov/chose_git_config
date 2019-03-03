package config

import (
	"bufio"
	"bytes"
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

func GetUserNamesAndEmail(filePaths chan string, userEmail, userNames chan []byte) error {
	wg := &sync.WaitGroup{}

	for path := range filePaths {
		absPath, err := preparePath(path)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(path string, wg *sync.WaitGroup) {
			defer wg.Done()
			users, emails, _ := ParseConfigFile(path)
			for _, email := range emails {
				userEmail <- email
			}
			for _, name := range users {
				userNames <- name
			}

		}(absPath, wg)
	}

	wg.Wait()
	return nil
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

func GetPathsByFileName(path, filename string) (chan string, chan error, chan struct{}) {
	resultsPaths := make(chan string)
	errors := make(chan error)
	quit := make(chan struct{})

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
		quit <- struct{}{}
		close(resultsPaths)
		close(errors)
		close(quit)
	}()

	return resultsPaths, errors, quit
}
