package validate

import (
	"testing"
)

func TestValidateVaultURI(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "https://testkv.vault.azure.net/",
			ExpectError: false,
		},
		{
			Input:       "https://testkv.vault.azure.cn/",
			ExpectError: false,
		},
		{
			Input:       "https://testkv.vault.usgovcloudapi.net/",
			ExpectError: false,
		},
		{
			Input:       "https://testkv.vault.microsoftazure.de/",
			ExpectError: false,
		},
		{
			Input:       "https://abc",
			ExpectError: true,
		},
		{
			Input:       "https://^",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := VaultURI(tc.Input, "")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Key Vault Uri to trigger a validation error for '%s'", tc.Input)
		}
	}
}
