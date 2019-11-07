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
		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		obj := BBcRelation{AssetGroupID: assetgroup}
		obj.SetIdLengthConf(&idLengthConfig)
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIDLength)
		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)

		obj.AddAsset(&u1, nil, "testString12345XXX").AddPointer(&txid1, &asid1).AddPointer(&txid2, nil)

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
		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIDLength)
		obj := BBcRelation{AssetGroupID: assetgroup}
		obj.SetIdLengthConf(&idLengthConfig)

		obj.AddAsset(&u1, nil, map[int]string{1: "aaa", 2: "bbb", 10: "asdfasdfasf;lakj;lkj;"})

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

	t.Run("simple creation (asset_raw)", func(t *testing.T) {
		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		asid := GetIdentifier("user1_789abcdef0123456789abcdef0", idLengthConfig.AssetIdLength)
		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)

		obj := BBcRelation{Version: 2, AssetGroupID: assetgroup}
		obj.SetIdLengthConf(&idLengthConfig)

		obj.AddPointer(&txid1, &asid1).AddPointer(&txid2, nil).AddAssetRaw(&asid, []byte("testString12345XXX"))

		t.Log("---------------Relation-----------------")
		t.Logf("id_length_config: %v", obj.IdLengthConf)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcRelation{Version:2}
		obj2.SetIdLengthConf(&idLengthConfig)
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length_config: %v", obj2.IdLengthConf)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.AssetGroupID, obj2.AssetGroupID) != 0 || bytes.Compare(obj.AssetRaw.AssetID, obj2.AssetRaw.AssetID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
		if bytes.Compare(obj.AssetRaw.AssetID, obj2.AssetRaw.AssetID) != 0 || bytes.Compare(obj.AssetRaw.AssetBody, obj2.AssetRaw.AssetBody) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("simple creation (asset_hash)", func(t *testing.T) {
		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIDLength)
		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIDLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIDLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIDLength)

		obj := BBcRelation{Version: 2, AssetGroupID: assetgroup}
		obj.SetIdLengthConf(&idLengthConfig)

		obj.AddPointer(&txid1, &asid1).AddPointer(&txid2, nil)

		for i := 0; i < 10; i++ {
			asid := GetIdentifier(fmt.Sprintf("asset_id_%d", i), idLengthConfig.AssetIdLength)
			obj.AddAssetHash(&asid)
		}

		t.Log("---------------Relation-----------------")
		t.Logf("id_length_config: %v", obj.IdLengthConf)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcRelation{Version:2}
		obj2.SetIdLengthConf(&idLengthConfig)
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length_config: %v", obj2.IdLengthConf)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.AssetGroupID, obj2.AssetGroupID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
		if obj2.AssetHash.AssetIdNum != 10 {
			t.Fatal("Not recovered correctly...")
		}
		for i := 0; i < 10; i++ {
			if bytes.Compare(obj.AssetHash.AssetIDs[i], obj2.AssetHash.AssetIDs[i]) != 0 {
				t.Fatal("Not recovered correctly...")
			}
		}
	})
}
