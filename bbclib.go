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

/*
Package bbclib is a library for defining BBcTransaction. This also provides serializer/deserializer and utilities for BBcTransaction object manipulation.


Serialization and deserialization

A BBcTransaction object contains various object, such as BBcEvent, BBcSignature.
In order to store a BBcTransaction object in DB or send it to other host, the object must be serialized.
Before serialization, the object is packed, meaning that it is transformed into binary format.
Then, the header is prepended to the packed data, resulting in a serialized data.
According to the header value, the packed data is compressed, so that you will get a smaller-sized serialized data.
Deserialization is the opposite transformation to serialization.

Utility functions

To build a BBcTransaction you need to create (new) objects you want to include. In many cases, it is a kind of common coding manner.
The utility functions are helpers to build a BBcTransaction with various objects.
*/
package bbclib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

// Header values for serialized data
const (
	FormatPlain = 0x0000
	FormatZlib  = 0x0010
)

type (
	BBcIdConfig struct {
		TransactionIdLength	  int
		UserIdLength          int
		AssetGroupIdLength    int
		AssetIdLength         int
		NonceLength           int
	}
)

const (
	defaultIDLength = 32
)

var IdLengthConfig = BBcIdConfig {
	TransactionIdLength: defaultIDLength,
	UserIdLength: defaultIDLength,
	AssetGroupIdLength: defaultIDLength,
	AssetIdLength: defaultIDLength,
	NonceLength: defaultIDLength,
}


// Configure various ID length
func ConfigureIdLength(conf *BBcIdConfig) {
	if conf.TransactionIdLength > 0 && conf.TransactionIdLength < 33 {
		IdLengthConfig.TransactionIdLength = conf.TransactionIdLength
	}
	if conf.UserIdLength > 0 && conf.UserIdLength < 33 {
		IdLengthConfig.UserIdLength = conf.UserIdLength
	}
	if conf.AssetGroupIdLength > 0 && conf.AssetGroupIdLength < 33 {
		IdLengthConfig.AssetGroupIdLength = conf.AssetGroupIdLength
	}
	if conf.AssetIdLength > 0 && conf.AssetIdLength < 33 {
		IdLengthConfig.AssetIdLength = conf.AssetIdLength
	}
	if conf.NonceLength > 0 && conf.NonceLength < 33 {
		IdLengthConfig.NonceLength = conf.NonceLength
	}
}

// Configure all kind of ID length with the same value
func ConfigureIdLengthAll(length int) {
	if length > 0 && length < 33 {
		IdLengthConfig.TransactionIdLength = length
		IdLengthConfig.UserIdLength = length
		IdLengthConfig.AssetGroupIdLength = length
		IdLengthConfig.AssetIdLength = length
		IdLengthConfig.NonceLength = length
	}
}

// Copy new config in refer to main
func UpdateIdLengthConfig(main, refer *BBcIdConfig) {
	if refer == nil {
		return
	}
	if refer.TransactionIdLength > 0 && main.TransactionIdLength != refer.TransactionIdLength {
		main.TransactionIdLength = refer.TransactionIdLength
	}
	if refer.UserIdLength > 0 && main.UserIdLength != refer.UserIdLength {
		main.UserIdLength = refer.UserIdLength
	}
	if refer.AssetGroupIdLength > 0 && main.AssetGroupIdLength != refer.AssetGroupIdLength {
		main.AssetGroupIdLength = refer.AssetGroupIdLength
	}
	if refer.AssetIdLength > 0 && main.AssetIdLength != refer.AssetIdLength {
		main.AssetIdLength = refer.AssetIdLength
	}
	if refer.NonceLength > 0 && main.NonceLength != refer.NonceLength {
		main.NonceLength = refer.NonceLength
	}
}

/*
Serialize BBcTransaction object into packed data

formatType = 0x0000: Packed data is simply used for serialized data.

formatType = 0x0010: Packed data is compressed using zlib, and the compressed data is used for serialized data.
*/
func Serialize(transaction *BBcTransaction, formatType uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	Put2byte(buf, formatType)
	dat, err := transaction.Pack()
	if err != nil {
		return nil, err
	}

	if formatType == FormatPlain {
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return nil, err
		}
	} else if formatType == FormatZlib {
		compressed := ZlibCompress(&dat)
		if err := binary.Write(buf, binary.LittleEndian, compressed); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// Deserialize BBcTransaction data with header
func Deserialize(dat []byte) (*BBcTransaction, error) {
	buf := bytes.NewBuffer(dat)

	formatType, err := Get2byte(buf)
	if err != nil {
		return nil, err
	}

	txdat, _, err := GetBytes(buf, len(dat)-2)
	if err != nil {
		return nil, err
	}

	if formatType == FormatPlain {
		txobj := BBcTransaction{}
		err2 := txobj.Unpack(&txdat)
		return &txobj, err2
	} else if formatType == FormatZlib {
		decompressed, err := ZlibDecompress(txdat)
		if err != nil {
			return nil, err
		}
		txobj := BBcTransaction{}
		err2 := txobj.Unpack(&decompressed)
		return &txobj, err2
	}

	return nil, errors.New("formatType not supported")
}


// MakeTransaction is a utility for making simple BBcTransaction object with BBcEvent, BBcRelation or/and BBcWitness
func MakeTransaction(eventNum, relationNum int, witness bool) *BBcTransaction {
	txobj := BBcTransaction{Version: 2}
	txobj.SetIdLengthConf(&IdLengthConfig)
	txobj.Timestamp = time.Now().UnixNano() / int64(time.Microsecond)

	for i := 0; i < eventNum; i++ {
		evt := BBcEvent{Version: txobj.Version}
		evt.SetIdLengthConf(&txobj.IdLengthConf)
		txobj.Events = append(txobj.Events, &evt)
	}

	for i := 0; i < relationNum; i++ {
		rtn := BBcRelation{Version: txobj.Version}
		rtn.SetIdLengthConf(&txobj.IdLengthConf)
		txobj.Relations = append(txobj.Relations, &rtn)
	}

	if witness {
		wit := BBcWitness{Version: txobj.Version}
		wit.SetIdLengthConf(&txobj.IdLengthConf)
		wit.Transaction = &txobj
		txobj.Witness = &wit
	}

	return &txobj
}
