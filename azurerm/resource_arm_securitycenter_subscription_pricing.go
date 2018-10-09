package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/2017-08-01-preview/security"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSecurityCenterSubscriptionPricing() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSecurityCenterSubscriptionPricingCreateUpdate,
		Read:   resourceArmSecurityCenterSubscriptionPricingRead,
		Update: resourceArmSecurityCenterSubscriptionPricingCreateUpdate,
		Delete: resourceArmSecurityCenterSubscriptionPricingDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"tier": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.Free),
					string(security.Standard),
				}, false),
			},
		},
	}
}

//NOTE: seems default is the only valid pricing name:
//Code="InvalidInputJson" Message="Pricing name '360k Sponsored' is not allowed. Expected 'default' for this scope."
func resourceArmSecurityCenterSubscriptionPricingCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).securityCenterPricingClient
	ctx := meta.(*ArmClient).StopContext

	pricing := security.Pricing{
		PricingProperties: &security.PricingProperties{
			PricingTier: security.PricingTier(d.Get("tier").(string)),
		},
	}

	_, err := client.UpdateSubscriptionPricing(ctx, "default", pricing)
	if err != nil {
		return fmt.Errorf("Error creating/updating Security Center Subscription pricing: %+v", err)
	}

	resp, err := client.GetSubscriptionPricing(ctx, "default")
	if err != nil {
		return fmt.Errorf("Error reading Security Center Subscription pricing: %+v", err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Security Center Subscription pricing ID is nil")
	}

	d.SetId(*resp.ID)

	return resourceArmSecurityCenterSubscriptionPricingRead(d, meta)
}

func resourceArmSecurityCenterSubscriptionPricingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).securityCenterPricingClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.GetSubscriptionPricing(ctx, "default")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Security Center Subscription was not found: %v", err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Security Center Subscription pricing: %+v", err)
	}

	if properties := resp.PricingProperties; properties != nil {
		d.Set("tier", properties.PricingTier)
	}

	return nil
}

func resourceArmSecurityCenterSubscriptionPricingDelete(_ *schema.ResourceData, _ interface{}) error {
	return nil //cannot be deleted
}
