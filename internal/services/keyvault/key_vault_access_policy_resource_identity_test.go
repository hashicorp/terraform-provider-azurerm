// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

// NOTE: Resource Identity is implemented manually for `azurerm_key_vault_access_policy` because the resource
// uses a custom parse.AccessPolicyId type that does not implement resourceids.ResourceId.
// This means the go:generate tool cannot be used and identity-based import is not supported.

func TestAccKeyVaultAccessPolicy_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyResource{}

	checkedFields := map[string]struct{}{
		"key_vault_id":   {},
		"object_id":      {},
		"application_id": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_key_vault_access_policy.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_key_vault_access_policy.test", tfjsonpath.New("key_vault_id"), tfjsonpath.New("key_vault_id")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_key_vault_access_policy.test", tfjsonpath.New("object_id"), tfjsonpath.New("object_id")),
				statecheck.ExpectIdentityValue("azurerm_key_vault_access_policy.test", tfjsonpath.New("application_id"), knownvalue.StringExact("")),
			},
		},
		data.ImportBlockWithIDStep(false),
	}, false)
}
