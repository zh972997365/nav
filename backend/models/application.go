package models

import (
	"time"

	"gorm.io/gorm"
)

type Application struct {
	ID          uint         `json:"id" gorm:"primary_key"`
	Title       string       `json:"title"`
	Icon        string       `json:"icon"`
	IconColor   string       `json:"icon_color"`
	Link        string       `json:"link"`
	Description string       `json:"description"`
	Status      string       `json:"status"`
	MenuID      uint         `json:"menu_id"`
	OrderID     int          `json:"order_id" gorm:"default:0"`
	CreatedAt   NullableTime `json:"created_at"`
	UpdatedAt   NullableTime `json:"updated_at"`
}

func (app *Application) BeforeCreate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	app.CreatedAt = NewNullableTime(time.Now().In(loc))
	app.UpdatedAt = NewNullableTime(time.Now().In(loc))
	return
}

func (app *Application) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	app.UpdatedAt = NewNullableTime(time.Now().In(loc))
	return
}
