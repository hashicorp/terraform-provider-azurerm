// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover"
)

func TestStorageMoverSmbMountEndpointCredentialValidation(t *testing.T) {
	arguments := storagemover.StorageMoverSmbMountEndpointResource{}.Arguments()
	for _, test := range []struct {
		name      string
		field     string
		value     string
		wantError bool
	}{
		{
			name:  "versionless-username",
			field: "username_key_vault_secret_id",
			value: "https://example.vault.azure.net/secrets/username",
		},
		{
			name:      "versioned-username",
			field:     "username_key_vault_secret_id",
			value:     "https://example.vault.azure.net/secrets/username/00000000000000000000000000000001",
			wantError: true,
		},
		{
			name:      "username-key-not-secret",
			field:     "username_key_vault_secret_id",
			value:     "https://example.vault.azure.net/keys/username",
			wantError: true,
		},
		{
			name:  "versionless-password",
			field: "password_key_vault_secret_id",
			value: "https://example.vault.azure.net/secrets/password",
		},
		{
			name:  "versioned-password",
			field: "password_key_vault_secret_id",
			value: "https://example.vault.azure.net/secrets/password/00000000000000000000000000000001",
		},
		{
			name:      "password-key-not-secret",
			field:     "password_key_vault_secret_id",
			value:     "https://example.vault.azure.net/keys/password",
			wantError: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, errors := arguments[test.field].ValidateFunc(test.value, test.field)
			if (len(errors) != 0) != test.wantError {
				t.Fatalf("expected validation error=%t, got %v", test.wantError, errors)
			}
		})
	}
}
