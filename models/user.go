package models

import (
	"time"

	"gorm.io/gorm"
)

// User 基础用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"size:50;uniqueIndex" json:"username"`
	Password  string         `gorm:"size:100" json:"-"`
	Email     string         `gorm:"size:100;uniqueIndex" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	RealName  string         `gorm:"size:50" json:"real_name"`
	Gender    string         `gorm:"size:10" json:"gender"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Role      string         `gorm:"size:20" json:"role"` // student, counselor, admin
	Status    string         `gorm:"size:20;default:'active'" json:"status"`
}

// Student 学生模型
type Student struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uint           `gorm:"uniqueIndex" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user"`
	StudentID    string         `gorm:"size:50;uniqueIndex" json:"student_id"`
	Grade        string         `gorm:"size:20" json:"grade"`
	Class        string         `gorm:"size:50" json:"class"`
	Major        string         `gorm:"size:100" json:"major"`
	College      string         `gorm:"size:100" json:"college"`
	Dormitory    string         `gorm:"size:50" json:"dormitory"`
	EmergContact string         `gorm:"size:50" json:"emerg_contact"`
	EmergPhone   string         `gorm:"size:20" json:"emerg_phone"`
}

// Counselor 心理咨询师模型
type Counselor struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	UserID        uint           `gorm:"uniqueIndex" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user"`
	EmployeeID    string         `gorm:"size:50;uniqueIndex" json:"employee_id"`
	Title         string         `gorm:"size:50" json:"title"`
	Qualification string         `gorm:"size:100" json:"qualification"`
	Specialty     string         `gorm:"size:255" json:"specialty"`
	Introduction  string         `gorm:"type:text" json:"introduction"`
	Office        string         `gorm:"size:100" json:"office"`
	Schedule      string         `gorm:"type:text" json:"schedule"`
}

// Admin 管理员模型
type Admin struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `gorm:"uniqueIndex" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user"`
	EmployeeID string         `gorm:"size:50;uniqueIndex" json:"employee_id"`
	Department string         `gorm:"size:100" json:"department"`
	Position   string         `gorm:"size:100" json:"position"`
	Permission string         `gorm:"type:text" json:"permission"`
}
