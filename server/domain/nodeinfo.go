package domain

import "encoding/json"

type NodeInfoReq struct {
	HardwareInfo   interface{} `json:"hardware_info"`
	IP             string      `json:"ip"`
	Latitude       float64     `json:"latitude"`
	Longitude      float64     `json:"longitude"`
	Location       string      `json:"location"`
	NodeId         string      `gorm:"primary_key" json:"node_id"`
	NodeType       string      `json:"node_type"`
	WorkerEndpoint string      `json:"worker_endpoint"`
	WorkerRegistered bool    `json:"worker_registered"`
	WorkerAddress  string      `json:"worker_address"`
	WorkerPubkey   string      `json:"worker_pubkey"`
	ValConsAddress string      `json:"val_cons_address"`
	LastUpdate     int64       `json:"last_update"`
}

type NodeInfo struct {
	HardwareInfo     string  `json:"hardware_info"`
	IP               string  `json:"ip"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Location         string  `json:"location"`
	NodeId           string  `gorm:"primary_key" json:"node_id"`
	NodeType         string  `json:"node_type"`
	WorkerEndpoint   string  `json:"worker_endpoint"`
	WorkerRegistered bool    `json:"worker_registered"`
	WorkerAddress    string  `json:"worker_address"`
	WorkerPubkey     string  `json:"worker_pubkey"`
	ValConsAddress   string  `json:"val_cons_address"`
	LastUpdate       int64   `json:"last_update"`
	IsValidator      bool    `json:"is_validator"`
	IsOnline         bool    `json:"is_online"`
}

type NodeInfoResponse struct {
	HardwareInfo   json.RawMessage `json:"hardware_info"`
	IP             string          `json:"ip"`
	Latitude       float64         `json:"latitude"`
	Longitude      float64         `json:"longitude"`
	Location       string          `json:"location"`
	NodeId         string          `json:"node_id"`
	NodeType       string          `json:"node_type"`
	WorkerEndpoint string          `json:"worker_endpoint"`
	WorkerRegistered bool    `json:"worker_registered"`
	WorkerAddress  string          `json:"worker_address"`
	WorkerPubkey   string          `json:"worker_pubkey"`
	ValConsAddress string          `json:"val_cons_address"`
	LastUpdate     int64           `json:"last_update"`
	IsValidator    bool            `json:"is_validator"`
	IsOnline       bool            `json:"is_online"`
}

type NodeInfoRepository interface {
	GetNodesInfo(limit, offset int64) ([]*NodeInfo, error)
	UpdateNodeById(body *NodeInfo) error
	DeleteNodeById(nodeId string) error
	CountTotalNodes() (int64, error)
}

type NodeInfoUsecase interface {
	GetNodesInfo(limit, offset int64) ([]*NodeInfoResponse, int64, error)
	UpdateNodeById(body *NodeInfoReq) error
	DeleteNodeById(nodeId string) error
}
