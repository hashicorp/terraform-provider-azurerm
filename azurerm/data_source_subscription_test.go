package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceAzureRMSubscription_current(t *testing.T) {
	resourceName := "data.azurerm_subscription.current"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSubscription_currentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "subscription_id"),
					testCheckAzureRMSubscriptionId(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttr(resourceName, "state", "Enabled"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSubscription_specific(t *testing.T) {
	resourceName := "data.azurerm_subscription.specific"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSubscription_specificConfig(os.Getenv("ARM_SUBSCRIPTION_ID")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "subscription_id"),
					testCheckAzureRMSubscriptionId(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "location_placement_id"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_id"),
					resource.TestCheckResourceAttrSet(resourceName, "spending_limit"),
				),
			},
		},
	})
}

func testCheckAzureRMSubscriptionId(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		attributeName := "subscription_id"
		subscriptionId := rs.Primary.Attributes[attributeName]
		client := testAccProvider.Meta().(*ArmClient)
		if subscriptionId != client.subscriptionId {
			return fmt.Errorf("%s: Attribute '%s' expected \"%s\", got \"%s\"", name, attributeName, client.subscriptionId, subscriptionId)
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
