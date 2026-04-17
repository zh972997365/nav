package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	LoginAt   NullableTime   `json:"login_at"`
	Status    string         `json:"status"`
	IsAdmin   bool           `json:"is_admin"`
	Token     string         `json:"token"`
	CreatedAt NullableTime   `json:"created_at"`
	UpdatedAt NullableTime   `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	user.CreatedAt = NewNullableTime(time.Now().In(loc))
	user.UpdatedAt = NewNullableTime(time.Now().In(loc))
	return
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	user.UpdatedAt = NewNullableTime(time.Now().In(loc))
	return
}
