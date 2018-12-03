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
	"reflect"
	"testing"
)

func TestCompressDecompress(t *testing.T) {
	original := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	comp := ZlibCompress(&original)
	t.Logf("compressed: %x\n", comp)

	decomp, err := ZlibDecompress(comp)
	if err != nil {
		t.Fatalf("failed to decompress (%v)", err)
	}
	t.Logf("compressed: %x\n", decomp)
	if !reflect.DeepEqual(original, decomp) {
		t.Fatal("failed to decompress (mismatch)")
	}
	t.Log("Succeeded")
}
