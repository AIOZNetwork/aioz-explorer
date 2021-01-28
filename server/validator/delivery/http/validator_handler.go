package http

import (
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

type ValidatorHandler struct {
	ctx      context.Context
	VUsecase domain.ValidatorUsecase
}

func NewValidatorHandler(ctx context.Context, g *echo.Group, vu domain.ValidatorUsecase) {
	handler := &ValidatorHandler{
		ctx:      ctx,
		VUsecase: vu,
	}
	g.GET("/validators/:address", handler.getValidators)
	//g.GET("/validators/info", handler.getValidatorsInfo)
	//g.POST("/validators/update", handler.updateValidators)
}

func (vh *ValidatorHandler) getValidators(c echo.Context) error {
	address := c.Param("address")
	if address == "" {
		return c.JSON(http.StatusBadRequest, Response{Messages: fmt.Sprintf("Invalid validator address"), Success: false})
	}
	resp, err := vh.VUsecase.GetValidatorByValAddress(vh.ctx, address)
	if err != nil {
		vh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusInternalServerError, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: 1})
}

/*
func (vh *ValidatorHandler) getValidatorsInfo(c echo.Context) error {
	limit := int64(10)
	offset := int64(0)
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		limit, _ = strconv.ParseInt(limitStr, 10, 64)
	}
	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		offset, _ = strconv.ParseInt(offsetStr, 10, 64)
	}

	resp, err := vh.VUsecase.GetValidatorsInfo(vh.ctx, limit, offset)
	if err != nil {
		vh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusInternalServerError, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true})
}

func (vh *ValidatorHandler) updateValidators(c echo.Context) error {
	// todo: waiting for example
	return c.JSON(http.StatusOK, Response{Success: true})
}
*/