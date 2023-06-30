package response

import (
	"errors"
)

//	error list of codes
var (
	ErrorInvalidStock = errors.New("NOT ENOUGH STOCK")
)

// list of codes for response
const (
	STATUSCODE_SUCCESSS = "200000"

	// bad request
	STATUSCODE_BADREQUEST   = "400000"
	STATUSCODE_INVALIDSTOCK = "400001"

	// internal server error
	STATUSCODE_INTERNAL_ERROR = "500000"
)

var customCodeMap = map[string]error{
	STATUSCODE_INVALIDSTOCK: ErrorInvalidStock,
}

func GetMessageFromCode(code string) string {
	if val, ok := customCodeMap[code]; ok {
		return val.Error()
	}

	if code == STATUSCODE_SUCCESSS {
		return "Success"
	}

	return "Invalid status code"
}

func GetErrorFromCode(code string) error {
	if val, ok := customCodeMap[code]; ok {
		return val
	}

	if code == STATUSCODE_SUCCESSS {
		return nil
	}

	return errors.New("Invalid status code")
}

// GetErrorCode return code for error
func GetErrorCode(err error) string {
	switch err {
	case ErrorInvalidStock:
		return STATUSCODE_INVALIDSTOCK

	case nil:
		return STATUSCODE_SUCCESSS

	default:
		return STATUSCODE_INTERNAL_ERROR
	}
}
