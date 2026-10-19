// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import "testing"

func TestNodePoolSecurityFlagsChanged(t *testing.T) {
	disabled := []any{map[string]any{"secure_boot_enabled": false, "vtpm_enabled": false}}
	enabled := []any{map[string]any{"secure_boot_enabled": true, "vtpm_enabled": true}}

	for _, tc := range []struct {
		name string
		old  []any
		new  []any
		want bool
	}{
		{name: "omitted to disabled", old: []any{}, new: disabled},
		{name: "disabled to omitted", old: disabled, new: []any{}},
		{name: "omitted to enabled", old: []any{}, new: enabled, want: true},
		{name: "enabled to omitted", old: enabled, new: []any{}, want: true},
		{name: "enabled to disabled", old: enabled, new: disabled, want: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := nodePoolSecurityFlagsChanged(tc.old, tc.new); got != tc.want {
				t.Fatalf("nodePoolSecurityFlagsChanged() = %t, want %t", got, tc.want)
			}
		})
	}
}
