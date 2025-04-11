package main

import (
	"log"
	"os"

	"ental-health-system/config"
	"ental-health-system/docs"
	"ental-health-system/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 心理健康系统 API
// @version 1.0
// @description 心理健康系统后端 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到.env文件，使用默认环境变量")
	}

	// 设置运行模式
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	// 初始化数据库连接
	config.InitDB()

	// 创建Gin实例
	r := gin.Default()

	// 配置CORS
	r.Use(config.CORSMiddleware())

	// 配置Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册路由
	routes.SetupRoutes(r)

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务器
	log.Printf("服务器启动在 http://localhost:%s", port)
	log.Printf("Swagger文档地址: http://localhost:%s/swagger/index.html", port)
	r.Run(":" + port)
}
