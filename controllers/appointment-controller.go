package controllers

import (
	"net/http"
	"strconv"
	"time"

	"ental-health-system/config"
	"ental-health-system/models"

	"github.com/gin-gonic/gin"
)

// CreateAppointmentRequest 创建预约请求
type CreateAppointmentRequest struct {
	CounselorID uint      `json:"counselor_id" binding:"required"`
	AppointTime time.Time `json:"appoint_time" binding:"required"`
	Duration    int       `json:"duration" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	Topic       string    `json:"topic" binding:"required"`
	Description string    `json:"description"`
}

// UpdateAppointmentRequest 更新预约请求
type UpdateAppointmentRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed canceled completed"`
	Remark string `json:"remark"`
}

// CreateCounselingRecordRequest 创建咨询记录请求
type CreateCounselingRecordRequest struct {
	Content   string `json:"content" binding:"required"`
	FollowUp  string `json:"follow_up"`
	IsPrivate bool   `json:"is_private"`
}

// CreateAppointment 创建咨询预约
// @Summary 创建咨询预约
// @Description 学生创建新的咨询预约
// @Tags 咨询预约
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAppointmentRequest true "预约信息"
// @Success 201 {object} models.Appointment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /appointments [post]
func CreateAppointment(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 绑定请求数据
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询学生信息
	var student models.Student
	result := config.DB.Where("user_id = ?", userID).First(&student)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "学生信息不存在"})
		return
	}

	// 查询咨询师信息
	var counselor models.Counselor
	result = config.DB.First(&counselor, req.CounselorID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "咨询师不存在"})
		return
	}

	// 创建预约
	appointment := models.Appointment{
		StudentID:   student.ID,
		CounselorID: counselor.ID,
		AppointTime: req.AppointTime,
		Duration:    req.Duration,
		Type:        req.Type,
		Topic:       req.Topic,
		Description: req.Description,
		Status:      "pending",
	}

	result = config.DB.Create(&appointment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建预约失败"})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

// GetAppointments 获取预约列表
// @Summary 获取预约列表
// @Description 根据用户角色获取相关的预约列表
// @Tags 咨询预约
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "预约状态筛选"
// @Success 200 {array} models.Appointment
// @Failure 401 {object} map[string]string
// @Router /appointments [get]
func GetAppointments(c *gin.Context) {
	// 从上下文中获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取查询参数
	status := c.Query("status")

	var appointments []models.Appointment
	query := config.DB

	// 根据状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 根据用户角色筛选
	switch userRole {
	case "student":
		// 查询学生ID
		var student models.Student
		result := config.DB.Where("user_id = ?", userID).First(&student)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "学生信息不存在"})
			return
		}
		query = query.Where("student_id = ?", student.ID)
	case "counselor":
		// 查询咨询师ID
		var counselor models.Counselor
		result := config.DB.Where("user_id = ?", userID).First(&counselor)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "咨询师信息不存在"})
			return
		}
		query = query.Where("counselor_id = ?", counselor.ID)
	case "admin":
		// 管理员可以查看所有预约
		break
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
		return
	}

	// 执行查询
	result := query.Preload("Student").Preload("Student.User").Preload("Counselor").Preload("Counselor.User").Find(&appointments)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取预约列表失败"})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

// GetAppointmentByID 获取预约详情
// @Summary 获取预约详情
// @Description 根据预约ID获取预约详细信息
// @Tags 咨询预约
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "预约ID"
// @Success 200 {object} models.Appointment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /appointments/{id} [get]
func GetAppointmentByID(c *gin.Context) {
	// 获取URL参数中的预约ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的预约ID"})
		return
	}

	// 从上下文中获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询预约
	var appointment models.Appointment
	result := config.DB.Preload("Student").Preload("Student.User").Preload("Counselor").Preload("Counselor.User").Preload("CounselingRecords").First(&appointment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "预约不存在"})
		return
	}

	// 检查权限
	if userRole != "admin" {
		var hasPermission bool
		if userRole == "student" {
			// 查询学生ID
			var student models.Student
			result := config.DB.Where("user_id = ?", userID).First(&student)
			if result.Error == nil && student.ID == appointment.StudentID {
				hasPermission = true
			}
		} else if userRole == "counselor" {
			// 查询咨询师ID
			var counselor models.Counselor
			result := config.DB.Where("user_id = ?", userID).First(&counselor)
			if result.Error == nil && counselor.ID == appointment.CounselorID {
				hasPermission = true
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问此预约"})
			return
		}
	}

	c.JSON(http.StatusOK, appointment)
}

// UpdateAppointment 更新预约状态
// @Summary 更新预约状态
// @Description 更新预约的状态信息
// @Tags 咨询预约
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "预约ID"
// @Param request body UpdateAppointmentRequest true "更新信息"
// @Success 200 {object} models.Appointment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /appointments/{id} [put]
func UpdateAppointment(c *gin.Context) {
	// 获取URL参数中的预约ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的预约ID"})
		return
	}

	// 绑定请求数据
	var req UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询预约
	var appointment models.Appointment
	result := config.DB.First(&appointment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "预约不存在"})
		return
	}

	// 检查权限
	var hasPermission bool
	if userRole == "admin" {
		hasPermission = true
	} else if userRole == "student" {
		// 查询学生ID
		var student models.Student
		result := config.DB.Where("user_id = ?", userID).First(&student)
		if result.Error == nil && student.ID == appointment.StudentID {
			// 学生只能取消自己的预约
			hasPermission = req.Status == "canceled"
		}
	} else if userRole == "counselor" {
		// 查询咨询师ID
		var counselor models.Counselor
		result := config.DB.Where("user_id = ?", userID).First(&counselor)
		if result.Error == nil && counselor.ID == appointment.CounselorID {
			// 咨询师可以确认、取消或完成预约
			hasPermission = true
		}
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新此预约"})
		return
	}

	// 更新预约状态
	appointment.Status = req.Status
	appointment.Remark = req.Remark

	result = config.DB.Save(&appointment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新预约失败"})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

// CreateCounselingRecord 创建咨询记录
// @Summary 创建咨询记录
// @Description 为指定预约创建咨询记录
// @Tags 咨询记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "预约ID"
// @Param request body CreateCounselingRecordRequest true "咨询记录"
// @Success 201 {object} models.CounselingRecord
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /appointments/{id}/records [post]
func CreateCounselingRecord(c *gin.Context) {
	// 获取URL参数中的预约ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的预约ID"})
		return
	}

	// 绑定请求数据
	var req CreateCounselingRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询预约
	var appointment models.Appointment
	result := config.DB.First(&appointment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "预约不存在"})
		return
	}

	// 检查权限（只有咨询师可以创建咨询记录）
	if userRole != "counselor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有咨询师可以创建咨询记录"})
		return
	}

	// 查询咨询师ID
	var counselor models.Counselor
	result = config.DB.Where("user_id = ?", userID).First(&counselor)
	if result.Error != nil || counselor.ID != appointment.CounselorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限为此预约创建咨询记录"})
		return
	}

	// 创建咨询记录
	record := models.CounselingRecord{
		AppointmentID: uint(id),
		Content:       req.Content,
		FollowUp:      req.FollowUp,
		IsPrivate:     req.IsPrivate,
		CreatedBy:     counselor.ID,
	}

	result = config.DB.Create(&record)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建咨询记录失败"})
		return
	}

	c.JSON(http.StatusCreated, record)
}
