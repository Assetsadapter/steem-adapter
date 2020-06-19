package txencoder

import (
	"encoding/hex"
	"time"

	"github.com/prometheus/common/log"
)

type TxEncoder interface {
	Encode() *[]byte
	Decode(offset int, data []byte) (int, error)
}

/*
	交易序列化组合方式：
	1byte ref_block_num +
	4byte ref_block_prefix +
	4byte expiration +
	1byte operations count +
	1byte operation_type +
	1byte from_data_len + from_data +
	1byte to_data_len +	to_data +
	8byte amount +
	1byte precision +
	7byte nai +
	1byte memo_data_len + memo_data +
    1byte extensions_count + extensions_data +
    1byte signatures_count + signatures_data
*/

type Transaction struct {
	RefBlockNum    []byte               // 2 byte, 当前区块号 & 0xFFFF 获取一个参考的区块
	RefBlockPrefix []byte               // 4 byte, 当前区块id的第四个字节开始取4个字节长度的值进行小端转换作为参考值
	Expiration     []byte               // 4 byte, 时间戳
	Operations     *[]TransferOperation // 交易的操作
	Extensions     *[]Extension         // 扩展
	Signatures     *[]Signature         // 签名
}

type Signature []byte

func newEmptyTransaction(refBlockNum uint16, expiration time.Time, refBlockPrefix string, op *[]RawTransferOperation) (*Transaction, error) {
	result := &Transaction{}
	result.RefBlockNum = uint16ToLittleEndianBytes(refBlockNum)
	refPrefix, err := hex.DecodeString(refBlockPrefix)
	if err != nil {
		log.Errorf("Reverse ref block prefix failed : %s", err.Error())
		return nil, err
	}
	result.RefBlockPrefix = refPrefix
	result.Expiration = uint32ToLittleEndianBytes(uint32(expiration.Unix()))
	ops, err := newEmptyTransferOperations(op)
	if err != nil {
		log.Errorf("Decoder transaction failed : %s", err.Error())
		return nil, err
	}
	result.Operations = ops
	result.Extensions = &[]Extension{}
	result.Signatures = &[]Signature{}
	return result, nil
}

func (tx *Transaction) Encode() *[]byte {
	bytesData := []byte{}
	bytesData = append(bytesData, tx.RefBlockNum...)
	bytesData = append(bytesData, tx.RefBlockPrefix...)
	bytesData = append(bytesData, tx.Expiration...)
	bytesData = append(bytesData, byte(len(*tx.Operations)))
	for _, op := range *tx.Operations {
		bytesData = append(bytesData, *op.Encode()...)
	}
	bytesData = append(bytesData, byte(len(*tx.Extensions)))
	bytesData = append(bytesData, byte(len(*tx.Signatures)))
	return &bytesData
}

func (tx *Transaction) Decode(offset int, data []byte) (int, error) {
	index := offset
	tx.RefBlockNum = data[index : index+2]
	index += 2
	tx.RefBlockPrefix = data[index : index+4]
	index += 4
	tx.Expiration = data[index : index+4]
	index += 4
	opsCount := int(data[index])
	index += 1
	tx.Operations = &[]TransferOperation{}
	for i := opsCount; i > 0; i-- {
		txOp := TransferOperation{}
		newOffset, err := txOp.Decode(index, data)
		if err != nil {
			return newOffset, err
		}
		*tx.Operations = append(*tx.Operations, txOp)
		index = newOffset
	}
	extenCount := int(data[index])
	index += 1
	tx.Extensions = &[]Extension{}
	if extenCount > 0 {
		e := Extension{}
		index, err := e.Decode(index, data)
		if err != nil {
			return index, err
		}
		*tx.Extensions = append(*tx.Extensions, e)
	}
	signCount := int(data[index])
	index += 1
	tx.Signatures = &[]Signature{}
	if signCount > 0 {
		for i := signCount; i > 0; i-- {
			*tx.Signatures = append(*tx.Signatures, data[index:index+signCount])
		}
	}
	return index, nil
}
