package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SecurityCenterSubscriptionPricingResource struct {
}

func TestAccSecurityCenterSubscriptionPricing_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_subscription_pricing", "test")
	r := SecurityCenterSubscriptionPricingResource{}

	// lintignore:AT001
	data.ResourceSequentialTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.tier("Standard", "AppServices"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Standard"),
			),
		},
		data.ImportStep(),
		{
			Config: r.tier("Free", "AppServices"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Free"),
			),
		},
		data.ImportStep(),
	})
}

func (SecurityCenterSubscriptionPricingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SecurityCenterSubscriptionPricingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.PricingClient.Get(ctx, id.ResourceType)
	if err != nil {
		return nil, fmt.Errorf("reading Security Center Subscription Pricing (%s): %+v", id.ResourceType, err)
	}

	return utils.Bool(resp.PricingProperties != nil), nil
}

func (SecurityCenterSubscriptionPricingResource) tier(tier string, resource_type string) string {
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
