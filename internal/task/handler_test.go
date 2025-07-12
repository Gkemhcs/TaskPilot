package task

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
	"github.com/Gkemhcs/taskpilot/internal/project"
	projectdb "github.com/Gkemhcs/taskpilot/internal/project/gen"
	taskdb "github.com/Gkemhcs/taskpilot/internal/task/gen"
	"github.com/Gkemhcs/taskpilot/internal/user"
	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func SetupNewTaskHandler() (*TaskHandler, *MockTaskRepo, *project.MockProjectRepo, *user.MockUserRepo, *auth.JWTManager, *logrus.Logger) {
	//mocking task repository
	taskMockRepo := new(MockTaskRepo)

	//mocking project repository
	projectMockRepo := new(project.MockProjectRepo)

	//mocking user repository
	userMockRepo := new(user.MockUserRepo)

	params := auth.CreateJwtManagerParams{
		AccessTokenDuration:  10 * time.Minute,
		RefreshTokenDuration: 10 * time.Hour,
		AccessTokenKey:       "rnk3mkrk3rk3rk3",
		RefreshTokenKey:      "21ieh12iei21eji12e",
	}
	jwtManager := auth.NewJWTManager(params)

	//initialising task service
	taskService := NewTaskService(taskMockRepo)

	//initialising taskproject service
	projectService := project.NewProjectService(projectMockRepo)

	//user service
	userService := user.NewUserService(userMockRepo)

	// logger setup and disable logging to console/io

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	taskHandler := NewTaskHandler(*taskService, userService, logger, projectService)
	return taskHandler, taskMockRepo, projectMockRepo, userMockRepo, jwtManager, logger

}

