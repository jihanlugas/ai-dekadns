package helper

import (
	"ai-dekadns/model"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func CreateResponseStatus(function, code, message string) *model.ResponseStatus {
	response := &model.ResponseStatus{
		ResponseCode:    fmt.Sprintf("%s-%s", function, code),
		ResponseMessage: message,
	}
	return response
}

func CreateValidationResponseStatus(function, code string, err error) *model.ResponseStatus {
	fmt.Println(err)
	var errMessage string
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		validationErrors := err.(validator.ValidationErrors)
		for _, x := range validationErrors {
			if errMessage != "" {
				errMessage = strings.Trim(errMessage, " ") + ","
			}
			if x.Tag() == "gt" || x.Tag() == "lt" || x.Tag() == "lte" || x.Tag() == "gte" {
				errMessage = fmt.Sprintf("%s %s %s", errMessage, "Invalid value of", x.Field())
			} else {
				errMessage = fmt.Sprintf("%s %s %s", errMessage, x.Field(), x.Tag())
			}
		}
	} else {
		errMessage = err.Error()
	}
	response := &model.ResponseStatus{
		ResponseCode:    fmt.Sprintf("%s-%s", function, code),
		ResponseMessage: strings.Trim(errMessage, " "),
	}
	return response
}
