// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client)
	groupClient := client.Subscription.SubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId := d.Get("subscription_id").(string)
	if subscriptionId == "" {
		subscriptionId = client.Account.SubscriptionId
	}

	id := commonids.NewSubscriptionID(subscriptionId)
	resp, err := groupClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("subscription_id", model.SubscriptionId)
		d.Set("display_name", model.DisplayName)
		d.Set("tenant_id", model.TenantId)
		d.Set("state", string(pointer.From(model.State)))
		if props := model.SubscriptionPolicies; props != nil {
			d.Set("location_placement_id", props.LocationPlacementId)
			d.Set("quota_id", props.QuotaId)
			d.Set("spending_limit", string(pointer.From(props.SpendingLimit)))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
