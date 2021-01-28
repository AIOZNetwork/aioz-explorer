package http

import (
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

type StatHandler struct {
	ctx         context.Context
	StatUsecase domain.StatisticUsecase
}

func NewStatHandler(ctx context.Context, g *echo.Group, su domain.StatisticUsecase) {
	handler := &StatHandler{
		ctx:         ctx,
		StatUsecase: su,
	}

	g.GET("/statistic", handler.getStatistic)
}

// getStatistic godoc
// @Summary get network statistic
// @Tags statistic
// @Description Retrieve blockchain stat
// @Accept  json
// @Produce  json
// @Success 201 {object} domain.Statistic "Recent blocks"
// @Router /statistic [get]
func (sh *StatHandler) getStatistic(c echo.Context) error {
	resp, err := sh.StatUsecase.GetStatistic(sh.ctx)
	if err != nil {
		sh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusInternalServerError, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Success: true, Data: resp})
}
