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
	"errors"
	"fmt"
	"reflect"
)

/*
BBcReference definition

The BBcReference is an input of UTXO (Unspent Transaction Output) structure and this object must accompanied by a BBcEvent object because it is an output of UTXO.

"AssetGroupID" distinguishes a type of asset, e.g., token-X, token-Y, Movie content, etc..
"TransactionID" is that of a certain transaction in the past. "EventIndexInRef" points to the BBcEvent object in the past BBcTransaction.
"SigIndices" is a mapping info between userID and the position (index) of the signature list in the BBcTransaction object.

"Transaction" is the pointer to the parent BBcTransaction object, and "RefTransaction" is the pointer to the past BBcTransaction object.

"IDLength", "Transaction", "RefTransaction" and "RefEvent" are not included in a packed data. They are for internal use only.
*/
type (
	BBcReference struct {
		IDLength        int
		AssetGroupID    []byte
		TransactionID   []byte
		EventIndexInRef uint16
		SigIndices      []int
		Transaction     *BBcTransaction
		RefTransaction  *BBcTransaction
		RefEvent        BBcEvent
	}
)

// Stringer outputs the content of the object
func (p *BBcReference) Stringer() string {
	ret := fmt.Sprintf("  asset_group_id: %x\n", p.AssetGroupID)
	ret += fmt.Sprintf("  transaction_id: %x\n", p.TransactionID)
	ret += fmt.Sprintf("  event_index_in_ref: %v\n", p.EventIndexInRef)
	ret += fmt.Sprintf("  sig_indices: %v\n", p.SigIndices)
	return ret
}

// SetTransaction links the BBcReference object to the parent transaction object
func (p *BBcReference) SetTransaction(txobj *BBcTransaction) {
	p.Transaction = txobj
}

// Add sets essential information to the BBcReference object
func (p *BBcReference) Add(assetGroupID *[]byte, refTransaction *BBcTransaction, eventIdx int) {
	if assetGroupID != nil {
		p.AssetGroupID = make([]byte, p.IDLength)
		copy(p.AssetGroupID, *assetGroupID)
	}
	if eventIdx > -1 {
		p.EventIndexInRef = uint16(eventIdx)
	}
	if refTransaction != nil {
		p.RefTransaction = refTransaction
		p.TransactionID = refTransaction.TransactionID[:p.IDLength]
		p.RefEvent = *p.RefTransaction.Events[p.EventIndexInRef]
	}
}

// AddApprover makes a memo for managing approvers who sign this BBcTransaction object
func (p *BBcReference) AddApprover(userID *[]byte) error {
	if p.Transaction == nil {
		return errors.New("transaction must be set")
	}
	if p.RefTransaction == nil {
		return errors.New("reference_transaction must be set")
	}

	flag := false
	for _, m := range p.RefEvent.MandatoryApprovers {
		if reflect.DeepEqual(m, *userID) {
			flag = true
			break
		}
	}
	if !flag {
		for _, o := range p.RefEvent.OptionApprovers {
			if reflect.DeepEqual(o, *userID) {
				flag = true
				break
			}
		}
	}
	if !flag {
		return errors.New("the user is not specified as approver")
	}

	idx := p.Transaction.GetSigIndex(*userID)
	p.SigIndices = append(p.SigIndices, idx)
	return nil
}

// AddSignature sets the BBcSignature object in the object
func (p *BBcReference) AddSignature(userID *[]byte, sig *BBcSignature) error {
	if p.Transaction == nil {
		return errors.New("transaction must be set")
	}
	p.Transaction.AddSignature(userID, sig)
	return nil
}

// Pack returns the binary data of the BBcReference object
func (p *BBcReference) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	PutBigInt(buf, &p.AssetGroupID, p.IDLength)
	PutBigInt(buf, &p.TransactionID, p.IDLength)
	Put2byte(buf, p.EventIndexInRef)
	Put2byte(buf, uint16(len(p.SigIndices)))
	for i := 0; i < len(p.SigIndices); i++ {
		Put2byte(buf, uint16(p.SigIndices[i]))
	}

	return buf.Bytes(), nil
}

// Unpack the BBcReference object to the binary data
func (p *BBcReference) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	p.AssetGroupID, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	p.TransactionID, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	p.EventIndexInRef, err = Get2byte(buf)
	if err != nil {
		return err
	}

	sigNum, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(sigNum); i++ {
		idx, err := Get2byte(buf)
		if err != nil {
			return err
		}
		p.SigIndices = append(p.SigIndices, int(idx))
	}

	return nil
}
