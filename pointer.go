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
	"fmt"
)

/*
BBcPointer definition

BBcPointer(s) are included in BBcRelation object. A BBcPointer object includes "TransactionID" and "AssetID" and
declares that the transaction has a certain relationship with the BBcTransaction and BBcAsset object specified by those IDs.

IDLength is not included in a packed data. It is for internal use only.
*/
type (
	BBcPointer struct {
		IdLengthConf  *BBcIdConfig
		TransactionID []byte
		AssetID       []byte
	}
)

// Stringer outputs the content of the object
func (p *BBcPointer) Stringer() string {
	ret := fmt.Sprintf("     transaction_id: %x\n", p.TransactionID)
	ret += fmt.Sprintf("     asset_id: %x\n", p.AssetID)
	return ret
}

// Set ID length configuration
func (p *BBcPointer) SetIdLengthConf(conf * BBcIdConfig) {
	p.IdLengthConf = conf
}

// Add sets essential information to the BBcPointer object
func (p *BBcPointer) Add(txid *[]byte, asid *[]byte) {
	if txid != nil {
		p.TransactionID = make([]byte, p.IdLengthConf.TransactionIdLength)
		copy(p.TransactionID, (*txid)[:p.IdLengthConf.TransactionIdLength])
	}
	if asid != nil {
		p.AssetID = make([]byte, p.IdLengthConf.AssetIdLength)
		copy(p.AssetID, (*asid)[:p.IdLengthConf.AssetIdLength])
	}
}

// Pack returns the binary data of the BBcPointer object
func (p *BBcPointer) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	PutBigInt(buf, &p.TransactionID, p.IdLengthConf.TransactionIdLength)

	if p.AssetID != nil {
		Put2byte(buf, 1)
	} else {
		Put2byte(buf, 0)
		return buf.Bytes(), nil
	}

	PutBigInt(buf, &p.AssetID, p.IdLengthConf.AssetIdLength)

	return buf.Bytes(), nil
}

// Unpack the BBcPointer object to the binary data
func (p *BBcPointer) Unpack(dat *[]byte) error {
	if p.IdLengthConf == nil {
		p.IdLengthConf = &BBcIdConfig{}
	}

	var err error
	buf := bytes.NewBuffer(*dat)

	p.TransactionID, p.IdLengthConf.TransactionIdLength, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	if val, err := Get2byte(buf); err != nil {
		return err
	} else if val == 0 {
		p.AssetID = nil
		return nil
	}

	p.AssetID, p.IdLengthConf.AssetIdLength, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	return nil
}
