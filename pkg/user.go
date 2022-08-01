package gophermarket

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

type User struct {
	Username string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"-"`
}

func (u User) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', tabwriter.AlignRight)

	fmt.Fprintln(w, "\nUsername\t", u.Username)
	fmt.Fprintln(w, "Password\t", u.Password)
	fmt.Fprintln(w, "Token\t", u.Token)

	if err := w.Flush(); err != nil {
		return err.Error()
	}

	return buf.String()
}
