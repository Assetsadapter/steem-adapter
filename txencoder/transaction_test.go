package txencoder

import (
	"encoding/hex"
	"testing"
	"time"
)

func TestTransaction_Encode(t *testing.T) {
	rawTx := RawTransaction{
		RefBlockNum:    486639 & 0xFFFF,
		RefBlockPrefix: "66f42506",
		Expiration:     time.Now(),
		Operations: &[]RawTransferOperation{
			{
				Type: 2,
				From: "initminer",
				To:   "exxexchange",
				Amount: RawAmount{
					Amount:    1000,
					Precision: 3,
					Nai:       "TBD",
				},
				Memo: "10000002",
			},
		},
	}
	byteTx, err := rawTx.Encode()
	if err != nil {
		t.Errorf("Encode raw transaction failed : %s", err.Error())
		t.Fail()
	}
	data := byteTx.Decode()
	t.Logf("Encode raw transaction result is : %s", hex.EncodeToString(*data))
}
