package elastic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/parse"
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
	id, err := parse.ElasticMonitorID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Elastic.MonitorClient.Get(ctx, id.ResourceGroup, id.MonitorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Elastic Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
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
		name = "test-tf-elastic-basic-%d"
		resource_group_name = azurerm_resource_group.test.name
		location = azurerm_resource_group.test.location
		user_info {
			email_address = "utkarshjain@microsoft.com"
		}
		sku {
			name = "staging_Monthly"
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
	sku {
		name = "staging_Monthly"
	}
	user_info {
		email_address = "utkarshjain@microsoft.com"
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
	sku {
		name = "staging_Monthly"
	}
	user_info {
		email_address = "utkarshjain@microsoft.com"
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
	sku {
		name = "staging_Monthly"
	}
	user_info {
		email_address = "utkarshjain@microsoft.com"
	}
	monitoring_status = false
	tags = {
		ENV = "Test"
	}
	}
`, r.template(data), data.RandomInteger%1000)
}
