package azure

import (
	"testing"
)

func TestAccAzureRMValidateKeyVaultChildID(t *testing.T) {
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
		warnings, err := ValidateKeyVaultChildId(tc.Input, "example")
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

func TestAccAzureRMValidateKeyVaultChildIDVersionOptional(t *testing.T) {
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
		warnings, err := ValidateKeyVaultChildIdVersionOptional(tc.Input, "example")
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

func TestAccAzureRMKeyVaultChild_parseID(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    KeyVaultChildID
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
			Expected: KeyVaultChildID{
				Name:            "bird",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "fdf067c93bbb4b22bff4d8b7a9a56217",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/certificates/hello/world",
			ExpectError: false,
			Expected: KeyVaultChildID{
				Name:            "hello",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "world",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/keys/castle/1492",
			ExpectError: false,
			Expected: KeyVaultChildID{
				Name:            "castle",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "1492",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		secretId, err := ParseKeyVaultChildID(tc.Input)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
			}

			return
		}

		if secretId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.KeyVaultBaseUrl != secretId.KeyVaultBaseUrl {
			t.Fatalf("Expected 'KeyVaultBaseUrl' to be '%s', got '%s' for ID '%s'", tc.Expected.KeyVaultBaseUrl, secretId.KeyVaultBaseUrl, tc.Input)
		}

		if tc.Expected.Name != secretId.Name {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Name, secretId.Name, tc.Input)
		}

		if tc.Expected.Version != secretId.Version {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Version, secretId.Version, tc.Input)
		}
	}
}

func TestAccAzureRMKeyVaultChild_parseIDVersionOptional(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    KeyVaultChildID
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
			Expected: KeyVaultChildID{
				Name:            "bird",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
			Expected: KeyVaultChildID{
				Name:            "bird",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "fdf067c93bbb4b22bff4d8b7a9a56217",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/certificates/hello/world",
			ExpectError: false,
			Expected: KeyVaultChildID{
				Name:            "hello",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "world",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/keys/castle/1492",
			ExpectError: false,
			Expected: KeyVaultChildID{
				Name:            "castle",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "1492",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		secretId, err := ParseKeyVaultChildIDVersionOptional(tc.Input)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
			}

			return
		}

		if secretId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.KeyVaultBaseUrl != secretId.KeyVaultBaseUrl {
			t.Fatalf("Expected 'KeyVaultBaseUrl' to be '%s', got '%s' for ID '%s'", tc.Expected.KeyVaultBaseUrl, secretId.KeyVaultBaseUrl, tc.Input)
		}

		if tc.Expected.Name != secretId.Name {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Name, secretId.Name, tc.Input)
		}

		if tc.Expected.Version != secretId.Version {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Version, secretId.Version, tc.Input)
		}
	}
}

func TestAccAzureRMKeyVaultChild_validateName(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "hello",
			ExpectError: false,
		},
		{
			Input:       "hello-world",
			ExpectError: false,
		},
		{
			Input:       "hello-world-21",
			ExpectError: false,
		},
		{
			Input:       "hello_world_21",
			ExpectError: true,
		},
		{
			Input:       "Hello-World",
			ExpectError: false,
		},
		{
			Input:       "20202020",
			ExpectError: false,
		},
		{
			Input:       "ABC123!@£",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := ValidateKeyVaultChildName(tc.Input, "")

		hasError := len(errors) > 0

		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Key Vault Child Name to trigger a validation error for '%s'", tc.Input)
		}
	}
}

func TestNewKeyVaultChildResourceID(t *testing.T) {
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
		id, err := NewKeyVaultChildResourceID(tc.keyVaultBaseUrl, childType, childName, childVersion)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for New Resource ID '%s': %+v", tc.keyVaultBaseUrl, err)
				return
			}
		}
		if id != tc.Expected {
			t.Fatalf("Expected id for %q to be %q, got %q", tc.keyVaultBaseUrl, tc.Expected, id)
		}
	}
}
