package main

import (
	"bytes"
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	mockFile := bytes.NewBuffer(make([]byte, 0, 1024))
	_, err := mockFile.WriteString("[user]\n\temail = user@test.com\n\tname = TestName")
	if err != nil {
		t.Errorf("[ERROR in TestParseConfigFile] During Writing To Mock File err is occurred: %v", err)
	}
	expectedUserName := []byte("TestName")
	expectedEmail := []byte("user@test.com")

	userNames, emails, err := parseConfigFile(mockFile)

	if err != nil {
		t.Errorf("[ERROR in TestParseConfigFile] During Parsing Mock File err is occurred: %q", err)
	}
	userName := userNames[0]
	if bytes.Compare(userName, expectedUserName) != 0 {
		t.Errorf("[ERROR in TestParseConfigFile]\nExpected UserName: %q\nActual UserName: %q", expectedUserName, userName)
	}
	email := emails[0]
	if bytes.Compare(email, expectedEmail) != 0 {
		t.Errorf("[ERROR in TestParseConfigFile]\nExpected UserName: %q\nActual UserName: %q", expectedEmail, email)
	}

}
