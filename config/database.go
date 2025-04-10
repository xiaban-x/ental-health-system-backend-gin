package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"ental-health-system/models"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	// 从环境变量获取数据库配置
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// 构建连接字符串
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// 配置GORM日志
	gormConfig := &gorm.Config{}
	if gin.Mode() == gin.DebugMode {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	// 连接数据库
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	log.Println("数据库连接成功")

	// 自动迁移数据库模型
	// migrateModels()
}

// 自动迁移数据库模型
func migrateModels() {
	log.Println("开始数据库迁移...")

	err := DB.AutoMigrate(
		&models.User{},         // 用户基础信息
		&models.Student{},      // 学生信息
		&models.Counselor{},    // 咨询师信息
		&models.Appointment{},  // 咨询预约
		&models.TimeSlot{},     // 咨询时间段
		&models.ExamPaper{},    // 试卷
		&models.ExamQuestion{}, // 试题
		&models.ExamRecord{},   // 考试记录
		&models.Resource{},     // 资源（文章、视频等）
		&models.ResourceTag{},  // 资源标签关联
		&models.Tag{},          // 标签
		&models.Feedback{},     // 用户反馈
		&models.Config{},       // 系统配置
		&models.Token{},        // 用户令牌
		&models.ChunkInfo{},    // 分片上传信息
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库迁移完成")
}
