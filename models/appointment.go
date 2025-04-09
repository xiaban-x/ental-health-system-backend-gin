package models

import (
	"time"

	"gorm.io/gorm"
)

// Appointment 咨询预约模型
type Appointment struct {
	ID                uint               `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	DeletedAt         gorm.DeletedAt     `gorm:"index" json:"-"`
	StudentID         uint               `json:"student_id"`
	Student           Student            `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	CounselorID       uint               `json:"counselor_id"`
	Counselor         Counselor          `gorm:"foreignKey:CounselorID" json:"counselor,omitempty"`
	AppointTime       time.Time          `json:"appoint_time"`
	Duration          int                `json:"duration"` // 预约时长（分钟）
	Type              string             `gorm:"size:50" json:"type"`
	Topic             string             `gorm:"size:100" json:"topic"`
	Description       string             `gorm:"type:text" json:"description"`
	Status            string             `gorm:"size:20;default:'pending'" json:"status"` // pending, confirmed, canceled, completed
	Remark            string             `gorm:"type:text" json:"remark"`
	CounselingRecords []CounselingRecord `gorm:"foreignKey:AppointmentID" json:"counseling_records,omitempty"`
}

// CounselingRecord 咨询记录模型
type CounselingRecord struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	AppointmentID uint           `json:"appointment_id"`
	Content       string         `gorm:"type:text" json:"content"`
	FollowUp      string         `gorm:"type:text" json:"follow_up"`
	IsPrivate     bool           `gorm:"default:false" json:"is_private"` // 是否仅咨询师可见
	CreatedBy     uint           `json:"created_by"`                      // 记录创建者ID（咨询师）
}
