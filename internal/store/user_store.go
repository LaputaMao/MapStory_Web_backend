// Package store my-web-app/internal/store/user_store.go
package store

import (
	"StoryMap/internal/model"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

// Create 创建一个新用户
func (s *UserStore) Create(user *model.User) error {
	return s.db.Create(user).Error
}

// FindByEmail 根据 Email 查找用户
func (s *UserStore) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (s *UserStore) Update(user *model.User) error {
	return s.db.Save(user).Error
}

// List 返回用户列表（分页和搜索）
func (s *UserStore) List(page, pageSize int, username string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// 创建一个查询构建器
	query := s.db.Model(&model.User{})

	// 如果提供了 username，添加模糊搜索条件
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	// 首先计算总数（在应用分页之前）
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 添加分页
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
