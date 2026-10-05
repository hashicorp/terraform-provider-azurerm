// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package windows_virtual_machine_test

func allocationType(static bool) string {
	if static {
		return "Static"
	}

	return "Dynamic"
}
