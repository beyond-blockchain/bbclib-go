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

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: ${SRCDIR}/libbbcsig.a -ldl
#include "libbbcsig.h"
*/
import "C"
import "unsafe"

/*
KeyPair definition

A KeyPair object hold a pair of private key and public key.
This object includes functions for sign and verify a signature. The sign/verify functions is realized by "libbbcsig".
*/
type (
	KeyPair struct {
		CurveType		int
		CompressionType	int
		Pubkey    		[]byte
		Privkey   		[]byte
	}
)

// Supported ECC curve type is SECP256k1 and Prime-256v1.
const (
	KeyTypeNotInitialized = 0
	KeyTypeEcdsaSECP256k1 = 1
	KeyTypeEcdsaP256v1    = 2

	DefaultCompressionMode = 4
)

// GenerateKeypair generates a new Key pair object with new private key and public key
func GenerateKeypair(curveType int, compressionMode int) KeyPair {
	pubkey := make([]byte, 100)
	privkey := make([]byte, 100)
	var lenPubkey, lenPrivkey C.int
	C.generate_keypair(C.int(curveType), C.uint8_t(compressionMode), &lenPubkey, (*C.uint8_t)(unsafe.Pointer(&pubkey[0])),
		&lenPrivkey, (*C.uint8_t)(unsafe.Pointer(&privkey[0])))
	return KeyPair{CurveType: curveType, Pubkey: pubkey[:lenPubkey], Privkey: privkey[:lenPrivkey]}
}

// GetPublicKeyUncompressed gets a public key (uncompressed) from private key
func (k *KeyPair) GetPublicKeyUncompressed() *[]byte {
	var lenPubkey C.int
	pubkey := make([]byte, 100)
	C.get_public_key_uncompressed(C.int(k.CurveType), C.int(len(k.Privkey)),
		(*C.uint8_t)(unsafe.Pointer(&k.Privkey[0])),
		&lenPubkey, (*C.uint8_t)(unsafe.Pointer(&pubkey[0])))
	k.Pubkey = pubkey[:lenPubkey]
	return &k.Pubkey
}

// GetPublicKeyCompressed gets a public key (compressed) from private key
func (k *KeyPair) GetPublicKeyCompressed() *[]byte {
	var lenPubkey C.int
	pubkey := make([]byte, 100)
	C.get_public_key_compressed(C.int(k.CurveType), C.int(len(k.Privkey)),
		(*C.uint8_t)(unsafe.Pointer(&k.Privkey[0])),
		&lenPubkey, (*C.uint8_t)(unsafe.Pointer(&pubkey[0])))
	k.Pubkey = pubkey[:lenPubkey]
	return &k.Pubkey
}


// ConvertFromPem imports PEM formatted private key
func (k *KeyPair) ConvertFromPem(pem string, compressionMode int) {
	pubkey := make([]byte, 100)
	privkey := make([]byte, 100)
	pemByte := ([]byte)(pem)

	var curveType, lenPubkey, lenPrivkey C.int
	C.convert_from_pem_with_curvetype(
		(*C.char)(unsafe.Pointer(&pemByte[0])), (C.uint8_t)(compressionMode), &curveType,
		&lenPubkey, (*C.uint8_t)(unsafe.Pointer(&pubkey[0])),
		&lenPrivkey, (*C.uint8_t)(unsafe.Pointer(&privkey[0])))
	k.Pubkey = pubkey[:lenPubkey]
	k.Privkey = privkey[:lenPrivkey]
	k.CurveType = int(curveType)
}

// ConvertFromPem imports DER formatted private key
func (k *KeyPair) ConvertFromDer(der []byte, compressionMode int) {
	pubkey := make([]byte, 100)
	privkey := make([]byte, 100)

	var curveType, lenPubkey, lenPrivkey C.int
	C.convert_from_der_with_curvetype(
		C.long(len(der)), (*C.uint8_t)(unsafe.Pointer(&der[0])), (C.uint8_t)(compressionMode), &curveType,
		&lenPubkey, (*C.uint8_t)(unsafe.Pointer(&pubkey[0])),
		&lenPrivkey, (*C.uint8_t)(unsafe.Pointer(&privkey[0])))
	k.Pubkey = pubkey[:lenPubkey]
	k.Privkey = privkey[:lenPrivkey]
	k.CurveType = int(curveType)
}

