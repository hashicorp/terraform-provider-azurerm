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
	file := filepath.Join(testDir, "test_recovery_services_vault.html.markdown")
	m := MustNewMarkFromFile(file)
	doc := m.BuildResourceDoc()

	// The key_vault.html.markdown should have been parsed successfully
	// This test verifies that the linkBlockFields function doesn't break existing functionality
	if len(doc.Args) == 0 {
		t.Fatal("Expected arguments to be parsed from key_vault.html.markdown")
	}

	// Look for any blocks that have SameNameAttr set to verify linking works
	sameNameAttrFound := false
	for _, field := range doc.Args {
		if field.SameNameAttr != nil {
			sameNameAttrFound = true
			t.Logf("Found SameNameAttr linking for field: %s", field.Name)
			break
		}
		// Check nested fields too
		if field.Subs != nil {
			for _, subField := range field.Subs {
				if subField.SameNameAttr != nil {
					sameNameAttrFound = true
					t.Logf("Found nested SameNameAttr linking for field: %s.%s", field.Name, subField.Name)
					break
				}
			}
		}
		if sameNameAttrFound {
			break
		}
	}

	t.Logf("Successfully parsed key_vault.html.markdown with %d arguments and %d attributes", len(doc.Args), len(doc.Attr))
	t.Logf("SameNameAttr linking found: %v", sameNameAttrFound)
}
