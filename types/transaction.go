package types

import (
	"github.com/Assetsadapter/bitshares-adapter/encoding"
	"github.com/pkg/errors"
)

/*
{
    "ref_block_num": 43371,
    "ref_block_prefix": 879685531,
    "expiration": "2020-06-05T06:25:06",
    "operations": [
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
    ],
    "extensions": [],
    "signatures": [
        "1f611f7bab3df325cfec4273b9fb56ad56f069cd325c2c35bc5e500416405d7cf85873b32793d11743dfd378fe9e1449e4e72b0215efa9645eda4937e0cdde98c2"
    ]
}
*/

type Transaction struct {
	RefBlockNum    uint16   `json:"ref_block_num"`
	RefBlockPrefix uint32   `json:"ref_block_prefix"`
	Expiration     Time     `json:"expiration"`
	Operations     []*Op    `json:"operations"`
	Signatures     []string `json:"signatures"`
	Extensions     []string `json:"extensions"`
	TransactionID  string
}

type Op struct {
	Type  string `json:"type"`
	Value V      `json:"value"`
}

type V struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount Amount `json:"amount"`
	Memo   string `json:"memo"`
}

type Amount struct {
	Amount    string `json:"amount"`
	Precision int    `json:"precision"`
	Nai       string `json:"nai"`
}

// Marshal implements encoding.Marshaller interface.
func (tx *Transaction) Marshal(encoder *encoding.Encoder) error {
	if len(tx.Operations) == 0 {
		return errors.New("no operation specified")
	}

	enc := encoding.NewRollingEncoder(encoder)

	enc.Encode(tx.RefBlockNum)
	enc.Encode(tx.RefBlockPrefix)
	enc.Encode(tx.Expiration)

	enc.EncodeUVarint(uint64(len(tx.Operations)))
	for _, op := range tx.Operations {
		enc.Encode(op)
	}

	// Extensions are not supported yet.
	enc.EncodeUVarint(0)
	return enc.Err()
}

// PushOperation can be used to add an operation into the encoding.
func (tx *Transaction) PushOperation(op Op) {
	tx.Operations = append(tx.Operations, &op)
}
