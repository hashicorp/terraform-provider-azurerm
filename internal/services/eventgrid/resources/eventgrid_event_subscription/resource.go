// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_event_subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/eventsubscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventGridEventSubscriptionCreateUpdate,
		Read:   resourceEventGridEventSubscriptionRead,
		Update: resourceEventGridEventSubscriptionCreateUpdate,
		Delete: resourceEventGridEventSubscriptionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
			if filterRaw := d.Get("advanced_filter"); len(filterRaw.([]any)) == 1 {
				filters := filterRaw.([]any)[0].(map[string]any)
				valueCount := 0
				for _, valRaw := range filters {
					for _, val := range valRaw.([]any) {
						v := val.(map[string]any)
						if values, ok := v["values"]; ok {
							valueCount += len(values.([]any))
						} else if _, ok := v["value"]; ok {
							valueCount++
						}
					}
				}
				if valueCount > 25 {
					return fmt.Errorf("the total number of `advanced_filter` values allowed on a single event subscription is 25, but %d are configured", valueCount)
				}
			}
			return nil
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := eventsubscriptions.ParseScopedEventSubscriptionID(id)
			return err
		}),

		Schema: eventgridEventSubscriptionSchema(),
	}
}
