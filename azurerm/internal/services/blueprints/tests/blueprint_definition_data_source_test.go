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
					resource.TestCheckResourceAttr(data.ResourceName, "target_scope", "subscription"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "time_created"),
				),
			},
		},
	})
}

// lintignore:AT001
func TestAccDataSourceBlueprintDefinition_basicAtManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_definition", "test")

	// TODO - Update when the AccTest environment is capable of supporting MG level testing. For now this will fail.

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlueprintDefinition_basicAtManagementGroup("SomeGroup"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "time_created"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "last_modified"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "target_scope"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
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
  name       = "simpleBlueprint"
  scope_type = "managementGroup"
  scope_name = %s
}

`, managementGroup)
}
