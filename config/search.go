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

func GetPathsByFileName2(path, filename string) (chan string, chan error, chan struct{}) {
	resultPaths := make(chan string)
	errors := make(chan error)
	quit := make(chan struct{})

	newDirs := make(chan string)

	wg := &sync.WaitGroup{}
	// todo: check the correctness of searching
	go func(wg *sync.WaitGroup, resultPaths chan string, errors chan error, newDirs chan string) {
		for dir := range newDirs {
			wg.Add(1)
			go func(wg *sync.WaitGroup, dir string) {
				defer wg.Done()
				getPathsByFileName(resultPaths, errors, newDirs, dir, filename)
			}(wg, dir)
		}
	}(wg, resultPaths, errors, newDirs)

	go func(wg *sync.WaitGroup, resultPaths chan string, errors chan error, newDirs chan string) {
		pathInfo, err := os.Stat(path)
		if err != nil {
			errors <- err
			return
		}

		if pathInfo.IsDir() {
			dir, err := os.Open(path)
			if err != nil {
				errors <- err
				return
			}
			dirNames, err := dir.Readdirnames(-1)
			if err != nil {
				errors <- err
				return
			}
			for _, dirName := range dirNames {
				newDirPath := filepath.Join(path, dirName)
				newDirs <- newDirPath
			}
			wg.Wait()
			quit <- struct{}{}
			close(resultPaths)
			close(newDirs)
			close(errors)
			close(quit)
		} else {
			getPathsByFileName(resultPaths, errors, newDirs, path, filename)
		}
	}(wg, resultPaths, errors, newDirs)

	return resultPaths, errors, quit
}

func getPathsByFileName(resultPaths chan string, errors chan error, newDirs chan string, path, filename string) {
	info, err := os.Stat(path)
	if err != nil {
		errors <- err
		return
	}
	switch info.IsDir() {
	case true:
		f, err := os.Open(path)
		defer func() {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}()
		if err != nil {
			errors <- err
			return
		}
		filesInfo, err := f.Readdir(-1)
		if err != nil {
			errors <- err
			return
		}

		for _, fileInfo := range filesInfo {
			name := fileInfo.Name()
			path := filepath.Join(path, name)
			if name == filename {
				resultPaths <- path
			}
			if fileInfo.IsDir() {
				newDirs <- path
			}
		}
	case false:
		if info.Name() == filename {
			resultPaths <- filepath.Join(path, info.Name())
		}
	}
}
