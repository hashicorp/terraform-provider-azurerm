// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestIotHubDeviceUpdateAccountName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"Aa0",
		"012345678901234567890123",
	}
	for _, v := range validNames {
		_, errors := IotHubDeviceUpdateAccountName(v, "example")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid IoT Hub Device Update Account Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid_name",
		"double-hyphen--invalid",
		"aa",
		"0123456789012345678901234",
	}
	for _, v := range invalidNames {
		_, errors := IotHubDeviceUpdateAccountName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid IoT Hub Device Update Account Name", v)
		}
	}
}
