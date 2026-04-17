package models

type SiteSettings struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	Title          string       `json:"title"`
	Icon           string       `json:"icon"`
	IconColor      string       `json:"icon_color"`
	Footer         string       `json:"footer"`
	Icp            string       `json:"icp"`
	ImageHostToken string       `json:"image_host_token"`
	CreatedAt      NullableTime `json:"created_at"`
	UpdatedAt      NullableTime `json:"updated_at"`
}
