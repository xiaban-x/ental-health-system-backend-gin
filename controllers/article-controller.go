package controllers

import (
	"net/http"
	"strconv"
	"time"

	"ental-health-system/config"
	"ental-health-system/models"

	"github.com/gin-gonic/gin"
)

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Summary     string `json:"summary"`
	Cover       string `json:"cover"`
	Category    string `json:"category" binding:"required"`
	Tags        string `json:"tags"`
	Status      string `json:"status" binding:"required,oneof=draft published"`
	IsTop       bool   `json:"is_top"`
	IsRecommend bool   `json:"is_recommend"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Summary     string `json:"summary"`
	Cover       string `json:"cover"`
	Category    string `json:"category"`
	Tags        string `json:"tags"`
	Status      string `json:"status" binding:"omitempty,oneof=draft published archived"`
	IsTop       *bool  `json:"is_top"`
	IsRecommend *bool  `json:"is_recommend"`
}

// CreateArticle 创建文章
// @Summary 创建文章
// @Description 创建新的文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateArticleRequest true "文章信息"
// @Success 201 {object} models.Article
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /articles [post]
func CreateArticle(c *gin.Context) {
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

	// 检查权限（只有咨询师和管理员可以创建文章）
	if userRole != "counselor" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限创建文章"})
		return
	}

	// 绑定请求数据
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建文章
	now := time.Now()
	article := models.Article{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		Cover:       req.Cover,
		Category:    req.Category,
		Tags:        req.Tags,
		AuthorID:    userID.(uint),
		Status:      req.Status,
		IsTop:       req.IsTop,
		IsRecommend: req.IsRecommend,
	}

	// 如果状态为已发布，设置发布时间
	if req.Status == "published" {
		article.PublishedAt = &now
	}

	result := config.DB.Create(&article)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	c.JSON(http.StatusCreated, article)
}

// GetArticles 获取文章列表
// @Summary 获取文章列表
// @Description 获取文章列表，支持分页和筛选
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param category query string false "文章分类"
// @Param status query string false "文章状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {array} models.Article
// @Router /articles [get]
func GetArticles(c *gin.Context) {
	// 获取查询参数
	category := c.Query("category")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 验证页码和每页数量
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	query := config.DB.Model(&models.Article{})

	// 默认只查询已发布的文章
	if status == "" {
		status = "published"
	}

	// 应用筛选条件
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != "all" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取文章列表
	var articles []models.Article
	result := query.Preload("Author").Order("is_top DESC, published_at DESC").Offset(offset).Limit(pageSize).Find(&articles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      articles,
	})
}

// GetArticleByID 获取文章详情
// @Summary 获取文章详情
// @Description 根据ID获取文章详细信息
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} models.Article
// @Failure 404 {object} map[string]string
// @Router /articles/{id} [get]
func GetArticleByID(c *gin.Context) {
	// 获取URL参数中的文章ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 查询文章
	var article models.Article
	result := config.DB.Preload("Author").First(&article, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 增加浏览次数
	config.DB.Model(&article).UpdateColumn("view_count", article.ViewCount+1)

	c.JSON(http.StatusOK, article)
}

// UpdateArticle 更新文章
// @Summary 更新文章
// @Description 更新文章信息
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Param request body UpdateArticleRequest true "更新信息"
// @Success 200 {object} models.Article
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /articles/{id} [put]
func UpdateArticle(c *gin.Context) {
	// 获取URL参数中的文章ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
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

	// 查询文章
	var article models.Article
	result := config.DB.First(&article, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	if userRole != "admin" && article.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新此文章"})
		return
	}

	// 绑定请求数据
	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新文章
	updates := make(map[string]interface{})

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}
	if req.Status != "" {
		updates["status"] = req.Status
		// 如果状态从草稿变为已发布，设置发布时间
		if req.Status == "published" && article.Status == "draft" {
			now := time.Now()
			updates["published_at"] = &now
		}
	}
	if req.IsTop != nil {
		updates["is_top"] = *req.IsTop
	}
	if req.IsRecommend != nil {
		updates["is_recommend"] = *req.IsRecommend
	}

	result = config.DB.Model(&article).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	// 获取更新后的文章
	config.DB.Preload("Author").First(&article, id)

	c.JSON(http.StatusOK, article)
}

// DeleteArticle 删除文章
// @Summary 删除文章
// @Description 删除指定ID的文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	// 获取URL参数中的文章ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
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

	// 查询文章
	var article models.Article
	result := config.DB.First(&article, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	if userRole != "admin" && article.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此文章"})
		return
	}

	// 删除文章（软删除）
	result = config.DB.Delete(&article)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章已删除"})
}
