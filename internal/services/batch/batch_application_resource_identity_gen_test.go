// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccBatchApplication_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_application", "test")
	r := BatchApplicationResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"batch_account_name":  {},
		"name":                {},
		"resource_group_name": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basicForResourceIdentity(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_batch_application.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_batch_application.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_batch_application.test", tfjsonpath.New("batch_account_name"), tfjsonpath.New("account_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_batch_application.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_batch_application.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
