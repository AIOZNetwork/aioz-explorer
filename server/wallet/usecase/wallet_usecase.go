package usecase

import (
	"encoding/json"
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mintkey"
	cmtypes "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"log"
	"sort"
	"strings"
	"swagger-server/context"
	"swagger-server/domain"
)

type walletUsecase struct {
	walletRepo    domain.WalletRepository
	delegatorRepo domain.DelegatorRepository
	validatorRepo domain.ValidatorRepository
	msgsRepo      domain.MsgsRepository
	cdc           *codec.Codec
}

const (
	PrivKeyAminoName = "tendermint/PrivKeySecp256k1"
	PubKeyAminoName  = "tendermint/PubKeySecp256k1"
)

func NewWalletUsecase(wr domain.WalletRepository, dr domain.DelegatorRepository,
	vr domain.ValidatorRepository, mr domain.MsgsRepository, cdc *codec.Codec) domain.WalletUsecase {
	return &walletUsecase{
		walletRepo:    wr,
		delegatorRepo: dr,
		validatorRepo: vr,
		msgsRepo:      mr,
		cdc:           cdc,
	}
}

func (w *walletUsecase) GetValidatorByAddress(ctx context.Context, valAddress string) (*domain.ValidatorResp, error) {
	resp, err := w.validatorRepo.GetValidatorByValAddress(ctx, valAddress)
	if err != nil {
		return nil, err
	}
	dlgs, err := w.delegatorRepo.GetDelegatorByValAddress(ctx, valAddress)
	if err != nil {
		return nil, err
	}
	coins, err := sdk.ParseDecCoins(resp.RewardPool)
	if err != nil {
		return nil, err
	}
	res := &domain.ValidatorResp{
		Address:       resp.Address,
		Tokens:        resp.Tokens,
		Power:         resp.Power,
		Jailed:        resp.Jailed,
		Status:        resp.Status,
		StakedInfo:    dlgs,
		IsActive:      resp.IsActive,
		Detail:        resp.Detail,
		Identity:      resp.Identity,
		Moniker:       resp.Moniker,
		Website:       resp.Website,
		Period:        resp.Period,
		RewardPool:    coins,
		ValConsAddr:   resp.ValConsAddr,
		ValConsPubkey: resp.ValConsPubkey,
	}
	return res, nil
}

func (w *walletUsecase) GetWalletByAddress(ctx context.Context, address string) (*domain.WalletResp, error) {
	resp, err := w.walletRepo.GetWalletByAddress(ctx, address)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return &domain.WalletResp{
				Address:       address,
				AccountNumber: 0,
				Sequence:      0,
				Coins:         sdk.NewCoins(),
			}, nil
		}
		return nil, err
	}
	coin, err := cmtypes.ParseCoins(resp.Coins)
	if err != nil {
		return nil, err
	}
	dlgs, err := w.delegatorRepo.GetDelegatorByAccAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	res := &domain.WalletResp{
		Address:          address,
		Coins:            coin,
		AccountNumber:    resp.AccountNumber,
		Sequence:         resp.Sequence,
		PubKey:           resp.PubKey,
		StakedInfo:       dlgs,
		OriginalVesting:  resp.OriginalVesting,
		DelegatedFree:    resp.DelegatedFree,
		DelegatedVesting: resp.DelegatedVesting,
		StartTime:        resp.StartTime,
		EndTime:          resp.EndTime,
	}

	return res, nil
}

func (w *walletUsecase) GetTxsByAddress(ctx context.Context, address string, limit, offset int) ([]*domain.MessageSend, int64, error) {
	resp, err := w.walletRepo.GetTxsByAddress(ctx, address, limit, offset)
	for _, r := range resp {
		if r.FromAddress == address {
			r.TransactionType = "send"
		} else if r.ToAddress == address {
			r.TransactionType = "receive"
		}
	}
	total, err := w.walletRepo.CountByTableNameWithCondition(ctx, domain.Table_message_send, "from_address = ? OR to_address = ?", address, address)
	if err != nil {
		return nil, -1, err
	}
	return resp, total, err
}

func (w *walletUsecase) GetTxsByAddressV2(ctx context.Context, address string, limit, offset int) ([]*domain.TxsResp, int64, error) {
	resp, err := w.msgsRepo.GetMsgsByAddress(ctx, address, limit, offset)
	total, err := w.walletRepo.CountByTableNameWithCondition(ctx, domain.Table_message_send, "from_address = ? OR to_address = ?", address, address)
	if err != nil {
		return nil, -1, err
	}
	ret := make([]*domain.TxsResp, 0)
	for _, r := range resp {
		var txtype string
		var coins sdk.Coins
		if sliceContains(address, r.AddressFrom) {
			txtype = "send"
		} else if sliceContains(address, r.AddressTo) {
			txtype = "receive"
		}
		err = json.Unmarshal([]byte(r.Amount), &coins)
		if err != nil {
			coins = sdk.NewCoins()
		}
		ret = append(ret, &domain.TxsResp{
			BlockHeight:      r.BlockHeight,
			BlockHash:        r.BlockHash,
			BlockTime:        r.BlockTime,
			TransactionHash:  r.TransactionHash,
			TransactionIndex: r.TransactionIndex,
			TransactionType:  txtype,
			MessageType:      r.MessageType,
			MessageIndex:     r.MessageIndex,
			Address:          r.Address,
			AddressFrom:      r.AddressFrom,
			AddressTo:        r.AddressTo,
			Amount:           coins,
			IsValid:          r.IsValid,
			PayloadMsg:       r.PayloadMsg,
			PayloadErr:       r.PayloadErr,
		})
	}
	return ret, total, err
}

