/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package steem

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/blocktree/go-owcrypt"

	"github.com/blocktree/go-owcdrivers/addressEncoder"

	"github.com/Assetsadapter/steem-adapter/addrdec"
)

func TestAddressDecoder_PrivateKeyToWIF(t *testing.T) {
	addrdec.Default.IsTestNet = false
	privKey, _ := hex.DecodeString("9f242fcfaaef51843f960e90e46806d7719a91c38f2c376f85d408a453849438")
	wif := addressEncoder.AddressEncode(privKey, addrdec.STM_mainnetPrivateWIF)
	fmt.Println(wif)
}

func TestAddressDecoder_AddressEncode(t *testing.T) {
	addrdec.Default.IsTestNet = false

	p2pk, _ := hex.DecodeString("033509c7153ab876dd7c305da694b90040e98bda49d58b066f567d766ba059879a")
	p2pkAddr, _ := addrdec.Default.AddressEncode(p2pk, false)
	t.Logf("p2pkAddr: %s", p2pkAddr)
}

func TestAddressDecoder_AddressDecode(t *testing.T) {

	addrdec.Default.IsTestNet = false

	p2pkAddr := "D8qHfnugKAgavULzVjQKyjqxMD7wBETN4s"
	p2pkHash, _ := addrdec.Default.AddressDecode(p2pkAddr)
	t.Logf("p2pkHash: %s", hex.EncodeToString(p2pkHash))

	p2shAddr := "Lb5hzBamSSS4xz2FJqAV2cL2bYMq6oDjJA"

	p2shHash, _ := addrdec.Default.AddressDecode(p2shAddr)
	t.Logf("p2shHash: %s", hex.EncodeToString(p2shHash))
}

func Test_a(t *testing.T) {
	password := "5K2NdiRgYG7FJanCjeszArx7ZKKCZu9kuGStY3a9aZxtnUdthDQ"
	// 需要生成的公钥的角色 生成方式为 账户名称 + 角色 + 密码 = 对应角色私钥
	roleMap := map[string]string{
		"owner":   "",
		"active":  "",
		"posting": "",
		"memo":    "",
	}
	addrdec.Default.IsTestNet = false
	for role, _ := range roleMap {
		// 生成指定角色的私钥
		rolePrivKey := privateKeyFormat("exx-exchange", password, role, t)
		// 使用角色私钥生成对应公钥
		pubKey, ret := owcrypt.GenPubkey(rolePrivKey, owcrypt.ECC_CURVE_SECP256K1)
		if ret != owcrypt.SUCCESS {
			t.Errorf("private key genery public key failed code is : %d", ret)
		}
		// 获取压缩公钥
		compPubKey := owcrypt.PointCompress(pubKey, owcrypt.ECC_CURVE_SECP256K1)
		// 转成 wif 格式的私钥
		key := addressEncoder.AddressEncode(rolePrivKey, addrdec.STM_mainnetPrivateWIF)
		// 生成角色私钥对应的地址
		addrdec.Default.IsTestNet = true
		roleAddr, err := addrdec.Default.AddressEncode(compPubKey)
		if err != nil {
			t.Errorf("encode key address failed : %s", err.Error())
			t.FailNow()
		}
		roleMap[role] = fmt.Sprintf("%s %s", roleAddr, key)
	}

	for k, v := range roleMap {
		fmt.Printf("key : %s, value : %s \n", k, v)
	}
}

// 通过密码生成各种角色的私钥 账户名称+角色+密码
func privateKeyFormat(account, key, role string, t *testing.T) []byte {
	sha := sha256.New()
	if _, err := sha.Write([]byte(account + role + key)); err != nil {
		t.Errorf("private key to role key failed : %s", err.Error())
		t.FailNow()
	}
	return sha.Sum(nil)
}
