// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managedidentity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func testAccFederatedIdentityCredential_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_federated_identity_credential", "test")
	r := FederatedIdentityCredentialResource{}

	checkedFields := map[string]struct{}{
		"name":                        {},
		"resource_group_name":         {},
		"subscription_id":             {},
		"user_assigned_identity_name": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_federated_identity_credential.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_federated_identity_credential.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_federated_identity_credential.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("user_assigned_identity_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_federated_identity_credential.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("user_assigned_identity_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_federated_identity_credential.test", tfjsonpath.New("user_assigned_identity_name"), tfjsonpath.New("user_assigned_identity_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, true)
}
