package eventgrid

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2021-12-01/eventgrid"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// EventSubscriptionEndpointType enumerates the values for event subscription endpoint types.
type EventSubscriptionEndpointType string

const (
	// AzureFunctionEndpoint ...
	AzureFunctionEndpoint EventSubscriptionEndpointType = "azure_function_endpoint"
	// EventHubEndpoint ...
	EventHubEndpoint EventSubscriptionEndpointType = "eventhub_endpoint"
	// EventHubEndpointID ...
	EventHubEndpointID EventSubscriptionEndpointType = "eventhub_endpoint_id"
	// HybridConnectionEndpoint ...
	HybridConnectionEndpoint EventSubscriptionEndpointType = "hybrid_connection_endpoint"
	// HybridConnectionEndpointID ...
	HybridConnectionEndpointID EventSubscriptionEndpointType = "hybrid_connection_endpoint_id"
	// ServiceBusQueueEndpointID ...
	ServiceBusQueueEndpointID EventSubscriptionEndpointType = "service_bus_queue_endpoint_id"
	// ServiceBusTopicEndpointID ...
	ServiceBusTopicEndpointID EventSubscriptionEndpointType = "service_bus_topic_endpoint_id"
	// StorageQueueEndpoint ...
	StorageQueueEndpoint EventSubscriptionEndpointType = "storage_queue_endpoint"
	// WebHookEndpoint ...
	WebHookEndpoint EventSubscriptionEndpointType = "webhook_endpoint"
)

func eventSubscriptionCustomizeDiffAdvancedFilter(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
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
}

func eventSubscriptionSchemaEventSubscriptionName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.All(
			validation.StringIsNotEmpty,
			validation.StringMatch(
				regexp.MustCompile("^[-a-zA-Z0-9]{3,64}$"),
				"EventGrid subscription name must be 3 - 64 characters long, contain only letters, numbers and hyphens.",
			),
		),
	}
}

func eventSubscriptionSchemaEventDeliverySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		ForceNew: true,
		Default:  string(eventgrid.EventDeliverySchemaEventGridSchema),
		ValidateFunc: validation.StringInSlice([]string{
			string(eventgrid.EventDeliverySchemaEventGridSchema),
			string(eventgrid.EventDeliverySchemaCloudEventSchemaV10),
			string(eventgrid.EventDeliverySchemaCustomInputSchema),
		}, false),
	}
}

func eventSubscriptionSchemaDeliveryProperty() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"header_name": {
					Type:             pluginsdk.TypeString,
					Required:         true,
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"Static",
						"Dynamic",
					}, false),
				},

				"value": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					Sensitive: true,
				},

				"source_field": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"secret": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func eventSubscriptionSchemaExpirationTimeUTC() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
}

func eventSubscriptionSchemaAzureFunctionEndpoint(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		MaxItems:      1,
		Optional:      true,
		ConflictsWith: conflictsWith,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"function_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
				"max_events_per_batch": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},
				"preferred_batch_size_in_kilobytes": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

func eventSubscriptionSchemaEventHubEndpointID(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeString,
		Optional:      true,
		Computed:      true,
		ConflictsWith: conflictsWith,
		ValidateFunc:  azure.ValidateResourceID,
	}
}

func eventSubscriptionSchemaEventHubEndpoint(conflictsWith []string) *pluginsdk.Schema {
	//lintignore:XS003
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		MaxItems:      1,
		Deprecated:    "Deprecated in favour of `" + "eventhub_endpoint_id" + "`",
		Optional:      true,
		Computed:      true,
		ConflictsWith: conflictsWith,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"eventhub_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
		},
	}
}

func eventSubscriptionSchemaHybridConnectionEndpointID(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeString,
		Optional:      true,
		Computed:      true,
		ConflictsWith: conflictsWith,
		ValidateFunc:  azure.ValidateResourceID,
	}
}

