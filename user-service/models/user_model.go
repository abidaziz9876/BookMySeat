package models

import "time"

type UserModel struct {
	ID            int64     `gorm:"column:id"`
	FirstName     string    `json:"first_name" gorm:"column:first_name"`
	LastName      string    `json:"last_name" gorm:"column:last_name"`
	Phone         string    `json:"phone" gorm:"column:phone"`
	Email         string    `json:"email" gorm:"column:email"`
	Password      string    `json:"password" gorm:"column:password"`
	Token         string    `json:"token"`
	Refresh_Token string    `josn:"refresh_token"`
	Created_At    time.Time `json:"created_at" gorm:"created_at"`
	Updated_At    time.Time `json:"updtaed_at" gorm:"updated_at"`
}
