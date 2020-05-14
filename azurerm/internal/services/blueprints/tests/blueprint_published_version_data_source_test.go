package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

// lintignore:AT001
func TestAccDataSourceBlueprintPublishedVersion_atSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_published_version", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlueprintPublishedVersion_atSubscription("testAcc_basicSubscription", "v0.1_testAcc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "target_scope", "subscription"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Acceptance Test stub for Blueprints at Subscription"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "time_created"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
				),
			},
		},
	})
}

func TestAccDataSourceBlueprintPublishedVersion_atRootManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_published_version", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlueprintPublishedVersion_atRootManagementGroup("testAcc_basicRootManagementGroup", "v0.1_testAcc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "target_scope", "subscription"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Acceptance Test stub for Blueprints at Root Management Group"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "time_created"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
				),
			},
		},
	})
}

func TestAccDataSourceBlueprintPublishedVersion_atChildManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_published_version", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlueprintPublishedVersion_atChildManagementGroup("testAcc_staticStubGroup", "testAcc_staticStubManagementGroup", "v0.1_testAcc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "target_scope", "subscription"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Acceptance Test stub for Blueprints at Child Management Group"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "time_created"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
				),
			},
		},
	})
}

func testAccDataSourceBlueprintPublishedVersion_atSubscription(bpName, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

data "azurerm_blueprint_published_version" "test" {
  subscription_id = data.azurerm_client_config.current.subscription_id
  blueprint_name  = "%s"
  version         = "%s"
}
`, bpName, version)
}

func testAccDataSourceBlueprintPublishedVersion_atRootManagementGroup(bpName, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

data "azurerm_blueprint_published_version" "test" {
  management_group = data.azurerm_client_config.current.tenant_id
  blueprint_name   = "%s"
  version          = "%s"
}
`, bpName, version)
}

func testAccDataSourceBlueprintPublishedVersion_atChildManagementGroup(mg, bpName, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

data "azurerm_blueprint_published_version" "test" {
  management_group = "%s"
  blueprint_name   = "%s"
  version          = "%s"
}
`, mg, bpName, version)
}
