// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceSubscription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSubscriptionRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*pluginsdk.Schema{
			"subscription_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenant_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location_placement_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"quota_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"spending_limit": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client)
	groupClient := client.Subscription.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId := d.Get("subscription_id").(string)
	if subscriptionId == "" {
		subscriptionId = client.Account.SubscriptionId
	}

	id := commonids.NewSubscriptionID(subscriptionId)
	resp, err := groupClient.Get(ctx, subscriptionId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("subscription_id", resp.SubscriptionID)
	d.Set("display_name", resp.DisplayName)
	d.Set("tenant_id", resp.TenantID)
	d.Set("state", resp.State)
	if props := resp.SubscriptionPolicies; props != nil {
		d.Set("location_placement_id", props.LocationPlacementID)
		d.Set("quota_id", props.QuotaID)
		d.Set("spending_limit", props.SpendingLimit)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
