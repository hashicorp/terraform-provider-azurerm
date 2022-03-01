package elastic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/sdk/2020-07-01/monitorsresource"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ElasticMonitorResource struct{}

func TestAccElasticMonitor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_monitor", "test")
	r := ElasticMonitorResource{}
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

func TestAccElasticMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_monitor", "test")
	r := ElasticMonitorResource{}
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

func TestAccElasticMonitor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_monitor", "test")
	r := ElasticMonitorResource{}
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

func TestAccElasticMonitor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_monitor", "test")
	r := ElasticMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
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

func TestAccElasticMonitor_logs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_monitor", "test")
	r := ElasticMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// this proves that we don't need to destroy the `logs` block separately
			Config: r.logs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccElasticMonitor_logsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_monitor", "test")
	r := ElasticMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// create with it
			Config: r.logs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// update it
			Config: r.logsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// remove just the `logs` block
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ElasticMonitorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := monitorsresource.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Elastic.MonitorClient.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ElasticMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_monitor" "test" {
  name                        = "acctest-estc%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-monthly-consumption_Monthly"
  elastic_cloud_email_address = "acctestuser-%[1]d@hashicorptest.com"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ElasticMonitorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_elastic_monitor" "import" {
  name                        = azurerm_elastic_monitor.test.name
  resource_group_name         = azurerm_elastic_monitor.test.resource_group_name
  location                    = azurerm_elastic_monitor.test.location
  sku_name                    = azurerm_elastic_monitor.test.sku_name
  elastic_cloud_email_address = azurerm_elastic_monitor.test.elastic_cloud_email_address
}
`, r.basic(data))
}

func (r ElasticMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_monitor" "test" {
  name                        = "acctest-estc%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-monthly-consumption_Monthly"
  elastic_cloud_email_address = "acctestuser-%[1]d@hashicorptest.com"

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ElasticMonitorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_monitor" "test" {
  name                        = "acctest-estc%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-monthly-consumption_Monthly"
  elastic_cloud_email_address = "acctestuser-%[1]d@hashicorptest.com"
  monitoring_enabled          = false

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ElasticMonitorResource) logs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_monitor" "test" {
  name                        = "acctest-estc%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-monthly-consumption_Monthly"
  elastic_cloud_email_address = "acctestuser-%[1]d@hashicorptest.com"

  logs {
    filtering_tag {
       action = "Include"
       name   = "TerraformAccTest"
       value  = "RandomValue%[1]d"
    }

    # NOTE: these are intentionally not set to true here for testing purposes
    send_activity_logs     = false
    send_azuread_logs      = false
    send_subscription_logs = false
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ElasticMonitorResource) logsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_monitor" "test" {
  name                        = "acctest-estc%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-monthly-consumption_Monthly"
  elastic_cloud_email_address = "acctestuser-%[1]d@hashicorptest.com"

  logs {
    filtering_tag {
       action = "Include"
       name   = "TerraformAccTest"
       value  = "UpdatedValue-%[1]d"
    }

    # NOTE: these are intentionally not set to true here for testing purposes
    send_activity_logs     = false
    send_azuread_logs      = false
    send_subscription_logs = false
  }

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
