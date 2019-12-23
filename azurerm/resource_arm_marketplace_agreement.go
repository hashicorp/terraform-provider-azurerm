package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMarketplaceAgreement() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMarketplaceAgreementCreateUpdate,
		Read:   resourceArmMarketplaceAgreementRead,
		Delete: resourceArmMarketplaceAgreementDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"offer": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"plan": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"publisher": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"license_text_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"privacy_policy_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmMarketplaceAgreementCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	offer := d.Get("offer").(string)
	plan := d.Get("plan").(string)
	publisher := d.Get("publisher").(string)

	log.Printf("[DEBUG] Retrieving the Marketplace Terms for Publisher %q / Offer %q / Plan %q", publisher, offer, plan)

	if features.ShouldResourcesBeImported() {
		agreement, err := client.Get(ctx, publisher, offer, plan)
		if err != nil {
			if !utils.ResponseWasNotFound(agreement.Response) {
				return fmt.Errorf("Error retrieving agreement for Publisher %q / Offer %q / Plan %q: %s", publisher, offer, plan, err)
			}
		}

		accepted := false
		if props := agreement.AgreementProperties; props != nil {
			if acc := props.Accepted; acc != nil {
				accepted = *acc
			}
		}

		if accepted {
			return tf.ImportAsExistsError("azurerm_marketplace_agreement", *agreement.ID)
		}
	}

	terms, err := client.Get(ctx, publisher, offer, plan)
	if err != nil {
		return fmt.Errorf("Error retrieving agreement for Publisher %q / Offer %q / Plan %q: %s", publisher, offer, plan, err)
	}
	if terms.AgreementProperties == nil {
		return fmt.Errorf("Error retrieving agreement for Publisher %q / Offer %q / Plan %q: AgreementProperties was nil", publisher, offer, plan)
	}

	terms.AgreementProperties.Accepted = utils.Bool(true)

	log.Printf("[DEBUG] Accepting the Marketplace Terms for Publisher %q / Offer %q / Plan %q", publisher, offer, plan)
	if _, err := client.Create(ctx, publisher, offer, plan, terms); err != nil {
		return fmt.Errorf("Error accepting Terms for Publisher %q / Offer %q / Plan %q: %s", publisher, offer, plan, err)
	}
	log.Printf("[DEBUG] Accepted the Marketplace Terms for Publisher %q / Offer %q / Plan %q", publisher, offer, plan)

	agreement, err := client.GetAgreement(ctx, publisher, offer, plan)
	if err != nil {
		return fmt.Errorf("Error retrieving agreement for Publisher %q / Offer %q / Plan %q: %s", publisher, offer, plan, err)
	}

	d.SetId(*agreement.ID)

	return resourceArmMarketplaceAgreementRead(d, meta)
}

func resourceArmMarketplaceAgreementRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	publisher := id.Path["agreements"]
	offer := id.Path["offers"]
	plan := id.Path["plans"]

	resp, err := client.Get(ctx, publisher, offer, plan)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Agreement was not found for Publisher %q / Offer %q / Plan %q", publisher, offer, plan)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving agreement for Publisher %q / Offer %q / Plan %q: %s", publisher, offer, plan, err)
	}

	d.Set("publisher", publisher)
	d.Set("offer", offer)
	d.Set("plan", plan)

	if props := resp.AgreementProperties; props != nil {
		d.Set("license_text_link", props.LicenseTextLink)
		d.Set("privacy_policy_link", props.PrivacyPolicyLink)
	}

	return nil
}

func resourceArmMarketplaceAgreementDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	publisher := id.Path["agreements"]
	offer := id.Path["offers"]
	plan := id.Path["plans"]

	if _, err := client.Cancel(ctx, publisher, offer, plan); err != nil {
		return fmt.Errorf("Error cancelling agreement for Publisher %q / Offer %q / Plan %q: %s", publisher, offer, plan, err)
	}

	return nil
}
