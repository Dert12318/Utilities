package custom

import (
	"time"

	timeConstant "github.com/Dert12318/Utilities/common/constant/time"
	"github.com/Dert12318/Utilities/common/functions"
	v9 "gopkg.in/go-playground/validator.v9"
)

func IsDateTime(fl v9.FieldLevel) bool {
	checkedValue := functions.ConvertReflectValueToString(fl.Field())

	_, err := time.Parse(timeConstant.DateTimeLayout, checkedValue)
	if err == nil {
		return true
	}

	_, err = time.Parse(timeConstant.DateTimeAltLayout, checkedValue)
	if err == nil {
		return true
	}

	return false
}
