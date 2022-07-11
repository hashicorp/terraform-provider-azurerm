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

type IntegrationRuntimeAzureResource struct{}

func TestAccSynapseIntegrationRuntimeAzure_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_integration_runtime_azure", "test")
	r := IntegrationRuntimeAzureResource{}

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

func TestAccSynapseIntegrationRuntimeAzure_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_integration_runtime_azure", "test")
	r := IntegrationRuntimeAzureResource{}

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

func TestAccSynapseIntegrationRuntimeAzure_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_integration_runtime_azure", "test")
	r := IntegrationRuntimeAzureResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseIntegrationRuntimeAzure_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_integration_runtime_azure", "test")
	r := IntegrationRuntimeAzureResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseIntegrationRuntimeAzure_autoResolve(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_integration_runtime_azure", "test")
	r := IntegrationRuntimeAzureResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoResolve(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r IntegrationRuntimeAzureResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r IntegrationRuntimeAzureResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_integration_runtime_azure" "test" {
  name                 = "azure-integration-runtime"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  location             = azurerm_resource_group.test.location
}
`, r.template(data))
}

func (r IntegrationRuntimeAzureResource) autoResolve(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_integration_runtime_azure" "test" {
  name                 = "azure-integration-runtime"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  location             = "AutoResolve"
}
`, r.template(data))
}

func (r IntegrationRuntimeAzureResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_integration_runtime_azure" "import" {
  name                 = azurerm_synapse_integration_runtime_azure.test.name
  synapse_workspace_id = azurerm_synapse_integration_runtime_azure.test.synapse_workspace_id
  location             = azurerm_synapse_integration_runtime_azure.test.location
}
`, r.basic(data))
}

func (r IntegrationRuntimeAzureResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_integration_runtime_azure" "test" {
  name                 = "azure-integration-runtime"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  location             = azurerm_resource_group.test.location

  compute_type     = "ComputeOptimized"
  core_count       = 16
  time_to_live_min = 10
  description      = "test"
}
`, r.template(data))
}

func (IntegrationRuntimeAzureResource) template(data acceptance.TestData) string {
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
