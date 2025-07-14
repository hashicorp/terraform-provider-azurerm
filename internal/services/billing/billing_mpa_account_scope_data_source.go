// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataSourceBillingMPAAccountScope() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBillingMPAAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"billing_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"customer_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func dataSourceBillingMPAAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	// (@jackofallops) - This is a helper Data Source until the Billing API is usable in the Azure SDK
	billingScopeMPAFmt := "/providers/Microsoft.Billing/billingAccounts/%s/customers/%s"

	d.SetId(fmt.Sprintf(billingScopeMPAFmt, d.Get("billing_account_name").(string), d.Get("customer_name").(string)))
	return nil
}
