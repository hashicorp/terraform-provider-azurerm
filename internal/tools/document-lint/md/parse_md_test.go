// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package md

import (
	"path/filepath"
	"runtime"
	"testing"
)

var testDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	testDir = filepath.Join(filepath.Dir(file), "testdata")
}

func Test_unmarshalFile(t *testing.T) {
	args := []struct {
		file    string
		itemNum int
		argsNum int
	}{
		{"key_vault.html.markdown", 64, 16},
		{"media_transform.html.markdown", 270, 5},
	}
	for _, arg := range args {
		file := filepath.Join(testDir, arg.file)
		m := MustNewMarkFromFile(file)
		if gotItems := len(m.Items); gotItems != arg.itemNum {
			t.Fatalf("%s expect item num; %d, got: %d", arg.file, gotItems, arg.itemNum)
		}
		doc := m.BuildResourceDoc()
		if gotArgs := len(doc.Args); gotArgs != arg.argsNum {
			t.Fatalf("`%s` expect arg num: %d, got: %d", arg.file, gotArgs, arg.argsNum)
		}
	}
}

func TestSameNameAttrLinking(t *testing.T) {
	// Test that linkBlockFields is working by using existing test file
	file := filepath.Join(testDir, "test_identity.html.markdown")
	m := MustNewMarkFromFile(file)
	doc := m.BuildResourceDoc()

	// Look for any blocks that have SameNameAttr set to verify linking works
	sameNameAttrFound := false
	for key, field := range doc.Args {
		if key == "identity" && field.SameNameAttr != nil {
			sameNameAttrFound = true
			break
		}
	}

	if !sameNameAttrFound {
		t.Fatalf("Expected to link sameNameAttr to Args, but not found")
	}
}
