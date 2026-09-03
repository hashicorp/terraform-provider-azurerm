// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccSynapseWorkspace_V0ToV1_501 tests the state migration path from an `id` with lowercased
// static segments to their canonicalized format. It uses v5.0.1 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccSynapseWorkspace_V0ToV1_501(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace", "test")
	r := SynapseWorkspaceResource{}

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-synapse-%[2]d/providers/Microsoft.Synapse/workspaces/acctestsw%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicV0(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Ensure import is a no-op to prevent the resource from normalizing the ID due to a combined CreateUpdate func
					plancheck.ExpectResourceAction(importedResourceName, plancheck.ResourceActionNoop),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acctestRG-synapse-%[2]d/providers/microsoft.synapse/workspaces/acctestsw%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicImported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-synapse-%[2]d/providers/Microsoft.Synapse/workspaces/acctestsw%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
	}, "5.0.1")
}

func (r SynapseWorkspaceResource) v0MigrationTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuaid%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r SynapseWorkspaceResource) basicV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_workspace" "test-import" {
  name                                 = "acctestsw%[2]d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  lifecycle {
    # ignore sql_administrator_login_password, it's not returned on imports
    ignore_changes = [sql_administrator_login_password]
  }
}

removed {
  from = azurerm_synapse_workspace.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_synapse_workspace.test-import
  id = "/subscriptions/%[3]s/resourcegroups/acctestRG-synapse-%[2]d/providers/microsoft.synapse/workspaces/acctestsw%[2]d"
}
`, r.v0MigrationTemplate(data), data.RandomInteger, data.Subscriptions.Primary)
}

func (r SynapseWorkspaceResource) basicImported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_workspace" "test-import" {
  name                                 = "acctestsw%[2]d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.v0MigrationTemplate(data), data.RandomInteger)
}
