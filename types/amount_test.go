package types

import (
	"testing"
)

func TestAssetAmount_UnmarshalJSON(t *testing.T) {
	//t.Run("amount uint64", func(t *testing.T) {
	//	data := `{
	//          "amount": "1.000 STEEM",
	//          "precision": 3,
	//		  "nai": "@@000000021"
	//        }`
	//	am := Amount{}
	//	require.NoError(t, json.Unmarshal([]byte(data), &am))
	//
	//	require.Equal(t, uint64(1), am.Amount)
	//	require.Equal(t, "STEEM", am.Asset)
	//	require.Equal(t, uint8(3), am.Precision)
	//	require.Equal(t, "@@000000021", am.Nai)
	//})
	//
	//t.Run("amount string", func(t *testing.T) {
	//	data := `12.000 SBD`
	//	am := Amount{}
	//	require.NoError(t, json.Unmarshal([]byte(data), &am))
	//
	//	require.Equal(t, uint64(12), am.Amount)
	//	require.Equal(t, "SBD", am.Asset)
	//})

	a := Amount{}
	data := `{"amount": "1.000 STEEM","precision": 3,"nai": "@@000000021"}`
	err := a.UnmarshalJSON([]byte(data))
	if err != nil {
		t.Errorf("unmarshal error : %s", err.Error())
		t.Fail()
	}
	t.Logf("Unmarshal result : %v", a)

}
