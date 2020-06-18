package txencoder

type Extension struct {
}

func (e *Extension) Decode() *[]byte {
	return nil
}

func (e *Extension) Encode(offset int, data []byte) (int, error) {
	return 0, nil
}
