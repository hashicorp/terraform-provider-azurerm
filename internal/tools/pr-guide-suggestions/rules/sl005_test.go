// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "testing"

func TestSL005(t *testing.T) {
	fs := runRule(t, sl005, wrap(`
			"s":         {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
			"i":         {Type: pluginsdk.TypeInt, Required: true, Description: "d"},
			"f":         {Type: pluginsdk.TypeFloat, Optional: true, Description: "d"},
			"validated": {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.StringIsNotEmpty},
			"diag":      {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateDiagFunc: someDiag},
			"boolean":   {Type: pluginsdk.TypeBool, Optional: true, Description: "d"},
			"computed":  {Type: pluginsdk.TypeString, Computed: true},
	`))

	for _, p := range []string{"s", "i", "f"} {
		if !flagged(fs, p) {
			t.Errorf("expected SL005 on unvalidated %q", p)
		}
	}
	for _, p := range []string{"validated", "diag", "boolean", "computed"} {
		if flagged(fs, p) {
			t.Errorf("did not expect SL005 on %q", p)
		}
	}
}
