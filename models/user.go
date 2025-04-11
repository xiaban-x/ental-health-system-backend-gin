package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户基础信息
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`                  // 用户名，唯一
	Password  string         `gorm:"size:100;not null" json:"-"`                                    // 密码，json中隐藏
	Name      string         `gorm:"size:50" json:"name"`                                           // 真实姓名
	Sex       string         `gorm:"size:10" json:"sex"`                                            // 性别
	Phone     string         `gorm:"size:20;uniqueIndex:idx_phone,where:phone <> ''" json:"phone"`  // 手机号，非空时唯一
	Email     string         `gorm:"size:100;uniqueIndex:idx_email,where:email <> ''" json:"email"` // 邮箱，非空时唯一
	Avatar    string         `gorm:"size:255" json:"avatar"`                                        // 头像URL
	Role      string         `gorm:"size:20;default:student" json:"role"`                           // 角色：student/counselor/admin
	Status    string         `gorm:"size:20;default:active" json:"status"`                          // 状态：active/inactive/blocked
	Remark    string         `gorm:"size:500" json:"remark"`                                        // 备注
	Student   *Student       `gorm:"foreignKey:UserID" json:"student,omitempty"`                    // 学生信息，一对一关系
	Counselor *Counselor     `gorm:"foreignKey:UserID" json:"counselor,omitempty"`                  // 咨询师信息，一对一关系
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`            // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`            // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                                // 软删除
}

// BeforeSave 在保存前对密码进行加密
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// ValidatePassword 验证密码
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Student 学生信息
type Student struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"uniqueIndex;not null" json:"user_id"`                                         // 关联的用户ID
	User           User           `gorm:"foreignKey:UserID" json:"-"`                                                  // 关联的用户信息
	StudentID      string         `gorm:"size:50;uniqueIndex:idx_student_id,where:student_id <> ''" json:"student_id"` // 学号，非空时唯一
	Major          string         `gorm:"size:100" json:"major"`                                                       // 专业
	ClassName      string         `gorm:"column:class_name;size:50" json:"class_name"`                                 // 班级
	Grade          string         `gorm:"size:20" json:"grade"`                                                        // 年级
	EnrollmentDate time.Time      `gorm:"column:enrollment_date" json:"enrollment_date"`                               // 入学日期
	GraduationDate time.Time      `gorm:"column:graduation_date" json:"graduation_date"`                               // 预计毕业日期
	Dormitory      string         `gorm:"size:50" json:"dormitory"`                                                    // 宿舍信息
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// Counselor 咨询师信息
type Counselor struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"uniqueIndex;not null" json:"user_id"`                                                               // 关联的用户ID
	User           User           `gorm:"foreignKey:UserID" json:"-"`                                                                        // 关联的用户信息
	Title          string         `gorm:"size:50" json:"title"`                                                                              // 职称
	Specialty      string         `gorm:"size:200" json:"specialty"`                                                                         // 专业领域
	Introduction   string         `gorm:"type:text" json:"introduction"`                                                                     // 个人简介
	Status         int            `gorm:"default:1" json:"status"`                                                                           // 状态：0-不可用 1-可用
	EmployeeID     string         `gorm:"column:employee_id;size:50;uniqueIndex:idx_employee_id,where:employee_id <> ''" json:"employee_id"` // 工号，非空时唯一
	Department     string         `gorm:"size:100" json:"department"`                                                                        // 所属部门
	OfficeLocation string         `gorm:"column:office_location;size:100" json:"office_location"`                                            // 办公室位置
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
