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

type SynapseIntegrationRuntimeSelfHostedResource struct{}

func TestAccSynapseIntegrationRuntimeSelfHosted_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_integration_runtime_self_hosted", "test")
	r := SynapseIntegrationRuntimeSelfHostedResource{}

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

func TestAccSynapseIntegrationRuntimeSelfHosted_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_integration_runtime_self_hosted", "test")
	r := SynapseIntegrationRuntimeSelfHostedResource{}

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

func (r SynapseIntegrationRuntimeSelfHostedResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IntegrationRuntimeID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.Synapse.IntegrationRuntimesClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}
	return utils.Bool(resp.ID != nil), nil
}

func (r SynapseIntegrationRuntimeSelfHostedResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "content"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestdf%d"
  location                             = azurerm_resource_group.test.location
  resource_group_name                  = azurerm_resource_group.test.name
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  managed_virtual_network_enabled      = true
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r SynapseIntegrationRuntimeSelfHostedResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_integration_runtime_self_hosted" "test" {
  name                 = "acctestSIR%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  description          = "test"
}
`, r.template(data), data.RandomInteger)
}

func (r SynapseIntegrationRuntimeSelfHostedResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_integration_runtime_self_hosted" "import" {
  name                 = azurerm_synapse_integration_runtime_self_hosted.test.name
  synapse_workspace_id = azurerm_synapse_integration_runtime_self_hosted.test.synapse_workspace_id
}
`, r.basic(data))
}