func eventSubscriptionSchemaHybridEndpoint(conflictsWith []string) *pluginsdk.Schema {
	//lintignore:XS003
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		MaxItems:      1,
		Deprecated:    "Deprecated in favour of `" + "hybrid_connection_endpoint_id" + "`",
		Optional:      true,
		Computed:      true,
		ConflictsWith: conflictsWith,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"hybrid_connection_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
		},
	}
}

func eventSubscriptionSchemaServiceBusQueueEndpointID(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeString,
		Optional:      true,
		ConflictsWith: conflictsWith,
		ValidateFunc:  azure.ValidateResourceID,
	}
}

func eventSubscriptionSchemaServiceBusTopicEndpointID(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeString,
		Optional:      true,
		ConflictsWith: conflictsWith,
		ValidateFunc:  azure.ValidateResourceID,
	}
}

func eventSubscriptionSchemaStorageQueueEndpoint(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		MaxItems:      1,
		Optional:      true,
		ConflictsWith: conflictsWith,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"storage_account_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
				"queue_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"queue_message_time_to_live_in_seconds": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

func eventSubscriptionSchemaWebHookEndpoint(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		MaxItems:      1,
		Optional:      true,
		ConflictsWith: conflictsWith,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"url": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsURLWithHTTPS,
				},
				"base_url": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"max_events_per_batch": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 5000),
				},
				"preferred_batch_size_in_kilobytes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 1024),
				},
				"active_directory_tenant_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"active_directory_app_id_or_uri": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func eventSubscriptionSchemaIncludedEventTypes() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func eventSubscriptionSchemaEnableAdvancedFilteringOnArrays() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
}

func eventSubscriptionSchemaSubjectFilter() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"subject_begins_with": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					AtLeastOneOf: []string{"subject_filter.0.subject_begins_with", "subject_filter.0.subject_ends_with", "subject_filter.0.case_sensitive"},
				},
				"subject_ends_with": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					AtLeastOneOf: []string{"subject_filter.0.subject_begins_with", "subject_filter.0.subject_ends_with", "subject_filter.0.case_sensitive"},
				},
				"case_sensitive": {
					Type:         pluginsdk.TypeBool,
					Optional:     true,
					AtLeastOneOf: []string{"subject_filter.0.subject_begins_with", "subject_filter.0.subject_ends_with", "subject_filter.0.case_sensitive"},
				},
			},
		},
	}
}

func eventSubscriptionSchemaAdvancedFilter() *pluginsdk.Schema {
	atLeastOneOf := []string{
		"advanced_filter.0.bool_equals", "advanced_filter.0.number_greater_than", "advanced_filter.0.number_greater_than_or_equals", "advanced_filter.0.number_less_than",
		"advanced_filter.0.number_less_than_or_equals", "advanced_filter.0.number_in", "advanced_filter.0.number_not_in", "advanced_filter.0.string_begins_with", "advanced_filter.0.string_not_begins_with",
		"advanced_filter.0.string_ends_with", "advanced_filter.0.string_not_ends_with", "advanced_filter.0.string_contains", "advanced_filter.0.string_not_contains", "advanced_filter.0.string_in",
		"advanced_filter.0.string_not_in", "advanced_filter.0.is_not_null", "advanced_filter.0.is_null_or_undefined", "advanced_filter.0.number_in_range", "advanced_filter.0.number_not_in_range",
	}
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"bool_equals": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"value": {
								Type:     pluginsdk.TypeBool,
								Required: true,
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"number_greater_than": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"value": {
								Type:     pluginsdk.TypeFloat,
								Required: true,
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"number_greater_than_or_equals": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"value": {
								Type:     pluginsdk.TypeFloat,
								Required: true,
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"number_less_than": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"value": {
								Type:     pluginsdk.TypeFloat,
								Required: true,
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"number_less_than_or_equals": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"value": {
								Type:     pluginsdk.TypeFloat,
								Required: true,
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"number_in": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeFloat,
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"number_not_in": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeFloat,
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"string_begins_with": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"string_not_begins_with": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"string_ends_with": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"string_not_ends_with": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"string_contains": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"string_not_contains": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"string_in": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"string_not_in": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"is_not_null": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"is_null_or_undefined": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"number_in_range": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type:     pluginsdk.TypeList,
									MinItems: 2,
									MaxItems: 2,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeFloat,
									},
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
				"number_not_in_range": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 25,
								Elem: &pluginsdk.Schema{
									Type:     pluginsdk.TypeList,
									MinItems: 2,
									MaxItems: 2,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeFloat,
									},
								},
							},
						},
					},
					AtLeastOneOf: atLeastOneOf,
				},
			},
		},
	}
}

