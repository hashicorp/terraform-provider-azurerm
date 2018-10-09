package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSecurityCenterSubscriptionPricing_update(t *testing.T) {
	resourceName := "azurerm_securitycenter_subscription_pricing.test"

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
		},
	})
}

func testCheckAzureRMSecurityCenterSubscriptionPricingExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		pricingName := rs.Primary.Attributes["pricings"]

		client := testAccProvider.Meta().(*ArmClient).securityCenterPricingClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
resource "azurerm_securitycenter_subscription_pricing" "test" {
    tier = "%s"
}
`, tier)
}
