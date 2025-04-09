package routes

import (
	"ental-health-system/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有路由
func SetupRoutes(r *gin.Engine) {
	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 预约相关路由
		appointmentRoutes := v1.Group("/appointments")
		{
			appointmentRoutes.POST("/", controllers.CreateAppointment)
			// appointmentRoutes.GET("/", controllers.GetUserAppointments)
			appointmentRoutes.PUT("/:id", controllers.UpdateAppointment)
		}

		// 文章相关路由
		articleRoutes := v1.Group("/articles")
		{
			// articleRoutes.GET("/", controllers.GetAllArticles)
			// articleRoutes.GET("/:id", controllers.GetArticle)
			articleRoutes.POST("/", controllers.CreateArticle)
		}
	}
}
