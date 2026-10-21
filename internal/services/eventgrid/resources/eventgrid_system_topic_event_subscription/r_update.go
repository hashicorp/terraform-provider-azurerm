// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_system_topic_event_subscription

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/eventsubscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventGridSystemTopicEventSubscriptionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptions
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := eventsubscriptions.NewSystemTopicEventSubscriptionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("system_topic").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.SystemTopicEventSubscriptionsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_eventgrid_system_topic_event_subscription", id.ID())
			}
		}
	}

	destination := helpers.ExpandEventSubscriptionDestination(d)
	if destination == nil {
		return fmt.Errorf("one of the following endpoint types must be specified to create an EventGrid System Topic Event Subscription: %q", possibleSystemTopicEventSubscriptionEndpointTypes())
	}

	filter, err := helpers.ExpandEventSubscriptionFilter(d)
	if err != nil {
		return fmt.Errorf("expanding `filters`: %+v", err)
	}

	deadLetterDestination := helpers.ExpandEventSubscriptionStorageBlobDeadLetterDestination(d)

	eventSubscriptionProperties := eventsubscriptions.EventSubscriptionProperties{
		Filter:              filter,
		RetryPolicy:         helpers.ExpandEventSubscriptionRetryPolicy(d),
		Labels:              pluginsdk.ExpandStringSlice(d.Get("labels").([]any)),
		EventDeliverySchema: pointer.ToEnum[eventsubscriptions.EventDeliverySchema](d.Get("event_delivery_schema").(string)),
		ExpirationTimeUtc:   pointer.To(d.Get("expiration_time_utc").(string)),
	}

	if v, ok := d.GetOk("delivery_identity"); ok {
		deliveryIdentityRaw := v.([]any)
		deliveryIdentity, err := helpers.ExpandEventSubscriptionIdentity(deliveryIdentityRaw)
		if err != nil {
			return fmt.Errorf("expanding `delivery_identity`: %+v", err)
		}

		eventSubscriptionProperties.DeliveryWithResourceIdentity = &eventsubscriptions.DeliveryWithResourceIdentity{
			Identity:    deliveryIdentity,
			Destination: destination,
		}
	} else {
		eventSubscriptionProperties.Destination = destination
	}

	if v, ok := d.GetOk("dead_letter_identity"); ok {
		if deadLetterDestination == nil {
			return fmt.Errorf("`dead_letter_identity`: `storage_blob_dead_letter_destination` must be specified")
		}
		deadLetterIdentityRaw := v.([]any)
		deadLetterIdentity, err := helpers.ExpandEventSubscriptionIdentity(deadLetterIdentityRaw)
		if err != nil {
			return fmt.Errorf("expanding `dead_letter_identity`: %+v", err)
		}

		eventSubscriptionProperties.DeadLetterWithResourceIdentity = &eventsubscriptions.DeadLetterWithResourceIdentity{
			Identity:              deadLetterIdentity,
			DeadLetterDestination: deadLetterDestination,
		}
	} else {
		eventSubscriptionProperties.DeadLetterDestination = deadLetterDestination
	}

	eventSubscription := eventsubscriptions.EventSubscription{
		Properties: &eventSubscriptionProperties,
	}

	if d.IsNewResource() {
		if err := client.SystemTopicEventSubscriptionsCreateOrUpdateCallbackThenPoll(ctx, id, eventSubscription, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.SystemTopicEventSubscriptionsCreateOrUpdateThenPoll(ctx, id, eventSubscription); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceEventGridSystemTopicEventSubscriptionRead(d, meta)
}
