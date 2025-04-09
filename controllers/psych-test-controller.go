package controllers

import (
	"net/http"
	"strconv"
	"time"

	"ental-health-system/config"
	"ental-health-system/models"

	"github.com/gin-gonic/gin"
)

// CreatePsychTestRequest 创建心理测评请求
// @Description 创建心理测评的请求结构
type CreatePsychTestRequest struct {
	Title       string `json:"title" binding:"required" example:"抑郁症筛查量表"`
	Description string `json:"description" binding:"required" example:"用于评估个体是否存在抑郁症状的量表"`
	Category    string `json:"category" binding:"required" example:"心理健康筛查"`
	Duration    int    `json:"duration" binding:"required" example:"30"`
	Questions   []struct {
		QuestionNo int    `json:"question_no" binding:"required" example:"1"`
		Content    string `json:"content" binding:"required" example:"您是否经常感到心情低落？"`
		Type       string `json:"type" binding:"required" example:"single"`
		Options    string `json:"options" example:"[\"从不\",\"偶尔\",\"经常\",\"总是\"]"`
		Score      int    `json:"score" example:"5"`
		Required   bool   `json:"required" example:"true"`
	} `json:"questions" binding:"required"`
}

// SubmitTestResultRequest 提交测评结果请求
// @Description 提交心理测评结果的请求结构
type SubmitTestResultRequest struct {
	Answers string `json:"answers" binding:"required" example:"{\"1\":\"经常\",\"2\":\"偶尔\"}"`
}

// CreatePsychTest 创建心理测评
// @Summary 创建心理测评
// @Description 创建新的心理测评问卷
// @Tags 心理测评
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePsychTestRequest true "测评信息"
// @Success 201 {object} models.PsychTest
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /psych-tests [post]
func CreatePsychTest(c *gin.Context) {
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

	// 检查权限（只有咨询师和管理员可以创建测评）
	if userRole != "counselor" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限创建测评"})
		return
	}

	// 绑定请求数据
	var req CreatePsychTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建测评
	psychTest := models.PsychTest{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Duration:    req.Duration,
		Status:      "active",
		CreatorID:   userID.(uint),
	}

	// 开始事务
	tx := config.DB.Begin()

	// 创建测评
	result := tx.Create(&psychTest)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建测评失败"})
		return
	}

	// 创建测评问题
	for _, q := range req.Questions {
		question := models.TestQuestion{
			TestID:     psychTest.ID,
			QuestionNo: q.QuestionNo,
			Content:    q.Content,
			Type:       q.Type,
			Options:    q.Options,
			Score:      q.Score,
			Required:   q.Required,
		}

		result = tx.Create(&question)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建测评问题失败"})
			return
		}
	}

	// 提交事务
	tx.Commit()

	// 返回创建的测评
	c.JSON(http.StatusCreated, psychTest)
}

// GetPsychTests 获取心理测评列表
// @Summary 获取心理测评列表
// @Description 获取所有可用的心理测评列表
// @Tags 心理测评
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category query string false "测评类别筛选"
// @Success 200 {array} models.PsychTest
// @Failure 401 {object} map[string]string
// @Router /psych-tests [get]
func GetPsychTests(c *gin.Context) {
	// 获取查询参数
	category := c.Query("category")

	var psychTests []models.PsychTest
	query := config.DB.Where("status = ?", "active")

	// 根据类别筛选
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 执行查询
	result := query.Find(&psychTests)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取测评列表失败"})
		return
	}

	c.JSON(http.StatusOK, psychTests)
}

// GetPsychTestByID 获取心理测评详情
// @Summary 获取心理测评详情
// @Description 根据ID获取心理测评详细信息，包括问题
// @Tags 心理测评
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "测评ID"
// @Success 200 {object} models.PsychTest
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /psych-tests/{id} [get]
func GetPsychTestByID(c *gin.Context) {
	// 获取URL参数中的测评ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的测评ID"})
		return
	}

	// 查询测评
	var psychTest models.PsychTest
	result := config.DB.Preload("Questions").First(&psychTest, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测评不存在"})
		return
	}

	c.JSON(http.StatusOK, psychTest)
}

// SubmitTestResult 提交测评结果
// @Summary 提交测评结果
// @Description 学生提交心理测评的答案和结果
// @Tags 心理测评
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "测评ID"
// @Param request body SubmitTestResultRequest true "测评结果"
// @Success 201 {object} models.TestResult
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /psych-tests/{id}/results [post]
func SubmitTestResult(c *gin.Context) {
	// 获取URL参数中的测评ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的测评ID"})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 绑定请求数据
	var req SubmitTestResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询测评
	var psychTest models.PsychTest
	result := config.DB.First(&psychTest, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测评不存在"})
		return
	}

	// 查询学生信息
	var student models.Student
	result = config.DB.Where("user_id = ?", userID).First(&student)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "学生信息不存在"})
		return
	}

	// 计算测评得分和结果（实际实现会更复杂，这里简化处理）
	score := 0
	resultText := "正常"
	analysis := "您的心理状态良好，请继续保持积极的生活态度。"
	recommendation := "建议定期参加心理健康活动，保持良好的心态。"

	// 创建测评结果
	testResult := models.TestResult{
		TestID:         uint(id),
		StudentID:      student.ID,
		Answers:        req.Answers,
		Score:          score,
		Result:         resultText,
		Analysis:       analysis,
		Recommendation: recommendation,
		CompletedAt:    time.Now(),
	}

	result = config.DB.Create(&testResult)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存测评结果失败"})
		return
	}

	c.JSON(http.StatusCreated, testResult)
}

// GetTestResults 获取测评结果列表
// @Summary 获取测评结果列表
// @Description 获取当前用户的测评结果列表
// @Tags 心理测评
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.TestResult
// @Failure 401 {object} map[string]string
// @Router /test-results [get]
func GetTestResults(c *gin.Context) {
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

	var testResults []models.TestResult
	query := config.DB

	// 根据用户角色筛选
	if userRole == "student" {
		// 查询学生ID
		var student models.Student
		result := config.DB.Where("user_id = ?", userID).First(&student)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "学生信息不存在"})
			return
		}
		query = query.Where("student_id = ?", student.ID)
	} else if userRole == "counselor" || userRole == "admin" {
		// 咨询师和管理员可以查看所有结果
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
		return
	}

	// 执行查询
	result := query.Preload("Student").Preload("Student.User").Find(&testResults)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取测评结果失败"})
		return
	}

	c.JSON(http.StatusOK, testResults)
}

// GetTestResultByID 获取测评结果详情
// @Summary 获取测评结果详情
// @Description 根据ID获取测评结果详细信息
// @Tags 心理测评
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "结果ID"
// @Success 200 {object} models.TestResult
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /test-results/{id} [get]
func GetTestResultByID(c *gin.Context) {
	// 获取URL参数中的结果ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结果ID"})
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

	// 查询测评结果
	var testResult models.TestResult
	result := config.DB.Preload("Student").Preload("Student.User").First(&testResult, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测评结果不存在"})
		return
	}

	// 检查权限
	if userRole == "student" {
		// 学生只能查看自己的测评结果
		var student models.Student
		result := config.DB.Where("user_id = ?", userID).First(&student)
		if result.Error != nil || student.ID != testResult.StudentID {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看此测评结果"})
			return
		}
	} else if userRole != "counselor" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看测评结果"})
		return
	}

	c.JSON(http.StatusOK, testResult)
}
