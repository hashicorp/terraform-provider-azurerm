package blueprints_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

// lintignore:AT001
func TestAccDataSourceBlueprintDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlueprintDefinition_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Acceptance Test stub for Blueprints at Subscription"),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "testAcc_basicSubscription"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "last_modified"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_scope", "subscription"), // Only subscriptions can be targets
					resource.TestCheckResourceAttrSet(data.ResourceName, "time_created"),
				),
			},
		},
	})
}

// lintignore:AT001
func TestAccDataSourceBlueprintDefinition_basicAtRootManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlueprintDefinition_basicAtRootManagementGroup(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", "testAcc_basicRootManagementGroup"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "time_created"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "last_modified"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_scope", "subscription"), // Only subscriptions can be targets
				),
			},
		},
	})
}

func TestAccDataSourceBlueprintDefinition_basicAtChildManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlueprintDefinition_basicAtManagementGroup("testAcc_staticStubGroup"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", "testAcc_staticStubManagementGroup"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "time_created"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "last_modified"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_scope", "subscription"), // Only subscriptions can be targets
				),
			},
		},
	})
}

func testAccDataSourceBlueprintDefinition_basic(data acceptance.TestData) string {
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

func testAccDataSourceBlueprintDefinition_basicAtManagementGroup(managementGroup string) string {
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

func testAccDataSourceBlueprintDefinition_basicAtRootManagementGroup() string {
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
