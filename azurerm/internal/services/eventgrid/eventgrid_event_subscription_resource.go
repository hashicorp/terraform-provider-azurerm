package eventgrid

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2020-04-01-preview/eventgrid"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventgrid/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func PossibleEventSubscriptionEndpointTypes() []string {
	return []string{
		string(AzureFunctionEndpoint),
		string(EventHubEndpoint),
		string(EventHubEndpointID),
		string(HybridConnectionEndpoint),
		string(HybridConnectionEndpointID),
		string(ServiceBusQueueEndpointID),
		string(ServiceBusTopicEndpointID),
		string(StorageQueueEndpoint),
		string(WebHookEndpoint),
	}
}

func resourceEventGridEventSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventGridEventSubscriptionCreateUpdate,
		Read:   resourceEventGridEventSubscriptionRead,
		Update: resourceEventGridEventSubscriptionCreateUpdate,
		Delete: resourceEventGridEventSubscriptionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.EventSubscriptionID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": eventSubscriptionSchemaEventSubscriptionName(),

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"event_delivery_schema": eventSubscriptionSchemaEventDeliverySchema(),

			"expiration_time_utc": eventSubscriptionSchemaExpirationTimeUTC(),

			"topic_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "This field has been updated to readonly field since Apr 25, 2019 so no longer has any affect and will be removed in version 3.0 of the provider.",
			},

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

			"eventhub_endpoint": eventSubscriptionSchemaEventHubEndpoint(
				utils.RemoveFromStringArray(
					PossibleEventSubscriptionEndpointTypes(),
					string(EventHubEndpoint),
				),
			),

			"hybrid_connection_endpoint_id": eventSubscriptionSchemaHybridConnectionEndpointID(
				utils.RemoveFromStringArray(
					PossibleEventSubscriptionEndpointTypes(),
					string(HybridConnectionEndpointID),
				),
			),

			"hybrid_connection_endpoint": eventSubscriptionSchemaHybridEndpoint(
				utils.RemoveFromStringArray(
					PossibleEventSubscriptionEndpointTypes(),
					string(HybridConnectionEndpoint),
				),
			),

			"service_bus_queue_endpoint_id": eventSubscriptionSchemaServiceBusQueueEndpointID(
				utils.RemoveFromStringArray(
					PossibleEventSubscriptionEndpointTypes(),
					string(ServiceBusQueueEndpointID),
				),
			),

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

			"storage_blob_dead_letter_destination": eventSubscriptionSchemaStorageBlobDeadletterDestination(),

			"retry_policy": eventSubscriptionSchemaRetryPolicy(),

			"labels": eventSubscriptionSchemaLabels(),
		},
	}
}

func resourceEventGridEventSubscriptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, scope, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_eventgrid_event_subscription", *existing.ID)
		}
	}

	destination := expandEventGridEventSubscriptionDestination(d)
	if destination == nil {
		return fmt.Errorf("One of the following endpoint types must be specificed to create an EventGrid Event Subscription: %q", PossibleEventSubscriptionEndpointTypes())
	}

	filter, err := expandEventGridEventSubscriptionFilter(d)
	if err != nil {
		return fmt.Errorf("expanding filters for EventGrid Event Subscription %q (Scope %q): %+v", name, scope, err)
	}

	expirationTime, err := expandEventGridExpirationTime(d)
	if err != nil {
		return fmt.Errorf("Error creating/updating EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
	}

	eventSubscriptionProperties := eventgrid.EventSubscriptionProperties{
		Destination:           destination,
		Filter:                filter,
		DeadLetterDestination: expandEventGridEventSubscriptionStorageBlobDeadLetterDestination(d),
		RetryPolicy:           expandEventGridEventSubscriptionRetryPolicy(d),
		Labels:                utils.ExpandStringSlice(d.Get("labels").([]interface{})),
		EventDeliverySchema:   eventgrid.EventDeliverySchema(d.Get("event_delivery_schema").(string)),
		ExpirationTimeUtc:     expirationTime,
	}

	eventSubscription := eventgrid.EventSubscription{
		EventSubscriptionProperties: &eventSubscriptionProperties,
	}

	log.Printf("[INFO] preparing arguments for AzureRM EventGrid Event Subscription creation with Properties: %+v.", eventSubscription)

	future, err := client.CreateOrUpdate(ctx, scope, name, eventSubscription)
	if err != nil {
		return fmt.Errorf("Error creating/updating EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for EventGrid Event Subscription %q (Scope %q) to become available: %s", name, scope, err)
	}

	read, err := client.Get(ctx, scope, name)
	if err != nil {
		return fmt.Errorf("Error retrieving EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read EventGrid Event Subscription %s (Scope %s) ID", name, scope)
	}

	d.SetId(*read.ID)

	return resourceEventGridEventSubscriptionRead(d, meta)
}

func resourceEventGridEventSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[WARN] EventGrid Event Subscription '%s' was not found (resource group '%s')", id.Name, id.Scope)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on EventGrid Event Subscription '%s': %+v", id.Name, err)
	}

	d.Set("name", resp.Name)
	d.Set("scope", id.Scope)

	if props := resp.EventSubscriptionProperties; props != nil {
		if props.ExpirationTimeUtc != nil {
			d.Set("expiration_time_utc", props.ExpirationTimeUtc.Format(time.RFC3339))
		}

		d.Set("event_delivery_schema", string(props.EventDeliverySchema))

		if azureFunctionEndpoint, ok := props.Destination.AsAzureFunctionEventSubscriptionDestination(); ok {
			if err := d.Set("azure_function_endpoint", flattenEventGridEventSubscriptionAzureFunctionEndpoint(azureFunctionEndpoint)); err != nil {
				return fmt.Errorf("Error setting `%q` for EventGrid Event Subscription %q (Scope %q): %s", "azure_function_endpoint", id.Name, id.Scope, err)
			}
		}
		if v, ok := props.Destination.AsEventHubEventSubscriptionDestination(); ok {
			if err := d.Set("eventhub_endpoint_id", v.ResourceID); err != nil {
				return fmt.Errorf("Error setting `%q` for EventGrid Event Subscription %q (Scope %q): %s", "eventhub_endpoint_id", id.Name, id.Scope, err)
			}

			if err := d.Set("eventhub_endpoint", flattenEventGridEventSubscriptionEventhubEndpoint(v)); err != nil {
				return fmt.Errorf("Error setting `%q` for EventGrid Event Subscription %q (Scope %q): %s", "eventhub_endpoint", id.Name, id.Scope, err)
			}
		}
		if v, ok := props.Destination.AsHybridConnectionEventSubscriptionDestination(); ok {
			if err := d.Set("hybrid_connection_endpoint_id", v.ResourceID); err != nil {
				return fmt.Errorf("Error setting `%q` for EventGrid Event Subscription %q (Scope %q): %s", "hybrid_connection_endpoint_id", id.Name, id.Scope, err)
			}

			if err := d.Set("hybrid_connection_endpoint", flattenEventGridEventSubscriptionHybridConnectionEndpoint(v)); err != nil {
				return fmt.Errorf("Error setting `%q` for EventGrid Event Subscription %q (Scope %q): %s", "hybrid_connection_endpoint", id.Name, id.Scope, err)
			}
		}
		if serviceBusQueueEndpoint, ok := props.Destination.AsServiceBusQueueEventSubscriptionDestination(); ok {
			if err := d.Set("service_bus_queue_endpoint_id", serviceBusQueueEndpoint.ResourceID); err != nil {
				return fmt.Errorf("Error setting `%q` for EventGrid Event Subscription %q (Scope %q): %s", "service_bus_queue_endpoint_id", id.Name, id.Scope, err)
			}
		}
		if serviceBusTopicEndpoint, ok := props.Destination.AsServiceBusTopicEventSubscriptionDestination(); ok {
			if err := d.Set("service_bus_topic_endpoint_id", serviceBusTopicEndpoint.ResourceID); err != nil {
				return fmt.Errorf("Error setting `%q` for EventGrid Event Subscription %q (Scope %q): %s", "service_bus_topic_endpoint_id", id.Name, id.Scope, err)
			}
		}
		if v, ok := props.Destination.AsStorageQueueEventSubscriptionDestination(); ok {
			if err := d.Set("storage_queue_endpoint", flattenEventGridEventSubscriptionStorageQueueEndpoint(v)); err != nil {
				return fmt.Errorf("Error setting `%q` for EventGrid Event Subscription %q (Scope %q): %s", "storage_queue_endpoint", id.Name, id.Scope, err)
			}
		}
		if v, ok := props.Destination.AsWebHookEventSubscriptionDestination(); ok {
			fullURL, err := client.GetFullURL(ctx, id.Scope, id.Name)
			if err != nil {
				return fmt.Errorf("Error making Read request on EventGrid Event Subscription full URL '%s': %+v", id.Name, err)
			}
			if err := d.Set("webhook_endpoint", flattenEventGridEventSubscriptionWebhookEndpoint(v, &fullURL)); err != nil {
				return fmt.Errorf("Error setting `%q` for EventGrid Event Subscription %q (Scope %q): %s", "webhook_endpoint", id.Name, id.Scope, err)
			}
		}

		if filter := props.Filter; filter != nil {
			d.Set("included_event_types", filter.IncludedEventTypes)
			if err := d.Set("subject_filter", flattenEventGridEventSubscriptionSubjectFilter(filter)); err != nil {
				return fmt.Errorf("Error setting `subject_filter` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
			}
			if err := d.Set("advanced_filter", flattenEventGridEventSubscriptionAdvancedFilter(filter)); err != nil {
				return fmt.Errorf("Error setting `advanced_filter` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
			}
		}

		if props.DeadLetterDestination != nil {
			if storageBlobDeadLetterDestination, ok := props.DeadLetterDestination.AsStorageBlobDeadLetterDestination(); ok {
				if err := d.Set("storage_blob_dead_letter_destination", flattenEventGridEventSubscriptionStorageBlobDeadLetterDestination(storageBlobDeadLetterDestination)); err != nil {
					return fmt.Errorf("Error setting `storage_blob_dead_letter_destination` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
				}
			}
		}

		if retryPolicy := props.RetryPolicy; retryPolicy != nil {
			if err := d.Set("retry_policy", flattenEventGridEventSubscriptionRetryPolicy(retryPolicy)); err != nil {
				return fmt.Errorf("Error setting `retry_policy` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
			}
		}

		if err := d.Set("labels", props.Labels); err != nil {
			return fmt.Errorf("Error setting `labels` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
		}
	}

	return nil
}

func resourceEventGridEventSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.Scope, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Event Grid Event Subscription %q: %+v", id.Name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Event Grid Event Subscription %q: %+v", id.Name, err)
	}

	return nil
}
