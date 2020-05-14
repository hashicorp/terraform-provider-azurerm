package tests

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
				Config: testAccDataSourceBlueprintDefinition_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Acceptance Test stub for Blueprints at Subscription"),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "testAcc_basicSubscription"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "last_modified"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_scope", "subscription"), // Only subscriptions can be targets
					resource.TestCheckResourceAttrSet(data.ResourceName, "time_created"),
					resource.TestCheckResourceAttr(data.ResourceName, "scope_type", "subscription"),
				),
			},
		},
	})
}

// lintignore:AT001
func TestAccDataSourceBlueprintDefinition_basicAtRootManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_definition", "test")

	// TODO - Update when the AccTest environment is capable of supporting MG level testing. For now this will fail.

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
					resource.TestCheckResourceAttr(data.ResourceName, "scope_type", "managementGroup"),
				),
			},
		},
	})
}

func TestAccDataSourceBlueprintDefinition_basicAtChildManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_definition", "test")

	// TODO - Update when the AccTest environment is capable of supporting MG level testing. For now this will fail.

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
					resource.TestCheckResourceAttr(data.ResourceName, "scope_type", "managementGroup"),
				),
			},
		},
	})
}

func testAccDataSourceBlueprintDefinition_basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_blueprint_definition" "test" {
  name       = "testAcc_basicSubscription"
  scope_type = "subscription"
  scope_name = data.azurerm_client_config.current.subscription_id
}

`
}

func testAccDataSourceBlueprintDefinition_basicAtManagementGroup(managementGroup string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_blueprint_definition" "test" {
  name       = "testAcc_staticStubManagementGroup"
  scope_type = "managementGroup"
  scope_name = "%s"
}

`, managementGroup)
}

func testAccDataSourceBlueprintDefinition_basicAtRootManagementGroup() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_blueprint_definition" "test" {
  name       = "testAcc_basicRootManagementGroup"
  scope_type = "managementGroup"
  scope_name = data.azurerm_client_config.current.tenant_id
}

`
}
