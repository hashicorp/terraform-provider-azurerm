// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type RoleManagementPolicyDataSource struct{}

func TestAccRoleManagementPolicyDataSource_managementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_management_policy", "test")
	r := RoleManagementPolicyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.managementGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccRoleManagementPolicyDataSource_resourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_management_policy", "test")
	r := RoleManagementPolicyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.resourceGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccRoleManagementPolicyDataSource_subscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_management_policy", "test")
	r := RoleManagementPolicyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.subscription(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccRoleManagementPolicyDataSource_resource(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_management_policy", "test")
	r := RoleManagementPolicyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.resource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func (RoleManagementPolicyDataSource) managementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {}

data "azurerm_client_config" "current" {
}

resource "azurerm_management_group" "test" {
  name = "acctest%[1]s"
}

data "azurerm_role_definition" "contributor" {
  name  = "Contributor"
  scope = azurerm_management_group.test.id
}

data "azurerm_role_management_policy" "test" {
  role_definition_id = data.azurerm_role_definition.contributor.id
  scope              = azurerm_management_group.test.id
}
`, data.RandomString)
}

func (RoleManagementPolicyDataSource) resourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]s"
  location = "%[2]s"
}

data "azurerm_role_definition" "contributor" {
  name  = "Contributor"
  scope = azurerm_resource_group.test.id
}

data "azurerm_role_management_policy" "test" {
  role_definition_id = data.azurerm_role_definition.contributor.id
  scope              = azurerm_resource_group.test.id
}
`, data.RandomString, data.Locations.Primary)
}

func (RoleManagementPolicyDataSource) subscription(data acceptance.TestData) string {
	return `
provider "azurerm" {}

data "azurerm_subscription" "test" {}

data "azurerm_role_definition" "contributor" {
  name  = "Contributor"
  scope = data.azurerm_subscription.test.id
}

data "azurerm_role_management_policy" "test" {
  role_definition_id = data.azurerm_role_definition.contributor.id
  scope              = data.azurerm_subscription.test.id
}
`
}

func (RoleManagementPolicyDataSource) resource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]s"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accteststg%[1]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

data "azurerm_role_definition" "contributor" {
  name  = "Contributor"
  scope = azurerm_resource_group.test.id
}

data "azurerm_role_management_policy" "test" {
  role_definition_id = data.azurerm_role_definition.contributor.id
  scope              = azurerm_storage_account.test.id
}
`, data.RandomString, data.Locations.Primary)
}
