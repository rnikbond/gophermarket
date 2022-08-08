package gophermarket

import (
	"fmt"
	"strings"
)

type User struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type Balance struct {
	Accrual   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (u User) String() string {

	builder := strings.Builder{}

	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("Username: %s\n", u.Username))
	builder.WriteString(fmt.Sprintf("Password: %s\n", u.Password))

	return builder.String()
}

func (b Balance) String() string {

	builder := strings.Builder{}

	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("Accrual: %f\n", b.Accrual))
	builder.WriteString(fmt.Sprintf("Withdrawn: %f\n", b.Withdrawn))

	return builder.String()
}
