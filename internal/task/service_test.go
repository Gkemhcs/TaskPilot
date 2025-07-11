package task

import (
	"context"
	"database/sql"
	"errors"

	"os"
	"testing"
	"time"

	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	taskdb "github.com/Gkemhcs/taskpilot/internal/task/gen"
	"github.com/lib/pq"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var taskService *TaskService
var mockRepo *MockTaskRepo

func TestMain(m *testing.M) {
	mockRepo = new(MockTaskRepo)
	taskService = NewTaskService(mockRepo)
	os.Exit(m.Run())
}

// Simulated duplicate error
func mockDuplicateError() error {
	return &pq.Error{Code: "23505"}
}
func TestGetStatus(t *testing.T) {
	testCases := []struct {
		testName       string
		status         string
		expectedOutput taskdb.TaskStatus
	}{
		{
			testName:       "valid in_progress",
			status:         "in_progress",
			expectedOutput: taskdb.TaskStatusINPROGRESS,
		},
		{
			testName:       "valid complete task",
			status:         "done",
			expectedOutput: taskdb.TaskStatusDONE,
		},
		{
			testName:       "invalid task status",
			status:         "working_fine",
			expectedOutput: taskdb.TaskStatusTODO,
		},
		{
			testName:       "invalid task status",
			status:         "working_fine",
			expectedOutput: taskdb.TaskStatusTODO,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			output := getStatus(tc.status)
			assert.Equal(t, output, tc.expectedOutput)
		})
	}
}

func TestGetPriority(t *testing.T) {
	testCases := []struct {
		testName       string
		priority       string
		expectedOutput taskdb.TaskPriority
	}{
		{
			testName:       "critical priority",
			priority:       "critical",
			expectedOutput: taskdb.TaskPriorityCRITICAL,
		},
		{
			testName:       "high priority",
			priority:       "high",
			expectedOutput: taskdb.TaskPriorityHIGH,
		},
		{
			testName:       "low priority",
			priority:       "low",
			expectedOutput: taskdb.TaskPriorityLOW,
		},
		{
			testName:       "invalid priority defaults to medium",
			priority:       "warning",
			expectedOutput: taskdb.TaskPriorityMEDIUM,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			actualPriority := getPriority(tc.priority)
			assert.Equal(t, tc.expectedOutput, actualPriority)
		})
	}
}

