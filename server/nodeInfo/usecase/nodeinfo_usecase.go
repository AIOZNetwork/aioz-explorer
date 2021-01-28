package usecase

import (
	"encoding/json"
	"log"
	"swagger-server/domain"
	"time"
)

type nodeInfoUsecase struct {
	nodeInfoRepo domain.NodeInfoRepository
}

func NewNodeInfoUsecase(niu domain.NodeInfoRepository) domain.NodeInfoUsecase {
	return &nodeInfoUsecase{
		nodeInfoRepo: niu,
	}
}

func (n *nodeInfoUsecase) GetNodesInfo(limit, offset int64) ([]*domain.NodeInfoResponse, int64, error) {
	resp, err := n.nodeInfoRepo.GetNodesInfo(limit, offset)
	if err != nil {
		return nil, -1, err
	}
	total, err := n.nodeInfoRepo.CountTotalNodes()
	if err != nil {
		return nil, -1, err
	}
	result := make([]*domain.NodeInfoResponse, 0)
	for _, v := range resp {
		var hwInfo json.RawMessage
		err = json.Unmarshal([]byte(v.HardwareInfo), &hwInfo)
		result = append(result, &domain.NodeInfoResponse{
			HardwareInfo:     hwInfo,
			IP:               v.IP,
			Latitude:         v.Latitude,
			Longitude:        v.Longitude,
			Location:         v.Location,
			NodeId:           v.NodeId,
			NodeType:         v.NodeType,
			WorkerEndpoint:   v.WorkerEndpoint,
			WorkerRegistered: v.WorkerRegistered,
			WorkerAddress:    v.WorkerAddress,
			WorkerPubkey:     v.WorkerPubkey,
			ValConsAddress:   v.ValConsAddress,
			LastUpdate:       v.LastUpdate,
			IsValidator:      v.IsValidator,
			IsOnline:         v.IsOnline,
		})
	}
	return result, total, nil
}

func (n *nodeInfoUsecase) UpdateNodeById(body *domain.NodeInfoReq) error {
	hwInfo, _ := json.Marshal(body.HardwareInfo)
	req := &domain.NodeInfo{
		HardwareInfo:     string(hwInfo),
		IP:               body.IP,
		Latitude:         body.Latitude,
		Longitude:        body.Longitude,
		Location:         body.Location,
		NodeId:           body.NodeId,
		NodeType:         body.NodeType,
		WorkerEndpoint:   body.WorkerEndpoint,
		WorkerRegistered: body.WorkerRegistered,
		WorkerAddress:    body.WorkerAddress,
		WorkerPubkey:     body.WorkerPubkey,
		ValConsAddress:   body.ValConsAddress,
		LastUpdate:       time.Now().Unix(),
	}
	return n.nodeInfoRepo.UpdateNodeById(req)
}

func (n *nodeInfoUsecase) DeleteNodeById(nodeId string) error {
	return n.nodeInfoRepo.DeleteNodeById(nodeId)
}
