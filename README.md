# steem-adapter

## 项目依赖库

- [openwallet](https://github.com/blocktree/openwallet.git)

## 如何测试

openwtester包下的测试用例已经集成了openwallet钱包体系，创建conf目录，新建STM.ini文件，编辑如下内容：

```ini
; Enable scanner block
isScan = true
; Custom account must be true, Otherwise the scan will not filter to the internal address
isCustomAccount = true
; Scanner notify cache block number
cacheBlockNum = 30
; Is't testNet
isTestNet = true
; Node url
serverAPI = "https://steemd.steemitdev.com/"
; wallet url
walletAPI = "http://192.168.4.182:9879"
; ChainID
chainID = "0000000000000000000000000000000000000000000000000000000000000000"
```