func TestCreateTaskHandler(t *testing.T) {
	taskHandler, taskMockRepo, projectMockRepo, userMockRepo, jwtManager, logger := SetupNewTaskHandler()

	jwtToken, err := jwtManager.GenerateAccessToken(1234, "test-user", "gudi@gmail")
	require.NoError(t, err)
	jwtMiddleware := middleware.JWTAuthMiddleware(logger, jwtManager)

	testCases := []struct {
		testName                   string
		requestBody                map[string]any
		mockUserRepoSetup          func()
		mockProjectRepo            func()
		mockTaskRepo               func()
		expectedStatusCode         int
		expectedUserServiceCall    bool
		expectedProjectServiceCall bool
		expectedTaskServiceCall    bool
	}{
		{
			testName: "valid task with missing status and priority",
			requestBody: map[string]any{
				"project_id":     4123,
				"title":          "Adding Navbar",
				"assignee_email": "gudi@gmail",
				"description":    "to add  a navbar on right of home page",
				"due_date":       time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
			},
			mockUserRepoSetup: func() {
				userMockRepo.On("GetUserByEmail", mock.Anything, "gudi@gmail").Return(
					userdb.User{
						ID:    1234,
						Email: "gudi@gmail",
						Name:  "test-user",
					}, nil)
			},
			mockProjectRepo: func() {
				projectMockRepo.On("GetProjectById", mock.Anything, int64(4123)).Return(
					projectdb.Project{
						Name: "project-1",
					}, nil)
			},
			mockTaskRepo: func() {
				params := taskdb.CreateTaskParams{
					ProjectID: 4123,
					AssigneeID: sql.NullInt64{
						Int64: 1234,
						Valid: true,
					},
					Title:       "Adding Navbar",
					Description: "to add  a navbar on right of home page",
					Status:      taskdb.TaskStatusTODO,
					Priority:    taskdb.TaskPriorityMEDIUM,
					DueDate: sql.NullTime{
						Time:  time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
						Valid: true,
					},
				}
				taskMockRepo.On("CreateTask", mock.Anything, params).Return(
					taskdb.Task{
						Title:     "Adding Navbar",
						ProjectID: 4123,
						AssigneeID: sql.NullInt64{
							Int64: 1234,
							Valid: true,
						},
					}, nil)
			},
			expectedStatusCode:         http.StatusCreated,
			expectedUserServiceCall:    true,
			expectedProjectServiceCall: true,
			expectedTaskServiceCall:    true,
		},
		{
			testName: "valid task",
			requestBody: map[string]any{
				"project_id":     4123,
				"title":          "Adding Navbar",
				"assignee_email": "gudi@gmail",
				"description":    "to add  a navbar on right of home page",
				"status":         "todo",
				"priority":       "high",
				"due_date":       time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
			},
			mockUserRepoSetup: func() {
				userMockRepo.On("GetUserByEmail", mock.Anything, "gudi@gmail").Return(
					userdb.User{
						ID:    1234,
						Email: "gudi@gmail",
						Name:  "test-user",
					}, nil)
			},
			mockProjectRepo: func() {
				projectMockRepo.On("GetProjectById", mock.Anything, int64(4123)).Return(
					projectdb.Project{
						Name: "project-1",
					}, nil)
			},
			mockTaskRepo: func() {
				params := taskdb.CreateTaskParams{
					ProjectID: 4123,
					AssigneeID: sql.NullInt64{
						Int64: 1234,
						Valid: true,
					},
					Title:       "Adding Navbar",
					Description: "to add  a navbar on right of home page",
					Status:      taskdb.TaskStatusTODO,
					Priority:    taskdb.TaskPriorityHIGH,
					DueDate: sql.NullTime{
						Time:  time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
						Valid: true,
					},
				}
				taskMockRepo.On("CreateTask", mock.Anything, params).Return(
					taskdb.Task{
						Title:     "Adding Navbar",
						ProjectID: 4123,
						AssigneeID: sql.NullInt64{
							Int64: 1234,
							Valid: true,
						},
					}, nil)
			},
			expectedStatusCode:         http.StatusCreated,
			expectedUserServiceCall:    true,
			expectedProjectServiceCall: true,
			expectedTaskServiceCall:    true,
		},
		{
			testName: "duplicate task",
			requestBody: map[string]any{
				"project_id":     4123,
				"title":          "Adding Navbar",
				"assignee_email": "gudi@gmail",
				"description":    "to add  a navbar on right of home page",
				"status":         "todo",
				"priority":       "high",
				"due_date":       time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
			},
			mockUserRepoSetup: func() {
				userMockRepo.On("GetUserByEmail", mock.Anything, "gudi@gmail").Return(
					userdb.User{
						ID:    1234,
						Email: "gudi@gmail",
						Name:  "test-user",
					}, nil)
			},
			mockProjectRepo: func() {
				projectMockRepo.On("GetProjectById", mock.Anything, int64(4123)).Return(
					projectdb.Project{
						Name: "project-1",
					}, nil)
			},
			mockTaskRepo: func() {
				params := taskdb.CreateTaskParams{
					ProjectID: 4123,
					AssigneeID: sql.NullInt64{
						Int64: 1234,
						Valid: true,
					},
					Title:       "Adding Navbar",
					Description: "to add  a navbar on right of home page",
					Status:      taskdb.TaskStatusTODO,
					Priority:    taskdb.TaskPriorityHIGH,
					DueDate: sql.NullTime{
						Time:  time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
						Valid: true,
					},
				}
				taskMockRepo.On("CreateTask", mock.Anything, params).Return(
					taskdb.Task{}, mockDuplicateError())
			},
			expectedStatusCode:         http.StatusBadRequest,
			expectedUserServiceCall:    true,
			expectedProjectServiceCall: true,
			expectedTaskServiceCall:    true,
		},

		{
			testName: "missing assignee email",
			requestBody: map[string]any{
				"project_id":  4123,
				"title":       "Adding Navbar",
				"description": "to add  a navbar on right of home page",
				"status":      "todo",
				"priority":    "high",
				"due_date":    time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
			},
			mockUserRepoSetup:          func() {},
			mockProjectRepo:            func() {},
			mockTaskRepo:               func() {},
			expectedStatusCode:         http.StatusBadRequest,
			expectedUserServiceCall:    false,
			expectedProjectServiceCall: false,
			expectedTaskServiceCall:    false,
		},
		{
			testName: "missing task due date",
			requestBody: map[string]any{
				"project_id":     4123,
				"title":          "Adding Navbar",
				"assignee_email": "gudi@gmail",
				"description":    "to add  a navbar on right of home page",
				"status":         "todo",
				"priority":       "high",
			},
			mockUserRepoSetup:          func() {},
			mockProjectRepo:            func() {},
			mockTaskRepo:               func() {},
			expectedStatusCode:         http.StatusBadRequest,
			expectedUserServiceCall:    false,
			expectedProjectServiceCall: false,
			expectedTaskServiceCall:    false,
		},
		{
			testName: "missing task title",
			requestBody: map[string]any{
				"project_id":     4123,
				"assignee_email": "gudi@gmail",
				"description":    "to add  a navbar on right of home page",
				"status":         "todo",
				"priority":       "high",
				"due_date":       time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
			},
			mockUserRepoSetup:          func() {},
			mockProjectRepo:            func() {},
			mockTaskRepo:               func() {},
			expectedStatusCode:         http.StatusBadRequest,
			expectedUserServiceCall:    false,
			expectedProjectServiceCall: false,
			expectedTaskServiceCall:    false,
		},
		{
			testName: "missing project id",
			requestBody: map[string]any{
				"title":          "Adding Navbar",
				"assignee_email": "gudi@gmail",
				"description":    "to add  a navbar on right of home page",
				"status":         "todo",
				"priority":       "high",
				"due_date":       time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
			},
			mockUserRepoSetup:          func() {},
			mockProjectRepo:            func() {},
			mockTaskRepo:               func() {},
			expectedStatusCode:         http.StatusBadRequest,
			expectedUserServiceCall:    false,
			expectedProjectServiceCall: false,
			expectedTaskServiceCall:    false,
		},
		{
			testName: "invalid project id",
			requestBody: map[string]any{
				"project_id":     41233,
				"title":          "Adding Navbar",
				"assignee_email": "gudi@gmail",
				"description":    "to add  a navbar on right of home page",
				"status":         "todo",
				"priority":       "high",
				"due_date":       time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
			},
			mockUserRepoSetup: func() {},
			mockProjectRepo: func() {
				projectMockRepo.On("GetProjectById", mock.Anything, int64(41233)).Return(
					projectdb.Project{},
					customErrors.ErrProjectIDNotExist)
			},
			mockTaskRepo:               func() {},
			expectedStatusCode:         http.StatusBadRequest,
			expectedUserServiceCall:    false,
			expectedProjectServiceCall: true,
			expectedTaskServiceCall:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			taskMockRepo.Calls = nil
			taskMockRepo.ExpectedCalls = nil

			userMockRepo.Calls = nil
			userMockRepo.ExpectedCalls = nil

			projectMockRepo.Calls = nil
			projectMockRepo.ExpectedCalls = nil

			tc.mockProjectRepo()
			tc.mockTaskRepo()
			tc.mockUserRepoSetup()

			body, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			// Create request
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/tasks/", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			// Recorder
			w := httptest.NewRecorder()

			// Set up full group routing like production
			r := gin.Default()
			apiGroup := r.Group("/api/v1")
			userGroup := apiGroup.Group("/tasks")

			userGroup.POST("/", jwtMiddleware, taskHandler.CreateTask)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, tc.expectedStatusCode)

			if tc.expectedProjectServiceCall {
				projectMockRepo.AssertCalled(t, "GetProjectById", mock.Anything, mock.Anything)
			} else {

			}

			if tc.expectedTaskServiceCall {
				taskMockRepo.AssertCalled(t, "CreateTask", mock.Anything, mock.Anything)
			} else {
				taskMockRepo.AssertNotCalled(t, "CreateTask", mock.Anything, mock.Anything)
			}
			if tc.expectedUserServiceCall {
				userMockRepo.AssertCalled(t, "GetUserByEmail", mock.Anything, mock.Anything)
			} else {
				userMockRepo.AssertNotCalled(t, "GetUserByEmail", mock.Anything, mock.Anything)
			}

		})
	}

}