func TestCreateTask(t *testing.T) {

	testCases := []struct {
		testName         string
		project          CreateTaskInput
		expectedParams   taskdb.CreateTaskParams
		expectedError    error
		expectedOutput   *taskdb.Task
		returnedError    error
		isProjectValid   bool
		expectedStatus   taskdb.TaskStatus
		expectedPriority taskdb.TaskPriority
	}{{
		testName: "valid task with high priority",
		project: CreateTaskInput{
			ProjectID:   1234,
			Title:       "Valid task with high priority and todo status",
			AssigneeID:  102,
			Description: "valid task1",
			Status:      "in_progress",
			Priority:    "high",
			DueDate:     time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
		},
		expectedParams: taskdb.CreateTaskParams{
			ProjectID: 1234,
			Title:     "Valid task with high priority and todo status",
			AssigneeID: sql.NullInt64{
				Int64: 102,
				Valid: true,
			},
			Description: "valid task1",
			Status:      taskdb.TaskStatusINPROGRESS,
			Priority:    taskdb.TaskPriorityHIGH,
			DueDate: sql.NullTime{
				Time:  time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
				Valid: true,
			},
		},
		expectedError: nil,
		returnedError: nil,
		expectedOutput: &taskdb.Task{
			ID:        1234,
			ProjectID: 1234,
			AssigneeID: sql.NullInt64{
				Int64: 102,
				Valid: true,
			},
			Title:       "Valid task with high priority and todo status",
			Description: "valid task1",
			Status:      taskdb.TaskStatusINPROGRESS,
			Priority:    taskdb.TaskPriorityHIGH,
			DueDate: sql.NullTime{
				Time:  time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
				Valid: true,
			},
		},
		isProjectValid:   true,
		expectedStatus:   taskdb.TaskStatusINPROGRESS,
		expectedPriority: taskdb.TaskPriorityHIGH,
	},
		{
			testName: "task with wrong values for status and priority",
			project: CreateTaskInput{
				ProjectID:   1234,
				Title:       "task with wrong values for status and priority",
				AssigneeID:  102,
				Description: "invalid task 1",
				Status:      "todoer",
				Priority:    "lower",
				DueDate:     time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
			},
			expectedParams: taskdb.CreateTaskParams{
				ProjectID: 1234,
				Title:     "task with wrong values for status and priority",
				AssigneeID: sql.NullInt64{
					Int64: 102,
					Valid: true,
				},
				Description: "invalid task 1",
				Status:      taskdb.TaskStatusTODO,
				Priority:    taskdb.TaskPriorityMEDIUM,
				DueDate: sql.NullTime{
					Time:  time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
					Valid: true,
				},
			},
			expectedError: nil,
			returnedError: nil,
			expectedOutput: &taskdb.Task{
				ID:        1234,
				ProjectID: 1234,
				AssigneeID: sql.NullInt64{
					Int64: 102,
					Valid: true,
				},
				Title:       "task with wrong values for status and priority",
				Description: "invalid task1",
				Status:      taskdb.TaskStatusTODO,
				Priority:    taskdb.TaskPriorityMEDIUM,
				DueDate: sql.NullTime{
					Time:  time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
					Valid: true,
				},
			},
			expectedStatus:   taskdb.TaskStatusTODO,
			expectedPriority: taskdb.TaskPriorityMEDIUM,
		},
		{
			testName: "duplicate task",
			project: CreateTaskInput{
				ProjectID:   1234,
				Title:       "Some duplicate task",
				AssigneeID:  101,
				Description: "already exists",
				Status:      "todo",
				Priority:    "medium",
				DueDate:     time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC),
			},
			expectedParams: taskdb.CreateTaskParams{
				ProjectID:   1234,
				Title:       "Some duplicate task",
				AssigneeID:  sql.NullInt64{Int64: 101, Valid: true},
				Description: "already exists",
				Status:      taskdb.TaskStatusTODO,
				Priority:    taskdb.TaskPriorityMEDIUM,
				DueDate:     sql.NullTime{Time: time.Date(2025, 7, 11, 14, 30, 0, 0, time.UTC), Valid: true},
			},
			expectedError:  customErrors.ErrTaskAlreadyExists, // this is your final expected output
			returnedError:  mockDuplicateError(),              // this is what the DB would return
			expectedOutput: &taskdb.Task{},                    // likely ignored if error happens
		},
		
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepo.Calls = nil
			mockRepo.ExpectedCalls = nil
			mockRepo.On("CreateTask", mock.Anything, tc.expectedParams).Return(*tc.expectedOutput, tc.returnedError)
			ctx := context.TODO()
			result, err := taskService.CreateTask(ctx, tc.project)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)

			} else {
				assert.NoError(t, err)
				assert.Equal(t, result.ProjectID, tc.expectedOutput.ProjectID)
				assert.Equal(t, result.Priority, tc.expectedPriority)
				assert.Equal(t, result.Status, tc.expectedStatus)
			}

			mockRepo.AssertCalled(t, "CreateTask", ctx, tc.expectedParams)
		})
	}

}


