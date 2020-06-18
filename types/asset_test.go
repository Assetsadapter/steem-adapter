package types

import (
	"testing"
)

func TestIsSupportSymbol(t *testing.T) {
	v := "HIVE"
	if IsSupportSymbol(v) {
		t.Log("supported token!")
		return
	}
	t.Error("not supported token!")
}
