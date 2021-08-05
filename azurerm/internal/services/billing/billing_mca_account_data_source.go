package billing

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func dataSourceBillingMCAAccountScope() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBillingMCAAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"billing_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"billing_profile_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"invoice_section_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func dataSourceBillingMCAAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	// (@jackofallops) - This is a helper Data Source until the Billing API is usable in the Azure SDK
	billingScopeMCAFmt := "/providers/Microsoft.Billing/billingAccounts/%s/billingProfiles/%s/invoiceSections/%s"

	d.SetId(fmt.Sprintf(billingScopeMCAFmt, d.Get("billing_account_name").(string), d.Get("billing_profile_name").(string), d.Get("invoice_section_name").(string)))
	return nil
}
