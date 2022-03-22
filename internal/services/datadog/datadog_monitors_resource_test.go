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
		data.ImportStep("user_info",
			"user_info.0.name",
			"user_info.0.email_address",
			"datadog_organization_properties",
			"datadog_organization_properties.0",
			"datadog_organization_properties.0.id",
			"datadog_organization_properties.0.name",
			"datadog_organization_properties.0.api_key",
			"datadog_organization_properties.0.application_key",
			"datadog_organization_properties.0.enterprise_app_id",
			"datadog_organization_properties.0.linking_auth_code",
			"datadog_organization_properties.0.linking_client_id",
			"datadog_organization_properties.0.redirect_uri"),
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
		data.ImportStep("user_info",
			"user_info.0.name",
			"user_info.0.email_address",
			"datadog_organization_properties",
			"datadog_organization_properties.0",
			"datadog_organization_properties.0.id",
			"datadog_organization_properties.0.name",
			"datadog_organization_properties.0.api_key",
			"datadog_organization_properties.0.application_key",
			"datadog_organization_properties.0.enterprise_app_id",
			"datadog_organization_properties.0.linking_auth_code",
			"datadog_organization_properties.0.linking_client_id",
			"datadog_organization_properties.0.redirect_uri"),
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
		data.ImportStep("user_info",
			"user_info.0.name",
			"user_info.0.email_address",
			"datadog_organization_properties",
			"datadog_organization_properties.0",
			"datadog_organization_properties.0.id",
			"datadog_organization_properties.0.name",
			"datadog_organization_properties.0.api_key",
			"datadog_organization_properties.0.application_key",
			"datadog_organization_properties.0.enterprise_app_id",
			"datadog_organization_properties.0.linking_auth_code",
			"datadog_organization_properties.0.linking_client_id",
			"datadog_organization_properties.0.redirect_uri"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user_info",
			"user_info.0.name",
			"user_info.0.email_address",
			"datadog_organization_properties",
			"datadog_organization_properties.0",
			"datadog_organization_properties.0.id",
			"datadog_organization_properties.0.name",
			"datadog_organization_properties.0.api_key",
			"datadog_organization_properties.0.application_key",
			"datadog_organization_properties.0.enterprise_app_id",
			"datadog_organization_properties.0.linking_auth_code",
			"datadog_organization_properties.0.linking_client_id",
			"datadog_organization_properties.0.redirect_uri"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user_info",
			"user_info.0.name",
			"user_info.0.email_address",
			"datadog_organization_properties",
			"datadog_organization_properties.0",
			"datadog_organization_properties.0.id",
			"datadog_organization_properties.0.name",
			"datadog_organization_properties.0.api_key",
			"datadog_organization_properties.0.application_key",
			"datadog_organization_properties.0.enterprise_app_id",
			"datadog_organization_properties.0.linking_auth_code",
			"datadog_organization_properties.0.linking_client_id",
			"datadog_organization_properties.0.redirect_uri"),
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
		name = "test-terraform-%d"
		resource_group_name = azurerm_resource_group.test.name
		location = "EAST US 2 EUAP"
		datadog_organization_properties {
			api_key = ""
			application_key = ""
		}
		user_info {
			name          = "Vidhi Kothari"
			email_address = "vidhi.kothari@microsoft.com"
		}
		sku {
			name = "Linked"
		}
		identity {
			type = "SystemAssigned"
		}
	}
`, r.template(data), data.RandomInteger%100)
}

func (r DatadogMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
	
	resource "azurerm_datadog_monitor" "test" {
		name = azurerm_datadog_monitor.test.name
		resource_group_name = azurerm_resource_group.test.name
		location = "EAST US 2 EUAP"
		datadog_organization_properties {
			api_key = ""
			application_key = ""
		}
		user_info {
			name          = "Vidhi Kothari"
			email_address = "vidhi.kothari@microsoft.com"
		}
		sku {
			name = "Linked"
		}
		identity {
			type = "SystemAssigned"
		}
		monitoring_status = flase
		tags = {
			ENV = "Test"
		}
	}
`, r.template(data))
}

func (r DatadogMonitorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

	resource "azurerm_datadog_monitor" "import" {
	name                = azurerm_datadog_monitor.test.name
	resource_group_name =  azurerm_resource_group.test.name
	location            = azurerm_datadog_monitor.test.location
	}
`, r.template(data))
}

func (r DatadogMonitorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

	resource "azurerm_datadog_monitor" "test" {
	name                = "test-terraform-%d"
	resource_group_name = azurerm_resource_group.test.name
	location            = azurerm_resource_group.test.location
	datadog_organization_properties {
		api_key = ""
		application_key = ""
		enterprise_app_id = ""
		linking_auth_code = ""
		linking_client_id = ""
		redirect_uri      = ""
	}
	identity {
		type = "SystemAssigned"
	}
	sku {
		name = "Linked"
	}
	user_info {
		name          = "Vidhi Kothari"
		email_address = "vidhi.kothari@microsoft.com"
		phone_number  = ""
	}
	monitoring_status = true
	tags = {
		ENV = "Test"
	}
	}
`, r.template(data), data.RandomInteger%100)
}
