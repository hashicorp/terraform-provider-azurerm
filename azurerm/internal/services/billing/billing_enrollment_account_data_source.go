package billing

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceBillingEnrollmentAccountScope() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBillingEnrollemntAccountRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"billing_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"enrollment_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func dataSourceBillingEnrollemntAccountRead(d *schema.ResourceData, meta interface{}) error {
	// (@jackofallops) - This is a helper Data Source until the Billing API is usable in the Azure SDK
	billingScopeEnrollmentFmt := "/providers/Microsoft.Billing/billingAccounts/%s/enrollmentAccounts/%s"

	d.SetId(fmt.Sprintf(billingScopeEnrollmentFmt, d.Get("billing_account_name").(string), d.Get("enrollment_account_name").(string)))
	return nil
}
