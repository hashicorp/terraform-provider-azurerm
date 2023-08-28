// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestCertificateContactsID(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "https://my-keyvault.vault.azure.net/certificates/contacts",
			ExpectError: false,
		},
		{
			// empty
			Input:       "",
			ExpectError: true,
		},
		{
			// missing suffix
			Input:       "https://my-keyvault.vault.azure.net",
			ExpectError: true,
		},
		{
			// wrong suffix
			Input:       "https://my-keyvault.vault.azure.net/certificates/contact",
			ExpectError: true,
		},
		{
			// additional item in path
			Input:       "https://my-keyvault.vault.azure.net/x/certificates/contact",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		warnings, err := CertificateContactsID(tc.Input, "example")
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, err)
			}

			return
		}

		if tc.ExpectError && len(warnings) == 0 {
			t.Fatalf("Got no errors for input %q but expected some", tc.Input)
		} else if !tc.ExpectError && len(warnings) > 0 {
			t.Fatalf("Got %d errors for input %q when didn't expect any", len(warnings), tc.Input)
		}
	}
}
