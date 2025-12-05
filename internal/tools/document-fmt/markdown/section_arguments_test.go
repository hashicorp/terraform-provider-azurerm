// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package markdown

import (
	"testing"
)

func TestNormalizeArgumentsContent(t *testing.T) {
	for _, val := range []string{
		"* `name` (Optional) - The name.",
		"* `id` (Required) - The ID.",
		"* `name` - (optional) The name.",
		"* `id` - (required) The ID.",
		"* `name`- (Required) The name.",
		"* `port`- (Required) The port. Defaults to `-1`.",
	} {
		res, _ := normalizeArgumentsContent([]string{val})
		t.Logf("from `%s` => `%s`", val, res[0])
	}
}
