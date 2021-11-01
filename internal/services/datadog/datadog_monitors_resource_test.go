package datadog_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DatadogMonitorResource struct{}

func TestAccDatadogMonitor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
	r := DatadogMonitorResource{}
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

func TestAccDatadogMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
	r := DatadogMonitorResource{}
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

func TestAccDatadogMonitor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
	r := DatadogMonitorResource{}
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

func TestAccDatadogMonitor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
	r := DatadogMonitorResource{}
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

func (r DatadogMonitorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DatadogMonitorID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Datadog.MonitorsClient.Get(ctx, id.ResourceGroup, id.MonitorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r DatadogMonitorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-datadog-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DatadogMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
	
	resource "azurerm_datadog_monitor" "test" {
		name = "test-terraform-6747642"
		resource_group_name = azurerm_resource_group.test.name
		location = azurerm_resource_group.test.location
		user_info {
			name          = "vidhi"
			email_address = "testtf@mpliftrelastic20210901outlo.onmicrosoft.com"
		}
		sku {
			name = "payg_v2_Monthly"
		}
		identity {
			type = "SystemAssigned"
		}
	}
`, r.template(data))
}

func (r DatadogMonitorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

	resource "azurerm_datadog_monitor" "import" {
	name                = azurerm_datadog_monitor.test.name
	resource_group_name = azurerm_datadog_monitor.test.resource_group_name
	location            = azurerm_datadog_monitor.test.location
	sku {
		name = "payg_v2_Monthly"
	}
	user_info {
		name          = "vidhi"
		email_address = "testtf@mpliftrelastic20210901outlo.onmicrosoft.com"
	}
	identity {
		type = "SystemAssigned"
	}
	}
`, r.template(data))
}

func (r DatadogMonitorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

	resource "azurerm_datadog_monitor" "test" {
	name                = "acctest-dm-%d"
	resource_group_name = azurerm_resource_group.test.name
	location            = azurerm_resource_group.test.location
	datadog_organization_properties {
		api_key           = ""
		application_key   = ""
		enterprise_app_id = ""
		linking_auth_code = ""
		linking_client_id = ""
		redirect_uri      = ""
	}
	identity {
		type = "SystemAssigned"
	}
	sku {
		name = "payg_v2_Monthly"
	}
	user_info {
		name          = "vidhi"
		email_address = "testtf@mpliftrelastic20210901outlo.onmicrosoft.com"
		phone_number  = ""
	}
	monitoring_status = false
	tags = {
		ENV = "Test"
	}
	}
`, r.template(data), data.RandomInteger)
}
