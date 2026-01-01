package seeds

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-music/utils/helpers"
	"gorm.io/gorm"
)

type UserSeeder struct {
	DB *gorm.DB
}

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{
		DB: db,
	}
}

func (s *UserSeeder) Seed(db *gorm.DB) error {
	user := schema.User{
		Name:     "Administrator",
		Email:    "admin@gmail.com",
		Password: helpers.Hash([]byte("password123")),
	}

	return db.Create(&user).Error
}

func (s *UserSeeder) Count() (int, error) {
	var count int64
	if err := s.DB.Model(&schema.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
