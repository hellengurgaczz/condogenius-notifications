package models

type Notification struct {
	Id      uint   `json:"id" gorm:"primaryKey"`
	Contact string `json:"contact" gorm:"not null"`
	Message string `json:"message" gorm:"not null"`
}
