package datadog_test

import (
	"context"
	"fmt"
	"os"
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
	if os.Getenv("ARM_TEST_DATADOG_API_KEY") == "" || os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_API_KEY and/or ARM_TEST_DATADOG_APPLICATION_KEY are not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
	r := DatadogMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user",
			"user.0.name",
			"user.0.email",
			"datadog_organization",
			"datadog_organization.0",
			"datadog_organization.0.id",
			"datadog_organization.0.name",
			"datadog_organization.0.api_key",
			"datadog_organization.0.application_key",
			"datadog_organization.0.enterprise_app_id",
			"datadog_organization.0.linking_auth_code",
			"datadog_organization.0.linking_client_id",
			"datadog_organization.0.redirect_uri"),
	})
}

func TestAccDatadogMonitor_requiresImport(t *testing.T) {
	if os.Getenv("ARM_TEST_DATADOG_API_KEY") == "" || os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_API_KEY and/or ARM_TEST_DATADOG_APPLICATION_KEY are not specified")
		return
	}
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
	if os.Getenv("ARM_TEST_DATADOG_API_KEY") == "" || os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_API_KEY and/or ARM_TEST_DATADOG_APPLICATION_KEY are not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
	r := DatadogMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user",
			"user.0.name",
			"user.0.email",
			"datadog_organization",
			"datadog_organization.0",
			"datadog_organization.0.id",
			"datadog_organization.0.name",
			"datadog_organization.0.api_key",
			"datadog_organization.0.application_key",
			"datadog_organization.0.enterprise_app_id",
			"datadog_organization.0.linking_auth_code",
			"datadog_organization.0.linking_client_id",
			"datadog_organization.0.redirect_uri"),
	})
}

func TestAccDatadogMonitor_update(t *testing.T) {
	if os.Getenv("ARM_TEST_DATADOG_API_KEY") == "" || os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_API_KEY and/or ARM_TEST_DATADOG_APPLICATION_KEY are not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor", "test")
	r := DatadogMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user",
			"user.0.name",
			"user.0.email",
			"datadog_organization",
			"datadog_organization.0",
			"datadog_organization.0.id",
			"datadog_organization.0.name",
			"datadog_organization.0.api_key",
			"datadog_organization.0.application_key",
			"datadog_organization.0.enterprise_app_id",
			"datadog_organization.0.linking_auth_code",
			"datadog_organization.0.linking_client_id",
			"datadog_organization.0.redirect_uri"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user",
			"user.0.name",
			"user.0.email",
			"datadog_organization",
			"datadog_organization.0",
			"datadog_organization.0.id",
			"datadog_organization.0.name",
			"datadog_organization.0.api_key",
			"datadog_organization.0.application_key",
			"datadog_organization.0.enterprise_app_id",
			"datadog_organization.0.linking_auth_code",
			"datadog_organization.0.linking_client_id",
			"datadog_organization.0.redirect_uri"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user",
			"user.0.name",
			"user.0.email",
			"datadog_organization",
			"datadog_organization.0",
			"datadog_organization.0.id",
			"datadog_organization.0.name",
			"datadog_organization.0.api_key",
			"datadog_organization.0.application_key",
			"datadog_organization.0.enterprise_app_id",
			"datadog_organization.0.linking_auth_code",
			"datadog_organization.0.linking_client_id",
			"datadog_organization.0.redirect_uri"),
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
  name                = "acctest-datadog-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "WEST US 2"
  datadog_organization {
    api_key         = %q
    application_key = %q
  }
  user {
    name  = "Test Datadog"
    email = "abc@xyz.com"
  }
  sku_name = "Linked"
  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger%100, os.Getenv("ARM_TEST_DATADOG_API_KEY"), os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY"))
}

func (r DatadogMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-datadog-%d"
  location = "%s"
}
resource "azurerm_datadog_monitor" "test" {
  name                = "acctest-datadog-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "WEST US 2"
  datadog_organization {
    api_key         = %q
    application_key = %q
  }
  user {
    name  = "Test Datadog"
    email = "abc@xyz.com"
  }
  sku_name = "Linked"
  identity {
    type = "SystemAssigned"
  }
  monitoring_enabled = false
  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger%100, os.Getenv("ARM_TEST_DATADOG_API_KEY"), os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY"))
}

func (r DatadogMonitorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
resource "azurerm_datadog_monitor" "import" {
  name                = azurerm_datadog_monitor.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_datadog_monitor.test.location
  datadog_organization {
    api_key         = %q
    application_key = %q
  }
  user {
    name  = "Test Datadog"
    email = "abc@xyz.com"
  }
  sku_name = "Linked"
  identity {
    type = "SystemAssigned"
  }
}
`, r.basic(data), os.Getenv("ARM_TEST_DATADOG_API_KEY"), os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY"))
}

func (r DatadogMonitorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
resource "azurerm_datadog_monitor" "test" {
  name                = "acctest-datadog-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datadog_organization {
    api_key           = %q
    application_key   = %q
    enterprise_app_id = ""
    linking_auth_code = ""
    linking_client_id = ""
    redirect_uri      = ""
  }
  identity {
    type = "SystemAssigned"
  }
  sku_name = "Linked"
  user {
    name         = "Test Datadog"
    email        = "abc@xyz.com"
    phone_number = ""
  }
  monitoring_enabled = true
  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger%100, os.Getenv("ARM_TEST_DATADOG_API_KEY"), os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY"))
}
