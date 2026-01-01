package schema

import (
	"git.dev.siap.id/kukuhkkh/app-music/utils/helpers"
)

type User struct {
	ID       uint64 `gorm:"primary_key;column:id" json:"id"`
	Name     string `gorm:"column:name;not null" json:"name"`
	Password string `gorm:"column:password;not null" json:"-"`
	Email    string `gorm:"column:email;unique;not null" json:"email"`
	Base
}

// ComparePassword compare password
func (u *User) ComparePassword(password string) bool {
	return helpers.ValidateHash(password, u.Password)
}
