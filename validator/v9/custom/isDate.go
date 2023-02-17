package custom

import (
	"time"

	timeConstant "github.com/Dert12318/Utilities/common/constant/time"
	"github.com/Dert12318/Utilities/common/functions"
	v9 "gopkg.in/go-playground/validator.v9"
)

func IsDate(fl v9.FieldLevel) bool {
	checkedValue := functions.ConvertReflectValueToString(fl.Field())

	_, err := time.Parse(timeConstant.DateLayout, checkedValue)
	if err != nil {
		return false
	}

	return true
}
