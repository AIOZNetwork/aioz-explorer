package ws

import (
	rpchttp "github.com/tendermint/tendermint/rpc/client"
)

func NewRpcConnection(remote, wsEndpoint string) (*rpchttp.HTTP, error) {
	client := rpchttp.NewHTTP(remote, wsEndpoint)
	err := client.Start()
	if err != nil {
		return nil, err
	}
	return client, nil
}
