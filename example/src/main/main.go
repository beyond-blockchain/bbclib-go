package main

import (
	"bytes"
	"fmt"
	"github.com/beyond-blockchain/bbclib-go"
	"io/ioutil"
)

func main() {
	idConf := bbclib.BBcIdConfig{16, 12, 10, 8, 24}
	bbclib.ConfigureIdLength(&idConf)

	assetGroupID := bbclib.GetIdentifierWithTimestamp("assetGroupID", 32)
	u1 := bbclib.GetIdentifierWithTimestamp("user1", 32)
	u2 := bbclib.GetIdentifierWithTimestamp("user2", 32)
	keypair1 := bbclib.GenerateKeypair(bbclib.KeyTypeEcdsaP256v1, 4)
	keypair2 := bbclib.GenerateKeypair(bbclib.KeyTypeEcdsaSECP256k1, 4)

	txobj := bbclib.MakeTransaction(3, 0, true)
	bbclib.AddEventAssetBodyString(txobj, 0, &assetGroupID, &u1, "teststring!!!!!")
	txobj.Events[0].AddMandatoryApprover(&u1)
	filedat, _ := ioutil.ReadFile("./asset_test.go")
	bbclib.AddEventAssetFile(txobj, 1, &assetGroupID, &u2, &filedat)
	txobj.Events[1].AddMandatoryApprover(&u2)
	datobj := map[string]string{"param1": "aaa", "param2": "bbb", "param3": "ccc"}
	bbclib.AddEventAssetBodyObject(txobj, 2, &assetGroupID, &u1, &datobj)
	txobj.Events[2].AddMandatoryApprover(&u1)

	txobj.Witness.AddWitness(&u1)
	txobj.Witness.AddWitness(&u2)

	bbclib.SignToTransaction(txobj, &u1, &keypair1)
	bbclib.SignToTransaction(txobj, &u2, &keypair2)

	fmt.Println("-------------transaction--------------")
	fmt.Printf("%v", txobj.Stringer())
	fmt.Println("--------------------------------------")

	dat, err := txobj.Pack()
	if err != nil {
		fmt.Printf("failed to serialize transaction object (%v)", err)
	}
	fmt.Printf("Packed data: %x", dat)

	obj2 := bbclib.BBcTransaction{}
	obj2.Unpack(&dat)
	obj2.Digest()
	if result := obj2.Signatures[0].Verify(obj2.TransactionID); !result {
		fmt.Println("Verification failed..")
	}

	if bytes.Compare(txobj.Events[0].Asset.AssetID, obj2.Events[0].Asset.AssetID) != 0 ||
		bytes.Compare(txobj.TransactionID, obj2.TransactionID) != 0 {
		fmt.Println("Not recovered correctly...")
	}
}
