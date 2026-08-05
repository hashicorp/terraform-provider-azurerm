// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestKubernetesClusterKeyVaultKeyIDValidation(t *testing.T) {
	validateFunc := (containers.Registration{}).
		SupportedResources()["azurerm_kubernetes_cluster"].
		Schema["key_management_service"].
		Elem.(*pluginsdk.Resource).
		Schema["key_vault_key_id"].
		ValidateFunc

	testCases := map[string]struct {
		input string
		valid bool
	}{
		"versioned": {
			input: "https://example.vault.azure.net/keys/key-name/0123456789abcdef0123456789abcdef",
			valid: true,
		},
		"versionless": {
			input: "https://example.vault.azure.net/keys/key-name",
			valid: true,
		},
		"invalid": {
			input: "not-a-key-vault-key-id",
			valid: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			_, errors := validateFunc(testCase.input, "key_vault_key_id")
			if testCase.valid && len(errors) != 0 {
				t.Fatalf("expected %q to be valid, got errors: %v", testCase.input, errors)
			}
			if !testCase.valid && len(errors) == 0 {
				t.Fatalf("expected %q to be invalid", testCase.input)
			}
		})
	}
}
