package txencoder

import (
	"encoding/hex"
	"time"
)

type RawTransaction struct {
	RefBlockNum    uint16
	RefBlockPrefix string
	Expiration     time.Time
	Operations     *[]RawOperation
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
