package txencoder

import (
	"encoding/hex"
	"time"
)

type RawTransaction struct {
	RefBlockNum    uint16
	RefBlockPrefix string
	Expiration     time.Time
	Operations     *[]RawTransferOperation
}

type RawTransferOperation struct {
	Type   uint8
	From   string
	To     string
	Amount RawAmount
	Memo   string
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
	rt.Operations = &[]RawTransferOperation{}
	for _, op := range *tx.Operations {
		txOp := RawTransferOperation{
			Type: uint8(op.Type),
			From: string(op.From),
			To:   string(op.To),
			Amount: RawAmount{
				Amount:    littleEndianBytesToUint64(op.Amount.Amount),
				Precision: op.Amount.Precision,
				Nai:       string(op.Amount.Nai),
			},
			Memo: string(op.Memo),
		}
		*rt.Operations = append(*rt.Operations, txOp)
	}
	return nil
}