func eventSubscriptionSchemaStorageBlobDeadletterDestination() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"storage_account_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
				"storage_blob_container_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func eventSubscriptionSchemaRetryPolicy() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_delivery_attempts": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 30),
				},
				"event_time_to_live": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 1440),
				},
			},
		},
	}
}

func eventSubscriptionSchemaLabels() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

func eventSubscriptionSchemaIdentity() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(eventgrid.EventSubscriptionIdentityTypeSystemAssigned),
						string(eventgrid.EventSubscriptionIdentityTypeUserAssigned),
					}, false),
				},
				"user_assigned_identity": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
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

func expandEventGridEventSubscriptionDestination(d *pluginsdk.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	if _, ok := d.GetOk("azure_function_endpoint"); ok {
		return expandEventGridEventSubscriptionAzureFunctionEndpoint(d)
	}

	if _, ok := d.GetOk("eventhub_endpoint_id"); ok {
		return expandEventGridEventSubscriptionEventhubEndpoint(d)
	} else if _, ok := d.GetOk("eventhub_endpoint"); ok {
		return expandEventGridEventSubscriptionEventhubEndpoint(d)
	}

	if _, ok := d.GetOk("hybrid_connection_endpoint_id"); ok {
		return expandEventGridEventSubscriptionHybridConnectionEndpoint(d)
	} else if _, ok := d.GetOk("hybrid_connection_endpoint"); ok {
		return expandEventGridEventSubscriptionHybridConnectionEndpoint(d)
	}

	if _, ok := d.GetOk("service_bus_queue_endpoint_id"); ok {
		return expandEventGridEventSubscriptionServiceBusQueueEndpoint(d)
	}

	if _, ok := d.GetOk("service_bus_topic_endpoint_id"); ok {
		return expandEventGridEventSubscriptionServiceBusTopicEndpoint(d)
	}

	if _, ok := d.GetOk("storage_queue_endpoint"); ok {
		return expandEventGridEventSubscriptionStorageQueueEndpoint(d)
	}

	if _, ok := d.GetOk("webhook_endpoint"); ok {
		return expandEventGridEventSubscriptionWebhookEndpoint(d)
	}

	return nil
}

func expandEventGridEventSubscriptionServiceBusQueueEndpoint(d *pluginsdk.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	endpoint := d.Get("service_bus_queue_endpoint_id")

	props := &eventgrid.ServiceBusQueueEventSubscriptionDestinationProperties{
		ResourceID: utils.String(endpoint.(string)),
	}

	deliveryMappings := expandDeliveryProperties(d)
	props.DeliveryAttributeMappings = &deliveryMappings

	return &eventgrid.ServiceBusQueueEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeServiceBusQueue,
		ServiceBusQueueEventSubscriptionDestinationProperties: props,
	}
}

func expandEventGridEventSubscriptionServiceBusTopicEndpoint(d *pluginsdk.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	endpoint := d.Get("service_bus_topic_endpoint_id")

	props := &eventgrid.ServiceBusTopicEventSubscriptionDestinationProperties{
		ResourceID: utils.String(endpoint.(string)),
	}

	deliveryMappings := expandDeliveryProperties(d)
	props.DeliveryAttributeMappings = &deliveryMappings

	return &eventgrid.ServiceBusTopicEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeServiceBusTopic,
		ServiceBusTopicEventSubscriptionDestinationProperties: props,
	}
}

