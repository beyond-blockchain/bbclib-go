package bbclib

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"testing"
)

// Serialized data from Python bbclib.py
var (
	assetGroupIDInTx = "5464b9653aa0100abd0dd1d402e80e0de7f21f5f23d890a83585291115a90a08"
	txdataEventRef   = "00000100000021dd035c0000000020000100ca00000020005464b9653aa0100abd0dd1d402e80e0de7f21f5f23d890a83585291115a90a080000010020009048feaeaf902a66879be3f0ee2e30a981df641b074f1fa901649002a9d065b2000000007a0000002000de36cf0094a8a7a80b4552de38d7d5de490086d60f395b468e937e1d8b9d95d020009048feaeaf902a66879be3f0ee2e30a981df641b074f1fa901649002a9d065b22000f7e4d7c82687e579662c69e952d22b26ea73eded26c363f7f0d68da8e5c500230000000000000c006576656e745f61737365743202004a00000020005464b9653aa0100abd0dd1d402e80e0de7f21f5f23d890a83585291115a90a082000573b5b63d6c7333f12ebff55330f2e06147438c633219f40c4de9688af3de3ef0000010000004a00000020005464b9653aa0100abd0dd1d402e80e0de7f21f5f23d890a83585291115a90a082000573b5b63d6c7333f12ebff55330f2e06147438c633219f40c4de9688af3de3ef01000100010000000100020000000000000002008d000000020000000802000004a8309fa78e3a9025668f82b4e07c7324693ed5b2c4fe65506c861189b53df39c75eb874b7de6773dd41a801357d3b7cca21ba5b189e9a4e5d262b77d1dc3a5c400020000924e10d1cfff15b0e28a25ebf2700392112beeb9abb137d8e06dc1443354c24a45355d1eb288c851848da9dc99b828526bd852d2fc528b9c5f3ae2c5417c808f8d000000020000000802000004862f5a212ab0db12d10e19f07a18a40248ac90f320061c27ff6f7cb87a0be8e2a231daf61077c2ec37dd9eee6e961e0fd7ca09fa965f62a7c39b7ce84821dc4500020000f43f16b5db01fdd0a33d3d5d9abf2cca9b2cb1bde5be4735faa935a6d3b77b3877607e538b75b1c09df1271958b5717d979d63cbe9e38d4b8f67254f961550f1"
	txidEventRef     = "667fd62ae54dd91e1138006d9d7cf9b4c11d27b297d3effb8e8fc1957fda1c4f"
	txdataRelation   = "1000789c6364606050bccb1c03a4181440040313c329303b246567aad50201aebdbc17af30bde0e37dfe493e5ef9c68415a6ad9a82a22bb93818195480aac2ada393af1d37b6177afd3fd4985f8f4da4c4e298b1e27c8723f7a675acb77dfc9e81a1066cdab5effb9397a4276c75fb7f71d794d72c9d192ecf59531accf535b8ff5df67df13d58816182c7bf75eb2768a5b5cf7efce19d9ec1cac6fb29d2ecfef22b19532630adbc90ba4981c1eaf0f93373fdcd62435664b3b8bace6e6239da2bcb72e8874579bd4df659311169b0f319f8188a5273124b32f3f3e2138b8b534b5e11e91b37a0aab4fa6b5a4f7d6fca095a30e4ceadf9b9e5a0acfaa6e997dfffeeeb3f38b5fe968c3f2350cd3db3f30c53562c5fc1ed1a74cfe2fad57b9e0c6dd7f82da3ddfa26d7c976cf9d7aa10a6cdfceb46793df2fc8f791b03d23afc4ba523cf1e704860235ff68e3cf73671aaf5868a7c0e0c073f2d005b93787231eddcfcd083b71bce19848a3cb6effee7cafda9d7398fa2e2b306c9fbfc9c0b240e4b62db7e57b5b852d1b4c5e4d9ac8cdf34337e2dd97836989896220df3232f03034a52c4e4c4a3e736271446414238317380e098727c89584dcc00836ab174c3230700009961506f397f7594d504deb6fdaf2a0a65825d3eeeaa623ff520372da043bb7da7e9e53fabaddbbf659b9ed15a906e1f0cbdbcf2c925ebab1f3e592a79792b6d7ca1e5e7a046492f1950d82731b576f49b8d4cae1b13c73fe79f6c99a0a7cea8b8dae1fe6bd6fb5b6c9329a57bddc2e9dafadf2d3a7f495b2ce515e3b4ef69e9cfa62a5ad48969765342b8a83daf4a314b536dc16bac827f9a14a620993c79a099f15d864d4ffe7d7eca8e27ef16891e1ad6f02e587de98df9df72e6f9a1cfff5539cbfa6c5272d3f3cbbe68587e21d57904962f1eb13255687ec133bb1554665e299c987b5ab676e6eaa9b77fb8d9f58e3bf8483bf5a63139994040feb6fb939efdc6497ed194fb38e447baceeb78893e0af0fd9c40300c7b55581"
	txidRelation     = "c390caecc3a4e46dc7f45db9fc4d56373d33dfe2f2692075f7e2f79e348915db"
	txobj            *BBcTransaction
	txobj2           *BBcTransaction
	txobj3           *BBcTransaction
	assetGroupID     []byte
	u1               []byte
	u2               []byte
	keypair1         KeyPair
	keypair2         KeyPair
)

