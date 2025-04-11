package controllers

import (
	"ental-health-system/config"
	"ental-health-system/models"
	"ental-health-system/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Name     string `json:"name"` // 可选
	Role     string `json:"role"` // 可选，默认为student
}

// @Summary 用户注册
// @Description 用户使用用户名和密码注册新账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param data body RegisterRequest true "注册信息"
// @Success 200 {object} map[string]interface{} "注册成功返回用户信息和token"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 409 {object} map[string]interface{} "用户名已存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的请求参数",
			"details": "用户名长度需在3-50之间，密码长度需在6-50之间",
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if result := config.DB.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 设置默认角色
	if req.Role == "" {
		req.Role = "student"
	}

	// 验证角色是否有效
	if req.Role != "student" && req.Role != "counselor" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户角色"})
		return
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Password: req.Password, // 密码会在 BeforeSave 钩子中自动加密
		Name:     req.Name,
		Role:     req.Role,
		Status:   "active",
	}

	// 保存用户
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 如果是学生，创建对应的学生记录
	if req.Role == "student" {
		student := models.Student{
			UserID: user.ID,
		}
		if err := config.DB.Create(&student).Error; err != nil {
			// 如果创建学生记录失败，删除已创建的用户
			config.DB.Delete(&user)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建学生信息失败"})
			return
		}
	}

	// 如果是咨询师，创建对应的咨询师记录
	if req.Role == "counselor" {
		counselor := models.Counselor{
			UserID: user.ID,
			Status: 1,
		}
		if err := config.DB.Create(&counselor).Error; err != nil {
			// 如果创建咨询师记录失败，删除已创建的用户
			config.DB.Delete(&user)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建咨询师信息失败"})
			return
		}
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	// 返回用户信息和token
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
			"role":     user.Role,
		},
	})
}

// @Summary 用户登录
// @Description 用户使用用户名和密码登录系统
// @Tags 认证
// @Accept json
// @Produce json
// @Param data body LoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{} "登录成功返回用户信息和token"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "用户名或密码错误"
// @Failure 403 {object} map[string]interface{} "账户已被禁用"
// @Router /login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}
	fmt.Printf("hello, %s\n", req)
	// 查找用户
	var user models.User
	result := config.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if !user.ValidatePassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 检查用户状态
	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{"error": "账户已被禁用"})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	// 获取过期时间（24小时后）
	expiresAt := time.Now().Add(24 * time.Hour)

	// 存储token到数据库
	tokenRecord := models.Token{
		UserID:    user.ID,
		Token:     token,
		Type:      "access",
		ExpiresAt: expiresAt,
		UserAgent: c.Request.UserAgent(),
		ClientIP:  c.ClientIP(),
	}

	if err := config.DB.Create(&tokenRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存token失败"})
		return
	}

	// 返回用户信息和token
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
			"role":     user.Role,
			"email":    user.Email,
		},
	})
}
