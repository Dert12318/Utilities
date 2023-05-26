package v9

import (
	"database/sql/driver"
	"reflect"

	"github.com/Dert12318/Utilities/common/types"
	"github.com/Dert12318/Utilities/validator"
	"github.com/Dert12318/Utilities/validator/v9/custom"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	v9 "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type (
	implementation struct {
		instance *v9.Validate
		trans    ut.Translator
	}
)

func New() (validator.Validator, error) {
	langEn := en.New()
	langId := id.New()
	uni := ut.New(langEn, langEn, langId)
	trans, _ := uni.GetTranslator("en")

	validate := v9.New()
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, err
	}

	// register all types.Null* types to use the ValidateValuer CustomTypeFunc
	validate.RegisterCustomTypeFunc(ValidateValuer, types.NullString{}, types.NullInt32{}, types.NullInt64{}, types.NullBool{}, types.NullFloat64{}, types.NullTime{})

	instance := &implementation{instance: validate, trans: trans}
	if err := instance.registerDefaultValidator(); err != nil {
		return nil, err
	}
	return instance, nil
}

func (i *implementation) registerDefaultValidator() error {
	if err := i.instance.RegisterValidation("date", custom.IsDate); err != nil {
		return err
	}
	if err := i.instance.RegisterValidation("datetime", custom.IsDateTime); err != nil {
		return err
	}

	return nil
}

func (i *implementation) RegisterValidation(tag string, fn func(fl v9.FieldLevel) bool) {
	i.instance.RegisterValidation(tag, fn)
}

func (i *implementation) RegisterStructValidation(fn func(sl v9.StructLevel), types interface{}) {
	i.instance.RegisterStructValidation(fn, types)
}

func (i *implementation) Validate(object interface{}) error {
	if err := i.instance.Struct(object); err != nil {
		return err
	}
	return nil
}

func (i *implementation) ValidateVar(object interface{}, constraint string) error {
	if err := i.instance.Var(object, constraint); err != nil {
		return err
	}
	return nil
}

// ValidateValuer implements validator.CustomTypeFunc
func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}
	return nil
}
