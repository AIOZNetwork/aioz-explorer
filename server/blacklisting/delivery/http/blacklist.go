package http

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
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

type BlacklistHandler struct {
	ctx       context.Context
	BLUsecase domain.BlacklistUsecase
}

func NewBlacklistHandler(ctx context.Context, g *echo.Group, us domain.BlacklistUsecase) {
	handler := &BlacklistHandler{
		ctx:       ctx,
		BLUsecase: us,
	}
	g.POST("/blacklist/add", handler.add2Blacklist)
	g.POST("/blacklist/remove", handler.removeFromBlacklist)
	g.POST("/whitelist/add", handler.add2Whitelist)
	g.POST("/whitelist/remove", handler.removeFromWhitelist)
}

func (bh *BlacklistHandler) add2Blacklist(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}
	var request struct {
		IPS []string `json:"ips"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		return c.JSON(http.StatusInternalServerError, "Cannot unmarshal request")
	}
	err = bh.BLUsecase.Add2Blacklist(bh.ctx, request.IPS)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	return c.JSON(http.StatusOK, "Add IP to black list successfully")
}

func (bh *BlacklistHandler) removeFromBlacklist(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}
	var request struct {
		IPS []string `json:"ips"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		return c.JSON(http.StatusInternalServerError, "Cannot unmarshal request")
	}
	err = bh.BLUsecase.UnbanIP(bh.ctx, request.IPS)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	return c.JSON(http.StatusOK, "Remove IP from black list successfully")
}

func (bh *BlacklistHandler) add2Whitelist(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}
	var request struct {
		IPS []string `json:"ips"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		return c.JSON(http.StatusInternalServerError, "Cannot unmarshal request")
	}
	err = bh.BLUsecase.Add2Whitelist(bh.ctx, request.IPS)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	return c.JSON(http.StatusOK, "Add IP to white list successfully")
}

func (bh *BlacklistHandler) removeFromWhitelist(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}
	var request struct {
		IPS []string `json:"ips"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		return c.JSON(http.StatusInternalServerError, "Cannot unmarshal request")
	}
	err = bh.BLUsecase.RemoveFromWhitelist(bh.ctx, request.IPS)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	return c.JSON(http.StatusOK, "Remove IP from whitelist successfully")
}
