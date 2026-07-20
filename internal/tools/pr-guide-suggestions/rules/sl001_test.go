// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "testing"

func TestSL001(t *testing.T) {
	fs := runRule(t, sl001, wrap(`
			"no_desc":  {Type: pluginsdk.TypeString, Optional: true},
			"empty":    {Type: pluginsdk.TypeString, Optional: true, Description: ""},
			"has_desc": {Type: pluginsdk.TypeString, Optional: true, Description: "x"},
			"computed": {Type: pluginsdk.TypeString, Computed: true},
			"helper":   commonschema.Location(),
			"block": {
				Type: pluginsdk.TypeList, Optional: true, MaxItems: 1,
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"child": {Type: pluginsdk.TypeString, Required: true},
				}},
			},
	`))

	// Missing or empty descriptions are flagged, including computed-only
	// attributes, blocks and nested children.
	for _, p := range []string{"no_desc", "empty", "computed", "block", "block.child"} {
		if !flagged(fs, p) {
			t.Errorf("expected SL001 on %q", p)
		}
	}
	// A non-empty description and an opaque helper call are not flagged.
	if flagged(fs, "has_desc") {
		t.Error("did not expect SL001 on has_desc")
	}
	if flagged(fs, "helper") {
		t.Error("did not expect SL001 on a helper-call property")
	}
}
