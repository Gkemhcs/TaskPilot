package importer

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"time"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	"github.com/Gkemhcs/taskpilot/internal/middleware"
	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func NewImportHandler(service *ImportService, logger *logrus.Logger) *ImportHandler {
	return &ImportHandler{
		service: service,
		logger:  logger,
	}
}

type ImportHandler struct {
	service *ImportService
	logger  *logrus.Logger
}

func RegisterImporterHandler(handler *ImportHandler, router *gin.RouterGroup, jwtManager *auth.JWTManager) {
	importerGroup := router.Group("/import")
	importerGroup.Use(middleware.JWTAuthMiddleware(handler.logger, jwtManager))
	{
		importerGroup.POST("/projects", handler.ImportProject)
		importerGroup.POST("/tasks", handler.ImportTask)
		importerGroup.GET("/status/:jobId", handler.GetStatus)
	}

}

// ImportProject handles the import of project data from an uploaded Excel file.
// @Summary Import projects from Excel file
// @Description Upload an Excel file to import projects in bulk
// @Tags Import
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file to import"
// @Success 201 {object} map[string]string "Import job created successfully"
// @Failure 400 {object} map[string]string "File is required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router       /api/v1/import/projects [post]
// @Security BearerAuth
func (importer *ImportHandler) ImportProject(c *gin.Context) {

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	val, exists := c.Get("userID")
	if !exists {
		importer.logger.Errorf("%v", customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusInternalServerError, customErrors.ErrUserIDNotFoundInContext.Error())
		return
	}

	userID, ok := val.(int)
	if !ok {
		importer.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, customErrors.ErrInvalidUserId.Error())
		return
	}
	// 1. Extract original file name and extension
	origFilename := fileHeader.Filename
	ext := filepath.Ext(origFilename)
	base := origFilename[:len(origFilename)-len(ext)]

	// 2. Add UUID or timestamp to make filename unique
	uniqueID := uuid.New().String()[:8] // or use time.Now().Unix()
	uniqueFilename := fmt.Sprintf("%s_%s%s", base, uniqueID, ext)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	jobID, err := importer.service.ImportProjectExcel(ctx, file, uniqueFilename, userID)
	if err != nil {
		importer.logger.Errorf("Error importing project: %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, map[string]string{
		"job_id":  jobID,
		"message": "Import job created successfully",
	})

}

// ImportTask handles the import of task data from an uploaded Excel file.
// @Summary Import tasks from Excel file
// @Description Upload an Excel file to import tasks in bulk
// @Tags Import
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file to import"
// @Success 201 {object} map[string]string "Import job created successfully"
// @Failure 400 {object} map[string]string "File is required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router       /api/v1/import/tasks [post]
// @Security BearerAuth
func (importer *ImportHandler) ImportTask(c *gin.Context) {

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	val, exists := c.Get("userID")
	if !exists {
		importer.logger.Errorf("%v", customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusInternalServerError, customErrors.ErrUserIDNotFoundInContext.Error())
		return
	}

	userID, ok := val.(int)
	if !ok {
		importer.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, customErrors.ErrInvalidUserId.Error())
		return
	}
	importer.logger.Info(userID, " : userid")
	// 1. Extract original file name and extension
	origFilename := fileHeader.Filename
	ext := filepath.Ext(origFilename)
	base := origFilename[:len(origFilename)-len(ext)]

	// 2. Add UUID or timestamp to make filename unique
	uniqueID := uuid.New().String()[:8] // or use time.Now().Unix()
	uniqueFilename := fmt.Sprintf("%s_%s%s", base, uniqueID, ext)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	jobID, err := importer.service.ImportTaskExcel(ctx, file, uniqueFilename, userID)
	if err != nil {
		importer.logger.Errorf("Error importing project: %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, map[string]string{
		"job_id":  jobID,
		"message": "Import job created successfully",
	})

}

// GetStatus returns the status of an import job by job ID.
// @Summary Get import job status
// @Description Get the status of a project or task import job
// @Tags Import
// @Produce json
// @Param jobId path string true "Job ID"
// @Success 200 {object} map[string]interface{} "Request succeeded"
// @Failure 400 {object} map[string]string "Bad request or job not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router       /api/v1/import/status/{jobId} [get]
// @Security BearerAuth
func (importer *ImportHandler) GetStatus(c *gin.Context) {

	job_id := c.Param("jobId")

	userId, ok := c.Get("userID")
	if !ok {
		importer.logger.Errorf("%v", customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrUserIDNotFoundInContext.Error())
		return
	}
	userID, ok := userId.(int)
	if !ok {
		importer.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, customErrors.ErrInvalidUserId.Error())
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	job, err := importer.service.Getstatus(ctx, job_id, userID)
	if err != nil {
		importer.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	importer.logger.Infof("%s job status was %s", job.ID, job.Status)
	utils.Success(c, http.StatusOK, map[string]any{
		"data":    job,
		"message": "request succeeded",
	})

}
