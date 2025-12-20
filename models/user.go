package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model           `json:"-"` // ซ่อน gorm.Model จาก JSON output
	Id                   uint       `json:"id" gorm:"primaryKey"`
	Username             string     `json:"username"`
	Name                 string     `json:"name"`
	Email                string     `json:"email" gorm:"unique"`
	PasswordHash         string     `json:"-"` // ซ่อน PasswordHash จาก JSON output
	PasswordResetToken   *string    `json:"-"`
	PasswordResetExpires *time.Time `json:"-"`
}
