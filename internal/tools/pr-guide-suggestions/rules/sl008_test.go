// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "testing"

func TestSL008Multiple(t *testing.T) {
	fs := runRule(t, sl008, wrap(`
			"sku_name":     {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
			"sku_capacity": {Type: pluginsdk.TypeInt, Optional: true, Description: "d"},
			"block": {
				Type: pluginsdk.TypeList, Optional: true, Description: "d",
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"sku_a": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
					"sku_b": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
				}},
			},
	`))

	if !flagged(fs, "sku_name") || !flagged(fs, "sku_capacity") {
		t.Error("expected SL008 on both top-level sku_* fields")
	}
	// Nested sku_* fields are out of scope even when multiple.
	for _, p := range []string{"block.sku_a", "block.sku_b"} {
		if flagged(fs, p) {
			t.Errorf("did not expect SL008 on nested %q", p)
		}
	}
}

func TestSL008Single(t *testing.T) {
	fs := runRule(t, sl008, wrap(`
			"sku_name": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
			"other":    {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
	`))
	if flagged(fs, "sku_name") {
		t.Error("did not expect SL008 on a lone sku_* field")
	}
}
