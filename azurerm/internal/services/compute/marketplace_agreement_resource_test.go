package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MarketplaceAgreementResource struct {
}

func TestAccMarketplaceAgreement(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":             testAccMarketplaceAgreement_basic,
			"requiresImport":    testAccMarketplaceAgreement_requiresImport,
			"agreementCanceled": testAccMarketplaceAgreement_agreementCanceled,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccMarketplaceAgreement_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_agreement", "test")
	r := MarketplaceAgreementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("license_text_link").Exists(),
				check.That(data.ResourceName).Key("privacy_policy_link").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func testAccMarketplaceAgreement_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_agreement", "test")
	r := MarketplaceAgreementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_marketplace_agreement"),
		},
	})
}

func testAccMarketplaceAgreement_agreementCanceled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_agreement", "test")
	r := MarketplaceAgreementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basicConfig,
			TestResource: r,
		}),
	})
}

func (t MarketplaceAgreementResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	publisher := id.Path["agreements"]
	offer := id.Path["offers"]
	plan := id.Path["plans"]

	resp, err := clients.Compute.MarketplaceAgreementsClient.Get(ctx, publisher, offer, plan)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Marketplace Agreement %q", id)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (MarketplaceAgreementResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	publisher := id.Path["agreements"]
	offer := id.Path["offers"]
	plan := id.Path["plans"]

	resp, err := client.Compute.MarketplaceAgreementsClient.Cancel(ctx, publisher, offer, plan)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, fmt.Errorf("marketplace agreement %q does not exist", id)
		}
		return nil, fmt.Errorf("canceling Marketplace Agreement : %+v", err)
	}

	return utils.Bool(true), nil
}

func (MarketplaceAgreementResource) basicConfig(_ acceptance.TestData) string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_marketplace_agreement" "test" {
  publisher = "barracudanetworks"
  offer     = "waf"
  plan      = "hourly"
}
`
}

func (r MarketplaceAgreementResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_marketplace_agreement" "import" {
  publisher = azurerm_marketplace_agreement.test.publisher
  offer     = azurerm_marketplace_agreement.test.offer
  plan      = azurerm_marketplace_agreement.test.plan
}
`, r.basicConfig(data))
}
