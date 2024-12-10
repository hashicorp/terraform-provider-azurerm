// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2024-02-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
	client := meta.(*clients.Client).Cdn.FrontDoorFirewallPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := waf.NewFrontDoorWebApplicationFirewallPolicyID(subscriptionId, resourceGroup, name)

	result, err := client.PoliciesGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(result.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := result.Model

	if model == nil {
		return fmt.Errorf("retrieving %s: 'model' was nil", id)
	}

	if model.Sku == nil {
		return fmt.Errorf("retrieving %s: 'model.Sku' was nil", id)
	}

	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: 'model.Properties' was nil", id)
	}

	props := model.Properties

	skuName := ""
	if sku := model.Sku; sku != nil {
		skuName = string(pointer.From(model.Sku.Name))
	}

	d.SetId(id.ID())
	d.Set("name", id.FrontDoorWebApplicationFirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("sku_name", skuName)

	if policy := props.PolicySettings; policy != nil {
		d.Set("enabled", pointer.From(policy.EnabledState) == waf.PolicyEnabledStateEnabled)
		d.Set("mode", pointer.From(policy.Mode))
		d.Set("redirect_url", policy.RedirectURL)
	}

	if err := d.Set("frontend_endpoint_ids", flattenFrontendEndpointLinkSlice(props.FrontendEndpointLinks)); err != nil {
		return fmt.Errorf("flattening 'frontend_endpoint_ids': %+v", err)
	}

	return nil
}
