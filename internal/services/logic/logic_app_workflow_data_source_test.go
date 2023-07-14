// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type LogicAppWorkflowDataSource struct{}

func TestAccLogicAppWorkflowDataSource_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
			),
		},
	})
}

func TestAccLogicAppWorkflowDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("parameters.%").HasValue("0"),
				check.That(data.ResourceName).Key("connector_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("connector_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccLogicAppWorkflowDataSource_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("parameters.%").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Source").HasValue("AcceptanceTests"),
			),
		},
	})
}

func (LogicAppWorkflowDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_logic_app_workflow" "test" {
  name                = azurerm_logic_app_workflow.test.name
  resource_group_name = azurerm_logic_app_workflow.test.resource_group_name
}
`, LogicAppWorkflowResource{}.empty(data))
}

func (LogicAppWorkflowDataSource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_logic_app_workflow" "test" {
  name                = azurerm_logic_app_workflow.test.name
  resource_group_name = azurerm_logic_app_workflow.test.resource_group_name
}
`, LogicAppWorkflowResource{}.tags(data))
}

func (LogicAppWorkflowDataSource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "workflow-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_logic_app_workflow" "test" {
  name                = azurerm_logic_app_workflow.test.name
  resource_group_name = azurerm_logic_app_workflow.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary)
}
