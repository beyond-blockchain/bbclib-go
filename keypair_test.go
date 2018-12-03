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
	"crypto/sha256"
	"testing"
)

func TestGenerateKeypair(t *testing.T) {
	for curvetype := 1; curvetype < 3; curvetype++ {
		t.Run("curvetype", func(t *testing.T) {
			keypair := GenerateKeypair(curvetype, defaultCompressionMode)
			t.Logf("keypair: %v", keypair)
			if len(keypair.Pubkey) != 65 {
				t.Fatal("fail to generate keypair")
			}
		})
	}
}

func TestKeyPair_Sign_and_Verify(t *testing.T) {
	digest := sha256.Sum256([]byte("aaaaaaaaaaa"))
	digest2 := sha256.Sum256([]byte("bbbbbbbbbbbbb"))
	t.Logf("SHA-256 digest : %x\n", digest)
	t.Logf("SHA-256 digest2: %x\n", digest2)

	for curvetype := 1; curvetype < 3; curvetype++ {
		keypair := GenerateKeypair(curvetype, defaultCompressionMode)
		keypair2 := GenerateKeypair(curvetype, defaultCompressionMode)
		t.Run("curvetype", func(t *testing.T) {
			t.Logf("Curvetype = %d", curvetype)
			if len(keypair.Pubkey) != 65 {
				t.Fatal("fail to generate keypair")
			}

			t.Logf("privkey   : %x\n", keypair.Privkey)
			sig1 := keypair.Sign(digest[:])
			t.Logf("signature : %x\n", sig1)
			if len(sig1) != 64 {
				t.Fatal("fail to sign")
			}
			result := keypair.Verify(digest[:], sig1)
			if !result {
				t.Fatal("fail to verify")
			}
			t.Log("[sig1] Verify succeeded")
			result = keypair.Verify(digest2[:], sig1)
			if result {
				t.Fatal("[invalid digest] Verify returns true but not correct...")
			}
			t.Log("[invalid digest] Verify failed as expected")

			t.Logf("privkey2  : %x\n", keypair2.Privkey)
			sig2 := keypair2.Sign(digest[:])
			t.Logf("signature2: %x\n", sig2)
			if len(sig2) != 64 {
				t.Fatal("fail to sign")
			}
			result = keypair2.Verify(digest[:], sig2)
			if !result {
				t.Fatal("fail to verify")
			}
			t.Log("[sig2] Verify succeeded")
			result = keypair2.Verify(digest2[:], sig2)
			if result {
				t.Fatal("[invalid digest] Verify returns true but not correct...")
			}
			t.Log("[invalid digest] Verify failed as expected")

			result = keypair2.Verify(digest[:], sig1)
			if result {
				t.Fatal("[swap] Verify returns true but not correct...")
			}
			t.Log("[swap] Verify failed as expected")
			result = keypair.Verify(digest[:], sig2)
			if result {
				t.Fatal("[swap] Verify returns true but not correct...")
			}
			t.Log("[swap] Verify failed as expected")
		})
	}
}

func TestVerifyBBcSignature(t *testing.T) {
	digest := sha256.Sum256([]byte("aaaaaaaaaaa"))

	for curvetype := 1; curvetype < 3; curvetype++ {
		t.Run("curvetype", func(t *testing.T) {
			keypair := GenerateKeypair(curvetype, defaultCompressionMode)
			sig := BBcSignature{}
			sig.SetPublicKey(uint32(curvetype), &keypair.Pubkey)
			signature := keypair.Sign(digest[:])
			sig.SetSignature(&signature)
			result1 := keypair.Verify(digest[:], sig.Signature)
			if !result1 {
				t.Fatal("fail to verify")
			}
			t.Log("Verify succeeded")

			result := VerifyBBcSignature(digest[:], &sig)
			if !result {
				t.Fatal("fail to verify")
			}
			t.Log("Verify succeeded")

			keypair2 := GenerateKeypair(curvetype, defaultCompressionMode)
			sig2 := BBcSignature{}
			sig2.SetPublicKey(uint32(curvetype), &keypair.Pubkey)
			signature2 := keypair2.Sign(digest[:])
			sig2.SetSignature(&signature2)
			result = sig2.Verify(digest[:])
			if result {
				t.Fatal("Verify returns true but not correct...")
			}
			t.Log("Verify failed as expected")
		})
	}
}

func TestKeyPair_ConvertFromPem(t *testing.T) {
	pem := "-----BEGIN EC PRIVATE KEY-----\nMHQCAQEEIIMVMPKLJqivgRDpRDaWJCOnob6s/+t4MdoFN/8PVkNSoAcGBSuBBAAK\noUQDQgAE/k1ZM/Ker1+N0+Lg5za0sJZeSAAeYwDEWnkgnkCynErs74G/tAnu/lcu\nk8kzAivYm8mitIpJJw1OdjCDJI457g==\n-----END EC PRIVATE KEY-----"
	keypair := KeyPair{CurveType: KeyTypeEcdsaSECP256k1}
	keypair.ConvertFromPem(pem, defaultCompressionMode)
	t.Logf("keypair: %v", keypair)

	if len(keypair.Privkey) != 32 {
		t.Fatal("failed to read private key in pem format")
	}
	t.Logf("private key: %x", keypair.Privkey)
}
