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


// AddRelationAssetFile sets a file digest to BBcAsset in BBcRelation and add it to a BBcTransaction object (old style, only for backward compatibility)
func AddRelationAssetFile(transaction *BBcTransaction, relationIdx int, assetGroupId, userId, assetFile *[]byte) {
	if transaction == nil {
		return
	}
	transaction.Relations[relationIdx].AssetGroupID = make([]byte, transaction.IdLengthConf.AssetGroupIdLength)
	copy(transaction.Relations[relationIdx].AssetGroupID, *assetGroupId)
	transaction.Relations[relationIdx].CreateAsset(userId, assetFile, nil)
}

// AddRelationAssetBodyString sets a string in BBcAsset in BBcRelation and add it to a BBcTransaction object (old style, only for backward compatibility)
func AddRelationAssetBodyString(transaction *BBcTransaction, relationIdx int, assetGroupId, userId *[]byte, body string) {
	if transaction == nil {
		return
	}
	transaction.Relations[relationIdx].AssetGroupID = make([]byte, transaction.IdLengthConf.AssetGroupIdLength)
	copy(transaction.Relations[relationIdx].AssetGroupID, *assetGroupId)
	transaction.Relations[relationIdx].CreateAsset(userId, nil, body)
}

// AddRelationAssetBodyObject sets an object (map[string]interface{}) in BBcAsset in BBcRelation, convert the info into msgpack, and add it in a BBcTransaction object (old style, only for backward compatibility)
func AddRelationAssetBodyObject(transaction *BBcTransaction, relationIdx int, assetGroupId, userId *[]byte, body interface{}) {
	if transaction == nil {
		return
	}
	transaction.Relations[relationIdx].AssetGroupID = make([]byte, transaction.IdLengthConf.AssetGroupIdLength)
	copy(transaction.Relations[relationIdx].AssetGroupID, *assetGroupId)
	transaction.Relations[relationIdx].CreateAsset(userId, nil, body)
}

// AddRelationAssetRawBody sets a data in BBcAssetRaw in BBcRelation and add it to a BBcTransaction object (old style, only for backward compatibility)
func AddRelationAssetRaw(transaction *BBcTransaction, relationIdx int, assetGroupId, assetId *[]byte, assetBody interface{}) {
	transaction.Relations[relationIdx].AssetGroupID = make([]byte, transaction.IdLengthConf.AssetGroupIdLength)
	copy(transaction.Relations[relationIdx].AssetGroupID, *assetGroupId)
	transaction.Relations[relationIdx].CreateAssetRaw(assetId, assetBody)
}

// AddRelationAssetHash sets assetIDs in BBcAssetHash in BBcRelation and add it to a BBcTransaction object (old style, only for backward compatibility)
func AddRelationAssetHash(transaction *BBcTransaction, relationIdx int, assetGroupId, assetIds *[]byte) {
	transaction.Relations[relationIdx].AssetGroupID = make([]byte, transaction.IdLengthConf.AssetGroupIdLength)
	copy(transaction.Relations[relationIdx].AssetGroupID, *assetGroupId)
	transaction.Relations[relationIdx].CreateAssetHash(assetIds)
}

// AddRelationPointer creates and includes a BBcPointer object in BBcRelation and then, add it in a BBcTransaction object (old style, only for backward compatibility)
func AddRelationPointer(transaction *BBcTransaction, relationIdx int, refTransactionId, refAssetId *[]byte) {
	if transaction == nil {
		return
	}
	transaction.Relations[relationIdx].CreatePointer(refTransactionId, refAssetId)
}

// AddPointerInRelation creates and includes a BBcPointer object in BBcRelation (old style, only for backward compatibility)
func AddPointerInRelation(relation *BBcRelation, refTransaction *BBcTransaction, refAssetId *[]byte) {
	relation.CreatePointer(&refTransaction.TransactionID, refAssetId)
}

// AddReference creates and includes a BBcReference object in a BBcTransaction object (old style, only for backward compatibility)
func AddReference(transaction *BBcTransaction, assetGroupId *[]byte, refTransaction *BBcTransaction, eventIdx int) {
	if transaction == nil || refTransaction == nil {
		return
	}
	transaction.CreateReference(assetGroupId, refTransaction, eventIdx)
}

// addInEvent is an internal function to add a BBcEvent object in a BBcTransaction object (old style, only for backward compatibility)
func addInEvent(transaction *BBcTransaction, eventIdx int, assetGroupID, userID *[]byte) {
	ast := BBcAsset{}
	transaction.Events[eventIdx].Add(assetGroupID, &ast)
	ast.Add(userID)
}

// AddEventAssetFile sets a file digest to a BBcAsset object in a BBcEvent object and then, add it in a BBcTransaction object (old style, only for backward compatibility)
func AddEventAssetFile(transaction *BBcTransaction, eventIdx int, assetGroupID, userID *[]byte, assetFile *[]byte) {
	if transaction == nil {
		return
	}
	addInEvent(transaction, eventIdx, assetGroupID, userID)
	if assetFile != nil {
		transaction.Events[eventIdx].Asset.AddFile(assetFile)
	}
}

// AddEventAssetBodyString sets a string to a BBcAsset object in a BBcEvent object and then, add it in a BBcTransaction object (old style, only for backward compatibility)
func AddEventAssetBodyString(transaction *BBcTransaction, eventIdx int, assetGroupID, userID *[]byte, body string) {
	if transaction == nil {
		return
	}
	addInEvent(transaction, eventIdx, assetGroupID, userID)
	if body != "" {
		transaction.Events[eventIdx].Asset.AddBodyString(body)
	}
}

// AddEventAssetBodyObject sets an object (map[string]interface{}) to a BBcAsset object in a BBcEvent object and then, add it in a BBcTransaction object (old style, only for backward compatibility)
func AddEventAssetBodyObject(transaction *BBcTransaction, eventIdx int, assetGroupID, userID *[]byte, body interface{}) {
	if transaction == nil {
		return
	}
	addInEvent(transaction, eventIdx, assetGroupID, userID)
	if body != "" {
		_ = transaction.Events[eventIdx].Asset.AddBodyObject(body)
	}
}

// MakeRelationWithAsset is a utility for making simple BBcTransaction object with BBcRelation with BBcAsset (old style, only for backward compatibility)
func MakeRelationWithAsset(assetGroupId, userId *[]byte, assetBodyString string, assetBodyObject interface{}, assetFile *[]byte) *BBcRelation {
	rtn := BBcRelation{}
	rtn.SetIdLengthConf(&IdLengthConfig)
	copy(rtn.AssetGroupID, *assetGroupId)
	if assetBodyString != "" {
		rtn.CreateAsset(userId, assetFile, assetBodyString)
	} else if assetBodyObject != nil {
		rtn.CreateAsset(userId, assetFile, assetBodyObject)
	}
	return &rtn
}

// SignToTransaction signs the transaction and append the BBcSignature object to it (old style, only for backward compatibility)
func SignToTransaction(transaction *BBcTransaction, userId *[]byte, keyPair *KeyPair) {
	transaction.Sign(userId, keyPair, false)
}
