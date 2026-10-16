// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerapps"
)

func TestValidateContainerAppRegistry(t *testing.T) {
	cases := []struct {
		Input Registry
		Valid bool
	}{
		{
			Input: Registry{
				Server:            "registry.example.com",
				UserName:          "user",
				PasswordSecretRef: "secretref",
			},
			Valid: true,
		},
		{
			Input: Registry{
				Server:   "registry.example.com",
				Identity: "identity",
			},
			Valid: true,
		},
		{
			Input: Registry{
				Server: "registry.example.com",
			},
			Valid: false,
		},
		{
			Input: Registry{
				Server:            "registry.example.com",
				UserName:          "user",
				PasswordSecretRef: "secretref",
				Identity:          "identity",
			},
			Valid: false,
		},
		{
			Input: Registry{
				Server:            "registry.example.com",
				PasswordSecretRef: "secretref",
			},
			Valid: false,
		},
		{
			Input: Registry{
				Server:   "registry.example.com",
				UserName: "user",
			},
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		err := ValidateContainerAppRegistry(tc.Input)
		valid := err == nil
		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %s", tc.Valid, valid, tc.Input)
		}
	}
}

func TestContainerEnvVarHashIdenticalEntries(t *testing.T) {
	a := map[string]interface{}{"name": "LOREM_VAR", "value": "ipsum"}
	b := map[string]interface{}{"name": "LOREM_VAR", "value": "ipsum"}
	if containerEnvVarHash(a) != containerEnvVarHash(b) {
		t.Fatal("expected equal hashes for identical entries")
	}
}

func TestContainerEnvVarHashDifferentValues(t *testing.T) {
	a := map[string]interface{}{"name": "LOREM_VAR", "value": "ipsum"}
	b := map[string]interface{}{"name": "LOREM_VAR", "value": "dolor"}
	if containerEnvVarHash(a) == containerEnvVarHash(b) {
		t.Fatal("expected different hashes for different values")
	}
}

func TestContainerEnvVarHashValueVsSecret(t *testing.T) {
	a := map[string]interface{}{"name": "LOREM_VAR", "value": "ipsum"}
	b := map[string]interface{}{"name": "LOREM_VAR", "secret_name": "sit-amet"}
	if containerEnvVarHash(a) == containerEnvVarHash(b) {
		t.Fatal("expected different hashes for value-backed vs secret-backed")
	}
}

func TestContainerEnvVarHashValueEqualsSecretName(t *testing.T) {
	a := map[string]interface{}{"name": "LOREM_VAR", "value": "consectetur", "secret_name": ""}
	b := map[string]interface{}{"name": "LOREM_VAR", "value": "", "secret_name": "consectetur"}
	if containerEnvVarHash(a) == containerEnvVarHash(b) {
		t.Fatal("expected different hashes when value equals secret_name")
	}
}

func TestContainerEnvVarHashDifferentSecretNames(t *testing.T) {
	a := map[string]interface{}{"name": "LOREM_VAR", "secret_name": "adipiscing"}
	b := map[string]interface{}{"name": "LOREM_VAR", "secret_name": "elit-sed"}
	if containerEnvVarHash(a) == containerEnvVarHash(b) {
		t.Fatal("expected different hashes for different secret_name values")
	}
}

func TestContainerEnvVarHashDifferentNames(t *testing.T) {
	a := map[string]interface{}{"name": "LOREM_VAR", "value": "ipsum"}
	b := map[string]interface{}{"name": "DOLOR_VAR", "value": "ipsum"}
	if containerEnvVarHash(a) == containerEnvVarHash(b) {
		t.Fatal("expected different hashes for different names")
	}
}

func TestContainerEnvVarHashCaseSensitive(t *testing.T) {
	a := map[string]interface{}{"name": "LOREM_VAR", "value": "ipsum"}
	b := map[string]interface{}{"name": "lorem_var", "value": "ipsum"}
	if containerEnvVarHash(a) == containerEnvVarHash(b) {
		t.Fatal("expected different hashes for different name casing")
	}
}

func TestContainerSecretHashSameNameDifferentValue(t *testing.T) {
	a := map[string]interface{}{"name": "lorem-secret", "value": "ipsum"}
	b := map[string]interface{}{"name": "lorem-secret", "value": "dolor"}
	if containerSecretHash(a) != containerSecretHash(b) {
		t.Fatal("expected equal hashes — secret identity is name-only")
	}
}

func TestContainerSecretHashDifferentNames(t *testing.T) {
	a := map[string]interface{}{"name": "lorem-secret"}
	b := map[string]interface{}{"name": "ipsum-secret"}
	if containerSecretHash(a) == containerSecretHash(b) {
		t.Fatal("expected different hashes for different secret names")
	}
}

func TestFlattenContainerAppSecretsNormalizesBlankKeyVaultURL(t *testing.T) {
	result := FlattenContainerAppSecrets(&containerapps.SecretsCollection{
		Value: []containerapps.ContainerAppSecret{
			{Name: pointer.To("lorem-secret"), KeyVaultURL: pointer.To(""), Value: pointer.To("ipsum")},
		},
	})
	if len(result) != 1 || result[0].KeyVaultSecretId != "" || result[0].Value != "ipsum" {
		t.Fatalf("expected value-backed secret with blank key_vault_secret_id, got %+v", result)
	}
}

func TestPreserveContainerAppSecretValues(t *testing.T) {
	current := []Secret{
		{Name: "lorem-secret", Value: ""},
		{Name: "dolor-secret", KeyVaultSecretId: "https://vault.example.net/secrets/amet", Value: ""},
	}
	prior := []Secret{
		{Name: "lorem-secret", Value: "consectetur"},
		{Name: "dolor-secret", KeyVaultSecretId: "https://vault.example.net/secrets/amet", Value: "ignored"},
	}

	result := PreserveContainerAppSecretValues(current, prior)

	if result[0].Value != "consectetur" {
		t.Fatalf("expected value-backed secret to be preserved, got %q", result[0].Value)
	}
	if result[1].Value != "" {
		t.Fatalf("expected key-vault-backed secret value to stay empty, got %q", result[1].Value)
	}
}
