package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	id            int    `json:"id" gorm:"primary_key`
	Username      string `json:"username"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password_hash string `json:"password_hash"`
}
