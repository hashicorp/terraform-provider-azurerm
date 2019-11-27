package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceAzureRMSubscription_current(t *testing.T) {
	resourceName := "data.azurerm_subscription.current"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSubscription_currentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "subscription_id"),
					testCheckAzureRMSubscriptionId(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "state", "Enabled"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSubscription_specific(t *testing.T) {
	resourceName := "data.azurerm_subscription.specific"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSubscription_specificConfig(os.Getenv("ARM_SUBSCRIPTION_ID")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "subscription_id"),
					testCheckAzureRMSubscriptionId(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttrSet(resourceName, "location_placement_id"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_id"),
					resource.TestCheckResourceAttrSet(resourceName, "spending_limit"),
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
		client := testAccProvider.Meta().(*ArmClient)
		if subscriptionId != client.Account.SubscriptionId {
			return fmt.Errorf("%s: Attribute '%s' expected \"%s\", got \"%s\"", resourceName, attributeName, client.Account.SubscriptionId, subscriptionId)
		}

		return nil
	}
}

const testAccDataSourceAzureRMSubscription_currentConfig = `
data "azurerm_subscription" "current" {}
`

func testAccDataSourceAzureRMSubscription_specificConfig(subscriptionId string) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "specific" {
  subscription_id = "%s"
}
`, subscriptionId)
}
