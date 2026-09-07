// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"
)

func TestSL010(t *testing.T) {
	fs := runRule(t, sl010, wrap(`
			"vm_count":     {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"max_size":     {Type: pluginsdk.TypeInt, Optional: true, Description: "d", ValidateFunc: validation.IntAtLeast(1)},
			"database_name": {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.StringIsNotEmpty},
	`))

	if f := at(fs, "vm_count"); f == nil {
		t.Error("expected SL010 on vm_count")
	} else if !strings.Contains(f.fix, `"virtual_machine_count"`) {
		t.Errorf("expected expansion to virtual_machine_count, got %q", f.fix)
	} else if f.rename != "virtual_machine_count" {
		t.Errorf("expected auto-fix rename to virtual_machine_count, got %q", f.rename)
	}
	if !flagged(fs, "max_size") {
		t.Error("expected SL010 on max_size")
	}
	if flagged(fs, "database_name") {
		t.Error("did not expect SL010 on a fully-spelled name")
	}
}
