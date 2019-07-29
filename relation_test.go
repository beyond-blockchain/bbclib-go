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

func TestRelationPackUnpack(t *testing.T) {
	var idLengthConfig = BBcIdConfig {
		TransactionIdLength: 32,
		UserIdLength: 32,
		AssetGroupIdLength: 32,
		AssetIdLength: 32,
		NonceLength: 32,
	}

	t.Run("simple creation (string asset)", func(t *testing.T) {
		obj := BBcRelation{}
		obj.SetIdLengthConf(&idLengthConfig)
		ptr1 := BBcPointer{}
		ptr2 := BBcPointer{}
		ast := BBcAsset{}

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		obj.Add(&assetgroup, &ast)
		obj.AddPointer(&ptr1)
		obj.AddPointer(&ptr2)

		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIDLength)
		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")
		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)
		ptr1.Add(&txid1, &asid1)
		ptr2.Add(&txid2, nil)

		t.Log("---------------Relation-----------------")
		t.Logf("id_length_config: %v", obj.IdLengthConf)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcRelation{}
		obj2.SetIdLengthConf(&idLengthConfig)
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length_config: %v", obj2.IdLengthConf)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.AssetGroupID, obj2.AssetGroupID) != 0 || bytes.Compare(obj.Asset.AssetID, obj2.Asset.AssetID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("simple creation (no pointer, msgpack asset)", func(t *testing.T) {
		ast := BBcAsset{}
		ast.SetIdLengthConf(&idLengthConfig)
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIDLength)
		ast.Add(&u1)
		ast.AddBodyObject(map[int]string{1: "aaa", 2: "bbb", 10: "asdfasdfasf;lakj;lkj;"})

		obj := BBcRelation{}
		obj.SetIdLengthConf(&idLengthConfig)
		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		obj.Add(&assetgroup, &ast)
		t.Log("---------------Relation-----------------")
		t.Logf("id_length_config: %v", obj.IdLengthConf)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcRelation{}
		obj2.SetIdLengthConf(&idLengthConfig)
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length_config: %v", obj2.IdLengthConf)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.AssetGroupID, obj2.AssetGroupID) != 0 || bytes.Compare(obj.Asset.AssetID, obj2.Asset.AssetID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}
