package main

import (
	"net/http"
	"strconv"

	"issueTracking/internal/db"

	"github.com/gin-gonic/gin"
)

func (a *application) getNotifications(c *gin.Context) {
	userID := c.GetInt("userId")
	notifications, err := a.models.Notifications.GetForUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if notifications == nil {
		notifications = []db.Notification{}
	}
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

func (a *application) markNotificationsRead(c *gin.Context) {
	userID := c.GetInt("userId")
	if err := a.models.Notifications.MarkAllRead(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func (a *application) markNotificationRead(c *gin.Context) {
	userID := c.GetInt("userId")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification id"})
		return
	}
	if err := a.models.Notifications.MarkRead(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (a *application) clearNotifications(c *gin.Context) {
	userID := c.GetInt("userId")
	if err := a.models.Notifications.DeleteAll(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All notifications cleared"})
}
