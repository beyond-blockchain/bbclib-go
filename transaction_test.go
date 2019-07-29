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
	u1 = GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIDLength)
	u2 = GetIdentifierWithTimestamp("user2", defaultIDLength)
	u3 = GetIdentifierWithTimestamp("user3", defaultIDLength)
	u4 = GetIdentifierWithTimestamp("user4", defaultIDLength)
	u5 = GetIdentifierWithTimestamp("user5", defaultIDLength)
	u6 = GetIdentifierWithTimestamp("user6", defaultIDLength)
)


func makeBaseTx(idconf BBcIdConfig) BBcTransaction {
	txobj := BBcTransaction{Version: 1, Timestamp: time.Now().UnixNano()}
	txobj.SetIdLengthConf(&idconf)
	keypair := GenerateKeypair(KeyTypeEcdsaP256v1, defaultCompressionMode)
	evt := BBcEvent{}
	txobj.AddEvent(&evt)
	crs := BBcCrossRef{}
	txobj.AddCrossRef(&crs)
	wit := BBcWitness{}
	txobj.AddWitness(&wit)

	ast := BBcAsset{}
	ast.SetIdLengthConf(&idconf)
	ast.Add(&u1)
	ast.AddBodyString("testString12345XXX")

	assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
	evt.Add(&assetgroup, &ast)

	evt.AddMandatoryApprover(&u1)
	evt.AddMandatoryApprover(&u2)
	evt.AddOptionParams(1, 2)
	evt.AddOptionApprover(&u3)
	evt.AddOptionApprover(&u4)

	dom := GetIdentifier("dummy domain", defaultIDLength)
	dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)
	crs.Add(&dom, &dummyTxid)

	wit.AddWitness(&u1)
	sig := BBcSignature{}
	sig.SetPublicKeyByKeypair(&keypair)
	signature, _ := txobj.Sign(&keypair)
	sig.SetSignature(&signature)
	wit.AddSignature(&u1, &sig)

	return txobj
}

func makeFollowTX(idconf BBcIdConfig, refTxObj *BBcTransaction) BBcTransaction {
	keypair := GenerateKeypair(KeyTypeEcdsaP256v1, defaultCompressionMode)
	txobj := BBcTransaction{Version: 1, Timestamp: time.Now().UnixNano()}
	txobj.SetIdLengthConf(&idconf)
	rtn := BBcRelation{}
	txobj.AddRelation(&rtn)
	ref := BBcReference{}
	txobj.AddReference(&ref)
	wit := BBcWitness{}
	txobj.AddWitness(&wit)
	crs := BBcCrossRef{}
	txobj.AddCrossRef(&crs)

	ast := BBcAsset{}
	ast.SetIdLengthConf(&idconf)
	ptr1 := BBcPointer{}
	ptr1.SetIdLengthConf(&idconf)
	ptr2 := BBcPointer{}
	ptr2.SetIdLengthConf(&idconf)

	assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
	rtn.Add(&assetgroup, &ast)
	rtn.AddPointer(&ptr1)
	rtn.AddPointer(&ptr2)

	ast.Add(&u1)
	ast.AddBodyString("testString12345XXX")

	txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
	txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
	asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)
	ptr1.Add(&txid1, &asid1)
	ptr2.Add(&txid2, nil)

	wit.AddWitness(&u5)
	ref.Add(&assetgroup, refTxObj, 0)
	wit.AddWitness(&u6)

	dom := GetIdentifier("dummy domain", defaultIDLength)
	dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)
	crs.Add(&dom, &dummyTxid)

	sig := BBcSignature{}
	sig.SetPublicKeyByKeypair(&keypair)
	signature, _ := txobj.Sign(&keypair)
	sig.SetSignature(&signature)
	ref.AddSignature(&u1, &sig)

	sig6 := BBcSignature{}
	sig6.SetPublicKeyByKeypair(&keypair)
	signature6, _ := txobj.Sign(&keypair)
	sig6.SetSignature(&signature6)
	txobj.AddSignature(&u6, &sig6)

	sig2 := BBcSignature{}
	sig2.SetPublicKeyByKeypair(&keypair)
	signature2, _ := txobj.Sign(&keypair)
	sig2.SetSignature(&signature2)
	ref.AddSignature(&u2, &sig2)

	sig5 := BBcSignature{}
	sig5.SetPublicKeyByKeypair(&keypair)
	signature5, _ := txobj.Sign(&keypair)
	sig5.SetSignature(&signature5)
	txobj.AddSignature(&u5, &sig5)

	sig4 := BBcSignature{}
	sig4.SetPublicKeyByKeypair(&keypair)
	signature3, _ := txobj.Sign(&keypair)
	sig4.SetSignature(&signature3)
	ref.AddSignature(&u4, &sig4)

	return txobj
}


