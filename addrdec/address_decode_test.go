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
	accountName := "exx-exchange"
	password := "5JnBXrZRh3KCemXMBGxToPViZPLumA1os8Cnx4hepLHLJtrrkGS"
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
		wifPriKey := addressEncoder.AddressEncode(priKey, STM_mainnetPrivateWIF)
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
