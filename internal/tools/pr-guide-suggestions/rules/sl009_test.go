// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"
)

func TestSL009(t *testing.T) {
	fs := runRule(t, sl009, wrap(`
			"size_mb":          {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"duration_seconds": {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"size_in_mb":       {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"timeout_ms":       commonschema.SomeHelper(),
			"count":            {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
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
	if flagged(fs, "count") {
		t.Error("did not expect SL009 on a unit-free name")
	}
}
