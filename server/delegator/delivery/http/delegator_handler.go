package http

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"swagger-server/context"
	"swagger-server/domain"
)

type Response struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	Messages interface{} `json:"messages"`
	Total    int64       `json:"total"`
}

type DelegatorHandler struct {
	ctx      context.Context
	DUsecase domain.DelegatorUsecase
}

func NewDelegatorHandler(ctx context.Context, g *echo.Group, du domain.DelegatorUsecase) {
	handler := &DelegatorHandler{
		ctx:      ctx,
		DUsecase: du,
	}
	g.GET("/delegators/:address", handler.getDelegators)
}

func (dh *DelegatorHandler) getDelegators(c echo.Context) error {
	address := c.Param("address")
	if address == "" {
		dh.ctx.Logger().WithError(errors.New("invalid acc address"))
		return c.JSON(http.StatusBadRequest, Response{Messages: fmt.Sprintf("Invalid acc address"), Success: false})
	}
	resp, total, err := dh.DUsecase.GetDelegatorByAccAddress(dh.ctx, address)
	if err != nil {
		dh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusInternalServerError, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: total})
}
