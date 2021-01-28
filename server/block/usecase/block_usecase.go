package usecase

import (
	"swagger-server/context"
	"swagger-server/domain"
)

type blockUsecase struct {
	blockRepo domain.BlockRepository
}

func NewBlockUsecase(b domain.BlockRepository) domain.BlockUsecase {
	return &blockUsecase{
		blockRepo: b,
	}
}

func (b *blockUsecase) GetByHeight(ctx context.Context, id string) (*domain.Block, error) {
	return b.blockRepo.GetByHeight(ctx, id)
}

func (b *blockUsecase) GetLatestBlocks(ctx context.Context, from, size int) ([]*domain.Block, int64, error) {
	blocks, err := b.blockRepo.GetLatestBlocks(ctx, from, size)
	if err != nil {
		return nil, -1, err
	}
	//total, err := b.blockRepo.CountWithCondition(ctx, "")
	//if err != nil {
	//	return nil, -1, err
	//}
	var total int64
	if len(blocks) == 0 {
		total = 0
	} else {
		total = blocks[0].Height
	}
	return blocks, total, nil
}

func (b *blockUsecase) GetByHash(ctx context.Context, hash string) (*domain.Block, error) {
	return b.blockRepo.GetByHash(ctx, hash)
}

func (b *blockUsecase) CountTotalBlocks(ctx context.Context) (int64, error) {
	return b.blockRepo.CountWithCondition(ctx, "")
}

func (b *blockUsecase) GetBestBlock(ctx context.Context) (*domain.Block, error) {
	return b.blockRepo.GetBestBlock(ctx)
}
