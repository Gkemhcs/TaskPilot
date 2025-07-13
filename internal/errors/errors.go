package errors

import (
	"errors"

	"github.com/lib/pq"
)


const (
    UniqueViolationErr = pq.ErrorCode("23505")
)
// Custom error variables for user authentication and request validation
var USER_NOT_FOUND = errors.New("user not found")

var INCORRECT_PASSWORD = errors.New("Enter the right password")

var MISSING_AUTHORIZATION_HEADER = errors.New("Authorization header is missing")
var INVALID_TOKEN_FORMAT = errors.New("invalid token format")

var INVALID_TOKEN = errors.New("InvalidToken Entered")

var INCORRECT_REQUEST_BODY = errors.New("incorrect request body format check for any missing parameters")

var MISSING_USER_NAME = errors.New("username is missing in request body")
var MISSING_PASSWORD = errors.New("password is missing in request body")

var ErrMismatchedPassword = errors.New("password doesnt match enter correct password")

var ErrUserAlreadyExists = errors.New("User with the same name already exists")
var ErrMissingEmail = errors.New("Email is missing in body")

var ErrMissingJwtToken = errors.New("authorization token missing")

var ErrUserNotExist = errors.New("User not found, first create the user")

var ErrInvalidProjectId = errors.New("Invalid Project ID")

var ErrInvalidUserId = errors.New("Invalid User Id set in context")
var ErrUserIDNotFoundInContext = errors.New("UserID is missing from context, try to pass access-token in body")
var ErrProjectNotExist = errors.New("sorry the project name you are searching for isn't found")
var ErrProjectIDNotExist = errors.New("sorry the project id you are searching for isn't found")
var ErrProjectsEmpty = errors.New("projects are empty")
var ErrProjectAlreadyExists=errors.New("project already exists")

var ErrTokenExpired = errors.New("access token has expired")

var ErrTaskNotFound=errors.New("sorry the task not found")

var ErrAssigneeMissingFromBody=errors.New("assignee email  is missing from body")
var ErrMissingDueDate=errors.New("Due date is missing from body")

var ErrorTaskTitleMissing=errors.New("task title is missing from request body")

var ErrMissingProjectID=errors.New("project id is missing from request body")

var ErrInvalidTaskID=errors.New("Invalid Task Id Entered")

var ErrTaskAlreadyExists=errors.New("Task Already exists")

var ErrTasksAreEmpty=errors.New("not tasks under the project id you mentioned")


var ErrParentProjectIDNotFound=errors.New("The corresponding project id  doesnt exist")



var ErrCreatingImportJob = errors.New("failed to create import job")

var ErrUploadingFile = errors.New("failed to upload file")

var ErrInvalidFileType = errors.New("invalid file type, only Excel files are allowed")

var ErrWhileEnqueuingImportJob = errors.New("failed to enqueue import job try after sometime")

var ErrGoogleApplicationCredentialsNotSet = errors.New("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")


var ErrLoadingServiceAccountFile=errors.New("failed to load service account file, check the path and permissions")


var ErrInvalidServiceAccountFile=errors.New("failed to unmarshal service account file, check the file format")


var ErrGeneratingSignedURL = errors.New("failed to generate signed URL for file download")

var ErrCreatingExportJob = errors.New("failed to create export job")


var ErrWhileEnqueuingExportJob = errors.New("failed to enqueue export job, try after sometime")

var ErrInvalidJobID=errors.New("invalid job id try passing valid id")