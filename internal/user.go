package gophermarket

import (
	"fmt"
	"strings"
)

type User struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

func (u User) String() string {

	builder := strings.Builder{}

	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("Username: %s\n", u.Username))
	builder.WriteString(fmt.Sprintf("Password: %s\n", u.Password))

	return builder.String()
}
