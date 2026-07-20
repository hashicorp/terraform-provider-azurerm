// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"
)

func TestSL013(t *testing.T) {
	fs := runRule(t, sl013, wrap(`
			"none_id":     {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
			"weak_id":     {Type: pluginsdk.TypeString, Required: true, Description: "d", ValidateFunc: validation.StringIsNotEmpty},
			"strong_id":   {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: commonids.ValidateSubscriptionID},
			"uuid_id":     {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.IsUUID},
			"computed_id": {Type: pluginsdk.TypeString, Computed: true},
			"helper_id":   commonschema.ResourceIDReferenceOptional(),
			"name":        {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
	`))

	if f := at(fs, "none_id"); f == nil {
		t.Error("expected SL013 on unvalidated *_id")
	} else if !strings.Contains(f.msg, "no validation") {
		t.Errorf("expected 'no validation' message, got %q", f.msg)
	}
	if f := at(fs, "weak_id"); f == nil {
		t.Error("expected SL013 on weakly-validated *_id")
	} else if !strings.Contains(f.msg, "weak validation") {
		t.Errorf("expected 'weak validation' message, got %q", f.msg)
	}
	for _, p := range []string{"strong_id", "uuid_id", "computed_id", "helper_id", "name"} {
		if flagged(fs, p) {
			t.Errorf("did not expect SL013 on %q", p)
		}
	}
}
