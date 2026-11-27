// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package marketplace_agreement

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMarketplaceAgreementCreate(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	subscriptionId := account.From(meta).SubscriptionId
	defer cancel()

	id := agreements.NewPlanID(subscriptionId, d.Get("publisher").(string), d.Get("offer").(string), d.Get("plan").(string))

	agreementId := agreements.NewOfferPlanID(id.SubscriptionId, id.PublisherId, id.OfferId, id.PlanId)

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		term, err := client.MarketplaceAgreementsGet(ctx, agreementId)
		if err != nil {
			if !response.WasNotFound(term.HttpResponse) {
				return fmt.Errorf("retrieving %s: %s", id, err)
			}
		}

		accepted := false
		if model := term.Model; model != nil {
			if props := model.Properties; props != nil {
				if acc := props.Accepted; acc != nil {
					accepted = *acc
				}
			}
		}
		if accepted {
			return tf.ImportAsExistsError("azurerm_marketplace_agreement", id.ID())
		}
	}

	resp, err := client.MarketplaceAgreementsGet(ctx, agreementId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: Model was nil", id)
	}

	terms := resp.Model
	if terms.Properties == nil {
		return fmt.Errorf("retrieving %s: AgreementProperties was nil", id)
	}

	terms.Properties.Accepted = pointer.To(true)

	log.Printf("[DEBUG] Accepting the Marketplace Terms for %s", id)
	if _, err := client.MarketplaceAgreementsCreate(ctx, agreementId, *terms); err != nil {
		return fmt.Errorf("accepting Terms for %s: %s", id, err)
	}
	log.Printf("[DEBUG] Accepted the Marketplace Terms for %s", id)

	d.SetId(id.ID())

	return resourceMarketplaceAgreementRead(d, meta)
}
