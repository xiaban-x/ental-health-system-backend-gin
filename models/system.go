package models

import (
	"time"

	"gorm.io/gorm"
)

// Token 用户令牌表
type Token struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`      // 关联的用户ID
	Token     string         `gorm:"type:text;not null" json:"token"`    // JWT令牌
	Type      string         `gorm:"size:20;default:access" json:"type"` // 令牌类型：access/refresh
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`         // 过期时间
	LastUsed  *time.Time     `json:"last_used"`                          // 最后使用时间
	UserAgent string         `gorm:"size:255" json:"user_agent"`         // 用户代理
	ClientIP  string         `gorm:"size:50" json:"client_ip"`           // 客户端IP
	IsRevoked bool           `gorm:"default:false" json:"is_revoked"`    // 是否已撤销
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Token) TableName() string {
	return "tokens"
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