func expandEventGridEventSubscriptionStorageQueueEndpoint(d *pluginsdk.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	props := d.Get("storage_queue_endpoint").([]interface{})[0].(map[string]interface{})
	storageAccountID := props["storage_account_id"].(string)
	queueName := props["queue_name"].(string)
	storageQueueEventSubscriptionDestinationProperties := &eventgrid.StorageQueueEventSubscriptionDestinationProperties{
		ResourceID: &storageAccountID,
		QueueName:  &queueName,
	}

	if v, ok := d.GetOk("storage_queue_endpoint.0.queue_message_time_to_live_in_seconds"); ok {
		queueMessageTimeToLiveInSeconds := int64(v.(int))
		storageQueueEventSubscriptionDestinationProperties.QueueMessageTimeToLiveInSeconds = &queueMessageTimeToLiveInSeconds
	}

	return eventgrid.StorageQueueEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeStorageQueue,
		StorageQueueEventSubscriptionDestinationProperties: storageQueueEventSubscriptionDestinationProperties,
	}
}

func expandEventGridEventSubscriptionEventhubEndpoint(d *pluginsdk.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	destinationProps := &eventgrid.EventHubEventSubscriptionDestinationProperties{}

	if _, ok := d.GetOk("eventhub_endpoint_id"); ok {
		endpoint := d.Get("eventhub_endpoint_id")
		destinationProps.ResourceID = utils.String(endpoint.(string))
	} else if _, ok := d.GetOk("eventhub_endpoint"); ok {
		ep := d.Get("eventhub_endpoint").([]interface{})
		if len(ep) == 0 || ep[0] == nil {
			return eventgrid.EventHubEventSubscriptionDestination{}
		}
		props := ep[0].(map[string]interface{})
		eventHubID := props["eventhub_id"].(string)
		destinationProps.ResourceID = &eventHubID
	}

	deliveryMappings := expandDeliveryProperties(d)
	destinationProps.DeliveryAttributeMappings = &deliveryMappings

	return eventgrid.EventHubEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeEventHub,
		EventHubEventSubscriptionDestinationProperties: destinationProps,
	}
}

func expandEventGridEventSubscriptionHybridConnectionEndpoint(d *pluginsdk.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	destinationProps := &eventgrid.HybridConnectionEventSubscriptionDestinationProperties{}

	if v, ok := d.GetOk("hybrid_connection_endpoint_id"); ok {
		destinationProps.ResourceID = utils.String(v.(string))
	} else if _, ok := d.GetOk("hybrid_connection_endpoint"); ok {
		ep := d.Get("hybrid_connection_endpoint").([]interface{})
		if len(ep) == 0 || ep[0] == nil {
			return eventgrid.HybridConnectionEventSubscriptionDestination{}
		}
		props := ep[0].(map[string]interface{})
		hybridConnectionID := props["hybrid_connection_id"].(string)
		destinationProps.ResourceID = &hybridConnectionID
	}

	deliveryMappings := expandDeliveryProperties(d)
	destinationProps.DeliveryAttributeMappings = &deliveryMappings

	return eventgrid.HybridConnectionEventSubscriptionDestination{
		EndpointType: eventgrid.EndpointTypeHybridConnection,
		HybridConnectionEventSubscriptionDestinationProperties: destinationProps,
	}
}

func expandEventGridEventSubscriptionAzureFunctionEndpoint(d *pluginsdk.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	input := d.Get("azure_function_endpoint")
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

	deliveryMappings := expandDeliveryProperties(d)
	props.DeliveryAttributeMappings = &deliveryMappings

	return azureFunctionDestination
}

