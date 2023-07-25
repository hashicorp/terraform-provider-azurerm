// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestIotHubDeviceUpdateInstanceName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"Aa0",
		"012345678901234567890123",
	}
	for _, v := range validNames {
		_, errors := IotHubDeviceUpdateInstanceName(v, "example")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid IoT Hub Device Update Instance Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid_name",
		"double-hyphen--invalid",
		"aa",
		"0123456789012345678901234",
	}
	for _, v := range invalidNames {
		_, errors := IotHubDeviceUpdateInstanceName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid IoT Hub Device Update Instance Name", v)
		}
	}
}
