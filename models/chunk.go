package models

import (
	"mime/multipart"
	"strconv"
	"time"
)

// ChunkInfo 分片信息实体类
type ChunkInfo struct {
	ID           uint                  `gorm:"primaryKey;column:id" json:"id"`                     // ID
	ChunkNumber  int                   `gorm:"column:chunk_number" json:"chunk_number"`            // 当前分片，从1开始
	ChunkSize    int                   `gorm:"column:chunk_size" json:"chunk_size"`                // 分片大小
	TotalSize    int                   `gorm:"column:total_size" json:"total_size"`                // 总大小
	Identifier   string                `gorm:"column:identifier" json:"identifier"`                // 文件标识
	Filename     string                `gorm:"column:filename" json:"filename"`                    // 文件名
	RelativePath string                `gorm:"column:relative_path" json:"relative_path"`          // 相对路径
	TotalChunks  int                   `gorm:"column:total_chunks" json:"total_chunks"`            // 总分片数
	FileType     string                `gorm:"column:file_type" json:"file_type"`                  // 文件类型
	ChunkPath    string                `gorm:"column:chunk_path" json:"chunk_path"`                // 分片在Minio中的路径
	Status       int                   `gorm:"column:status;default:0" json:"status"`              // 状态：0-上传中，1-上传完成
	CreatedAt    time.Time             `gorm:"column:created_at;autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt    time.Time             `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"` // 更新时间
	File         *multipart.FileHeader `gorm:"-" json:"-"`                                         // 分片文件（不存入数据库，仅用于传输）
}

// TableName 指定表名
func (ChunkInfo) TableName() string {
	return "chunk_info"
}

// SetChunkNumber 设置分片号，支持类型转换
func (c *ChunkInfo) SetChunkNumber(value interface{}) {
	switch v := value.(type) {
	case string:
		if num, err := strconv.Atoi(v); err == nil {
			c.ChunkNumber = num
		}
	case int:
		c.ChunkNumber = v
	case float64:
		c.ChunkNumber = int(v)
	}
}

// SetChunkSize 设置分片大小，支持类型转换
func (c *ChunkInfo) SetChunkSize(value interface{}) {
	switch v := value.(type) {
	case string:
		if num, err := strconv.Atoi(v); err == nil {
			c.ChunkSize = num
		}
	case int:
		c.ChunkSize = v
	case float64:
		c.ChunkSize = int(v)
	}
}

// SetTotalSize 设置总大小，支持类型转换
func (c *ChunkInfo) SetTotalSize(value interface{}) {
	switch v := value.(type) {
	case string:
		if num, err := strconv.Atoi(v); err == nil {
			c.TotalSize = num
		}
	case int:
		c.TotalSize = v
	case float64:
		c.TotalSize = int(v)
	}
}

// SetTotalChunks 设置总分片数，支持类型转换
func (c *ChunkInfo) SetTotalChunks(value interface{}) {
	switch v := value.(type) {
	case string:
		if num, err := strconv.Atoi(v); err == nil {
			c.TotalChunks = num
		}
	case int:
		c.TotalChunks = v
	case float64:
		c.TotalChunks = int(v)
	}
}
