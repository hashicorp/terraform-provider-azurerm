// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package marketplace_agreement

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Uses this data source to access information about an existing Marketplace Agreement.",
		Read:        dataSourceMarketplaceAgreementRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"offer": {
				Description:  "The Offer of the Marketplace Image.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"plan": {
				Description:  "The Plan of the Marketplace Image.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"publisher": {
				Description:  "The Publisher of the Marketplace Image.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"accepted": {
				Description: "Whether the Marketplace Agreement has been accepted.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
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

func dataSourceMarketplaceAgreementRead(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).MarketplaceAgreementsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	// The Resource ID for this is the Plan ID, however we have to retrieve information about the signed plan
	id := agreements.NewPlanID(subscriptionId, d.Get("publisher").(string), d.Get("offer").(string), d.Get("plan").(string))

	getId := agreements.NewOfferPlanID(id.SubscriptionId, id.PublisherId, id.OfferId, id.PlanId)
	term, err := client.MarketplaceAgreementsGet(ctx, getId)
	if err != nil {
		if response.WasNotFound(term.HttpResponse) {
			return fmt.Errorf("%s was not found", getId)
		}

		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	d.SetId(id.ID())
	if model := term.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("license_text_link", props.LicenseTextLink)
			d.Set("privacy_policy_link", props.PrivacyPolicyLink)
			d.Set("accepted", pointer.From(props.Accepted))
		}
	}
	return nil
}
