// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package markdown

import (
	"testing"
)

func TestRegSubMatch(t *testing.T) {
	idx := tryBlockHeadReg.FindStringSubmatchIndex("`traffic_analytics` supports the following:")
	t.Logf("%v", idx)

	for _, val := range []string{
		"  * `abc`  def",
		"* `abc` -  something  here.  ",
	} {
		res := removeRedundantSpace(val)
		t.Logf("from `%s` => `%s`", val, res)
	}
}
