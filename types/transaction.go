package types

import (
	"github.com/Assetsadapter/bitshares-adapter/encoding"
	"github.com/pkg/errors"
)

/*
{
	"ref_block_num": 90,
	"ref_block_prefix": 1997317362,
	"expiration": "2020-04-07T02:45:51",
	"operations": [
		{
			"type": "transfer_operation",
			"value": {
				"from": "initminer",
				"to": "leor",
				"amount": {
					"amount": "1",
					"precision": 3,
					"nai": "@@000000013"
				},
				"memo": "test cli_wallet"
			}
		}
	],
	"extensions": [],
	"signatures": [
		"1f68b0915353ef12aa77697849f538d488d1d27fc9eee123ce32744c487b71e3e44ba3027f8cde7ab90a4aecf58cf37b506f2196e0de5d90dd6b93c607ccccc361"
	]
}
*/

type Transaction struct {
	RefBlockNum    uint16     `json:"ref_block_num"`
	RefBlockPrefix uint32     `json:"ref_block_prefix"`
	Expiration     Time       `json:"expiration"`
	Operations     Operations `json:"operations"`
	Signatures     []string   `json:"signatures"`
	Extensions     []string   `json:"extensions"`
	TransactionID  string
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
func (tx *Transaction) PushOperation(op Operation) {
	tx.Operations = append(tx.Operations, op)
}
