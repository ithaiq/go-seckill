package models

import "fmt"

type UserModel struct {
	Id   int
	Name string
}

func (u *UserModel) String() string {
	return fmt.Sprintf("id:%d name:%s", u.Id, u.Name)
}
