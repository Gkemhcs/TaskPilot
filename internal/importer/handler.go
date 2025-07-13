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


type  ImportHandler struct{
	service *ImportService
	logger *logrus.Logger

}


func RegisterImporterHandler(handler *ImportHandler ,router *gin.RouterGroup,jwtManager *auth.JWTManager){
	importerGroup := router.Group("/import")
	importerGroup.Use(middleware.JWTAuthMiddleware(handler.logger,jwtManager))
	{
		importerGroup.POST("/projects", handler.ImportProject)
		importerGroup.POST("/tasks", handler.ImportTask)
		importerGroup.GET("/status/:jobId",handler.GetStatus)
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
		utils.Error(c, http.StatusInternalServerError,customErrors.ErrUserIDNotFoundInContext.Error())
		return
	}

	userID, ok := val.(int)
	if !ok {
		importer.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError,customErrors.ErrInvalidUserId.Error() )
		return
	}
	// 1. Extract original file name and extension
	origFilename := fileHeader.Filename
	ext := filepath.Ext(origFilename)
	base := origFilename[:len(origFilename)-len(ext)]

	// 2. Add UUID or timestamp to make filename unique
	uniqueID := uuid.New().String()[:8] // or use time.Now().Unix()
	uniqueFilename := fmt.Sprintf("%s_%s%s", base, uniqueID, ext)

	ctx,cancel:=context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
    jobID, err := importer.service.ImportProjectExcel(ctx, file, uniqueFilename,userID)
    if err != nil {
		importer.logger.Errorf("Error importing project: %v", err)
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }
	utils.Success(c,http.StatusCreated,map[string]string{
		"job_id": jobID,
		"message": "Import job created successfully",
	})





}




func (importer *ImportHandler) ImportTask(c *gin.Context) {	

file, fileHeader, err := c.Request.FormFile("file")
    if err != nil {
		utils.Error(c,http.StatusBadRequest, "file is required")
        return
    }
    defer file.Close()


	val, exists := c.Get("userID")
	if !exists {
		importer.logger.Errorf("%v", customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusInternalServerError,customErrors.ErrUserIDNotFoundInContext.Error() )
		return
	}

	userID, ok := val.(int)
	if !ok {
		importer.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, customErrors.ErrInvalidUserId.Error())
		return
	}
	importer.logger.Info(userID," : userid")
	// 1. Extract original file name and extension
	origFilename := fileHeader.Filename
	ext := filepath.Ext(origFilename)
	base := origFilename[:len(origFilename)-len(ext)]

	// 2. Add UUID or timestamp to make filename unique
	uniqueID := uuid.New().String()[:8] // or use time.Now().Unix()
	uniqueFilename := fmt.Sprintf("%s_%s%s", base, uniqueID, ext)

	ctx,cancel:=context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
    jobID, err := importer.service.ImportTaskExcel(ctx, file, uniqueFilename,userID)
    if err != nil {
		importer.logger.Errorf("Error importing project: %v", err)
        c.JSON(http.StatusInternalServerError,err.Error())
        return
    }
	utils.Success(c,http.StatusCreated,map[string]string{
		"job_id": jobID,
		"message": "Import job created successfully",
	})



}


func(importer *ImportHandler)GetStatus(c *gin.Context){

	job_id:=c.Param("jobId")

	userId,ok:=c.Get("userID")
	if !ok{
		importer.logger.Errorf("%v",customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c,http.StatusBadRequest,customErrors.ErrUserIDNotFoundInContext.Error())
		return 
	}
		userID, ok := userId.(int)
	if !ok {
		importer.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, customErrors.ErrInvalidUserId.Error())
		return
	}
	ctx,cancel:=context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()
	job,err:=importer.service.Getstatus(ctx,job_id,userID)
	if err!=nil{
		importer.logger.Errorf("%v",err)
		utils.Error(c,http.StatusBadRequest,err.Error())
		return 
	}
	importer.logger.Infof("%s job status was %s",job.ID,job.Status)
	utils.Success(c,http.StatusOK,map[string]any{
		"data":job ,
		"message":"request succeeded",
	})




}
