// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SynapseRoleAssignmentResource struct{}

func TestAccSynapseRoleAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_role_assignment", "test")
	r := SynapseRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseRoleAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_role_assignment", "test")
	r := SynapseRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSynapseRoleAssignment_sparkPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_role_assignment", "test")
	r := SynapseRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sparkPool(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SynapseRoleAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.RoleAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	workspaceName, _, err := parse.SynapseScope(id.Scope)
	if err != nil {
		return nil, err
	}

	env := client.Account.Environment
	suffix, ok := env.Synapse.DomainSuffix()
	if !ok {
		return nil, fmt.Errorf("could not determine the domain suffix for synapse in environment %q: %+v", env.Name, env.Storage)
	}

	roleAssignmentsClient, err := client.Synapse.RoleAssignmentsClient(workspaceName, *suffix)
	if err != nil {
		return nil, err
	}
	resp, err := roleAssignmentsClient.GetRoleAssignmentByID(ctx, id.DataPlaneAssignmentId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Synapse RoleAssignment (Resource Group %q): %+v", workspaceName, err)
	}

	return utils.Bool(true), nil
}

func (r SynapseRoleAssignmentResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_role_assignment" "test" {
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  role_name            = "Synapse SQL Administrator"
  principal_id         = data.azurerm_client_config.current.object_id

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template)
}

func (r SynapseRoleAssignmentResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_role_assignment" "import" {
  synapse_workspace_id = azurerm_synapse_role_assignment.test.synapse_workspace_id
  role_name            = azurerm_synapse_role_assignment.test.role_name
  principal_id         = azurerm_synapse_role_assignment.test.principal_id
}
`, config)
}

func (r SynapseRoleAssignmentResource) sparkPool(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_spark_pool" "test" {
  name                 = "acctestSSP%s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"
  node_count           = 3
  spark_version        = "3.4"
}

resource "azurerm_synapse_role_assignment" "test" {
  synapse_spark_pool_id = azurerm_synapse_spark_pool.test.id
  role_name             = "Synapse Contributor"
  principal_id          = data.azurerm_client_config.current.object_id

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomString)
}

func (r SynapseRoleAssignmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_synapse_firewall_rule" "test" {
  name                 = "AllowAll"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}

data "azurerm_client_config" "current" {}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}
