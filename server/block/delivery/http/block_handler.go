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

type BlockHandler struct {
	ctx      context.Context
	BUsecase domain.BlockUsecase
}

func NewBlockHandler(ctx context.Context, g *echo.Group, us domain.BlockUsecase) {
	handler := &BlockHandler{
		ctx:      ctx,
		BUsecase: us,
	}
	g.GET("/blocks", handler.getLatestBlocks)
	g.GET("/block/hash/:hash", handler.getBlockByHash)
	g.GET("/block/height/:height", handler.getBlockByHeight)
}

// getLatestBlocks godoc
// @Summary get recent blocks
// @Tags block
// @Description Retrieve recent blocks
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 201 {object} domain.Block "Recent blocks"
// @Router /blocks [get]
func (b *BlockHandler) getLatestBlocks(c echo.Context) error {
	var err error
	limit := int64(10)
	offset := int64(0)
	if c.QueryParam("limit") != "" {
		param := c.QueryParam("limit")
		limit, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			b.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid limit param"})
		}
	}
	if c.QueryParam("offset") != "" {
		param := c.QueryParam("offset")
		offset, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			b.ctx.Logger().WithError(err)
			return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid offset param"})
		}
	}
	resp, total, err := b.BUsecase.GetLatestBlocks(b.ctx, int(offset), int(limit))
	if err != nil {
		b.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: total})
}

// getBlockByHash godoc
// @Summary get block by block hash
// @Tags block
// @Description Retrieve specific block by block hash
// @Accept  json
// @Produce  json
// @Param hash path string true "Block hash"
// @Success 201 {object} domain.Block "Block data"
// @Failure 400 "Block hash was not found"
// @Router /block/hash/{hash} [get]
func (b *BlockHandler) getBlockByHash(c echo.Context) error {
	hash := c.Param("hash")
	if hash == "" {
		b.ctx.Logger().WithError(errors.New("invalid block hash"))
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid block hash"})
	}
	resp, err := b.BUsecase.GetByHash(b.ctx, hash)
	if err != nil {
		b.ctx.Logger().WithError(err)
		return c.JSON(http.StatusNotFound, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: 1})
}

// getBlockByHeight godoc
// @Summary get block by block height
// @Tags block
// @Description Retrieve specific block by height
// @Accept  json
// @Produce  json
// @Param height path string true "Block height"
// @Success 201 {object} domain.Block "Block data"
// @Failure 400 "Block height was not found"
// @Router /block/height/{height} [get]
func (b *BlockHandler) getBlockByHeight(c echo.Context) error {
	height := c.Param("height")
	if height == "" {
		b.ctx.Logger().WithError(errors.New("invalid block height"))
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: "Invalid block height"})
	}
	resp, err := b.BUsecase.GetByHeight(b.ctx, height)
	if err != nil {
		b.ctx.Logger().WithError(err)
		return c.JSON(http.StatusNotFound, Response{Messages: err.Error(), Success: false})
	}
	return c.JSON(http.StatusOK, Response{Data: resp, Success: true, Total: 1})
}
