package eventgrid

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2020-04-01-preview/eventgrid"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventgrid/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func enpointPropertyNames() []string {
	return []string{
		"azure_function_endpoint",
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
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "This field has been updated to readonly field since Apr 25, 2019 so no longer has any affect and will be removed in version 3.0 of the provider.",
			},

			"azure_function_endpoint": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: utils.RemoveFromStringArray(enpointPropertyNames(), "azure_function_endpoint"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"max_events_per_batch": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"preferred_batch_size_in_kilobytes": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
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
						"base_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_events_per_batch": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 5000),
						},
						"preferred_batch_size_in_kilobytes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 1024),
						},
						"active_directory_tenant_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"active_directory_app_id_or_uri": {
							Type:     schema.TypeString,
							Optional: true,
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
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bool_equals": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"value": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"number_greater_than": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"value": {
										Type:     schema.TypeFloat,
										Required: true,
									},
								},
							},
						},
						"number_greater_than_or_equals": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"value": {
										Type:     schema.TypeFloat,
										Required: true,
									},
								},
							}},
						"number_less_than": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"value": {
										Type:     schema.TypeFloat,
										Required: true,
									},
								},
							},
						},
						"number_less_than_or_equals": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"value": {
										Type:     schema.TypeFloat,
										Required: true,
									},
								},
							},
						},
						"number_in": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 5,
										Elem: &schema.Schema{
											Type: schema.TypeFloat,
										},
									},
								},
							},
						},
						"number_not_in": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 5,
										Elem: &schema.Schema{
											Type: schema.TypeFloat,
										},
									},
								},
							},
						},
						"string_begins_with": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 5,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"string_ends_with": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 5,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"string_contains": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 5,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"string_in": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 5,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"string_not_in": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 5,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
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
		},
	}
}

func resourceArmEventGridEventSubscriptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("One of the following endpoint types must be specificed to create an EventGrid Event Subscription: %q", enpointPropertyNames())
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
	if v, ok := d.GetOk("azure_function_endpoint"); ok {
		return expandEventGridEventSubscriptionAzureFunctionEndpoint(v)
	}

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

	if v, ok := d.GetOk("webhook_endpoint"); ok {
		return expandEventGridEventSubscriptionWebhookEndpoint(v)
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

func expandEventGridEventSubscriptionAzureFunctionEndpoint(input interface{}) eventgrid.BasicEventSubscriptionDestination {
	configs := input.([]interface{})

	props := eventgrid.AzureFunctionEventSubscriptionDestinationProperties{}
	azureFunctionDestination := &eventgrid.AzureFunctionEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeAzureFunction,
		AzureFunctionEventSubscriptionDestinationProperties: &props,
	}

	if len(configs) == 0 {
		return azureFunctionDestination
	}

	config := configs[0].(map[string]interface{})

	if v, ok := config["function_id"]; ok && v != "" {
		props.ResourceID = utils.String(v.(string))
	}

	if v, ok := config["max_events_per_batch"]; ok && v != 0 {
		props.MaxEventsPerBatch = utils.Int32(int32(v.(int)))
	}

	if v, ok := config["preferred_batch_size_in_kilobytes"]; ok && v != 0 {
		props.PreferredBatchSizeInKilobytes = utils.Int32(int32(v.(int)))
	}

	return azureFunctionDestination
}

func expandEventGridEventSubscriptionWebhookEndpoint(input interface{}) eventgrid.BasicEventSubscriptionDestination {
	configs := input.([]interface{})

	props := eventgrid.WebHookEventSubscriptionDestinationProperties{}
	webhookDestination := &eventgrid.WebHookEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeWebHook,
		WebHookEventSubscriptionDestinationProperties: &props,
	}

	if len(configs) == 0 {
		return webhookDestination
	}

	config := configs[0].(map[string]interface{})

	if v, ok := config["url"]; ok && v != "" {
		props.EndpointURL = utils.String(v.(string))
	}

	if v, ok := config["max_events_per_batch"]; ok && v != 0 {
		props.MaxEventsPerBatch = utils.Int32(int32(v.(int)))
	}

	if v, ok := config["preferred_batch_size_in_kilobytes"]; ok && v != 0 {
		props.PreferredBatchSizeInKilobytes = utils.Int32(int32(v.(int)))
	}

	if v, ok := config["active_directory_tenant_id"]; ok && v != "" {
		props.AzureActiveDirectoryTenantID = utils.String(v.(string))
	}

	if v, ok := config["active_directory_app_id_or_uri"]; ok && v != "" {
		props.AzureActiveDirectoryApplicationIDOrURI = utils.String(v.(string))
	}

	return webhookDestination
}

