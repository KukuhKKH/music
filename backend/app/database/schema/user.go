package schema

import (
	"git.dev.siap.id/kukuhkkh/app-music/utils/helpers"
)

type User struct {
	ID       uint64 `gorm:"primary_key;column:id"`
	Name     string `gorm:"column:name;not null"`
	Password string `gorm:"column:password;not null"`
	Email    string `gorm:"column:email;unique;not null"`
	Base
}

// ComparePassword compare password
func (u *User) ComparePassword(password string) bool {
	return helpers.ValidateHash(password, u.Password)
}
