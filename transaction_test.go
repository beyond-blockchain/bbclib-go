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

package bbclib

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

var idLengthConfig = BBcIdConfig {
	TransactionIdLength: 32,
	UserIdLength: 32,
	AssetGroupIdLength: 32,
	AssetIdLength: 32,
	NonceLength: 32,
}

var (
	txtest_u1 = GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIDLength)
	txtest_u2 = GetIdentifierWithTimestamp("user2", defaultIDLength)
	txtest_u3 = GetIdentifierWithTimestamp("user3", defaultIDLength)
	txtest_u4 = GetIdentifierWithTimestamp("user4", defaultIDLength)
	txtest_u5 = GetIdentifierWithTimestamp("user5", defaultIDLength)
	txtest_u6 = GetIdentifierWithTimestamp("user6", defaultIDLength)
)


func makeBaseTx(idconf BBcIdConfig) BBcTransaction {
	txobj := BBcTransaction{Version: 1, Timestamp: time.Now().UnixNano()}
	txobj.SetIdLengthConf(&idconf)
	keyPair, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)

	assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
	dom := GetIdentifier("dummy domain", defaultIDLength)
	dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)
	txobj.CreateEvent(assetgroup, []int{0}).AddMandatoryApprover(&txtest_u1).AddMandatoryApprover(&txtest_u2).AddOptionParams(1, 2).AddOptionApprover(&txtest_u3).AddOptionApprover(&txtest_u4).AddAsset(&txtest_u1, nil, "testString12345XXX")
	txobj.AddCrossRef(&dom, &dummyTxid)
	txobj.AddWitness(&txtest_u1)
	txobj.AddSignature(&txtest_u1, keyPair, false)
	return txobj
}

func makeFollowTX(idconf BBcIdConfig, refTxObj *BBcTransaction) BBcTransaction {
	keypair, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
	txobj := BBcTransaction{Version: 1, Timestamp: time.Now().UnixNano()}
	txobj.SetIdLengthConf(&idconf)

	assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
	txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
	txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
	asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)
	dom := GetIdentifier("dummy domain", defaultIDLength)
	dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)

	txobj.AddCrossRef(&dom, &dummyTxid)
	txobj.CreateRelation(assetgroup).AddPointer(&txid1, &asid1).AddPointer(&txid2, nil).AddAsset(&txtest_u1, nil, "testString12345XXX")
	txobj.AddReference(&assetgroup, refTxObj, 0)
	txobj.AddWitness(&txtest_u5).AddWitness(&txtest_u6)

	txobj.AddSignature(&txtest_u1, keypair, false)
	txobj.AddSignature(&txtest_u2, keypair, false)
	txobj.AddSignature(&txtest_u4, keypair, false)
	txobj.AddSignature(&txtest_u5, keypair, false)
	txobj.AddSignature(&txtest_u6, keypair, false)
	return txobj
}

func makeFollowTXWithAssetRaw(idconf BBcIdConfig, refTxObj *BBcTransaction) BBcTransaction {
	keypair, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
	txobj := BBcTransaction{Version: 2, Timestamp: time.Now().UnixNano()}
	txobj.SetIdLengthConf(&idconf)

	assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
	asid := GetIdentifier("user1_789abcdef0123456789abcdef0", idLengthConfig.AssetIdLength)
	txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
	txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
	asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)
	dom := GetIdentifier("dummy domain", defaultIDLength)
	dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)

	txobj.AddCrossRef(&dom, &dummyTxid)
	txobj.CreateRelation(assetgroup).AddPointer(&txid1, &asid1).AddPointer(&txid2, nil).AddAssetRaw(&asid,"testString12345XXX")
	txobj.AddReference(&assetgroup, refTxObj, 0)
	txobj.AddWitness(&txtest_u5).AddWitness(&txtest_u6)

	txobj.AddSignature(&txtest_u1, keypair, false)
	txobj.AddSignature(&txtest_u2, keypair, false)
	txobj.AddSignature(&txtest_u4, keypair, false)
	txobj.AddSignature(&txtest_u5, keypair, false)
	txobj.AddSignature(&txtest_u6, keypair, false)
	return txobj
}