func expandEventGridEventSubscriptionFilter(d *schema.ResourceData) (*eventgrid.EventSubscriptionFilter, error) {
	filter := &eventgrid.EventSubscriptionFilter{}

	if includedEvents, ok := d.GetOk("included_event_types"); ok {
		filter.IncludedEventTypes = utils.ExpandStringSlice(includedEvents.([]interface{}))
	}

	if v, ok := d.GetOk("subject_filter"); ok {
		if v.([]interface{})[0] != nil {
			config := v.([]interface{})[0].(map[string]interface{})
			subjectBeginsWith := config["subject_begins_with"].(string)
			subjectEndsWith := config["subject_ends_with"].(string)
			caseSensitive := config["case_sensitive"].(bool)

			filter.SubjectBeginsWith = &subjectBeginsWith
			filter.SubjectEndsWith = &subjectEndsWith
			filter.IsSubjectCaseSensitive = &caseSensitive
		}
	}

	if advancedFilter, ok := d.GetOk("advanced_filter"); ok {
		advancedFilters := make([]eventgrid.BasicAdvancedFilter, 0)
		for filterKey, filterSchema := range advancedFilter.([]interface{})[0].(map[string]interface{}) {
			for _, options := range filterSchema.([]interface{}) {
				if filter, err := expandAdvancedFilter(filterKey, options.(map[string]interface{})); err == nil {
					advancedFilters = append(advancedFilters, filter)
				} else {
					return nil, err
				}
			}
		}
		filter.AdvancedFilters = &advancedFilters
	}

	return filter, nil
}

func expandAdvancedFilter(operatorType string, config map[string]interface{}) (eventgrid.BasicAdvancedFilter, error) {
	k := config["key"].(string)

	switch operatorType {
	case "bool_equals":
		v := config["value"].(bool)
		return eventgrid.BoolEqualsAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeBoolEquals, Value: &v}, nil
	case "number_greater_than":
		v := config["value"].(float64)
		return eventgrid.NumberGreaterThanAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeNumberGreaterThan, Value: &v}, nil
	case "number_greater_than_or_equals":
		v := config["value"].(float64)
		return eventgrid.NumberGreaterThanOrEqualsAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeNumberGreaterThanOrEquals, Value: &v}, nil
	case "number_less_than":
		v := config["value"].(float64)
		return eventgrid.NumberLessThanAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeNumberLessThan, Value: &v}, nil
	case "number_less_than_or_equals":
		v := config["value"].(float64)
		return eventgrid.NumberLessThanOrEqualsAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeNumberLessThanOrEquals, Value: &v}, nil
	case "number_in":
		v := utils.ExpandFloatSlice(config["values"].([]interface{}))
		return eventgrid.NumberInAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeNumberIn, Values: v}, nil
	case "number_not_in":
		v := utils.ExpandFloatSlice(config["values"].([]interface{}))
		return eventgrid.NumberNotInAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeNumberIn, Values: v}, nil
	case "string_begins_with":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringBeginsWithAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringBeginsWith, Values: v}, nil
	case "string_ends_with":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringEndsWithAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringEndsWith, Values: v}, nil
	case "string_contains":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringContainsAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringContains, Values: v}, nil
	case "string_in":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringInAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringIn, Values: v}, nil
	case "string_not_in":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringNotInAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringNotIn, Values: v}, nil
	default:
		return nil, fmt.Errorf("Invalid `advanced_filter` operator_type %q used", operatorType)
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

	hybridConnectionId := ""
	if input.ResourceID != nil {
		hybridConnectionId = *input.ResourceID
	}

	return []interface{}{
		map[string]interface{}{
			"hybrid_connection_id": hybridConnectionId,
		},
	}
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

func flattenEventGridEventSubscriptionAzureFunctionEndpoint(input *eventgrid.AzureFunctionEventSubscriptionDestination) []interface{} {
	results := make([]interface{}, 0)

	if input == nil {
		return results
	}

	functionID := ""
	if input.ResourceID != nil {
		functionID = *input.ResourceID
	}

	maxEventsPerBatch := 0
	if input.MaxEventsPerBatch != nil {
		maxEventsPerBatch = int(*input.MaxEventsPerBatch)
	}

	preferredBatchSize := 0
	if input.PreferredBatchSizeInKilobytes != nil {
		preferredBatchSize = int(*input.PreferredBatchSizeInKilobytes)
	}

	return append(results, map[string]interface{}{
		"function_id":                       functionID,
		"max_events_per_batch":              maxEventsPerBatch,
		"preferred_batch_size_in_kilobytes": preferredBatchSize,
	})
}

