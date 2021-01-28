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

type SentinelDataConnectorThreatIntelligenceTaxiiResource struct{}

func TestAccAzureRMSentinelDataConnectorThreatIntelligenceTaxii_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_threat_intelligence_taxii", "test")
	r := SentinelDataConnectorThreatIntelligenceTaxiiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("taxii_server_username", "taxii_server_password"),
	})
}

func TestAccAzureRMSentinelDataConnectorThreatIntelligenceTaxii_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_threat_intelligence_taxii", "test")
	r := SentinelDataConnectorThreatIntelligenceTaxiiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.update(data, "foo", 107),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("taxii_server_username", "taxii_server_password"),
		{
			Config: r.update(data, "bar", 135),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("taxii_server_username", "taxii_server_password"),
	})
}

func TestAccAzureRMSentinelDataConnectorThreatIntelligenceTaxii_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_threat_intelligence_taxii", "test")
	r := SentinelDataConnectorThreatIntelligenceTaxiiResource{}

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

func (r SentinelDataConnectorThreatIntelligenceTaxiiResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Sentinel.DataConnectorsClient

	id, err := parse.DataConnectorID(state.ID)
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.WorkspaceName, id.Name); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Sentinel Data Connector Threat Intelligence Taxii %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r SentinelDataConnectorThreatIntelligenceTaxiiResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_threat_intelligence_taxii" "test" {
  name                       = "accTestDC-%[2]d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  display_name               = "accTestDC-%[2]d"
  taxii_server_api_root      = "https://limo.anomali.com/api/v1/taxii2/feeds"
  taxii_server_collection_id = 107
  taxii_server_username      = "guest"
  taxii_server_password      = "guest"
}
`, template, data.RandomInteger)
}

func (r SentinelDataConnectorThreatIntelligenceTaxiiResource) update(data acceptance.TestData, displayName string, collectionId int) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "test" {}

resource "azurerm_sentinel_data_connector_threat_intelligence_taxii" "test" {
  name                       = "accTestDC-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  display_name               = "%s"
  taxii_server_api_root      = "https://limo.anomali.com/api/v1/taxii2/feeds"
  taxii_server_collection_id = %d
  taxii_server_username      = "guest"
  taxii_server_password      = "guest"
  tenant_id                  = data.azurerm_client_config.test.tenant_id
}
`, template, data.RandomInteger, displayName, collectionId)
}

func (r SentinelDataConnectorThreatIntelligenceTaxiiResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_threat_intelligence_taxii" "import" {
  name                       = azurerm_sentinel_data_connector_threat_intelligence_taxii.test.name
  log_analytics_workspace_id = azurerm_sentinel_data_connector_threat_intelligence_taxii.test.log_analytics_workspace_id
  display_name               = azurerm_sentinel_data_connector_threat_intelligence_taxii.test.display_name
  taxii_server_api_root      = azurerm_sentinel_data_connector_threat_intelligence_taxii.test.taxii_server_api_root
  taxii_server_collection_id = azurerm_sentinel_data_connector_threat_intelligence_taxii.test.taxii_server_collection_id
}
`, template)
}

func (r SentinelDataConnectorThreatIntelligenceTaxiiResource) template(data acceptance.TestData) string {
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
