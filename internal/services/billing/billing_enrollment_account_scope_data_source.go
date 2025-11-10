// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataSourceBillingEnrollmentAccountScope() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBillingEnrollemntAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"billing_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"enrollment_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func dataSourceBillingEnrollemntAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	// (@jackofallops) - This is a helper Data Source until the Billing API is usable in the Azure SDK
	billingScopeEnrollmentFmt := "/providers/Microsoft.Billing/billingAccounts/%s/enrollmentAccounts/%s"

	d.SetId(fmt.Sprintf(billingScopeEnrollmentFmt, d.Get("billing_account_name").(string), d.Get("enrollment_account_name").(string)))
	return nil
}
