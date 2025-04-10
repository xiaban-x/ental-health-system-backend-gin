package routes

import (
	"ental-health-system/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// 用户相关路由
		users := v1.Group("/users")
		{
			users.POST("/", controllers.CreateUser)
			users.GET("/", controllers.GetUserList)
			users.GET("/:id", controllers.GetUserByID)
			users.PUT("/:id", controllers.UpdateUser)
			users.DELETE("/:id", controllers.DeleteUser)
		}

		// 预约相关路由
		appointments := v1.Group("/appointments")
		{
			appointments.POST("/", controllers.CreateAppointment)
			appointments.GET("/", controllers.GetAppointmentList)
			appointments.GET("/:id", controllers.GetAppointmentByID)
			appointments.PUT("/:id", controllers.UpdateAppointment)
			appointments.DELETE("/:id", controllers.DeleteAppointment)
		}
	}
}
