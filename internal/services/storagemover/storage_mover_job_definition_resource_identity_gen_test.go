// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccStorageMoverJobDefinition_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_job_definition", "test")
	r := StorageMoverJobDefinitionResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"project_name":        {},
		"resource_group_name": {},
		"storage_mover_name":  {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_storage_mover_job_definition.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_storage_mover_job_definition.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_mover_job_definition.test", tfjsonpath.New("project_name"), tfjsonpath.New("storage_mover_project_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_mover_job_definition.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("storage_mover_project_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_mover_job_definition.test", tfjsonpath.New("storage_mover_name"), tfjsonpath.New("storage_mover_project_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_mover_job_definition.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("storage_mover_project_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
