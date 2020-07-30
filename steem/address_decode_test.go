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
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/blocktree/go-owcdrivers/addressEncoder"

	"github.com/Assetsadapter/steem-adapter/addrdec"
)

func TestAddressDecoder_PrivateKeyToWIF(t *testing.T) {
	addrdec.Default.IsTestNet = false
	// exx-exchange private key : 7eeb056aa7801f8c84ebe3dd90891f0d42f4a4451deee67fd8a43d1f9640ec13
	// exx-exchange active role private key : 89569453d9cd2b7649b20f107c8f878cf7a6b5a0812ef0b1e4a42530b55f1a3a
	privKey, _ := hex.DecodeString("d05958af013e4825acf787c5139f4f111e2b1c11840c9ea58490b0a1dd96efd6")
	wif := addressEncoder.AddressEncode(privKey, addrdec.STM_PrivateWIF)
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
