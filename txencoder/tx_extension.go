package txencoder

type Extension struct {
}

func (e *Extension) Encode() *[]byte {
	return nil
}

func (e *Extension) Decode(offset int, data []byte) (int, error) {
	return 0, nil
}
