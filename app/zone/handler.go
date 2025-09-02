package zone

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

	ctxMessage := "Page Zone"

	req := new(request.PageZone)

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
		Message:    "Success Get Zone",
		Data:       result,
	}

	c.JSON(http.StatusOK, response)

}
func (h Handler) Create(c *gin.Context) {
	var err error

	ctxMessage := "Create Zone"

	req := new(request.CreateZone)

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
		Message:    "Success Create Zone",
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) GetById(c *gin.Context) {
	var err error

	ctxMessage := "Get Zone"

	id := c.Param("id")

	data, err := h.usecase.GetById(c, id)
	if err != nil {
		response := helper.CreateResponseStatus(ctxMessage, "002", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := &model.Response{
		Status:     true,
		StatusCode: 200,
		Message:    "Success Get Zone",
		Data:       data,
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) Delete(c *gin.Context) {
	var err error

	ctxMessage := "DELETE Zone"

	id := c.Param("id")

	err = h.usecase.Delete(c, id)
	if err != nil {
		response := helper.CreateResponseStatus(ctxMessage, "002", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := &model.Response{
		Status:     true,
		StatusCode: 200,
		Message:    "Success Delete Zone",
	}

	c.JSON(http.StatusOK, response)
}
