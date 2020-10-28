package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccDataSourceAzureRMSubscription_current(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subscription", "current")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSubscription_currentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "subscription_id"),
					testCheckAzureRMSubscriptionId(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "tenant_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "tags.%"),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "Enabled"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSubscription_specific(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subscription", "specific")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSubscription_specificConfig(os.Getenv("ARM_SUBSCRIPTION_ID")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "subscription_id"),
					testCheckAzureRMSubscriptionId(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "tenant_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location_placement_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "quota_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "spending_limit"),
				),
			},
		},
	})
}

func testCheckAzureRMSubscriptionId(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		attributeName := "subscription_id"
		subscriptionId := rs.Primary.Attributes[attributeName]
		client := acceptance.AzureProvider.Meta().(*clients.Client)
		if subscriptionId != client.Account.SubscriptionId {
			return fmt.Errorf("%s: Attribute '%s' expected \"%s\", got \"%s\"", resourceName, attributeName, client.Account.SubscriptionId, subscriptionId)
		}

		return nil
	}
}

const testAccDataSourceAzureRMSubscription_currentConfig = `
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}
`

func testAccDataSourceAzureRMSubscription_specificConfig(subscriptionId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  subscription_id = "%s"
}

data "azurerm_subscription" "specific" {
  subscription_id = "%s"
}
`, subscriptionId, subscriptionId)
}
