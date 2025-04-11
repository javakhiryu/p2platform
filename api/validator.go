package api

import (
	"p2platform/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

var validSource validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if source, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedSource(source)
	}
	return false
}