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

var (
	tw *WalletManager
)

func init() {
	tw = testNewWalletManager()
}

func testNewWalletManager() *WalletManager {
	wm := NewWalletManager(nil)
	//wm.Config.ServerAPI = "http://api.bts.ai/rpc"
	wm.Config.ServerAPI =  "http://127.0.0.1:9876"
	//wm.Config.WalletAPI= "http://1.wallet.info/btsws"

	wm.Api = NewWalletClient(wm.Config.ServerAPI, wm.Config.WalletAPI, false)
	return wm
}
