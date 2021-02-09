package blueprints_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type BlueprintDefinitionDataSource struct {
}

// lintignore:AT001
func TestAccBlueprintDefinitionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_definition", "test")
	r := BlueprintDefinitionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("description").HasValue("Acceptance Test stub for Blueprints at Subscription"),
				check.That(data.ResourceName).Key("name").HasValue("testAcc_basicSubscription"),
				check.That(data.ResourceName).Key("last_modified").Exists(),
				check.That(data.ResourceName).Key("target_scope").HasValue("subscription"),
				check.That(data.ResourceName).Key("time_created").Exists(),
			),
		},
	})
}

// lintignore:AT001
func TestAccBlueprintDefinitionDataSource_basicAtRootManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_definition", "test")
	r := BlueprintDefinitionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicAtRootManagementGroup(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("testAcc_basicRootManagementGroup"),
				check.That(data.ResourceName).Key("time_created").Exists(),
				check.That(data.ResourceName).Key("last_modified").Exists(),
				check.That(data.ResourceName).Key("target_scope").HasValue("subscription"),
			),
		},
	})
}

func TestAccBlueprintDefinitionDataSource_basicAtChildManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_definition", "test")
	r := BlueprintDefinitionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicAtManagementGroup("testAcc_staticStubGroup"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("testAcc_staticStubManagementGroup"),
				check.That(data.ResourceName).Key("time_created").Exists(),
				check.That(data.ResourceName).Key("last_modified").Exists(),
				check.That(data.ResourceName).Key("target_scope").HasValue("subscription"),
			),
		},
	})
}

func (BlueprintDefinitionDataSource) basic(data acceptance.TestData) string {
	subscription := data.Client().SubscriptionIDAlt
	return fmt.Sprintf(`
provider "azurerm" {
  subscription_id = "%s"
  features {}
}

data "azurerm_subscription" "current" {}

data "azurerm_blueprint_definition" "test" {
  name     = "testAcc_basicSubscription"
  scope_id = data.azurerm_subscription.current.id
}

`, subscription)
}

func (BlueprintDefinitionDataSource) basicAtManagementGroup(managementGroup string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_management_group" "test" {
  name = "%s"
}

data "azurerm_blueprint_definition" "test" {
  name     = "testAcc_staticStubManagementGroup"
  scope_id = data.azurerm_management_group.test.id
}

`, managementGroup)
}

func (BlueprintDefinitionDataSource) basicAtRootManagementGroup() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_management_group" "root" {
  name = data.azurerm_client_config.current.tenant_id
}

data "azurerm_blueprint_definition" "test" {
  name     = "testAcc_basicRootManagementGroup"
  scope_id = data.azurerm_management_group.root.id
}
`
}
