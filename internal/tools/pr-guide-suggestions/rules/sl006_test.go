// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "testing"

func TestSL006(t *testing.T) {
	fs := runRule(t, sl006, wrap(`
			"loose": {
				Type: pluginsdk.TypeList, Optional: true, Description: "d",
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"a": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
					"b": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
				}},
			},
			"has_required": {
				Type: pluginsdk.TypeList, Optional: true, Description: "d",
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"a": {Type: pluginsdk.TypeString, Required: true, Description: "d"},
				}},
			},
			"at_least_one": {
				Type: pluginsdk.TypeList, Optional: true, Description: "d",
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"a": {Type: pluginsdk.TypeString, Optional: true, Description: "d", AtLeastOneOf: []string{"a", "b"}},
					"b": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
				}},
			},
			"exactly_one": {
				Type: pluginsdk.TypeList, Optional: true, Description: "d",
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"a": {Type: pluginsdk.TypeString, Optional: true, Description: "d", ExactlyOneOf: []string{"a", "b"}},
					"b": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
				}},
			},
			"opaque_child": {
				Type: pluginsdk.TypeList, Optional: true, Description: "d",
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"a": commonschema.Location(),
				}},
			},
	`))

	if !flagged(fs, "loose") {
		t.Error("expected SL006 on an all-optional block")
	}
	for _, p := range []string{"has_required", "at_least_one", "exactly_one", "opaque_child"} {
		if flagged(fs, p) {
			t.Errorf("did not expect SL006 on %q", p)
		}
	}
}
