// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "testing"

func TestSL007(t *testing.T) {
	fs := runRule(t, sl007, wrap(`
			"list_arr": {Type: pluginsdk.TypeList, Optional: true, Description: "d", Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString}},
			"set_arr":  {Type: pluginsdk.TypeSet, Required: true, Description: "d", Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString}},
			"bounded":  {Type: pluginsdk.TypeList, Optional: true, Description: "d", MaxItems: 5, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString}},
			"block": {
				Type: pluginsdk.TypeList, Optional: true, Description: "d",
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"a": {Type: pluginsdk.TypeString, Required: true, Description: "d"},
				}},
			},
			"computed": {Type: pluginsdk.TypeList, Computed: true, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString}},
			"scalar":   {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
	`))

	if !flagged(fs, "list_arr") || !flagged(fs, "set_arr") {
		t.Error("expected SL007 on unbounded scalar arrays")
	}
	for _, p := range []string{"bounded", "block", "computed", "scalar"} {
		if flagged(fs, p) {
			t.Errorf("did not expect SL007 on %q", p)
		}
	}
}
