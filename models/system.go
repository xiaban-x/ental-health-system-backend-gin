package models

import (
	"time"
)

// Token 令牌表
type Token struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	Username  string    `json:"username"`
	TableName string    `gorm:"column:table_name" json:"table_name"`
	Role      string    `json:"role"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `gorm:"column:expired_at" json:"expired_at"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// Config 配置表
type Config struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Feedback 用户反馈
type Feedback struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Reply     string    `json:"reply"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
