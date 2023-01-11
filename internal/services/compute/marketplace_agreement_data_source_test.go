package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MarketplaceAgreementDataSource struct{}

func TestAccDataSourceMarketplaceAgreement_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_marketplace_agreement", "test")
	r := MarketplaceAgreementDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("license_text_link").Exists(),
				check.That(data.ResourceName).Key("privacy_policy_link").Exists(),
			),
		},
	})
}

func (MarketplaceAgreementDataSource) basic() string {
	return fmt.Sprintf(`
%s

data "azurerm_marketplace_agreement" "test" {
  publisher = "barracudanetworks"
  offer     = "waf"
  plan      = "hourly"
}
`, MarketplaceAgreementResource{}.basic("waf"))
}
