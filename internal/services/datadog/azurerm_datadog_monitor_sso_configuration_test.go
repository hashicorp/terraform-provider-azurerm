package datadog_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/singlesignon"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SSODatadogMonitorResource struct{}

func TestAccDatadogMonitorSSO_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_DATADOG_API_KEY") == "" || os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_API_KEY and/or ARM_TEST_DATADOG_APPLICATION_KEY are not specified")
		return
	}
	if os.Getenv("ARM_TEST_ENTERPRISE_APP_ID") == "" {
		t.Skip("Skipping as Enterprise App Id for SAML is not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_sso_configuration", "test")
	r := SSODatadogMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on_enabled").HasValue("Enable"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatadogMonitorSSO_update(t *testing.T) {
	if os.Getenv("ARM_TEST_DATADOG_API_KEY") == "" || os.Getenv("ARM_TEST_DATADOG_APPLICATION_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_DATADOG_API_KEY and/or ARM_TEST_DATADOG_APPLICATION_KEY are not specified")
		return
	}
	if os.Getenv("ARM_TEST_ENTERPRISE_APP_ID") == "" {
		t.Skip("Skipping as Enterprise App Id for SAML is not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_sso_configuration", "test")
	r := SSODatadogMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on_enabled").HasValue("Enable"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on_enabled").HasValue("Disable"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("single_sign_on_enabled").HasValue("Enable"),
			),
		},
		data.ImportStep(),
	})
}

func (r SSODatadogMonitorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := singlesignon.ParseSingleSignOnConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Datadog.SingleSignOn.ConfigurationsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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
  single_sign_on_enabled    = "Enable"
  enterprise_application_id = %q
}
`, r.template(data), os.Getenv("ARM_TEST_ENTERPRISE_APP_ID"))
}

func (r SSODatadogMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_datadog_monitor_sso_configuration" "test" {
  datadog_monitor_id        = azurerm_datadog_monitor.test.id
  single_sign_on_enabled    = "Disable"
  enterprise_application_id = %q
}
`, r.template(data), os.Getenv("ARM_TEST_ENTERPRISE_APP_ID"))
}
