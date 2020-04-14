/*
 * Copyright 2018 The OpenWallet Authors
 * This file is part of the OpenWallet library.
 *
 * The OpenWallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The OpenWallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package steem

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/Assetsadapter/steem-adapter/encoding"
	"github.com/Assetsadapter/steem-adapter/types"
	"github.com/blocktree/openwallet/common"
	"github.com/blocktree/openwallet/crypto"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/tidwall/gjson"
)

type Asset struct {
	ID                 types.ObjectID `json:"id"`
	Symbol             string         `json:"symbol"`
	Precision          uint8          `json:"precision"`
	Issuer             string         `json:"issuer"`
	DynamicAssetDataID string         `json:"dynamic_asset_data_id"`
}

type BlockHeader struct {
	TransactionMerkleRoot string            `json:"transaction_merkle_root"`
	Previous              string            `json:"previous"`
	Timestamp             types.Time        `json:"timestamp"`
	Witness               string            `json:"witness"`
	Extensions            []json.RawMessage `json:"extensions"`
	WitnessSignature      string            `json:"witness_signature"`
}

func NewBlockHeader(result *gjson.Result) *BlockHeader {
	obj := BlockHeader{}
	json.Unmarshal([]byte(result.Raw), &obj)
	return &obj
}

func (block *BlockHeader) Serialize() ([]byte, error) {
	var b bytes.Buffer
	encoder := encoding.NewEncoder(&b)

	if err := encoder.Encode(block); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (block *BlockHeader) CalculateID() (string, error) {
	var msgBuffer bytes.Buffer

	// Write the serialized transaction.
	rawTx, err := block.Serialize()
	if err != nil {
		return "", err
	}

	if _, err := msgBuffer.Write(rawTx); err != nil {
		return "", errors.Wrap(err, "failed to write serialized block header")
	}

	msgBytes := msgBuffer.Bytes()

	// Compute the digest.
	digest := sha256.Sum224(msgBytes)

	id := hex.EncodeToString(digest[:])
	length := 40
	if len(id) < 40 {
		length = len(id)
	}
	return id[:length], nil
}

// MarshalBlockHeader implements encoding.Marshaller interface.
func (block *BlockHeader) Marshal(encoder *encoding.Encoder) error {

	enc := encoding.NewRollingEncoder(encoder)

	enc.Encode(block.TransactionMerkleRoot)
	enc.Encode(block.Previous)
	enc.Encode(block.Timestamp)
	enc.Encode(block.Witness)
	enc.Encode(block.WitnessSignature)

	// Extensions are not supported yet.
	enc.EncodeUVarint(0)
	return enc.Err()
}

/*
"block": {
	"previous": "0000005af2a40c77b87897e01a7430cc36ea578a",
	"timestamp": "2020-04-07T02:45:24",
	"witness": "initminer",
	"transaction_merkle_root": "989ff103edf047ae94e0f7e7e09cd1c0ccf7f4ed",
	"extensions": [],
	"witness_signature": "20a50f54af6fac486990ddacfb4d64dcaf8f53df8b6a9c9e6624f25ef67f1ea50728d93d11353183037312e686be47fcdcc45472e79f66701715090e64ed36a29e",
	"transactions": [
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
	],
	"block_id": "0000005b1feae084bd3bd9190c5ad53d372f5996",
	"signing_key": "TST6LLegbAgLAy28EHrffBVuANFWcFgmqRMW13wBmTExqFE9SCkg4",
	"transaction_ids": [
		"d92701a99c44fb9e0ee944497311775eaeeedb94"
	]
}
*/

type Block struct {
	Previous              string               `json:"previous"`
	Timestamp             types.Time           `json:"timestamp"`
	Witness               string               `json:"witness"`
	TransactionMerkleRoot string               `json:"transaction_merkle_root"`
	Extensions            []json.RawMessage    `json:"extensions"`
	WitnessSignature      string               `json:"witness_signature"`
	Transactions          []*types.Transaction `json:"transactions"`
	BlockID               string               `json:"block_id"`
	SigningKey            string               `json:"signing_key"`
	TransactionIDs        []string             `json:"transaction_ids"`
	Height                uint64
}

