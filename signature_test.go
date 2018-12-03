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
)

func TestSignaturePackUnpack(t *testing.T) {

	t.Run("simple creation (set by keypair)", func(t *testing.T) {
		keypair := GenerateKeypair(KeyTypeEcdsaP256v1, defaultCompressionMode)
		sig := BBcSignature{}
		sig.SetPublicKeyByKeypair(&keypair)
		signature := GetRandomValue(64)
		sig.SetSignature(&signature)

		t.Log("---------------signature-----------------")
		t.Logf("%v", sig.Stringer())
		t.Log("--------------------------------------")

		dat, err := sig.Pack()
		if err != nil {
			t.Fatalf("failed to serialize signature object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcSignature{}
		obj2.Unpack(&dat)
		t.Log("---------------transaction-----------------")
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(sig.Signature, obj2.Signature) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("simple creation (set with raw public key)", func(t *testing.T) {
		sig := BBcSignature{}
		keydat := GetRandomValue(65)
		sig.SetPublicKey(KeyTypeEcdsaP256v1, &keydat)
		signature := GetRandomValue(64)
		sig.SetSignature(&signature)

		t.Log("---------------signature-----------------")
		t.Logf("%v", sig.Stringer())
		t.Log("--------------------------------------")

		dat, err := sig.Pack()
		if err != nil {
			t.Fatalf("failed to serialize signature object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := RecoverSignatureObject(&dat)
		t.Log("---------------transaction-----------------")
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(sig.Signature, obj2.Signature) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}
