// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"
)

func TestSL002(t *testing.T) {
	fs := runRule(t, sl002, wrap(`
			"single": {
				Type: pluginsdk.TypeList, Optional: true, MaxItems: 1, Description: "d",
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"origin": {Type: pluginsdk.TypeString, Required: true, Description: "d"},
				}},
			},
			"toggle": {
				Type: pluginsdk.TypeList, Optional: true, MaxItems: 1, Description: "d",
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"enabled": {Type: pluginsdk.TypeBool, Optional: true, Description: "d"},
				}},
			},
			"multi": {
				Type: pluginsdk.TypeList, Optional: true, MaxItems: 1, Description: "d",
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"a": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
					"b": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
				}},
			},
			"unbounded": {
				Type: pluginsdk.TypeList, Optional: true, Description: "d",
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"a": {Type: pluginsdk.TypeString, Required: true, Description: "d"},
				}},
			},
	`))

	single := at(fs, "single")
	if single == nil {
		t.Fatal("expected SL002 on single-property block")
	}
	if !strings.Contains(single.fix, `"single_origin"`) {
		t.Errorf("expected flatten fix naming single_origin, got %q", single.fix)
	}

	toggle := at(fs, "toggle")
	if toggle == nil {
		t.Fatal("expected SL002 on the enabled toggle block")
	}
	if !strings.Contains(toggle.fix, `boolean "toggle_enabled"`) {
		t.Errorf("expected boolean toggle fix, got %q", toggle.fix)
	}

	if flagged(fs, "multi") {
		t.Error("did not expect SL002 on a multi-property block")
	}
	if flagged(fs, "unbounded") {
		t.Error("did not expect SL002 without MaxItems:1")
	}
}
