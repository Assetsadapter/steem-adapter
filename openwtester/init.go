package openwtester

import (
	"github.com/Assetsadapter/steem-adapter/steem"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openw"
)

func init() {
	//注册钱包管理工具
	log.Notice("Wallet Manager Load Successfully.")
	cache := steem.NewCacheManager()

	openw.RegAssets(steem.Symbol, steem.NewWalletManager(&cache))
}
