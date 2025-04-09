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
	migrateModels()
}

// 自动迁移数据库模型
func migrateModels() {
	log.Println("开始数据库迁移...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Counselor{},
		&models.Admin{},
		&models.PsychTest{},
		&models.TestQuestion{},
		&models.TestResult{},
		&models.Appointment{},
		&models.CounselingRecord{},
		&models.Article{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库迁移完成")
}
