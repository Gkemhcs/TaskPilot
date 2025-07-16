package project

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	projectdb "github.com/Gkemhcs/taskpilot/internal/project/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
)


var projectService  *ProjectService 
var mockRepo *MockProjectRepo
func TestMain(m *testing.M){
	mockRepo=new(MockProjectRepo)
	projectService=NewProjectService(mockRepo)
	os.Exit(m.Run())

}
func TestCreateProject(t *testing.T) {
	testCases := []struct {
		name            string
		ctx             context.Context
		input           Project
		expectedParams  projectdb.CreateProjectParams
		expectedProject projectdb.Project
		expectError     bool
	}{
		{
			name: "Valid green project",
			ctx:  context.TODO(),
			input: Project{
				Name:        "Green Project",
				Description: "Green is good",
				Color:       "green",
				User:        101,
			},
			expectedParams: projectdb.CreateProjectParams{
				UserID: 101,
				Name:   "Green Project",
				Description: sql.NullString{
					String: "Green is good",
					Valid:  true,
				},
				Color: projectdb.NullProjectColor{
					ProjectColor: mapColor("green"),
					Valid:        true,
				},
			},
			expectedProject: projectdb.Project{ID: 1, UserID: 101, Name: "Green Project"},
			expectError:     false,
		},
		{
			name: "Valid Yellow project",
			ctx:  context.TODO(),
			input: Project{
				Name:        "Yellow Project",
				Description: "Yellow is good",
				Color:       "yellow",
				User:        101,
			},
			expectedParams: projectdb.CreateProjectParams{
				UserID: 101,
				Name:   "Yellow Project",
				Description: sql.NullString{
					String: "Yellow is good",
					Valid:  true,
				},
				Color: projectdb.NullProjectColor{
					ProjectColor: mapColor("yellow"),
					Valid:        true,
				},
			},
			expectedProject: projectdb.Project{ID: 1, UserID: 101, Name: "Yellow Project"},
			expectError:     false,
		},
		{
			name: "Valid Default project",
			ctx:  context.TODO(),
			input: Project{
				Name:        "Default Project",
				Description: "Default is good",
				Color:       "red",
				User:        101,
			},
			expectedParams: projectdb.CreateProjectParams{
				UserID: 101,
				Name:   "Default Project",
				Description: sql.NullString{
					String: "Default is good",
					Valid:  true,
				},
				Color: projectdb.NullProjectColor{
					ProjectColor: mapColor("red"),
					Valid:        true,
				},
			},
			expectedProject: projectdb.Project{ID: 1, UserID: 101, Name: "Default Project"},
			expectError:     false,
		},
		{
			name: "Empty Description",
			ctx:  context.TODO(),
			input: Project{
				Name:        "Empty Project",
				Description: "",
				Color:       "red",
				User:        101,
			},
			expectedParams: projectdb.CreateProjectParams{
				UserID: 101,
				Name:   "Empty Project",
				Description: sql.NullString{
					String: "",
					Valid:  false,
				},
				Color: projectdb.NullProjectColor{
					ProjectColor: mapColor("red"),
					Valid:        true,
				},
			},
			expectedProject: projectdb.Project{ID: 1, UserID: 101, Name: "Empty Project"},
			expectError:     false,
		},
		{
			name: "Empty Color",
			ctx:  context.TODO(),
			input: Project{
				Name:        "Empty Color",
				Description: "koti",
				Color:        "",
				User:        101,
			},
			expectedParams: projectdb.CreateProjectParams{
				UserID: 101,
				Name:   "Empty Color",
				Description: sql.NullString{
					String: "koti",
					Valid:  true,
				},
				Color: projectdb.NullProjectColor{
					ProjectColor: mapColor(""),
					Valid:        true,
				},
			},
			expectedProject: projectdb.Project{ID: 1, UserID: 101, Name: "Empty Color"},
			expectError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expectError {
				mockRepo.On("CreateProject", mock.Anything, tc.expectedParams).
					Return(tc.expectedProject, nil)
			}else{
				mockRepo.On("CreateProject",mock.Anything,tc.expectedParams).Return(nil,errors.New("error"))
			}
			result, err := projectService.CreateProject(tc.ctx, tc.input)

			if tc.expectError {
				assert.Error(t, err)
				
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedProject.ID, result.ID)
				assert.Equal(t, tc.expectedProject.Name, result.Name)
				
				if tc.input.Description==""{
					assert.Equal(t,tc.expectedParams.Description.Valid,false)
				}

				if tc.input.Color==""{
					expectedColor:= projectdb.NullProjectColor(projectdb.NullProjectColor{ProjectColor:"RED", Valid:true})
					assert.Equal(t,tc.expectedParams.Color,expectedColor)
				}
				
				mockRepo.AssertCalled(t, "CreateProject", mock.Anything, tc.expectedParams)
			}
		})
	}
}

