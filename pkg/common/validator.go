package common

import (
	"fmt"
	"unicode"

	validator "gopkg.in/go-playground/validator.v9"
)

//NewValidator instanciates a new validator.
func NewValidator() *validator.Validate {
	answer := validator.New()
	var err error
	err = answer.RegisterValidation("kubeFriendlyName", kubeFriendlyName)
	if err != nil {
		LogExit("Cannot create validator:" + err.Error())
	}
	err = answer.RegisterValidation("javaPackageName", javaPackageName)
	if err != nil {
		LogExit("Cannot create validator:" + err.Error())
	}
	err = answer.RegisterValidation("startsWithUpperCaseAsciiAlpha", startsWithUpperCaseASCIIAlpha)
	if err != nil {
		LogExit("Cannot create validator:" + err.Error())
	}
	err = answer.RegisterValidation("startsWithLowerCaseAsciiAlpha", startsWithLowerCaseASCIIAlpha)
	if err != nil {
		LogExit("Cannot create validator:" + err.Error())
	}
	return answer
}

func kubeFriendlyName(fl validator.FieldLevel) bool {
	valueStr := fmt.Sprintf("%v", fl.Field())
	firstRun := []rune(valueStr[0:1])[0]
	if !unicode.IsLetter(firstRun) {
		return false
	}
	if unicode.IsUpper(firstRun) {
		return false
	}
	for _, val := range valueStr[1:] {
		runVal := rune(val)

		if !unicode.IsLetter(runVal) {
			if !unicode.IsNumber(runVal) {
				if runVal != '-' {
					return false
				}
			}
		} else {
			if unicode.IsUpper(runVal) {
				return false
			}
		}
	}
	return true
}

func javaPackageName(fl validator.FieldLevel) bool {
	valueStr := fmt.Sprintf("%v", fl.Field())
	firstRun := []rune(valueStr[0:1])[0]
	if !unicode.IsLetter(firstRun) {
		return false
	}
	if unicode.IsUpper(firstRun) {
		return false
	}
	for _, val := range valueStr[1:] {
		runVal := rune(val)
		if !unicode.IsLetter(runVal) {
			if !unicode.IsNumber(runVal) {
				if runVal != '.' {
					return false
				}
			}
		} else {
			if unicode.IsUpper(runVal) {
				return false
			}
		}
	}
	return true
}

func startsWithUpperCaseASCIIAlpha(fl validator.FieldLevel) bool {
	valueStr := fmt.Sprintf("%v", fl.Field())
	firstRun := []rune(valueStr[0:1])[0]
	if !unicode.IsLetter(firstRun) {
		return false
	}
	if !unicode.IsUpper(firstRun) {
		return false
	}
	return true
}
func startsWithLowerCaseASCIIAlpha(fl validator.FieldLevel) bool {
	valueStr := fmt.Sprintf("%v", fl.Field())
	firstRun := []rune(valueStr[0:1])[0]
	if !unicode.IsLetter(firstRun) {
		return false
	}
	if !unicode.IsLower(firstRun) {
		return false
	}
	return true
}
