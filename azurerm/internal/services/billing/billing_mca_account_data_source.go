package billing

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceBillingMCAAccountScope() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBillingMCAAccountRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"billing_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"billing_profile_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"invoice_section_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func dataSourceBillingMCAAccountRead(d *schema.ResourceData, meta interface{}) error {
	// (@jackofallops) - This is a helper Data Source until the Billing API is usable in the Azure SDK
	billingScopeMCAFmt := "/providers/Microsoft.Billing/billingAccounts/%s/billingProfiles/%s/invoiceSections/%s"

	d.SetId(fmt.Sprintf(billingScopeMCAFmt, d.Get("billing_account_name").(string), d.Get("billing_profile_name").(string), d.Get("invoice_section_name").(string)))
	return nil
}