func TestGetProjectById(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	projectService := NewProjectService(mockRepo)

	testCases := []struct {
		name          string
		projectID     int
		mockSetup     func()
		expectedError error
		expectedProj  *projectdb.Project
	}{
		{
			name:      "Project found successfully",
			projectID: 1,
			mockSetup: func() {
				mockRepo.On("GetProjectById", mock.Anything, int64(1)).
					Return(projectdb.Project{
						ID:     1,
						UserID: 101,
						Name:   "Test Project",
					}, nil)
			},
			expectedError: nil,
			expectedProj: &projectdb.Project{
				ID:     1,
				UserID: 101,
				Name:   "Test Project",
			},
		},
		{
			name:      "Project not found",
			projectID: 2,
			mockSetup: func() {
				mockRepo.On("GetProjectById", mock.Anything, int64(2)).
					Return(projectdb.Project{}, sql.ErrNoRows)
			},
			expectedError: customErrors.ErrProjectIDNotExist,
			expectedProj:  nil,
		},
		{
			name:      "Database error",
			projectID: 3,
			mockSetup: func() {
				mockRepo.On("GetProjectById", mock.Anything, int64(3)).
					Return(projectdb.Project{}, errors.New("db error"))
			},
			expectedError: errors.New("db error"),
			expectedProj:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			
			tc.mockSetup()

			project, err := projectService.GetProjectById(context.TODO(), tc.projectID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedProj, project)
			}

			mockRepo.AssertCalled(t, "GetProjectById", mock.Anything, int64(tc.projectID))
		})
	}
}


func TestGetProjectsByUserId(t *testing.T) {
	testCases := []struct {
		name           string
		userID          int
		mockSetup      func()
		expectedError  error
		expectedResult []projectdb.Project
	}{
		{
			name:  "Valid User ID",
			userID: 23,
			mockSetup: func() {
				mockRepo.On("GetProjectsByUserId",context.TODO(), int32(23)).Return(
					[]projectdb.Project{
						{
							ID:     1,
							UserID: 101,
							Name:   "Test Project",
						},
					}, nil,
				)
			},
			expectedError:  nil,
			expectedResult: []projectdb.Project{
				{
					ID:     1,
					UserID: 101,
					Name:   "Test Project",
				},
			},
		},
		{
			name:"Non Existent User Id",
			userID:3899,
			mockSetup: func() {
				mockRepo.On("GetProjectsByUserId",context.TODO(),int32(3899)).Return(
					[]projectdb.Project{},customErrors.ErrUserNotExist,
				)

			},
			expectedError: customErrors.ErrUserNotExist ,
			expectedResult: []projectdb.Project{},
		},
		{
			name:"Error",
			userID:02020,
			mockSetup: func() {
				mockRepo.On("GetProjectsByUserId",context.TODO(),int32(02020)).Return(
					[]projectdb.Project{},errors.New("db error"),
				)

			},
			expectedError: errors.New("db error") ,
			expectedResult: []projectdb.Project{},
		},
	}


	for _,tc := range testCases{
		t.Run(tc.name, func(t *testing.T) {
		
		tc.mockSetup()
		projects, err := projectService.GetProjectsByUserId(context.TODO(), tc.userID)

		if tc.expectedError != nil {
			assert.Equal(t, tc.expectedError.Error(), err.Error()) // compare error values properly
			assert.Empty(t, projects)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, projects)
		}

		mockRepo.AssertCalled(t, "GetProjectsByUserId", mock.Anything, int32(tc.userID))
	})
	}
	// You may want to add the test loop here to run the test cases
}

