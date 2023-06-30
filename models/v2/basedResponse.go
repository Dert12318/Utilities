package models

import (
	"strings"
)

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Code       int         `json:"code"`
	CustomCode string      `json:"customCode,omitempty"`
	Message    string      `json:"message"`
	Errors     interface{} `json:"errors,omitempty"`
}

func BuildSuccessResponse(message string, code int, data interface{}) SuccessResponse {
	res := SuccessResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, code int, err error) ErrorResponse {
	errorMessage := err.Error()

	splitError := strings.Split(errorMessage, "\n")
	res := ErrorResponse{
		Code:    code,
		Message: message,
		Errors:  splitError,
	}
	return res
}

func BuildCustomError(message, customCode string, code int, err error) ErrorResponse {
	errorMessage := err.Error()

	splitError := strings.Split(errorMessage, "\n")
	res := ErrorResponse{
		Code:       code,
		CustomCode: customCode,
		Message:    message,
		Errors:     splitError,
	}
	return res
}
