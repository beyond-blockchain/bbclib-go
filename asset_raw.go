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
	"encoding/binary"
	"fmt"
)

/*
BBcAssetRaw definition

"IDLength" is not included in a packed data. It is for internal use only.

"AssetID" is externally calculated digest value.
The length of "AssetID" is defined by "IDLength".
*/
type (
	BBcAssetRaw struct {
		IdLengthConf      *BBcIdConfig
		AssetID           []byte
		AssetBodySize     uint16
		AssetBody         []byte
	}
)

// Stringer outputs the content of the object
func (p *BBcAssetRaw) Stringer() string {
	ret := "  AssetRaw:\n"
	ret += fmt.Sprintf("     asset_id: %x\n", p.AssetID)
	ret += fmt.Sprintf("     body_size: %d\n", p.AssetBodySize)
	ret += fmt.Sprintf("     body: %v\n", p.AssetBody)
	return ret
}

// Set ID length configuration
func (p *BBcAssetRaw) SetIdLengthConf(conf * BBcIdConfig) {
	p.IdLengthConf = conf
}

// AddBodyString sets a string data in the BBcAsset object
func (p *BBcAssetRaw) AddBody(assetID *[]byte, assetBody interface{}) {
	if assetID != nil {
		p.AssetID = make([]byte, p.IdLengthConf.AssetIdLength)
		copy(p.AssetID, *assetID)
	}
	switch assetBody.(type) {
	case string:
		p.AssetBody = []byte(assetBody.(string))
		p.AssetBodySize = uint16(len(p.AssetBody))
		break
	case []byte:
		p.AssetBody = assetBody.([]byte)
		p.AssetBodySize = uint16(len(p.AssetBody))
		break
	}
}


// Pack returns the binary data of the BBcAsset object
func (p *BBcAssetRaw) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)
	PutBigInt(buf, &p.AssetID, p.IdLengthConf.AssetIdLength)
	Put2byte(buf, p.AssetBodySize)
	if p.AssetBodySize > 0 {
		if err := binary.Write(buf, binary.LittleEndian, p.AssetBody); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// Unpack the BBcAsset object to the binary data
func (p *BBcAssetRaw) Unpack(dat *[]byte) error {
	if p.IdLengthConf == nil {
		p.IdLengthConf = &BBcIdConfig{}
	}

	var err error
	buf := bytes.NewBuffer(*dat)

	p.AssetID, p.IdLengthConf.AssetIdLength, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	p.AssetBodySize, err = Get2byte(buf)
	if err != nil {
		return err
	}
	p.AssetBody, _, err = GetBytes(buf, int(p.AssetBodySize))

	return err
}
