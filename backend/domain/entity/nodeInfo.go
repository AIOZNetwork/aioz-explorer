package entity

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
	LastUpdate       int64   `json:"last_update" gorm:"index:idx_node_info_last_update"`
}
