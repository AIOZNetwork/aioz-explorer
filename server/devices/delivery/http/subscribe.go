package http

import (
	"context"
	"log"
	"math/rand"
	"swagger-server/ws/message"
	"time"
)

func (d *DeviceHandler) subscribeWallets(wallets []string) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)
	_, err := d.wsc.Send(context.Background(), message.ClientSubscription{
		CID:     -5,
		MsgId:   mid,
		MsgType: "wallet.subscribe",
		MsgData: wallets,
	})
	if err != nil {
		log.Println(err)
	}
}

func (d *DeviceHandler) unsubscribeWallets(wallets []string) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)
	_, err := d.wsc.Send(context.Background(), message.ClientSubscription{
		CID:     -5,
		MsgId:   mid,
		MsgType: "wallet.unsubscribe",
		MsgData: wallets,
	})
	if err != nil {
		log.Println(err)
	}
}
