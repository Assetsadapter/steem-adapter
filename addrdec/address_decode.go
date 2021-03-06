package addrdec

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/openwallet"

	"github.com/blocktree/go-owcdrivers/addressEncoder"
)

var (
	STMMainNetPublicKeyPrefixCompat = "STM"
	STMTestNetPublicKeyPrefixCompat = "TST"

	//STM stuff
	STM_mainnetPublic        = addressEncoder.AddressType{"base58", addressEncoder.BTCAlphabet, "ripemd160", "", 33, []byte(STMMainNetPublicKeyPrefixCompat), nil}
	STM_testnetPublic        = addressEncoder.AddressType{"base58", addressEncoder.BTCAlphabet, "ripemd160", "", 33, []byte(STMTestNetPublicKeyPrefixCompat), nil}
	STM_PrivateWIF           = addressEncoder.AddressType{"base58", addressEncoder.BTCAlphabet, "doubleSHA256", "", 32, []byte{0x80}, nil}
	STM_PrivateWIFCompressed = addressEncoder.AddressType{"base58", addressEncoder.BTCAlphabet, "doubleSHA256", "", 32, []byte{0x80}, []byte{0x01}}

	Default = AddressDecoderV2{}
)

//AddressDecoderV2
type AddressDecoderV2 struct {
	*openwallet.AddressDecoderV2Base
	IsTestNet bool
}

//NewAddressDecoder 地址解析器
func NewAddressDecoderV2(isTestNet bool) *AddressDecoderV2 {
	decoder := AddressDecoderV2{
		IsTestNet:isTestNet,
	}
	return &decoder
}

// AddressDecode decode address
func (dec *AddressDecoderV2) AddressDecode(pubKey string, opts ...interface{}) ([]byte, error) {
	var pubKeyMaterial string
	if strings.HasPrefix(pubKey, STMMainNetPublicKeyPrefixCompat) { // "STM"
		pubKeyMaterial = pubKey[len(STMMainNetPublicKeyPrefixCompat):] // strip "STM"
	} else if strings.HasPrefix(pubKey, STMTestNetPublicKeyPrefixCompat) { // "TST"
		pubKeyMaterial = pubKey[len(STMTestNetPublicKeyPrefixCompat):] // strip "TST"
	} else {
		return nil, fmt.Errorf("public key should start with [%q]", STMMainNetPublicKeyPrefixCompat)
	}
	ret, err := addressEncoder.Base58Decode(pubKeyMaterial, addressEncoder.NewBase58Alphabet(STM_mainnetPublic.Alphabet))
	if err != nil {
		return nil, addressEncoder.ErrorInvalidAddress
	}
	if addressEncoder.VerifyChecksum(ret, STM_mainnetPublic.ChecksumType) == false {
		return nil, addressEncoder.ErrorInvalidAddress
	}
	return ret[:len(ret)-4], nil
}

// AddressEncode encode address
func (dec *AddressDecoderV2) AddressEncode(hash []byte, opts ...interface{}) (string, error) {
	pubType := STM_mainnetPublic
	if dec.IsTestNet {
		pubType = STM_testnetPublic
	}
	data := addressEncoder.CatData(hash, addressEncoder.CalcChecksum(hash, pubType.ChecksumType))
	return string(pubType.Prefix) + addressEncoder.EncodeData(data, pubType.EncodeType, pubType.Alphabet), nil
}

func GetRoleCompressKey(accountName, role, password string, curveType uint32) ([]byte, error) {
	priKey, err := CalculateAccountRolePrivateKey(accountName, role, password)
	if err != nil {
		return nil, err
	}
	comPriKey, err := GetCompPubKey(priKey, curveType)
	if err != nil {
		return nil, err
	}
	return comPriKey, nil
}

// 计算角色私钥
func CalculateAccountRolePrivateKey(accountName, role, password string) ([]byte, error) {
	if 0 == len(accountName) || 0 == len(role) || 0 == len(password) {
		return nil, fmt.Errorf("invalied args")
	}
	priKey := sha256.Sum256([]byte(accountName + role + password))
	return priKey[:], nil
}

// 获取角色的压缩公钥
func GetCompPubKey(priKey []byte, curveType uint32) ([]byte, error) {
	pubKye, ret := owcrypt.GenPubkey(priKey, curveType)
	if ret != owcrypt.SUCCESS {
		return nil, fmt.Errorf("GenPubKey failed code is : %d", ret)
	}
	compPubKey := owcrypt.PointCompress(pubKye, curveType)
	return compPubKey, nil
}