func TestBBcLibUtilitiesTx1(t *testing.T) {
	assetGroupID = GetIdentifierWithTimestamp("assetGroupID", defaultIDLength)
	u1 = GetIdentifierWithTimestamp("user1", defaultIDLength)
	u2 = GetIdentifierWithTimestamp("user2", defaultIDLength)
	keypair1 = GenerateKeypair(KeyTypeEcdsaP256v1, defaultCompressionMode)
	keypair2 = GenerateKeypair(KeyTypeEcdsaSECP256k1, defaultCompressionMode)

	t.Run("MakeTransaction and events", func(t *testing.T) {
		txobj = MakeTransaction(3, 0, true, 32)
		AddEventAssetBodyString(txobj, 0, &assetGroupID, &u1, "teststring!!!!!")
		txobj.Events[0].AddMandatoryApprover(&u1)
		filedat, _ := ioutil.ReadFile("./asset_test.go")
		AddEventAssetFile(txobj, 1, &assetGroupID, &u2, &filedat)
		txobj.Events[1].AddMandatoryApprover(&u2)
		datobj := map[string]string{"param1": "aaa", "param2": "bbb", "param3": "ccc"}
		AddEventAssetBodyObject(txobj, 2, &assetGroupID, &u1, &datobj)
		txobj.Events[2].AddMandatoryApprover(&u1)

		txobj.Witness.AddWitness(&u1)
		txobj.Witness.AddWitness(&u2)

		SignToTransaction(txobj, &u1, &keypair1)
		SignToTransaction(txobj, &u2, &keypair2)

		t.Log("-------------transaction--------------")
		t.Logf("%v", txobj.Stringer())
		t.Log("--------------------------------------")

		dat, err := txobj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)
		obj2.Digest()
		if result := obj2.Signatures[0].Verify(obj2.TransactionID); !result {
			t.Fatal("Verification failed..")
		}

		if bytes.Compare(txobj.Events[0].Asset.AssetID, obj2.Events[0].Asset.AssetID) != 0 ||
			bytes.Compare(txobj.TransactionID, obj2.TransactionID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}

func TestBBcLibUtilitiesTx2(t *testing.T) {
	t.Run("MakeTransaction and events/reference", func(t *testing.T) {
		txobj2 = MakeTransaction(2, 0, true, 32)
		AddEventAssetBodyString(txobj2, 0, &assetGroupID, &u1, "teststring!!!!!")
		filedat, _ := ioutil.ReadFile("./crossref_test.go")
		AddEventAssetFile(txobj2, 1, &assetGroupID, &u2, &filedat)

		AddReference(txobj2, &assetGroupID, txobj, 0)
		txobj2.References[0].AddApprover(&u1)
		AddReference(txobj2, &assetGroupID, txobj, 1)
		txobj2.References[1].AddApprover(&u2)

		crs := BBcCrossRef{}
		txobj2.AddCrossRef(&crs)
		dom := GetIdentifier("dummy domain", defaultIDLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIDLength)
		crs.Add(&dom, &dummyTxid)

		SignToTransaction(txobj2, &u1, &keypair1)
		SignToTransaction(txobj2, &u2, &keypair2)

		t.Log("-------------transaction--------------")
		t.Logf("%v", txobj2.Stringer())
		t.Log("--------------------------------------")

		dat, err := txobj2.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)
		obj2.Digest()
		if result := obj2.Signatures[0].Verify(obj2.TransactionID); !result {
			t.Fatal("Verification failed..")
		}

		if bytes.Compare(txobj2.Events[0].Asset.AssetID, obj2.Events[0].Asset.AssetID) != 0 ||
			bytes.Compare(txobj2.TransactionID, obj2.TransactionID) != 0 ||
			len(obj2.References[0].SigIndices) != 1 || len(obj2.References[1].SigIndices) != 1 {
			t.Fatal("Not recovered correctly...")
		}
	})
}
func TestBBcLibUtilitiesTx3(t *testing.T) {
	t.Run("MakeTransaction and relations", func(t *testing.T) {
		txobj3 = MakeTransaction(0, 3, true, 32)
		AddRelationAssetBodyString(txobj3, 0, &assetGroupID, &u1, "teststring!!!!!")
		filedat, _ := ioutil.ReadFile("./crossref_test.go")
		AddRelationAssetFile(txobj3, 1, &assetGroupID, &u2, &filedat)
		datobj := map[string]string{"param1": "aaa", "param2": "bbb", "param3": "ccc"}
		AddRelationAssetBodyObject(txobj3, 2, &assetGroupID, &u1, &datobj)

		datobj2 := map[string]string{"param1": "lll", "param2": "gggg", "param3": "ddd"}
		rtn := MakeRelationWithAsset(&assetGroupID, &u2, "", &datobj2, nil, 32)
		txobj3.AddRelation(rtn)

		AddRelationPointer(txobj3, 0, &txobj.TransactionID, nil)
		AddRelationPointer(txobj3, 1, &txobj2.TransactionID, &txobj2.Events[0].Asset.AssetID)
		AddPointerInRelation(rtn, txobj, &txobj.Events[2].Asset.AssetID)

		txobj3.Witness.AddWitness(&u1)
		txobj3.Witness.AddWitness(&u2)

		SignToTransaction(txobj3, &u1, &keypair1)
		SignToTransaction(txobj3, &u2, &keypair2)

		t.Log("-------------transaction--------------")
		t.Logf("%v", txobj3.Stringer())
		t.Log("--------------------------------------")

		dat, err := txobj3.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)
		obj2.Digest()

		result, _ := obj2.VerifyAll()
		if !result {
			t.Fatal("Verification failed..")
		}

		if bytes.Compare(txobj3.Relations[0].Asset.AssetID, obj2.Relations[0].Asset.AssetID) != 0 ||
			bytes.Compare(txobj3.TransactionID, obj2.TransactionID) != 0 ||
			len(txobj3.Witness.SigIndices) != 2 || len(obj2.Witness.SigIndices) != 2 {
			t.Fatal("Not recovered correctly...")
		}
	})
}

