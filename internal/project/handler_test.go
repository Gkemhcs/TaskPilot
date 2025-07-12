package project

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	"github.com/Gkemhcs/taskpilot/internal/middleware"
	projectdb "github.com/Gkemhcs/taskpilot/internal/project/gen"
	"github.com/Gkemhcs/taskpilot/internal/task"
	taskdb "github.com/Gkemhcs/taskpilot/internal/task/gen"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func SetupNewProjectHandler() (*ProjectHandler, *auth.JWTManager, *task.MockTaskRepo, *MockProjectRepo, *logrus.Logger) {
	//creating auth jwt manager to generate access tokens and validate using middleware
	params := auth.CreateJwtManagerParams{
		AccessTokenDuration:  10 * time.Minute,
		RefreshTokenDuration: 10 * time.Hour,
		AccessTokenKey:       "rnk3mkrk3rk3rk3",
		RefreshTokenKey:      "21ieh12iei21eji12e",
	}
	jwtManager := auth.NewJWTManager(params)

	// mocking a project repo
	projectMockRepo := new(MockProjectRepo)

	//mocking  a task repo
	taskMockRepo := new(task.MockTaskRepo)

	//initialising the task query service

	taskQueryService := task.NewTaskService(taskMockRepo)
	// initialising project service to attach to handler
	projectService := NewProjectService(projectMockRepo)

	// logger setup and disable logging to console/io

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	// initialising project handler
	projectHandler := NewProjectHandler(logger, projectService, taskQueryService)

	return projectHandler, jwtManager, taskMockRepo, projectMockRepo, logger

}

// Simulated duplicate error
func mockDuplicateError() error {
	return &pq.Error{Code: "23505"}
}

func TestCreateProjectHandler(t *testing.T) {

	projectHandler, jwtManager, _, projectMockRepo, logger := SetupNewProjectHandler()

	jwtToken, err := jwtManager.GenerateAccessToken(1234, "test-user", "gudi@gmail.com")
	require.NoError(t, err)
	jwtMiddleware := middleware.JWTAuthMiddleware(logger, jwtManager)

	testCases := []struct {
		testName            string
		requestBody         map[string]string
		mockSetup           func()
		expectedServiceCall bool
		expectedError       error
	}{
		{
			testName: "Valid project id",
			requestBody: map[string]string{
				"name":        "implementing rbac for taskpilot",
				"color":       "green",
				"description": "we need to implement rbac as a part of this project",
			},
			mockSetup: func() {
				projectParams := projectdb.CreateProjectParams{
					UserID:      int32(1234),
					Name:        "implementing rbac for taskpilot",
					Description: sql.NullString{String: "we need to implement rbac as a part of this project", Valid: true},
					Color:       projectdb.NullProjectColor{ProjectColor: projectdb.ProjectColorGREEN, Valid: true},
				}

				projectMockRepo.On("CreateProject", mock.Anything, projectParams).Return(
					projectdb.Project{
						UserID:      1234,
						Name:        "implementing rbac for taskpilot",
						Description: sql.NullString{String: "we need to implement rbac as a part of this project", Valid: true},
						Color:       projectdb.NullProjectColor{ProjectColor: projectdb.ProjectColorGREEN, Valid: true},
					}, nil)
			},
			expectedServiceCall: true,
			expectedError:       nil,
		},
		{
			testName: "duplicate  entry project ",
			requestBody: map[string]string{
				"name":        "implementing rbac for taskpilot",
				"color":       "green",
				"description": "we need to implement rbac as a part of this project",
			},
			mockSetup: func() {
				projectParams := projectdb.CreateProjectParams{
					UserID:      int32(1234),
					Name:        "implementing rbac for taskpilot",
					Description: sql.NullString{String: "we need to implement rbac as a part of this project", Valid: true},
					Color:       projectdb.NullProjectColor{ProjectColor: projectdb.ProjectColorGREEN, Valid: true},
				}

				projectMockRepo.On("CreateProject", mock.Anything, projectParams).Return(
					projectdb.Project{}, mockDuplicateError())
			},
			expectedServiceCall: true,
			expectedError:       customErrors.ErrProjectAlreadyExists},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			mockRepo.Calls = nil
			mockRepo.ExpectedCalls = nil

			tc.mockSetup()

			body, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			// Create request
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/projects/", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			// Recorder
			w := httptest.NewRecorder()

			// Set up full group routing like production
			r := gin.Default()
			apiGroup := r.Group("/api/v1")
			userGroup := apiGroup.Group("/projects")

			userGroup.POST("/", jwtMiddleware, projectHandler.CreateProject)
			r.ServeHTTP(w, req)
			if tc.expectedError != nil {

				assert.Error(t, tc.expectedError, w.Code)
			} else {
				assert.Equal(t, http.StatusCreated, w.Code)
			}
			if tc.expectedServiceCall {
				projectMockRepo.AssertCalled(t, "CreateProject", mock.Anything, mock.Anything)

			} else {
				projectMockRepo.AssertNotCalled(t, "CreateProject", mock.Anything, mock.Anything)

			}

		})
	}

}

