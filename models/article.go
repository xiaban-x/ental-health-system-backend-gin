package models

import (
	"time"
)

// Article 文章模型
type Article struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
	Title       string     `json:"title" gorm:"not null"`
	Content     string     `json:"content" gorm:"type:text;not null"`
	Summary     string     `json:"summary" gorm:"type:text"`
	Cover       string     `json:"cover"`
	Category    string     `json:"category" gorm:"not null"`
	Tags        string     `json:"tags"`
	Status      string     `json:"status" gorm:"default:draft"`
	ViewCount   int        `json:"view_count" gorm:"default:0"`
	IsTop       bool       `json:"is_top" gorm:"default:false"`
	IsRecommend bool       `json:"is_recommend" gorm:"default:false"`
	AuthorID    uint       `json:"author_id" gorm:"not null"`
	Author      User       `json:"author" gorm:"foreignKey:AuthorID"`
	PublishedAt *time.Time `json:"published_at"`
}
