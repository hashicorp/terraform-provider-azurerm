// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
	id, err := agreements.ParsePlanID(state.ID)
	if err != nil {
		return nil, err
	}

	agreementId := agreements.NewOfferPlanID(id.SubscriptionId, id.PublisherId, id.OfferId, id.PlanId)
	resp, err := clients.Compute.MarketplaceAgreementsClient.MarketplaceAgreementsGet(ctx, agreementId)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Marketplace Agreement %q", id)
	}

	if resp.Model == nil {
		return pointer.To(false), fmt.Errorf("retrieving %s, %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if accept := props.Accepted; accept != nil && *accept {
				return utils.Bool(true), nil
			}
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
		subscriptionId := clients.Account.SubscriptionId
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(15*time.Minute))
		defer cancel()

		idGet := agreements.NewOfferPlanID(subscriptionId, "barracudanetworks", offer, "hourly")
		idCancel := agreements.NewPlanID(subscriptionId, "barracudanetworks", offer, "hourly")
		existing, err := client.MarketplaceAgreementsGet(ctx, idGet)
		if err != nil {
			return err
		}

		if model := existing.Model; model != nil {
			if props := model.Properties; props != nil {
				if accepted := props.Accepted; accepted != nil && *accepted {
					resp, err := client.MarketplaceAgreementsCancel(ctx, idCancel)
					if err != nil {
						if response.WasNotFound(resp.HttpResponse) {
							return fmt.Errorf("marketplace agreement %q does not exist", idGet)
						}
						return fmt.Errorf("canceling %s: %+v", idGet, err)
					}
				}
			}
		}

		return nil
	}
}
