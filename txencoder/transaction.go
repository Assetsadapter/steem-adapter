package txencoder

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	NAI_STEEM = "STEEM"
	NAI_SDB   = "SDB"
	NAI_TESTS = "TESTS"
	NAI_VESTS = "VESTS"
	NAI_HDB   = "HDB"
	NAI_HIVE  = "HIVE"
)

type RawTransaction struct {
	RefBlockNum    uint16          // 参考的区块号
	RefBlockPrefix string          // 参考区块id
	Expiration     time.Time       // 交易到期时间
	Operations     *[]RawOperation // 交易操作
	Signature      *[]string       // 交易签名
}

type RawOperation interface {
	OpType() OpType
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
	Type   OpType
	From   string
	To     string
	Amount RawAmount
	Memo   string
}

func (txOp *RawTransferOperation) OpType() OpType {
	return txOp.Type
}

type RawAmount struct {
	Amount    uint64
	Precision uint8
	Nai       string
}

func GetPrecision(coin string) uint8 {
	switch coin {
	case NAI_VESTS:
		return 6
	case NAI_HDB:
		fallthrough
	case NAI_HIVE:
		fallthrough
	case NAI_SDB:
		fallthrough
	case NAI_STEEM:
		fallthrough
	case NAI_TESTS:
		return 3
	}
	return 3
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

func (rt *RawTransaction) decode(tx Transaction) error {
	rt.RefBlockNum = littleEndianBytesToUint16(tx.RefBlockPrefix)
	rt.RefBlockPrefix = hex.EncodeToString(tx.RefBlockPrefix)
	rt.Expiration = time.Unix(int64(littleEndianBytesToUint32(tx.Expiration)), 8)
	rt.Operations = &[]RawOperation{}
	for _, op := range *tx.Operations {
		rawOp := op.(TxEncoder).DecodeRaw().(*RawTransferOperation)
		*rt.Operations = append(*rt.Operations, rawOp)
	}
	return nil
}
