package dao

import (
	"myblog/internal/model"

	"gorm.io/gorm"
)

type TagDAO struct {
	db *gorm.DB
}

func NewTagDAO(db *gorm.DB) *TagDAO {
	return &TagDAO{db: db}
}

func (d *TagDAO) FirstOrCreate(tagName string) (*model.Tag, error) {
	var tag model.Tag
	err := d.db.Where(model.Tag{Name: tagName}).FirstOrCreate(&tag).Error
	return &tag, err
}
