package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

//NOTE: seems default is the only valid pricing name:
//Code="InvalidInputJson" Message="Pricing name 'kt's price' is not allowed. Expected 'default' for this scope."
const securityCenterSubscriptionPricingName = "default"

func resourceArmSecurityCenterSubscriptionPricing() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSecurityCenterSubscriptionPricingUpdate,
		Read:   resourceArmSecurityCenterSubscriptionPricingRead,
		Update: resourceArmSecurityCenterSubscriptionPricingUpdate,
		Delete: resourceArmSecurityCenterSubscriptionPricingDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"tier": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.Free),
					string(security.Standard),
				}, false),
			},
		},
	}
}

func resourceArmSecurityCenterSubscriptionPricingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).SecurityCenter.PricingClient
	ctx := meta.(*ArmClient).StopContext

	name := securityCenterSubscriptionPricingName

	// not doing import check as afaik it always exists (cannot be deleted)
	// all this resource does is flip a boolean

	pricing := security.Pricing{
		PricingProperties: &security.PricingProperties{
			PricingTier: security.PricingTier(d.Get("tier").(string)),
		},
	}

	if _, err := client.UpdateSubscriptionPricing(ctx, name, pricing); err != nil {
		return fmt.Errorf("Error creating/updating Security Center Subscription pricing: %+v", err)
	}

	resp, err := client.GetSubscriptionPricing(ctx, name)
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
	client := meta.(*ArmClient).SecurityCenter.PricingClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.GetSubscriptionPricing(ctx, securityCenterSubscriptionPricingName)
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
	log.Printf("[DEBUG] Security Center Subscription deletion invocation")
	return nil //cannot be deleted.
}