func TestBBcLibSerializeDeserialize(t *testing.T) {
	t.Run("simple serialize and deserialize", func(t *testing.T) {
		dat, err := Serialize(txobj2, FormatPlain)
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Serialized data: %x", dat)
		t.Logf("Serialized data size: %d", len(dat))

		obj2, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IDLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(txobj2.TransactionID, obj2.TransactionID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("serialize and deserialize with zlib", func(t *testing.T) {
		dat, err := Serialize(txobj3, FormatZlib)
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Serialized data: %x", dat)
		t.Logf("Serialized data size: %d", len(dat))

		obj2, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IDLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(txobj3.TransactionID, obj2.TransactionID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}

func TestBBcLibSerializeDeserializePythonData(t *testing.T) {
	t.Run("deserialize txdata genarated (type0) by python bbclib", func(t *testing.T) {
		dat, _ := hex.DecodeString(txdataEventRef)
		txobj4, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", txobj4.IDLength)
		t.Logf("%v", txobj4.Stringer())
		t.Log("--------------------------------------")

		txidOrg, _ := hex.DecodeString(txidEventRef)
		if bytes.Compare(txobj4.TransactionID, txidOrg) != 0 {
			t.Fatal("Not recovered correctly...1")
		}
		asgidOrg, _ := hex.DecodeString(assetGroupIDInTx)
		if bytes.Compare(txobj4.Events[0].AssetGroupID, asgidOrg) != 0 {
			t.Fatal("Not recovered correctly...2")
		}
	})

	t.Run("deserialize txdata genarated (type0x0010) by python bbclib", func(t *testing.T) {
		dat, _ := hex.DecodeString(txdataRelation)
		txobj5, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", txobj5.IDLength)
		t.Logf("%v", txobj5.Stringer())
		t.Log("--------------------------------------")

		txidOrg, _ := hex.DecodeString(txidRelation)
		if bytes.Compare(txobj5.TransactionID, txidOrg) != 0 {
			t.Fatal("Not recovered correctly...1")
		}
		asgidOrg, _ := hex.DecodeString(assetGroupIDInTx)
		if bytes.Compare(txobj5.Relations[0].AssetGroupID, asgidOrg) != 0 {
			t.Fatal("Not recovered correctly...2")
		}
	})
}
