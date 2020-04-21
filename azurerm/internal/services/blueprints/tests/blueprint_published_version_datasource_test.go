package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceBlueprintPublishedVersion_basicSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_published_version", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlueprintPublishedVersion_basicSubscription("testAcc_basicSubscription", "testAcc"),
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

func testAccDataSourceBlueprintPublishedVersion_basicSubscription(bpName, version string) string {
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
