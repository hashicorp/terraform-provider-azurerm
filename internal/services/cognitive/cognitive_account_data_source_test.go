// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CognitiveAccountDataSource struct{}

func TestAccCognitiveAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cognitive_account", "test")
	r := CognitiveAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("outbound_network_access_restricted").HasValue("false"),
				check.That(data.ResourceName).Key("project_management_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Acceptance").HasValue("Test"),
				check.That(data.ResourceName).Key("custom_subdomain_name").IsNotEmpty(),
			),
		},
	})
}

func TestAccCognitiveAccountDataSource_aiServices_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cognitive_account", "test")
	r := CognitiveAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.aiServices_complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("outbound_network_access_restricted").HasValue("false"),
				check.That(data.ResourceName).Key("dynamic_throttling_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("project_management_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("S0"),
				check.That(data.ResourceName).Key("tags.Acceptance").HasValue("Test"),
				check.That(data.ResourceName).Key("custom_subdomain_name").IsNotEmpty(),
				check.That(data.ResourceName).Key("network_acls.0.default_action").HasValue("Deny"),
				check.That(data.ResourceName).Key("network_acls.0.virtual_network_rules.0.ignore_missing_vnet_service_endpoint").HasValue("false"),
				check.That(data.ResourceName).Key("network_acls.0.virtual_network_rules.0.subnet_id").IsNotEmpty(),
			),
		},
	})
}

func TestAccCognitiveAccountDataSource_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cognitive_account", "test")
	r := CognitiveAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("Face"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
			),
		},
	})
}

func (CognitiveAccountDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Face"
  sku_name            = "S0"

  tags = {
    Acceptance = "Test"
  }
}

data "azurerm_cognitive_account" "test" {
  name                = azurerm_cognitive_account.test.name
  resource_group_name = azurerm_cognitive_account.test.resource_group_name
}
`, CognitiveAccountDataSource{}.template(data), data.RandomInteger)
}

func (CognitiveAccountDataSource) aiServices_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_cognitive_account" "test" {
  name                = azurerm_cognitive_account.test.name
  resource_group_name = azurerm_cognitive_account.test.resource_group_name
}
`, CognitiveAccountResource{}.aiServices_complete(data))
}

func (CognitiveAccountDataSource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Face"
  sku_name            = "S0"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Acceptance = "Test"
  }
}

data "azurerm_cognitive_account" "test" {
  name                = azurerm_cognitive_account.test.name
  resource_group_name = azurerm_cognitive_account.test.resource_group_name
}
`, CognitiveAccountDataSource{}.template(data), data.RandomInteger)
}

func (CognitiveAccountDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
