package txencoder

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestTransaction_Encode(t *testing.T) {
	refBlockId := "0007db3dfa79a5e7f81f3f620517db1dfd13656c"
	rawTx := RawTransaction{
		RefBlockNum:    514877 & 0xFFFF,
		RefBlockPrefix: refBlockId[8:16],
		Expiration:     time.Now(),
		Operations: &[]RawOperation{
			&RawTransferOperation{
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
	data := "ef6c66f425061d7cec5e010209696e69746d696e65720b65787865786368616e6765e80300000000000003544244000000000831303030303030320000"
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

func TestParseRefBlockPrefixToNumber(t *testing.T) {
	refBlockPrefix := "f7860b82"
	revByte, err := reverseHexToBytes(refBlockPrefix)
	if err != nil {
		t.Errorf("reverseHexToBytes failed : %s", err.Error())
		t.FailNow()
	}
	prefix, err := strconv.ParseUint(hex.EncodeToString(revByte), 16, 64)
	if err != nil {
		t.Errorf("Parse uint failed : %s", err.Error())
		t.FailNow()
	}
	fmt.Printf("refBlockPrefix number is : %d \n", prefix)
}
