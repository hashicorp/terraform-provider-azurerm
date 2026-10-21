// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_system_topic_event_subscription

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/eventsubscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventGridSystemTopicEventSubscriptionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptions
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := eventsubscriptions.ParseSystemTopicEventSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SystemTopicEventSubscriptionsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	fullUrlResp, err := client.SystemTopicEventSubscriptionsGetFullURL(ctx, *id)
	if err != nil {
		// unexpected status 400 with error: InvalidRequest: Destination type of the event subscription XXXX
		// is StorageQueue which doesn't support full URL. Only webhook destination type supports full URL.
		if !response.WasBadRequest(fullUrlResp.HttpResponse) {
			return fmt.Errorf("retrieving full url for %s: %+v", *id, err)
		}
	}

	d.Set("name", id.EventSubscriptionName)
	d.Set("system_topic", id.SystemTopicName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			expirationTimeUtc := ""
			if props.ExpirationTimeUtc != nil {
				t, err := props.GetExpirationTimeUtcAsTime()
				if err == nil {
					expirationTimeUtc = t.Format(time.RFC3339)
				}
			}
			d.Set("expiration_time_utc", expirationTimeUtc)
			d.Set("event_delivery_schema", pointer.FromEnum(props.EventDeliverySchema))

			destination := props.Destination
			deliveryIdentityFlattened := make([]any, 0)
			if deliveryIdentity := props.DeliveryWithResourceIdentity; deliveryIdentity != nil {
				destination = deliveryIdentity.Destination
				deliveryIdentityFlattened = helpers.FlattenEventSubscriptionIdentity(deliveryIdentity.Identity)
			}
			if err := d.Set("delivery_identity", deliveryIdentityFlattened); err != nil {
				return fmt.Errorf("setting `delivery_identity`: %+v", err)
			}

			existingMappingsFromState := helpers.ExpandEventSubscriptionDeliveryAttributeMappings(d.Get("delivery_property").([]any))
			if err := d.Set("delivery_property", helpers.FlattenEventSubscriptionDeliveryAttributeMappings(destination, existingMappingsFromState)); err != nil {
				return fmt.Errorf("setting `delivery_property` for %s: %+v", *id, err)
			}

			if err := d.Set("azure_function_endpoint", helpers.FlattenEventSubscriptionDestinationAzureFunction(destination)); err != nil {
				return fmt.Errorf("setting `azure_function_endpoint` for %s: %+v", *id, err)
			}

			d.Set("eventhub_id", helpers.FlattenEventSubscriptionDestinationEventHub(destination))
			d.Set("hybrid_connection_id", helpers.FlattenEventSubscriptionDestinationHybridConnection(destination))
			d.Set("service_bus_queue_id", helpers.FlattenEventSubscriptionDestinationServiceBusQueueEndpoint(destination))
			d.Set("service_bus_topic_id", helpers.FlattenEventSubscriptionDestinationServiceBusTopicEndpoint(destination))
			if err := d.Set("storage_queue_endpoint", helpers.FlattenEventSubscriptionDestinationStorageQueueEndpoint(destination)); err != nil {
				return fmt.Errorf("setting `storage_queue_endpoint` for %s: %+v", *id, err)
			}
			if err := d.Set("webhook_endpoint", helpers.FlattenEventSubscriptionWebhookEndpoint(destination, fullUrlResp.Model)); err != nil {
				return fmt.Errorf("setting `webhook_endpoint` for %s: %+v", *id, err)
			}

			deadLetterDestination := props.DeadLetterDestination
			deadLetterIdentityFlattened := make([]any, 0)
			if deadLetterIdentity := props.DeadLetterWithResourceIdentity; deadLetterIdentity != nil {
				deadLetterDestination = deadLetterIdentity.DeadLetterDestination
				deadLetterIdentityFlattened = helpers.FlattenEventSubscriptionIdentity(deadLetterIdentity.Identity)
			}
			if err := d.Set("dead_letter_identity", deadLetterIdentityFlattened); err != nil {
				return fmt.Errorf("setting `dead_letter_identity`: %+v", err)
			}
			if err := d.Set("storage_blob_dead_letter_destination", helpers.FlattenEventSubscriptionStorageBlobDeadLetterDestination(deadLetterDestination)); err != nil {
				return fmt.Errorf("setting `storage_blob_dead_letter_destination`: %+v", err)
			}

			enableAdvancedFilteringOnArrays := false
			includedEventTypes := make([]string, 0)
			if filter := props.Filter; filter != nil {
				enableAdvancedFilteringOnArrays = pointer.From(filter.EnableAdvancedFilteringOnArrays)
				includedEventTypes = pointer.From(filter.IncludedEventTypes)
			}
			d.Set("advanced_filtering_on_arrays_enabled", enableAdvancedFilteringOnArrays)
			d.Set("included_event_types", includedEventTypes)
			if err := d.Set("advanced_filter", helpers.FlattenEventSubscriptionAdvancedFilter(props.Filter)); err != nil {
				return fmt.Errorf("setting `advanced_filter` for %s: %+v", *id, err)
			}
			if err := d.Set("retry_policy", helpers.FlattenEventSubscriptionRetryPolicy(props.RetryPolicy)); err != nil {
				return fmt.Errorf("setting `retry_policy` for %s: %+v", *id, err)
			}
			if err := d.Set("subject_filter", helpers.FlattenEventSubscriptionSubjectFilter(props.Filter)); err != nil {
				return fmt.Errorf("setting `subject_filter` for %s: %+v", *id, err)
			}

			d.Set("labels", props.Labels)
		}
	}

	return nil
}

func possibleSystemTopicEventSubscriptionEndpointTypes() []string {
	return []string{
		string(helpers.AzureFunctionEndpoint),
		string(helpers.EventHubID),
		string(helpers.HybridConnectionID),
		string(helpers.ServiceBusQueueID),
		string(helpers.ServiceBusTopicID),
		string(helpers.StorageQueueEndpoint),
		string(helpers.WebHookEndpoint),
	}
}
