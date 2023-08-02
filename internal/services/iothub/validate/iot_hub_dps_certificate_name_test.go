// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestIoTHubDpsCertificateName(t *testing.T) {
	validNames := []string{
		"validName123",
		"cert_a-1.cer",
		"1234567890123456789012345678901234567890123456789012345678901234",
	}
	for _, v := range validNames {
		_, errors := IoTHubDpsCertificateName(v, "example")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid IoT Hub DPS Certificate Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"",
		"invalid!",
		"!@£",
		"12345678901234567890123456789012345678901234567890123456789012345",
	}
	for _, v := range invalidNames {
		_, errors := IoTHubDpsCertificateName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid IoT Hub DPS Certificate Name", v)
		}
	}
}
