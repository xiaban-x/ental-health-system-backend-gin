package models

import (
	"time"
)

// Resource 资源表
type Resource struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	URL         string    `json:"url"`
	CoverImage  string    `gorm:"column:cover_image" json:"cover_image"`
	Type        string    `json:"type"`
	Duration    int       `json:"duration"`
	Size        int64     `json:"size"`
	Format      string    `json:"format"`
	AuthorID    int       `gorm:"column:author_id" json:"author_id"`
	AuthorName  string    `gorm:"column:author_name" json:"author_name"`
	ViewCount   int       `gorm:"column:view_count" json:"view_count"`
	LikeCount   int       `gorm:"column:like_count" json:"like_count"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// Tag 标签
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// ResourceTag 资源标签关联
type ResourceTag struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ResourceID int       `gorm:"column:resource_id" json:"resource_id"`
	TagID      int       `gorm:"column:tag_id" json:"tag_id"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
