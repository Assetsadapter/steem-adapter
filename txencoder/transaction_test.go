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
	data := byteTx.Encode()
	t.Logf("Encode raw transaction result is : %s", hex.EncodeToString(*data))
}

func TestTransaction_Decode(t *testing.T) {
	data := "ef6c66f425060b1fec5e010209696e69746d696e65720b65787865786368616e6765e80300000000000003544244000000000831303030303030320000"
	rawTx := RawTransaction{}
	dataByte, err := hex.DecodeString(data)
	if err != nil {
		t.Errorf("Decode data hex failed : %s", err.Error())
	}
	err = rawTx.Decoder(dataByte)
	if err != nil {
		t.Errorf("Decode transaction failed : %s", err.Error())
		t.Fail()
	}
	t.Logf("Decode transaction result : %v", rawTx)
}