// ReadX509 imports X.509 public key certificate
func (k *KeyPair) ReadX509(cert string, compressionMode int) {
	pubkey := make([]byte, 100)
	certByte := ([]byte)(cert)

	var curveType, lenPubkey C.int
	C.read_x509_with_curvetype(
		(*C.char)(unsafe.Pointer(&certByte[0])), (C.uint8_t)(compressionMode), &curveType,
		&lenPubkey, (*C.uint8_t)(unsafe.Pointer(&pubkey[0])))
	k.Pubkey = pubkey[:lenPubkey]
	k.CurveType = int(curveType)
}

// VerifyX509 verifies the public key's legitimacy
func (k *KeyPair) CheckX509(cert string, privkey string) bool {
	certByte := ([]byte)(cert)
	privkeyByte := ([]byte)(privkey)
	result := C.verify_x509( (*C.char)(unsafe.Pointer(&certByte[0])), (*C.char)(unsafe.Pointer(&privkeyByte[0])))
	return result == 1
}

// Sign to a given digest
func (k *KeyPair) Sign(digest []byte) []byte {
	sigR := make([]byte, 100)
	sigS := make([]byte, 100)
	var lenSigR, lenSigS C.uint
	C.sign(C.int(k.CurveType), C.int(len(k.Privkey)), (*C.uint8_t)(unsafe.Pointer(&k.Privkey[0])),
		C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
		(*C.uint8_t)(unsafe.Pointer(&sigR[0])), (*C.uint8_t)(unsafe.Pointer(&sigS[0])),
		(*C.uint)(&lenSigR), (*C.uint)(&lenSigS))

	if lenSigR < 32 {
		zeros := make([]byte, 32-lenSigR)
		for i := range zeros {
			zeros[i] = 0
		}
		sigR = append(zeros, sigR[:32]...)
	}
	if lenSigS < 32 {
		zeros := make([]byte, 32-lenSigS)
		for i := range zeros {
			zeros[i] = 0
		}
		sigS = append(zeros, sigS[:32]...)
	}
	//sig := make([]byte, lenSigR+lenSigS)
	sig := append(sigR[:32], sigS[:32]...)

	return sig
}

// Verify a given digest with signature
func (k *KeyPair) Verify(digest []byte, sig []byte) bool {
	result := C.verify(C.int(k.CurveType), C.int(len(k.Pubkey)), (*C.uint8_t)(unsafe.Pointer(&(k.Pubkey[0]))),
		C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
		C.int(len(sig)), (*C.uint8_t)(unsafe.Pointer(&sig[0])))
	return result == 1
}

// OutputDer outputs DER formatted private key
func (k *KeyPair) OutputDer() []byte {
	der := make([]byte, 1024)
	C.output_der(C.int(k.CurveType), C.int(len(k.Privkey)), (*C.uint8_t)(unsafe.Pointer(&k.Privkey[0])), (*C.uint8_t)(unsafe.Pointer(&der[0])))
	return der
}

// OutputDer outputs PEM formatted private key
func (k *KeyPair) OutputPem() string {
	pem := make([]byte, 2048)
	C.output_pem(C.int(k.CurveType), C.int(len(k.Privkey)), (*C.uint8_t)(unsafe.Pointer(&k.Privkey[0])), (*C.uint8_t)(unsafe.Pointer(&pem[0])))
	return string(pem)
}

// OutputDer outputs DER formatted private key
func (k *KeyPair) OutputPublicKeyDer() []byte {
	der := make([]byte, 8192)
	C.output_public_key_der(C.int(k.CurveType), C.int(len(k.Pubkey)), (*C.uint8_t)(unsafe.Pointer(&k.Pubkey[0])), (*C.uint8_t)(unsafe.Pointer(&der[0])))
	return der
}

// OutputDer outputs PEM formatted private key
func (k *KeyPair) OutputPublicKeyPem() string {
	pem := make([]byte, 32768)
	C.output_public_key_pem(C.int(k.CurveType), C.int(len(k.Pubkey)), (*C.uint8_t)(unsafe.Pointer(&k.Pubkey[0])), (*C.uint8_t)(unsafe.Pointer(&pem[0])))
	return string(pem)
}

// VerifyBBcSignature verifies a given digest with BBcSignature object
func VerifyBBcSignature(digest []byte, sig *BBcSignature) bool {
	result := C.verify(C.int(sig.KeyType), C.int(len(sig.Pubkey)), (*C.uint8_t)(unsafe.Pointer(&sig.Pubkey[0])),
		C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
		C.int(len(sig.Signature)), (*C.uint8_t)(unsafe.Pointer(&sig.Signature[0])))
	return result == 1
}
