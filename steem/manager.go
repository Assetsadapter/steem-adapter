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
	"fmt"

	"github.com/Assetsadapter/steem-adapter/addrdec"
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
)

type WalletManager struct {
	openwallet.AssetsAdapterBase

	Api             *WalletClient                   // 节点客户端
	Config          *WalletConfig                   // 节点配置
	Decoder         openwallet.AddressDecoder       //地址编码器
	DecoderV2       openwallet.AddressDecoderV2     //地址编码器V2
	TxDecoder       openwallet.TransactionDecoder   //交易单编码器
	Log             *log.OWLogger                   //日志工具
	ContractDecoder openwallet.SmartContractDecoder //智能合约解析器
	Blockscanner    *StmBlockScanner                //区块扫描器
	CacheManager    openwallet.ICacheManager        //缓存管理器
	WebsocketAPI    string                          //steem WebsocketAPI
	ChainId         string
}

func NewWalletManager(cacheManager openwallet.ICacheManager) *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig(Symbol)
	wm.Api = NewWalletClient(wm.Config.ServerAPI, wm.Config.WalletAPI, false)
	wm.Blockscanner = NewBlockScanner(&wm)
	wm.Decoder = NewAddressDecoder(&wm)
	conf, _ := config.NewConfig("ini", fmt.Sprintf("conf/%s.ini", Symbol))
	isTestNet, _ := conf.Bool("isTestNet")
	wm.DecoderV2 = addrdec.NewAddressDecoderV2(isTestNet)
	wm.TxDecoder = NewTransactionDecoder(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())
	wm.CacheManager = cacheManager
	wm.ContractDecoder = NewContractDecoder(&wm)

	return &wm
}
