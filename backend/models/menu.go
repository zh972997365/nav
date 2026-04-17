package models

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID         uint          `json:"id" gorm:"primary_key"`
	Title      string        `json:"title"`
	Icon       string        `json:"icon"`
	IconColor  string        `json:"icon_color"`
	Status     string        `json:"status"`
	OrderID    int           `json:"order_id" gorm:"default:0"`
	ParentID   *uint         `json:"parent_id"`
	Apps       []Application `json:"apps" gorm:"foreignkey:MenuID"`
	CreatedAt  NullableTime  `json:"created_at"`
	UpdatedAt  NullableTime  `json:"updated_at"`
	gorm.Model `gorm:"-"`
}

func (menu *Menu) BeforeCreate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	menu.CreatedAt = NewNullableTime(time.Now().In(loc))
	menu.UpdatedAt = NewNullableTime(time.Now().In(loc))
	return
}

func (menu *Menu) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	menu.UpdatedAt = NewNullableTime(time.Now().In(loc))
	return
}
