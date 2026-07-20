// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"
)

func TestSL011(t *testing.T) {
	fs := runRule(t, sl011, wrap(`
			"is_enabled":  {Type: pluginsdk.TypeBool, Optional: true, Description: "d"},
			"is_name":     {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.StringIsNotEmpty},
			"enabled":     {Type: pluginsdk.TypeBool, Optional: true, Description: "d"},
			"is_helper":   commonschema.SomeBool(),
	`))

	if f := at(fs, "is_enabled"); f == nil {
		t.Error("expected SL011 on is_enabled")
	} else if !strings.Contains(f.fix, `"enabled"`) {
		t.Errorf("expected rename to enabled, got %q", f.fix)
	}
	// Only booleans; a string named is_* is not flagged.
	if flagged(fs, "is_name") {
		t.Error("did not expect SL011 on a non-boolean is_* property")
	}
	if flagged(fs, "enabled") {
		t.Error("did not expect SL011 without the is_ prefix")
	}
	// Opaque helper value: type unknown, so not flagged.
	if flagged(fs, "is_helper") {
		t.Error("did not expect SL011 on a helper-valued property")
	}
}
