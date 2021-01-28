package http

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"io/ioutil"
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

type NodeHandler struct {
	ctx         context.Context
	NodeUsecase domain.NodeInfoUsecase
}

func NewNodeHandler(ctx context.Context, g *echo.Group, niu domain.NodeInfoUsecase) {
	handler := &NodeHandler{
		ctx:         ctx,
		NodeUsecase: niu,
	}

	g.GET("/node_info", handler.getNodeInfo)
	g.POST("/node_info/update", handler.updateNodeInfo)
}

// getNodeInfo godoc
// @Summary get node information
// @Tags node
// @Description Retrieve node information
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 201 {object} domain.NodeInfo "Recent blocks"
// @Router /node_info [get]
func (n *NodeHandler) getNodeInfo(c echo.Context) error {
	limit := int64(1000)
	offset := int64(0)
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		limit, _ = strconv.ParseInt(limitStr, 10, 64)
	}
	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		offset, _ = strconv.ParseInt(offsetStr, 10, 64)
	}
	resp, total, err := n.NodeUsecase.GetNodesInfo(limit, offset)
	if err != nil {
		n.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error(), Total: total})
	}
	return c.JSON(http.StatusOK, Response{Success: true, Data: resp})
}

// updateNodeInfo godoc
// @Summary update node with information sent from themselves
// @Tags node
// @Description Update node information
// @Accept  json
// @Produce  json
// @Success 200 "OK"
// @Router /node_info/update [post]
func (n *NodeHandler) updateNodeInfo(c echo.Context) error {
	var req *domain.NodeInfoReq
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		n.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	if err := json.Unmarshal(reqBody, &req); err != nil {
		n.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	if req.NodeId == "" {
		n.ctx.Logger().WithError(errors.New("node ID is empty"))
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: errors.New("node ID is empty")})
	}
	if err := n.NodeUsecase.UpdateNodeById(req); err != nil {
		n.ctx.Logger().WithError(err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	return c.JSON(http.StatusOK, Response{Success: true})
}
