package models

import (
	"strings"
)

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data"`
}

type SuccessPageResponse struct {
	SuccessResponse
	Page interface{} `json:"page,omitempty"`
}

type Page struct {
	Offset    int   `json:"offset"`
	Limit     int   `json:"limit"`
	TotalData int64 `json:"total_data"`
}

func BuildSuccessResponse(message string, status int, data interface{}) SuccessResponse {
	res := SuccessResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, status int, err error, data interface{}) ErrorResponse {
	errorMessage := err.Error()

	splitError := strings.Split(errorMessage, "\n")
	res := ErrorResponse{
		Status:  status,
		Message: message,
		Errors:  splitError,
		Data:    data,
	}
	return res
}

func BuildSuccessPageResponse(status int, message string, data interface{}, args ...interface{}) SuccessPageResponse {
	res := SuccessPageResponse{}
	res.Status = status
	res.Message = message
	if data != "" && data != nil {
		res.Data = data
	}
	if len(args) > 0 {
		res.Page = args[0]
	}
	return res
}
