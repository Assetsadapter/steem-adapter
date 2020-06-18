package types

type OperationSerialization interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type TransferSerialization struct {
	Type   uint8
	From   string
	To     string
	Amount string
	Memo   string
}

func (tx *TransferSerialization) Marshal() (ret []byte, err error) {

	return nil, nil
}

func (tx *TransferSerialization) Unmarshal([]byte) error {
	return nil
}
