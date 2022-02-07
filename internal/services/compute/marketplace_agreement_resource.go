package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
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
			_, err := parse.PlanID(id)
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

	id := parse.NewPlanID(subscriptionId, d.Get("publisher").(string), d.Get("offer").(string), d.Get("plan").(string))

	log.Printf("[DEBUG] retrieving %s", id)

	term, err := client.Get(ctx, id.AgreementName, id.OfferName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(term.Response) {
			return fmt.Errorf("retrieving %s: %s", id, err)
		}
	}

	accepted := false
	if props := term.AgreementProperties; props != nil {
		if acc := props.Accepted; acc != nil {
			accepted = *acc
		}
	}

	if accepted {
		agreement, err := client.GetAgreement(ctx, id.AgreementName, id.OfferName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(agreement.Response) {
				return fmt.Errorf("retrieving %s: %s", id, err)
			}
		}
		return tf.ImportAsExistsError("azurerm_marketplace_agreement", id.ID())
	}

	terms, err := client.Get(ctx, id.AgreementName, id.OfferName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %s", id, err)
	}
	if terms.AgreementProperties == nil {
		return fmt.Errorf("retrieving %s: AgreementProperties was nil", id)
	}

	terms.AgreementProperties.Accepted = utils.Bool(true)

	log.Printf("[DEBUG] Accepting the Marketplace Terms for %s", id)
	if _, err := client.Create(ctx, id.AgreementName, id.OfferName, id.Name, terms); err != nil {
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

	id, err := parse.PlanID(d.Id())
	if err != nil {
		return err
	}

	term, err := client.Get(ctx, id.AgreementName, id.OfferName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(term.Response) {
			log.Printf("[DEBUG] The Marketplace Terms was not found for Publisher %q / Offer %q / Plan %q", id.AgreementName, id.OfferName, id.Name)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving the Marketplace Terms for Publisher %q / Offer %q / Plan %q: %s", id.AgreementName, id.OfferName, id.Name, err)
	}

	d.Set("publisher", id.AgreementName)
	d.Set("offer", id.OfferName)
	d.Set("plan", id.Name)

	if props := term.AgreementProperties; props != nil {
		if accepted := props.Accepted != nil && *props.Accepted; !accepted {
			// if props.Accepted is not true, the agreement does not exist
			d.SetId("")
		}
		d.Set("license_text_link", props.LicenseTextLink)
		d.Set("privacy_policy_link", props.PrivacyPolicyLink)
	}

	return nil
}

func resourceMarketplaceAgreementDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PlanID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Cancel(ctx, id.AgreementName, id.OfferName, id.Name); err != nil {
		return fmt.Errorf("cancelling agreement for Publisher %q / Offer %q / Plan %q: %s", id.AgreementName, id.OfferName, id.Name, err)
	}

	return nil
}
