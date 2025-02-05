// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func PrivateKeyContentTest(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// not a private key
			Input: "not a private key",
			Valid: false,
		},

		{
			// valid
			Input: "-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABDiBJ8mdc\nO9BK5YPbSF7OeJAAAAGAAAAAEAAAAzAAAAC3NzaC1lZDI1NTE5AAAAICsAOqrYOIKq3/JQ\nKWWyWRawxxawwKL4jDTlYKOAI1Y/AAAAsGQC60Q7yg9dNjX5QL+JW4HieedOv7bNOvMXqh\n0UAyCqD1liZmvKRZaQQF4qt4jE2+SdBjwz+WsBpn3y7kxgrTy7xaHqx/V5l87n41qzpY/W\nXgUlPKzKDxqsJVk6tI+5jYhNx8s7Kc6Pd1QZQ09WaDwzw0Ag+7j6nLGxrRJ8+WeNX8Qm0f\nV1Ft8Bs0Bo8Rc7wRtQOgHmgBdZE6i5swaUmZM7jNV4NvOhJxnMkvygiiLM\n-----END OPENSSH PRIVATE KEY-----\n",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := PrivateSSHKey(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