func makeFollowTXWithAssetHash(idconf BBcIdConfig, refTxObj *BBcTransaction) BBcTransaction {
	keypair, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
	txobj := BBcTransaction{Version: 2, Timestamp: time.Now().UnixNano()}
	txobj.SetIdLengthConf(&idconf)

	assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
	txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
	txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
	asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)
	dom := GetIdentifier("dummy domain", defaultIDLength)
	dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)

	txobj.AddCrossRef(&dom, &dummyTxid)
	txobj.CreateRelation(assetgroup).AddPointer(&txid1, &asid1).AddPointer(&txid2, nil)
	for i := 0; i < 3; i++ {
		a := GetIdentifier(fmt.Sprintf("asset_id_%d", i), idLengthConfig.AssetIdLength)
		txobj.Relations[0].AddAssetHash(&a)
	}
	txobj.AddReference(&assetgroup, refTxObj, 0)
	txobj.AddWitness(&txtest_u5).AddWitness(&txtest_u6)

	txobj.AddSignature(&txtest_u1, keypair, false)
	txobj.AddSignature(&txtest_u2, keypair, false)
	txobj.AddSignature(&txtest_u4, keypair, false)
	txobj.AddSignature(&txtest_u5, keypair, false)
	txobj.AddSignature(&txtest_u6, keypair, false)
	return txobj
}


