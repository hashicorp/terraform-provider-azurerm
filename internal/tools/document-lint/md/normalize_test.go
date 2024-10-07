// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package md

import (
	"os"
	"path"
	"testing"
)

func Test_fixFileNormalize(t *testing.T) {
	t.Skipf("skip normalize unit test")
	dir, err := os.ReadDir(ResourceDir())
	_ = err
	for _, en := range dir {
		if en.IsDir() {
			continue
		}
		fullPath := path.Join(ResourceDir(), en.Name())
		FixFileNormalize(fullPath)
	}
}

func TestMDFile(t *testing.T) {
	file := "automation_watcher.html.markdown"
	FixFileNormalize(path.Join(ResourceDir(), file))
}

func TestRegSubMatch(t *testing.T) {
	idx := oldBlockHeadReg.FindStringSubmatchIndex("`traffic_analytics` supports the following:")
	t.Logf("%v", idx)

	for _, val := range []string{
		"  * `abc`  def",
		"* `abc` -  something  here.  ",
	} {
		res := removeRedundantSpace(val)
		t.Logf("from `%s` => `%s`", val, res)
	}
}
