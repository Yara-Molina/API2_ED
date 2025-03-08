package controllers

import (
	"net/http"
	domain "notifications/src/domain/entities"
	"notifications/src/domain/repositories"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	repo repositories.NotificationRepository
}

func NewNotificationController(repo repositories.NotificationRepository) *NotificationController {
	return &NotificationController{repo: repo}
}

func (ctrl *NotificationController) GetAll(c *gin.Context) {
	notifications, err := ctrl.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener notificaciones"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func (ctrl *NotificationController) Create(c *gin.Context) {
	var notification domain.Notification

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	if err := ctrl.repo.Create(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la notificación"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Notificación creada exitosamente"})
}
