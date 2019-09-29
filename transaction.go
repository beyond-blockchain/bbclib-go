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
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"time"
)

/*
BBcTransaction definition

BBcTransaction is just a container of various objects.

Events, References, Relations and Signatures are list of BBcEvent, BBcReference, BBcRelation and BBcSignature objects, respectively.
"digestCalculating", "TransactionBaseDigest", "TransactionData" and "SigIndexedUsers" are not included in the packed data. They are internal use only.

Calculating TransactionID

How to calculate the TransactionID of the transaction is a little bit complicated, meaning that 2-step manner.
This is because inter-domain transaction authenticity (i.e., CrossReference) can be conducted in secure manner.
By presenting TransactionBaseDigest (see below) to an outer-domain, the domain user can confirm the existence of the transaction in the past.
(no need to present whole transaction data including the asset information).

1st step:
  * Pack info (from version to Witness) by packBase()
  * Calculate SHA256 digest of the packed info. This value is TransactionBaseDigest.

2nd step:
  * Pack BBcCrossRef object to get packed data by packCrossRef()
  * Concatenate TransactionBaseDigest and the packed BBcCrossRef
  * Calculate SHA256 digest of the concatenated data. This value is TransactionID
*/
type (
	BBcTransaction struct {
		digestCalculating     bool
		IdLengthConf          BBcIdConfig
		TransactionID         []byte
		TransactionBaseDigest []byte
		TransactionData       []byte
		SigIndexedUsers       [][]byte
		Version               uint32
		Timestamp             int64
		TransactionIdLength   int
		Events                []*BBcEvent
		References            []*BBcReference
		Relations             []*BBcRelation
		Witness               *BBcWitness
		Crossref              *BBcCrossRef
		Signatures            []*BBcSignature
	}
)

// Stringer outputs the content of the object
func (p *BBcTransaction) Stringer() string {
	var ret string
	ret = "------- Dump of the transaction data ------\n"
	ret += fmt.Sprintf("* transaction_id: %x\n", p.TransactionID)
	ret += fmt.Sprintf("version: %d\n", p.Version)
	ret += fmt.Sprintf("timestamp: %d\n", p.Timestamp)
	if p.Version != 0 {
		ret += fmt.Sprintf("transaction_id_length: %d\n", p.IdLengthConf.TransactionIdLength)
	}

	ret += fmt.Sprintf("Event[]: %d\n", len(p.Events))
	for i := range p.Events {
		ret += fmt.Sprintf(" [%d]\n", i)
		ret += p.Events[i].Stringer()
	}

	ret += fmt.Sprintf("Reference[]: %d\n", len(p.References))
	for i := range p.References {
		ret += fmt.Sprintf(" [%d]\n", i)
		ret += p.References[i].Stringer()
	}

	ret += fmt.Sprintf("Relation[]: %d\n", len(p.Relations))
	for i := range p.Relations {
		ret += fmt.Sprintf(" [%d]\n", i)
		ret += p.Relations[i].Stringer()
	}

	if p.Witness != nil {
		ret += p.Witness.Stringer()
	} else {
		ret += "Witness: None\n"
	}

	if p.Crossref != nil {
		ret += p.Crossref.Stringer()
	} else {
		ret += "Cross_Ref: None\n"
	}

	ret += fmt.Sprintf("Signature[]: %d\n", len(p.Signatures))
	for i := range p.Signatures {
		ret += fmt.Sprintf(" [%d]\n", i)
		ret += p.Signatures[i].Stringer()
	}
	return ret
}

// Set ID length configuration
func (p *BBcTransaction) SetIdLengthConf(conf * BBcIdConfig) {
	p.IdLengthConf.TransactionIdLength = conf.TransactionIdLength
	p.IdLengthConf.UserIdLength = conf.UserIdLength
	p.IdLengthConf.AssetGroupIdLength = conf.AssetGroupIdLength
	p.IdLengthConf.AssetIdLength = conf.AssetIdLength
	p.IdLengthConf.NonceLength = conf.NonceLength
	p.TransactionIdLength = conf.TransactionIdLength
}

