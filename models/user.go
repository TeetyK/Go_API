package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model   `json:"-"` // ซ่อน gorm.Model จาก JSON output
	Id           uint       `json:"id" gorm:"primaryKey"`
	Username     string     `json:"username" gorm:"unique"`
	Email        string     `json:"email" gorm:"unique"`
	Name         string     `json:"name"`
	Role         string     `json:"role"`
	PasswordHash string     `json:"-"` // ซ่อน PasswordHash จาก JSON output
}
