package txsigner

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestSignCanonical(t *testing.T) {
	privKye := "3350b49f9fae9bc43486dcbc112e441bed205dbccb24aea2abbf1c814af6d8b5"
	v := "ethereum"
	sha := sha256.New()
	sha.Write([]byte(v))
	vHash := sha.Sum(nil)
	priKeyBytes, err := hex.DecodeString(privKye)
	if err != nil {
		t.FailNow()
	}
	sig, err := SignCanonical(priKeyBytes, vHash)
	println(hex.EncodeToString(sig))
	if err != nil {
		t.FailNow()
	}
	//"7001514412b7521992e48c4450de8ff51a18e475b06f6b88315e2354b8ba9962 04691c639deebfc869a8995597d7377e169716337bcd8c4bd72c3541011e1ffa"
	//"7912f50819764de81ab7791ab3d62f8dabe84c2fdb2f17d76465d28f8a968f73 55fbb6cd8dfc7545b6258d4b032753b2074232b07f3911822b37f024cd101166 00"
}
