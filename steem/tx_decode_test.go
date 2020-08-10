package steem

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/Assetsadapter/steem-adapter/addrdec"
	"github.com/Assetsadapter/steem-adapter/txsigner"

	"github.com/Assetsadapter/steem-adapter/txencoder"
)

func TestSha256ChainId(t *testing.T) {
	chainId := "0000000000000000000000000000000000000000000000000000000000000000"
	sha := sha256.New()
	chainIdBytes, err := hex.DecodeString(chainId)
	if err != nil {
		t.Errorf("ChainId parse to hex failed : %s", err.Error())
		t.FailNow()
	}
	_, err = sha.Write(chainIdBytes)
	if err != nil {
		t.Errorf("Write chainId to sha256 converter failed : %s", err.Error())
	}
	result := sha.Sum(nil)
	t.Logf("ChainId sha256 result : %s", hex.EncodeToString(result))
}

func TestDigest(t *testing.T) {
	tv, err := time.Parse("2006-01-02T15:04:05", "2020-07-23T03:25:42")
	if err != nil {
		t.Errorf("failed : %s", err.Error())
		t.FailNow()
	}
	println(fmt.Sprintf("%d", tv.Unix()))
}

func TestTransactionDecoder_SignRawTransaction(t *testing.T) {
	tv, err := time.Parse("2006-01-02T15:04:05", "2020-07-28T07:26:00")
	if err != nil {
		t.Errorf("failed : %s", err.Error())
		t.FailNow()
	}

	op := txencoder.RawTransferOperation{
		Type: 2,
		From: "exx-withdraw",
		To:   "leor",
		Amount: txencoder.RawAmount{
			Amount:    1000,
			Precision: 3,
			Nai:       "TESTS",
		},
		Memo: "",
	}
	rawTx := txencoder.RawTransaction{
		RefBlockNum:    16645,
		RefBlockPrefix: 0xc97878c3,
		Expiration:     tv,
		Operations:     &[]txencoder.RawOperation{&op},
		Extensions:     &[]txencoder.Extension{},
		Signatures:     &[]string{},
	}
	tx, err := rawTx.Encode()
	if err != nil {
		t.Errorf("failed : %s", err.Error())
		t.FailNow()
	}
	digest, err := tx.Digest("08c5839c0f1c1a0acae7f2e33978a21168b2c1b5f78059f902bc0c3977fff163")
	if err != nil {
		t.Errorf("failed : %s", err.Error())
		t.FailNow()
	}

	rolePriKey, err := addrdec.CalculateAccountRolePrivateKey("exx-withdraw", "active", "5KQ3aSN53ci5QQEn8ebm8MK9fQ6ra5JZ6VaaBrf1c8wu7FR3QKX")
	if err != nil {
		t.Errorf("failed : %s", err.Error())
		t.FailNow()
	}

	sigened, err := txsigner.SignCanonical(rolePriKey, digest)
	if err != nil {
		t.Errorf("failed : %s", err.Error())
		t.FailNow()
	}

	fmt.Println(hex.EncodeToString(sigened))
}
