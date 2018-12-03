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

import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
)

/*
BBcSignature definition

The BBcSignature holds public key and signature. The signature is for the TransactionID of the transaction object.
*/
type (
	BBcSignature struct {
		KeyType      uint32
		Pubkey       []byte
		PubkeyLen    uint32
		Signature    []byte
		SignatureLen uint32
	}
)

// Stringer outputs the content of the object
func (p *BBcSignature) Stringer() string {
	if p.KeyType == KeyTypeNotInitialized {
		return "  Not initialized\n"
	}
	ret := fmt.Sprintf("  key_type: %d\n", p.KeyType)
	ret += fmt.Sprintf("  signature: %x\n", p.Signature)
	ret += fmt.Sprintf("  pubkey: %x\n", p.Pubkey)
	return ret
}

// SetPublicKey sets signature binary in the object
func (p *BBcSignature) SetPublicKey(keyType uint32, pubkey *[]byte) {
	p.KeyType = keyType
	p.Pubkey = *pubkey
	p.PubkeyLen = uint32(len(p.Pubkey) * 8)
}

// SetPublicKeyByKeypair sets public key (in keypair object) in the object
func (p *BBcSignature) SetPublicKeyByKeypair(keypair *KeyPair) {
	p.KeyType = uint32(keypair.CurveType)
	p.Pubkey = keypair.Pubkey
	p.PubkeyLen = uint32(len(p.Pubkey) * 8)
}

// SetSignature sets signature binary in the object
func (p *BBcSignature) SetSignature(sig *[]byte) {
	p.Signature = *sig
	p.SignatureLen = uint32(len(p.Signature) * 8)
}

// Verify the TransactionID of the parent BBcTransaction object with the signature in the object
func (p *BBcSignature) Verify(digest []byte) bool {
	return VerifyBBcSignature(digest, p)
}

// Pack returns the binary data of the BBcSignature object
func (p *BBcSignature) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	Put4byte(buf, p.KeyType)
	if p.KeyType == KeyTypeNotInitialized {
		return buf.Bytes(), nil
	}

	Put4byte(buf, p.PubkeyLen)
	if err := binary.Write(buf, binary.LittleEndian, p.Pubkey); err != nil {
		return nil, err
	}

	Put4byte(buf, p.SignatureLen)
	if err := binary.Write(buf, binary.LittleEndian, p.Signature); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Unpack the BBcSignature object to the binary data
func (p *BBcSignature) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	keyType, err := Get4byte(buf)
	if err != nil {
		return err
	}
	if keyType == 0 {
		return nil
	}
	p.KeyType = uint32(keyType)

	p.PubkeyLen, err = Get4byte(buf)
	if err != nil {
		return err
	}
	p.Pubkey = make([]byte, int(p.PubkeyLen/8))
	p.Pubkey, _ = GetBytes(buf, int(p.PubkeyLen/8))

	p.SignatureLen, err = Get4byte(buf)
	if err != nil {
		return err
	}
	p.Signature = make([]byte, int(p.SignatureLen/8))
	p.Signature, _ = GetBytes(buf, int(p.SignatureLen/8))

	return nil
}

// RecoverSignatureObject is a utility for recovering signature data into BBcSignature object
func RecoverSignatureObject(dat *[]byte) *BBcSignature {
	sig := BBcSignature{}
	sig.Unpack(dat)
	return &sig
}
