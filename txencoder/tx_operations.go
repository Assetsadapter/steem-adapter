package txencoder

type TransferOperation struct {
	Type   byte   // 1 byte
	From   []byte // data len
	To     []byte // data len
	Amount *Amount
	Memo   []byte // data len
}

func newEmptyTransferOperations(data *[]RawTransferOperation) (*[]TransferOperation, error) {
	ops := &[]TransferOperation{}
	for _, operation := range *data {
		result := TransferOperation{}
		result.Type = operation.Type
		result.From = []byte(operation.From)
		result.To = []byte(operation.To)
		_amount, err := newEmptyAmount(&operation.Amount)
		if err != nil {
			return nil, err
		}
		result.Amount = _amount
		result.Memo = []byte(operation.Memo)
		*ops = append(*ops, result)
	}

	return ops, nil
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
