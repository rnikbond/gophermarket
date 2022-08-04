package gophermarket

import (
	"fmt"
)

type User struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

func (u User) String() string {

	s := "\n"
	s += fmt.Sprintf("Username: %s\n", u.Username)
	s += fmt.Sprintf("Password: %s\n", u.Password)

	return s
}
