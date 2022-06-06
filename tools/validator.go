package tools

import "github.com/go-playground/validator/v10"

var onlyCapitalAndLower validator.Func = func(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		for _, ch := range value {
			if ch < 'A' || ch > 'z' {
				return false
			}
		}
		return true
	}
	return false
}

var mustCapitalLowerNum validator.Func = func(fl validator.FieldLevel) bool {
	var flag1, flag2, flag3 = false, false, false
	if value, ok := fl.Field().Interface().(string); ok {
		for _, ch := range value {
			if ch >= 'A' && ch <= 'Z' {
				flag1 = true
			} else if ch >= 'a' && ch <= 'z' {
				flag2 = true
			} else if ch >= '0' && ch <= '9' {
				flag3 = true
			}
			if flag1 && flag2 && flag3 {
				return true
			}
		}
		return false
	}
	return false
}
