module swagger-server

go 1.12

require (
	aioz.io/go-aioz v0.0.0-00010101000000-000000000000
	firebase.google.com/go/v4 v4.0.0
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/cosmos/cosmos-sdk v0.39.1
	github.com/cosmos/go-bip39 v0.0.0-20200817134856-d632e0d11689
	github.com/didip/tollbooth v4.0.2+incompatible
	github.com/didip/tollbooth_echo v0.0.0-20201202024403-6a73caa03064
	github.com/go-openapi/spec v0.19.9 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/labstack/echo/v4 v4.1.17
	github.com/lib/pq v1.3.0
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rcrowley/go-metrics v0.0.0-20190704165056-9c2d0518ed81
	github.com/rgamba/evtwebsocket v0.0.0-20181029234908-48b8cd9f8616 // indirect
	github.com/rs/zerolog v1.15.0
	github.com/sacOO7/go-logger v0.0.0-20180719173527-9ac9add5a50d // indirect
	github.com/sacOO7/gowebsocket v0.0.0-20180719182212-1436bb906a4e // indirect
	github.com/sacOO7/socketcluster-client-go v1.0.0
	github.com/shopspring/decimal v0.0.0-20200227202807-02e2044944cc
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/assertions v1.1.1 // indirect
	github.com/spf13/viper v1.7.1
	github.com/swaggo/echo-swagger v1.0.0
	github.com/swaggo/swag v1.6.7
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.8
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 // indirect
	golang.org/x/net v0.0.0-20201006153459-a7d1128ccaa0 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20201006155630-ac719f4daadf // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools v0.0.0-20201007032633-0806396f153e // indirect
	google.golang.org/api v0.33.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gorm.io/driver/postgres v1.0.2
	gorm.io/gorm v1.20.8
	gorm.io/hints v0.0.0-20201206014637-b313d3d1dc5e
)

replace (
	aioz.io/crypto => 10.0.0.50/aioz-network/crypto.git v0.0.0-20191109082407-2b3bd2b9fc42
	aioz.io/go-aioz => 10.0.0.50/aioz-network/go-aioz.git v0.0.0-20201019083412-b6e9231df1f6
	github.com/cosmos/cosmos-sdk => 10.0.0.50/aioz-network/cosmos-sdk.git v0.0.0-20200925040103-f05ac1cf00f5
	github.com/tendermint/go-amino => 10.0.0.50/aioz-network/go-amino.git v0.0.0-20191217044337-b3e919f22633
	github.com/tendermint/iavl => 10.0.0.50/aioz-network/iavl.git v0.0.0-20191204024104-d91b764d4bd5
	github.com/tendermint/tendermint => 10.0.0.50/aioz-network/tendermint.git v0.0.0-20200925040514-ae8ab33b094d
)
