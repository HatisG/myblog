package dao

import (
	"myblog/internal/model"

	"gorm.io/gorm"
)

type PostDAO struct {
	db *gorm.DB
}

func NewPostDAO(db *gorm.DB) *PostDAO {
	return &PostDAO{db: db}
}

func (d *PostDAO) Create(post *model.Post) error {
	return d.db.Create(post).Error
}

func (d *PostDAO) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := d.db.Preload("Tags").First(&post, id).Error
	return &post, err
}

func (d *PostDAO) FindAll(page, pageSize int, tagName string) ([]model.Post, int64, error) {
	offset := (page - 1) * pageSize
	query := d.db.Model(&model.Post{}).Preload("Tags")

	if tagName != "" {
		query = query.Joins("join post_tags on post_tags.post_id = posts.id").
			Joins("join tags on tags.id = post_tags.tag_id").
			Where("tags.name = ?", tagName)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var posts []model.Post
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&posts).Error
	return posts, total, err
}

func (d *PostDAO) Update(id uint, updates map[string]interface{}) error {
	return d.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (d *PostDAO) Delete(id uint) error {
	return d.db.Delete(&model.Post{}, id).Error
}

func (d *PostDAO) ClearTags(post *model.Post) error {
	return d.db.Model(post).Association("Tags").Clear()
}

func (d *PostDAO) AppendTags(post *model.Post, tags []model.Tag) error {
	return d.db.Model(post).Association("Tags").Append(tags)
}

func (d *PostDAO) Begin() *gorm.DB {
	return d.db.Begin()
}
