// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"
)

func TestSL004(t *testing.T) {
	fs := runRule(t, sl004, wrap(`
			"none_mode":  {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.StringInSlice([]string{"None", "Auto"}, false)},
			"off_mode":   {Type: pluginsdk.TypeString, Required: true, Description: "d", ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false)},
			"converted":  {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.StringInSlice([]string{string(pkg.ChannelNone)}, false)},
			"wrapped":    {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.Any(privatezones.ValidateID, validation.StringInSlice([]string{"System", "None"}, false))},
			"clean":      {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.StringInSlice([]string{"Standard", "Premium"}, false)},
			"computed":   {Type: pluginsdk.TypeString, Computed: true, ValidateFunc: validation.StringInSlice([]string{"None"}, false)},
			"enum_func":  {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.StringInSlice(possibleValues(), false)},
	`))

	// Literal "None", "Off", a string(...None) enum constant, and a StringInSlice
	// nested inside validation.Any are all detected.
	for _, p := range []string{"none_mode", "off_mode", "converted", "wrapped"} {
		if !flagged(fs, p) {
			t.Errorf("expected SL004 on %q", p)
		}
	}
	if f := at(fs, "none_mode"); f != nil && !strings.Contains(f.msg, "None") {
		t.Errorf("message should name the accepted value, got %q", f.msg)
	}
	// A clean enum, a computed-only property, and an enum sourced from an opaque
	// function are not flagged.
	for _, p := range []string{"clean", "computed", "enum_func"} {
		if flagged(fs, p) {
			t.Errorf("did not expect SL004 on %q", p)
		}
	}
}
