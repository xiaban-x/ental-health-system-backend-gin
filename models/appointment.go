package models

import (
	"time"
)

// Appointment 咨询预约
type Appointment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        int       `gorm:"column:user_id" json:"user_id"`
	Username      string    `json:"username"`
	CounselorID   int       `gorm:"column:counselor_id" json:"counselor_id"`
	CounselorName string    `gorm:"column:counselor_name" json:"counselor_name"`
	TimeSlotID    int       `gorm:"column:time_slot_id" json:"time_slot_id"`
	StartTime     time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime       time.Time `gorm:"column:end_time" json:"end_time"`
	Status        string    `json:"status"`
	Reason        string    `json:"reason"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TimeSlot 咨询时间段
type TimeSlot struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CounselorID int       `gorm:"column:counselor_id" json:"counselor_id"`
	StartTime   time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime     time.Time `gorm:"column:end_time" json:"end_time"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