func flattenEventGridEventSubscriptionWebhookEndpoint(input *eventgrid.WebHookEventSubscriptionDestination, fullURL *eventgrid.EventSubscriptionFullURL) []interface{} {
	results := make([]interface{}, 0)

	if input == nil {
		return results
	}

	webhookURL := ""
	if fullURL != nil {
		webhookURL = *fullURL.EndpointURL
	}

	webhookBaseURL := ""
	if input.EndpointBaseURL != nil {
		webhookBaseURL = *input.EndpointBaseURL
	}

	maxEventsPerBatch := 0
	if input.MaxEventsPerBatch != nil {
		maxEventsPerBatch = int(*input.MaxEventsPerBatch)
	}

	preferredBatchSizeInKilobytes := 0
	if input.PreferredBatchSizeInKilobytes != nil {
		preferredBatchSizeInKilobytes = int(*input.PreferredBatchSizeInKilobytes)
	}

	azureActiveDirectoryTenantID := ""
	if input.AzureActiveDirectoryTenantID != nil {
		azureActiveDirectoryTenantID = *input.AzureActiveDirectoryTenantID
	}

	azureActiveDirectoryApplicationIDOrURI := ""
	if input.AzureActiveDirectoryApplicationIDOrURI != nil {
		azureActiveDirectoryApplicationIDOrURI = *input.AzureActiveDirectoryApplicationIDOrURI
	}

	return append(results, map[string]interface{}{
		"url":                               webhookURL,
		"base_url":                          webhookBaseURL,
		"max_events_per_batch":              maxEventsPerBatch,
		"preferred_batch_size_in_kilobytes": preferredBatchSizeInKilobytes,
		"active_directory_tenant_id":        azureActiveDirectoryTenantID,
		"active_directory_app_id_or_uri":    azureActiveDirectoryApplicationIDOrURI,
	})
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

func flattenEventGridEventSubscriptionAdvancedFilter(input *eventgrid.EventSubscriptionFilter) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.AdvancedFilters == nil {
		return results
	}

	boolEquals := make([]interface{}, 0)
	numberGreaterThan := make([]interface{}, 0)
	numberGreaterThanOrEquals := make([]interface{}, 0)
	numberLessThan := make([]interface{}, 0)
	numberLessThanOrEquals := make([]interface{}, 0)
	numberIn := make([]interface{}, 0)
	numberNotIn := make([]interface{}, 0)
	stringBeginsWith := make([]interface{}, 0)
	stringEndsWith := make([]interface{}, 0)
	stringContains := make([]interface{}, 0)
	stringIn := make([]interface{}, 0)
	stringNotIn := make([]interface{}, 0)

	for _, item := range *input.AdvancedFilters {
		switch f := item.(type) {
		case eventgrid.BoolEqualsAdvancedFilter:
			v := interface{}(f.Value)
			boolEquals = append(boolEquals, flattenValue(f.Key, &v))
		case eventgrid.NumberGreaterThanAdvancedFilter:
			v := interface{}(f.Value)
			numberGreaterThan = append(numberGreaterThan, flattenValue(f.Key, &v))
		case eventgrid.NumberGreaterThanOrEqualsAdvancedFilter:
			v := interface{}(f.Value)
			numberGreaterThanOrEquals = append(numberGreaterThanOrEquals, flattenValue(f.Key, &v))
		case eventgrid.NumberLessThanAdvancedFilter:
			v := interface{}(f.Value)
			numberLessThan = append(numberLessThan, flattenValue(f.Key, &v))
		case eventgrid.NumberLessThanOrEqualsAdvancedFilter:
			v := interface{}(f.Value)
			numberLessThanOrEquals = append(numberLessThanOrEquals, flattenValue(f.Key, &v))
		case eventgrid.NumberInAdvancedFilter:
			v := utils.FlattenFloatSlice(f.Values)
			numberIn = append(numberIn, flattenValues(f.Key, &v))
		case eventgrid.NumberNotInAdvancedFilter:
			v := utils.FlattenFloatSlice(f.Values)
			numberNotIn = append(numberNotIn, flattenValues(f.Key, &v))
		case eventgrid.StringBeginsWithAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringBeginsWith = append(stringBeginsWith, flattenValues(f.Key, &v))
		case eventgrid.StringEndsWithAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringEndsWith = append(stringEndsWith, flattenValues(f.Key, &v))
		case eventgrid.StringContainsAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringContains = append(stringContains, flattenValues(f.Key, &v))
		case eventgrid.StringInAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringIn = append(stringIn, flattenValues(f.Key, &v))
		case eventgrid.StringNotInAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringNotIn = append(stringNotIn, flattenValues(f.Key, &v))
		}
	}

	return []interface{}{
		map[string][]interface{}{
			"bool_equals":                   boolEquals,
			"number_greater_than":           numberGreaterThan,
			"number_greater_than_or_equals": numberGreaterThanOrEquals,
			"number_less_than":              numberLessThan,
			"number_less_than_or_equals":    numberLessThanOrEquals,
			"number_in":                     numberIn,
			"number_not_in":                 numberNotIn,
			"string_begins_with":            stringBeginsWith,
			"string_ends_with":              stringEndsWith,
			"string_contains":               stringContains,
			"string_in":                     stringIn,
			"string_not_in":                 stringNotIn,
		},
	}
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

func flattenValue(inputKey *string, inputValue *interface{}) map[string]interface{} {
	key := ""
	if inputKey != nil {
		key = *inputKey
	}
	var value interface{}
	if inputValue != nil {
		value = inputValue
	}

	return map[string]interface{}{
		"key":   key,
		"value": value,
	}
}

func flattenValues(inputKey *string, inputValues *[]interface{}) map[string]interface{} {
	key := ""
	if inputKey != nil {
		key = *inputKey
	}
	values := make([]interface{}, 0)
	if inputValues != nil {
		values = *inputValues
	}

	return map[string]interface{}{
		"key":    key,
		"values": values,
	}
}
