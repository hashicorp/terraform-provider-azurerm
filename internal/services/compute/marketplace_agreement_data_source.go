package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2021-01-01/agreements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceMarketplaceAgreement() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMarketplaceAgreementRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"offer": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"plan": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"publisher": {
				Type:         pluginsdk.TypeString,
				Required:     true,
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

func dataSourceMarketplaceAgreementRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := agreements.NewOfferPlanID(subscriptionId, d.Get("publisher").(string), d.Get("offer").(string), d.Get("plan").(string))

	log.Printf("[DEBUG] retrieving %s", id)

	term, err := client.MarketplaceAgreementsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(term.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	d.SetId(id.ID())
	if model := term.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("license_text_link", props.LicenseTextLink)
			d.Set("privacy_policy_link", props.PrivacyPolicyLink)
		}
	}
	return nil
}
