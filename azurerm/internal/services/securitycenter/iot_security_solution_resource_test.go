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
  unmasked_ip_logging_enabled = true
  log_analytics_workspace_id  = azurerm_log_analytics_workspace.test.id
  export                      = ["RawEvents"]

  recommendation {
    iot_acr_authentication_enabled               = false
    iot_agent_send_unutilized_msg_enabled        = false
    iot_baseline_enabled                         = false
    iot_edge_hub_mem_optimize_enabled            = false
    iot_edge_logging_option_enabled              = false
    iot_inconsistent_module_settings_enabled     = false
    iot_install_agent_enabled                    = false
    iot_ip_filter_deny_all_enabled               = false
    iot_ip_filter_permissive_rule_enabled        = false
    iot_open_ports_enabled                       = false
    iot_permissive_firewall_policy_enabled       = false
    iot_permissive_input_firewall_rules_enabled  = false
    iot_permissive_output_firewall_rules_enabled = false
    iot_privileged_docker_options_enabled        = false
    iot_shared_credentials_enabled               = false
    iot_vulnerable_tls_cipher_suite_enabled      = false
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