func expandDeliveryProperties(d *pluginsdk.ResourceData) []eventgrid.BasicDeliveryAttributeMapping {
	var basicDeliveryAttributeMapping []eventgrid.BasicDeliveryAttributeMapping

	deliveryMappingsConfig, deliveryMappingsExists := d.GetOk("delivery_property")
	if !deliveryMappingsExists {
		return basicDeliveryAttributeMapping
	}

	input := deliveryMappingsConfig.([]interface{})
	if len(input) == 0 {
		return basicDeliveryAttributeMapping
	}

	for _, r := range input {
		mappingBlock := r.(map[string]interface{})

		if mappingBlock["type"].(string) == "Static" {
			basicDeliveryAttributeMapping = append(basicDeliveryAttributeMapping, eventgrid.StaticDeliveryAttributeMapping{
				Name: utils.String(mappingBlock["header_name"].(string)),
				Type: eventgrid.TypeStatic,
				StaticDeliveryAttributeMappingProperties: &eventgrid.StaticDeliveryAttributeMappingProperties{
					Value:    utils.String(mappingBlock["value"].(string)),
					IsSecret: utils.Bool(mappingBlock["secret"].(bool)),
				},
			})
		} else if mappingBlock["type"].(string) == "Dynamic" {
			basicDeliveryAttributeMapping = append(basicDeliveryAttributeMapping, eventgrid.DynamicDeliveryAttributeMapping{
				Name: utils.String(mappingBlock["header_name"].(string)),
				Type: eventgrid.TypeDynamic,
				DynamicDeliveryAttributeMappingProperties: &eventgrid.DynamicDeliveryAttributeMappingProperties{
					SourceField: utils.String(mappingBlock["source_field"].(string)),
				},
			})
		}
	}

	return basicDeliveryAttributeMapping
}

func expandEventGridEventSubscriptionWebhookEndpoint(d *pluginsdk.ResourceData) eventgrid.BasicEventSubscriptionDestination {
	input := d.Get("webhook_endpoint")
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

	deliveryMappings := expandDeliveryProperties(d)
	props.DeliveryAttributeMappings = &deliveryMappings

	return webhookDestination
}

func expandEventGridEventSubscriptionFilter(d *pluginsdk.ResourceData) (*eventgrid.EventSubscriptionFilter, error) {
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

	if v, ok := d.GetOk("advanced_filtering_on_arrays_enabled"); ok {
		filter.EnableAdvancedFilteringOnArrays = utils.Bool(v.(bool))
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
		return eventgrid.NumberNotInAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeNumberNotIn, Values: v}, nil
	case "string_begins_with":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringBeginsWithAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringBeginsWith, Values: v}, nil
	case "string_not_begins_with":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringNotBeginsWithAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringNotBeginsWith, Values: v}, nil
	case "string_ends_with":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringEndsWithAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringEndsWith, Values: v}, nil
	case "string_not_ends_with":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringNotEndsWithAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringNotEndsWith, Values: v}, nil
	case "string_contains":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringContainsAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringContains, Values: v}, nil
	case "string_not_contains":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringNotContainsAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringNotContains, Values: v}, nil
	case "string_in":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringInAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringIn, Values: v}, nil
	case "string_not_in":
		v := utils.ExpandStringSlice(config["values"].([]interface{}))
		return eventgrid.StringNotInAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeStringNotIn, Values: v}, nil
	case "is_not_null":
		return eventgrid.IsNotNullAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeIsNotNull}, nil
	case "is_null_or_undefined":
		return eventgrid.IsNullOrUndefinedAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeIsNullOrUndefined}, nil
	case "number_in_range":
		v := utils.ExpandFloatRangeSlice(config["values"].([]interface{}))
		return eventgrid.NumberInRangeAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeNumberInRange, Values: v}, nil
	case "number_not_in_range":
		v := utils.ExpandFloatRangeSlice(config["values"].([]interface{}))
		return eventgrid.NumberNotInRangeAdvancedFilter{Key: &k, OperatorType: eventgrid.OperatorTypeNumberNotInRange, Values: v}, nil
	default:
		return nil, fmt.Errorf("Invalid `advanced_filter` operator_type %q used", operatorType)
	}
}

