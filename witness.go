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
)

/*
BBcWitness definition

The BBcWitness has the mapping info between the userIDs and BBcSignature objects.
This object should be used if BBcRelation is used or a certain user wants to sign to the transaction in some reason.

"UserIDs" is the list of userID, and "SigIndices" is a mapping info between userID and the position (index) of the signature list in the BBcTransaction object.

"Transaction" is the pointer to the parent BBcTransaction object.

"IDLength" and "Transaction" are not included in a packed data. They are for internal use only.
*/
type (
	BBcWitness struct {
		IDLength    int
		UserIDs     [][]byte
		SigIndices  []int
		Transaction *BBcTransaction
	}
)

// Stringer outputs the content of the object
func (p *BBcWitness) Stringer() string {
	ret := "Witness:\n"
	if p.UserIDs != nil {
		for i := range p.UserIDs {
			ret += fmt.Sprintf(" [%d]\n", i)
			ret += fmt.Sprintf(" user_id: %x\n", p.UserIDs[i])
			ret += fmt.Sprintf(" sig_index: %d\n", p.SigIndices[i])
		}
	} else {
		ret += "  None (invalid)\n"
	}
	return ret
}

// SetTransaction links the BBcWitness object to the parent transaction object
func (p *BBcWitness) SetTransaction(txobj *BBcTransaction) {
	p.Transaction = txobj
}

// AddWitness makes a memo for managing signer who sign this BBcTransaction object
// This must be done before AddSignature.
func (p *BBcWitness) AddWitness(userID *[]byte) error {
	if p.Transaction == nil {
		return errors.New("transaction must be set")
	}
	p.UserIDs = append(p.UserIDs, (*userID)[:p.IDLength])
	idx := p.Transaction.GetSigIndex(*userID)
	p.SigIndices = append(p.SigIndices, idx)
	return nil
}

// AddSignature sets the BBcSignature to the parent BBcTransaction and the position in the Signatures list in BBcTransaction is based on the UserID
func (p *BBcWitness) AddSignature(userID *[]byte, sig *BBcSignature) error {
	if p.Transaction == nil {
		return errors.New("transaction must be set")
	}
	p.Transaction.AddSignature(userID, sig)
	return nil
}

// Pack returns the binary data of the BBcWitness object
func (p *BBcWitness) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	Put2byte(buf, uint16(len(p.UserIDs)))
	for i := 0; i < len(p.UserIDs); i++ {
		PutBigInt(buf, &p.UserIDs[i], p.IDLength)
		Put2byte(buf, uint16(p.SigIndices[i]))
	}

	return buf.Bytes(), nil
}

// Unpack the BBcWitness object to the binary data
func (p *BBcWitness) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	userNum, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(userNum); i++ {
		userID, err2 := GetBigInt(buf)
		if err2 != nil {
			return err2
		}
		p.UserIDs = append(p.UserIDs, userID)

		idx, err2 := Get2byte(buf)
		if err2 != nil {
			return err2
		}
		p.SigIndices = append(p.SigIndices, int(idx))
	}

	return nil
}
