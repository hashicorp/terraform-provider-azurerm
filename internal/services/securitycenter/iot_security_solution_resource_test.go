// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type IotSecuritySolutionResource struct{}

func TestAccIotSecuritySolution_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_solution", "test")
	r := IotSecuritySolutionResource{}

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

func TestAccIotSecuritySolution_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_solution", "test")
	r := IotSecuritySolutionResource{}

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

func TestAccIotSecuritySolution_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_solution", "test")
	r := IotSecuritySolutionResource{}

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

func TestAccIotSecuritySolution_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_solution", "test")
	r := IotSecuritySolutionResource{}

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

func TestAccIotSecuritySolution_additionalWorkspace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_solution", "test")
	r := IotSecuritySolutionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.additionalWorkspace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateAdditionalWorkspace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (IotSecuritySolutionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IotSecuritySolutionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.IotSecuritySolutionClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Security Center Iot Security Solution %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return pointer.To(resp.ID != nil), nil
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

  depends_on = [azurerm_role_assignment.iot-assign]
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

  depends_on = [azurerm_role_assignment.iot-assign]
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

# "AzureSecurityOfThings" and "Security" will be created automatically by service, so we create them manually in case the resource group can't be deleted.

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "AzureSecurityOfThings"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/AzureSecurityOfThings"
  }
}

resource "azurerm_log_analytics_solution" "test2" {
  solution_name         = "Security"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/Security"
  }
}

resource "azurerm_iot_security_solution" "test" {
  name                       = "acctest-Iot-Security-Solution-%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  display_name               = "Iot Security Solution"
  iothub_ids                 = [azurerm_iothub.test.id]
  enabled                    = true
  log_unmasked_ips_enabled   = true
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  events_to_export           = ["RawEvents"]
  disabled_data_sources      = ["TwinData"]

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

  query_for_resources    = "where type != \"microsoft.devices/iothubs\" | where name contains \"iot\""
  query_subscription_ids = [data.azurerm_client_config.test.subscription_id]

  tags = {
    "Env" : "Staging"
  }

  depends_on = [azurerm_role_assignment.iot-assign, azurerm_role_assignment.law-assign]
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r IotSecuritySolutionResource) additionalWorkspace(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-law-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_iot_security_solution" "test" {
  name                = "acctest-Iot-Security-Solution-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  display_name        = "Iot Security Solution"
  iothub_ids          = [azurerm_iothub.test.id]

  additional_workspace {
    data_types   = ["Alerts"]
    workspace_id = azurerm_log_analytics_workspace.test.id
  }

  depends_on = [azurerm_role_assignment.iot-assign]
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r IotSecuritySolutionResource) updateAdditionalWorkspace(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-law-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_workspace" "test2" {
  name                = "acctest-law2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_iot_security_solution" "test" {
  name                = "acctest-Iot-Security-Solution-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  display_name        = "Iot Security Solution"
  iothub_ids          = [azurerm_iothub.test.id]

  additional_workspace {
    data_types   = ["Alerts", "RawEvents"]
    workspace_id = azurerm_log_analytics_workspace.test2.id
  }

  depends_on = [azurerm_role_assignment.iot-assign, azurerm_role_assignment.law-assign]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
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

data "azuread_service_principal" "iotsec" {
  display_name = "Azure Security for IoT"
}

resource "azurerm_role_assignment" "iot-assign" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Contributor"
  principal_id         = data.azuread_service_principal.iotsec.object_id
}

resource "azurerm_role_assignment" "law-assign" {
  scope                = azurerm_iothub.test.id
  role_definition_name = "Log Analytics Contributor"
  principal_id         = data.azuread_service_principal.iotsec.object_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
