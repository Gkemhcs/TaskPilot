package importer

import (
	"context"
	"net/http"
	"time"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	"github.com/Gkemhcs/taskpilot/internal/middleware"
	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


func NewImportHandler(service ImportService, logger *logrus.Logger) *ImportHandler {
	return &ImportHandler{
		service: service,
		logger:  logger,

	}
}


type  ImportHandler struct{
	service ImportService
	logger *logrus.Logger

}


func RegisterImporterHandler(handler *ImportHandler ,router *gin.RouterGroup,jwtManager *auth.JWTManager){
	importerGroup := router.Group("/import")
	importerGroup.Use(middleware.JWTAuthMiddleware(handler.logger,jwtManager))
	{
		importerGroup.POST("/projects", handler.ImportProject)
	}
	

}

func(importer *ImportHandler) ImportProject(c *gin.Context){

	file, fileHeader, err := c.Request.FormFile("file")
    if err != nil {
		utils.Error(c,http.StatusBadRequest, "file is required")
        return
    }
    defer file.Close()


	val, exists := c.Get("userID")
	if !exists {
		importer.logger.Errorf("%v", customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusInternalServerError, "unauthenticated: user ID not found")
		return
	}

	userID, ok := val.(int)
	if !ok {
		importer.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, "invalid user ID type")
		return
	}


	ctx,cancel:=context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
    jobID, err := importer.service.ImportProjectExcel(ctx, file, fileHeader.Filename,userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
	utils.Success(c,http.StatusCreated,map[string]string{
		"job_id": jobID,
		"message": "Import job created successfully",
	})



}
