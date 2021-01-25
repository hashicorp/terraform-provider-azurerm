package parse

import "testing"

func TestNewNestedItemID(t *testing.T) {
	childType := "keys"
	childName := "test"
	childVersion := "testVersionString"
	cases := []struct {
		Scenario        string
		keyVaultBaseUrl string
		Expected        string
		ExpectError     bool
	}{
		{
			Scenario:        "empty values",
			keyVaultBaseUrl: "",
			Expected:        "",
			ExpectError:     true,
		},
		{
			Scenario:        "valid, no port",
			keyVaultBaseUrl: "https://test.vault.azure.net",
			Expected:        "https://test.vault.azure.net/keys/test/testVersionString",
			ExpectError:     false,
		},
		{
			Scenario:        "valid, with port",
			keyVaultBaseUrl: "https://test.vault.azure.net:443",
			Expected:        "https://test.vault.azure.net/keys/test/testVersionString",
			ExpectError:     false,
		},
	}
	for _, tc := range cases {
		id, err := NewNestedItemID(tc.keyVaultBaseUrl, childType, childName, childVersion)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for New Resource ID '%s': %+v", tc.keyVaultBaseUrl, err)
				return
			}
		}
		if id.ID() != tc.Expected {
			t.Fatalf("Expected id for %q to be %q, got %q", tc.keyVaultBaseUrl, tc.Expected, id)
		}
	}
}
