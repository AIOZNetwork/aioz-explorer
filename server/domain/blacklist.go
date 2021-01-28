package domain

import "swagger-server/context"

type BlacklistUsecase interface {
	UnbanIP(ctx context.Context, ips []string) error
	Add2Blacklist(ctx context.Context, ips []string) error
	Add2Whitelist(ctx context.Context, ips []string) error
	RemoveFromWhitelist(ctx context.Context, ips []string) error
}


type BlacklistRepository interface {
	UnbanIP(ctx context.Context, ips []string) error
	Add2Blacklist(ctx context.Context, ips []string) error
	Add2Whitelist(ctx context.Context, ips []string) error
	RemoveFromWhitelist(ctx context.Context, ips []string) error
}