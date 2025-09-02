package types

import (
	"ai-dekadns/model"
	"ai-dekadns/request"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	Page(c *gin.Context, req request.PageType) (pagination *model.Pagination, err error)
}

type usecase struct {
	typeRepo Repository
}

func (u usecase) Page(c *gin.Context, req request.PageType) (pagination *model.Pagination, err error) {
	page, _ := strconv.Atoi(req.Page)
	limit, _ := strconv.Atoi(req.Limit)
	pageReq := model.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  "name asc", // case types sort di hardcode
	}

	pagination, err = u.typeRepo.Page(req, pageReq)
	if err != nil {
		return pagination, err
	}

	return pagination, nil
}

func NewUsecase(typeRepo Repository) Usecase {
	return &usecase{
		typeRepo: typeRepo,
	}
}