// AddEvent adds the BBcEvent object in the transaction object
func (p *BBcTransaction) AddEvent(obj *BBcEvent) {
	obj.SetIdLengthConf(&p.IdLengthConf)
	p.Events = append(p.Events, obj)
}

// AddReference adds the BBcReference object in the transaction object
func (p *BBcTransaction) AddReference(obj *BBcReference) {
	obj.SetIdLengthConf(&p.IdLengthConf)
	p.References = append(p.References, obj)
	obj.Transaction = p
}

// AddRelation adds the BBcRelation object in the transaction object
func (p *BBcTransaction) AddRelation(obj *BBcRelation) {
	obj.SetIdLengthConf(&p.IdLengthConf)
	obj.SetVersion(p.Version)
	p.Relations = append(p.Relations, obj)
}

// AddWitness sets the BBcWitness object in the transaction object
func (p *BBcTransaction) AddWitness(obj *BBcWitness) {
	obj.SetIdLengthConf(&p.IdLengthConf)
	p.Witness = obj
	obj.Transaction = p
}

// AddCrossRef sets the BBcCrossRef object in the transaction object
func (p *BBcTransaction) AddCrossRef(obj *BBcCrossRef) {
	obj.SetIdLengthConf(&p.IdLengthConf)
	p.Crossref = obj
}

// AddSignature adds the BBcSignature object for the specified userID in the transaction object
func (p *BBcTransaction) AddSignature(userID *[]byte, sig *BBcSignature) {
	uid := make([]byte, int(p.IdLengthConf.UserIdLength))
	copy(uid, *userID)
	for i := range p.SigIndexedUsers {
		if reflect.DeepEqual(p.SigIndexedUsers[i], uid) {
			p.Signatures[i] = sig
			return
		}
	}
	p.SigIndexedUsers = append(p.SigIndexedUsers, uid)
	p.Signatures = append(p.Signatures, sig)
}

// GetSigIndex reserves and returns the position (index) of the corespondent userID in the signature list
func (p *BBcTransaction) GetSigIndex(userID []byte) int {
	uid := make([]byte, int(p.IdLengthConf.UserIdLength))
	copy(uid, userID)
	var i = -1
	for i = range p.SigIndexedUsers {
		if reflect.DeepEqual(p.SigIndexedUsers[i], uid) {
			return i
		}
	}
	p.SigIndexedUsers = append(p.SigIndexedUsers, uid)
	sig := BBcSignature{}
	p.Signatures = append(p.Signatures, &sig)
	return i + 1
}

// SetSigIndex simply sets the index of signature list for the specified userID
func (p *BBcTransaction) SetSigIndex(userID []byte, idx int) {
	var i = -1
	for i = range p.SigIndexedUsers {
		if reflect.DeepEqual(p.SigIndexedUsers[i], userID) {
			return
		}
	}
	if len(p.SigIndexedUsers)-1 < idx {
		for i:=0; i<idx+1; i++ {
			val := make([]byte, p.IdLengthConf.UserIdLength)
			p.SigIndexedUsers = append(p.SigIndexedUsers, val)
		}
	}
	p.SigIndexedUsers[idx] = userID
}

// Sign TransactionID using private key in the given keypair
func (p *BBcTransaction) Sign(keypair *KeyPair) ([]byte, error) {
	digest := p.Digest()
	signature := keypair.Sign(digest)
	if signature == nil {
		return nil, errors.New("fail to sign")
	}
	return signature, nil
}

// SignAndAdd signs to TransactionID and add BBcSignature in the transaction
func (p *BBcTransaction) SignAndAdd(keypair *KeyPair, userID []byte, noPubkey bool) error {
	digest := p.Digest()
	signature := keypair.Sign(digest)
	if signature == nil {
		return errors.New("fail to sign")
	}
	idx := p.GetSigIndex(userID)
	p.Signatures[idx].KeyType = uint32(keypair.CurveType)
	p.Signatures[idx].SetSignature(&signature)
	if ! noPubkey {
		p.Signatures[idx].Pubkey = keypair.Pubkey
		p.Signatures[idx].PubkeyLen = uint32(len(keypair.Pubkey) * 8)
	}
	return nil
}

