package main

import (
	"fmt"
	"github.com/beyond-blockchain/bbclib-go"
	//"main/bbclib"
	"io/ioutil"
	"time"
)


func makeTransaction(num int) *bbclib.BBcTransaction {
	assetGroupID := bbclib.GetIdentifierWithTimestamp("assetGroupID", 32)
	u1 := bbclib.GetIdentifierWithTimestamp("user1", 32)
	u2 := bbclib.GetIdentifierWithTimestamp("user2", 32)
	keypair1, _ := bbclib.GenerateKeypair(bbclib.KeyTypeEcdsaP256v1, 4)
	keypair2, _ := bbclib.GenerateKeypair(bbclib.KeyTypeEcdsaP256v1, 4)
	//keypair1 := bbclib.GenerateKeypair(bbclib.KeyTypeEcdsaP256v1, 4)
	//keypair2 := bbclib.GenerateKeypair(bbclib.KeyTypeEcdsaP256v1, 4)

	txobj := bbclib.MakeTransaction(3, 0, true)
	bodyString := fmt.Sprintf("teststring!!!!:%d", num)
	bbclib.AddEventAssetBodyString(txobj, 0, &assetGroupID, &u1, bodyString)
	txobj.Events[0].AddMandatoryApprover(&u1)
	filedat, _ := ioutil.ReadFile("./asset_test.go")
	bbclib.AddEventAssetFile(txobj, 1, &assetGroupID, &u2, &filedat)
	txobj.Events[1].AddMandatoryApprover(&u2)
	datobj := map[string]string{"param1": "aaa", "param2": "bbb", "param3": string(num)}
	bbclib.AddEventAssetBodyObject(txobj, 2, &assetGroupID, &u1, &datobj)
	txobj.Events[2].AddMandatoryApprover(&u1)

	txobj.Witness.AddWitness(&u1)
	txobj.Witness.AddWitness(&u2)

	bbclib.SignToTransaction(txobj, &u1, keypair1)
	bbclib.SignToTransaction(txobj, &u2, keypair2)
	return txobj
}


func main() {
	//idConf := bbclib.BBcIdConfig{16, 12, 10, 8, 24}
	//bbclib.ConfigureIdLength(&idConf)
	start := time.Now().Unix()
	for i := 0; i<10000; i++ {
		txobj := makeTransaction(i)
		if result := txobj.Signatures[0].Verify(txobj.Digest()); !result {
			fmt.Println("Verification failed..")
		}
		if result := txobj.Signatures[1].Verify(txobj.Digest()); !result {
			fmt.Println("Verification failed..")
		}
	}
	end := time.Now().Unix()
	fmt.Printf("Time: %vs\n", end-start)
}
