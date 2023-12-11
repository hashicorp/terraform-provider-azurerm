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
