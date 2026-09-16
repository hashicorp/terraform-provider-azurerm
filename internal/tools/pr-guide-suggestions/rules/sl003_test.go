// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"
)

func TestSL003(t *testing.T) {
	fs := runRule(t, sl003, wrap(`
			"max_scalar": {Type: pluginsdk.TypeString, Optional: true, Description: "d", MaxItems: 2},
			"min_scalar": {Type: pluginsdk.TypeInt, Optional: true, Description: "d", MinItems: 1},
			"plain":      {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
			"list":       {Type: pluginsdk.TypeList, Optional: true, Description: "d", MaxItems: 2, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString}},
			"set":        {Type: pluginsdk.TypeSet, Optional: true, Description: "d", MinItems: 1, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString}},
	`))

	if !flagged(fs, "max_scalar") || !flagged(fs, "min_scalar") {
		t.Error("expected SL003 on scalars declaring MaxItems/MinItems")
	}
	if f := at(fs, "max_scalar"); f != nil && !strings.Contains(f.fix, "remove MinItems/MaxItems") {
		t.Errorf("expected removal fix, got %q", f.fix)
	}
	for _, p := range []string{"plain", "list", "set"} {
		if flagged(fs, p) {
			t.Errorf("did not expect SL003 on %q", p)
		}
	}
}

func TestSL003Severity(t *testing.T) {
	if sl003.Severity != Error {
		t.Errorf("SL003 should be error severity, got %q", sl003.Severity)
	}
}
