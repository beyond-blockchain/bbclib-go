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
	"testing"
)

func TestAssetRawPackUnpack(t *testing.T) {
	var idLengthConfig = BBcIdConfig {
		TransactionIdLength: 32,
		UserIdLength: 32,
		AssetGroupIdLength: 32,
		AssetIdLength: 32,
		NonceLength: 32,
	}

	t.Run("simple creation (string)", func(t *testing.T) {
		obj := BBcAssetRaw{}
		obj.SetIdLengthConf(&idLengthConfig)
		asid := GetIdentifier("user1_789abcdef0123456789abcdef0", idLengthConfig.AssetIdLength)
		obj.AddBody(&asid, "testString12345XXX")
		t.Log("--------------------------------------")
		t.Logf("id_length_config: %v", obj.IdLengthConf)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcAssetRaw{}
		obj2.SetIdLengthConf(&idLengthConfig)
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length_config: %v", obj2.IdLengthConf)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.AssetID, obj2.AssetID) != 0 || bytes.Compare(obj.AssetBody, obj2.AssetBody) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("simple creation (byte)", func(t *testing.T) {
		obj := BBcAssetRaw{}
		obj.SetIdLengthConf(&idLengthConfig)
		asid := GetIdentifier("user1_789abcdef0123456789abcdef0", idLengthConfig.AssetIdLength)
		obj.AddBody(&asid, []byte("testString12345XXX"))
		t.Log("--------------------------------------")
		t.Logf("id_length_config: %v", obj.IdLengthConf)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcAssetRaw{}
		obj2.SetIdLengthConf(&idLengthConfig)
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length_config: %v", obj2.IdLengthConf)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.AssetID, obj2.AssetID) != 0 || bytes.Compare(obj.AssetBody, obj2.AssetBody) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}
