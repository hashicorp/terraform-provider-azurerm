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

func (r SSODatadogMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	















resource "azurerm_datadog_monitor_sso_configuration" "test" {
  datadog_monitor_id        = "/subscriptions/5a611eed-e33a-44e8-92b1-3f6bf835905e/resourceGroups/acctest-datadog/providers/Microsoft.Datadog/monitors/test-terraform-acctests"
  singlesignon_state        = "Enable"
  enterprise_application_id = "183bc0b4-c560-4a55-8b7e-3eac5ad18774"
}
`)
}

func (r SSODatadogMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_datadog_monitor_sso_configuration" "test" {
  datadog_monitor_id        = "/subscriptions/5a611eed-e33a-44e8-92b1-3f6bf835905e/resourceGroups/acctest-datadog/providers/Microsoft.Datadog/monitors/test-terraform-acctests"
  singlesignon_state        = "Disable"
  enterprise_application_id = "183bc0b4-c560-4a55-8b7e-3eac5ad18774"
}
`)
}
