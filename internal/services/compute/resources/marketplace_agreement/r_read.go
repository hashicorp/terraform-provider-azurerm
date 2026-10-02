// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package marketplace_agreement

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMarketplaceAgreementRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agreements.ParsePlanID(d.Id())
	if err != nil {
		return err
	}

	agreementId := agreements.NewOfferPlanID(id.SubscriptionId, id.PublisherId, id.OfferId, id.PlanId)
	term, err := client.MarketplaceAgreementsGet(ctx, agreementId)
	if err != nil {
		if response.WasNotFound(term.HttpResponse) {
			log.Printf("[DEBUG] The Marketplace Terms was not found for %s", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving the Marketplace Terms for %s: %s", id, err)
	}

	d.Set("publisher", id.PublisherId)
	d.Set("offer", id.OfferId)
	d.Set("plan", id.PlanId)

	if model := term.Model; model != nil {
		if props := model.Properties; props != nil {
			if accepted := props.Accepted != nil && *props.Accepted; !accepted {
				// if props.Accepted is not true, the agreement does not exist
				d.SetId("")
			}
			d.Set("license_text_link", props.LicenseTextLink)
			d.Set("privacy_policy_link", props.PrivacyPolicyLink)
		}
	}
	return nil
}
