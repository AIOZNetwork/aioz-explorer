package usecase

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
	"swagger-server/context"
	"swagger-server/domain"
)

type msgsUsecase struct {
	msgsRepo domain.MsgsRepository
}

func NewTxsUsecase(t domain.MsgsRepository) domain.MsgsUsecase {
	return &msgsUsecase{
		msgsRepo: t,
	}
}

func (mu *msgsUsecase) GetMsgSendHistory(ctx context.Context, heightFrom, heightTo, limit, offset int) ([]*domain.MessageSend, int64, error) {
	resp, err := mu.msgsRepo.GetMsgSendHistory(ctx, heightFrom, heightTo, limit, offset)
	if err != nil {
		return nil, -1, err
	}
	total, err := mu.msgsRepo.CountByTableNameWithCondition(ctx, domain.Table_message_send, "block_height >= ? AND block_height <= ?", heightFrom, heightTo)
	if err != nil {
		return nil, -1, err
	}
	return resp, total, nil
}

func (mu *msgsUsecase) GetMsgsByAddress(ctx context.Context, address string, limit, offset int) ([]*domain.Txs, int64, error) {
	resp, err := mu.msgsRepo.GetMsgsByAddress(ctx, address, limit, offset)
	if err != nil {
		return nil, -1, err
	}
	total, err := mu.msgsRepo.CountByTableNameWithCondition(ctx, domain.Table_txs, "address = ?", address)
	if err != nil {
		return nil, -1, err
	}
	return resp, total, nil
}

func (mu *msgsUsecase) GetRecentMsgs(ctx context.Context, limit, offset int) ([]*domain.Txs, int64, error) {
	mgs, err := mu.msgsRepo.GetRecentMsgs(ctx, limit, offset)
	if err != nil {
		return nil, -1, err
	}
	total, err := mu.msgsRepo.CountByTableNameWithCondition(ctx, domain.Table_txs, "")
	if err != nil {
		return nil, -1, err
	}
	return mgs, total, nil
}

func (mu *msgsUsecase) GetMsgsByTxHash(ctx context.Context, txHash string) ([]*domain.Txs, error) {
	return mu.msgsRepo.GetMsgsByTxHash(ctx, txHash)
}

func (mu *msgsUsecase) GetVolumeSend24h(ctx context.Context) (string, error) {
	//result, err := mu.msgsRepo.PaginateVolumeSend24h(ctx)
	result, err := mu.msgsRepo.GetVolumeSend24h(ctx)
	if err != nil {
		return "", err
	}
	total := sdk.NewCoins()
	for _, r := range result {
		var coins sdk.Coins
		err = json.Unmarshal([]byte(r.Amount), &coins)
		//coins, err := sdk.ParseCoins(r.Amount)
		if err != nil {
			log.Println(err)
			continue
		}
		total = total.Add(coins)
	}
	b, err := json.Marshal(total)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
