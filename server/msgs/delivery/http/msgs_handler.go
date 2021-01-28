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

type MsgsHandler struct {
	ctx         context.Context
	MsgsUsecase domain.MsgsUsecase
}

func NewTxsHandler(ctx context.Context, g *echo.Group, tu domain.MsgsUsecase) {
	handler := &MsgsHandler{
		ctx:         ctx,
		MsgsUsecase: tu,
	}

	g.GET("/msgs", handler.getRecentMsgs)
	g.GET("/msgs/:address", handler.getMsgsOfWalletVersion2)
	g.GET("/msgs/send", handler.getMsgSendHistory) // api cho bede Dai
}

// getRecentMsgs godoc
// @Summary get recent messages
// @Tags message
// @Description Retrieve recent messages
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 201 {object} domain.Txs "Recent messages"
// @Router /msgs [get]
func (t *MsgsHandler) getRecentMsgs(c echo.Context) error {
	var err error
	limit := int64(10)
	offset := int64(0)
	if c.QueryParam("limit") != "" {
		param := c.QueryParam("limit")
		limit, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			t.ctx.Logger().WithError(errors.New("invalid limit param"))
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid limit param"})
		}
	}
	if c.QueryParam("offset") != "" {
		param := c.QueryParam("offset")
		offset, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			t.ctx.Logger().WithError(errors.New("invalid offset param"))
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid offset param"})
		}
	}
	resp, total, err := t.MsgsUsecase.GetRecentMsgs(t.ctx, int(limit), int(offset))
	if err != nil {
		t.ctx.Logger().WithError(err)
		return c.JSON(http.StatusInternalServerError, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: total})
}

// getMsgsOfWalletVersion2 godoc
// @Summary get recent messages
// @Tags message
// @Description Retrieve recent messages of a wallet
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Param address path string true "Wallet address"
// @Success 201 {object} domain.Txs "Recent messages of wallet"
// @Router /msgs/{address} [get]
func (t *MsgsHandler) getMsgsOfWalletVersion2(c echo.Context) error {
	address := c.Param("address")
	if address == "" {
		t.ctx.Logger().WithError(errors.New("invalid wallet address"))
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid wallet address"})
	}
	var err error
	limit := int64(10)
	offset := int64(0)
	if c.QueryParam("limit") != "" {
		param := c.QueryParam("limit")
		limit, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			t.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid limit param"})
		}
	}
	if c.QueryParam("offset") != "" {
		param := c.QueryParam("offset")
		offset, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			t.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid offset param"})
		}
	}
	resp, total, err := t.MsgsUsecase.GetMsgsByAddress(t.ctx, address, int(limit), int(offset))
	if err != nil {
		t.ctx.Logger().WithError(err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: total})
}

// getMsgSendHistory godoc
// @Summary get history messages send
// @Tags message
// @Description Retrieve history messages send
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 201 {object} domain.MessageSend "History messages send"
// @Router /msgs/send [get]
func (t *MsgsHandler) getMsgSendHistory(c echo.Context) error {
	heightFrom := c.QueryParam("from")
	if heightFrom == "" {
		t.ctx.Logger().WithError(errors.New("invalid block height from"))
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid block height from"})
	}
	heightTo := c.QueryParam("to")
	if heightTo == "" {
		t.ctx.Logger().WithError(errors.New("invalid block height to"))
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid block height to"})
	}
	var err error
	limit := int64(10)
	offset := int64(0)
	if c.QueryParam("limit") != "" {
		param := c.QueryParam("limit")
		limit, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			t.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid limit param"})
		}
	}
	if c.QueryParam("offset") != "" {
		param := c.QueryParam("offset")
		offset, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			t.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid offset param"})
		}
	}
	hf, err := strconv.ParseInt(heightFrom, 10, 64)
	if err != nil {
		t.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid height from param"})
	}
	ht, err := strconv.ParseInt(heightTo, 10, 64)
	if err != nil {
		t.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid height to param"})
	}
	resp, total, err := t.MsgsUsecase.GetMsgSendHistory(t.ctx, int(hf), int(ht), int(limit), int(offset))
	if err != nil {
		t.ctx.Logger().WithError(err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: total})
}
