package record

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

func (h Handler) Create(c *gin.Context) {
	var err error

	ctxMessage := "Create Record"

	req := new(request.CreateRecord)

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

	err = h.usecase.Create(c, *req)
	if err != nil {
		response := helper.CreateResponseStatus(ctxMessage, "002", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := &model.Response{
		Status:     true,
		StatusCode: 200,
		Message:    "Success Create Record",
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) Update(c *gin.Context) {
	var err error

	ctxMessage := "Update Record"

	req := new(request.UpdateRecord)

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

	err = h.usecase.Update(c, *req)
	if err != nil {
		response := helper.CreateResponseStatus(ctxMessage, "002", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := &model.Response{
		Status:     true,
		StatusCode: 200,
		Message:    "Success Update Record",
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) Delete(c *gin.Context) {
	var err error

	ctxMessage := "Delete Record"

	req := new(request.DeleteRecord)

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

	err = h.usecase.Delete(c, *req)
	if err != nil {
		response := helper.CreateResponseStatus(ctxMessage, "002", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := &model.Response{
		Status:     true,
		StatusCode: 200,
		Message:    "Success Delete Record",
	}

	c.JSON(http.StatusOK, response)
}
