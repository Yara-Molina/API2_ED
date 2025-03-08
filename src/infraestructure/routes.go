package infraestructure

import (
	"notifications/src/application"
	"notifications/src/infraestructure/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, notificationUseCase *application.ProcessLoanUseCase) {
	// Configuración CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Crear instancia del controlador
	notificationController := controllers.NewNotificationController(notificationUseCase.Repo)

	// Definir rutas
	router.GET("/notifications", notificationController.GetAll)
	router.POST("/notifications", notificationController.Create)
}
