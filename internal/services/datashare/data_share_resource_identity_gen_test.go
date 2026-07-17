// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datashare_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccDataShare_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share", "test")
	r := DataShareResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"account_name":        {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_data_share.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_data_share.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_share.test", tfjsonpath.New("account_name"), tfjsonpath.New("account_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_share.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("account_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_share.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("account_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