// VerifyAll verifies TransactionID with all BBcSignature objects in the transaction
func (p *BBcTransaction) VerifyAll() (bool, int) {
	digest := p.Digest()
	for i := range p.Signatures {
		if p.Signatures[i].KeyType == KeyTypeNotInitialized {
			continue
		}
		if ret := VerifyBBcSignature(digest, p.Signatures[i]); !ret {
			return false, i
		}
	}
	return true, -1
}

// Digest calculates TransactionID of the BBcTransaction object
func (p *BBcTransaction) Digest() []byte {
	p.digestCalculating = true
	if p.TransactionID == nil {
		p.TransactionID = make([]byte, p.TransactionIdLength)
	}
	buf := new(bytes.Buffer)

	err := p.packBase(buf)
	if err != nil {
		p.digestCalculating = false
		return nil
	}

	buf = new(bytes.Buffer)
	if err = binary.Write(buf, binary.LittleEndian, p.TransactionBaseDigest); err != nil {
		p.digestCalculating = false
		return nil
	}

	err = p.packCrossRef(buf)
	if err != nil {
		p.digestCalculating = false
		return nil
	}

	digest := sha256.Sum256(buf.Bytes())
	p.TransactionID = digest[:p.TransactionIdLength]
	p.digestCalculating = false
	return digest[:]
}

// packCrossRef packs only BBcCrossRef object in binary data
func (p *BBcTransaction) packCrossRef(buf *bytes.Buffer) error {
	if p.Crossref != nil {
		dat, err := p.Crossref.Pack()
		if err != nil {
			return err
		}
		Put2byte(buf, 1)
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return err
		}
	} else {
		Put2byte(buf, 0)
	}
	return nil
}

// packBase packs the base part of BBcTransaction object in binary data (from version to witness)
func (p *BBcTransaction) packBase(buf *bytes.Buffer) error {
	Put4byte(buf, p.Version)
	if p.Timestamp == 0 {
		p.Timestamp = time.Now().UnixNano() / int64(time.Microsecond)
	}
	Put8byte(buf, p.Timestamp)
	Put2byte(buf, uint16(p.TransactionIdLength))

	Put2byte(buf, uint16(len(p.Events)))
	for _, obj := range p.Events {
		dat, err := obj.Pack()
		if err != nil {
			return err
		}
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return err
		}
	}

	Put2byte(buf, uint16(len(p.References)))
	for _, obj := range p.References {
		dat, err := obj.Pack()
		if err != nil {
			return err
		}
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return err
		}
	}

	Put2byte(buf, uint16(len(p.Relations)))
	for _, obj := range p.Relations {
		dat, err := obj.Pack()
		if err != nil {
			return err
		}
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return err
		}
	}

	if p.Witness != nil {
		dat, err := p.Witness.Pack()
		if err != nil {
			return err
		}
		Put2byte(buf, 1)
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return err
		}
	} else {
		Put2byte(buf, 0)
	}

	digest := sha256.Sum256(buf.Bytes())
	p.TransactionBaseDigest = digest[:]

	return nil
}

// Pack BBcTransaction object in binary data
func (p *BBcTransaction) Pack() ([]byte, error) {
	if !p.digestCalculating && p.TransactionID == nil {
		p.Digest()
	}

	if p.Version == 0 {
		return nil, errors.New("not support version=0 transaction")
	}

	buf := new(bytes.Buffer)
	err := p.packBase(buf)
	if err != nil {
		return nil, err
	}
	err = p.packCrossRef(buf)
	if err != nil {
		return nil, err
	}

	Put2byte(buf, uint16(len(p.Signatures)))
	for _, obj := range p.Signatures {
		dat, err := obj.Pack()
		if err != nil {
			return nil, err
		}
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return nil, err
		}
	}

	p.TransactionData = buf.Bytes()
	return p.TransactionData, nil
}

// unpackHeader unpacks the header part of the binary data
func (p *BBcTransaction) unpackHeader(buf *bytes.Buffer) error {
	var err error
	p.Version, err = Get4byte(buf)
	if err != nil {
		return err
	}

	p.Timestamp, err = Get8byte(buf)
	if err != nil {
		return err
	}

	idLen, err := Get2byte(buf)
	if err != nil {
		return err
	}
	p.IdLengthConf.TransactionIdLength = int(idLen)
	p.TransactionIdLength = int(idLen)
	return nil
}

