package main

import "fmt"

type User struct {
	UserName  []byte
	UserEmail []byte
}

func (user *User) UserRepresentation() string {
	return fmt.Sprintf("[user]\n\tuser = %s\n\temail = %s", user.UserName, user.UserEmail)
}
