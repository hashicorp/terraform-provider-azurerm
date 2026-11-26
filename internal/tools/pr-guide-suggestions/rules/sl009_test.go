// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"
)

func TestSL009(t *testing.T) {
	fs := runRule(t, sl009, wrap(`
			"size_mb":                  {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"duration_seconds":         {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"size_in_mb":               {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"data_storage_size_in_gbs": {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"data_storage_size_in_gb":  {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"timeout_in_second":        {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"timeout_in_seconds":       {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"timeout_ms":               commonschema.SomeHelper(),
			"count":                    {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
	`))

	if f := at(fs, "size_mb"); f == nil {
		t.Error("expected SL009 on size_mb")
	} else if !strings.Contains(f.fix, `"size_in_mb"`) {
		t.Errorf("expected rename to size_in_mb, got %q", f.fix)
	}
	if !flagged(fs, "duration_seconds") {
		t.Error("expected SL009 on duration_seconds")
	}
	// Name-only rule: applies even to helper-valued properties.
	if !flagged(fs, "timeout_ms") {
		t.Error("expected SL009 on helper-valued timeout_ms")
	}
	if flagged(fs, "size_in_mb") {
		t.Error("did not expect SL009 on the already-correct size_in_mb")
	}
	if f := at(fs, "data_storage_size_in_gbs"); f == nil {
		t.Error("expected SL009 on the plural data size suffix in data_storage_size_in_gbs")
	} else if f.rename != "data_storage_size_in_gb" {
		t.Errorf("expected auto-fix rename to data_storage_size_in_gb, got %q", f.rename)
	}
	if flagged(fs, "data_storage_size_in_gb") {
		t.Error("did not expect SL009 on the singular data size suffix in data_storage_size_in_gb")
	}
	if f := at(fs, "timeout_in_second"); f == nil {
		t.Error("expected SL009 on the singular time suffix in timeout_in_second")
	} else if f.rename != "timeout_in_seconds" {
		t.Errorf("expected auto-fix rename to timeout_in_seconds, got %q", f.rename)
	}
	if flagged(fs, "timeout_in_seconds") {
		t.Error("did not expect SL009 on the plural time suffix in timeout_in_seconds")
	}
	if flagged(fs, "count") {
		t.Error("did not expect SL009 on a unit-free name")
	}
}
