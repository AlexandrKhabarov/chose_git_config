package main

import (
	"bytes"
	"testing"
)

func TestWriteFromTo(t *testing.T) {
	mockFile := bytes.NewBuffer(make([]byte, 0, 1024))
	_, err := mockFile.WriteString("row1\nrow2\nrow3\n[user]\n\temail = user@test.com\n\tname = TestName")
	var (
		from int64 = 15
		to   int64 = 61
	)
	if err != nil {
		t.Errorf("[ERROR in TestWriteFromTo] During Writing To Mock File err is occurred: %v", err)
	}
	text := []byte("[user]\n\temail = test@user.com\n\tname = Name")
	expectedBufferContent := []byte("row1\nrow2\nrow3\n[user]\n\temail = test@user.com\n\tname = Name")
	err = writeFromTo(mockFile, from, to, text)
	bufferContent := mockFile.Bytes()

	if err != nil {
		t.Errorf("[ERROR in TestWriteFromTo] Error is occured: %v", err)
	}
	if bytes.Compare(expectedBufferContent, bufferContent) != 0 {
		t.Errorf("[ERROR in TestWriteFromTo]\nExpected buffer content: %q\nActual byffer content: %q", expectedBufferContent, bufferContent)
	}

}

func TestGetWritingOffset(t *testing.T) {
	mockFile := bytes.NewBuffer(make([]byte, 0, 1024))
	_, err := mockFile.WriteString("row1\nrow2\nrow3\n[user]\n\temail = user@test.com\n\tname = TestName")
	var (
		expectedStart int64 = 15
		expectedEnd   int64 = 61
	)
	if err != nil {
		t.Errorf("[ERROR in TestGetWritingOffset] During Writing To Mock File err is occurred: %v", err)
	}
	start, end := getWritingRange(mockFile)

	if expectedStart != start {
		t.Errorf("[ERROR in TestGetWritingOffset]\nExpected start: %v\nActual Start: %v", expectedStart, start)
	}

	if expectedEnd != end {
		t.Errorf("[ERROR in TestGetWritingOffset]\nExpected end: %v\nActual end: %v", expectedEnd, end)
	}
}
