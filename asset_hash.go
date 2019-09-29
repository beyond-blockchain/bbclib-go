/*
Copyright (c) 2019 Zettant Inc.

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
BBcAssetHash definition

"IDLength" is not included in a packed data. It is for internal use only.

Multiple AssetIDs can be cotained in the list.
The length of "AssetID" is defined by "IDLength".
*/
type (
	BBcAssetHash struct {
		IdLengthConf      *BBcIdConfig
		AssetIdNum        uint16
		AssetIDs          [][]byte
	}
)

// Stringer outputs the content of the object
func (p *BBcAssetHash) Stringer() string {
	ret := "  AssetHash:\n"
	ret += fmt.Sprintf("     num_of_asset_ids: %d\n", p.AssetIdNum)
	ret += "  AssetIDs:\n"
	if p.AssetIDs != nil {
		for _, a := range p.AssetIDs {
			ret += fmt.Sprintf("    - %x\n", a)
		}
	} else {
		ret += "    - None\n"
	}
	return ret
}

// Set ID length configuration
func (p *BBcAssetHash) SetIdLengthConf(conf * BBcIdConfig) {
	p.IdLengthConf = conf
}

// AddAssetId sets a string data in the BBcAsset object
func (p *BBcAssetHash) AddAssetId(assetId *[]byte) {
	asid := make([]byte, p.IdLengthConf.AssetIdLength)
	copy(asid, *assetId)
	p.AssetIDs = append(p.AssetIDs, asid)
	p.AssetIdNum += 1
}

// Pack returns the binary data of the BBcAsset object
func (p *BBcAssetHash) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)
	Put2byte(buf, p.AssetIdNum)
	for i := 0; i < int(p.AssetIdNum); i++ {
		PutBigInt(buf, &p.AssetIDs[i], p.IdLengthConf.AssetIdLength)
	}
	return buf.Bytes(), nil
}

// Unpack the BBcAsset object to the binary data
func (p *BBcAssetHash) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)
	p.AssetIdNum, err = Get2byte(buf)
	if err != nil {
		return err
	}

	for i := 0; i < int(p.AssetIdNum); i++ {
		assetId, ulen, err := GetBigInt(buf)
		if err != nil {
			return err
		}
		p.IdLengthConf.UserIdLength = ulen
		p.AssetIDs = append(p.AssetIDs, assetId)
	}
	return nil
}
