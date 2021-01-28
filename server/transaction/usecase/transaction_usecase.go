package usecase

import (
	"swagger-server/context"
	"swagger-server/domain"
)

type transactionUsecase struct {
	transactionRepo domain.TransactionRepository
	msgsRepo        domain.MsgsRepository
}

func NewTransactionUsecase(t domain.TransactionRepository, m domain.MsgsRepository) domain.TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: t,
		msgsRepo:        m,
	}
}

func (t *transactionUsecase) GetTransaction(ctx context.Context, hash string) (*domain.TransactionResp, error) {
	tx, err := t.transactionRepo.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}
	msgs, err := t.msgsRepo.GetMsgsByTxHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	resp := &domain.TransactionResp{
		BlockHeight: tx.BlockHeight,
		BlockTime:   tx.BlockTime,
		Hash:        tx.Hash,
		TxIndex:     tx.TxIndex,
		Gas:         tx.Gas,
		Memo:        tx.Memo,
		Fee:         tx.Fee,
		Signatures:  tx.Signatures,
		IsValid:     tx.IsValid,
		Payload:     msgs,
	}
	return resp, nil
}

func (t *transactionUsecase) GetTransactions(ctx context.Context, limit, offset int64) ([]*domain.TransactionResp, int64, error) {
	txs, err := t.transactionRepo.GetTransactions(ctx, limit, offset)
	if err != nil {
		return nil, -1, err
	}
	resp := make([]*domain.TransactionResp, 0)
	for _, tx := range txs {
		msgs, err := t.msgsRepo.GetMsgsByTxHash(ctx, tx.Hash)
		if err != nil {
			return nil, -1, err
		}
		resp = append(resp, &domain.TransactionResp{
			BlockHeight: tx.BlockHeight,
			BlockTime:   tx.BlockTime,
			Hash:        tx.Hash,
			TxIndex:     tx.TxIndex,
			Gas:         tx.Gas,
			Memo:        tx.Memo,
			Fee:         tx.Fee,
			Signatures:  tx.Signatures,
			IsValid:     tx.IsValid,
			Payload:     msgs,
		})
	}

	total, err := t.transactionRepo.CountWithCondition(ctx, "")
	if err != nil {
		return nil, -1, err
	}
	return resp, total, nil
}

func (t *transactionUsecase) GetTransactionsByBlockHeight(ctx context.Context, height, limit, offset int64) ([]*domain.TransactionResp, int64, error) {
	txs, err := t.transactionRepo.GetTransactionsByBlockHeight(ctx, height, limit, offset)
	if err != nil {
		return nil, -1, err
	}
	resp := make([]*domain.TransactionResp, 0)
	for _, tx := range txs {
		msgs, err := t.msgsRepo.GetMsgsByTxHash(ctx, tx.Hash)
		if err != nil {
			return nil, -1, err
		}
		resp = append(resp, &domain.TransactionResp{
			BlockHeight: tx.BlockHeight,
			BlockTime:   tx.BlockTime,
			Hash:        tx.Hash,
			TxIndex:     tx.TxIndex,
			Gas:         tx.Gas,
			Memo:        tx.Memo,
			Fee:         tx.Fee,
			Signatures:  tx.Signatures,
			IsValid:     tx.IsValid,
			Payload:     msgs,
		})
	}
	total, err := t.transactionRepo.CountWithCondition(ctx, "block_height = ?", height)
	if err != nil {
		return nil, -1, err
	}
	return resp, total, nil
}

func (t *transactionUsecase) CountTotalTx(ctx context.Context) (int64, error) {
	return t.transactionRepo.CountWithCondition(ctx, "")
}
