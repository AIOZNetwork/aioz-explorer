package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"swagger-server/context"
	"swagger-server/domain"
)

type Response struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	Messages interface{} `json:"messages"`
	Total    int64       `json:"total"`
}

type TransactionHandler struct {
	ctx      context.Context
	TUsecase domain.TransactionUsecase
}

func NewTransactionHandler(ctx context.Context, g *echo.Group, us domain.TransactionUsecase) {
	handler := &TransactionHandler{
		ctx:      ctx,
		TUsecase: us,
	}
	g.GET("/transaction/:hash", handler.getTransactionByHash)
	g.GET("/transactions", handler.getTransactions)
	g.GET("/transactions/height/:height", handler.getTransactionsByHeight)
}

// getTransactionByHash godoc
// @Summary get transaction by hash
// @Tags transaction
// @Description Retrieve transaction detail
// @Accept  json
// @Produce  json
// @Param hash path string true "Transaction hash"
// @Success 201 {object} domain.Transaction "Transaction detail"
// @Router /transaction/{hash} [get]
func (th *TransactionHandler) getTransactionByHash(c echo.Context) error {
	hash := c.Param("hash")
	if hash == "" {
		th.ctx.Logger().WithError(errors.New("invalid transaction hash"))
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid transaction hash"})
	}
	resp, err := th.TUsecase.GetTransaction(th.ctx, hash)
	if err != nil {
		th.ctx.Logger().WithError(err)
		return c.JSON(http.StatusNotFound, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: 1})
}

// getTransactions godoc
// @Summary get recent transactions
// @Tags transaction
// @Description Retrieve recent transactions
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 201 {object} domain.Transaction "Recent transactions"
// @Router /transactions [get]
func (th *TransactionHandler) getTransactions(c echo.Context) error {
	var err error
	limit := int64(10)
	offset := int64(0)
	if c.QueryParam("limit") != "" {
		param := c.QueryParam("limit")
		limit, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			th.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
		}
	}
	if c.QueryParam("offset") != "" {
		param := c.QueryParam("offset")
		offset, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			th.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
		}
	}
	resp, total, err := th.TUsecase.GetTransactions(th.ctx, limit, offset)
	if err != nil {
		th.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: total})
}

// getTransactionsByHeight godoc
// @Summary get transactions by block height
// @Tags transaction
// @Description Retrieve transactions by block height
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Param height path string true "Block height"
// @Success 201 {object} domain.Transaction "Transaction detail"
// @Router /transaction/height/{height} [get]
func (th *TransactionHandler) getTransactionsByHeight(c echo.Context) error {
	var err error
	height := c.Param("height")
	if height == "" {
		th.ctx.Logger().WithField("/transaction/height/{height}", "Invalid block height")
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid block height"})
	}
	limit := int64(10)
	if c.QueryParam("limit") != "" {
		param := c.QueryParam("limit")
		limit, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			th.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
		}
	}
	offset := int64(0)
	if c.QueryParam("offset") != "" {
		param := c.QueryParam("offset")
		offset, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			th.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
		}
	}
	h, _ := strconv.ParseInt(height, 10, 64)
	resp, total, err := th.TUsecase.GetTransactionsByBlockHeight(th.ctx, h, limit, offset)
	if err != nil {
		th.ctx.Logger().WithError(err)
		return c.JSON(http.StatusNotFound, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: total})
}
