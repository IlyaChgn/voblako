package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	StatusOk = 200

	StatusBadRequest   = 400
	StatusUnauthorized = 401
	StatusForbidden    = 403

	StatusInternalServerError = 500
)

const (
	ErrInternalServer = "Server error"

	ErrAuthorized    = "User already authorized"
	ErrNotAuthorized = "User not authorized"
	ErrForbidden     = "User have no access to this content"

	ErrWrongPasswordFormat = "Password must have length between 8 and 32 symbols"
	ErrDoNotMatch          = "Passwords do not match"
	ErrWrongCredentials    = "Wrong credentials"
	ErrAlreadyExists       = "User with this email already exists"

	ErrWrongFilename = "Filename must have length between 1 and 50"

	ErrBadJSON          = "Wrong JSON format"
	ErrBadForm          = "Wrong form format"
	ErrInvalidID        = "Invalid ID format"
	ErrInvalidURLParams = "Invalid URL params"
)

type ErrResponse struct {
	Status string `json:"status"`
}

func newErrResponse(status string) *ErrResponse {
	return &ErrResponse{
		Status: status,
	}
}

func sendResponse(writer http.ResponseWriter, response any) {
	serverResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Something went wrong while marshalling JSON", err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}

	_, err = writer.Write(serverResponse)
	if err != nil {
		log.Println("Something went wrong while senddng response", err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}
}

func SendOkResponse(writer http.ResponseWriter, body any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(StatusOk)

	sendResponse(writer, body)
}

func SendErrResponse(writer http.ResponseWriter, code int, status string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)

	response := newErrResponse(status)

	sendResponse(writer, response)
}
