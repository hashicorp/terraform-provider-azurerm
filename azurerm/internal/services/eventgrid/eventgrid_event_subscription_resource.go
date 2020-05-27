package eventgrid

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2020-04-01-preview/eventgrid"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventgrid/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func enpointPropertyNames() []string {
	return []string{
		"eventhub_endpoint",
		"eventhub_endpoint_id",
		"hybrid_connection_endpoint",
		"hybrid_connection_endpoint_id",
		"service_bus_queue_endpoint_id",
		"service_bus_topic_endpoint_id",
		"storage_queue_endpoint",
		"webhook_endpoint",
	}
}

func resourceArmEventGridEventSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventGridEventSubscriptionCreateUpdate,
		Read:   resourceArmEventGridEventSubscriptionRead,
		Update: resourceArmEventGridEventSubscriptionCreateUpdate,
		Delete: resourceArmEventGridEventSubscriptionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.EventGridEventSubscriptionID(id)
			return err
		}),

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

			"topic_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"eventhub_endpoint_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: utils.RemoveFromStringArray(enpointPropertyNames(), "eventhub_endpoint_id"),
				ValidateFunc:  azure.ValidateResourceID,
			},

			"eventhub_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Deprecated:    "Deprecated in favour of `" + "eventhub_endpoint_id" + "`",
				Optional:      true,
				Computed:      true,
				ConflictsWith: utils.RemoveFromStringArray(enpointPropertyNames(), "eventhub_endpoint"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eventhub_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"hybrid_connection_endpoint_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: utils.RemoveFromStringArray(enpointPropertyNames(), "hybrid_connection_endpoint_id"),
				ValidateFunc:  azure.ValidateResourceID,
			},

			"hybrid_connection_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Deprecated:    "Deprecated in favour of `" + "hybrid_connection_endpoint_id" + "`",
				Optional:      true,
				Computed:      true,
				ConflictsWith: utils.RemoveFromStringArray(enpointPropertyNames(), "hybrid_connection_endpoint"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hybrid_connection_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"service_bus_queue_endpoint_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: utils.RemoveFromStringArray(enpointPropertyNames(), "service_bus_queue_endpoint_id"),
				ValidateFunc:  azure.ValidateResourceID,
			},

			"service_bus_topic_endpoint_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: utils.RemoveFromStringArray(enpointPropertyNames(), "service_bus_topic_endpoint_id"),
				ValidateFunc:  azure.ValidateResourceID,
			},

			"storage_queue_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: utils.RemoveFromStringArray(enpointPropertyNames(), "storage_queue_endpoint"),
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

			"webhook_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: utils.RemoveFromStringArray(enpointPropertyNames(), "webhook_endpoint"),
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
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"value": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						"operator_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(eventgrid.OperatorTypeBoolEquals),
								string(eventgrid.OperatorTypeNumberGreaterThan),
								string(eventgrid.OperatorTypeNumberGreaterThanOrEquals),
								string(eventgrid.OperatorTypeNumberIn),
								string(eventgrid.OperatorTypeNumberLessThan),
								string(eventgrid.OperatorTypeNumberLessThanOrEquals),
								string(eventgrid.OperatorTypeNumberNotIn),
								string(eventgrid.OperatorTypeStringBeginsWith),
								string(eventgrid.OperatorTypeStringContains),
								string(eventgrid.OperatorTypeStringEndsWith),
								string(eventgrid.OperatorTypeStringIn),
								string(eventgrid.OperatorTypeStringNotIn),
							}, false),
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
		return fmt.Errorf("One of the following endpoint types must be specificed to create an EventGrid Event Subscription: %q", enpointPropertyNames())
	}

	filter := expandEventGridEventSubscriptionFilter(d)

	advancedFilters, err := expandEventGridEventSubscriptionAdvancedFilter(d)
	if err != nil {
		return fmt.Errorf("Error creating/updating EventGrid Event Subscription %q (Scope %q): %s Advanced Filters", name, scope, err)
	}
	if advancedFilters != nil {
		filter.AdvancedFilters = advancedFilters
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

	return resourceArmEventGridEventSubscriptionRead(d, meta)
}

func resourceArmEventGridEventSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventGridEventSubscriptionID(d.Id())
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

		if props.Topic != nil && *props.Topic != "" {
			d.Set("topic_name", props.Topic)
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
		if _, ok := props.Destination.AsWebHookEventSubscriptionDestination(); ok {
			fullURL, err := client.GetFullURL(ctx, id.Scope, id.Name)
			if err != nil {
				return fmt.Errorf("Error making Read request on EventGrid Event Subscription full URL '%s': %+v", id.Name, err)
			}
			if err := d.Set("webhook_endpoint", flattenEventGridEventSubscriptionWebhookEndpoint(&fullURL)); err != nil {
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

func resourceArmEventGridEventSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.EventSubscriptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventGridEventSubscriptionID(d.Id())
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

func expandEventGridExpirationTime(d *schema.ResourceData) (*date.Time, error) {
	if expirationTimeUtc, ok := d.GetOk("expiration_time_utc"); ok {
		if expirationTimeUtc == "" {
			return nil, nil
		}

		parsedExpirationTimeUtc, err := date.ParseTime(time.RFC3339, expirationTimeUtc.(string))
		if err != nil {
			return nil, err
		}

		return &date.Time{Time: parsedExpirationTimeUtc}, nil
	}

	return nil, nil
}

func expandEventGridEventSubscriptionDestination(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	if v, ok := d.GetOk("eventhub_endpoint_id"); ok {
		return &eventgrid.EventHubEventSubscriptionDestination{
			EndpointType: eventgrid.EndpointTypeEventHub,
			EventHubEventSubscriptionDestinationProperties: &eventgrid.EventHubEventSubscriptionDestinationProperties{
				ResourceID: utils.String(v.(string)),
			},
		}
	} else if _, ok := d.GetOk("eventhub_endpoint"); ok {
		return expandEventGridEventSubscriptionEventhubEndpoint(d)
	}

	if v, ok := d.GetOk("hybrid_connection_endpoint_id"); ok {
		return &eventgrid.HybridConnectionEventSubscriptionDestination{
			EndpointType: eventgrid.EndpointTypeHybridConnection,
			HybridConnectionEventSubscriptionDestinationProperties: &eventgrid.HybridConnectionEventSubscriptionDestinationProperties{
				ResourceID: utils.String(v.(string)),
			},
		}
	} else if _, ok := d.GetOk("hybrid_connection_endpoint"); ok {
		return expandEventGridEventSubscriptionHybridConnectionEndpoint(d)
	}

	if v, ok := d.GetOk("service_bus_queue_endpoint_id"); ok {
		return &eventgrid.ServiceBusQueueEventSubscriptionDestination{
			EndpointType: eventgrid.EndpointTypeServiceBusQueue,
			ServiceBusQueueEventSubscriptionDestinationProperties: &eventgrid.ServiceBusQueueEventSubscriptionDestinationProperties{
				ResourceID: utils.String(v.(string)),
			},
		}
	}

	if v, ok := d.GetOk("service_bus_topic_endpoint_id"); ok {
		return &eventgrid.ServiceBusTopicEventSubscriptionDestination{
			EndpointType: eventgrid.EndpointTypeServiceBusTopic,
			ServiceBusTopicEventSubscriptionDestinationProperties: &eventgrid.ServiceBusTopicEventSubscriptionDestinationProperties{
				ResourceID: utils.String(v.(string)),
			},
		}
	}

	if _, ok := d.GetOk("storage_queue_endpoint"); ok {
		return expandEventGridEventSubscriptionStorageQueueEndpoint(d)
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

	return eventgrid.StorageQueueEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeStorageQueue,
		StorageQueueEventSubscriptionDestinationProperties: &eventgrid.StorageQueueEventSubscriptionDestinationProperties{
			ResourceID: &storageAccountID,
			QueueName:  &queueName,
		},
	}
}

func expandEventGridEventSubscriptionEventhubEndpoint(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	props := d.Get("eventhub_endpoint").([]interface{})[0].(map[string]interface{})
	eventHubID := props["eventhub_id"].(string)

	return eventgrid.EventHubEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeEventHub,
		EventHubEventSubscriptionDestinationProperties: &eventgrid.EventHubEventSubscriptionDestinationProperties{
			ResourceID: &eventHubID,
		},
	}
}

func expandEventGridEventSubscriptionHybridConnectionEndpoint(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	props := d.Get("hybrid_connection_endpoint").([]interface{})[0].(map[string]interface{})
	hybridConnectionID := props["hybrid_connection_id"].(string)

	return eventgrid.HybridConnectionEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeHybridConnection,
		HybridConnectionEventSubscriptionDestinationProperties: &eventgrid.HybridConnectionEventSubscriptionDestinationProperties{
			ResourceID: &hybridConnectionID,
		},
	}
}

func expandEventGridEventSubscriptionWebhookEndpoint(d *schema.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	props := d.Get("webhook_endpoint").([]interface{})[0].(map[string]interface{})
	url := props["url"].(string)

	return eventgrid.WebHookEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeWebHook,
		WebHookEventSubscriptionDestinationProperties: &eventgrid.WebHookEventSubscriptionDestinationProperties{
			EndpointURL: &url,
		},
	}
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

func expandEventGridEventSubscriptionAdvancedFilter(d *schema.ResourceData) (*[]eventgrid.BasicAdvancedFilter, error) {
	advFilters := d.Get("advanced_filter").([]interface{})
	advancedFilters := make([]eventgrid.BasicAdvancedFilter, 0, len(advFilters))

	for _, advFilter := range advFilters {
		advfilterconfig := advFilter.(map[string]interface{})
		key := advfilterconfig["key"].(string)
		operatorType := advfilterconfig["operator_type"].(string)
		value := advfilterconfig["value"].(string)
		values := utils.ExpandStringSlice(advfilterconfig["values"].([]interface{}))

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeBoolEquals)) == 0 {
			boolEquals := &eventgrid.BoolEqualsAdvancedFilter{}
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return nil, err
			}
			boolEquals.Value = &boolValue
			boolEquals.Key = &key
			boolEquals.OperatorType = eventgrid.OperatorTypeBoolEquals
			advancedFilters = append(advancedFilters, boolEquals)
		}

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeNumberGreaterThan)) == 0 {
			numberGreaterThan := &eventgrid.NumberGreaterThanAdvancedFilter{}
			numberValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			numberGreaterThan.Value = &numberValue
			numberGreaterThan.Key = &key
			numberGreaterThan.OperatorType = eventgrid.OperatorTypeNumberGreaterThan
			advancedFilters = append(advancedFilters, numberGreaterThan)
		}

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeNumberGreaterThanOrEquals)) == 0 {
			numberGreaterThanEquals := &eventgrid.NumberGreaterThanOrEqualsAdvancedFilter{}
			numberValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			numberGreaterThanEquals.Value = &numberValue
			numberGreaterThanEquals.Key = &key
			numberGreaterThanEquals.OperatorType = eventgrid.OperatorTypeNumberGreaterThanOrEquals
			advancedFilters = append(advancedFilters, numberGreaterThanEquals)
		}

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeNumberIn)) == 0 {
			numberIn := &eventgrid.NumberInAdvancedFilter{}
			floatValues, err := sliceAtof(*values)
			if err != nil {
				return nil, err
			}
			numberIn.Values = &floatValues
			numberIn.Key = &key
			numberIn.OperatorType = eventgrid.OperatorTypeNumberIn
			advancedFilters = append(advancedFilters, numberIn)
		}

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeNumberLessThan)) == 0 {
			numberLessThan := &eventgrid.NumberLessThanAdvancedFilter{}
			numberValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			numberLessThan.Value = &numberValue
			numberLessThan.Key = &key
			numberLessThan.OperatorType = eventgrid.OperatorTypeNumberLessThan
			advancedFilters = append(advancedFilters, numberLessThan)
		}

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeNumberLessThanOrEquals)) == 0 {
			numberLessThanEquals := &eventgrid.NumberLessThanOrEqualsAdvancedFilter{}
			numberValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			numberLessThanEquals.Value = &numberValue
			numberLessThanEquals.Key = &key
			numberLessThanEquals.OperatorType = eventgrid.OperatorTypeNumberLessThanOrEquals
			advancedFilters = append(advancedFilters, numberLessThanEquals)
		}

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeNumberNotIn)) == 0 {
			numberNotIn := &eventgrid.NumberNotInAdvancedFilter{}
			floatValues, err := sliceAtof(*values)
			if err != nil {
				return nil, err
			}
			numberNotIn.Values = &floatValues
			numberNotIn.Key = &key
			numberNotIn.OperatorType = eventgrid.OperatorTypeNumberNotIn
			advancedFilters = append(advancedFilters, numberNotIn)
		}

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeStringBeginsWith)) == 0 {
			stringBeginsWith := &eventgrid.StringBeginsWithAdvancedFilter{}
			stringBeginsWith.Values = values
			stringBeginsWith.Key = &key
			stringBeginsWith.OperatorType = eventgrid.OperatorTypeStringBeginsWith
			advancedFilters = append(advancedFilters, stringBeginsWith)
		}

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeStringContains)) == 0 {
			stringContains := &eventgrid.StringContainsAdvancedFilter{}
			stringContains.Values = values
			stringContains.Key = &key
			stringContains.OperatorType = eventgrid.OperatorTypeStringContains
			advancedFilters = append(advancedFilters, stringContains)
		}

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeStringEndsWith)) == 0 {
			stringEndsWith := &eventgrid.StringEndsWithAdvancedFilter{}
			stringEndsWith.Values = values
			stringEndsWith.Key = &key
			stringEndsWith.OperatorType = eventgrid.OperatorTypeStringEndsWith
			advancedFilters = append(advancedFilters, stringEndsWith)
		}

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeStringIn)) == 0 {
			stringIn := &eventgrid.StringEndsWithAdvancedFilter{}
			stringIn.Values = values
			stringIn.Key = &key
			stringIn.OperatorType = eventgrid.OperatorTypeStringIn
			advancedFilters = append(advancedFilters, stringIn)
		}

		if strings.Compare(operatorType, string(eventgrid.OperatorTypeStringNotIn)) == 0 {
			stringNotIn := &eventgrid.StringEndsWithAdvancedFilter{}
			stringNotIn.Values = values
			stringNotIn.Key = &key
			stringNotIn.OperatorType = eventgrid.OperatorTypeStringNotIn
			advancedFilters = append(advancedFilters, stringNotIn)
		}
	}
	return &advancedFilters, nil
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

