package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MarketplaceAgreementResource struct{}

func TestAccMarketplaceAgreement_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_agreement", "test")
	r := MarketplaceAgreementResource{}
	offer := "waf"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientWithoutResource(r.cancelExistingAgreement(offer)),
			),
		},
		{
			Config: r.basic(offer),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("license_text_link").Exists(),
				check.That(data.ResourceName).Key("privacy_policy_link").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMarketplaceAgreement_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_agreement", "test")
	r := MarketplaceAgreementResource{}
	offer := "barracuda-ng-firewall"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientWithoutResource(r.cancelExistingAgreement(offer)),
			),
		},
		{
			Config: r.basic(offer),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(offer),
			ExpectError: acceptance.RequiresImportError("azurerm_marketplace_agreement"),
		},
	})
}

func (t MarketplaceAgreementResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PlanID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.MarketplaceAgreementsClient.Get(ctx, id.AgreementName, id.OfferName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Marketplace Agreement %q", id)
	}

	if resp.ID == nil {
		return utils.Bool(false), nil
	}

	if props := resp.AgreementProperties; props != nil {
		if accept := props.Accepted; accept != nil && *accept {
			return utils.Bool(true), nil
		}
	}

	return utils.Bool(false), nil
}

func (MarketplaceAgreementResource) basic(offer string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_marketplace_agreement" "test" {
  publisher = "barracudanetworks"
  offer     = "%s"
  plan      = "hourly"
}
`, offer)
}

func (r MarketplaceAgreementResource) requiresImport(offer string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_marketplace_agreement" "import" {
  publisher = azurerm_marketplace_agreement.test.publisher
  offer     = azurerm_marketplace_agreement.test.offer
  plan      = azurerm_marketplace_agreement.test.plan
}
`, r.basic(offer))
}

func (MarketplaceAgreementResource) empty() string {
	return `
provider "azurerm" {
  features {}
}
`
}

func (r MarketplaceAgreementResource) cancelExistingAgreement(offer string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		client := clients.Compute.MarketplaceAgreementsClient
		id := parse.NewPlanID(client.SubscriptionID, "barracudanetworks", offer, "hourly")

		existing, err := client.Get(ctx, id.AgreementName, id.OfferName, id.Name)
		if err != nil {
			return err
		}

		if props := existing.AgreementProperties; props != nil {
			if accepted := props.Accepted; accepted != nil && *accepted {
				resp, err := client.Cancel(ctx, id.AgreementName, id.OfferName, id.Name)
				if err != nil {
					if utils.ResponseWasNotFound(resp.Response) {
						return fmt.Errorf("marketplace agreement %q does not exist", id)
					}
					return fmt.Errorf("canceling Marketplace Agreement : %+v", err)
				}
			}
		}

		return nil
	}
}
