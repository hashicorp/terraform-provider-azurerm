// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_system_topic_event_subscription

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/eventsubscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventGridSystemTopicEventSubscriptionCreateUpdate,
		Read:   resourceEventGridSystemTopicEventSubscriptionRead,
		Update: resourceEventGridSystemTopicEventSubscriptionCreateUpdate,
		Delete: resourceEventGridSystemTopicEventSubscriptionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := eventsubscriptions.ParseSystemTopicEventSubscriptionID(id)
			return err
		}),

		Schema: eventgridSystemTopicEventSubscriptionSchema(),
	}
}