func flattenEventGridEventSubscriptionEventhubEndpoint(input *eventgrid.EventHubEventSubscriptionDestination) []interface{} {
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

	return []interface{}{result}
}

func flattenEventGridEventSubscriptionAdvancedFilter(filter *eventgrid.EventSubscriptionFilter) []interface{} {
	if filter.AdvancedFilters == nil {
		return nil
	}

	filterResult := make([]interface{}, 0, len(*filter.AdvancedFilters))
	for _, advancedFilter := range *filter.AdvancedFilters {
		advFilter := make(map[string]interface{})
		if boolEqualsFilter, _ := advancedFilter.AsBoolEqualsAdvancedFilter(); boolEqualsFilter != nil {
			advFilter["key"] = boolEqualsFilter.Key
			advFilter["operator_type"] = boolEqualsFilter.OperatorType
			advFilter["value"] = strconv.FormatBool(*boolEqualsFilter.Value)
		}

		if numberGreaterThanFilter, _ := advancedFilter.AsNumberGreaterThanAdvancedFilter(); numberGreaterThanFilter != nil {
			advFilter["key"] = numberGreaterThanFilter.Key
			advFilter["operator_type"] = numberGreaterThanFilter.OperatorType
			advFilter["value"] = strconv.FormatFloat(*numberGreaterThanFilter.Value, 'f', 0, 64)
		}

		if numberGreaterThanOrEqualsFilter, _ := advancedFilter.AsNumberGreaterThanOrEqualsAdvancedFilter(); numberGreaterThanOrEqualsFilter != nil {
			advFilter["key"] = numberGreaterThanOrEqualsFilter.Key
			advFilter["operator_type"] = numberGreaterThanOrEqualsFilter.OperatorType
			advFilter["value"] = strconv.FormatFloat(*numberGreaterThanOrEqualsFilter.Value, 'f', 0, 64)
		}

		if numberInFilter, _ := advancedFilter.AsNumberInAdvancedFilter(); numberInFilter != nil {
			advFilter["key"] = numberInFilter.Key
			advFilter["operator_type"] = numberInFilter.OperatorType
			advFilter["values"] = sliceFtoa(*numberInFilter.Values)
		}

		if numberLessThanFilter, _ := advancedFilter.AsNumberLessThanAdvancedFilter(); numberLessThanFilter != nil {
			advFilter["key"] = numberLessThanFilter.Key
			advFilter["operator_type"] = numberLessThanFilter.OperatorType
			advFilter["value"] = strconv.FormatFloat(*numberLessThanFilter.Value, 'f', 0, 64)
		}

		if numberLessThanOrEqualsFilter, _ := advancedFilter.AsNumberLessThanOrEqualsAdvancedFilter(); numberLessThanOrEqualsFilter != nil {
			advFilter["key"] = numberLessThanOrEqualsFilter.Key
			advFilter["operator_type"] = numberLessThanOrEqualsFilter.OperatorType
			advFilter["value"] = strconv.FormatFloat(*numberLessThanOrEqualsFilter.Value, 'f', 0, 64)
		}

		if numberNotInFilter, _ := advancedFilter.AsNumberNotInAdvancedFilter(); numberNotInFilter != nil {
			advFilter["key"] = numberNotInFilter.Key
			advFilter["operator_type"] = numberNotInFilter.OperatorType
			advFilter["values"] = sliceFtoa(*numberNotInFilter.Values)
		}

		if stringBeginsWithFilter, _ := advancedFilter.AsStringBeginsWithAdvancedFilter(); stringBeginsWithFilter != nil {
			advFilter["key"] = stringBeginsWithFilter.Key
			advFilter["operator_type"] = stringBeginsWithFilter.OperatorType
			advFilter["values"] = stringBeginsWithFilter.Values
		}

		if stringContainsFilter, _ := advancedFilter.AsStringContainsAdvancedFilter(); stringContainsFilter != nil {
			advFilter["key"] = stringContainsFilter.Key
			advFilter["operator_type"] = stringContainsFilter.OperatorType
			advFilter["values"] = stringContainsFilter.Values
		}

		if stringEndsWithFilter, _ := advancedFilter.AsStringEndsWithAdvancedFilter(); stringEndsWithFilter != nil {
			advFilter["key"] = stringEndsWithFilter.Key
			advFilter["operator_type"] = stringEndsWithFilter.OperatorType
			advFilter["values"] = stringEndsWithFilter.Values
		}

		if stringInFilter, _ := advancedFilter.AsStringInAdvancedFilter(); stringInFilter != nil {
			advFilter["key"] = stringInFilter.Key
			advFilter["operator_type"] = stringInFilter.OperatorType
			advFilter["values"] = stringInFilter.Values
		}

		filterResult = append(filterResult, advFilter)
	}

	return filterResult
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

func sliceAtof(strvalues []string) ([]float64, error) {
	floatvalues := make([]float64, 0, len(strvalues))
	for _, a := range strvalues {
		i, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return floatvalues, err
		}
		floatvalues = append(floatvalues, i)
	}
	return floatvalues, nil
}

func sliceFtoa(floatvalues []float64) []string {
	valuesText := []string{}

	for _, number := range floatvalues {
		text := strconv.FormatFloat(number, 'f', 0, 64)
		valuesText = append(valuesText, text)
	}
	return valuesText
}
