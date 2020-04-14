# steem-adapter

## 项目依赖库

- [openwallet](https://github.com/blocktree/openwallet.git)

## 如何测试

openwtester包下的测试用例已经集成了openwallet钱包体系，创建conf目录，新建STM.ini文件，编辑如下内容：

```ini

#wallet api url
; is enable block chain scanner
isScan = true
isCustomAccount = true
#core api url
ServerAPI = "http://127.0.0.1:8090"
; wallet api url
walletAPI = "http://127.0.0.1:8093"
# ChainID
ChainID = "18dcf0a285365fc58b71f18b3d3fec954aa0c141c44e4e5cb4cf777b9eab274e"
# MemoPrivateKey 需要将充值账号和提币账号的memokey设置成同一个 简化逻辑
MemoPrivateKey = "5JtEo2yUL4rFsKT1DvEFk2bVdqMUiyxCBZSCwt4WyTvqjGDb7mU"
```