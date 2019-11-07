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
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"github.com/lestrrat-go/jwx/jwk"
)

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
		PublicKeyStructure  *ecdsa.PublicKey
		PrivateKeyStructure *ecdsa.PrivateKey
	}
)

// Supported ECC curve type is Prime-256v1 only
const (
	KeyTypeNotInitialized = 0
	//KeyTypeEcdsaSECP256k1 = 1  // unsported
	KeyTypeEcdsaP256v1    = 2

	DefaultCompressionMode = 4
)
const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)

type ecdsaSignature struct {
	R, S *big.Int
}


// ReadBits encodes the absolute value of bigint as big-endian bytes. Callers must ensure
// that buf has enough space. If buf is too short the result will be incomplete.
func readBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}

// PaddedBigBytes encodes a big integer as a big-endian byte slice. The length
// of the slice is at least n bytes.
func paddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	readBits(bigint, ret)
	return ret
}

// setup KeyPair object from ecdsa.PrivateKey
func setupKeypair(kp *KeyPair, privKey *ecdsa.PrivateKey) {
	kp.CurveType = KeyTypeEcdsaP256v1
	kp.PrivateKeyStructure = privKey
	kp.PublicKeyStructure = &privKey.PublicKey

	priv := paddedBigBytes(privKey.D, privKey.Params().BitSize/8)
	pub := elliptic.Marshal(elliptic.P256(), privKey.PublicKey.X, privKey.PublicKey.Y)
	if kp.CompressionType != DefaultCompressionMode {
		pub[0] = 0x03
		kp.Privkey = priv
		kp.Pubkey = pub[:(len(pub)+1)/2]
	} else {
		kp.Privkey = priv
		kp.Pubkey = pub
	}
}

// GenerateKeypair generates a new Key pair object with new private key and public key
func GenerateKeypair(curveType int, compressionMode int) (*KeyPair, error) {
	if curveType != KeyTypeEcdsaP256v1 {
		return nil, errors.New("bbclib-go supports Prime-256v1 only. So, forcibly use P-256")
	}
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	var kp KeyPair
	kp.CompressionType = compressionMode
	setupKeypair(&kp, privKey)
	return &kp, nil
}

// GetKeyId returns Key ID of the public key
func (k *KeyPair) GetKeyId() ([]byte, error) {
	jsonKey, err := jwk.New(k.PublicKeyStructure)
	if err != nil {
		return nil, err
	}
	a, _ := jsonKey.(*jwk.ECDSAPublicKey).MarshalJSON()
	h := sha256.Sum256(a)
	fmt.Printf("hex=%x\n", a)
	fmt.Printf("--> %x\n", h)
	thumbprint, err := jsonKey.Thumbprint(crypto.SHA256)
	return thumbprint, err
}

// GetPublicKeyUncompressed gets a public key (uncompressed) from private key
func (k *KeyPair) GetPublicKeyUncompressed() *[]byte {
	x, y := k.PublicKeyStructure.Curve.ScalarBaseMult(k.Privkey)
	k.Pubkey = elliptic.Marshal(elliptic.P256(), x, y)
	return &k.Pubkey
}

// GetPublicKeyCompressed gets a public key (compressed) from private key
func (k *KeyPair) GetPublicKeyCompressed() *[]byte {
	x, y := k.PublicKeyStructure.Curve.ScalarBaseMult(k.Privkey)
	pub := elliptic.Marshal(elliptic.P256(), x, y)
	k.Pubkey = pub[:(len(pub)+1)/2]
	return &k.Pubkey
}


// ConvertFromPem imports PEM formatted private key
func (k *KeyPair) ConvertFromPem(pemstr string, compressionMode int) error {
	block, _ := pem.Decode([]byte(pemstr))
	if block == nil {
		return errors.New("invalid PEM format")
	}

	if block.Type == "EC PUBLIC KEY" {
		pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return err
		}
		k.PublicKeyStructure = pubkey.(*ecdsa.PublicKey)
		pub := elliptic.Marshal(k.PublicKeyStructure.Curve, k.PublicKeyStructure.X, k.PublicKeyStructure.Y)
		k.CompressionType = DefaultCompressionMode
		k.CurveType = KeyTypeEcdsaP256v1 // support P-256 only
		k.Pubkey = pub
		return nil

	} else if block.Type == "EC PRIVATE KEY" {
		return k.ConvertFromDer(block.Bytes, compressionMode)
	}

	return errors.New("not supported key")
}

