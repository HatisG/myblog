package service

import (
	"errors"
	"myblog/internal/dao"
	"myblog/internal/model"
	"strings"

	"gorm.io/gorm"
)

type PostService struct {
	postDAO *dao.PostDAO
	tagDAO  *dao.TagDAO
	db      *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{
		postDAO: dao.NewPostDAO(db),
		tagDAO:  dao.NewTagDAO(db),
		db:      db,
	}
}

type CreatePostReq struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags"`
}

type UpdatePostReq struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags"`
}

func (s *PostService) Create(req *CreatePostReq) (*model.Post, error) {
	post := &model.Post{
		Title:   req.Title,
		Content: req.Content,
	}

	// 处理标签
	var tags []model.Tag
	for _, tagName := range req.Tags {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}
		tag, err := s.tagDAO.FirstOrCreate(tagName)
		if err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}
	post.Tags = tags

	err := s.postDAO.Create(post)
	return post, err
}

func (s *PostService) GetByID(id uint) (*model.Post, error) {
	return s.postDAO.FindByID(id)
}

func (s *PostService) List(page, pageSize int, tagName string) ([]model.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.postDAO.FindAll(page, pageSize, tagName)
}

func (s *PostService) Update(id uint, req *UpdatePostReq) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新文章基本信息
	updates := map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
	}
	if err := tx.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 获取文章
	var post model.Post
	if err := tx.Preload("Tags").First(&post, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 清除旧标签
	if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// 添加新标签
	var tags []model.Tag
	for _, tagName := range req.Tags {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}
		var tag model.Tag
		if err := tx.Where(model.Tag{Name: tagName}).FirstOrCreate(&tag).Error; err != nil {
			tx.Rollback()
			return err
		}
		tags = append(tags, tag)
	}

	if len(tags) > 0 {
		if err := tx.Model(&post).Association("Tags").Append(tags); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *PostService) Delete(id uint) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取文章
	var post model.Post
	if err := tx.First(&post, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}

	// 清除标签关联
	if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// 删除文章
	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
