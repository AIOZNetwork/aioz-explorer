package usecase

import (
	"swagger-server/context"
	"swagger-server/domain"
)

type blacklistUsecase struct {
	blacklistRepo domain.BlacklistRepository
}

func NewBlacklistUsecase(bl domain.BlacklistRepository) domain.BlacklistUsecase {
	return &blacklistUsecase{
		blacklistRepo: bl,
	}
}

func (b *blacklistUsecase) Add2Blacklist(ctx context.Context, ips []string) error {
	return b.blacklistRepo.Add2Blacklist(ctx, ips)
}

func (b *blacklistUsecase) UnbanIP(ctx context.Context, ips []string) error {
	return b.blacklistRepo.UnbanIP(ctx, ips)
}

func (b *blacklistUsecase) Add2Whitelist(ctx context.Context, ips []string) error {
	return b.blacklistRepo.Add2Whitelist(ctx, ips)
}

func (b *blacklistUsecase) RemoveFromWhitelist(ctx context.Context, ips []string) error {
	return b.blacklistRepo.RemoveFromWhitelist(ctx, ips)
}
