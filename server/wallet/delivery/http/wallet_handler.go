package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"swagger-server/context"
	"swagger-server/domain"
)

type Response struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	Messages interface{} `json:"messages"`
	Total    int64       `json:"total"`
}

type WalletHandler struct {
	ctx      context.Context
	WUsecase domain.WalletUsecase
}

func NewWalletHandler(ctx context.Context, g *echo.Group, wu domain.WalletUsecase) {
	handler := &WalletHandler{
		ctx:      ctx,
		WUsecase: wu,
	}
	g.GET("/wallet/:address", handler.getWallet)
	g.GET("/wallet/txs/:address", handler.getTxsByWallet)
	g.POST("/wallet/contacts", handler.getContactsByWallet)
	g.POST("/key/new", handler.newKey)
	g.POST("/key/recover", handler.recoverKey)
	g.POST("/key/encrypt", handler.encryptKey)
	g.POST("/key/decrypt", handler.decryptKey)
}

// getWallet godoc
// @Summary get wallet information
// @Tags wallet
// @Description Retrieve information of wallet
// @Accept  json
// @Produce  json
// @Param address path string true "Wallets address"
// @Success 200
// @Router /wallet/{address} [get]
func (wh *WalletHandler) getWallet(c echo.Context) error {
	address := c.Param("address")
	if address == "" {
		wh.ctx.Logger().WithError(errors.New("invalid wallet address"))
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid wallet address"})
	}
	tokens := strings.Split(address, "-")
	if len(tokens) < 2 {
		wh.ctx.Logger().WithError(errors.New("invalid wallet address"))
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: fmt.Sprintf("Invalid address")})
	}
	switch tokens[0] {
	case domain.ValidatorAddressFormat:
		resp, err := wh.WUsecase.GetValidatorByAddress(wh.ctx, address)
		if err != nil {
			wh.ctx.Logger().WithError(err)
			return c.JSON(http.StatusNotFound, Response{Success: false, Messages: err.Error()})
		}
		return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: 1})
	case domain.AccAddressFormat:
		resp, err := wh.WUsecase.GetWalletByAddress(wh.ctx, address)
		if err != nil {
			wh.ctx.Logger().WithError(err)
			return c.JSON(http.StatusNotFound, Response{Success: false, Messages: err.Error()})
		}
		return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: 1})
	default:
		wh.ctx.Logger().WithError(errors.New("invalid address"))
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: fmt.Sprintf("Invalid address")})
	}
}

// getTxsByWallet godoc
// @Summary get wallet's transactions
// @Tags wallet
// @Description Retrieve all transactions of wallet
// @Accept  json
// @Produce  json
// @Param address path string true "Wallets address"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 201 {object} domain.Txs "Transactions history"
// @Router /wallet/txs/{address} [get]
func (wh *WalletHandler) getTxsByWallet(c echo.Context) error {
	var err error
	address := c.Param("address")
	if address == "" {
		wh.ctx.Logger().WithError(errors.New("invalid wallet address"))
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid address"})
	}
	limit := int64(10)
	if c.QueryParam("limit") != "" {
		param := c.QueryParam("limit")
		limit, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			wh.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
		}
	}
	offset := int64(0)
	if c.QueryParam("offset") != "" {
		param := c.QueryParam("offset")
		offset, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			wh.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
		}
	}
	resp, total, err := wh.WUsecase.GetTxsByAddressV2(wh.ctx, address, int(limit), int(offset))
	if err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: total})
}

// newKey godoc
// @Summary create a key for user
// @Tags key
// @Description Create a key for user
// @Accept  json
// @Produce  json
// @Param key body domain.RequestNewKey true "Data to create key"
// @Success 201 {object} domain.KeyResponse "The created key"
// @Failure 400
// @Router /key/new [post]
func (wh *WalletHandler) newKey(c echo.Context) error {
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	var req *domain.RequestNewKey
	if err := json.Unmarshal(reqBody, &req); err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	resp, err := wh.WUsecase.CreateWallet(wh.ctx, req)
	if err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: 1})
}

// recoverKey godoc
// @Summary recover key from mnemonic
// @Tags key
// @Description recover key from mnemonic
// @Accept  json
// @Produce  json
// @Param key body domain.RequestRecoverKey true "Data to recover key"
// @Success 201 {object} domain.KeyResponse "The recovered key"
// @Failure 400
// @Router /key/recover [post]
func (wh *WalletHandler) recoverKey(c echo.Context) error {
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	var req *domain.RequestRecoverKey
	if err := json.Unmarshal(reqBody, &req); err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	resp, err := wh.WUsecase.RecoverWallet(wh.ctx, req)
	if err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}

	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: 1})
}

// encryptKey godoc
// @Summary encrypt key with password
// @Tags key
// @Description encrypt key with password
// @Accept  json
// @Produce  json
// @Param key body domain.RequestEncrKey true "Password to encrypt key"
// @Success 201 {object} domain.KeyResponse "The encrypted key"
// @Failure 400
// @Router /key/encrypt [post]
func (wh *WalletHandler) encryptKey(c echo.Context) error {
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	var req *domain.RequestEncrKey
	if err := json.Unmarshal(reqBody, &req); err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	resp, err := wh.WUsecase.EncryptKey(wh.ctx, req)
	if err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}

	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: 1})
}

// decryptKey godoc
// @Summary decrypt key with password
// @Tags key
// @Description decrypt key with password
// @Accept  json
// @Produce  json
// @Param key body domain.RequestDecrKey true "Password to decrypt key"
// @Success 201 {object} domain.KeyResponse "The decrypted key"
// @Failure 400
// @Router /key/decrypt [post]
func (wh *WalletHandler) decryptKey(c echo.Context) error {
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	var req *domain.RequestDecrKey
	if err := json.Unmarshal(reqBody, &req); err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	resp, err := wh.WUsecase.DecryptKey(wh.ctx, req)
	if err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}

	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: 1})
}

// getContactsByWallet godoc
// @Summary get related wallets of given addresses
// @Tags wallet
// @Description get related wallets of given address
// @Accept  json
// @Produce  json
// @Param addresses body []string true "List of addresses"
// @Success 201 {object} domain.ContactsWallet "The decrypted key"
// @Failure 400
// @Router /wallet/contacts [post]
func (wh *WalletHandler) getContactsByWallet(c echo.Context) error {
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	var req struct {
		Addresses []string `json:"addresses"`
	}
	if err := json.Unmarshal(reqBody, &req); err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}

	resp, err := wh.WUsecase.GetContactsWallet(wh.ctx, req.Addresses)
	if err != nil {
		wh.ctx.Logger().WithError(err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true})
}
