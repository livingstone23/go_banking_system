package api



import (
	"go_banking_system/util"
	"github.com/go-playground/validator/v10"
)



var validCurrency validator.Func = func(fielLevel validator.FieldLevel) bool {
	if currency, ok := fielLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
