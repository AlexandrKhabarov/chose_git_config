package config

import "fmt"

type User struct {
	UserName  []byte
	UserEmail []byte
}

func (user *User) UserRepresentation() string {
	return fmt.Sprintf("[user]\n\tname = %s\n\temail = %s", user.UserName, user.UserEmail)
}
