package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SentinelDataConnectorAzureActiveDirectoryResource struct{}

func TestAccAzureRMSentinelDataConnectorAzureActiveDirectory_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_azure_active_directory", "test")
	r := SentinelDataConnectorAzureActiveDirectoryResource{}

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

func TestAccAzureRMSentinelDataConnectorAzureActiveDirectory_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_azure_active_directory", "test")
	r := SentinelDataConnectorAzureActiveDirectoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMSentinelDataConnectorAzureActiveDirectory_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_azure_active_directory", "test")
	r := SentinelDataConnectorAzureActiveDirectoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMSentinelDataConnectorAzureActiveDirectory_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_azure_active_directory", "test")
	r := SentinelDataConnectorAzureActiveDirectoryResource{}

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

func (r SentinelDataConnectorAzureActiveDirectoryResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Sentinel.DataConnectorsClient

	id, err := parse.DataConnectorID(state.ID)
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.WorkspaceName, id.Name); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Sentinel Data Connector Azure Active Directory %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r SentinelDataConnectorAzureActiveDirectoryResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_azure_active_directory" "test" {
  name                       = "accTestDC-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}
`, template, data.RandomInteger)
}

func (r SentinelDataConnectorAzureActiveDirectoryResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "test" {}

resource "azurerm_sentinel_data_connector_azure_active_directory" "test" {
  name                       = "accTestDC-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  tenant_id                  = data.azurerm_client_config.test.tenant_id
  alerts_enabled             = false
}
`, template, data.RandomInteger)
}

func (r SentinelDataConnectorAzureActiveDirectoryResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_azure_active_directory" "import" {
  name                       = azurerm_sentinel_data_connector_azure_active_directory.test.name
  log_analytics_workspace_id = azurerm_sentinel_data_connector_azure_active_directory.test.log_analytics_workspace_id
}
`, template)
}

func (r SentinelDataConnectorAzureActiveDirectoryResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
