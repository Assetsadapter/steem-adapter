package types

import (
	"encoding/hex"
	"encoding/json"

	"fmt"
	"reflect"

	"github.com/Assetsadapter/steem-adapter/encoding"
	"github.com/tidwall/gjson"
)

type Operation interface {
	Type() string
}

type Operations []Operation

func (ops *Operations) UnmarshalJSON(b string) (err error) {
	r := gjson.Get(b, "operations")
	rs := r.Array()
	for _, r := range rs {
		t := r.Get("type")
		_type := t.String()
		val, err := unmarshalOperation(_type, []byte(r.Raw))
		if err != nil {
			return err
		}
		*ops = append(*ops, val)
	}
	return nil
}

type operationTuple struct {
	Type string
	Data Operation
}

func (op *operationTuple) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		op.Type,
		op.Data,
	})
}

func (ops Operations) MarshalJSON() ([]byte, error) {
	tuples := make([]*operationTuple, 0, len(ops))
	for _, op := range ops {
		tuples = append(tuples, &operationTuple{
			Type: op.Type(),
			Data: op,
		})
	}
	return json.Marshal(tuples)
}

func unmarshalOperation(opType string, obj []byte) (Operation, error) {
	op, ok := knownOperations[opType]
	if !ok {
		// operation is unknown wrap it as an unknown operation
		val := UnknownOperation{
			kind: opType,
			Data: obj,
		}
		return &val, nil
	} else {
		val := reflect.New(op).Interface()
		if err := json.Unmarshal(obj, val); err != nil {
			return nil, err
		}
		return val.(Operation), nil
	}
}

var knownOperations = map[string]reflect.Type{
	"transfer_operation": reflect.TypeOf(TransferOperation{}),
	//"limit_order_create": reflect.TypeOf(LimitOrderCreateOperation{}),
	//"limit_order_cancel": reflect.TypeOf(LimitOrderCancelOperation{}),
}

// UnknownOperation
type UnknownOperation struct {
	kind string
	Data json.RawMessage
}

func (op *UnknownOperation) Type() string { return op.kind }

// NewTransferOperation returns a new instance of TransferOperation
func NewTransferOperation(from, to, memo string, amount Amount) *TransferOperation {
	op := &TransferOperation{
		Value{
			From: from,
			To:   to,
			Amount: Amount{
				Amount:    amount.Amount,
				Precision: amount.Precision,
				Nai:       amount.Nai,
			},
			Memo: memo,
		},
	}
	return op
}

// TransferOperation
type TransferOperation struct {
	Value Value `json:"value"`
}

type Value struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount Amount `json:"amount"`
	Memo   string `json:"memo"`
}

func (op *TransferOperation) Type() string { return "transfer_operation" }

func (op *TransferOperation) Marshal(encoder *encoding.Encoder) error {
	enc := encoding.NewRollingEncoder(encoder)

	enc.Encode(op.Type())
	enc.Encode(op.Value.From)
	enc.Encode(op.Value.To)
	enc.Encode(op.Value.Amount)
	enc.Encode(op.Value.Memo)

	enc.EncodeUVarint(0)
	return enc.Err()
}

type Buffer []byte

func (p Buffer) String() string {
	return hex.EncodeToString(p)
}

func (p *Buffer) FromString(data string) error {
	buf, err := hex.DecodeString(data)
	if err != nil {
		return fmt.Errorf("DecodeString: %v", err)
	}

	*p = buf
	return nil
}

func (p Buffer) Bytes() []byte {
	return p
}

func (p Buffer) Marshal(encoder *encoding.Encoder) error {
	enc := encoding.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(len(p)))
	enc.Encode(p.Bytes())
	return enc.Err()
}

func (p Buffer) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Buffer) UnmarshalJSON(data []byte) error {
	var b string
	if err := json.Unmarshal(data, &b); err != nil {
		return fmt.Errorf("Unmarshal: %s", err.Error())
	}

	return p.FromString(b)
}
