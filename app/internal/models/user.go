package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Login            string `json:"login"`
	Password         string `json:"password"`
	AuthToken        string `json:"auth_token" gorm:"column:auth_token"`
	AuthTokenExpired int64  `json:"auth_token_expired" gorm:"column:auth_token_expired"`
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{db: db}
}

type UserModel struct {
	db *gorm.DB
}
type UserSearchOptions struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	AuthToken string `json:"auth_token"`
}

func (bm *UserModel) GetUser(opt UserSearchOptions) (*User, *gorm.DB) {
	var getUser User
	db := bm.db
	if opt.Login != "" {
		db = db.Where("Login=?", opt.Login)
	}
	if opt.Password != "" {
		db = db.Where("Password=?", opt.Password)
	}
	if opt.AuthToken != "" {
		db = db.Where("AuthToken=?", opt.AuthToken)
	}

	db = db.Find(&getUser)
	if getUser.Login == "" {
		return nil, db
	}
	return &getUser, db
}