func TestGetProjectByIdHandler(t *testing.T) {
	projectHandler, jwtManager, _, projectMockRepo, logger := SetupNewProjectHandler()

	jwtToken, err := jwtManager.GenerateAccessToken(1234, "test-user", "gudi@gmail.com")

	require.NoError(t, err)
	jwtMiddleware := middleware.JWTAuthMiddleware(logger, jwtManager)

	testCases := []struct {
		testName            string
		projectId           any
		mockSetup           func()
		expectedStatusCode  int
		expectedServiceCall bool
	}{
		{
			testName:  "Valid Project Id",
			projectId: 234,
			mockSetup: func() {
				projectMockRepo.On("GetProjectById", mock.Anything, int64(234)).Return(
					projectdb.Project{
						ID:   234,
						Name: "valid project",
					}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedServiceCall: true,
		},
		{
			testName:  "InValid Project Id",
			projectId: "agsgs",
			mockSetup: func() {

			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedServiceCall: false,
		},
		{
			testName:  "NonExistent Project Id",
			projectId: 46,
			mockSetup: func() {
				projectMockRepo.On("GetProjectById", mock.Anything, int64(46)).Return(
					projectdb.Project{},
					customErrors.ErrProjectIDNotExist)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedServiceCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			projectMockRepo.Calls = nil
			projectMockRepo.ExpectedCalls = nil

			tc.mockSetup()

			// Create request
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/projects/%v", tc.projectId), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			// Recorder
			w := httptest.NewRecorder()

			// Set up full group routing like production
			r := gin.Default()
			apiGroup := r.Group("/api/v1")
			userGroup := apiGroup.Group("/projects")

			userGroup.GET("/:id", jwtMiddleware, projectHandler.GetProjectById)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)

			if tc.expectedServiceCall {
				projectMockRepo.AssertCalled(t, "GetProjectById", mock.Anything, mock.Anything)

			} else {
				projectMockRepo.AssertNotCalled(t, "GetProjectById", mock.Anything, mock.Anything)

			}

		})
	}

}

func TestGetProjectByNameHandler(t *testing.T) {

	projectHandler, jwtManager, _, projectMockRepo, logger := SetupNewProjectHandler()

	jwtToken, err := jwtManager.GenerateAccessToken(1234, "test-user", "gudi@gmail.com")

	require.NoError(t, err)
	jwtMiddleware := middleware.JWTAuthMiddleware(logger, jwtManager)

	testCases := []struct {
		testName            string
		projectname         string
		mockSetup           func()
		expectedStatusCode  int
		expectedServiceCall bool
	}{
		{
			testName:    "valid project name",
			projectname: "project-1",
			mockSetup: func() {
				params := projectdb.GetProjectByNameParams{
					Name:   "project-1",
					UserID: int32(1234),
				}
				projectMockRepo.On("GetProjectByName", mock.Anything, params).Return(
					projectdb.Project{
						Name:   "project-1",
						UserID: 1234,
					}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedServiceCall: true,
		},
		{
			testName:    "Non existent project",
			projectname: "project-not-found",
			mockSetup: func() {
				params := projectdb.GetProjectByNameParams{
					Name:   "project-not-found",
					UserID: int32(1234),
				}
				projectMockRepo.On("GetProjectByName", mock.Anything, params).Return(
					projectdb.Project{}, sql.ErrNoRows)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedServiceCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			projectMockRepo.Calls = nil
			projectMockRepo.ExpectedCalls = nil

			tc.mockSetup()

			// Create request
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/projects/names?name=%s", tc.projectname), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			// Recorder
			w := httptest.NewRecorder()

			// Set up full group routing like production
			r := gin.Default()
			apiGroup := r.Group("/api/v1")
			userGroup := apiGroup.Group("/projects")

			userGroup.GET("/names", jwtMiddleware, projectHandler.GetProjectByName)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
			if tc.expectedServiceCall {
				projectMockRepo.AssertCalled(t, "GetProjectByName", mock.Anything, mock.Anything)
			} else {
				projectMockRepo.AssertNotCalled(t, "GetProjectByName", mock.Anything, mock.Anything)
			}

		})
	}
}

func TestGetProjectsByUserIdHandler(t *testing.T) {
	projectHandler, jwtManager, _, projectMockRepo, logger := SetupNewProjectHandler()

	jwtToken, err := jwtManager.GenerateAccessToken(1234, "test-user", "gudi@gmail.com")

	require.NoError(t, err)
	jwtMiddleware := middleware.JWTAuthMiddleware(logger, jwtManager)

	testCases := []struct {
		testName            string
		mockSetup           func()
		expectedStatusCode  int
		expectedServiceCall bool
	}{
		{
			testName: "user containing projects",
			mockSetup: func() {
				projectMockRepo.On("GetProjectsByUserId", mock.Anything, int32(1234)).Return(
					[]projectdb.Project{
						{
							ID:   123,
							Name: "project-1",
						},
						{
							ID:   23,
							Name: "project-2",
						},
					}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedServiceCall: true,
		},
		{
			testName: "user doesnt contain any projects",
			mockSetup: func() {
				projectMockRepo.On("GetProjectsByUserId", mock.Anything, int32(1234)).Return([]projectdb.Project{}, sql.ErrNoRows)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedServiceCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			projectMockRepo.Calls = nil
			projectMockRepo.ExpectedCalls = nil

			tc.mockSetup()

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/projects/", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			// Recorder
			w := httptest.NewRecorder()

			// Set up full group routing like production
			r := gin.Default()
			apiGroup := r.Group("/api/v1")
			userGroup := apiGroup.Group("/projects")

			userGroup.GET("/", jwtMiddleware, projectHandler.GetProjectsByUserId)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
			if tc.expectedServiceCall {
				projectMockRepo.AssertCalled(t, "GetProjectsByUserId", mock.Anything, mock.Anything)
			} else {
				projectMockRepo.AssertNotCalled(t, "GetProjectsByUserId", mock.Anything, mock.Anything)
			}

		})
	}

}

func TestDeleteProjectHandler(t *testing.T) {
	projectHandler, jwtManager, _, projectMockRepo, logger := SetupNewProjectHandler()

	jwtToken, err := jwtManager.GenerateAccessToken(1234, "test-user", "gudi@gmail.com")

	require.NoError(t, err)
	jwtMiddleware := middleware.JWTAuthMiddleware(logger, jwtManager)
	testCases := []struct {
		testName            string
		projectId           int
		mockSetup           func()
		expectedServiceCall bool
		expectedStatusCode  int
	}{
		{
			testName:  "Valid Projectid",
			projectId: 234,
			mockSetup: func() {
				projectMockRepo.On("DeleteProject", mock.Anything, int64(234)).Return(nil)
			},
			expectedServiceCall: true,
			expectedStatusCode:  http.StatusOK,
		},
		{
			testName:  "NonExistent Projectid",
			projectId: 23,
			mockSetup: func() {
				projectMockRepo.On("DeleteProject", mock.Anything, int64(23)).Return(sql.ErrNoRows)
			},
			expectedServiceCall: true,
			expectedStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			gin.SetMode(gin.TestMode)
			projectMockRepo.Calls = nil
			projectMockRepo.ExpectedCalls = nil
			tc.mockSetup()

			// Create request
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/projects/%v", tc.projectId), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			// Recorder
			w := httptest.NewRecorder()

			// Set up full group routing like production
			r := gin.Default()
			apiGroup := r.Group("/api/v1")
			userGroup := apiGroup.Group("/projects")

			userGroup.DELETE("/:id", jwtMiddleware, projectHandler.DeleteProject)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)

			if tc.expectedServiceCall {
				projectMockRepo.AssertCalled(t, "DeleteProject", mock.Anything, mock.Anything)

			} else {
				projectMockRepo.AssertNotCalled(t, "DeleteProject", mock.Anything, mock.Anything)

			}

		})

	}
}



func TestGetTasksByProjectIDHandler(t *testing.T) {
	projectHandler, jwtManager, taskMockRepo, _, logger := SetupNewProjectHandler()

	jwtToken, err := jwtManager.GenerateAccessToken(1234, "test-user", "gudi@gmail.com")

	require.NoError(t, err)
	jwtMiddleware := middleware.JWTAuthMiddleware(logger, jwtManager)
	testCases := []struct {
		testName            string
		projectId           any
		mockSetup           func()
		expectedServiceCall bool
		expectedStatusCode  int
	}{
		{
			testName:  "valid project id ",
			projectId: 24,
			mockSetup: func() {
				taskMockRepo.On("GetTasksByProjectId", mock.Anything, int64(24)).Return(
					[]taskdb.Task{
						{
							Title: "task-1",
						},
						{
							Title: "task-2",
						},
					}, nil)
			},
			expectedServiceCall: true,
			expectedStatusCode:  http.StatusOK,
		},
		{
			testName:  "valid project id ",
			projectId: "try",
			mockSetup: func() {
				
			},
			expectedServiceCall: false,
			expectedStatusCode:  http.StatusBadRequest,
		},
		{
			testName:  "no tasks under request project",
			projectId: 24,
			mockSetup: func() {
				taskMockRepo.On("GetTasksByProjectId", mock.Anything, int64(24)).Return(
					[]taskdb.Task{}, 
					sql.ErrNoRows)
			},
			expectedServiceCall: true,
			expectedStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			taskMockRepo.Calls = nil
			taskMockRepo.ExpectedCalls = nil
			tc.mockSetup()

			// Create request
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/projects/%v/tasks", tc.projectId), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			// Recorder
			w := httptest.NewRecorder()

			// Set up full group routing like production
			r := gin.Default()
			apiGroup := r.Group("/api/v1")
			userGroup := apiGroup.Group("/projects")

			userGroup.GET("/:id/tasks", jwtMiddleware, projectHandler.GetTasksByProjectID)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)

			if tc.expectedServiceCall {
				taskMockRepo.AssertCalled(t, "GetTasksByProjectId", mock.Anything, mock.Anything)

			} else {
				taskMockRepo.AssertNotCalled(t, "GetTasksByProjectId", mock.Anything, mock.Anything)

			}
		})

	}
}
