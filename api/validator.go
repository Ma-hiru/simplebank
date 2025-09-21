package api

import (
	"github.com/Ma-hiru/simplebank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	var currency, ok = fieldLevel.Field().Interface().(string)
	if !ok {
		return ok
	}

	return util.IsSupportCurrency(currency)
}
