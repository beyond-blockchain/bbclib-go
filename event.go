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
	"encoding/binary"
	"errors"
	"fmt"
)

/*
BBcEvent definition

BBcEvent expresses an output of UTXO (Unspent Transaction Output) structure.

"AssetGroupID" distinguishes a type of asset, e.g., token-X, token-Y, Movie content, etc..

"ReferenceIndices" has the index numbers in BBcReference object list in the transaction object.
It expresses that this BBcEvent object has a certain relationship with the BBcReference objects specified by ReferenceIndices.
This would be used in the case that the transaction object has multiple BBcReference objects.

BBcEvent designates Mandatory/Option Approvers to be signers to a BBcTransaction in the future, which use the asset in the BBcEvent.
As for "OptionApprovers", it is enough that some of them give sign to the BBcTransaction. The number of signers to be included is "OptionApproverNumNumerator".

Asset is the most important part of the BBcTransaction. The BBcAsset object includes the digital asset to be protected by BBc-1.

"IDLength" is not included in a packed data. It is for internal use only.
*/
type (
	BBcEvent struct {
		IDLength                     int
		AssetGroupID                 []byte
		ReferenceIndices             []int
		MandatoryApprovers           [][]byte
		OptionApproverNumNumerator   uint16
		OptionApproverNumDenominator uint16
		OptionApprovers              [][]byte
		Asset                        *BBcAsset
	}
)

// Stringer outputs the content of the object
func (p *BBcEvent) Stringer() string {
	ret := fmt.Sprintf("  asset_group_id: %x\n", p.AssetGroupID)
	if p.ReferenceIndices != nil {
		ret += fmt.Sprintf("  reference_indices: %v\n", p.ReferenceIndices)
	} else {
		ret += fmt.Sprintf("  reference_indices: None\n")
	}
	ret += "  mandatory_approvers:\n"
	if p.MandatoryApprovers != nil {
		for _, a := range p.MandatoryApprovers {
			ret += fmt.Sprintf("    - %x\n", a)
		}
	} else {
		ret += "    - None\n"
	}
	ret += "  option_approvers:\n"
	if p.OptionApprovers != nil {
		for _, o := range p.OptionApprovers {
			ret += fmt.Sprintf("    - %x\n", o)
		}
	} else {
		ret += "    - None\n"
	}
	ret += fmt.Sprintf("  option_approver_num_numerator: %d\n", p.OptionApproverNumNumerator)
	ret += fmt.Sprintf("  option_approver_num_denominator: %d\n", p.OptionApproverNumDenominator)
	if p.Asset != nil {
		ret += p.Asset.Stringer()
	} else {
		ret += fmt.Sprintf("  Asset: None\n")
	}
	return ret
}

// Add sets essential information to the BBcEvent object
func (p *BBcEvent) Add(assetGroupID *[]byte, asset *BBcAsset) {
	if assetGroupID != nil {
		p.AssetGroupID = make([]byte, p.IDLength)
		copy(p.AssetGroupID, (*assetGroupID)[:p.IDLength])
	}
	if asset != nil {
		p.Asset = asset
		p.Asset.IDLength = p.IDLength
	}
}

// AddReferenceIndex sets an index to ReferenceIndices of the BBcEvent object
func (p *BBcEvent) AddReferenceIndex(relIndex int) {
	if relIndex != -1 {
		p.ReferenceIndices = append(p.ReferenceIndices, relIndex)
	}
}

// AddOptionParams sets values to OptionApproverNumNumerator and OptionApproverNumDenominator in the BBcEvent object
func (p *BBcEvent) AddOptionParams(numerator int, denominator int) {
	p.OptionApproverNumNumerator = uint16(numerator)
	p.OptionApproverNumDenominator = uint16(denominator)
}

// AddMandatoryApprover sets userID in MandatoryApprover list of the BBcEvent object
func (p *BBcEvent) AddMandatoryApprover(userID *[]byte) {
	uid := make([]byte, p.IDLength)
	copy(uid, *userID)
	p.MandatoryApprovers = append(p.MandatoryApprovers, uid)
}

// AddOptionApprover sets userID in OptionApprover list of the BBcEvent object
func (p *BBcEvent) AddOptionApprover(userID *[]byte) {
	uid := make([]byte, p.IDLength)
	copy(uid, *userID)
	p.OptionApprovers = append(p.OptionApprovers, uid)
}

// Pack returns the binary data of the BBcEvent object
func (p *BBcEvent) Pack() ([]byte, error) {
	if len(p.OptionApprovers) != int(p.OptionApproverNumDenominator) {
		return nil, errors.New("num of option approvers must be equal to OptionApproverNumDenominator")
	}
	buf := new(bytes.Buffer)

	PutBigInt(buf, &p.AssetGroupID, p.IDLength)

	Put2byte(buf, uint16(len(p.ReferenceIndices)))
	for i := 0; i < len(p.ReferenceIndices); i++ {
		Put2byte(buf, uint16(p.ReferenceIndices[i]))
	}

	Put2byte(buf, uint16(len(p.MandatoryApprovers)))
	for i := 0; i < len(p.MandatoryApprovers); i++ {
		PutBigInt(buf, &p.MandatoryApprovers[i], p.IDLength)
	}

	Put2byte(buf, p.OptionApproverNumNumerator)
	Put2byte(buf, p.OptionApproverNumDenominator)
	for i := 0; i < int(p.OptionApproverNumDenominator); i++ {
		PutBigInt(buf, &p.OptionApprovers[i], p.IDLength)
	}

	if p.Asset != nil {
		ast, err := p.Asset.Pack()
		if err != nil {
			return nil, err
		}
		Put4byte(buf, uint32(binary.Size(ast)))
		if err := binary.Write(buf, binary.LittleEndian, ast); err != nil {
			return nil, err
		}
	} else {
		Put4byte(buf, 0)
	}

	return buf.Bytes(), nil
}

// unpackApprovers unpacks the approver part of the binary data
func (p *BBcEvent) unpackApprovers(buf *bytes.Buffer) error {
	numMandatory, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(numMandatory); i++ {
		userID, err2 := GetBigInt(buf)
		if err2 != nil {
			return err2
		}
		p.MandatoryApprovers = append(p.MandatoryApprovers, userID)
	}

	p.OptionApproverNumNumerator, err = Get2byte(buf)
	if err != nil {
		return err
	}
	p.OptionApproverNumDenominator, err = Get2byte(buf)
	if err != nil {
		return err
	}

	for i := 0; i < int(p.OptionApproverNumDenominator); i++ {
		userID, err2 := GetBigInt(buf)
		if err2 != nil {
			return err2
		}
		p.OptionApprovers = append(p.OptionApprovers, userID)
	}
	return nil
}

// Unpack the binary data to the BBcEvent object
func (p *BBcEvent) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	p.AssetGroupID, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	numReferences, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(numReferences); i++ {
		idx, err2 := Get2byte(buf)
		if err2 != nil {
			return err2
		}
		p.ReferenceIndices = append(p.ReferenceIndices, int(idx))
	}

	if err = p.unpackApprovers(buf); err != nil {
		return err
	}

	assetSize, err := Get4byte(buf)
	if err != nil {
		return err
	}
	if assetSize > 0 {
		ast, err := GetBytes(buf, int(assetSize))
		if err != nil {
			return err
		}
		p.Asset = &BBcAsset{IDLength: p.IDLength}
		p.Asset.Unpack(&ast)
	}

	return nil
}
