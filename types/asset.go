package types

import (
	"errors"
	"fmt"
	"strings"
)

type Asset struct {
	Amount    string `json:"amount"`
	Symbol    string `json:"symbol"`
	Nai       string `json:"nai"`
	Precision uint   `json:"precision"`
}

var (
	// 资金的单位，包括测试链与主链
	supportSymbol = map[string]string{
		"TESTS": "",
		"TBD":   "",
		"STEEM": "",
		"SBD":   "",
		"HBD":   "",
		"HIVE":  "",
		"VESTS": "",
	}

	// 三种币的标识符
	symbolToken = map[string]string{
		"STM":   "@@000000021",
		"SBD":   "@@000000013",
		"VESTS": "@@000000037",
	}
)

func UnMarshal(amount string) (*Asset, error) {
	spl := strings.Split(amount, " ")
	if 2 != len(spl) {
		return nil, errors.New(fmt.Sprintf("invalid amount=%s", amount))
	}
	if IsSupportSymbol(spl[1]) {
		return nil, errors.New(fmt.Sprint("not support symbol"))
	}
	asset := &Asset{
		Amount: spl[0],
		Symbol: spl[1],
	}
	asset.GetPrecision()
	asset.GetNai()
	return asset, nil
}

func Marshal(asset *Asset) string {
	return fmt.Sprintf("%s %s", asset.Amount, asset.Symbol)
}

func (a *Asset) GetPrecision() {
	switch a.Symbol {
	case "TESTS":
		fallthrough
	case "HIVE":
		fallthrough
	case "STEEM":
		fallthrough
	case "TBD":
		fallthrough
	case "SBD":
		fallthrough
	case "HBD":
		a.Precision = 3
	case "VESTS":
		a.Precision = 6
	}
}

func (a *Asset) GetNai() {
	switch a.Symbol {
	case "TESTS":
		fallthrough
	case "HIVE":
		fallthrough
	case "STEEM":
		a.Nai = symbolToken["STM"]
	case "TBD":
		fallthrough
	case "SBD":
		fallthrough
	case "HBD":
		a.Nai = symbolToken["SBD"]
	case "VESTS":
		a.Nai = symbolToken["VESTS"]
	}
}

func IsSupportSymbol(symbol string) bool {
	_, ok := supportSymbol[symbol]
	return ok
}
