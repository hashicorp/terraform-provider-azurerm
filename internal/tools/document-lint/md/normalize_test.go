// Copyright IBM Corp. 2014, 2025
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
		FixFileNormalize(path.Join(ResourceDir(), en.Name()))
	}
}

func TestMDFile(t *testing.T) {
	file := "automation_watcher.html.markdown"
	FixFileNormalize(path.Join(ResourceDir(), file))
}

func TestRegSubMatch(t *testing.T) {
	t.Logf("%v", oldBlockHeadReg.FindStringSubmatchIndex("`traffic_analytics` supports the following:"))

	for _, val := range []string{
		"  * `abc`  def",
		"* `abc` -  something  here.  ",
	} {
		res := removeRedundantSpace(val)
		t.Logf("from `%s` => `%s`", val, res)
	}
}
