package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Assetsadapter/bitshares-adapter/encoding"
	"github.com/pkg/errors"
)

/*
 {
	"type": "transfer_operation",
	"value": {
		"from": "porter",
		"to": "aaronnn47",
		"amount": {
			"amount": "1117817",
			"precision": 3,
			"nai": "@@000000021"
		},
		"memo": "Ported balance"
	}
}
*/

type Operation interface {
	Type() OpType
}

type Operations []Operation

func (ops *Operations) UnmarshalJSON(b []byte) (err error) {
	// unmarshal array
	var o []json.RawMessage
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}

	// foreach operation
	for _, op := range o {
		var kv []json.RawMessage
		if err := json.Unmarshal(op, &kv); err != nil {
			return err
		}

		if len(kv) != 2 {
			return errors.New("invalid operation format: should be name, value")
		}

		var opType uint16
		if err := json.Unmarshal(kv[0], &opType); err != nil {
			return err
		}

		val, err := unmarshalOperation(OpType(opType), kv[1])
		if err != nil {
			return err
		}

		*ops = append(*ops, val)
	}

	return nil
}

type operationTuple struct {
	Type OpType
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

func unmarshalOperation(opType OpType, obj json.RawMessage) (Operation, error) {
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

var knownOperations = map[OpType]reflect.Type{
	TransferOpType:         reflect.TypeOf(TransferOperation{}),
	LimitOrderCreateOpType: reflect.TypeOf(LimitOrderCreateOperation{}),
	LimitOrderCancelOpType: reflect.TypeOf(LimitOrderCancelOperation{}),
}

// UnknownOperation
type UnknownOperation struct {
	kind OpType
	Data json.RawMessage
}

func (op *UnknownOperation) Type() OpType { return op.kind }

// NewTransferOperation returns a new instance of TransferOperation
func NewTransferOperation(from, to, amount string, fee AssetAmount) *TransferOperation {
	op := &TransferOperation{
		Value{
			From: from,
			To:   to,
			Amount: AssetAmount{
				//Amount:    amount,
				Precision: 3,
				Nai:       "",
			},
			Memo: "",
		},
	}

	return op
}

/*
{
    "type": "transfer_operation",
    "value": {
        "from": "initminer",
        "to": "leor",
        "amount": {
            "amount": "1000",
            "precision": 3,
            "nai": "@@000000013"
        },
        "memo": "990909"
    }
}
*/

// TransferOperation
type TransferOperation struct {
	Value Value `json:"value"`
}

type Value struct {
	From   string      `json:"from"`
	To     string      `json:"to"`
	Amount AssetAmount `json:"amount"`
	Memo   string      `json:"memo"`
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
		return fmt.Errorf("Unmarshal: %v", err)
	}

	return p.FromString(b)
}

type Memo struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Nonce   string `json:"nonce"`
	Message Buffer `json:"message"`
}

func (m Memo) Marshal(encoder *encoding.Encoder) error {
	enc := encoding.NewRollingEncoder(encoder)
	enc.Encode(m.From)
	enc.Encode(m.To)
	enc.Encode(m.Nonce)
	enc.Encode(m.Message)
	return enc.Err()
}

func (op *TransferOperation) Type() OpType { return TransferOpType }

func (op *TransferOperation) Marshal(encoder *encoding.Encoder) error {
	enc := encoding.NewRollingEncoder(encoder)

	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Value.From)
	enc.Encode(op.Value.To)
	enc.Encode(op.Value.Amount)
	enc.Encode(op.Value.Memo)

	//Memo?
	// enc.EncodeUVarint(0)
	//Extensions
	enc.EncodeUVarint(0)
	return enc.Err()
}

// LimitOrderCreateOperation
type LimitOrderCreateOperation struct {
	Fee          AssetAmount       `json:"fee"`
	Seller       ObjectID          `json:"seller"`
	AmountToSell AssetAmount       `json:"amount_to_sell"`
	MinToReceive AssetAmount       `json:"min_to_receive"`
	Expiration   Time              `json:"expiration"`
	FillOrKill   bool              `json:"fill_or_kill"`
	Extensions   []json.RawMessage `json:"extensions"`
}

func (op *LimitOrderCreateOperation) Marshal(encoder *encoding.Encoder) error {
	enc := encoding.NewRollingEncoder(encoder)

	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Fee)
	enc.Encode(op.Seller)
	enc.Encode(op.AmountToSell)
	enc.Encode(op.MinToReceive)
	enc.Encode(op.Expiration)
	enc.EncodeBool(op.FillOrKill)

	//extensions
	enc.EncodeUVarint(0)
	return enc.Err()
}

func (op *LimitOrderCreateOperation) Type() OpType { return LimitOrderCreateOpType }

// LimitOrderCancelOpType
type LimitOrderCancelOperation struct {
	Fee              AssetAmount       `json:"fee"`
	FeePayingAccount ObjectID          `json:"fee_paying_account"`
	Order            ObjectID          `json:"order"`
	Extensions       []json.RawMessage `json:"extensions"`
}

func (op *LimitOrderCancelOperation) Marshal(encoder *encoding.Encoder) error {
	enc := encoding.NewRollingEncoder(encoder)

	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Fee)
	enc.Encode(op.FeePayingAccount)
	enc.Encode(op.Order)

	// extensions
	enc.EncodeUVarint(0)
	return enc.Err()
}

func (op *LimitOrderCancelOperation) Type() OpType { return LimitOrderCancelOpType }

// FillOrderOpType
type FillOrderOperation struct {
	Order   ObjectID
	Account ObjectID
	Pays    AssetAmount
	Recives AssetAmount
	Fee     AssetAmount
	Price   Price
	IsMaker bool
}

func (op *FillOrderOperation) Type() OpType { return FillOrderOpType }
