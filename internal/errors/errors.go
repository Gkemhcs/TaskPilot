package errors

import "errors"

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

var ErrUserAlreadyExists=errors.New("User with the same name already exists")