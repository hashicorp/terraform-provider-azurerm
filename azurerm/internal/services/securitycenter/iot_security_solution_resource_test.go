package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type IotSecuritySolutionResource struct {
}

func TestAccIotSecuritySolution_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_solution", "test")
	r := IotSecuritySolutionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotSecuritySolution_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_solution", "test")
	r := IotSecuritySolutionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccIotSecuritySolution_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_solution", "test")
	r := IotSecuritySolutionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotSecuritySolution_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_solution", "test")
	r := IotSecuritySolutionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (IotSecuritySolutionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.IotSecuritySolutionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.IotSecuritySolutionClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Security Center Iot Security Solution %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r IotSecuritySolutionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iot_security_solution" "test" {
  name                = "acctest-Iot-Security-Solution-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  display_name        = "Iot Security Solution"
  iothub_ids          = [azurerm_iothub.test.id]
}
`, r.template(data), data.RandomInteger)
}

func (r IotSecuritySolutionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iot_security_solution" "import" {
  name                = azurerm_iot_security_solution.test.name
  resource_group_name = azurerm_iot_security_solution.test.resource_group_name
  location            = azurerm_iot_security_solution.test.location
  display_name        = azurerm_iot_security_solution.test.display_name
  iothub_ids          = [azurerm_iothub.test.id]
}
`, r.basic(data))
}

func (r IotSecuritySolutionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "test" {}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_iot_security_solution" "test" {
  name                        = "acctest-Iot-Security-Solution-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  display_name                = "Iot Security Solution"
  iothub_ids                  = [azurerm_iothub.test.id]
  enabled                     = true
  log_unmasked_ips_enabled = true
  log_analytics_workspace_id  = azurerm_log_analytics_workspace.test.id
  events_to_export                      = ["RawEvents"]

  recommendations_enabled {
    acr_authentication               = false
    agent_send_unutilized_msg        = false
    baseline                         = false
    edge_hub_mem_optimize            = false
    edge_logging_option              = false
    inconsistent_module_settings     = false
    install_agent                    = false
    ip_filter_deny_all               = false
    ip_filter_permissive_rule        = false
    open_ports                       = false
    permissive_firewall_policy       = false
    permissive_input_firewall_rules  = false
    permissive_output_firewall_rules = false
    privileged_docker_options        = false
    shared_credentials               = false
    vulnerable_tls_cipher_suite      = false
  }

  user_defined_resource {
    query            = "where type != \"microsoft.devices/iothubs\" | where name contains \"iot\""
    subscription_ids = [data.azurerm_client_config.test.subscription_id]
  }

  tags = {
    "Env" : "Staging"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r IotSecuritySolutionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-security-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
