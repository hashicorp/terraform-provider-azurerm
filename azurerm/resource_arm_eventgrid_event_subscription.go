package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2018-09-15-preview/eventgrid"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"event_delivery_schema": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(eventgrid.EventGridSchema),
				ValidateFunc: validation.StringInSlice([]string{
					string(eventgrid.CloudEventV01Schema),
					string(eventgrid.CustomInputSchema),
					string(eventgrid.EventGridSchema),
				}, false),
			},

			"topic_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"storage_queue_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"eventhub_endpoint", "hybrid_connection_endpoint", "webhook_endpoint"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_account_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"queue_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.URLIsHTTPS,
						},
					},
				},
			},

			"included_event_types": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"subject_filter": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subject_begins_with": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subject_ends_with": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"case_sensitive": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"storage_blob_dead_letter_destination": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_account_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"storage_blob_container_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"retry_policy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_delivery_attempts": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 30),
						},
						"event_time_to_live": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 1440),
						},
					},
				},
			},

			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmEventGridEventSubscriptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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
		Destination:           destination,
		Filter:                expandEventGridEventSubscriptionFilter(d),
		DeadLetterDestination: expandEventGridEventSubscriptionStorageBlobDeadLetterDestination(d),
		RetryPolicy:           expandEventGridEventSubscriptionRetryPolicy(d),
		Labels:                utils.ExpandStringSlice(d.Get("labels").([]interface{})),
		EventDeliverySchema:   eventgrid.EventDeliverySchema(d.Get("event_delivery_schema").(string)),
	}

	if v, ok := d.GetOk("topic_name"); ok {
		eventSubscriptionProperties.Topic = utils.String(v.(string))
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
	client := meta.(*ArmClient).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

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
		d.Set("event_delivery_schema", string(props.EventDeliverySchema))

		if props.Topic != nil && *props.Topic != "" {
			d.Set("topic_name", props.Topic)
		}

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
		if _, ok := props.Destination.AsWebHookEventSubscriptionDestination(); ok {
			fullURL, err := client.GetFullURL(ctx, scope, name)
			if err != nil {
				return fmt.Errorf("Error making Read request on EventGrid Event Subscription full URL '%s': %+v", name, err)
			}
			if err := d.Set("webhook_endpoint", flattenEventGridEventSubscriptionWebhookEndpoint(&fullURL)); err != nil {
				return fmt.Errorf("Error setting `webhook_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
			}
		}

		if filter := props.Filter; filter != nil {
			d.Set("included_event_types", filter.IncludedEventTypes)
			if err := d.Set("subject_filter", flattenEventGridEventSubscriptionSubjectFilter(filter)); err != nil {
				return fmt.Errorf("Error setting `subject_filter` for EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
			}
		}

		if props.DeadLetterDestination != nil {
			if storageBlobDeadLetterDestination, ok := props.DeadLetterDestination.AsStorageBlobDeadLetterDestination(); ok {
				if err := d.Set("storage_blob_dead_letter_destination", flattenEventGridEventSubscriptionStorageBlobDeadLetterDestination(storageBlobDeadLetterDestination)); err != nil {
					return fmt.Errorf("Error setting `storage_blob_dead_letter_destination` for EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
				}
			}
		}

		if retryPolicy := props.RetryPolicy; retryPolicy != nil {
			if err := d.Set("retry_policy", flattenEventGridEventSubscriptionRetryPolicy(retryPolicy)); err != nil {
				return fmt.Errorf("Error setting `retry_policy` for EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
			}
		}

		if err := d.Set("labels", props.Labels); err != nil {
			return fmt.Errorf("Error setting `labels` for EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
		}
	}

	return nil
}

func resourceArmEventGridEventSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

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
	props := d.Get("webhook_endpoint").([]interface{})[0].(map[string]interface{})
	url := props["url"].(string)

	webhookEndpoint := eventgrid.WebHookEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeWebHook,
		WebHookEventSubscriptionDestinationProperties: &eventgrid.WebHookEventSubscriptionDestinationProperties{
			EndpointURL: &url,
		},
	}
	return webhookEndpoint
}

func expandEventGridEventSubscriptionFilter(d *schema.ResourceData) *eventgrid.EventSubscriptionFilter {
	filter := &eventgrid.EventSubscriptionFilter{}

	if includedEvents, ok := d.GetOk("included_event_types"); ok {
		filter.IncludedEventTypes = utils.ExpandStringSlice(includedEvents.([]interface{}))
	}

	if subjectFilter, ok := d.GetOk("subject_filter"); ok {
		config := subjectFilter.([]interface{})[0].(map[string]interface{})
		subjectBeginsWith := config["subject_begins_with"].(string)
		subjectEndsWith := config["subject_ends_with"].(string)
		caseSensitive := config["case_sensitive"].(bool)

		filter.SubjectBeginsWith = &subjectBeginsWith
		filter.SubjectEndsWith = &subjectEndsWith
		filter.IsSubjectCaseSensitive = &caseSensitive
	}

	return filter
}

func expandEventGridEventSubscriptionStorageBlobDeadLetterDestination(d *schema.ResourceData) eventgrid.BasicDeadLetterDestination {
	if v, ok := d.GetOk("storage_blob_dead_letter_destination"); ok {
		dest := v.([]interface{})[0].(map[string]interface{})
		resourceID := dest["storage_account_id"].(string)
		blobName := dest["storage_blob_container_name"].(string)
		return eventgrid.StorageBlobDeadLetterDestination{
			EndpointType: eventgrid.EndpointTypeStorageBlob,
			StorageBlobDeadLetterDestinationProperties: &eventgrid.StorageBlobDeadLetterDestinationProperties{
				ResourceID:        &resourceID,
				BlobContainerName: &blobName,
			},
		}
	}
	return nil
}

func expandEventGridEventSubscriptionRetryPolicy(d *schema.ResourceData) *eventgrid.RetryPolicy {
	if v, ok := d.GetOk("retry_policy"); ok {
		dest := v.([]interface{})[0].(map[string]interface{})
		maxDeliveryAttempts := dest["max_delivery_attempts"].(int)
		eventTimeToLive := dest["event_time_to_live"].(int)
		return &eventgrid.RetryPolicy{
			MaxDeliveryAttempts:      utils.Int32(int32(maxDeliveryAttempts)),
			EventTimeToLiveInMinutes: utils.Int32(int32(eventTimeToLive)),
		}
	}
	return nil
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

func flattenEventGridEventSubscriptionWebhookEndpoint(input *eventgrid.EventSubscriptionFullURL) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.EndpointURL != nil {
		result["url"] = *input.EndpointURL
	}

	return []interface{}{result}
}

func flattenEventGridEventSubscriptionSubjectFilter(filter *eventgrid.EventSubscriptionFilter) []interface{} {
	if (filter.SubjectBeginsWith != nil && *filter.SubjectBeginsWith == "") && (filter.SubjectEndsWith != nil && *filter.SubjectEndsWith == "") {
		return nil
	}
	result := make(map[string]interface{})

	if filter.SubjectBeginsWith != nil {
		result["subject_begins_with"] = *filter.SubjectBeginsWith
	}

	if filter.SubjectEndsWith != nil {
		result["subject_ends_with"] = *filter.SubjectEndsWith
	}

	if filter.IsSubjectCaseSensitive != nil {
		result["case_sensitive"] = *filter.IsSubjectCaseSensitive
	}

	return []interface{}{result}
}

func flattenEventGridEventSubscriptionStorageBlobDeadLetterDestination(dest *eventgrid.StorageBlobDeadLetterDestination) []interface{} {
	if dest == nil {
		return nil
	}
	result := make(map[string]interface{})

	if dest.ResourceID != nil {
		result["storage_account_id"] = *dest.ResourceID
	}

	if dest.BlobContainerName != nil {
		result["storage_blob_container_name"] = *dest.BlobContainerName
	}

	return []interface{}{result}
}

func flattenEventGridEventSubscriptionRetryPolicy(retryPolicy *eventgrid.RetryPolicy) []interface{} {
	result := make(map[string]interface{})

	if v := retryPolicy.EventTimeToLiveInMinutes; v != nil {
		result["event_time_to_live"] = int(*v)
	}

	if v := retryPolicy.MaxDeliveryAttempts; v != nil {
		result["max_delivery_attempts"] = int(*v)
	}

	return []interface{}{result}
}
