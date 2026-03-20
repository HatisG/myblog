package model

type Tag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:50;uniqueIndex;not null" json:"name"`
}
