package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func UpdateUserInfo(path string, text []byte) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		return err
	}

	start, end := getWritingRange(f)

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	content, err := getNewContent(f, start, end, text)
	if err != nil {
		return err
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func getNewContent(rw io.ReadWriter, from, to int, text []byte) ([]byte, error) {
	content, err := ioutil.ReadAll(rw)
	if err != nil {
		return nil, err
	}
	textLen := len(text)
	contentTailLen := len(content[to:])

	buf := make([]byte, from+textLen+contentTailLen+1)
	copy(buf[:from], content[:from])
	copy(buf[from:from+textLen], text)
	copy(buf[from+textLen:from+textLen+contentTailLen], content[to:])
	buf[len(buf)-1] = '\n'

	return buf, nil
}

func getWritingRange(file io.Reader) (int, int) {
	// TODO: Check For buffer Overloading in ReadBytes methods
	var (
		r     = bufio.NewReader(file)
		start = 0
		end   = 0
	)
LOOP:
	for {
		line, err := r.ReadBytes('\n')
		if err == io.EOF {
			break LOOP
		}
		if bytes.Contains(line, []byte("[user]")) {
			end = start
			end += len(line)
			for {
				line, err := r.ReadBytes('\n')
				if err == io.EOF {
					end += len(line)
					break LOOP
				}
				if bytes.Contains(line, []byte("[")) {
					break LOOP
				} else {
					end += len(line)
				}
			}
		} else {
			start += len(line)
		}
	}
	return start, end
}
