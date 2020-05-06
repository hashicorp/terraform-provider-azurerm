package eventgrid

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2020-01-01-preview/eventgrid"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func getEnpointTypes() []string {
	return []string{"webhook_endpoint", "storage_queue_endpoint", "eventhub_endpoint", "hybrid_connection_endpoint", "service_bus_queue_endpoint", "service_bus_topic_endpoint", "azure_function_endpoint"}
}

// RemoveFromStringArray removes all matching values from a string array
func RemoveFromStringArray(elements []string, remove string) []string {
	for i, v := range elements {
		if v == remove {
			return append(elements[:i], elements[i+1:]...)
		}
	}
	return elements
}

// AdvancedFilterDiffSuppressFunc performs a type relaxed diff
func AdvancedFilterDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if strings.EqualFold(old, new) {
		return true
	} else if o, err := strconv.ParseFloat(old, 64); err == nil {
		n, err := strconv.ParseFloat(new, 64)
		if err == nil {
			return o == n
		}
	} else if o, err := strconv.ParseBool(old); err == nil {
		n, err := strconv.ParseBool(new)
		if err == nil {
			return o == n
		}
	}
	return false
}

// ParseFloat tries to convert a string to a float64 value
func ParseFloat(value string) (*float64, error) {
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return &f, nil
	}
	return nil, fmt.Errorf("Value %q is not a float number", value)
}

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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"event_delivery_schema": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(eventgrid.EventGridSchema),
				ValidateFunc: validation.StringInSlice([]string{
					string(eventgrid.EventGridSchema),
					string(eventgrid.CloudEventSchemaV10),
					string(eventgrid.CustomInputSchema),
				}, false),
			},

			"expiration_time_utc": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_queue_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: RemoveFromStringArray(getEnpointTypes(), "storage_queue_endpoint"),
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
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"eventhub_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: RemoveFromStringArray(getEnpointTypes(), "eventhub_endpoint"),
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
				ConflictsWith: RemoveFromStringArray(getEnpointTypes(), "hybrid_connection_endpoint"),
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
				ConflictsWith: RemoveFromStringArray(getEnpointTypes(), "webhook_endpoint"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},
					},
				},
			},

			"service_bus_queue_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: RemoveFromStringArray(getEnpointTypes(), "service_bus_queue_endpoint"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_bus_queue_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"service_bus_topic_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: RemoveFromStringArray(getEnpointTypes(), "service_bus_topic_endpoint"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_bus_queue_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"azure_function_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: RemoveFromStringArray(getEnpointTypes(), "azure_function_endpoint"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"azure_function_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"included_event_types": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
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

			"advanced_filter": {
				Type:     schema.TypeList,
				MaxItems: 5,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"operator_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(eventgrid.OperatorTypeAdvancedFilter),
								string(eventgrid.OperatorTypeBoolEquals),
								string(eventgrid.OperatorTypeNumberGreaterThan),
								string(eventgrid.OperatorTypeNumberGreaterThanOrEquals),
								string(eventgrid.OperatorTypeNumberLessThan),
								string(eventgrid.OperatorTypeNumberLessThanOrEquals),
								string(eventgrid.OperatorTypeNumberIn),
								string(eventgrid.OperatorTypeNumberNotIn),
								string(eventgrid.OperatorTypeStringBeginsWith),
								string(eventgrid.OperatorTypeStringContains),
								string(eventgrid.OperatorTypeStringEndsWith),
								string(eventgrid.OperatorTypeStringIn),
								string(eventgrid.OperatorTypeStringNotIn),
							}, false),
						},
						"value": {
							Type:             schema.TypeString,
							Optional:         true,
							ConflictsWith:    []string{"advanced_filter.values"},
							DiffSuppressFunc: AdvancedFilterDiffSuppressFunc,
						},
						"values": {
							Type:          schema.TypeList,
							MinItems:      1,
							MaxItems:      5,
							Optional:      true,
							ConflictsWith: []string{"advanced_filter.value"},
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateFunc:     validation.StringIsNotEmpty,
								DiffSuppressFunc: AdvancedFilterDiffSuppressFunc,
							},
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
							ValidateFunc: validation.StringIsNotEmpty,
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

			"topic_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmEventGridEventSubscriptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
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
		return fmt.Errorf("One of the following endpoint types must be specificed to create an EventGrid Event Subscription: %q", getEnpointTypes())
	}

	filter, err := expandEventGridEventSubscriptionFilter(d)
	if err != nil {
		return fmt.Errorf("Error creating/updating EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
	}

	parsedTime, err := date.ParseTime(time.RFC3339, d.Get("expiration_time_utc").(string))
	if err != nil {
		return fmt.Errorf("Error creating/updating EventGrid Event Subscription %q (Scope %q): %s", name, scope, err)
	}

	expirationTime := date.Time{Time: parsedTime}

	eventSubscriptionProperties := eventgrid.EventSubscriptionProperties{
		Destination:           destination,
		Filter:                filter,
		DeadLetterDestination: expandEventGridEventSubscriptionStorageBlobDeadLetterDestination(d),
		RetryPolicy:           expandEventGridEventSubscriptionRetryPolicy(d),
		Labels:                utils.ExpandStringSlice(d.Get("labels").([]interface{})),
		EventDeliverySchema:   eventgrid.EventDeliverySchema(d.Get("event_delivery_schema").(string)),
		ExpirationTimeUtc:     &expirationTime,
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
	client := meta.(*clients.Client).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureEventGridEventSubscriptionID(d.Id())
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
		d.Set("expiration_time_utc", props.ExpirationTimeUtc.Format(time.RFC3339))

		if props.Topic != nil && *props.Topic != "" {
			d.Set("topic_name", props.Topic)
		}

		if storageQueueEndpoint, ok := props.Destination.AsStorageQueueEventSubscriptionDestination(); ok {
			if err := d.Set("storage_queue_endpoint", flattenEventGridEventSubscriptionStorageQueueEndpoint(storageQueueEndpoint)); err != nil {
				return fmt.Errorf("Error setting `storage_queue_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
			}
		}
		if eventHubEndpoint, ok := props.Destination.AsEventHubEventSubscriptionDestination(); ok {
			if err := d.Set("eventhub_endpoint", flattenEventGridEventSubscriptionEventHubEndpoint(eventHubEndpoint)); err != nil {
				return fmt.Errorf("Error setting `eventhub_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
			}
		}
		if hybridConnectionEndpoint, ok := props.Destination.AsHybridConnectionEventSubscriptionDestination(); ok {
			if err := d.Set("hybrid_connection_endpoint", flattenEventGridEventSubscriptionHybridConnectionEndpoint(hybridConnectionEndpoint)); err != nil {
				return fmt.Errorf("Error setting `hybrid_connection_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
			}
		}
		if serviceBusQueueEndpoint, ok := props.Destination.AsServiceBusQueueEventSubscriptionDestination(); ok {
			if err := d.Set("service_bus_queue_endpoint", flattenEventGridEventSubscriptionServiceBusQueueEndpoint(serviceBusQueueEndpoint)); err != nil {
				return fmt.Errorf("Error setting `service_bus_queue_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
			}
		}
		if serviceBusTopicEndpoint, ok := props.Destination.AsServiceBusTopicEventSubscriptionDestination(); ok {
			if err := d.Set("service_bus_topic_endpoint", flattenEventGridEventSubscriptionServiceBusTopicEndpoint(serviceBusTopicEndpoint)); err != nil {
				return fmt.Errorf("Error setting `service_bus_topic_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
			}
		}
		if azureFunctionEndpoint, ok := props.Destination.AsAzureFunctionEventSubscriptionDestination(); ok {
			if err := d.Set("azure_function_endpoint", flattenEventGridEventSubscriptionAzureFunctionEndpoint(azureFunctionEndpoint)); err != nil {
				return fmt.Errorf("Error setting `azure_function_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
			}
		}
		if _, ok := props.Destination.AsWebHookEventSubscriptionDestination(); ok {
			fullURL, err := client.GetFullURL(ctx, id.Scope, id.Name)
			if err != nil {
				return fmt.Errorf("Error making Read request on EventGrid Event Subscription full URL '%s': %+v", id.Name, err)
			}
			if err := d.Set("webhook_endpoint", flattenEventGridEventSubscriptionWebhookEndpoint(&fullURL)); err != nil {
				return fmt.Errorf("Error setting `webhook_endpoint` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
			}
		}

		if filter := props.Filter; filter != nil {
			d.Set("included_event_types", filter.IncludedEventTypes)
			if err := d.Set("subject_filter", flattenEventGridEventSubscriptionSubjectFilter(filter)); err != nil {
				return fmt.Errorf("Error setting `subject_filter` for EventGrid Event Subscription %q (Scope %q): %s", id.Name, id.Scope, err)
			}

			if err := d.Set("advanced_filter", flattenEventGridEventSubscriptionAdvancedFilter(filter.AdvancedFilters)); err != nil {
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

func resourceArmEventGridEventSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureEventGridEventSubscriptionID(d.Id())
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

	if _, ok := d.GetOk("service_bus_queue_endpoint"); ok {
		return expandEventGridEventSubscriptionServiceBusQueueEndpoint(d)
	}

	if _, ok := d.GetOk("service_bus_topic_endpoint"); ok {
		return expandEventGridEventSubscriptionServiceBusTopicEndpoint(d)
	}

	if _, ok := d.GetOk("azure_function"); ok {
		return expandEventGridEventSubscriptionAzureFunctionEndpoint(d)
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

func expandEventGridEventSubscriptionServiceBusQueueEndpoint(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	props := d.Get("service_bus_queue_endpoint").([]interface{})[0].(map[string]interface{})
	serviceBusQueueID := props["service_bus_queue_id"].(string)

	serviceBusQueueEndpoint := eventgrid.ServiceBusQueueEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeServiceBusQueue,
		ServiceBusQueueEventSubscriptionDestinationProperties: &eventgrid.ServiceBusQueueEventSubscriptionDestinationProperties{
			ResourceID: &serviceBusQueueID,
		},
	}
	return serviceBusQueueEndpoint
}

func expandEventGridEventSubscriptionServiceBusTopicEndpoint(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	props := d.Get("service_bus_topic_endpoint").([]interface{})[0].(map[string]interface{})
	serviceBusTopicID := props["service_bus_topic_id"].(string)

	serviceBusTopicEndpoint := eventgrid.ServiceBusTopicEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeServiceBusTopic,
		ServiceBusTopicEventSubscriptionDestinationProperties: &eventgrid.ServiceBusTopicEventSubscriptionDestinationProperties{
			ResourceID: &serviceBusTopicID,
		},
	}
	return serviceBusTopicEndpoint
}

func expandEventGridEventSubscriptionAzureFunctionEndpoint(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	props := d.Get("azure_function_endpoint").([]interface{})[0].(map[string]interface{})
	azureFunctionResourceID := props["azure_function_id"].(string)

	azureFunctionEndpoint := eventgrid.AzureFunctionEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeAzureFunction,
		AzureFunctionEventSubscriptionDestinationProperties: &eventgrid.AzureFunctionEventSubscriptionDestinationProperties{
			ResourceID: &azureFunctionResourceID,
		},
	}
	return azureFunctionEndpoint
}

func expandEventGridEventSubscriptionFilter(d *schema.ResourceData) (*eventgrid.EventSubscriptionFilter, error) {
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

	if advancedFilter, ok := d.GetOk("advanced_filter"); ok {
		advancedFilters := make([]eventgrid.BasicAdvancedFilter, 0)
		for _, v := range advancedFilter.([]interface{}) {
			config := v.(map[string]interface{})

			if filter, err := expandAdvancedFilter(config); err == nil {
				advancedFilters = append(advancedFilters, filter)
			} else {
				return nil, err
			}
		}
		filter.AdvancedFilters = &advancedFilters
	}

	return filter, nil
}

func expandAdvancedFilter(config map[string]interface{}) (eventgrid.BasicAdvancedFilter, error) {
	operatorType := config["operator_type"].(string)
	key := config["key"].(string)
	value := config["value"].(string)
	values := utils.ExpandStringSlice(config["values"].([]interface{}))

	switch operatorType {
	case string(eventgrid.OperatorTypeAdvancedFilter),
		string(eventgrid.OperatorTypeBoolEquals),
		string(eventgrid.OperatorTypeNumberGreaterThan),
		string(eventgrid.OperatorTypeNumberGreaterThanOrEquals),
		string(eventgrid.OperatorTypeNumberLessThan),
		string(eventgrid.OperatorTypeNumberLessThanOrEquals):
		// Workaround as long as nested schema validation is not working as expected (see https://github.com/hashicorp/terraform-plugin-sdk/issues/71)
		if values != nil && len(*values) > 0 {
			return nil, fmt.Errorf("Conflicting field for `advanced_filter` (key=%s, operator_type=%s): values", key, operatorType)
		}
		if &value == nil || len(value) == 0 {
			return nil, fmt.Errorf("Missing value for`advanced_filter` (key=%s, operator_type=%s, value=%q)", key, operatorType, value)
		}
		return expandScalarAdvancedFilter(key, operatorType, value)
	case
		string(eventgrid.OperatorTypeNumberIn),
		string(eventgrid.OperatorTypeNumberNotIn),
		string(eventgrid.OperatorTypeStringBeginsWith),
		string(eventgrid.OperatorTypeStringContains),
		string(eventgrid.OperatorTypeStringEndsWith),
		string(eventgrid.OperatorTypeStringIn),
		string(eventgrid.OperatorTypeStringNotIn):
		// Workaround as long as nested schema validation is not working as expected (see https://github.com/hashicorp/terraform-plugin-sdk/issues/71)
		if &value != nil && len(value) > 0 {
			return nil, fmt.Errorf("Conflicting field for `advanced_filter` (key=%s, operator_type=%s): value", key, operatorType)
		}
		if len(*values) == 0 {
			return nil, fmt.Errorf("Missing values for `advanced_filter` (key=%s, operator_type=%s, values=%q)", key, operatorType, values)
		}
		return expandArrayAdvancedFilter(key, operatorType, *values)
	default:
		return nil, fmt.Errorf("Invalid `advanced_filter` operator_type %s used", operatorType)
	}
}

func expandArrayAdvancedFilter(key string, operatorType string, values []string) (eventgrid.BasicAdvancedFilter, error) {
	switch operatorType {
	case string(eventgrid.OperatorTypeNumberIn):
		var numbers = []float64{}
		for _, v := range values {
			if f, err := ParseFloat(v); err == nil {
				numbers = append(numbers, *f)
			} else {
				return nil, err
			}
		}
		return eventgrid.NumberInAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeNumberIn, Values: &numbers}, nil
	case string(eventgrid.OperatorTypeNumberNotIn):
		var numbers = []float64{}
		for _, v := range values {
			if f, err := ParseFloat(v); err == nil {
				numbers = append(numbers, *f)
			} else {
				return nil, err
			}
		}
		return eventgrid.NumberNotInAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeNumberIn, Values: &numbers}, nil
	case string(eventgrid.OperatorTypeStringIn):
		return eventgrid.StringInAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeStringIn, Values: &values}, nil
	case string(eventgrid.OperatorTypeStringNotIn):
		return eventgrid.StringNotInAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeStringNotIn, Values: &values}, nil
	case string(eventgrid.OperatorTypeStringBeginsWith):
		return eventgrid.StringBeginsWithAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeStringBeginsWith, Values: &values}, nil
	case string(eventgrid.OperatorTypeStringEndsWith):
		return eventgrid.StringEndsWithAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeStringEndsWith, Values: &values}, nil
	case string(eventgrid.OperatorTypeStringContains):
		return eventgrid.StringContainsAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeStringContains, Values: &values}, nil
	default:
		return nil, nil
	}
}

func expandScalarAdvancedFilter(key string, operatorType string, value string) (eventgrid.BasicAdvancedFilter, error) {
	switch operatorType {
	case string(eventgrid.OperatorTypeNumberLessThan):
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return eventgrid.NumberLessThanAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeNumberLessThan, Value: &v}, nil
		}
		return nil, fmt.Errorf("Value %q is not a float number", value)
	case string(eventgrid.OperatorTypeNumberGreaterThan):
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return eventgrid.NumberGreaterThanAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeNumberGreaterThan, Value: &v}, nil
		}
		return nil, fmt.Errorf("Value %q is not a float number", value)
	case string(eventgrid.OperatorTypeNumberLessThanOrEquals):
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return eventgrid.NumberLessThanOrEqualsAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeNumberLessThanOrEquals, Value: &v}, nil
		}
		return nil, fmt.Errorf("Value %q is not a float number", value)
	case string(eventgrid.OperatorTypeNumberGreaterThanOrEquals):
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return eventgrid.NumberGreaterThanOrEqualsAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeNumberGreaterThanOrEquals, Value: &v}, nil
		}
		return nil, fmt.Errorf("Value %q is not a float number", value)
	case string(eventgrid.OperatorTypeBoolEquals):
		v, err := strconv.ParseBool(value)
		if err == nil {
			return eventgrid.BoolEqualsAdvancedFilter{Key: &key, OperatorType: eventgrid.OperatorTypeBoolEquals, Value: &v}, nil
		}
		return nil, fmt.Errorf("Value %q is not a bool", value)
	default:
		return nil, nil
	}
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

