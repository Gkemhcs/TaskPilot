package user

import (
	"context"
	"errors"
	"os"
	"testing"

	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockRepo *MockUserRepo
var userService *UserService

// TestMain sets up the mock repository and user service before running all tests.
// This function is called automatically by the Go test runner.
func TestMain(m *testing.M) {
	mockRepo = new(MockUserRepo)
	userService = NewUserService(mockRepo)
	os.Exit(m.Run())
}

func TestGeneratePasswordHash(t *testing.T) {
	password := "woeormwmro"
	hash, err := userService.GeneratePasswordHash(password)

	assert.NotEmpty(t, hash)
	assert.NoError(t, err)
	assert.NoError(t, userService.CheckHashedPassword(password, hash))

}
func getHashedPassword(password string) string {
	hashedPassword, err := userService.GeneratePasswordHash(password)
	if err != nil {
		return ""
	}
	return hashedPassword

}
func TestCheckHashedPassword(t *testing.T) {
	testCases := []struct {
		name           string
		password       string
		hashedPassword string
		expectedError  bool
	}{
		{name: "valid password hash",
			password:       "abcd234",
			hashedPassword: getHashedPassword("abcd234"),
			expectedError:  false,
		}, {
			name:           "invalid password hash",
			password:       "abcd1234",
			hashedPassword: "1o1ok1o3ko1",
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userService.CheckHashedPassword(tc.password, tc.hashedPassword)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		inputName     string
		inputEmail    string
		inputPassword string
		mockSetup     func()
		expectedError error
	}{
		{
			name:          "Success",
			inputName:     "Koti",
			inputEmail:    "koti@example.com",
			inputPassword: "oeo33o",
			mockSetup: func() {
				mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("userdb.CreateUserParams")).
					Return(userdb.User{ID: 1, Name: "Koti"}, nil)
			},
			expectedError: nil,
		},

		{
			name:          "User already exists (duplicate key)",
			inputName:     "Koti",
			inputEmail:    "koti@example.com",
			inputPassword: "r1ir3irnk",
			mockSetup: func() {
				pqErr := &pq.Error{
					Code:       "23505",
					Constraint: "users_pkey",
				}
				mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("userdb.CreateUserParams")).
					Return(userdb.User{}, pqErr)
			},
			expectedError: customErrors.ErrUserAlreadyExists,
		},
		{
			name:          "Other DB error",
			inputName:     "Koti",
			inputEmail:    "koti@gkemhcs.com",
			inputPassword: "j1ijojo",
			mockSetup: func() {
				mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("userdb.CreateUserParams")).
					Return(userdb.User{}, errors.New("db failure"))
			},
			expectedError: errors.New("db failure"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			err := userService.CreateUser(context.TODO(), tc.inputName, tc.inputEmail, tc.inputPassword)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.ExpectedCalls = nil // reset mock after each test case
		})
	}

}

func TestLoginUser(t *testing.T) {
	testCases := []struct {
		name           string
		email          string
		password       string
		hashedPassword string
		mockSetup      func()
		expectedError  error
		expectedResult *userdb.User
	}{
		{
			name:           "Existing user",
			email:          "gkemhcs@gmail.com",
			password:       "gkemhcs",
			hashedPassword: getHashedPassword("gkemhcs"),
			mockSetup: func() {
				mockRepo.On("GetUserByEmail", mock.Anything, "gkemhcs@gmail.com").Return(
					userdb.User{
						Name:           "gkemhcs",
						HashedPassword: getHashedPassword("gkemhcs"),
						Email:          "gkemhcs@gmail.com",
					}, nil)
			},
			expectedError: nil,
			expectedResult: &userdb.User{
				Name:  "gkemhcs",
				Email: "gkemhcs@gmail.com",
			},
		}, {
			name:           "Non existing User",
			email:          "gkemhcs@yahoo",
			password:       "gkemhcs",
			hashedPassword: getHashedPassword("gkemhcs"),
			mockSetup: func() {
				mockRepo.On("GetUserByEmail", mock.Anything, "gkemhcs@yahoo").Return(
					userdb.User{},
					customErrors.USER_NOT_FOUND,
				)
			},
			expectedError:  customErrors.USER_NOT_FOUND,
			expectedResult: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			user, err := userService.LoginUser(context.TODO(), tc.email, tc.password)
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult.Email, user.Email)
				assert.Equal(t, tc.expectedResult.Name, user.Name)
			}
			mockRepo.AssertCalled(t, "GetUserByEmail", mock.Anything, tc.email)

		})
	}

}
