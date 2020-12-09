package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SynapseManagedPrivateEndpointResource struct{}

func TestAccSynapseManagedPrivateEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseManagedPrivateEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_managed_private_endpoint", "test")
	r := SynapseManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r SynapseManagedPrivateEndpointResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ManagedPrivateEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	environment := client.Account.Environment
	managedPrivateEndpointsClient, err := client.Synapse.ManagedPrivateEndpointsClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return nil, err
	}
	resp, err := managedPrivateEndpointsClient.Get(ctx, id.ManagedVirtualNetworkName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Synapse Managed Private Endpoints (Workspace %q / Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r SynapseManagedPrivateEndpointResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_synapse_managed_private_endpoint" "test" {
  name                 = "acctestEndpoint%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  target_resource_id   = azurerm_storage_account.test_endpoint.id
  subresource_name     = "blob"

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r SynapseManagedPrivateEndpointResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_synapse_managed_private_endpoint" "import" {
  name                 = azurerm_synapse_managed_private_endpoint.test.name
  synapse_workspace_id = azurerm_synapse_managed_private_endpoint.test.synapse_workspace_id
  target_resource_id   = azurerm_synapse_managed_private_endpoint.test.target_resource_id
  subresource_name     = azurerm_synapse_managed_private_endpoint.test.subresource_name
}
`, config)
}

func (r SynapseManagedPrivateEndpointResource) template(data acceptance.TestData) string {
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

resource "azurerm_storage_account" "test_endpoint" {
  name                     = "acctestacce%s"
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
  managed_virtual_network_enabled      = true
}

resource "azurerm_synapse_firewall_rule" "test" {
  name                 = "AllowAll"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
}
