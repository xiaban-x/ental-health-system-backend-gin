package models

import (
	"time"
)

// ExamPaper 试卷表
type ExamPaper struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Time        int       `json:"time"`
	Status      int       `json:"status"`
	UserID      int       `gorm:"column:user_id" json:"user_id"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// ExamQuestion 试题表
type ExamQuestion struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PaperID      int       `gorm:"column:paper_id" json:"paper_id"`
	PaperName    string    `gorm:"column:paper_name" json:"paper_name"`
	QuestionName string    `gorm:"column:question_name" json:"question_name"`
	Options      string    `json:"options"`
	Score        int       `json:"score"`
	Answer       string    `json:"answer"`
	Analysis     string    `json:"analysis"`
	Type         int       `json:"type"`
	Sequence     int       `json:"sequence"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// ExamRecord 考试记录
type ExamRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     int       `gorm:"column:user_id" json:"user_id"`
	PaperID    int       `gorm:"column:paper_id" json:"paper_id"`
	TotalScore int       `gorm:"column:total_score" json:"total_score"`
	Feedback   string    `json:"feedback"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
