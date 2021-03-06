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
	"bytes"
	"encoding/hex"
	"testing"
)

func TestGenerateKeypair(t *testing.T) {
	t.Run("simply generate keypair", func(t *testing.T) {
		keypair, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
		t.Logf("keypair: %v", keypair)
		if len(keypair.Pubkey) != 65 {
			t.Fatal("fail to generate keypair")
		}
	})
}

func TestKeyPair_Sign_and_Verify(t *testing.T) {
	digest := sha256.Sum256([]byte("aaaaaaaaaaa"))
	digest2 := sha256.Sum256([]byte("bbbbbbbbbbbbb"))
	t.Logf("SHA-256 digest : %x\n", digest)
	t.Logf("SHA-256 digest2: %x\n", digest2)

	keypair, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
	keypair2, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
	t.Run("curvetype = KeyTypeEcdsaP256v1 only", func(t *testing.T) {
		t.Logf("Curvetype = %d", KeyTypeEcdsaP256v1)
		if len(keypair.Pubkey) != 65 {
			t.Fatal("fail to generate keypair")
		}

		t.Logf("privkey   : %x\n", keypair.Privkey)
		sig1 := keypair.Sign(digest[:])
		t.Logf("signature : %x (%d)\n", sig1, len(sig1))
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
		t.Logf("signature2: %x (%d)\n", sig2, len(sig2))
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

func TestVerifyBBcSignature(t *testing.T) {
	digest := sha256.Sum256([]byte("aaaaaaaaaaa"))

	t.Run("curvetype", func(t *testing.T) {
		keypair, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
		sig := BBcSignature{}
		sig.SetPublicKey(uint32(KeyTypeEcdsaP256v1), &keypair.Pubkey)
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

		keypair2, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
		sig2 := BBcSignature{}
		sig2.SetPublicKey(uint32(KeyTypeEcdsaP256v1), &keypair.Pubkey)
		signature2 := keypair2.Sign(digest[:])
		sig2.SetSignature(&signature2)
		result = sig2.Verify(digest[:])
		if result {
			t.Fatal("Verify returns true but not correct...")
		}
		t.Log("Verify failed as expected")
	})
}

func TestKeyPair_ConvertFromPem(t *testing.T) {
	pem := "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIJ6T6fh4uQCmxImVaRueeTuZ4rSvcJtO1mRQgolFs4fgoAoGCCqGSM49\nAwEHoUQDQgAEDbE3poddh6YhzVAkOO6edq9VimbnD5t46eu/9CW6Y2C3uQxaBV39\nJFLML26BlpFqRD/W+GgqpTgTfbWq0cJj+A==\n-----END EC PRIVATE KEY-----"
	keypair := KeyPair{CurveType: KeyTypeEcdsaP256v1}
	err := keypair.ConvertFromPem(pem, DefaultCompressionMode)
	t.Logf("error: %v", err)
	t.Logf("keypair: %v", keypair)

	if len(keypair.Privkey) != 32 {
		t.Fatal("failed to read private key in pem format")
	}
	t.Logf("private key: %x", keypair.Privkey)
}

func TestKeyPair_Export_Import(t *testing.T) {
	keypair, _ := GenerateKeypair(KeyTypeEcdsaP256v1, DefaultCompressionMode)
	der := keypair.OutputDer()
	pem, _ := keypair.OutputPem()

	keypair2 := KeyPair{CurveType: KeyTypeNotInitialized}
	_ = keypair2.ConvertFromDer(der, DefaultCompressionMode)
	if bytes.Compare(keypair.Privkey, keypair2.Privkey) != 0 {
		t.Fatal("export or import is failed (DER)")
	}
	if keypair.CurveType != keypair2.CurveType {
		t.Fatal("curve type cannot be obtained (DER)")
	}

	pubkey2 := keypair2.GetPublicKeyCompressed()
	t.Logf("public key (compressed): %v", pubkey2)

	keypair3 := KeyPair{CurveType: KeyTypeNotInitialized}
	_ = keypair3.ConvertFromPem(pem, DefaultCompressionMode)
	if bytes.Compare(keypair.Privkey, keypair3.Privkey) != 0 {
		t.Fatal("export or import is failed (PEM)")
	}
	if keypair.CurveType != keypair3.CurveType {
		t.Fatal("curve type cannot be obtained (PEM)")
	}

	pubkey3 := keypair3.GetPublicKeyCompressed()
	t.Logf("public key (compressed): %v", pubkey3)
}


func TestKeyId(t *testing.T) {
	t.Run("KeyID calculation", func(t *testing.T) {
		pem := "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIDSt1IOhS5ZmY6nkX/Wh7pT+Y45TmYxrwoc1pG72v387oAoGCCqGSM49\nAwEHoUQDQgAEdEsjD2i2LytHOjNxxc9PbFeqQ89aMLOfmdBbEoSOhZBukJ52EqQM\nhOdgHqyqD4hEyYxgDu3uIbKat+lEZEhb3Q==\n-----END EC PRIVATE KEY-----"
		keyIdHex := "f7211e7e0db043a29fc6624006bf8dfaac9fafe65b0c6e0dfd573d81e95bd83e"
		keypair := KeyPair{CurveType: KeyTypeEcdsaP256v1}
		err := keypair.ConvertFromPem(pem, DefaultCompressionMode)
		if err != nil {
			t.Fatal("Fail to import pem key")
		}

		keyId, _ := keypair.GetKeyId()
		t.Logf("keyID=%x\n", keyId)
		if hex.EncodeToString(keyId) != keyIdHex {
			t.Fatal("invalid keyId")
		}
	})
}
