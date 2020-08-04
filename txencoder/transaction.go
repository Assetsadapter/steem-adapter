package txencoder

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	NAI_STEEM = "STEEM"
	NAI_SDB   = "SDB"
	NAI_TESTS = "TESTS"
	NAI_VESTS = "VESTS"
	NAI_HDB   = "HDB"
	NAI_TDB   = "TBD"
	NAI_HIVE  = "HIVE"
)

type RawTransaction struct {
	RefBlockNum    uint16          `json:"ref_block_num"`    // 参考的区块号
	RefBlockPrefix uint32          `json:"ref_block_prefix"` // 参考区块id
	Expiration     time.Time       `json:"expiration"`       // 交易到期时间
	Operations     *[]RawOperation `json:"operations"`       // 交易操作
	Extensions     *[]Extension    `json:"extensions"`       // 交易扩展
	Signatures     *[]string       `json:"signatures"`       // 交易签名
}

type RawOperation interface {
	OpType() OpType
	ParseToBroadcastJson() interface{}
}

func NewRawOperation(opType OpType) RawOperation {
	switch opType {
	case Vote:
	case Comment:
	case Transfer:
		return &RawTransferOperation{}
	}
	return nil
}

type RawTransferOperation struct {
	Type   OpType    `json:",omitempty"`
	From   string    `json:",omitempty"`
	To     string    `json:",omitempty"`
	Amount RawAmount `json:",omitempty"`
	Memo   string    `json:",omitempty"`
}

/*
[
        {
            "ref_block_num": 56885,
            "ref_block_prefix": 2621693493,
            "expiration": "2020-07-22T16:40:11",
            "operations": [
                [
                    "transfer",
                    {
                        "from": "exx-withdraw",
                        "to": "leor",
                        "amount": "1.000 TESTS",
                        "memo": ""
                    }
                ]
            ],
            "extensions": [],
            "signatures": [
                "1f6b0253d1613241ce4b88d7b71d06306cba2bd64a8bcce9aa1ae2e559620bc4ae470356488c081db77d3b342cfc31b1235dcb3ddfbea5e338db3df3f7634b2b84"
            ]
        }
    ]
*/

type TransferJson struct {
	RefBlockNum    uint16        `json:"ref_block_num"`
	RefBlockPrefix uint32        `json:"ref_block_prefix"`
	Expiration     string        `json:"expiration"`
	Operations     []interface{} `json:"operations"`
	Extensions     []string      `json:"extensions"`
	Signatures     []string      `json:"signatures"`
}

type AmountJson struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Memo   string `json:"memo"`
}

func (txOp *RawTransferOperation) OpType() OpType {
	return txOp.Type
}

/*
[
	"transfer",
	{
		"from": "hiveio",
		"to": "alice",
		"amount": {
			"amount": "10",
			"precision": 3,
			"nai": "@@000000021"
		},
		"memo": "Thanks for all the fish."
	}
]
*/
func (txOp *RawTransferOperation) ParseToBroadcastJson() interface{} {

	opJson := []interface{}{"transfer"}
	amount := decimal.New(int64(txOp.Amount.Amount), -int32(txOp.Amount.Precision))
	amountJson := fmt.Sprintf("%s %s", amount.StringFixed(int32(txOp.Amount.Precision)), RemoveFillChar(txOp.Amount.Nai))
	opJson = append(opJson, AmountJson{
		From:   txOp.From,
		To:     txOp.To,
		Amount: amountJson,
		Memo:   txOp.Memo,
	})

	//result := `["transfer", {"from": "` + txOp.From + `", "to": "` + txOp.To + `", "amount": "` + amountJson + `", "memo": "` + txOp.Memo + `"}]`
	//bcJson := fmt.Sprintf("[\"transfer\", {\"from\": \"%s\", \"to\": \"%s\", \"amount\": \"%s\", \"memo\": \"%s\"}]",
	//	txOp.From,
	//	txOp.To,
	//	amountJson,
	//	txOp.Memo)
	return opJson
}

type RawAmount struct {
	Amount    uint64 `json:"amount"`
	Precision uint8  `json:"precision"`
	Nai       string `json:"nai"`
}

/*
{
	"amount": "10",
	"precision": 3,
	"nai": "@@000000021"
}
*/
func (a *RawAmount) ParseToBroadcastJson() interface{} {
	result := `{"amount": "` + fmt.Sprintf("%d", a.Amount) + `", "precision": ` + fmt.Sprintf("%d", a.Precision) + `, "nai": "` + GetNai(a.Nai) + `"}`
	return result
}

func GetPrecision(coin string) uint8 {

	switch {
	case strings.HasPrefix(coin, NAI_VESTS):
		return 6
	case strings.HasPrefix(coin, NAI_HDB):
		fallthrough
	case strings.HasPrefix(coin, NAI_HIVE):
		fallthrough
	case strings.HasPrefix(coin, NAI_SDB):
		fallthrough
	case strings.HasPrefix(coin, NAI_STEEM):
		fallthrough
	case strings.HasPrefix(coin, NAI_TESTS):
		return 3
	}
	return 3
}

