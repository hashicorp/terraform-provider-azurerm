package validate

import (
	"testing"
)

func TestNestedItemId(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/certificates/hello/world",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/keys/castle/1492",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		warnings, err := NestedItemId(tc.Input, "example")
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

func TestNestedItemIdWithOptionalVersion(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/certificates/hello/world",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/keys/castle/1492",
			ExpectError: false,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		warnings, err := NestedItemIdWithOptionalVersion(tc.Input, "example")
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