func TestGetProjectByName(t *testing.T) {
	testCases := []struct {
		name           string
		userID 			int 
		projectName 	string
		mockSetup      func(params projectdb.GetProjectByNameParams)
		expectedError  error
		expectedResult *projectdb.Project
	}{
		{
			name: "Valid Project Name",
			userID:344,
			projectName: "Valid project name",
			mockSetup: func(params projectdb.GetProjectByNameParams) {
				mockRepo.On("GetProjectByName", context.TODO(), params).Return(projectdb.Project{
					Name:   "metrics dashboard setup",
					UserID: 23,
					Description: sql.NullString{
						String: "metrics is useful for monitoring",
						Valid:  true,
					},
				}, nil)
			},
			expectedError:  nil,
			expectedResult: &projectdb.Project{
				Name:   "metrics dashboard setup",
				UserID: 23,
				Description: sql.NullString{
					String: "metrics is useful for monitoring",
					Valid:  true,
				},
			},
		},
		{
			name: "Non Existent Project Name",
			userID: 672,
			projectName: "Non Existent Project Name",
			mockSetup: func(params projectdb.GetProjectByNameParams) {
				mockRepo.On("GetProjectByName", context.TODO(), params).Return(projectdb.Project{}, customErrors.ErrProjectNotExist)

			},
			expectedError:  customErrors.ErrProjectNotExist,
			expectedResult:nil,
			},
			{
			name: "Db Error",
			userID: 672,
			projectName: "Sample Project",
			mockSetup: func(params projectdb.GetProjectByNameParams) {
				mockRepo.On("GetProjectByName", context.TODO(), params).Return(projectdb.Project{}, errors.New("db error"))

			},
			expectedError:  errors.New("db error"),
			expectedResult:nil,
			},
			
		}
	

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params:=projectdb.GetProjectByNameParams{
				UserID: int32(tc.userID),
				Name: tc.projectName,
			}
			tc.mockSetup(params)
			result, err := projectService.GetProjectByName(context.TODO(),tc.projectName,tc.userID)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
			mockRepo.AssertCalled(t, "GetProjectByName", mock.Anything, params)
		})
	}
}



func TestUpdateProject(t *testing.T){


	testCases:=[]struct{
			testName string 
			req UpdateProjectRequest
			expectedParams projectdb.UpdateProjectParams
			expectedError error
	}{
			{
				testName: "Valid project update",
				req: UpdateProjectRequest{
					ProjectID: 1234,
					Name: func( s string)*string{return &s}("starpilot project"),
					Color: func( s string)*string{return &s}("yellow"),
				},
				expectedParams: projectdb.UpdateProjectParams{
					Name: sql.NullString{
						String: "starpilot project",
						Valid: true,
					},
					ID: 1234,
					Color: projectdb.NullProjectColor{
						ProjectColor: projectdb.ProjectColorYELLOW,
						Valid: true,
					},
				},
				expectedError: nil ,
			},
			{
				testName: "db connection error",
				req: UpdateProjectRequest{
					ProjectID: 1234,
					Name: func( s string)*string{return &s}("starpilot project"),
					Color: func( s string)*string{return &s}("yellow"),
				},
				expectedParams: projectdb.UpdateProjectParams{
					Name: sql.NullString{
						String: "starpilot project",
						Valid: true,
					},
					ID: 1234,
					Color: projectdb.NullProjectColor{
						ProjectColor: projectdb.ProjectColorYELLOW,
						Valid: true,
					},
				},
				expectedError: errors.New("db error") ,
			},

			
	}

	for _,tc := range testCases{
		t.Run(tc.testName,func(t *testing.T){
				mockRepo.Calls=nil 
				mockRepo.ExpectedCalls=nil 

				mockRepo.On("UpdateProject",mock.Anything,tc.expectedParams).Return(tc.expectedError)


				err:=projectService.UpdateProject(context.TODO(),tc.req)

				assert.Equal(t,err,tc.expectedError)

				mockRepo.AssertCalled(t,"UpdateProject",mock.Anything,tc.expectedParams)


		})
	}

}