func (w *walletUsecase) GetContactsWallet(ctx context.Context, addresses []string) ([]*domain.ContactsWallet, error) {
	resp := SortedMap{
		M: make(map[string]*struct {
			Send          int64
			Receive       int64
			VolumeSend    sdk.Coins
			VolumeReceive sdk.Coins
		}),
		S: make([]string, 0),
	}
	for _, a := range addresses {
		rowsFrom, err := w.walletRepo.GetContactsByWallet(ctx, a, "from")
		if err != nil {
			return nil, err
		}
		for rowsFrom.Next() {
			var addr string
			var cnt int64
			if err = rowsFrom.Scan(&addr, &cnt); err != nil {
				log.Println(err)
				continue
			}
			volRcv, err := w.GetVolumeFromWalletToWallet(ctx, addr, a)
			if err != nil {
				log.Println(err)
				continue
			}
			if _, ok := resp.M[addr]; !ok {
				resp.M[addr] = &struct {
					Send          int64
					Receive       int64
					VolumeSend    sdk.Coins
					VolumeReceive sdk.Coins
				}{
					Send:          0,
					Receive:       cnt,
					VolumeReceive: volRcv,
				}
			} else {
				resp.M[addr].Receive += cnt
				resp.M[addr].VolumeReceive = resp.M[addr].VolumeReceive.Add(volRcv)
			}
		}

		rowsTo, err := w.walletRepo.GetContactsByWallet(ctx, a, "to")
		if err != nil {
			return nil, err
		}
		for rowsTo.Next() {
			var addr string
			var cnt int64
			if err = rowsTo.Scan(&addr, &cnt); err != nil {
				log.Println(err)
				continue
			}
			volSnd, err := w.GetVolumeFromWalletToWallet(ctx, a, addr)
			if err != nil {
				log.Println(err)
				continue
			}
			if _, ok := resp.M[addr]; !ok {
				resp.M[addr] = &struct {
					Send          int64
					Receive       int64
					VolumeSend    sdk.Coins
					VolumeReceive sdk.Coins
				}{
					Send:       cnt,
					Receive:    0,
					VolumeSend: volSnd,
				}
			} else {
				resp.M[addr].Send += cnt
				resp.M[addr].VolumeSend = resp.M[addr].VolumeSend.Add(volSnd)
			}
		}
	}
	resp.S = sortedKeys(resp.M)
	result := make([]*domain.ContactsWallet, 0)
	for _, s := range resp.S {
		result = append(result, &domain.ContactsWallet{
			Address:        s,
			TotalTxs:       resp.M[s].Send + resp.M[s].Receive,
			AmountSent:     resp.M[s].VolumeSend,
			AmountReceived: resp.M[s].VolumeReceive,
		})
	}
	return result, nil
}

func (w *walletUsecase) CreateWallet(ctx context.Context, req *domain.RequestNewKey) (*domain.KeyResponse, error) {
	if len(req.Password) > 0 && len(req.Password) < domain.MinPassLength {
		return nil, errors.New("password must be at least 8 characters")
	}

	entropySeed, err := bip39.NewEntropy(domain.MnemonicEntropySize)
	if err != nil {
		return nil, err
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed[:])
	if err != nil {
		return nil, err
	}

	seed := bip39.NewSeed(mnemonic, domain.DefaultBIP39Passphrase)
	hdPath := hd.NewFundraiserParams(req.Account, sdk.GetConfig().GetCoinType(), req.Index)

	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hdPath.String())
	if err != nil {
		return nil, err
	}

	privKey := secp256k1.PrivKeySecp256k1(derivedPriv)
	pubKey := privKey.PubKey()

	if len(req.Password) == 0 {
		return &domain.KeyResponse{
			Address:  sdk.AccAddress(pubKey.Address()).String(),
			PubKey:   pubKey.Bytes(),
			PrivKey:  string(w.cdc.MustMarshalJSON(privKey)),
			Mnemonic: mnemonic,
		}, nil
	}

	return &domain.KeyResponse{
		Address:   sdk.AccAddress(pubKey.Address()).String(),
		PubKey:    pubKey.Bytes(),
		PrivKey:   string(w.cdc.MustMarshalJSON(privKey)),
		Mnemonic:  mnemonic,
		PrivArmor: mintkey.EncryptArmorPrivKey(privKey, req.Password),
	}, nil
}