func TestTransactionPackUnpackSimple(t *testing.T) {
	t.Run("simple creation (with relation)", func(t *testing.T) {
		keypair := GenerateKeypair(KeyTypeEcdsaP256v1, defaultCompressionMode)
		txobj := BBcTransaction{Version: 1, Timestamp: time.Now().UnixNano()}
		txobj.SetIdLengthConf(&idLengthConfig)
		rtn := BBcRelation{}
		txobj.AddRelation(&rtn)
		wit := BBcWitness{}
		txobj.AddWitness(&wit)
		crs := BBcCrossRef{}
		txobj.AddCrossRef(&crs)

		ast := BBcAsset{}
		ptr1 := BBcPointer{}
		ptr2 := BBcPointer{}

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		rtn.Add(&assetgroup, &ast)
		rtn.AddPointer(&ptr1)
		rtn.AddPointer(&ptr2)

		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")

		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)
		ptr1.Add(&txid1, &asid1)
		ptr2.Add(&txid2, nil)

		wit.AddWitness(&u1)
		wit.AddWitness(&u2)

		dom := GetIdentifier("dummy domain", defaultIDLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)
		crs.Add(&dom, &dummyTxid)

		sig := BBcSignature{}
		sig.SetPublicKeyByKeypair(&keypair)
		signature, _ := txobj.Sign(&keypair)
		sig.SetSignature(&signature)
		wit.AddSignature(&u1, &sig)

		sig2 := BBcSignature{}
		sig2.SetPublicKeyByKeypair(&keypair)
		signature2, _ := txobj.Sign(&keypair)
		sig2.SetSignature(&signature2)
		wit.AddSignature(&u2, &sig2)

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
		obj2.SetIdLengthConf(&idLengthConfig)
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
		obj2.SetIdLengthConf(&idLengthConfig)
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
		keypair := GenerateKeypair(KeyTypeEcdsaP256v1, defaultCompressionMode)
		txobj3 := BBcTransaction{Version: 1, Timestamp: time.Now().UnixNano()}
		txobj3.SetIdLengthConf(&idLengthConfig)
		evt := BBcEvent{}
		txobj3.AddEvent(&evt)
		ref := BBcReference{}
		txobj3.AddReference(&ref)
		crs := BBcCrossRef{}
		txobj3.AddCrossRef(&crs)

		ast := BBcAsset{}
		ast.SetIdLengthConf(&idLengthConfig)
		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		evt.Add(&assetgroup, &ast)

		evt.AddMandatoryApprover(&u1)
		evt.AddOptionParams(0, 0)

		ref.Add(&assetgroup, &txobj2, 0)

		dom := GetIdentifier("dummy domain", defaultIDLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)
		crs.Add(&dom, &dummyTxid)

		sig := BBcSignature{}
		sig.SetPublicKeyByKeypair(&keypair)
		signature, err := txobj3.Sign(&keypair)
		if err != nil {
			t.Fatal(err)
		}
		sig.SetSignature(&signature)
		ref.AddSignature(&u1, &sig)

		sig2 := BBcSignature{}
		sig2.SetPublicKeyByKeypair(&keypair)
		signature2, err := txobj3.Sign(&keypair)
		if err != nil {
			t.Fatal(err)
		}
		sig2.SetSignature(&signature2)
		ref.AddSignature(&u2, &sig2)

		sig3 := BBcSignature{}
		sig3.SetPublicKeyByKeypair(&keypair)
		signature3, err := txobj3.Sign(&keypair)
		if err != nil {
			t.Fatal(err)
		}
		sig3.SetSignature(&signature3)
		ref.AddSignature(&u4, &sig3)

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
		obj2.SetIdLengthConf(&idLengthConfig)
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
		keypair := GenerateKeypair(KeyTypeEcdsaP256v1, defaultCompressionMode)
		txobj4 := BBcTransaction{Version: 1, Timestamp: time.Now().UnixNano()}
		txobj4.SetIdLengthConf(&idLengthConfig)
		rtn := BBcRelation{}
		txobj4.AddRelation(&rtn)
		ref := BBcReference{}
		txobj4.AddReference(&ref)
		wit := BBcWitness{}
		txobj4.AddWitness(&wit)
		crs := BBcCrossRef{}
		txobj4.AddCrossRef(&crs)

		ast := BBcAsset{}
		ast.SetIdLengthConf(&idLengthConfig)
		ptr1 := BBcPointer{}
		ptr1.SetIdLengthConf(&idLengthConfig)
		ptr2 := BBcPointer{}
		ptr2.SetIdLengthConf(&idLengthConfig)

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		rtn.Add(&assetgroup, &ast)
		rtn.AddPointer(&ptr1)
		rtn.AddPointer(&ptr2)

		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")

		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)
		ptr1.Add(&txid1, &asid1)
		ptr2.Add(&txid2, nil)

		wit.AddWitness(&u5)
		ref.Add(&assetgroup, &txobj2, 0)
		wit.AddWitness(&u6)

		dom := GetIdentifier("dummy domain", defaultIDLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)
		crs.Add(&dom, &dummyTxid)

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
		obj4.SetIdLengthConf(&idLengthConfig)
		obj4.Unpack(&dat)
		t.Log("---------------transaction-----------------")
		t.Logf("%v", obj4.Stringer())
		t.Log("--------------------------------------")
		signum = len(obj4.Signatures)
		if signum != 5 {
			t.Fatal("Invalid number of signatures")
		}
		refObj := obj4.References[0]
		refObj.Add(nil, &txobj2, -1)

		sig := BBcSignature{}
		sig.SetPublicKeyByKeypair(&keypair)
		signature, err := obj4.Sign(&keypair)
		if err != nil {
			t.Fatal(err)
		}
		sig.SetSignature(&signature)
		refObj.AddSignature(&u1, &sig)
		ref.AddSignature(&u1, &sig)

		sig6 := BBcSignature{}
		sig6.SetPublicKeyByKeypair(&keypair)
		signature6, err := obj4.Sign(&keypair)
		if err != nil {
			t.Fatal(err)
		}
		sig6.SetSignature(&signature6)
		obj4.AddSignature(&u6, &sig6)
		txobj4.AddSignature(&u6, &sig6)

		sig2 := BBcSignature{}
		sig2.SetPublicKeyByKeypair(&keypair)
		signature2, err := obj4.Sign(&keypair)
		if err != nil {
			t.Fatal(err)
		}
		sig2.SetSignature(&signature2)
		refObj.AddSignature(&u2, &sig2)
		ref.AddSignature(&u2, &sig2)

		sig5 := BBcSignature{}
		sig5.SetPublicKeyByKeypair(&keypair)
		signature5, err := obj4.Sign(&keypair)
		if err != nil {
			t.Fatal(err)
		}
		sig5.SetSignature(&signature5)
		obj4.AddSignature(&u5, &sig5)
		txobj4.AddSignature(&u5, &sig5)

		sig4 := BBcSignature{}
		sig4.SetPublicKeyByKeypair(&keypair)
		signature3, err := obj4.Sign(&keypair)
		if err != nil {
			t.Fatal(err)
		}
		sig4.SetSignature(&signature3)
		refObj.AddSignature(&u4, &sig4)
		ref.AddSignature(&u4, &sig4)

		t.Log("---------------transaction-----------------")
		t.Logf("%v", obj4.Stringer())
		t.Log("--------------------------------------")
		signum = len(obj4.Signatures)
		if signum != 5 {
			t.Fatal("Invalid number of signatures")
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
			t.Fatal("Not created correctly...")
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

