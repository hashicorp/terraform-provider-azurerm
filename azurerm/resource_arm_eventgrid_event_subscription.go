package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2018-09-15-preview/eventgrid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmEventGridEventSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventGridEventSubscriptionCreateUpdate,
		Read:   resourceArmEventGridEventSubscriptionRead,
		Update: resourceArmEventGridEventSubscriptionCreateUpdate,
		Delete: resourceArmEventGridEventSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"storage_queue_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"eventhub_endpoint", "hybrid_connection_endpoint", "webhook_endpoint"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_account_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"queue_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"eventhub_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"storage_queue_endpoint", "hybrid_connection_endpoint", "webhook_endpoint"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eventhub_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"hybrid_connection_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"storage_queue_endpoint", "eventhub_endpoint", "webhook_endpoint"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hybrid_connection_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"webhook_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"storage_queue_endpoint", "eventhub_endpoint", "hybrid_connection_endpoint"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmEventGridEventSubscriptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventGridEventSubscriptionsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
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
		return fmt.Errorf("One of `webhook_endpoint`, eventhub_endpoint` `hybrid_connection_endpoint` or `storage_queue_endpoint` must be specificed to create an EventGrid Event Subscription")
	}

	eventSubscriptionProperties := eventgrid.EventSubscriptionProperties{
		Destination: destination,
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

	return resourceArmEventGridEventSubscriptionRead(d, meta)
}

func resourceArmEventGridEventSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventGridEventSubscriptionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureEventGridEventSubscriptionID(d.Id())
	if err != nil {
		return err
	}
	scope := id.Scope
	name := id.Name

	resp, err := client.Get(ctx, scope, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] EventGrid Event Subscription '%s' was not found (resource group '%s')", name, scope)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on EventGrid Event Subscription '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("scope", scope)

	if props := resp.EventSubscriptionProperties; props != nil {
		if storageQueueEndpoint, ok := props.Destination.AsStorageQueueEventSubscriptionDestination(); ok {
			if err := d.Set("storage_queue_endpoint", flattenEventGridEventSubscriptionStorageQueueEndpoint(storageQueueEndpoint)); err != nil {
				return fmt.Errorf("Error setting `storage_queue_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
			}
		}
		if eventHubEndpoint, ok := props.Destination.AsEventHubEventSubscriptionDestination(); ok {
			if err := d.Set("eventhub_endpoint", flattenEventGridEventSubscriptionEventHubEndpoint(eventHubEndpoint)); err != nil {
				return fmt.Errorf("Error setting `eventhub_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
			}
		}
		if hybridConnectionEndpoint, ok := props.Destination.AsHybridConnectionEventSubscriptionDestination(); ok {
			if err := d.Set("hybrid_connection_endpoint", flattenEventGridEventSubscriptionHybridConnectionEndpoint(hybridConnectionEndpoint)); err != nil {
				return fmt.Errorf("Error setting `hybrid_connection_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
			}
		}
		if webhookEndpoint, ok := props.Destination.AsWebHookEventSubscriptionDestination(); ok {
			if err := d.Set("webhook_endpoint", flattenEventGridEventSubscriptionWebhookEndpoint(webhookEndpoint)); err != nil {
				return fmt.Errorf("Error setting `webhook_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
			}
		}
	}

	return nil
}

func resourceArmEventGridEventSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventGridEventSubscriptionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureEventGridEventSubscriptionID(d.Id())
	if err != nil {
		return err
	}
	scope := id.Scope
	name := id.Name

	future, err := client.Delete(ctx, scope, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Event Grid Event Subscription %q: %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Event Grid Event Subscription %q: %+v", name, err)
	}

	return nil
}

type AzureEventGridEventSubscriptionID struct {
	Scope string
	Name  string
}

func parseAzureEventGridEventSubscriptionID(id string) (*AzureEventGridEventSubscriptionID, error) {
	segments := strings.Split(id, "/providers/Microsoft.EventGrid/eventSubscriptions/")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected ID to be in the format `{scope}/providers/Microsoft.EventGrid/eventSubscriptions/{name} - got %d segments", len(segments))
	}

	scope := segments[0]
	name := segments[1]
	eventSubscriptionID := AzureEventGridEventSubscriptionID{
		Scope: scope,
		Name:  name,
	}
	return &eventSubscriptionID, nil
}

func expandEventGridEventSubscriptionDestination(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	if _, ok := d.GetOk("storage_queue_endpoint"); ok {
		return expandEventGridEventSubscriptionStorageQueueEndpoint(d)
	}

	if _, ok := d.GetOk("eventhub_endpoint"); ok {
		return expandEventGridEventSubscriptionEventHubEndpoint(d)
	}

	if _, ok := d.GetOk("hybrid_connection_endpoint"); ok {
		return expandEventGridEventSubscriptionHybridConnectionEndpoint(d)
	}

	if _, ok := d.GetOk("webhook_endpoint"); ok {
		return expandEventGridEventSubscriptionWebhookEndpoint(d)
	}

	return nil
}

func expandEventGridEventSubscriptionStorageQueueEndpoint(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	props := d.Get("storage_queue_endpoint").([]interface{})[0].(map[string]interface{})
	storageAccountID := props["storage_account_id"].(string)
	queueName := props["queue_name"].(string)

	storageQueueEndpoint := eventgrid.StorageQueueEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeStorageQueue,
		StorageQueueEventSubscriptionDestinationProperties: &eventgrid.StorageQueueEventSubscriptionDestinationProperties{
			ResourceID: &storageAccountID,
			QueueName:  &queueName,
		},
	}
	return storageQueueEndpoint
}

func expandEventGridEventSubscriptionEventHubEndpoint(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	props := d.Get("eventhub_endpoint").([]interface{})[0].(map[string]interface{})
	eventHubID := props["eventhub_id"].(string)

	eventHubEndpoint := eventgrid.EventHubEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeEventHub,
		EventHubEventSubscriptionDestinationProperties: &eventgrid.EventHubEventSubscriptionDestinationProperties{
			ResourceID: &eventHubID,
		},
	}
	return eventHubEndpoint
}

func expandEventGridEventSubscriptionHybridConnectionEndpoint(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	props := d.Get("hybrid_connection_endpoint").([]interface{})[0].(map[string]interface{})
	hybridConnectionID := props["hybrid_connection_id"].(string)

	hybridConnectionEndpoint := eventgrid.HybridConnectionEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeHybridConnection,
		HybridConnectionEventSubscriptionDestinationProperties: &eventgrid.HybridConnectionEventSubscriptionDestinationProperties{
			ResourceID: &hybridConnectionID,
		},
	}
	return hybridConnectionEndpoint
}

func expandEventGridEventSubscriptionWebhookEndpoint(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	props := d.Get("hybrid_connection_endpoint").([]interface{})[0].(map[string]interface{})
	url := props["url"].(string)

	webhookEndpoint := eventgrid.WebHookEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeWebHook,
		WebHookEventSubscriptionDestinationProperties: &eventgrid.WebHookEventSubscriptionDestinationProperties{
			EndpointURL: &url,
		},
	}
	return webhookEndpoint
}

func flattenEventGridEventSubscriptionStorageQueueEndpoint(input *eventgrid.StorageQueueEventSubscriptionDestination) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.ResourceID != nil {
		result["storage_account_id"] = *input.ResourceID
	}
	if input.QueueName != nil {
		result["queue_name"] = *input.QueueName
	}

	return []interface{}{result}
}

func flattenEventGridEventSubscriptionEventHubEndpoint(input *eventgrid.EventHubEventSubscriptionDestination) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.ResourceID != nil {
		result["eventhub_id"] = *input.ResourceID
	}

	return []interface{}{result}
}

func flattenEventGridEventSubscriptionHybridConnectionEndpoint(input *eventgrid.HybridConnectionEventSubscriptionDestination) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.ResourceID != nil {
		result["eventhub_id"] = *input.ResourceID
	}

	return []interface{}{result}
}

func flattenEventGridEventSubscriptionWebhookEndpoint(input *eventgrid.WebHookEventSubscriptionDestination) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.EndpointURL != nil {
		result["url"] = *input.EndpointURL
	}

	return []interface{}{result}
}
