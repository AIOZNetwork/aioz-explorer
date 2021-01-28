package http

import (
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"swagger-server/context"
	"swagger-server/domain"
	"swagger-server/ws"
)

type Response struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	Messages interface{} `json:"messages"`
	Total    int64       `json:"total"`
}

type DeviceHandler struct {
	ctx      context.Context
	DUsecase domain.DeviceUsecase
	wsc      *ws.WSClient
	app      *firebase.App
	RespCh   map[string]*struct {
		SubCh   chan string
		UnsubCh chan string
	}
	MapResponseWalletToken map[string]string
}

func NewDeviceHandler(ctx context.Context, g *echo.Group, a *firebase.App, du domain.DeviceUsecase, wsclient *ws.WSClient) {
	handler := &DeviceHandler{
		ctx:      ctx,
		DUsecase: du,
		wsc:      wsclient,
		app:      a,
		RespCh: make(map[string]*struct {
			SubCh   chan string
			UnsubCh chan string
		}),
		MapResponseWalletToken: make(map[string]string),
	}
	handler.subscribeOnStartup()

	// should clean code
	g.POST("/device/register", handler.registerDevice)
	g.POST("/device/remove", handler.removeDevice)
	//g.POST("/device/enable", handler.enableDevice)
	//g.POST("/device/disable", handler.disableDevice)
}

// registerDevice godoc
// @Summary register new device
// @Tags notification
// @Description register a device to receive notification
// @Accept  json
// @Produce  json
// @Param request body domain.SetPNTokenReq true "Data to create new token"
// @Success 200
// @Failure 400
// @Router /device/register [post]
func (d *DeviceHandler) registerDevice(c echo.Context) error {
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		d.ctx.Logger().WithField("/device/register", err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	var req *domain.SetPNTokenReq
	if err := json.Unmarshal(reqBody, &req); err != nil {
		d.ctx.Logger().WithField("/device/register", err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	d.ctx.Logger().WithField("/api/device/register", req)
	err = d.DUsecase.SaveDevice(d.ctx, req)
	if err != nil {
		d.ctx.Logger().WithField("/device/register", err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	// subscribe and update map token - wallets
	d.subscribeWallets(req.Wallets)
	// replace list wallet subscribed by new one
	d.wsc.MapWalletNotification[req.PnToken] = req.Wallets
	return c.JSON(http.StatusOK, Response{Success: true})
}

// removeDevice godoc
// @Summary remove device token
// @Tags notification
// @Description Remove device from database
// @Accept  json
// @Produce  json
// @Param device body domain.RemovePNTokenReq true "Token key to remove"
// @Success 200
// @Failure 400
// @Router /device/remove [post]
func (d *DeviceHandler) removeDevice(c echo.Context) error {
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		d.ctx.Logger().WithField("/device/remove", err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	var req *domain.RemovePNTokenReq
	if err := json.Unmarshal(reqBody, &req); err != nil {
		d.ctx.Logger().WithField("/device/remove", err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	if err := d.DUsecase.DeleteDevice(d.ctx, req); err != nil {
		d.ctx.Logger().WithField("/device/remove", err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	// delete record in map and unsubscribe
	delete(d.wsc.MapWalletNotification, req.PnToken)
	d.unsubscribeWallets(req.Wallets)
	return c.JSON(http.StatusOK, Response{Success: true})
}

// enableDevice godoc
// @Summary enable device to receive notifications
// @Tags notification
// @Description Enable device to start receiving notification
// @Accept  json
// @Produce  json
// @Param device body string true "Token key to enable"
// @Success 200
// @Failure 400
// @Router /device/enable [post]
func (d *DeviceHandler) enableDevice(c echo.Context) error {
	var req struct {
		Token string `json:"token"`
	}
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		d.ctx.Logger().WithField("/device/enable", err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	if err := json.Unmarshal(reqBody, &req); err != nil {
		d.ctx.Logger().WithField("/device/enable", err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	err = d.DUsecase.UpdateDevice(d.ctx, req.Token, domain.Device_status_active)
	if err != nil {
		d.ctx.Logger().WithField("/device/enable", err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	wallets, err := d.DUsecase.GetWalletsByToken(d.ctx, req.Token)
	if err != nil {
		d.ctx.Logger().WithField("/device/enable", err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	d.subscribeWallets(wallets)
	return c.JSON(http.StatusOK, Response{Success: true})
}

// disableDevice godoc
// @Summary disable device to receive notifications
// @Tags notification
// @Description disable device to stop receiving notification
// @Accept  json
// @Produce  json
// @Param device body string true "Token key to disable"
// @Success 200
// @Failure 400
// @Router /device/disable [post]
func (d *DeviceHandler) disableDevice(c echo.Context) error {
	var req struct {
		Token string `json:"token"`
	}
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		d.ctx.Logger().WithField("/device/disable", err)
		return c.JSON(http.StatusBadRequest, Response{Success: false, Messages: err.Error()})
	}
	if err := json.Unmarshal(reqBody, &req); err != nil {
		d.ctx.Logger().WithField("/device/disable", err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	err = d.DUsecase.UpdateDevice(d.ctx, req.Token, domain.Device_status_deactive)
	if err != nil {
		d.ctx.Logger().WithField("/device/disable", err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	wallets, err := d.DUsecase.GetWalletsByToken(d.ctx, req.Token)
	if err != nil {
		d.ctx.Logger().WithField("/device/disable", err)
		return c.JSON(http.StatusInternalServerError, Response{Success: false, Messages: err.Error()})
	}
	delete(d.wsc.MapWalletNotification, req.Token)
	d.unsubscribeWallets(wallets)
	return c.JSON(http.StatusOK, Response{Success: true})
}

func (d *DeviceHandler) subscribeOnStartup() {
	listDevices, err := d.DUsecase.GetDevices(d.ctx)
	if err != nil {
		d.ctx.Logger().WithField("subscribeOnStartup", err)
	}

	for _, v := range listDevices {
		if _, ok := d.wsc.MapWalletNotification[v.PnToken]; !ok {
			d.wsc.MapWalletNotification[v.PnToken] = make([]string, 0)
			d.wsc.MapWalletNotification[v.PnToken] = append(
				d.wsc.MapWalletNotification[v.PnToken], v.Wallet)
		} else {
			d.wsc.MapWalletNotification[v.PnToken] = append(
				d.wsc.MapWalletNotification[v.PnToken], v.Wallet)
		}
	}
	for _, v := range d.wsc.MapWalletNotification {
		d.subscribeWallets(v)
	}
}

func (d *DeviceHandler) ListenNotificationOnReconnect() {
	var err error
	attempt := 0
	listDevices := make([]*domain.PnTokenDevice, 0)
	for {
		listDevices, err = d.DUsecase.GetDevices(d.ctx)
		if err != nil {
			//todo: retry
			d.ctx.Logger().Error("Failed to get list device FCM successfully")
			attempt++
			continue
		} else {
			d.ctx.Logger().Info("Get list device FCM successfully")
			break
		}
	}

	for _, v := range listDevices {
		if _, ok := d.wsc.MapWalletNotification[v.PnToken]; !ok {
			d.wsc.MapWalletNotification[v.PnToken] = make([]string, 0)
			d.wsc.MapWalletNotification[v.PnToken] = append(
				d.wsc.MapWalletNotification[v.PnToken], v.Wallet)
		} else {
			d.wsc.MapWalletNotification[v.PnToken] = append(
				d.wsc.MapWalletNotification[v.PnToken], v.Wallet)
		}
	}
	for _, v := range d.wsc.MapWalletNotification {
		d.subscribeWallets(v)
	}
}
