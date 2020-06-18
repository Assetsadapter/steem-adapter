package encoding_

type Marshaller interface {
	Marshal(*Encoder) error
}
