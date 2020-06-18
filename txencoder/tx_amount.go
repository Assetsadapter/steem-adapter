package txencoder

type Amount struct {
	Amount    []byte // 8 byte
	Precision byte   // 1 byte
	Nai       []byte // 7 byte
}

func newEmptyAmount(amount *RawAmount) (*Amount, error) {
	result := &Amount{}
	result.Amount = uint64ToLittleEndianBytes(amount.Amount)
	result.Precision = amount.Precision
	//naiByte, err := hex.DecodeString(amount.Nai)
	//if err != nil {
	//	log.Errorf("decoder nai failed : %s", err.Error())
	//	return nil, err
	//}
	result.Nai = fillNai([]byte(amount.Nai))
	return result, nil
}

// nai 需要用0补足7个byte
func fillNai(nai []byte) []byte {
	if len(nai) == 7 {
		return nai
	}
	fillLen := 7 - len(nai)
	for i := 0; i < fillLen; i++ {
		nai = append(nai, byte(0))
	}
	return nai
}

func (a *Amount) Decode() *[]byte {
	bytesData := []byte{}
	bytesData = append(bytesData, (*a).Amount...)
	bytesData = append(bytesData, (*a).Precision)
	bytesData = append(bytesData, (*a).Nai...)
	return &bytesData
}

func (a *Amount) Encode(offset int, data []byte) (int, error) {
	index := offset
	a.Amount = data[index : index+8]
	index += 8
	a.Precision = data[index]
	index += 1
	a.Nai = data[index : index+7]
	index += 7
	return index, nil
}