func TestGetTaskByIDHandler(t *testing.T){

	taskHandler, taskMockRepo, _, _, jwtManager, logger := SetupNewTaskHandler()

	jwtToken, err := jwtManager.GenerateAccessToken(1234, "test-user", "gudi@gmail")
	require.NoError(t, err)
	jwtMiddleware := middleware.JWTAuthMiddleware(logger, jwtManager)
	
	testCases:=[]struct{
		testName string 
		taskID int 
		mockSetup func()
		expectedStatusCode int 
		expectedServiceCall bool 
	}{
		{
			testName: "Valid task id",
			taskID:101,
			mockSetup: func() {
				taskMockRepo.On("GetTaskById",mock.Anything,int64(101)).Return(
					taskdb.Task{
						Title: "task-1",
						ID:101,},nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedServiceCall: true,
		},
		{
			testName: "Non Existent Task ID",
			taskID:121,
			mockSetup: func() {
				taskMockRepo.On("GetTaskById",mock.Anything,int64(121)).Return(
					taskdb.Task{},sql.ErrNoRows)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedServiceCall: true,
		},
	}


	for _,tc := range testCases{
			t.Run(tc.testName,func(t *testing.T){

				gin.SetMode(gin.TestMode)
				taskMockRepo.Calls=nil 
				taskMockRepo.ExpectedCalls=nil 


				tc.mockSetup()

			// Create request
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tasks/%v", tc.taskID), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			// Recorder
			w := httptest.NewRecorder()

			// Set up full group routing like production
			r := gin.Default()
			apiGroup := r.Group("/api/v1")
			userGroup := apiGroup.Group("/tasks")

			userGroup.GET("/:id", jwtMiddleware, taskHandler.GetTaskByID)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)

			if tc.expectedServiceCall {
				taskMockRepo.AssertCalled(t, "GetTaskById", mock.Anything, mock.Anything)

			} else {
				taskMockRepo.AssertNotCalled(t, "GetTaskById", mock.Anything, mock.Anything)

			}



			})

	}



}



func DeleteTaskHandler(t *testing.T){


	taskHandler, taskMockRepo, _, _, jwtManager, logger := SetupNewTaskHandler()

	jwtToken, err := jwtManager.GenerateAccessToken(1234, "test-user", "gudi@gmail.com")

	require.NoError(t, err)
	jwtMiddleware := middleware.JWTAuthMiddleware(logger, jwtManager)

	testCases:=[]struct{
		testName string 
		taskID int 
		mockSetup func()
		expectedStatusCode int 
		expectedServiceCall bool 

	}{	
		{
			testName: "existing task",
			taskID: 102,
			mockSetup: func() {
				taskMockRepo.On("DeleteTask",mock.Anything,int64(102)).Return(1,nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedServiceCall: true,
		},
		{
			testName: "non-existing task",
			taskID: 104,
			mockSetup: func() {
				taskMockRepo.On("DeleteTask",mock.Anything,int64(104)).Return(0,nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedServiceCall: true,
		},
	}

	for _,tc := range testCases{
		t.Run(tc.testName,func( t *testing.T){

			gin.SetMode(gin.TestMode)
			taskMockRepo.Calls=nil 
			taskMockRepo.ExpectedCalls=nil 


			tc.mockSetup()

			// Create request
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/tasks/%v", tc.taskID), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			// Recorder
			w := httptest.NewRecorder()

			// Set up full group routing like production
			r := gin.Default()
			apiGroup := r.Group("/api/v1")
			userGroup := apiGroup.Group("/tasks")

			userGroup.DELETE("/:id", jwtMiddleware, taskHandler.DeleteTask)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)

			if tc.expectedServiceCall {
				taskMockRepo.AssertCalled(t, "DeleteTask", mock.Anything, mock.Anything)

			} else {
				taskMockRepo.AssertNotCalled(t, "DeleteTask", mock.Anything, mock.Anything)

			}

		})
	}

}


