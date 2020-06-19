package txencoder

import "errors"

type Operation interface {
	InitData(data *RawOperation) error
}

type TransferOperation struct {
	Type   byte   // 1 byte
	From   []byte // data len
	To     []byte // data len
	Amount *Amount
	Memo   []byte // data len
}

func NewOperation(opType OpType) Operation {
	switch opType {
	case Vote:
	case Comment:
	case Transfer:
		return &TransferOperation{}
	}
	return nil
}

func (txOp *TransferOperation) InitData(data *RawOperation) error {
	rawTxOp, ok := (*data).(*RawTransferOperation)
	if !ok {
		return errors.New("Init data failed : invalid raw transfer operation data ")
	}
	txOp.Type = byte(rawTxOp.Type)
	txOp.From = []byte(rawTxOp.From)
	txOp.To = []byte(rawTxOp.To)
	_amount, err := newEmptyAmount(&rawTxOp.Amount)
	if err != nil {
		return err
	}
	txOp.Amount = _amount
	txOp.Memo = []byte(rawTxOp.Memo)
	return nil
}

func (txOp *TransferOperation) Encode() *[]byte {
	bytesData := []byte{}
	bytesData = append(bytesData, txOp.Type)
	bytesData = append(bytesData)
	bytesData = append(bytesData, byte(len(txOp.From)))
	bytesData = append(bytesData, txOp.From...)
	bytesData = append(bytesData, byte(len(txOp.To)))
	bytesData = append(bytesData, txOp.To...)
	bytesData = append(bytesData, *txOp.Amount.Encode()...)
	bytesData = append(bytesData, byte(len(txOp.Memo)))
	bytesData = append(bytesData, txOp.Memo...)
	return &bytesData
}

func (txOp *TransferOperation) Decode(offset int, data []byte) (int, error) {
	index := offset
	txOp.Type = data[index]
	index += 1
	fromLen := int(data[index])
	index += 1
	txOp.From = data[index : index+fromLen]
	index += fromLen
	toLen := int(data[index])
	index += 1
	txOp.To = data[index : index+toLen]
	index += toLen
	amount := &Amount{}
	index, err := amount.Decode(index, data)
	if err != nil {
		return index, err
	}
	txOp.Amount = amount
	memoLen := int(data[index])
	index += 1
	txOp.Memo = data[index : index+memoLen]
	index += memoLen
	return index, nil
}

func (txOp *TransferOperation) DecodeRaw() interface{} {
	ret := &RawTransferOperation{
		Type:   OpType(txOp.Type),
		From:   string(txOp.From),
		To:     string(txOp.To),
		Amount: txOp.Amount.DecodeRaw().(RawAmount),
		Memo:   string(txOp.Memo),
	}
	return ret
}
