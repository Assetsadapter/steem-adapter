package types

import (
	"encoding/json"
	"strings"

	"github.com/Assetsadapter/steem-adapter/encoding"
	"github.com/shopspring/decimal"
)

type Amount struct {
	Amount    uint64 `json:"amount"`
	Asset     string `json:"asset"`
	Precision uint8  `json:"precision"`
	Nai       string `json:"nai"`
}

func (aa Amount) Marshal(encoder *encoding.Encoder) error {
	enc := encoding.NewRollingEncoder(encoder)
	enc.EncodeLittleEndianUInt64(aa.Amount)
	enc.EncodeNumber(aa.Precision)
	enc.Encode(FillNai(aa.Nai))
	return enc.Err()
}

func (aa *Amount) UnmarshalJSON(b []byte) (err error) {
	stringCase := struct {
		Amount    string `json:"amount"`
		Precision uint8  `json:"precision"`
		Nai       string `json:"nai"`
	}{}

	if err = json.Unmarshal(b, &stringCase); err == nil {
		spl := strings.Split(stringCase.Amount, " ")
		d, err := ConvertAmountToDecimal(spl[0], 0)
		if err != nil {
			return err
		}
		aa.Precision = stringCase.Precision
		aa.Amount = uint64(d.IntPart())
		aa.Nai = stringCase.Nai
		return err
	}

	return err
}

func ConvertAmountToDecimal(amount string, precision int) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(amount)
	if err != nil {
		return decimal.Zero, err
	}
	d = d.Shift(int32(precision))
	return d, nil
}

// nai 用 x00 补全 7 位
func FillNai(nai string) string {
	var (
		src     = "\x00"
		fillLen = 7 - len(nai)
	)
	filled := nai
	for i := 0; i < fillLen; i++ {
		filled += src
	}
	return filled
}
