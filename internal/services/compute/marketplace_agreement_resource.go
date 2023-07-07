// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMarketplaceAgreement() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMarketplaceAgreementCreateUpdate,
		Read:   resourceMarketplaceAgreementRead,
		Delete: resourceMarketplaceAgreementDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := agreements.ParsePlanID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"offer": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"plan": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"publisher": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"license_text_link": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"privacy_policy_link": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMarketplaceAgreementCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := agreements.NewPlanID(subscriptionId, d.Get("publisher").(string), d.Get("offer").(string), d.Get("plan").(string))

	log.Printf("[DEBUG] retrieving %s", id)

	agreementId := agreements.NewOfferPlanID(id.SubscriptionId, id.PublisherId, id.OfferId, id.PlanId)
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

	terms.Properties.Accepted = utils.Bool(true)

	log.Printf("[DEBUG] Accepting the Marketplace Terms for %s", id)
	if _, err := client.MarketplaceAgreementsCreate(ctx, agreementId, *terms); err != nil {
		return fmt.Errorf("accepting Terms for %s: %s", id, err)
	}
	log.Printf("[DEBUG] Accepted the Marketplace Terms for %s", id)

	d.SetId(id.ID())

	return resourceMarketplaceAgreementRead(d, meta)
}

func resourceMarketplaceAgreementRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceMarketplaceAgreementDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agreements.ParsePlanID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.MarketplaceAgreementsCancel(ctx, *id); err != nil {
		return fmt.Errorf("cancelling agreement for %s: %s", *id, err)
	}

	return nil
}
