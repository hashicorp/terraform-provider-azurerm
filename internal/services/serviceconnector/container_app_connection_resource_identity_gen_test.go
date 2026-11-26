// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package serviceconnector_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccContainerAppConnection_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_connection", "test")
	r := ContainerAppConnectionResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"container_app_name":  {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.storageBlob(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_container_app_connection.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_container_app_connection.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_container_app_connection.test", tfjsonpath.New("container_app_name"), tfjsonpath.New("container_app_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_container_app_connection.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("container_app_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_container_app_connection.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("container_app_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
