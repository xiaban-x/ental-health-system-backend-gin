package routes

import (
	"ental-health-system/config"
	"ental-health-system/controllers"

	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetupRoutes(r *gin.Engine) {
	// 创建限流器
	limiter := config.NewIPRateLimiter(rate.Every(1*time.Second), 3)

	// 全局中间件
	r.Use(config.LoggerMiddleware())
	r.Use(config.CORSMiddleware())

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 公开路由（不需要认证）
		public := v1.Group("")
		public.Use(config.RateLimitMiddleware(limiter))
		{
			public.POST("/login", controllers.Login)
			public.POST("/register", controllers.Register)
		}

		// 需要认证的路由
		auth := v1.Group("")
		auth.Use(config.JWTMiddleware())
		{
			// 用户相关路由
			users := auth.Group("/users")
			{
				// users.GET("/profile", controllers.GetUserProfile)
				// users.PUT("/profile", controllers.UpdateUserProfile)

				// 管理员专用路由
				admin := users.Group("")
				admin.Use(config.RoleAuthMiddleware("admin"))
				{
					admin.GET("", controllers.GetUserList)
					admin.POST("", controllers.CreateUser)
					admin.PUT("/:id", controllers.UpdateUser)
					admin.DELETE("/:id", controllers.DeleteUser)
				}
			}

			// 学生专用路由
			student := auth.Group("/student")
			student.Use(config.RoleAuthMiddleware("student"))
			{
				// 添加学生专用路由
			}

			// 咨询师专用路由
			counselor := auth.Group("/counselor")
			counselor.Use(config.RoleAuthMiddleware("counselor"))
			{
				// 添加咨询师专用路由
			}

			// 预约相关路由
			appointments := auth.Group("/appointments")
			{
				appointments.POST("/", controllers.CreateAppointment)
				appointments.GET("/", controllers.GetAppointmentList)
				appointments.GET("/:id", controllers.GetAppointmentByID)
				appointments.PUT("/:id", controllers.UpdateAppointment)
				appointments.DELETE("/:id", controllers.DeleteAppointment)
			}
		}
	}
}
