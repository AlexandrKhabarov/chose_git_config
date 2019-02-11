package main

import (
	"testing"
)

func TestUserStringRepresentation(t *testing.T) {
	user := User{UserName: []byte("UserName"), UserEmail: []byte("UserEmail")}
	expectedStringUser := "[user]\n\tuser = UserName\n\temail = UserEmail"
	stringUser := user.UserRepresentation()

	if stringUser != expectedStringUser {
		t.Errorf("[ERROR in TestUserStringRepresentation]\nExpected: %v\nActual: %v", expectedStringUser, stringUser)
	}
}