// unpackEvent unpacks the events part of the binary data
func (p *BBcTransaction) unpackEvent(buf *bytes.Buffer) error {
	num, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(num); i++ {
		size, err2 := Get4byte(buf)
		if err2 != nil {
			return err2
		}
		data, _, err2 := GetBytes(buf, int(size))
		if err2 != nil {
			return err2
		}
		obj := BBcEvent{}
		obj.SetIdLengthConf(&p.IdLengthConf)
		obj.Unpack(&data)
		p.Events = append(p.Events, &obj)
	}
	return nil
}

// unpackReference unpacks the references part of the binary data
func (p *BBcTransaction) unpackReference(buf *bytes.Buffer) error {
	num, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(num); i++ {
		size, err2 := Get4byte(buf)
		if err2 != nil {
			return err2
		}
		data, _, err2 := GetBytes(buf, int(size))
		if err2 != nil {
			return err2
		}
		obj := BBcReference{}
		obj.SetIdLengthConf(&p.IdLengthConf)
		obj.SetTransaction(p)
		obj.Unpack(&data)
		p.References = append(p.References, &obj)
	}
	return nil
}

// unpackRelation unpacks the relations part of the binary data
func (p *BBcTransaction) unpackRelation(buf *bytes.Buffer) error {
	num, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(num); i++ {
		size, err2 := Get4byte(buf)
		if err2 != nil {
			return err2
		}
		data, _, _ := GetBytes(buf, int(size))
		obj := BBcRelation{Version: p.Version}
		obj.SetIdLengthConf(&p.IdLengthConf)
		obj.Unpack(&data)
		p.Relations = append(p.Relations, &obj)
	}
	return nil
}

// unpackWitness unpacks the witness part of the binary data
func (p *BBcTransaction) unpackWitness(buf *bytes.Buffer) error {
	num, err := Get2byte(buf)
	if err != nil {
		return err
	}
	if num > 0 {
		size, err2 := Get4byte(buf)
		if err2 != nil {
			return err2
		}
		data, _, _ := GetBytes(buf, int(size))
		p.Witness = &BBcWitness{}
		p.Witness.SetIdLengthConf(&p.IdLengthConf)
		p.Witness.SetTransaction(p)
		p.Witness.Unpack(&data)
	}
	return nil
}

// unpackCrossRef unpacks the crossref part of the binary data
func (p *BBcTransaction) unpackCrossRef(buf *bytes.Buffer) error {
	num, err := Get2byte(buf)
	if err != nil {
		return err
	}
	if num > 0 {
		size, err2 := Get4byte(buf)
		if err2 != nil {
			return err2
		}
		dat, _, err2 := GetBytes(buf, int(size))
		if err2 != nil {
			return err2
		}
		p.Crossref = &BBcCrossRef{}
		p.Crossref.SetIdLengthConf(&p.IdLengthConf)
		p.Crossref.Unpack(&dat)
	}
	return nil
}

// unpackSignature unpacks the signatures part of the binary data
func (p *BBcTransaction) unpackSignature(buf *bytes.Buffer) error {
	num, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(num); i++ {
		size, err2 := Get4byte(buf)
		if err2 != nil {
			return err2
		}
		data, _, err2 := GetBytes(buf, int(size))
		if err2 != nil {
			return err2
		}
		obj := BBcSignature{}
		obj.Unpack(&data)
		p.Signatures = append(p.Signatures, &obj)
	}
	return nil
}

// Unpack binary data to BBcTransaction object
func (p *BBcTransaction) Unpack(dat *[]byte) error {
	buf := bytes.NewBuffer(*dat)

	if err := p.unpackHeader(buf); err != nil {
		return err
	}

	if err := p.unpackEvent(buf); err != nil {
		return err
	}

	if err := p.unpackReference(buf); err != nil {
		return err
	}

	if err := p.unpackRelation(buf); err != nil {
		return err
	}

	if err := p.unpackWitness(buf); err != nil {
		return err
	}

	if err := p.unpackCrossRef(buf); err != nil {
		return err
	}

	if err := p.unpackSignature(buf); err != nil {
		return err
	}

	p.Digest()
	return nil
}
