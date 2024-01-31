// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSubscriptions() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSubscriptionsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"display_name_prefix": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
			"display_name_contains": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
			"subscriptions": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"subscription_id": {
							Type:     pluginsdk.TypeString,
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
				},
			},
		},
	}
}

func dataSourceSubscriptionsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	armClient := meta.(*clients.Client)
	subClient := armClient.Subscription.SubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	displayNamePrefix := strings.ToLower(d.Get("display_name_prefix").(string))
	displayNameContains := strings.ToLower(d.Get("display_name_contains").(string))

	// ListComplete returns an iterator struct
	results, err := subClient.ListComplete(ctx)
	if err != nil {
		return fmt.Errorf("listing subscriptions: %+v", err)
	}

	subscriptions := make([]interface{}, 0)
	for _, item := range results.Items {
		// check if the display name prefix matches the given input
		if displayNamePrefix != "" {
			if item.DisplayName == nil || !strings.HasPrefix(strings.ToLower(*item.DisplayName), displayNamePrefix) {
				// the display name does not match the given prefix
				continue
			}
		}
		// check if the display name matches the 'contains' comparison
		if displayNameContains != "" {
			if item.DisplayName == nil || !strings.Contains(strings.ToLower(*item.DisplayName), displayNameContains) {
				// the display name does not match the contains check
				continue
			}
		}

		quotaId := ""
		locationPlacementId := ""
		spendingLimit := ""
		if policies := item.SubscriptionPolicies; policies != nil {
			locationPlacementId = pointer.From(policies.LocationPlacementId)
			quotaId = pointer.From(policies.QuotaId)
			if policies.SpendingLimit != nil {
				spendingLimit = string(*policies.SpendingLimit)
			}
		}

		subscriptions = append(subscriptions, map[string]interface{}{
			"display_name":          pointer.From(item.DisplayName),
			"id":                    pointer.From(item.Id),
			"location_placement_id": locationPlacementId,
			"quota_id":              quotaId,
			"spending_limit":        spendingLimit,
			"state":                 string(pointer.From(item.State)),
			"subscription_id":       pointer.From(item.SubscriptionId),
			"tags":                  tags.Flatten(item.Tags),
			"tenant_id":             pointer.From(item.TenantId),
		})
	}

	d.SetId("subscriptions-" + armClient.Account.TenantId)
	if err = d.Set("subscriptions", subscriptions); err != nil {
		return fmt.Errorf("setting `subscriptions`: %+v", err)
	}

	return nil
}
