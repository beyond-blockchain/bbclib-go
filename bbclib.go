/*
Copyright (c) 2018 Zettant Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Package bbclib is a library for defining BBcTransaction. This also provides serializer/deserializer and utilities for BBcTransaction object manipulation.


Serialization and deserialization

A BBcTransaction object contains various object, such as BBcEvent, BBcSignature.
In order to store a BBcTransaction object in DB or send it to other host, the object must be serialized.
Before serialization, the object is packed, meaning that it is transformed into binary format.
Then, the header is prepended to the packed data, resulting in a serialized data.
According to the header value, the packed data is compressed, so that you will get a smaller-sized serialized data.
Deserialization is the opposite transformation to serialization.

Utility functions

To build a BBcTransaction you need to create (new) objects you want to include. In many cases, it is a kind of common coding manner.
The utility functions are helpers to build a BBcTransaction with various objects.
*/
package bbclib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// Header values for serialized data
const (
	FormatPlain = 0x0000
	FormatZlib  = 0x0010
)

const (
	defaultIDLength = 32
)

/*
Serialize BBcTransaction object into packed data

formatType = 0x0000: Packed data is simply used for serialized data.

formatType = 0x0010: Packed data is compressed using zlib, and the compressed data is used for serialized data.
*/
func Serialize(transaction *BBcTransaction, formatType uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	Put2byte(buf, formatType)
	dat, err := transaction.Pack()
	if err != nil {
		return nil, err
	}

	if formatType == FormatPlain {
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return nil, err
		}
	} else if formatType == FormatZlib {
		compressed := ZlibCompress(&dat)
		if err := binary.Write(buf, binary.LittleEndian, compressed); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// Deserialize BBcTransaction data with header
func Deserialize(dat []byte) (*BBcTransaction, error) {
	buf := bytes.NewBuffer(dat)

	formatType, err := Get2byte(buf)
	if err != nil {
		return nil, err
	}

	txdat, err := GetBytes(buf, len(dat)-2)
	if err != nil {
		return nil, err
	}

	if formatType == FormatPlain {
		txobj := BBcTransaction{}
		err2 := txobj.Unpack(&txdat)
		return &txobj, err2
	} else if formatType == FormatZlib {
		decompressed, err := ZlibDecompress(txdat)
		if err != nil {
			return nil, err
		}
		txobj := BBcTransaction{}
		err2 := txobj.Unpack(&decompressed)
		return &txobj, err2
	}

	return nil, errors.New("formatType not supported")
}

// MakeTransaction is a utility for making simple BBcTransaction object with BBcEvent, BBcRelation or/and BBcWitness
func MakeTransaction(eventNum, relationNum int, witness bool, idLength int) *BBcTransaction {
	txobj := BBcTransaction{Version: 1, IDLength: idLength}
	txobj.Timestamp = time.Now().UnixNano() / int64(time.Microsecond)

	for i := 0; i < eventNum; i++ {
		evt := BBcEvent{}
		txobj.AddEvent(&evt)
	}

	for i := 0; i < relationNum; i++ {
		rtn := BBcRelation{}
		txobj.AddRelation(&rtn)
	}

	if witness {
		wit := BBcWitness{}
		txobj.AddWitness(&wit)
	}

	return &txobj
}

// addInRelation is an internal function to create a BBcAsset and add it to  BBcRelation object and then a BBcTransaction object
func addInRelation(transaction *BBcTransaction, relationIdx int, assetGroupID, userID *[]byte) {
	ast := BBcAsset{}
	transaction.Relations[relationIdx].Add(assetGroupID, &ast)
	ast.Add(userID)
}

// AddRelationAssetFile sets a file digest to BBcAsset in BBcRelation and add it to a BBcTransaction object
func AddRelationAssetFile(transaction *BBcTransaction, relationIdx int, assetGroupID, userID, assetFile *[]byte) {
	if transaction == nil {
		return
	}
	addInRelation(transaction, relationIdx, assetGroupID, userID)
	if assetFile != nil {
		transaction.Relations[relationIdx].Asset.AddFile(assetFile)
	}
}

// AddRelationAssetBodyString sets a string in BBcAsset in BBcRelation and add it to a BBcTransaction object
func AddRelationAssetBodyString(transaction *BBcTransaction, relationIdx int, assetGroupID, userID *[]byte, body string) {
	if transaction == nil {
		return
	}
	addInRelation(transaction, relationIdx, assetGroupID, userID)
	if body != "" {
		transaction.Relations[relationIdx].Asset.AddBodyString(body)
	}
}

// AddRelationAssetBodyObject sets an object (map[string]interface{}) in BBcAsset in BBcRelation, convert the info into msgpack, and add it in a BBcTransaction object
func AddRelationAssetBodyObject(transaction *BBcTransaction, relationIdx int, assetGroupID, userID *[]byte, body interface{}) {
	if transaction == nil {
		return
	}
	addInRelation(transaction, relationIdx, assetGroupID, userID)
	if body != nil {
		if err := transaction.Relations[relationIdx].Asset.AddBodyObject(body); err != nil {
			fmt.Println("Fail to exec AddRelationAssetBodyObject")
		}
	}
}

// AddRelationPointer creates and includes a BBcPointer object in BBcRelation and then, add it in a BBcTransaction object
func AddRelationPointer(transaction *BBcTransaction, relationIdx int, refTransactionID, refAssetID *[]byte) {
	if transaction == nil {
		return
	}
	ptr := BBcPointer{}
	transaction.Relations[relationIdx].AddPointer(&ptr)
	ptr.Add(refTransactionID, refAssetID)
}

// AddPointerInRelation creates and includes a BBcPointer object in BBcRelation
func AddPointerInRelation(relation *BBcRelation, refTransaction *BBcTransaction, refAssetID *[]byte) {
	ptr := BBcPointer{}
	relation.AddPointer(&ptr)
	ptr.Add(&refTransaction.TransactionID, refAssetID)
}

// AddReference creates and includes a BBcReference object in a BBcTransaction object
func AddReference(transaction *BBcTransaction, assetGroupID *[]byte, refTransaction *BBcTransaction, eventIdx int) {
	if transaction == nil || refTransaction == nil {
		return
	}
	if refTransaction.TransactionID == nil {
		refTransaction.Digest()
	}
	ref := BBcReference{}
	transaction.AddReference(&ref)
	ref.Add(assetGroupID, refTransaction, eventIdx)
}

// addInEvent is an internal function to add a BBcEvent object in a BBcTransaction object
func addInEvent(transaction *BBcTransaction, eventIdx int, assetGroupID, userID *[]byte) {
	ast := BBcAsset{}
	transaction.Events[eventIdx].Add(assetGroupID, &ast)
	ast.Add(userID)
}

// AddEventAssetFile sets a file digest to a BBcAsset object in a BBcEvent object and then, add it in a BBcTransaction object
func AddEventAssetFile(transaction *BBcTransaction, eventIdx int, assetGroupID, userID *[]byte, assetFile *[]byte) {
	if transaction == nil {
		return
	}
	addInEvent(transaction, eventIdx, assetGroupID, userID)
	if assetFile != nil {
		transaction.Events[eventIdx].Asset.AddFile(assetFile)
	}
}

// AddEventAssetBodyString sets a string to a BBcAsset object in a BBcEvent object and then, add it in a BBcTransaction object
func AddEventAssetBodyString(transaction *BBcTransaction, eventIdx int, assetGroupID, userID *[]byte, body string) {
	if transaction == nil {
		return
	}
	addInEvent(transaction, eventIdx, assetGroupID, userID)
	if body != "" {
		transaction.Events[eventIdx].Asset.AddBodyString(body)
	}
}

// AddEventAssetBodyObject sets an object (map[string]interface{}) to a BBcAsset object in a BBcEvent object and then, add it in a BBcTransaction object
func AddEventAssetBodyObject(transaction *BBcTransaction, eventIdx int, assetGroupID, userID *[]byte, body interface{}) {
	if transaction == nil {
		return
	}
	addInEvent(transaction, eventIdx, assetGroupID, userID)
	if body != "" {
		transaction.Events[eventIdx].Asset.AddBodyObject(body)
	}
}

// MakeRelationWithAsset is a utility for making simple BBcTransaction object with BBcRelation with BBcAsset
func MakeRelationWithAsset(assetGroupID, userID *[]byte, assetBodyString string, assetBodyObject interface{}, assetFile *[]byte, idLength int) *BBcRelation {
	rtn := BBcRelation{IDLength: idLength}
	ast := BBcAsset{IDLength: idLength}
	ast.Add(userID)
	rtn.Add(assetGroupID, &ast)
	if assetFile != nil {
		ast.AddFile(assetFile)
	}
	if assetBodyString != "" {
		ast.AddBodyString(assetBodyString)
	} else if assetBodyObject != nil {
		ast.AddBodyObject(assetBodyObject)
	}
	return &rtn
}

// SignToTransaction signs the transaction and append the BBcSignature object to it
func SignToTransaction(transaction *BBcTransaction, userID *[]byte, keypair *KeyPair) {
	sig := BBcSignature{}
	sig.SetPublicKeyByKeypair(keypair)
	signature, _ := transaction.Sign(keypair)
	sig.SetSignature(&signature)
	transaction.AddSignature(userID, &sig)
}