func expandEventGridEventSubscriptionStorageBlobDeadLetterDestination(d *pluginsdk.ResourceData) eventgrid.BasicDeadLetterDestination {
	if v, ok := d.GetOk("storage_blob_dead_letter_destination"); ok {
		dest := v.([]interface{})[0].(map[string]interface{})
		resourceID := dest["storage_account_id"].(string)
		blobName := dest["storage_blob_container_name"].(string)
		return eventgrid.StorageBlobDeadLetterDestination{
			EndpointType: eventgrid.EndpointTypeBasicDeadLetterDestinationEndpointTypeStorageBlob,
			StorageBlobDeadLetterDestinationProperties: &eventgrid.StorageBlobDeadLetterDestinationProperties{
				ResourceID:        &resourceID,
				BlobContainerName: &blobName,
			},
		}
	}

	return nil
}

func expandEventGridEventSubscriptionRetryPolicy(d *pluginsdk.ResourceData) *eventgrid.RetryPolicy {
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

func expandEventGridEventSubscriptionIdentity(input []interface{}) (*eventgrid.EventSubscriptionIdentity, error) {
	if len(input) == 0 || input[0] == nil {
		return &eventgrid.EventSubscriptionIdentity{
			Type: eventgrid.EventSubscriptionIdentityType("None"),
		}, nil
	}

	identity := input[0].(map[string]interface{})
	identityType := eventgrid.EventSubscriptionIdentityType(identity["type"].(string))
	eventgridIdentity := eventgrid.EventSubscriptionIdentity{
		Type: identityType,
	}

	userAssignedIdentity := identity["user_assigned_identity"].(string)
	if identityType == eventgrid.EventSubscriptionIdentityTypeUserAssigned {
		eventgridIdentity.UserAssignedIdentity = utils.String(userAssignedIdentity)
	} else if len(userAssignedIdentity) > 0 {
		return nil, fmt.Errorf("`user_assigned_identity` can only be specified when `type` is `UserAssigned`; but `type` is currently %q", identityType)
	}

	return &eventgridIdentity, nil
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

func flattenDeliveryProperties(d *pluginsdk.ResourceData, input *[]eventgrid.BasicDeliveryAttributeMapping) []interface{} {
	if input == nil {
		return nil
	}

	deliveryProperties := make([]interface{}, len(*input))
	for i, element := range *input {
		attributeMapping := make(map[string]interface{})

		if staticMapping, ok := element.AsStaticDeliveryAttributeMapping(); ok {
			attributeMapping["type"] = staticMapping.Type
			if staticMapping.Name != nil {
				attributeMapping["header_name"] = staticMapping.Name
			}
			if staticMapping.IsSecret != nil {
				attributeMapping["secret"] = staticMapping.IsSecret
			}

			if *staticMapping.IsSecret {
				// If this is a secret, the Azure API just returns a value of 'Hidden',
				// so we need to lookup the value that was provided from config to return
				propertiesFromConfig := expandDeliveryProperties(d)
				for _, v := range propertiesFromConfig {
					if configMap, ok := v.AsStaticDeliveryAttributeMapping(); ok {
						if *configMap.Name == *staticMapping.Name {
							if configMap.Value != nil {
								attributeMapping["value"] = configMap.Value
							}
							break
						}
					}
				}
			} else {
				attributeMapping["value"] = staticMapping.Value
			}
		} else if dynamicMapping, ok := element.AsDynamicDeliveryAttributeMapping(); ok {
			attributeMapping["type"] = dynamicMapping.Type
			if dynamicMapping.Name != nil {
				attributeMapping["header_name"] = dynamicMapping.Name
			}
			if dynamicMapping.SourceField != nil {
				attributeMapping["source_field"] = dynamicMapping.SourceField
			}
		}

		deliveryProperties[i] = attributeMapping
	}

	return deliveryProperties
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

	if input.QueueMessageTimeToLiveInSeconds != nil {
		result["queue_message_time_to_live_in_seconds"] = *input.QueueMessageTimeToLiveInSeconds
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
	numberInRange := make([]interface{}, 0)
	numberNotInRange := make([]interface{}, 0)
	stringBeginsWith := make([]interface{}, 0)
	stringNotBeginsWith := make([]interface{}, 0)
	stringEndsWith := make([]interface{}, 0)
	stringNotEndsWith := make([]interface{}, 0)
	stringContains := make([]interface{}, 0)
	stringNotContains := make([]interface{}, 0)
	stringIn := make([]interface{}, 0)
	stringNotIn := make([]interface{}, 0)
	isNotNull := make([]interface{}, 0)
	isNullOrUndefined := make([]interface{}, 0)

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
		case eventgrid.StringNotBeginsWithAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringNotBeginsWith = append(stringNotBeginsWith, flattenValues(f.Key, &v))
		case eventgrid.StringEndsWithAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringEndsWith = append(stringEndsWith, flattenValues(f.Key, &v))
		case eventgrid.StringNotEndsWithAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringNotEndsWith = append(stringNotEndsWith, flattenValues(f.Key, &v))
		case eventgrid.StringContainsAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringContains = append(stringContains, flattenValues(f.Key, &v))
		case eventgrid.StringNotContainsAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringNotContains = append(stringNotContains, flattenValues(f.Key, &v))
		case eventgrid.StringInAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringIn = append(stringIn, flattenValues(f.Key, &v))
		case eventgrid.StringNotInAdvancedFilter:
			v := utils.FlattenStringSlice(f.Values)
			stringNotIn = append(stringNotIn, flattenValues(f.Key, &v))
		case eventgrid.NumberInRangeAdvancedFilter:
			v := utils.FlattenFloatRangeSlice(f.Values)
			numberInRange = append(numberInRange, flattenRangeValues(f.Key, &v))
		case eventgrid.NumberNotInRangeAdvancedFilter:
			v := utils.FlattenFloatRangeSlice(f.Values)
			numberNotInRange = append(numberNotInRange, flattenRangeValues(f.Key, &v))
		case eventgrid.IsNotNullAdvancedFilter:
			isNotNull = append(isNotNull, flattenKey(f.Key))
		case eventgrid.IsNullOrUndefinedAdvancedFilter:
			isNullOrUndefined = append(isNullOrUndefined, flattenKey(f.Key))
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
			"number_in_range":               numberInRange,
			"number_not_in_range":           numberNotInRange,
			"string_begins_with":            stringBeginsWith,
			"string_not_begins_with":        stringNotBeginsWith,
			"string_ends_with":              stringEndsWith,
			"string_not_ends_with":          stringNotEndsWith,
			"string_contains":               stringContains,
			"string_not_contains":           stringNotContains,
			"string_in":                     stringIn,
			"string_not_in":                 stringNotIn,
			"is_not_null":                   isNotNull,
			"is_null_or_undefined":          isNullOrUndefined,
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

func flattenRangeValues(inputKey *string, inputValues *[][]interface{}) map[string]interface{} {
	key := ""
	if inputKey != nil {
		key = *inputKey
	}
	values := make([]interface{}, 0)
	if inputValues != nil {
		for _, item := range *inputValues {
			values = append(values, item)
		}
	}

	return map[string]interface{}{
		"key":    key,
		"values": values,
	}
}

func flattenKey(inputKey *string) map[string]interface{} {
	key := ""
	if inputKey != nil {
		key = *inputKey
	}

	return map[string]interface{}{
		"key": key,
	}
}

func flattenEventGridEventSubscriptionIdentity(input *eventgrid.EventSubscriptionIdentity) []interface{} {
	if input == nil || string(input.Type) == "None" {
		return []interface{}{}
	}

	result := map[string]interface{}{
		"type": string(input.Type),
	}

	if input.UserAssignedIdentity != nil {
		result["user_assigned_identity"] = *input.UserAssignedIdentity
	}

	return []interface{}{
		result,
	}
}