func flattenEventGridEventSubscriptionServiceBusQueueEndpoint(input *eventgrid.ServiceBusQueueEventSubscriptionDestination) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.ResourceID != nil {
		result["service_bus_queue_id"] = *input.ResourceID
	}

	return []interface{}{result}
}

func flattenEventGridEventSubscriptionServiceBusTopicEndpoint(input *eventgrid.ServiceBusTopicEventSubscriptionDestination) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.ResourceID != nil {
		result["service_bus_topic_id"] = *input.ResourceID
	}

	return []interface{}{result}
}

func flattenEventGridEventSubscriptionAzureFunctionEndpoint(input *eventgrid.AzureFunctionEventSubscriptionDestination) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.ResourceID != nil {
		result["azure_function_id"] = *input.ResourceID
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

func flattenEventGridEventSubscriptionAdvancedFilter(input *[]eventgrid.BasicAdvancedFilter) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var key interface{}
		var operatorType eventgrid.OperatorType
		var value interface{}
		var values interface{}

		switch f := item.(type) {
		case eventgrid.AdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
		case eventgrid.BoolEqualsAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			value = strconv.FormatBool(*f.Value)
		case eventgrid.NumberGreaterThanAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			value = strconv.FormatFloat(*f.Value, 'f', -1, 64)
		case eventgrid.NumberGreaterThanOrEqualsAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			value = strconv.FormatFloat(*f.Value, 'f', -1, 64)
		case eventgrid.NumberLessThanAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			value = strconv.FormatFloat(*f.Value, 'f', -1, 64)
		case eventgrid.NumberInAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			var numbers = []string{}
			for _, f := range *f.Values {
				number := strconv.FormatFloat(f, 'f', -1, 64)
				numbers = append(numbers, number)
			}
			values = utils.FlattenStringSlice(&numbers)
		case eventgrid.NumberNotInAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			var numbers = []string{}
			for _, f := range *f.Values {
				number := strconv.FormatFloat(f, 'f', -1, 64)
				numbers = append(numbers, number)
			}
			values = utils.FlattenStringSlice(&numbers)
		case eventgrid.StringBeginsWithAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			values = utils.FlattenStringSlice(f.Values)
		case eventgrid.StringContainsAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			values = utils.FlattenStringSlice(f.Values)
		case eventgrid.StringEndsWithAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			values = utils.FlattenStringSlice(f.Values)
		case eventgrid.StringInAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			values = utils.FlattenStringSlice(f.Values)
		case eventgrid.StringNotInAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			values = utils.FlattenStringSlice(f.Values)
		case eventgrid.NumberLessThanOrEqualsAdvancedFilter:
			key = *f.Key
			operatorType = f.OperatorType
			value = strconv.FormatFloat(*f.Value, 'f', -1, 64)
		}

		results = append(results, map[string]interface{}{
			"key":           key,
			"operator_type": operatorType,
			"value":         value,
			"values":        values,
		})
	}

	return results
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
