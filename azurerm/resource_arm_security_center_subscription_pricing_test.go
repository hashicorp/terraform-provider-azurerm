package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMSecurityCenterSubscriptionPricing_update(t *testing.T) {
	resourceName := "azurerm_security_center_subscription_pricing.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSecurityCenterSubscriptionPricing_tier("Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterSubscriptionPricingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tier", "Standard"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMSecurityCenterSubscriptionPricing_tier("Free"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterSubscriptionPricingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tier", "Free"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMSecurityCenterSubscriptionPricingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).SecurityCenter.PricingClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		pricingName := rs.Primary.Attributes["pricings"]

		resp, err := client.GetSubscriptionPricing(ctx, pricingName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Security Center Subscription Pricing %q was not found: %+v", pricingName, err)
			}

			return fmt.Errorf("Bad: GetSubscriptionPricing: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSecurityCenterSubscriptionPricing_tier(tier string) string {
	return fmt.Sprintf(`
resource "azurerm_security_center_subscription_pricing" "test" {
  tier = "%s"
}
`, tier)
}
