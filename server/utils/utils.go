package utils

import (
	"fmt"
	cmtypes "github.com/cosmos/cosmos-sdk/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
	"strings"
	"time"
)

func createInt64P(x int64) *int64 {
	return &x
}

func GetTimeFromBlockHeight(c *rpchttp.HTTP, height int64) int64 {
	block, err := c.Block(createInt64P(height))
	if err != nil {
		return time.Now().Unix()
	}
	return block.Block.Time.Unix()
}

func BuildSendQuery(wallet string) string {
	return fmt.Sprintf(`tm.event='Tx' AND message.sender='%v'`, wallet)
}

func BuildReceiveQuery(wallet string) string {
	return fmt.Sprintf(`tm.event='Tx' AND transfer.recipient='%v'`, wallet)
}

func RemoveFromSlice(s []string, t string) []string {
	for i, e := range s {
		if e == t {
			s[i] = s[len(s)-1]
			// We do not need to put s[i] at the end, as it will be discarded anyway
			return s[:len(s)-1]
		}
	}
	return s
}

func AddToSlice(s []string, t string) []string {
	if !CheckExistsInSlice(s, t) {
		s = append(s, t)
		return s
	} else {
		return s
	}
}

func CheckExistsInSlice(clice []string, elem string) bool {
	for _, v := range clice {
		if elem == v {
			return true
		}
	}
	return false
}

func CrossCheckSliceInSlice(source, target []string) bool {
	for _, i := range source {
		for _, j := range target {
			if i == j {
				return true
			}
		}
	}
	return false
}

func RemoveTrailingZerosFromDecCoins(dc cmtypes.DecCoins) string {
	out := ""
	for _, c := range dc {
		tokens := strings.Split(c.Amount.String(), ".")
		if len(tokens) < 2 {
			out += fmt.Sprintf("%v,", c.String())
			continue
		}
		s := tokens[1]
		for i := len(s) - 1; i >= 0; i-- {
			if s[i] == '0' {
				s = s[:i]
			} else {
				break
			}
		}
		if s == "" {
			out += tokens[0] + " " + c.Denom + ","
		} else {
			out += strings.Join([]string{tokens[0], s}, ".") + " " + c.Denom + ","
		}
	}
	resp := out[:len(out)-1]
	return resp
}

func RemoveTrailingZerosFromDec(d string) string {
	tokens := strings.Split(d, ".")
	if len(tokens) < 2 {
		return tokens[0]
	}
	s := tokens[1]
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '0' {
			s = s[:i]
		} else {
			break
		}
	}
	if s == "" {
		return tokens[0]
	} else {
		return strings.Join([]string{tokens[0], s}, ".")
	}
}
