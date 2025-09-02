package types

import (
	"ai-dekadns/helper"
	"ai-dekadns/model"
	"ai-dekadns/request"
	"ai-dekadns/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leebenson/conform"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) Handler {
	return Handler{usecase: usecase}
}

func (h Handler) Page(c *gin.Context) {
	var err error

	ctxMessage := "Page Type"

	req := new(request.PageType)

	err = c.Bind(req)
	if err != nil {
		response := helper.CreateResponseStatus(ctxMessage, "001", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//trims whitespaces
	conform.Strings(req)

	// validate struct
	v := validator.Validator
	err = v.Struct(req)
	if err != nil {
		errValidator := validator.MapValidationErrors(err)
		response := &model.Response{
			Status:     false,
			StatusCode: http.StatusBadRequest,
			Message:    "error validation",
			Data:       errValidator,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	result, err := h.usecase.Page(c, *req)
	if err != nil {
		response := helper.CreateResponseStatus(ctxMessage, "002", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := &model.Response{
		Status:     true,
		StatusCode: 200,
		Message:    "Success Get Type",
		Data:       result,
	}

	c.JSON(http.StatusOK, response)

}
