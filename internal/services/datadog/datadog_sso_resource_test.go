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
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_sso", "test")
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

// func TestAccDatadogMonitorSSO_update(t *testing.T) {
// 	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_sso", "test")
// 	r := SSODatadogMonitorResource{}
// 	data.ResourceTest(t, r, []acceptance.TestStep{
// 		{
// 			Config: r.basic(data),
// 			Check: acceptance.ComposeTestCheckFunc(
// 				check.That(data.ResourceName).ExistsInAzure(r),
// 				check.That(data.ResourceName).Key("singlesignon_state").HasValue("Enable"),
// 			),
// 		},
// 		data.ImportStep(),
// 		{
// 			Config: r.update(data),
// 			Check: acceptance.ComposeTestCheckFunc(
// 				check.That(data.ResourceName).ExistsInAzure(r),
// 				check.That(data.ResourceName).Key("singlesignon_state").HasValue("Disable"),
// 			),
// 		},
// 		data.ImportStep(),
// 		{
// 			Config: r.basic(data),
// 			Check: acceptance.ComposeTestCheckFunc(
// 				check.That(data.ResourceName).ExistsInAzure(r),
// 				check.That(data.ResourceName).Key("singlesignon_state").HasValue("Enable"),
// 			),
// 		},
// 		data.ImportStep(),
// 	})
// }

func (r SSODatadogMonitorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DatadogSingleSignOnConfigurationsID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Datadog.SingleSignOnConfigurationsClient.Get(ctx, id.ResourceGroup, id.MonitorName, id.SingleSignOnConfigurationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		//return nil, fmt.Errorf("retrieving Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

// func (r SSODatadogMonitorResource) template(data acceptance.TestData) string {
// 	return fmt.Sprintf(`
// provider "azurerm" {
//   features {}
// }

// data "azurerm_resource_group" "test" {
//   name     = "acctest-datadog"
// }
// `)
// }

func (r SSODatadogMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	
	resource "azurerm_datadog_monitor_sso" "test" {
		name = "test-terraform-6747642"
		resource_group_name = "acctest-datadog"
		singlesignon_state = "Enable"
		enterpriseapp_id = "183bc0b4-c560-4a55-8b7e-3eac5ad18774"
	}
`)
}

func (r SSODatadogMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`

	resource "azurerm_datadog_monitor_sso" "test" {
		name = "test-terraform-6747642"
		resource_group_name = "acctest-datadog"
		singlesignon_state = "Disable"
		enterpriseapp_id = "183bc0b4-c560-4a55-8b7e-3eac5ad18774"
	}
`)
}
