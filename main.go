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

// @title           学生心理健康管理系统 API
// @version         1.0
// @description     学生心理健康管理系统的RESTful API服务
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 使用Bearer令牌进行身份验证
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
