package types

import "testing"

func TestOperations_UnmarshalJSON(t *testing.T) {
	data := `{"operations":[{"type":"transfer_operation","value":{"from":"initminer","to":"leor","amount":{"amount":"1000","precision":3,"nai":"@@000000013"},"memo":"123456789"}}]}`

	ops := Operations{}
	err := ops.UnmarshalJSON(data)
	if err != nil {
		t.Errorf("Opreations unmarshal failed : %s", err.Error())
		t.Fail()
	}
	t.Logf("Opreations unmarshal value : %v", ops)
}

func TestOperations_MarshalJSON(t *testing.T) {
	ops := Operations{
		&TransferOperation{
			Value: Value{
				From: "initminer",
				To:   "exxexchange",
				Amount: Amount{
					Amount:    1000,
					Asset:     "",
					Precision: 3,
					Nai:       "@@000000013",
				},
				Memo: "10000002",
			},
		},
	}
	bytes, err := ops.MarshalJSON()
	if err != nil {
		t.Errorf("Operations marshal failed : %s", err.Error())
		return
	}
	t.Logf("Operations marshal value : %v", string(bytes))
}
