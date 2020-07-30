package addrdec

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/go-owcrypt"
)

func TestAddressDecoderV2_AddressDecode(t *testing.T) {
	type fields struct {
		IsTestNet bool
	}
	type args struct {
		addr string
		opts []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "bts bech32", fields: fields{IsTestNet: false},
			args:    args{addr: "BTS4txUYJW7JwMhYYP3uePSfj4XfHqUnbLdzpPmGNQFHvaNjw4Ps2"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := &AddressDecoderV2{
				IsTestNet: tt.fields.IsTestNet,
			}
			got, err := dec.AddressDecode(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddressDecoderV2.AddressDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				// t.Errorf("AddressDecoderV2.AddressDecode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateAccountRolePrivateKey(t *testing.T) {
	//tagetActivePrivateKey := "bebfaedefcca779bcbc042b37f6212988c1e1fd1e38ef9e6f8cf061065967e64"
	accountName := "exx-withdraw"
	password := "5KQ3aSN53ci5QQEn8ebm8MK9fQ6ra5JZ6VaaBrf1c8wu7FR3QKX"
	rolesMap := map[string]string{
		"owner":   "",
		"active":  "",
		"posting": "",
		"memo":    "",
	}
	Default.IsTestNet = true
	for k, _ := range rolesMap {
		priKey, err := CalculateAccountRolePrivateKey(accountName, k, password)
		if err != nil {
			t.Errorf("calculateAccountRolePrivateKey failed : %s", err.Error())
			t.FailNow()
		}
		wifPriKey := addressEncoder.AddressEncode(priKey, STM_PrivateWIF)
		comPriKey, err := GetCompPubKey(priKey, owcrypt.ECC_CURVE_SECP256K1)
		if err != nil {
			t.Errorf("GetCompPubKey failed : %s", err.Error())
			t.FailNow()
		}
		address, err := Default.AddressEncode(comPriKey)
		if err != nil {
			t.Errorf("Address decode failed : %s", err.Error())
			t.FailNow()
		}
		rolesMap[k] = fmt.Sprintf("%s %s", wifPriKey, address)
	}
	for k, v := range rolesMap {
		fmt.Printf("role: %s, value: %s \n", k, v)
	}
}
