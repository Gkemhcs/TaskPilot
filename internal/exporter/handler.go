package exporter

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	"github.com/Gkemhcs/taskpilot/internal/middleware"
	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func NewExportHandler(service *ExportService, logger *logrus.Logger) *ExportHandler {
	return &ExportHandler{
		service: service,
		logger:  logger,
	}
}

type ExportHandler struct {
	service *ExportService
	logger  *logrus.Logger
}

func RegisterExportHandler(handler *ExportHandler, router *gin.RouterGroup, jwtManager *auth.JWTManager) {
	exportGroup := router.Group("/export")
	exportGroup.Use(middleware.JWTAuthMiddleware(handler.logger, jwtManager))
	{
		exportGroup.POST("/projects", handler.ExportProject)
		exportGroup.POST("/tasks", handler.ExportTask)
		exportGroup.GET("/status/:jobId", handler.GetExportStatus)
	}

}

// ExportProject handles the creation of a new project export job.
// @Summary      Export project to Excel
// @Description  Creates an export job for a project's data in Excel format
// @Tags         export
// @Accept       json
// @Produce      json
// @Param        request  body      ExportTaskRequest true  "Export project request"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      500   {object}  utils.ErrorResponse
// @Router       /api/v1/export/projects [post]
// @Security BearerAuth
func (exporter *ExportHandler) ExportProject(c *gin.Context) {

	var request ExportTaskRequest

	err := c.ShouldBindJSON(&request)

	val, exists := c.Get("userID")
	if !exists {
		exporter.logger.Errorf("%v", customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusInternalServerError, "unauthenticated: user ID not found")
		return
	}

	userID, ok := val.(int)
	if !ok {
		exporter.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, "invalid user ID type")
		return
	}
	// 1. Extract original file name and extension
	filename := "project-" + strconv.Itoa(userID) + ".xlsx"
	ext := filepath.Ext(filename)
	base := filename[:len(filename)-len(ext)]

	// 2. Add UUID or timestamp to make filename unique
	uniqueID := uuid.New().String()[:8] // or use time.Now().Unix()
	uniqueFilename := fmt.Sprintf("%s_%s%s", base, uniqueID, ext)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	jobID, err := exporter.service.ExportProjectExcel(ctx, uniqueFilename, userID)
	if err != nil {
		exporter.logger.Errorf("Error exporting project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Success(c, http.StatusCreated, map[string]string{
		"job_id":  jobID,
		"message": " Task Export job created successfully",
	})

}

// ExportTask handles the creation of a new task export job.
// @Summary      Export tasks to Excel
// @Description  Creates an export job for tasks of a project in Excel format
// @Tags         export
// @Accept       json
// @Produce      json
// @Param        request  body      ExportTaskRequest true  "Export task request"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      500   {object}  utils.ErrorResponse
// @Router       /api/v1/export/tasks [post]
// @Security BearerAuth
func (exporter *ExportHandler) ExportTask(c *gin.Context) {

	var request ExportTaskRequest

	err := c.ShouldBindJSON(&request)
	val, exists := c.Get("userID")
	if !exists {
		exporter.logger.Errorf("%v", customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusInternalServerError, "unauthenticated: user ID not found")
		return
	}

	userID, ok := val.(int)
	if !ok {
		exporter.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, "invalid user ID type")
		return
	}
	// 1. Extract original file name and extension
	filename := "task-" + strconv.Itoa(userID) + ".xlsx"
	ext := filepath.Ext(filename)
	base := filename[:len(filename)-len(ext)]

	// 2. Add UUID or timestamp to make filename unique
	uniqueID := uuid.New().String()[:8] // or use time.Now().Unix()
	uniqueFilename := fmt.Sprintf("%s_%s%s", base, uniqueID, ext)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	jobID, err := exporter.service.ExportTaskExcel(ctx, uniqueFilename, userID, request.ProjectID)
	if err != nil {
		exporter.logger.Errorf("Error exporting project: %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, map[string]string{
		"job_id":  jobID,
		"message": "Project Export job created successfully",
	})

}

// GetExportStatus returns the status of an export job by job ID.
// @Summary      Get export job status
// @Description  Retrieves the status and details of an export job
// @Tags         export
// @Produce      json
// @Param        jobId  path      string  true  "Export job ID"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      500   {object}  utils.ErrorResponse
// @Router       /api/v1/export/status/{jobId} [get]
// @Security BearerAuth
func (exporter *ExportHandler) GetExportStatus(c *gin.Context) {

	jobId := c.Param("jobId")
	userId, ok := c.Get("userID")
	if !ok {
		exporter.logger.Errorf("%v", customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrUserIDNotFoundInContext.Error())
		return
	}
	userID, ok := userId.(int)
	if !ok {
		exporter.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, customErrors.ErrInvalidUserId.Error())
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	exportJob, err := exporter.service.GetExportStatus(ctx, userID, jobId)
	if err != nil {
		exporter.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	exporter.logger.Infof("%s job status was %s", exportJob.ID, exportJob.Status)
	utils.Success(c, http.StatusOK, map[string]any{
		"data":    exportJob,
		"message": "request succeeded",
	})
}
