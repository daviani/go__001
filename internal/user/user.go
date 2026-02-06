package user

import "fmt"

type User struct {
	Name string
	Age  int
}

func (u User) Greet() string {
	return fmt.Sprintf("Hello, I'm %s", u.Name)
}
