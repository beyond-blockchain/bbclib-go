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
	"fmt"
)

/*
BBcRelation definition

The BBcRelation holds the asset (by BBcAsset) and the relationship with the other transaction/asset (by BBcPointer).
Different from UTXO, state information or account-type information can be expressed by using this object.
If you want to include signature(s) according to the contents of BBcRelation object, BBcWitness should be included in the transaction object.

"AssetGroupID" distinguishes a type of asset, e.g., token-X, token-Y, Movie content, etc..
"Pointers" is a list of BBcPointers object. "Asset" is a BBcAsset object.

"IDLength" is not included in a packed data. It is for internal use only.
*/
type (
	BBcRelation struct {
		IDLength     int
		AssetGroupID []byte
		Pointers     []*BBcPointer
		Asset        *BBcAsset
	}
)

// Stringer outputs the content of the object
func (p *BBcRelation) Stringer() string {
	ret := fmt.Sprintf("  asset_group_id: %x\n", p.AssetGroupID)
	if p.Pointers != nil {
		ret += fmt.Sprintf("  Pointers[]: %d\n", len(p.Pointers))
		for i := range p.Pointers {
			ret += fmt.Sprintf("   [%d]\n", i)
			ret += p.Pointers[i].Stringer()
		}
	} else {
		ret += fmt.Sprintf("  Pointers[]: None\n")
	}
	if p.Asset != nil {
		ret += p.Asset.Stringer()
	} else {
		ret += fmt.Sprintf("  Asset: None\n")
	}
	return ret
}

// Add sets essential information (assetGroupID and BBcAsset object) to the BBcRelation object
func (p *BBcRelation) Add(assetGroupID *[]byte, asset *BBcAsset) {
	if assetGroupID != nil {
		p.AssetGroupID = make([]byte, p.IDLength)
		copy(p.AssetGroupID, *assetGroupID)
	}
	if asset != nil {
		p.Asset = asset
		p.Asset.IDLength = p.IDLength
	}
}

// AddPointer sets the BBcPointer object in the object
func (p *BBcRelation) AddPointer(pointer *BBcPointer) {
	pointer.IDLength = p.IDLength
	p.Pointers = append(p.Pointers, pointer)
}

// Pack returns the binary data of the BBcRelation object
func (p *BBcRelation) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	PutBigInt(buf, &p.AssetGroupID, p.IDLength)

	Put2byte(buf, uint16(len(p.Pointers)))
	for _, p := range p.Pointers {
		dat, er := p.Pack()
		if er != nil {
			return nil, er
		}
		Put2byte(buf, uint16(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return nil, err
		}
	}
	if p.Asset != nil {
		ast, er := p.Asset.Pack()
		if er != nil {
			return nil, er
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

// Unpack the BBcRelation object to the binary data
func (p *BBcRelation) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	p.AssetGroupID, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	numPointers, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(numPointers); i++ {
		size, err2 := Get2byte(buf)
		if err2 != nil {
			return err2
		}
		ptr, _ := GetBytes(buf, int(size))
		pointer := BBcPointer{IDLength: p.IDLength}
		pointer.Unpack(&ptr)
		p.Pointers = append(p.Pointers, &pointer)
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
