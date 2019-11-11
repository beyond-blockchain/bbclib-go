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

	bodyString := fmt.Sprintf("teststring!!!!:%d", num)
	filedat, _ := ioutil.ReadFile("./asset_test.go")
	datobj := map[string]string{"param1": "aaa", "param2": "bbb", "param3": string(num)}

	txobj := bbclib.MakeTransaction(3, 0, true)
	txobj.Events[0].SetAssetGroup(&assetGroupID).AddMandatoryApprover(&u1).CreateAsset(&u1, nil, bodyString)
	txobj.Events[1].SetAssetGroup(&assetGroupID).AddMandatoryApprover(&u2).CreateAsset(&u2, &filedat, nil)
	txobj.Events[2].SetAssetGroup(&assetGroupID).AddMandatoryApprover(&u1).CreateAsset(&u1, nil, &datobj)
	txobj.AddWitness(&u1).AddWitness(&u2)

	txobj.Sign(&u1, keypair1, false).Sign(&u2, keypair2, false)
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
