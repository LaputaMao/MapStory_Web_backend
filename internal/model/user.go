// Package model my-web-app/internal/model/user.go
package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 对应数据库中的 users 表
type User struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	Username      string     `gorm:"type:varchar(100);not null;unique" json:"username"`
	Email         string     `gorm:"type:varchar(100);not null;unique" json:"email"`
	Password      string     `gorm:"type:varchar(255);not null" json:"-"` // json:"-" 确保密码哈希不被序列化返回给前端
	Role          string     `gorm:"type:varchar(50);not null;default:'user'" json:"role"`
	CreatedAt     time.Time  `json:"createdAt"`
	LastLoginTime *time.Time `json:"lastLoginTime"` // 使用指针类型以允许 NULL 值
}

// BeforeCreate 是一个 GORM Hook，在创建记录之前自动调用
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return
}

// CheckPasswordHash 验证密码是否正确
func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
