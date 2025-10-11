// Package model my-web-app/internal/model/storymap.go
package model

import (
	"encoding/json"
	"time"
)

// StoryMap 对应数据库中的 story_maps 表
type StoryMap struct {
	ID              uint            `gorm:"primarykey"`
	Title           string          `gorm:"type:varchar(255);not null"`
	Subtitle        string          `gorm:"type:varchar(255)"`
	TitleBackground string          `gorm:"type:varchar(512)"` // 存储标题背景图的URL
	Content         json.RawMessage `gorm:"type:json"`         // 存储富文本编辑器的完整JSON
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
