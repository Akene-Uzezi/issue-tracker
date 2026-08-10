package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"issueTracking/internal/db"

	"github.com/gin-gonic/gin"
)

func (a *application) deathReport(c *gin.Context) {
	var deathReport db.DeathReport
	if err := c.ShouldBindJSON(&deathReport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A bad request was sent"})
		return
	}
	context := c.Request.Context()
	err := a.models.DeathReport.InsertDeathReport(context, &deathReport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notify superadmins, admins and managers about the new death report so it
	// surfaces in both frontends regardless of which app they are using.
	created, err := a.models.Notifications.CreateForRoles(context, []string{"superadmin", "admin", "manager"}, db.Notification{
		Type:      "death_report",
		Title:     "New death report submitted",
		Message:   fmt.Sprintf("Death report %s (%s) was submitted and is awaiting review.", deathReport.Ref, deathReport.Department),
		Ref:       deathReport.Ref,
		RelatedID: deathReport.ID,
	})
	if err != nil {
		log.Printf("ERROR: failed to create death report notifications: %v", err)
	} else {
		log.Printf("INFO: created %d death report notification(s)", created)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "The death has been reported"})
}

func (a *application) updateDeathReport(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "superadmin" && userRole != "admin" && userRole != "manager" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "you are unauthorized to view this"})
		return
	}
	var updateRequest db.DeathReport
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("A bad request was passed: %v", err.Error())})
		return
	}
	ctx := c.Request.Context()
	existingReport, err := a.models.DeathReport.SearchByID(ctx, updateRequest.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingReport == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := a.models.DeathReport.UpdateDeathReport(ctx, &updateRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The death report has been updated"})
}

func (a *application) getDeathReports(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "superadmin" && userRole != "admin" && userRole != "manager" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "you are unauthorized to view this"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit
	ctx := c.Request.Context()
	if dateFrom == "" && dateTo == "" {
		deathReports, totalPages, totalItems, err := a.models.DeathReport.GetDeathReports(ctx, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deathReports": PaginatedDeathReportResponse{
			Data: deathReports,
			Pagination: PaginationMeta{
				CurrentPage: page,
				PageSize:    limit,
				TotalItems:  totalItems,
				TotalPages:  totalPages,
			},
		}})
		return
	}
	deatReports, totalPages, totalItems, err := a.models.DeathReport.FetchDeathReportsByDate(ctx, limit, offset, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deathReports": PaginatedDeathReportResponse{
		Data: deatReports,
		Pagination: PaginationMeta{
			CurrentPage: page,
			PageSize:    limit,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}})
}

func (a *application) searchDeathReport(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "superadmin" && userRole != "admin" && userRole != "manager" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "you are unauthorized to view this"})
		return
	}
	searchQuery := c.Query("searchQuery")
	if searchQuery == "" {
		a.getDeathReports(c)
		return
	}
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")
	ctx := c.Request.Context()
	if dateFrom == "" && dateTo == "" {
		deathReports, err := a.models.DeathReport.SearchDeathReports(ctx, searchQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deathReports": deathReports})
		return
	}
	deathReports, err := a.models.DeathReport.SearchDeathReportsByDate(ctx, searchQuery, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deathReports": deathReports})
}
