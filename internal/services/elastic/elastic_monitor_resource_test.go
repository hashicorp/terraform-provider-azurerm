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

func (r ElasticMonitorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctest-elastic-%d"
  location = "%s"
}
`, data.RandomInteger%1000, data.Locations.Primary)
}

func (r ElasticMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_elastic_monitor" "test" {
  name                = "test-tf-elastic-basic-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "staging_Monthly"

  user_info {
    email_address = "ElasticTerraformTesting@mpliftrelastic20211117outlo.onmicrosoft.com"
  }
}
`, r.template(data), data.RandomInteger%1000)
}

func (r ElasticMonitorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
resource "azurerm_elastic_monitor" "testImport" {
  name                = azurerm_elastic_monitor.test.name
  resource_group_name = azurerm_elastic_monitor.test.resource_group_name
  location            = azurerm_elastic_monitor.test.location
  sku_name            = "staging_Monthly"

  user_info {
    email_address = "ElasticTerraformTesting@mpliftrelastic20211117outlo.onmicrosoft.com"
  }
}
`, r.basic(data))
}

func (r ElasticMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
resource "azurerm_elastic_monitor" "test" {
  name                = "test-tf-elastic-basic-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "staging_Monthly"

  user_info {
    email_address = "ElasticTerraformTesting@mpliftrelastic20211117outlo.onmicrosoft.com"
  }
  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger%1000)
}

func (r ElasticMonitorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
resource "azurerm_elastic_monitor" "test" {
  name                = "test-tf-elastic-complete-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "staging_Monthly"
  monitoring_status   = false

  user_info {
    email_address = "ElasticTerraformTesting@mpliftrelastic20211117outlo.onmicrosoft.com"
  }
  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger%1000)
}