func (w *walletUsecase) RecoverWallet(ctx context.Context, req *domain.RequestRecoverKey) (*domain.KeyResponse, error) {
	if len(req.Password) > 0 && len(req.Password) < domain.MinPassLength {
		return nil, errors.New("password must be at least 8 characters")
	}

	seed, err := bip39.NewSeedWithErrorChecking(req.Mnemonic, keys.DefaultBIP39Passphrase)
	if err != nil {
		return nil, err
	}

	hdPath := hd.NewFundraiserParams(req.Account, sdk.GetConfig().GetCoinType(), req.Index)

	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hdPath.String())
	if err != nil {
		return nil, err
	}

	privKey := secp256k1.PrivKeySecp256k1(derivedPriv)
	pubKey := privKey.PubKey()

	if len(req.Password) == 0 {
		return &domain.KeyResponse{
			Address: sdk.AccAddress(pubKey.Address()).String(),
			PubKey:  pubKey.Bytes(),
			PrivKey: string(w.cdc.MustMarshalJSON(privKey)),
		}, nil
	}

	return &domain.KeyResponse{
		Address:   sdk.AccAddress(pubKey.Address()).String(),
		PubKey:    pubKey.Bytes(),
		PrivKey:   string(w.cdc.MustMarshalJSON(privKey)),
		PrivArmor: mintkey.EncryptArmorPrivKey(privKey, req.Password),
	}, nil
}

func (w *walletUsecase) EncryptKey(ctx context.Context, req *domain.RequestEncrKey) (*domain.KeyResponse, error) {
	if len(req.Password) < domain.MinPassLength {
		return nil, errors.New("password must be at least 8 characters")
	}
	var privKey crypto.PrivKey
	if err := w.cdc.UnmarshalJSON([]byte(req.PrivKey), &privKey); err != nil {
		return nil, err
	}

	pubKey := privKey.PubKey()

	return &domain.KeyResponse{
		Address:   sdk.AccAddress(pubKey.Address()).String(),
		PubKey:    pubKey.Bytes(),
		PrivArmor: mintkey.EncryptArmorPrivKey(privKey, req.Password),
	}, nil
}

func (w *walletUsecase) DecryptKey(ctx context.Context, req *domain.RequestDecrKey) (*domain.KeyResponse, error) {
	if len(req.Password) < domain.MinPassLength {
		return nil, errors.New("password must be at least 8 characters")
	}
	privKey, err := mintkey.UnarmorDecryptPrivKey(req.PrivArmor, req.Password)
	if err != nil {
		return nil, err
	}

	pubKey := privKey.PubKey()

	return &domain.KeyResponse{
		Address: sdk.AccAddress(pubKey.Address()).String(),
		PubKey:  pubKey.Bytes(),
		PrivKey: string(w.cdc.MustMarshalJSON(privKey)),
	}, nil
}

func (w *walletUsecase) GetTopWallet(ctx context.Context, topN int64) ([]*domain.WalletResp, error) {
	resp, err := w.walletRepo.GetTopWallet(ctx, topN)
	res := make([]*domain.WalletResp, 0)
	for _, r := range resp {
		coin, err := cmtypes.ParseCoins(r.Coins)
		if err != nil {
			return nil, err
		}
		res = append(res, &domain.WalletResp{
			Address:       r.Address,
			Coins:         coin,
			AccountNumber: r.AccountNumber,
			Sequence:      r.Sequence,
			PubKey:        r.PubKey,
		})
	}
	return res, err
}

func (w *walletUsecase) CountTotalWallets(ctx context.Context) (int64, error) {
	return w.walletRepo.CountByTableNameWithCondition(ctx, domain.Table_wallet_address, "")
}

func (w *walletUsecase) GetVolumeFromWalletToWallet(ctx context.Context, from, to string) (sdk.Coins, error) {
	volume := sdk.NewCoins()
	msg, err := w.walletRepo.GetVolumeFromWalletToWallet(ctx, from, to)
	if err != nil {
		return nil, err
	}
	for _, a := range msg {
		coins, _ := sdk.ParseCoins(a.Amount)
		volume = volume.Add(coins)
	}
	return volume, nil
}

type SortedMap struct {
	M map[string]*struct {
		Send          int64
		Receive       int64
		VolumeSend    sdk.Coins
		VolumeReceive sdk.Coins
	}
	S []string
}

func (sm *SortedMap) Len() int {
	return len(sm.M)
}

func (sm *SortedMap) Less(i, j int) bool {
	return sm.M[sm.S[i]].Send+sm.M[sm.S[i]].Receive > sm.M[sm.S[j]].Send+sm.M[sm.S[j]].Receive
}

func (sm *SortedMap) Swap(i, j int) {
	sm.S[i], sm.S[j] = sm.S[j], sm.S[i]
}

func sortedKeys(m map[string]*struct {
	Send          int64
	Receive       int64
	VolumeSend    sdk.Coins
	VolumeReceive sdk.Coins
}) []string {
	sm := new(SortedMap)
	sm.M = m
	sm.S = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.S[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.S
}

func sliceContains(input string, array []string) bool {
	for _, a := range array {
		if a == input {
			return true
		}
	}
	return false
}