func NewBlock(height uint32, result *gjson.Result) *Block {
	obj := Block{}
	json.Unmarshal([]byte(result.Raw), &obj)
	obj.Height = uint64(height)
	return &obj
}

func (block *Block) CalculateID() error {
	header := BlockHeader{}
	header.TransactionMerkleRoot = block.TransactionMerkleRoot
	header.Previous = block.Previous
	header.Timestamp = block.Timestamp
	header.Witness = block.Witness
	header.Extensions = block.Extensions
	header.WitnessSignature = block.WitnessSignature

	id, err := header.CalculateID()
	if err != nil {
		return err
	}
	block.BlockID = id
	return nil
}

func NewTransaction(result *gjson.Result, transactionID string) (*types.Transaction, error) {
	obj := types.Transaction{}
	err := json.Unmarshal([]byte(result.Raw), &obj)
	obj.TransactionID = transactionID
	return &obj, err
}

//UnscanRecord 扫描失败的区块及交易
type UnscanRecord struct {
	ID          string `storm:"id"` // primary key
	BlockHeight uint64
	TxID        string
	Reason      string
}

//NewUnscanRecord new UnscanRecord
func NewUnscanRecord(height uint64, txID, reason string) *UnscanRecord {
	obj := UnscanRecord{}
	obj.BlockHeight = height
	obj.TxID = txID
	obj.Reason = reason
	obj.ID = common.Bytes2Hex(crypto.SHA256([]byte(fmt.Sprintf("%d_%s", height, txID))))
	return &obj
}

// ParseHeader 区块链头
func ParseHeader(b *Block) *openwallet.BlockHeader {
	obj := openwallet.BlockHeader{}

	//解析josn
	obj.Merkleroot = b.TransactionMerkleRoot
	obj.Hash = b.BlockID
	obj.Previousblockhash = b.Previous
	obj.Height = b.Height
	obj.Time = uint64(b.Timestamp.Unix())
	obj.Symbol = Symbol
	return &obj
}

type BlockchainInfo struct {
	HeadBlockNum             uint64    `json:"head_block_number"`
	HeadBlockID              string    `json:"head_block_id"`
	LastIrreversibleBlockNum uint64    `json:"last_irreversible_block_num"`
	Timestamp                time.Time `json:"time"`

	/*
		{
			"id": "2.1.0",
			"head_block_number": 1544081,
			"head_block_id": "00178f912d70e9ed3539f2acfba4752dee5d77bb",
			"time": "2019-07-17T04:09:40",
			"current_witness": "1.6.8",
			"next_maintenance_time": "2019-07-18T00:00:00",
			"last_budget_time": "2019-07-17T00:00:00",
			"witness_budget": 0,
			"accounts_registered_this_interval": 2,
			"recently_missed_count": 0,
			"current_aslot": 1672768,
			"recent_slots_filled": "340282366920938463463374607431768211455",
			"dynamic_flags": 0,
			"last_irreversible_block_num": 1544074
		}
	*/
}

const TimeLayout = `2006-01-02T15:04:05`

func NewBlockchainInfo(result *gjson.Result) *BlockchainInfo {
	obj := BlockchainInfo{}
	obj.HeadBlockNum = result.Get("head_block_number").Uint()
	obj.HeadBlockID = result.Get("head_block_id").String()
	obj.LastIrreversibleBlockNum = result.Get("last_irreversible_block_num").Uint()
	obj.Timestamp, _ = time.ParseInLocation(TimeLayout, result.Get("time").String(), time.UTC)
	return &obj
}

type Balance struct {
	AssetID types.ObjectID `json:"asset_id"`
	Amount  string         `json:"amount"`
}

func NewBalance(result *gjson.Result) *Balance {
	arr := result.Array()
	for _, item := range arr {
		obj := Balance{}
		obj.Amount = item.Get("amount").String()
		obj.AssetID = types.MustParseObjectID(item.Get("asset_id").String())
		return &obj
	}
	return nil
}

type BroadcastResponse struct {
	ID string `json:"id"`
}