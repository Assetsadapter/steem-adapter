package types

import (
	"encoding/json"
	"strconv"

	"github.com/Assetsadapter/bitshares-adapter/encoding"
)

type Price struct {
	Base  AssetAmount `json:"base"`
	Quote AssetAmount `json:"quote"`
}

type AssetAmount struct {
	Amount    uint64 `json:"amount"`
	Precision uint   `json:"precision"`
	Nai       string `json:"nai"`
}

func (aa AssetAmount) Marshal(encoder *encoding.Encoder) error {
	enc := encoding.NewRollingEncoder(encoder)
	enc.EncodeLittleEndianUInt64(aa.Amount)
	enc.Encode(aa.Nai)
	return enc.Err()
}

// RPC client might return asset amount as uint64 or string,
// therefore a custom unmarshaller is used
func (aa *AssetAmount) UnmarshalJSON(b []byte) (err error) {
	stringCase := struct {
		Amount  string   `json:"amount"`
		AssetID ObjectID `json:"asset_id"`
	}{}

	uint64Case := struct {
		Amount  uint64   `json:"amount"`
		AssetID ObjectID `json:"asset_id"`
	}{}

	if err = json.Unmarshal(b, &uint64Case); err == nil {
		aa.Amount = uint64Case.Amount
		return nil
	}

	// failed on uint64, try string
	if err = json.Unmarshal(b, &stringCase); err == nil {
		aa.Amount, err = strconv.ParseUint(stringCase.Amount, 10, 64)
		return err
	}

	return err
}
