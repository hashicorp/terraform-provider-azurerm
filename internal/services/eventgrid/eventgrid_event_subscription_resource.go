package eventgrid

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2021-12-01/eventgrid" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func PossibleEventSubscriptionEndpointTypes() []string {
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

		CustomizeDiff: pluginsdk.CustomizeDiffShim(eventSubscriptionCustomizeDiffAdvancedFilter),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EventSubscriptionID(id)
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
					PossibleEventSubscriptionEndpointTypes(),
					string(AzureFunctionEndpoint),
				),
			),

			"eventhub_endpoint_id": eventSubscriptionSchemaEventHubEndpointID(
				utils.RemoveFromStringArray(
					PossibleEventSubscriptionEndpointTypes(),
					string(EventHubEndpointID),
				),
			),

			"hybrid_connection_endpoint_id": eventSubscriptionSchemaHybridConnectionEndpointID(
				utils.RemoveFromStringArray(
					PossibleEventSubscriptionEndpointTypes(),
					string(HybridConnectionEndpointID),
				),
			),

			// TODO: this can become `service_bus_queue_id` in 4.0
			"service_bus_queue_endpoint_id": eventSubscriptionSchemaServiceBusQueueEndpointID(
				utils.RemoveFromStringArray(
					PossibleEventSubscriptionEndpointTypes(),
					string(ServiceBusQueueEndpointID),
				),
			),

			// TODO: this can become `service_bus_topic_id` in 4.0
			"service_bus_topic_endpoint_id": eventSubscriptionSchemaServiceBusTopicEndpointID(
				utils.RemoveFromStringArray(
					PossibleEventSubscriptionEndpointTypes(),
					string(ServiceBusTopicEndpointID),
				),
			),

			"storage_queue_endpoint": eventSubscriptionSchemaStorageQueueEndpoint(
				utils.RemoveFromStringArray(
					PossibleEventSubscriptionEndpointTypes(),
					string(StorageQueueEndpoint),
				),
			),

			"webhook_endpoint": eventSubscriptionSchemaWebHookEndpoint(
				utils.RemoveFromStringArray(
					PossibleEventSubscriptionEndpointTypes(),
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
	client := meta.(*clients.Client).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewEventSubscriptionID(d.Get("scope").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.Scope, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_eventgrid_event_subscription", id.ID())
		}
	}

	destination := expandEventGridEventSubscriptionDestination(d)
	if destination == nil {
		return fmt.Errorf("One of the following endpoint types must be specificed to create an EventGrid Event Subscription: %q", PossibleEventSubscriptionEndpointTypes())
	}

	filter, err := expandEventGridEventSubscriptionFilter(d)
	if err != nil {
		return fmt.Errorf("expanding filters for %s: %+v", id, err)
	}

	expirationTime, err := expandEventGridExpirationTime(d)
	if err != nil {
		return fmt.Errorf("expanding `expiration_time` for %s: %+v", id, err)
	}

	deadLetterDestination := expandEventGridEventSubscriptionStorageBlobDeadLetterDestination(d)

	eventSubscriptionProperties := eventgrid.EventSubscriptionProperties{
		Filter:              filter,
		RetryPolicy:         expandEventGridEventSubscriptionRetryPolicy(d),
		Labels:              utils.ExpandStringSlice(d.Get("labels").([]interface{})),
		EventDeliverySchema: eventgrid.EventDeliverySchema(d.Get("event_delivery_schema").(string)),
		ExpirationTimeUtc:   expirationTime,
	}

	if v, ok := d.GetOk("delivery_identity"); ok {
		deliveryIdentityRaw := v.([]interface{})
		deliveryIdentity, err := expandEventGridEventSubscriptionIdentity(deliveryIdentityRaw)
		if err != nil {
			return fmt.Errorf("expanding `delivery_identity`: %+v", err)
		}

		eventSubscriptionProperties.DeliveryWithResourceIdentity = &eventgrid.DeliveryWithResourceIdentity{
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
		deadLetterIdentityRaw := v.([]interface{})
		deadLetterIdentity, err := expandEventGridEventSubscriptionIdentity(deadLetterIdentityRaw)
		if err != nil {
			return fmt.Errorf("expanding `dead_letter_identity`: %+v", err)
		}

		eventSubscriptionProperties.DeadLetterWithResourceIdentity = &eventgrid.DeadLetterWithResourceIdentity{
			Identity:              deadLetterIdentity,
			DeadLetterDestination: deadLetterDestination,
		}
	} else {
		eventSubscriptionProperties.DeadLetterDestination = deadLetterDestination
	}

	eventSubscription := eventgrid.EventSubscription{
		EventSubscriptionProperties: &eventSubscriptionProperties,
	}

	future, err := client.CreateOrUpdate(ctx, id.Scope, id.Name, eventSubscription)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceEventGridEventSubscriptionRead(d, meta)
}

func resourceEventGridEventSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Scope, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("scope", id.Scope)

	if props := resp.EventSubscriptionProperties; props != nil {
		if props.ExpirationTimeUtc != nil {
			d.Set("expiration_time_utc", props.ExpirationTimeUtc.Format(time.RFC3339))
		}

		d.Set("event_delivery_schema", string(props.EventDeliverySchema))

		destination := props.Destination
		deliveryIdentityFlattened := make([]interface{}, 0)
		if deliveryIdentity := props.DeliveryWithResourceIdentity; deliveryIdentity != nil {
			destination = deliveryIdentity.Destination
			deliveryIdentityFlattened = flattenEventGridEventSubscriptionIdentity(deliveryIdentity.Identity)
		}
		if err := d.Set("delivery_identity", deliveryIdentityFlattened); err != nil {
			return fmt.Errorf("setting `delivery_identity` for %s: %+v", *id, err)
		}

		if azureFunctionEndpoint, ok := destination.AsAzureFunctionEventSubscriptionDestination(); ok {
			if err := d.Set("azure_function_endpoint", flattenEventGridEventSubscriptionAzureFunctionEndpoint(azureFunctionEndpoint)); err != nil {
				return fmt.Errorf("setting `azure_function_endpoint` for %s: %+v", *id, err)
			}

			if azureFunctionEndpoint.DeliveryAttributeMappings != nil {
				if err := d.Set("delivery_property", flattenDeliveryProperties(d, azureFunctionEndpoint.DeliveryAttributeMappings)); err != nil {
					return fmt.Errorf("setting `delivery_property` for %s: %+v", *id, err)
				}
			}
		}
		if v, ok := destination.AsEventHubEventSubscriptionDestination(); ok {
			if err := d.Set("eventhub_endpoint_id", v.ResourceID); err != nil {
				return fmt.Errorf("setting `eventhub_endpoint_id` for %s: %+v", *id, err)
			}

			if v.DeliveryAttributeMappings != nil {
				if err := d.Set("delivery_property", flattenDeliveryProperties(d, v.DeliveryAttributeMappings)); err != nil {
					return fmt.Errorf("setting `delivery_property` for %s: %+v", *id, err)
				}
			}
		}
		if v, ok := destination.AsHybridConnectionEventSubscriptionDestination(); ok {
			if err := d.Set("hybrid_connection_endpoint_id", v.ResourceID); err != nil {
				return fmt.Errorf("setting `hybrid_connection_endpoint_id` for %s: %+v", *id, err)
			}

			if v.DeliveryAttributeMappings != nil {
				if err := d.Set("delivery_property", flattenDeliveryProperties(d, v.DeliveryAttributeMappings)); err != nil {
					return fmt.Errorf("setting `delivery_property` for %s: %+v", *id, err)
				}
			}
		}
		if serviceBusQueueEndpoint, ok := destination.AsServiceBusQueueEventSubscriptionDestination(); ok {
			if err := d.Set("service_bus_queue_endpoint_id", serviceBusQueueEndpoint.ResourceID); err != nil {
				return fmt.Errorf("setting `service_bus_queue_endpoint_id` for %s: %+v", *id, err)
			}
		}
		if serviceBusTopicEndpoint, ok := destination.AsServiceBusTopicEventSubscriptionDestination(); ok {
			if err := d.Set("service_bus_topic_endpoint_id", serviceBusTopicEndpoint.ResourceID); err != nil {
				return fmt.Errorf("setting `service_bus_topic_endpoint_id` for %s: %+v", *id, err)
			}
			if serviceBusTopicEndpoint.DeliveryAttributeMappings != nil {
				if err := d.Set("delivery_property", flattenDeliveryProperties(d, serviceBusTopicEndpoint.DeliveryAttributeMappings)); err != nil {
					return fmt.Errorf("setting `delivery_property` for %s: %+v", *id, err)
				}
			}
		}
		if v, ok := destination.AsStorageQueueEventSubscriptionDestination(); ok {
			if err := d.Set("storage_queue_endpoint", flattenEventGridEventSubscriptionStorageQueueEndpoint(v)); err != nil {
				return fmt.Errorf("setting `storage_queue_endpoint` for %s: %+v", *id, err)
			}
		}
		if v, ok := destination.AsWebHookEventSubscriptionDestination(); ok {
			fullURL, err := client.GetFullURL(ctx, id.Scope, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving Full WebHook URL for %s: %+v", *id, err)
			}
			if err := d.Set("webhook_endpoint", flattenEventGridEventSubscriptionWebhookEndpoint(v, &fullURL)); err != nil {
				return fmt.Errorf("setting `webhook_endpoint` for %s: %+v", *id, err)
			}

			if v.DeliveryAttributeMappings != nil {
				if err := d.Set("delivery_property", flattenDeliveryProperties(d, v.DeliveryAttributeMappings)); err != nil {
					return fmt.Errorf("setting `delivery_property` for %s: %+v", *id, err)
				}
			}
		}

		deadLetterDestination := props.DeadLetterDestination
		deadLetterIdentityFlattened := make([]interface{}, 0)
		if deadLetterIdentity := props.DeadLetterWithResourceIdentity; deadLetterIdentity != nil {
			deadLetterDestination = deadLetterIdentity.DeadLetterDestination
			deadLetterIdentityFlattened = flattenEventGridEventSubscriptionIdentity(deadLetterIdentity.Identity)
		}
		if err := d.Set("dead_letter_identity", deadLetterIdentityFlattened); err != nil {
			return fmt.Errorf("setting `dead_letter_identity` for %s: %+v", *id, err)
		}

		if deadLetterDestination != nil {
			if storageBlobDeadLetterDestination, ok := deadLetterDestination.AsStorageBlobDeadLetterDestination(); ok {
				if err := d.Set("storage_blob_dead_letter_destination", flattenEventGridEventSubscriptionStorageBlobDeadLetterDestination(storageBlobDeadLetterDestination)); err != nil {
					return fmt.Errorf("Error setting `storage_blob_dead_letter_destination` for %s: %+v", *id, err)
				}
			}
		}

		if filter := props.Filter; filter != nil {
			d.Set("included_event_types", filter.IncludedEventTypes)
			d.Set("advanced_filtering_on_arrays_enabled", filter.EnableAdvancedFilteringOnArrays)
			if err := d.Set("subject_filter", flattenEventGridEventSubscriptionSubjectFilter(filter)); err != nil {
				return fmt.Errorf("setting `subject_filter` for %s: %+v", *id, err)
			}
			if err := d.Set("advanced_filter", flattenEventGridEventSubscriptionAdvancedFilter(filter)); err != nil {
				return fmt.Errorf("setting `advanced_filter` for %s: %+v", *id, err)
			}
		}

		if props.DeadLetterDestination != nil {
			if storageBlobDeadLetterDestination, ok := props.DeadLetterDestination.AsStorageBlobDeadLetterDestination(); ok {
				if err := d.Set("storage_blob_dead_letter_destination", flattenEventGridEventSubscriptionStorageBlobDeadLetterDestination(storageBlobDeadLetterDestination)); err != nil {
					return fmt.Errorf("setting `storage_blob_dead_letter_destination` for %s: %+v", *id, err)
				}
			}
		}

		if retryPolicy := props.RetryPolicy; retryPolicy != nil {
			if err := d.Set("retry_policy", flattenEventGridEventSubscriptionRetryPolicy(retryPolicy)); err != nil {
				return fmt.Errorf("setting `retry_policy` for %s: %+v", *id, err)
			}
		}

		if err := d.Set("labels", props.Labels); err != nil {
			return fmt.Errorf("setting `labels` for %s: %+v", *id, err)
		}
	}

	return nil
}

func resourceEventGridEventSubscriptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.Scope, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}
