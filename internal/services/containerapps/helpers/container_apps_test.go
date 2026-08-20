// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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

// TestFlattenContainerAppSecrets_keyVaultValueIsEmptyString verifies that a
// Key Vault-backed secret is always flattened with value="" rather than a
// non-empty string. This ensures the flattened state is consistent with the
// schema Default:"" added to fix #32923 (null vs "" TypeSet hash mismatch).
func TestFlattenContainerAppSecrets_keyVaultValueIsEmptyString(t *testing.T) {
	input := &containerapps.SecretsCollection{
		Value: []containerapps.ContainerAppSecret{
			{
				Name:        pointer.To("kv-secret"),
				Identity:    pointer.To("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id"),
				KeyVaultURL: pointer.To("https://example.vault.azure.net/secrets/my-secret"),
				Value:       nil, // Azure never returns the value for KV-backed secrets
			},
		},
	}

	result := FlattenContainerAppSecrets(input)

	if len(result) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(result))
	}

	if result[0].Value != "" {
		t.Errorf("expected value to be empty string for KV-backed secret, got %q", result[0].Value)
	}
}

// TestFlattenContainerAppJobSecrets_keyVaultValueIsEmptyString is the same
// check for azurerm_container_app_job (uses jobs.JobSecretsCollection).
func TestFlattenContainerAppJobSecrets_keyVaultValueIsEmptyString(t *testing.T) {
	input := &jobs.JobSecretsCollection{
		Value: []jobs.Secret{
			{
				Name:        pointer.To("kv-secret"),
				Identity:    pointer.To("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id"),
				KeyVaultURL: pointer.To("https://example.vault.azure.net/secrets/my-secret"),
				Value:       nil,
			},
		},
	}

	result := FlattenContainerAppJobSecrets(input)

	if len(result) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(result))
	}

	if result[0].Value != "" {
		t.Errorf("expected value to be empty string for KV-backed secret, got %q", result[0].Value)
	}
}

// TestSecretsSchema_valueDefaultIsEmptyString verifies that the schema default
// for the value field is "" so that Terraform treats an omitted value the same
// way the flatten functions do, preventing the TypeSet null/empty hash mismatch
// described in #32923.
func TestSecretsSchema_valueDefaultIsEmptyString(t *testing.T) {
	s := SecretsSchema()
	resource, ok := s.Elem.(*pluginsdk.Resource)
	if !ok {
		t.Fatal("expected Elem to be *pluginsdk.Resource")
	}

	valueSchema, ok := resource.Schema["value"]
	if !ok {
		t.Fatal("expected 'value' key in secret schema")
	}

	if valueSchema.Default != "" {
		t.Errorf("expected Default to be empty string, got %v", valueSchema.Default)
	}
}