func TestTransactionPackUnpackSimple(t *testing.T) {
	t.Run("simple creation (with relation)", func(t *testing.T) {
		keypair, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
		txobj := BBcTransaction{Version: 1, Timestamp: time.Now().UnixNano()}
		txobj.SetIdLengthConf(&idLengthConfig)

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)
		dom := GetIdentifier("dummy domain", defaultIDLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)

		txobj.AddCrossRef(&dom, &dummyTxid)
		txobj.CreateRelation(assetgroup).AddPointer(&txid1, &asid1).AddPointer(&txid2, nil).AddAsset(&txtest_u1, nil, "testString12345XXX")
		txobj.AddWitness(&txtest_u1).AddWitness(&txtest_u2)
		txobj.AddSignature(&txtest_u1, keypair, false)
		txobj.AddSignature(&txtest_u2, keypair, false)

		/*
		t.Log("---------------transaction-----------------")
		t.Logf("%v", txobj.Stringer())
		t.Log("--------------------------------------")
		 */

		dat, err := txobj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		//t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		_ = obj2.Unpack(&dat)

		/*
		t.Log("---------------transaction-----------------")
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")
		 */

		//obj2.Digest()
		if result := obj2.Signatures[0].Verify(obj2.TransactionID); !result {
			t.Fatal("Verification failed..")
		}

		if bytes.Compare(txobj.Relations[0].Asset.AssetID, obj2.Relations[0].Asset.AssetID) != 0 ||
			bytes.Compare(txobj.TransactionID, obj2.TransactionID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}

func TestTransactionPackUnpackSimpleWithEvent(t *testing.T) {
	txobj2 := makeBaseTx(idLengthConfig)

	t.Run("simple creation (with event)", func(t *testing.T) {
		/*
		t.Log("---------------transaction-----------------")
		t.Logf("%v", txobj2.Stringer())
		t.Log("--------------------------------------")
		 */

		dat, err := txobj2.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		//t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)
		/*
		t.Log("---------------transaction-----------------")
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")
		*/

		obj2.Digest()
		if result := obj2.Signatures[0].Verify(obj2.TransactionID); !result {
			t.Fatal("Verification failed..")
		}

		if bytes.Compare(txobj2.Events[0].Asset.AssetID, obj2.Events[0].Asset.AssetID) != 0 ||
			bytes.Compare(txobj2.TransactionID, obj2.TransactionID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("simple creation (with event/reference)", func(t *testing.T) {
		keypair, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
		txobj3 := BBcTransaction{Version: 1, Timestamp: time.Now().UnixNano()}
		txobj3.SetIdLengthConf(&idLengthConfig)

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		dom := GetIdentifier("dummy domain", defaultIDLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)

		txobj3.AddCrossRef(&dom, &dummyTxid)
		txobj3.CreateEvent(assetgroup, []int{0}).AddMandatoryApprover(&txtest_u1).AddMandatoryApprover(&txtest_u2).AddOptionParams(0, 0).AddAsset(&txtest_u1, nil, "testString12345XXX")
		txobj3.AddReference(&assetgroup, &txobj2, 0)
		txobj3.AddWitness(&txtest_u1).AddWitness(&txtest_u3)
		txobj3.AddSignature(&txtest_u1, keypair, false)
		txobj3.AddSignature(&txtest_u2, keypair, false)
		txobj3.AddSignature(&txtest_u2, keypair, false)
		//ref.AddSignature(&txtest_u2, &sig2)

		/*
			t.Log("---------------transaction-----------------")
			t.Logf("%v", txobj3.Stringer())
			t.Log("--------------------------------------")
		*/

		dat, err := txobj3.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		//t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)
		/*
			t.Log("---------------transaction-----------------")
			t.Logf("%v", obj2.Stringer())
			t.Log("--------------------------------------")
		*/

		obj2.Digest()
		if result := obj2.Signatures[0].Verify(obj2.TransactionID); !result {
			t.Fatal("Verification failed..signature[0]")
		}
		if result := obj2.Signatures[1].Verify(obj2.TransactionID); !result {
			t.Fatal("Verification failed..signature[1]")
		}

		if bytes.Compare(txobj3.Events[0].Asset.AssetID, obj2.Events[0].Asset.AssetID) != 0 ||
			bytes.Compare(txobj3.TransactionID, obj2.TransactionID) != 0 ||
			len(obj2.References[0].SigIndices) != 3 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("pack/unpack between creating a tx and signing to it", func(t *testing.T) {
		keypair, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
		txobj4 := BBcTransaction{Version: 1, Timestamp: time.Now().UnixNano()}
		txobj4.SetIdLengthConf(&idLengthConfig)

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)
		dom := GetIdentifier("dummy domain", defaultIDLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)

		txobj4.AddCrossRef(&dom, &dummyTxid)
		txobj4.CreateRelation(assetgroup).AddPointer(&txid1, &asid1).AddPointer(&txid2, nil).AddAsset(&txtest_u1, nil, "testString12345XXX")
		txobj4.AddReference(&assetgroup, &txobj2, 0)
		txobj4.AddWitness(&txtest_u5).AddWitness(&txtest_u6)

		/*
			t.Log("---------------transaction-----------------")
			t.Logf("%v", txobj4.Stringer())
			t.Log("--------------------------------------")
		*/
		signum := len(txobj4.Signatures)
		if signum != 5 {
			t.Fatal("Invalid number of signatures")
		}

		dat, err := txobj4.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		//t.Logf("Packed data: %x", dat)

		obj4 := BBcTransaction{}
		obj4.Unpack(&dat)
		/*
		t.Log("---------------transaction-----------------")
		t.Logf("%v", obj4.Stringer())
		t.Log("--------------------------------------")
		 */
		signum = len(obj4.Signatures)
		if signum != 5 {
			t.Fatal("Invalid number of signatures")
		}
		refObj := obj4.References[0]
		refObj.Add(nil, &txobj2, -1)

		obj4.AddSignature(&txtest_u1, keypair, false)
		obj4.AddSignature(&txtest_u2, keypair, false)
		obj4.AddSignature(&txtest_u4, keypair, false)  // [5]
		obj4.AddSignature(&txtest_u5, keypair, false)
		obj4.AddSignature(&txtest_u6, keypair, false)
		/*
		t.Log("---------------transaction-----------------")
		t.Logf("%v", obj4.Stringer())
		t.Log("--------------------------------------")
		*/

		signum = len(obj4.Signatures)
		if signum != 5 {
			t.Fatal("Invalid number of signatures")
		}

		d1 := txobj4.Digest()
		d2 := obj4.Digest()
		if bytes.Compare(d1, d2) != 0 {
			t.Fatal("transaction_id mismatch")
		}

		if ret, i := txobj4.VerifyAll(); !ret {
			t.Fatalf("Invalid signature at idx=%d\n", i)
		}
		if ret, i := obj4.VerifyAll(); !ret {
			t.Fatalf("Invalid signature at idx=%d\n", i)
		}

		dat2_1, err := txobj4.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		dat2_2, err := obj4.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		if bytes.Compare(dat2_1, dat2_2) != 0 {
			t.Log("The order of signatures are different (no problem)")
		}
	})
}

func TestTransactionWithVariousIdLength(t *testing.T) {
	t.Run("transaction_id length is 10)", func(t *testing.T) {
		idconf := idLengthConfig
		idconf.TransactionIdLength = 10
		txprev := makeBaseTx(idconf)
		txobj := makeFollowTX(idconf, &txprev)
		if txobj.TransactionIdLength != 10 {
			t.Fatalf("Invalid transaction_id length field (%d != 10) idconf=%v\n", txobj.TransactionIdLength, idconf)
		}
		if len(txobj.Relations[0].Pointers[0].TransactionID) != 10 {
			t.Fatal("Invalid transaction_id length (!= 10)")
		}
	})

	t.Run("user_id length is 10)", func(t *testing.T) {
		idconf := idLengthConfig
		idconf.UserIdLength = 10
		txprev := makeBaseTx(idconf)
		txobj := makeFollowTX(idconf, &txprev)
		if len(txobj.Witness.UserIDs[0]) != 10 {
			t.Fatalf("Invalid user_id length in Witness (%d != 10) idconf=%v\n", len(txobj.Witness.UserIDs[0]), idconf)
		}
		if len(txobj.Relations[0].Asset.UserID) != 10 {
			t.Fatalf("Invalid user_id length in Asset of Relation (%d != 10) idconf=%v\n", len(txobj.Relations[0].Asset.UserID), idconf)
		}
		if len(txprev.Events[0].Asset.UserID) != 10 {
			t.Fatalf("Invalid user_id length in Asset of Event (%d != 10) idconf=%v\n", len(txprev.Events[0].Asset.UserID), idconf)
		}
	})

	t.Run("asset_group_id length is 10)", func(t *testing.T) {
		idconf := idLengthConfig
		idconf.AssetGroupIdLength = 10
		txprev := makeBaseTx(idconf)
		txobj := makeFollowTX(idconf, &txprev)
		if len(txobj.Relations[0].AssetGroupID) != 10 {
			t.Fatalf("Invalid asset_group_id length in Relation (%d != 10) idconf=%v\n", len(txobj.Relations[0].AssetGroupID), idconf)
		}
		if len(txprev.Events[0].AssetGroupID) != 10 {
			t.Fatalf("Invalid asset_group_id length in Event (%d != 10) idconf=%v\n", len(txprev.Events[0].AssetGroupID), idconf)
		}
	})

	t.Run("asset_id length is 10)", func(t *testing.T) {
		idconf := idLengthConfig
		idconf.AssetIdLength = 10
		txprev := makeBaseTx(idconf)
		txobj := makeFollowTX(idconf, &txprev)
		if len(txobj.Relations[0].Asset.AssetID) != 10 {
			t.Fatalf("Invalid asset_id length in Asset of Relation (%d != 10) idconf=%v\n", len(txobj.Relations[0].Asset.AssetID), idconf)
		}
		if len(txprev.Events[0].Asset.AssetID) != 10 {
			t.Fatalf("Invalid asset_id length in Asset of Event (%d != 10) idconf=%v\n", len(txprev.Events[0].Asset.AssetID), idconf)
		}
	})

	t.Run("nonce length is 10)", func(t *testing.T) {
		idconf := idLengthConfig
		idconf.NonceLength = 10
		txprev := makeBaseTx(idconf)
		txobj := makeFollowTX(idconf, &txprev)
		if len(txobj.Relations[0].Asset.Nonce) != 10 {
			t.Fatalf("Invalid nonce length in Asset of Relation (%d != 10) idconf=%v\n", len(txobj.Relations[0].Asset.AssetID), idconf)
		}
		if len(txprev.Events[0].Asset.Nonce) != 10 {
			t.Fatalf("Invalid nonce length in Asset of Event (%d != 10) idconf=%v\n", len(txprev.Events[0].Asset.AssetID), idconf)
		}
	})

	t.Run("all lengths are 10)", func(t *testing.T) {
		idconf := BBcIdConfig{10,10,10,10,10}
		txprev := makeBaseTx(idconf)
		txobj := makeFollowTX(idconf, &txprev)
		if txobj.TransactionIdLength != 10 {
			t.Fatalf("Invalid transaction_id length field (%d != 10) idconf=%v\n", txobj.TransactionIdLength, idconf)
		}
		if len(txobj.Relations[0].Pointers[0].TransactionID) != 10 {
			t.Fatal("Invalid transaction_id length (!= 10)")
		}
		if len(txobj.Witness.UserIDs[0]) != 10 {
			t.Fatalf("Invalid user_id length in Witness (%d != 10) idconf=%v\n", len(txobj.Witness.UserIDs[0]), idconf)
		}
		if len(txobj.Relations[0].Asset.UserID) != 10 {
			t.Fatalf("Invalid user_id length in Asset of Relation (%d != 10) idconf=%v\n", len(txobj.Relations[0].Asset.UserID), idconf)
		}
		if len(txprev.Events[0].Asset.UserID) != 10 {
			t.Fatalf("Invalid user_id length in Asset of Event (%d != 10) idconf=%v\n", len(txprev.Events[0].Asset.UserID), idconf)
		}
		if len(txobj.Relations[0].AssetGroupID) != 10 {
			t.Fatalf("Invalid asset_group_id length in Relation (%d != 10) idconf=%v\n", len(txobj.Relations[0].AssetGroupID), idconf)
		}
		if len(txprev.Events[0].AssetGroupID) != 10 {
			t.Fatalf("Invalid asset_group_id length in Event (%d != 10) idconf=%v\n", len(txprev.Events[0].AssetGroupID), idconf)
		}
		if len(txobj.Relations[0].Asset.AssetID) != 10 {
			t.Fatalf("Invalid asset_id length in Asset of Relation (%d != 10) idconf=%v\n", len(txobj.Relations[0].Asset.AssetID), idconf)
		}
		if len(txprev.Events[0].Asset.AssetID) != 10 {
			t.Fatalf("Invalid asset_id length in Asset of Event (%d != 10) idconf=%v\n", len(txprev.Events[0].Asset.AssetID), idconf)
		}
		if len(txobj.Relations[0].Asset.Nonce) != 10 {
			t.Fatalf("Invalid nonce length in Asset of Relation (%d != 10) idconf=%v\n", len(txobj.Relations[0].Asset.AssetID), idconf)
		}
		if len(txprev.Events[0].Asset.Nonce) != 10 {
			t.Fatalf("Invalid nonce length in Asset of Event (%d != 10) idconf=%v\n", len(txprev.Events[0].Asset.AssetID), idconf)
		}
	})
}

func TestTransactionWithAssetRawAndAssetHash(t *testing.T) {
	t.Run("transaction with BBcAssetRaw", func(t *testing.T) {
		idconf := idLengthConfig
		idconf.TransactionIdLength = 10
		txprev := makeBaseTx(idconf)
		txobj := makeFollowTXWithAssetRaw(idconf, &txprev)

		dat, err := txobj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		//t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)

		d1 := txobj.Digest()
		d2 := obj2.Digest()
		if bytes.Compare(d1, d2) != 0 {
			t.Fatal("transaction_id mismatch")
		}

		result, _ := obj2.VerifyAll()
		if !result {
			t.Fatal("Verification failed..")
		}

		if bytes.Compare(txobj.Relations[0].AssetRaw.AssetID, obj2.Relations[0].AssetRaw.AssetID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
		if bytes.Compare(txobj.Relations[0].AssetRaw.AssetBody, obj2.Relations[0].AssetRaw.AssetBody) != 0 {
			t.Fatal("Not recovered correctly...")
		}
		if bytes.Compare(txobj.Relations[0].Pointers[0].TransactionID, obj2.Relations[0].Pointers[0].TransactionID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
		if bytes.Compare(txobj.Witness.UserIDs[0], obj2.Witness.UserIDs[0]) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("transaction with BBcAssetHash", func(t *testing.T) {
		idconf := idLengthConfig
		idconf.TransactionIdLength = 10
		txprev := makeBaseTx(idconf)
		txobj := makeFollowTXWithAssetHash(idconf, &txprev)

		dat, err := txobj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		//t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)

		d1 := txobj.Digest()
		d2 := obj2.Digest()
		/*
			t.Log("---------------transaction-----------------")
			t.Logf("%v", txobj.Stringer())
			t.Log("--------------------------------------")
		*/
		if bytes.Compare(d1, d2) != 0 {
			t.Fatal("transaction_id mismatch")
		}

		result, _ := obj2.VerifyAll()
		if !result {
			t.Fatal("Verification failed..")
		}

		for i := 0; i < 3; i++ {
			if bytes.Compare(txobj.Relations[0].AssetHash.AssetIDs[i], obj2.Relations[0].AssetHash.AssetIDs[i]) != 0 {
				t.Fatal("Not recovered correctly...")
			}
		}
		if bytes.Compare(txobj.Relations[0].Pointers[0].TransactionID, obj2.Relations[0].Pointers[0].TransactionID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
		if bytes.Compare(txobj.Witness.UserIDs[0], obj2.Witness.UserIDs[0]) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}
