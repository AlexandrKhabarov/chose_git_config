package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func UpdateUserInfo(path string, text []byte) error {
	// todo: Fix bug with growing file by path
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return err
	}
	start, end := getWritingRange(f)
	f.Close()
	f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	writeFromTo(f, start, end, text)

	return nil
}

func writeFromTo(rw io.ReadWriter, from, to int64, text []byte) error {
	content, err := ioutil.ReadAll(rw)
	if err != nil {
		return err
	}
	textLen := int64(len(text))
	contentTailLen := int64(len(content[to:]))

	buf := make([]byte, from+textLen+contentTailLen+1)
	copy(buf[:from], content[:from])
	copy(buf[from:from+textLen], text)
	copy(buf[from+textLen:from+textLen+contentTailLen], content[to:])
	buf[len(buf)-1] = '\n'
	rw.Write(buf)
	return nil
}

func getWritingRange(file io.Reader) (int64, int64) {
	// todo: Check For buffer Overloading in ReadBytes methods
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
	return int64(start), int64(end)
}
