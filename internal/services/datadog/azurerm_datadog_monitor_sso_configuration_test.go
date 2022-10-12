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

type SSODatadogMonitorResource struct{}

func TestAccDatadogMonitorSSO_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_sso_configuration", "test")
	r := SSODatadogMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("singlesignon_state").HasValue("Enable"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatadogMonitorSSO_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_sso_configuration", "test")
	r := SSODatadogMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("singlesignon_state").HasValue("Enable"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("singlesignon_state").HasValue("Disable"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("singlesignon_state").HasValue("Enable"),
			),
		},
		data.ImportStep(),
	})
}

func (r SSODatadogMonitorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DatadogSingleSignOnConfigurationsID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Datadog.SingleSignOnConfigurationsClient.Get(ctx, id.ResourceGroup, id.MonitorName, id.SingleSignOnConfigurationName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}

	if *resp.Properties.EnterpriseAppID == "" {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r SSODatadogMonitorResource) template(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger%100, os.Getenv("ARM_TEST_DATADOG_API_KEY"), os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY"))
}

func (r SSODatadogMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_datadog_monitor_sso_configuration" "test" {
  datadog_monitor_id        = azurerm_datadog_monitor.test.id
  singlesignon_state        = "Enable"
  enterprise_application_id = "183bc0b4-c560-4a55-8b7e-3eac5ad18774"
}
`, r.template(data))
}

func (r SSODatadogMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_datadog_monitor_sso_configuration" "test" {
  datadog_monitor_id        = azurerm_datadog_monitor.test.id
  singlesignon_state        = "Disable"
  enterprise_application_id = "183bc0b4-c560-4a55-8b7e-3eac5ad18774"
}
`, r.template(data))
}
