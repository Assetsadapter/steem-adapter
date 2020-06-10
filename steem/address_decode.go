/*
 * Copyright 2018 The OpenWallet Authors
 * This file is part of the OpenWallet library.
 *
 * The OpenWallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The OpenWallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package steem

import (
	"github.com/Assetsadapter/steem-adapter/addrdec"
	"github.com/blocktree/go-owcdrivers/addressEncoder"
)

type addressDecoder struct {
	wm *WalletManager //钱包管理者
}

var (
	STM_mainnetPrivateWIFCompressed = addressEncoder.AddressType{"base58", addressEncoder.BTCAlphabet, "doubleSHA256", "", 32, []byte{0x80}, nil}
)

//NewAddressDecoder 地址解析器
func NewAddressDecoder(wm *WalletManager) *addressDecoder {
	decoder := addressDecoder{}
	decoder.wm = wm
	return &decoder
}

//PrivateKeyToWIF 私钥转WIF
func (decoder *addressDecoder) PrivateKeyToWIF(priv []byte, isTestnet bool) (string, error) {
	//var private_key = this.toBuffer();
	//// checksum includes the version
	//private_key = Buffer.concat([new Buffer([0x80]), private_key]);
	//var checksum = hash.sha256(private_key);
	//checksum = hash.sha256(checksum);
	//checksum = checksum.slice(0, 4);
	//var private_wif = Buffer.concat([private_key, checksum]);
	//return base58.encode(private_wif);
	wif := addressEncoder.AddressEncode(priv, STM_mainnetPrivateWIFCompressed)
	return wif, nil
}

//PublicKeyToAddress 公钥转地址
func (decoder *addressDecoder) PublicKeyToAddress(pub []byte, isTestnet bool) (string, error) {
	address, err := addrdec.Default.AddressEncode(pub)
	return address, err
}

//RedeemScriptToAddress 多重签名赎回脚本转地址
func (decoder *addressDecoder) RedeemScriptToAddress(pubs [][]byte, required uint64, isTestnet bool) (string, error) {
	return "", nil
}

//WIFToPrivateKey WIF转私钥
func (decoder *addressDecoder) WIFToPrivateKey(wif string, isTestnet bool) ([]byte, error) {
	priv, err := addrdec.Default.AddressDecode(wif)
	if err != nil {
		return nil, err
	}
	return priv, nil
}
