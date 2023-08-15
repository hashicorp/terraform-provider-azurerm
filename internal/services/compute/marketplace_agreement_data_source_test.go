// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	offer := "barracuda-email-security-gateway"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: MarketplaceAgreementResource{}.empty(),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientWithoutResource(MarketplaceAgreementResource{}.cancelExistingAgreement(offer)),
			),
		},
		{
			Config: r.basic(offer),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("license_text_link").Exists(),
				check.That(data.ResourceName).Key("privacy_policy_link").Exists(),
			),
		},
	})
}

func (MarketplaceAgreementDataSource) basic(offer string) string {
	return fmt.Sprintf(`
%s

data "azurerm_marketplace_agreement" "test" {
  publisher = "barracudanetworks"
  offer     = "%s"
  plan      = "hourly"
}
`, MarketplaceAgreementResource{}.basic(offer), offer)
}
