// Package model my-web-app/internal/model/image.go
package model

import "time"

// Image 对应数据库中的 images 表
type Image struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Path      string    `gorm:"type:varchar(255);not null"`
	Format    string    `gorm:"type:varchar(50)"`
	CreatedAt time.Time // GORM 会自动处理创建时间
}
