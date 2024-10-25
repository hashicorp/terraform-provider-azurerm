// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/eventsubscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func possibleEventSubscriptionEndpointTypes() []string {
	return []string{
		string(AzureFunctionEndpoint),
		string(EventHubEndpointID),
		string(HybridConnectionEndpointID),
		string(ServiceBusQueueEndpointID),
		string(ServiceBusTopicEndpointID),
		string(StorageQueueEndpoint),
		string(WebHookEndpoint),
	}
}

func resourceEventGridEventSubscription() *pluginsdk.Resource {
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

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(_ context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
			if filterRaw := d.Get("advanced_filter"); len(filterRaw.([]interface{})) == 1 {
				filters := filterRaw.([]interface{})[0].(map[string]interface{})
				valueCount := 0
				for _, valRaw := range filters {
					for _, val := range valRaw.([]interface{}) {
						v := val.(map[string]interface{})
						if values, ok := v["values"]; ok {
							valueCount += len(values.([]interface{}))
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

		Schema: map[string]*pluginsdk.Schema{
			"name": eventSubscriptionSchemaEventSubscriptionName(),

			"scope": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"event_delivery_schema": eventSubscriptionSchemaEventDeliverySchema(),

			"expiration_time_utc": eventSubscriptionSchemaExpirationTimeUTC(),

			"azure_function_endpoint": eventSubscriptionSchemaAzureFunctionEndpoint(
				utils.RemoveFromStringArray(
					possibleEventSubscriptionEndpointTypes(),
					string(AzureFunctionEndpoint),
				),
			),

			"eventhub_endpoint_id": eventSubscriptionSchemaEventHubEndpointID(
				utils.RemoveFromStringArray(
					possibleEventSubscriptionEndpointTypes(),
					string(EventHubEndpointID),
				),
			),

			"hybrid_connection_endpoint_id": eventSubscriptionSchemaHybridConnectionEndpointID(
				utils.RemoveFromStringArray(
					possibleEventSubscriptionEndpointTypes(),
					string(HybridConnectionEndpointID),
				),
			),

			// TODO: this can become `service_bus_queue_id` in 4.0
			"service_bus_queue_endpoint_id": eventSubscriptionSchemaServiceBusQueueEndpointID(
				utils.RemoveFromStringArray(
					possibleEventSubscriptionEndpointTypes(),
					string(ServiceBusQueueEndpointID),
				),
			),

			// TODO: this can become `service_bus_topic_id` in 4.0
			"service_bus_topic_endpoint_id": eventSubscriptionSchemaServiceBusTopicEndpointID(
				utils.RemoveFromStringArray(
					possibleEventSubscriptionEndpointTypes(),
					string(ServiceBusTopicEndpointID),
				),
			),

			"storage_queue_endpoint": eventSubscriptionSchemaStorageQueueEndpoint(
				utils.RemoveFromStringArray(
					possibleEventSubscriptionEndpointTypes(),
					string(StorageQueueEndpoint),
				),
			),

			"webhook_endpoint": eventSubscriptionSchemaWebHookEndpoint(
				utils.RemoveFromStringArray(
					possibleEventSubscriptionEndpointTypes(),
					string(WebHookEndpoint),
				),
			),

			"included_event_types": eventSubscriptionSchemaIncludedEventTypes(),

			"subject_filter": eventSubscriptionSchemaSubjectFilter(),

			"advanced_filter": eventSubscriptionSchemaAdvancedFilter(),

			"delivery_identity": eventSubscriptionSchemaIdentity(),

			"dead_letter_identity": eventSubscriptionSchemaIdentity(),

			"storage_blob_dead_letter_destination": eventSubscriptionSchemaStorageBlobDeadletterDestination(),

			"retry_policy": eventSubscriptionSchemaRetryPolicy(),

			"labels": eventSubscriptionSchemaLabels(),

			"advanced_filtering_on_arrays_enabled": eventSubscriptionSchemaEnableAdvancedFilteringOnArrays(),

			"delivery_property": eventSubscriptionSchemaDeliveryProperty(),
		},
	}
}

func resourceEventGridEventSubscriptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptions
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := eventsubscriptions.NewScopedEventSubscriptionID(d.Get("scope").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_eventgrid_event_subscription", id.ID())
		}
	}

	destination := expandEventSubscriptionDestination(d)
	if destination == nil {
		return fmt.Errorf("one of the following endpoint types must be specificed to create an EventGrid Event Subscription: %q", possibleEventSubscriptionEndpointTypes())
	}

	filter, err := expandEventSubscriptionFilter(d)
	if err != nil {
		return fmt.Errorf("expanding filters for %s: %+v", id, err)
	}

	deadLetterDestination := expandEventSubscriptionStorageBlobDeadLetterDestination(d)

	properties := eventsubscriptions.EventSubscriptionProperties{
		ExpirationTimeUtc:   pointer.To(d.Get("expiration_time_utc").(string)),
		EventDeliverySchema: pointer.To(eventsubscriptions.EventDeliverySchema(d.Get("event_delivery_schema").(string))),
		Filter:              filter,
		Labels:              utils.ExpandStringSlice(d.Get("labels").([]interface{})),
		RetryPolicy:         expandEventSubscriptionRetryPolicy(d),
	}

	if v, ok := d.GetOk("delivery_identity"); ok {
		deliveryIdentityRaw := v.([]interface{})
		deliveryIdentity, err := expandEventSubscriptionIdentity(deliveryIdentityRaw)
		if err != nil {
			return fmt.Errorf("expanding `delivery_identity`: %+v", err)
		}

		properties.DeliveryWithResourceIdentity = &eventsubscriptions.DeliveryWithResourceIdentity{
			Identity:    deliveryIdentity,
			Destination: destination,
		}
	} else {
		properties.Destination = destination
	}

	if v, ok := d.GetOk("dead_letter_identity"); ok {
		if deadLetterDestination == nil {
			return fmt.Errorf("`dead_letter_identity`: `storage_blob_dead_letter_destination` must be specified")
		}
		deadLetterIdentityRaw := v.([]interface{})
		deadLetterIdentity, err := expandEventSubscriptionIdentity(deadLetterIdentityRaw)
		if err != nil {
			return fmt.Errorf("expanding `dead_letter_identity`: %+v", err)
		}

		properties.DeadLetterWithResourceIdentity = &eventsubscriptions.DeadLetterWithResourceIdentity{
			Identity:              deadLetterIdentity,
			DeadLetterDestination: deadLetterDestination,
		}
	} else {
		properties.DeadLetterDestination = deadLetterDestination
	}

	payload := eventsubscriptions.EventSubscription{
		Properties: &properties,
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceEventGridEventSubscriptionRead(d, meta)
}

func resourceEventGridEventSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptions
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := eventsubscriptions.ParseScopedEventSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	fullUrlResp, err := client.GetFullURL(ctx, *id)
	if err != nil {
		// unexpected status 400 with error: InvalidRequest: Destination type of the event subscription XXXX
		// is StorageQueue which doesn't support full URL. Only webhook destination type supports full URL.
		if !response.WasBadRequest(fullUrlResp.HttpResponse) {
			return fmt.Errorf("retrieving full url for %s: %+v", *id, err)
		}
	}

	d.Set("name", id.EventSubscriptionName)
	d.Set("scope", id.Scope)

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

			eventDeliverySchema := ""
			if props.EventDeliverySchema != nil {
				eventDeliverySchema = string(*props.EventDeliverySchema)
			}
			d.Set("event_delivery_schema", eventDeliverySchema)

			destination := props.Destination
			deliveryIdentityFlattened := make([]interface{}, 0)
			if deliveryIdentity := props.DeliveryWithResourceIdentity; deliveryIdentity != nil {
				destination = deliveryIdentity.Destination
				deliveryIdentityFlattened = flattenEventSubscriptionIdentity(deliveryIdentity.Identity)
			}
			if err := d.Set("delivery_identity", deliveryIdentityFlattened); err != nil {
				return fmt.Errorf("setting `delivery_identity` for %s: %+v", *id, err)
			}

			existingMappingsFromState := expandEventSubscriptionDeliveryAttributeMappings(d.Get("delivery_property").([]interface{}))
			deliveryMappings := flattenEventSubscriptionDeliveryAttributeMappings(destination, existingMappingsFromState)
			if err := d.Set("delivery_property", deliveryMappings); err != nil {
				return fmt.Errorf("setting `delivery_property` for %s: %+v", *id, err)
			}

			if err := d.Set("azure_function_endpoint", flattenEventSubscriptionDestinationAzureFunction(destination)); err != nil {
				return fmt.Errorf("setting `azure_function_endpoint` for %s: %+v", *id, err)
			}

			d.Set("eventhub_endpoint_id", flattenEventSubscriptionDestinationEventHub(destination))
			d.Set("hybrid_connection_endpoint_id", flattenEventSubscriptionDestinationHybridConnection(destination))
			d.Set("service_bus_queue_endpoint_id", flattenEventSubscriptionDestinationServiceBusQueueEndpoint(destination))
			d.Set("service_bus_topic_endpoint_id", flattenEventSubscriptionDestinationServiceBusTopicEndpoint(destination))
			if err := d.Set("storage_queue_endpoint", flattenEventSubscriptionDestinationStorageQueueEndpoint(destination)); err != nil {
				return fmt.Errorf("setting `storage_queue_endpoint` for %s: %+v", *id, err)
			}
			if err := d.Set("webhook_endpoint", flattenEventSubscriptionWebhookEndpoint(destination, fullUrlResp.Model)); err != nil {
				return fmt.Errorf("setting `webhook_endpoint` for %s: %+v", *id, err)
			}

			deadLetterDestination := props.DeadLetterDestination
			deadLetterIdentityFlattened := make([]interface{}, 0)
			if deadLetterIdentity := props.DeadLetterWithResourceIdentity; deadLetterIdentity != nil {
				deadLetterDestination = deadLetterIdentity.DeadLetterDestination
				deadLetterIdentityFlattened = flattenEventSubscriptionIdentity(deadLetterIdentity.Identity)
			}
			if err := d.Set("dead_letter_identity", deadLetterIdentityFlattened); err != nil {
				return fmt.Errorf("setting `dead_letter_identity` for %s: %+v", *id, err)
			}
			if err := d.Set("storage_blob_dead_letter_destination", flattenEventSubscriptionStorageBlobDeadLetterDestination(deadLetterDestination)); err != nil {
				return fmt.Errorf("setting `storage_blob_dead_letter_destination` for %s: %+v", *id, err)
			}

			enableAdvancedFilteringOnArrays := false
			includedEventTypes := make([]string, 0)
			if filter := props.Filter; filter != nil {
				enableAdvancedFilteringOnArrays = pointer.From(filter.EnableAdvancedFilteringOnArrays)
				includedEventTypes = pointer.From(filter.IncludedEventTypes)
			}
			d.Set("advanced_filtering_on_arrays_enabled", enableAdvancedFilteringOnArrays)
			d.Set("included_event_types", includedEventTypes)
			if err := d.Set("advanced_filter", flattenEventSubscriptionAdvancedFilter(props.Filter)); err != nil {
				return fmt.Errorf("setting `advanced_filter` for %s: %+v", *id, err)
			}
			if err := d.Set("retry_policy", flattenEventSubscriptionRetryPolicy(props.RetryPolicy)); err != nil {
				return fmt.Errorf("setting `retry_policy` for %s: %+v", *id, err)
			}
			if err := d.Set("subject_filter", flattenEventSubscriptionSubjectFilter(props.Filter)); err != nil {
				return fmt.Errorf("setting `subject_filter` for %s: %+v", *id, err)
			}

			d.Set("labels", props.Labels)
		}
	}

	return nil
}

func resourceEventGridEventSubscriptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptions
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := eventsubscriptions.ParseScopedEventSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
