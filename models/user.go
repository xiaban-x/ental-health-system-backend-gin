package models

import (
	"time"
)

// User 用户基础信息
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Sex       string    `json:"sex"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// Student 学生信息
type Student struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         int       `gorm:"column:user_id" json:"user_id"`
	StudentID      string    `gorm:"column:student_id" json:"student_id"`
	Major          string    `json:"major"`
	ClassName      string    `gorm:"column:class_name" json:"class_name"`
	Grade          string    `json:"grade"`
	EnrollmentDate time.Time `gorm:"column:enrollment_date" json:"enrollment_date"`
	GraduationDate time.Time `gorm:"column:graduation_date" json:"graduation_date"`
	Dormitory      string    `json:"dormitory"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// Counselor 咨询师信息
type Counselor struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         int       `gorm:"column:user_id" json:"user_id"`
	Title          string    `json:"title"`
	Specialty      string    `json:"specialty"`
	Introduction   string    `json:"introduction"`
	Status         int       `json:"status"`
	EmployeeID     string    `gorm:"column:employee_id" json:"employee_id"`
	Department     string    `json:"department"`
	OfficeLocation string    `gorm:"column:office_location" json:"office_location"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