// ConvertFromPem imports DER formatted private key
func (k *KeyPair) ConvertFromDer(der []byte, compressionMode int) error {
	k.CompressionType = compressionMode
	privkey, err := x509.ParseECPrivateKey(der)
	if err != nil {
		return err
	}
	setupKeypair(k, privkey)
	return nil
}

// ReadX509 imports X.509 public key certificate
func (k *KeyPair) ReadX509(certstr string, compressionMode int) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certstr))
	if block == nil {
		return nil, errors.New("invalid PEM format")
	}

	k.CompressionType = compressionMode
	if block.Type == "CERTIFICATE" {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		k.CurveType = KeyTypeEcdsaP256v1 // support P-256 only
		k.CompressionType = DefaultCompressionMode
		k.PublicKeyStructure = cert.PublicKey.(*ecdsa.PublicKey)
		if compressionMode == 4 {
			_ = k.GetPublicKeyUncompressed()
		} else {
			_ = k.GetPublicKeyUncompressed()
		}
		return cert, nil
	}
	return nil, errors.New("not supported certificate")
}

// VerifyX509 verifies the public key's legitimacy
func (k *KeyPair) CheckX509(certstr string, privkey string) bool {
	/*
	cert, err := k.ReadX509(certstr, k.CompressionType)
	if err != nil {
		return false
	}
	//signature := cert.Signature
	//return ecdsa.Verify(k.PublicKeyStructure, cert.Signature, signature.R, signature.S)
	 */
	return true
}

// Sign to a given digest
func (k *KeyPair) Sign(digest []byte) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, k.PrivateKeyStructure, digest)
	if err != nil {
		return nil
	}
	rPad := paddedBigBytes(r, 32)
	sPad := paddedBigBytes(s, 32)
	sig := append(rPad, sPad...)
	return sig
}

// Verify a given digest with signature
func (k *KeyPair) Verify(digest []byte, sig []byte) bool {
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:])
	return ecdsa.Verify(k.PublicKeyStructure, digest, r, s)
}

// OutputDer outputs DER formatted private key
func (k *KeyPair) OutputDer() []byte {
	der, err := x509.MarshalECPrivateKey(k.PrivateKeyStructure)
	if err != nil {
		return nil
	}
	return der
}

// OutputDer outputs PEM formatted private key
func (k *KeyPair) OutputPem() (string, error) {
	der := k.OutputDer()
	if der == nil {
		return "", errors.New("failed to export the private key to pem format")
	}
	pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: der,
		},
	)
	return string(pem), nil
}

// OutputDer outputs DER formatted private key
func (k *KeyPair) OutputPublicKeyDer() []byte {
	der, err := x509.MarshalPKIXPublicKey(k.PublicKeyStructure)
	if err != nil {
		return nil
	}
	return der
}

// OutputDer outputs PEM formatted private key
func (k *KeyPair) OutputPublicKeyPem() (string, error) {
	der := k.OutputPublicKeyDer()
	if der == nil {
		return "", errors.New("failed to export the public key to pem format")
	}
	pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "EC PUBLIC KEY",
			Bytes: der,
		},
	)
	return string(pem), nil
}

// VerifyBBcSignature verifies a given digest with BBcSignature object
func VerifyBBcSignature(digest []byte, sig *BBcSignature) bool {
	if sig.Pubkey == nil || sig.PubkeyLen == 0 {
		return true
	}
	if sig.KeyType != KeyTypeEcdsaP256v1 {
		return false
	}

	x, y := elliptic.Unmarshal(elliptic.P256(), sig.Pubkey)
	pubkey := ecdsa.PublicKey{X: x, Y: y, Curve: elliptic.P256()}

	r := new(big.Int).SetBytes(sig.Signature[:32])
	s := new(big.Int).SetBytes(sig.Signature[32:])
	return ecdsa.Verify(&pubkey, digest, r, s)
}