func TestGetTaskByID(t *testing.T){

	testCases:=[]struct{
			testName string 
			taskID int 
			expectedTaskID int64
			returnOutput taskdb.Task
			expectedOutput *taskdb.Task
			expectedError error 
			returnError error 
			
	}{
		{
			testName: "valid task id",
			taskID: 101,
			expectedTaskID: 101,
			returnOutput: taskdb.Task{
				ID: 101,
				Title: "valid task",
			},
			expectedOutput: &taskdb.Task{
				ID: 101,
				Title: "valid task",
			},
			expectedError: nil,
			returnError: nil,
		},
		{
			testName: "invalid task id",
			taskID: 1001,
			expectedTaskID: 1001,
			returnOutput: taskdb.Task{},
			expectedOutput: nil ,
			expectedError: customErrors.ErrTaskNotFound,
			returnError: sql.ErrNoRows,
		},
		{
			testName: "db error",
			taskID: 1001,
			expectedTaskID: 1001,
			returnOutput: taskdb.Task{},
			expectedOutput: nil ,
			expectedError: errors.New("db connection failed"),
			returnError: errors.New("db connection failed"),
		},
		}

		for _,tc := range testCases{
			t.Run(tc.testName,func(t *testing.T){
				mockRepo.Calls=nil 
				mockRepo.ExpectedCalls=nil 
				mockRepo.On("GetTaskById",mock.Anything,tc.expectedTaskID).Return(tc.returnOutput,tc.returnError)


				ctx:=context.TODO()
				result,err:=taskService.GetTaskByID(ctx,tc.taskID)
				if tc.returnError!=nil{
					assert.Equal(t,err,tc.expectedError)
				}else{
					assert.Equal(t,tc.expectedOutput.Title,result.Title)
				}
				mockRepo.AssertCalled(t,"GetTaskById",mock.Anything,tc.expectedTaskID)

			})
		}
}

func TestGetTasksByProjectId(t *testing.T){

	testCases:=[]struct{
		testName string 
		projectID int 
		expectedProjectID int64
		returnResult []taskdb.Task
		expectedResult []taskdb.Task
		returnError error 
		expectedError error 
	}{
		{
			testName: "Valid Project Id",
			projectID: 24,
			expectedProjectID: 24,
			returnResult: []taskdb.Task{
				{
					Title: "task1",
				},{
					Title:"task2",
				},

			},
			expectedResult: []taskdb.Task{
				{
					Title: "task1",
				},{
					Title:"task2",
				},

			},
			returnError: nil,
			expectedError: nil,
		},
		{
			testName: "Empty Tasks",
			projectID: 200,
			expectedProjectID: 200,
			returnResult: nil,
			expectedResult: nil,
			expectedError: customErrors.ErrTasksAreEmpty,
			returnError: sql.ErrNoRows,
		},
		{
			testName: "DB Error",
			projectID: 200,
			expectedProjectID: 200,
			returnResult: nil,
			expectedResult: nil,
			expectedError:errors.New("db error"),
			returnError: errors.New("db error"),
		},


	}

	for _,tc := range testCases{
		t.Run(tc.testName,func(t *testing.T){
			mockRepo.Calls=nil 
			mockRepo.ExpectedCalls=nil 
			mockRepo.On("GetTasksByProjectId",mock.Anything,tc.expectedProjectID).Return(tc.returnResult,tc.returnError)

			result,err:=taskService.GetTasksByProjectID(context.TODO(),tc.projectID)
			if tc.expectedError!=nil{
				assert.Equal(t,err,tc.expectedError)
			}else{
				assert.NoError(t,err)
				assert.Equal(t,result,tc.expectedResult)
			}
			mockRepo.AssertCalled(t,"GetTasksByProjectId",mock.Anything,tc.expectedProjectID)
		})
	}
}


func TestDeleteTask(t *testing.T){

	testCases:=[]struct{
		testName string 
		taskID int 
		expectedTaskID int64 
		expectedErr error 
		returnError error 
		returnResult int64

	}{
			{
				testName:"valid task id",
				taskID: 103,
				expectedTaskID: 103,
				returnResult: 1,
				expectedErr: nil,
				returnError: nil,
			},
			{
				testName:"invalid task id",
				taskID: 103,
				expectedTaskID: 103,
				returnResult: 0,
				expectedErr: customErrors.ErrTaskNotFound,
				returnError: nil,
			},
	}


	for _,tc := range testCases{
		t.Run(tc.testName,func(t *testing.T){
			mockRepo.Calls=nil 
			mockRepo.ExpectedCalls=nil 
			mockRepo.On("DeleteTask",mock.Anything,tc.expectedTaskID).Return(tc.returnResult,tc.returnError)
			err:=taskService.DeleteTask(context.TODO(),tc.taskID)
			if tc.expectedErr!=nil{
				assert.Equal(t,tc.expectedErr,err)
			}else{
				assert.NoError(t,err)
			}
			mockRepo.AssertCalled(t,"DeleteTask",mock.Anything,tc.expectedTaskID)
		})
	}

}