func GetNai(coin string) string {
	switch {
	case strings.HasPrefix(coin, NAI_TESTS):
		fallthrough
	case strings.HasPrefix(coin, NAI_HIVE):
		fallthrough
	case strings.HasPrefix(coin, NAI_STEEM):
		return "@@000000021"
	case strings.HasPrefix(coin, NAI_SDB):
		fallthrough
	case strings.HasPrefix(coin, NAI_HDB):
		fallthrough
	case strings.HasPrefix(coin, NAI_TDB):
		return "@@000000013"
	case strings.HasPrefix(coin, NAI_VESTS):
		return "@@000000037"
	}
	return ""
}

func RemoveFillChar(nai string) string {
	switch {
	case strings.HasPrefix(nai, NAI_TDB):
		return NAI_TDB
	case strings.HasPrefix(nai, NAI_VESTS):
		return NAI_VESTS
	case strings.HasPrefix(nai, NAI_HDB):
		return NAI_HDB
	case strings.HasPrefix(nai, NAI_SDB):
		return NAI_SDB
	case strings.HasPrefix(nai, NAI_STEEM):
		return NAI_STEEM
	case strings.HasPrefix(nai, NAI_HIVE):
		return NAI_HIVE
	case strings.HasPrefix(nai, NAI_TESTS):
		return NAI_TESTS
	}
	return ""
}

func CreateAmount(amount, nai string) (*RawAmount, error) {
	spl := strings.Split(amount, " ")
	if 2 != len(spl) {
		return nil, fmt.Errorf("Invalid amount ")
	}
	amo, err := strconv.ParseUint(spl[0], 10, 64)
	if err != nil {
		return nil, err
	}
	return &RawAmount{
		Amount:    amo,
		Precision: GetPrecision(spl[1]),
		Nai:       spl[1],
	}, nil
}

func (rt *RawTransaction) Encode() (*Transaction, error) {
	tx, err := newEmptyTransaction(rt.RefBlockNum, rt.Expiration, rt.RefBlockPrefix, rt.Operations)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (rt *RawTransaction) Decoder(data []byte) error {
	tx := Transaction{}
	_, err := tx.Decode(0, data)
	if err != nil {
		return err
	}
	err = rt.decode(tx)
	if err != nil {
		return err
	}
	return nil
}

const layout = "2006-01-02T15:04:05"

func (rt *RawTransaction) decode(tx Transaction) error {
	rt.RefBlockNum = littleEndianBytesToUint16(tx.RefBlockNum)
	rt.RefBlockPrefix = littleEndianBytesToUint32(tx.RefBlockPrefix)
	rt.Expiration = time.Unix(int64(littleEndianBytesToUint32(tx.Expiration)), 8).UTC()

	rt.Operations = &[]RawOperation{}
	if len(*tx.Operations) > 0 {
		for _, op := range *tx.Operations {
			rawOp := op.(TxEncoder).DecodeRaw().(*RawTransferOperation)
			*rt.Operations = append(*rt.Operations, rawOp)
		}
	}
	rt.Extensions = &[]Extension{}
	if len(*tx.Extensions) > 0 {

	}
	rt.Signatures = &[]string{}
	if len(*tx.Signatures) > 0 {
		for _, signature := range *tx.Signatures {
			rawSign := hex.EncodeToString(signature)
			*rt.Signatures = append(*rt.Signatures, rawSign)
		}
	}
	return nil
}

/*
"trx":{
	"ref_block_num":1097,
	"ref_block_prefix":2181793527,
	"expiration":"2016-03-24T18:00:21",
	"operations":[
		[
			"transfer",
			{
				"from": "hiveio",
				"to": "alice",
				"amount": {
					"amount": "10",
					"precision": 3,
					"nai": "@@000000021"
				},
				"memo": "Thanks for all the fish."
			}
		]
	],
	"extensions":[],
	"signatures":[]
}
*/

func (rt *RawTransaction) ParseToBroadcastJson() interface{} {
	//revRef, err := reverseHexToBytes(rt.RefBlockPrefix)
	//if err != nil {
	//	panic(fmt.Sprintf("RefBlockPrefix reverseHexToBytes failed : %s", err.Error()))
	//}
	//refBlockPrefix, err := strconv.ParseUint(hex.EncodeToString(revRef), 16, 64)
	//
	//
	//if err != nil {
	//	panic(fmt.Sprintf("Parse refBlockPrefix failed : %s", err.Error()))
	//}
	ops := []interface{}{}
	for _, op := range *rt.Operations {
		ops = append(ops, op.ParseToBroadcastJson())
	}
	txJson := TransferJson{
		RefBlockNum:    rt.RefBlockNum,
		RefBlockPrefix: rt.RefBlockPrefix,
		Expiration:     rt.Expiration.Format("2006-01-02T15:04:05"),
		Operations:     ops,
		Extensions:     []string{},
		Signatures:     *rt.Signatures,
	}

	return txJson
}
