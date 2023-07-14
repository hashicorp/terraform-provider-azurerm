// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-11-01/frontdoor" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceCdnFrontDoorFirewallPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCdnFrontDoorFirewallPolicyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorFirewallPolicyName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"redirect_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"frontend_endpoint_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceCdnFrontDoorFirewallPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorLegacyFirewallPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFrontDoorFirewallPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	skuName := ""
	if sku := resp.Sku; sku != nil {
		skuName = string(sku.Name)
	}

	d.SetId(id.ID())
	d.Set("name", id.FrontDoorWebApplicationFirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("sku_name", skuName)

	if properties := resp.WebApplicationFirewallPolicyProperties; properties != nil {
		if policy := properties.PolicySettings; policy != nil {
			d.Set("enabled", policy.EnabledState == frontdoor.PolicyEnabledStateEnabled)
			d.Set("mode", string(policy.Mode))
			d.Set("redirect_url", policy.RedirectURL)
		}

		if err := d.Set("frontend_endpoint_ids", flattenFrontendEndpointLinkSlice(properties.FrontendEndpointLinks)); err != nil {
			return fmt.Errorf("flattening 'frontend_endpoint_ids': %+v", err)
		}
	}

	return nil
}
