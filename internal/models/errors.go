package models

import "errors"

var (
	MarshallingSessionError = errors.New("marshalling session error")
	AddToRedisError         = errors.New("error occurred while adding data to Redis")
	DeleteFromRedisError    = errors.New("error occurred while removing data from Redis")
	SessionNotExistsError   = errors.New("session does not exist")

	IncorrectPasswordLen = errors.New("incorrect password length")
	PasswordsNotMatch    = errors.New("passwords do not match")

	UserNotExists     = errors.New("user does not exist")
	UserAlreadyExists = errors.New("user already exists")

	PermissionDeniedError = errors.New("permission denied")
	InvalidInputError     = errors.New("invalid input")
	InvalidFilenameError  = errors.New("invalid filename")
)
