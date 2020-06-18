package txencoder

import "time"

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
