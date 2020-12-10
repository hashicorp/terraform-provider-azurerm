package securitycenter_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMSecurityCenterSubscriptionPricing_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_subscription_pricing", "test")

	// lintignore:AT001
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSecurityCenterSubscriptionPricing_tier("Standard", "AppServices"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterSubscriptionPricingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "Standard"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSecurityCenterSubscriptionPricing_tier("Free", "AppServices"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterSubscriptionPricingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "Free"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSecurityCenterSubscriptionPricingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.PricingClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		pricingName := rs.Primary.Attributes["pricings"]

		resp, err := client.Get(ctx, pricingName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Security Center Subscription Pricing %q was not found: %+v", pricingName, err)
			}

			return fmt.Errorf("Bad: Get: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSecurityCenterSubscriptionPricing_tier(tier string, resource_type string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_subscription_pricing" "test" {
  tier          = "%s"
  resource_type = "%s"
}
`, tier, resource_type)
}
