package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"swagger-server/context"
	"swagger-server/domain"
	"swagger-server/utils"
)

type Response struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	Messages interface{} `json:"messages"`
	Total    int64       `json:"total"`
}

type StakingHandler struct {
	ctx      context.Context
	SUsecase domain.StakingUsecase
}

func NewStakingHandler(ctx context.Context, g *echo.Group, su domain.StakingUsecase) {
	handler := &StakingHandler{
		ctx:      ctx,
		SUsecase: su,
	}

	g.GET("/staking/validators", handler.getTopValidators)
	g.GET("/staking/wallets", handler.getTopStakingWallets)
	g.GET("/staking/tokens", handler.getTotalStakes)
}

// getTopValidators godoc
// @Summary get top validators
// @Tags staking
// @Description Retrieve top validators
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 201 {object} domain.Validator "Top validators"
// @Router /staking/validators [get]
func (sh *StakingHandler) getTopValidators(c echo.Context) error {
	limit := int64(10)
	offset := int64(0)
	var err error
	if c.QueryParam("limit") != "" {
		param := c.QueryParam("limit")
		limit, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			sh.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Messages: err.Error(), Success: false})
		}
	}
	if c.QueryParam("offset") != "" {
		param := c.QueryParam("offset")
		offset, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			sh.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Messages: err.Error(), Success: false})
		}
	}
	resp, total, err := sh.SUsecase.GetTopValidators(sh.ctx, int(limit), int(offset))
	if err != nil {
		sh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: total})
}

// getTopStakingWallets godoc
// @Summary get top staking wallets
// @Tags staking
// @Description Retrieve top staking wallets
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 201 {object} domain.Delegator "Top delegators"
// @Router /staking/wallets [get]
func (sh *StakingHandler) getTopStakingWallets(c echo.Context) error {
	limit := int64(10)
	offset := int64(0)
	var err error
	if c.QueryParam("limit") != "" {
		param := c.QueryParam("limit")
		limit, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			sh.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Messages: err.Error(), Success: false})
		}
	}
	if c.QueryParam("offset") != "" {
		param := c.QueryParam("offset")
		offset, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			sh.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Messages: err.Error(), Success: false})
		}
	}
	resp, total, err := sh.SUsecase.GetTopStakingWallets(sh.ctx, int(limit), int(offset))
	if err != nil {
		sh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: total})
}

// getTotalStakes godoc
// @Summary get total staked tokens
// @Tags staking
// @Description Retrieve total staked tokens
// @Accept  json
// @Produce  json
// @Success 201 "Total staked tokens"
// @Router /staking/tokens [get]
func (sh *StakingHandler) getTotalStakes(c echo.Context) error {
	resp, err := sh.SUsecase.GetTotalStakes(sh.ctx)
	if err != nil {
		sh.ctx.Logger().WithField("/staking/tokens", err)
		return c.JSON(http.StatusBadRequest, Response{Messages: err.Error(), Success: false})
	}
	result := utils.RemoveTrailingZerosFromDec(resp)
	return c.JSON(http.StatusOK, Response{Data: result, Success: true})
}
