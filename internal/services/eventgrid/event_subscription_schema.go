package eventgrid

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/eventsubscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	serviceBusQueues "github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queues"
	serviceBusTopics "github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// EventSubscriptionEndpointType enumerates the values for event subscription endpoint types.
type EventSubscriptionEndpointType string

const (
	// AzureFunctionEndpoint ...
	AzureFunctionEndpoint EventSubscriptionEndpointType = "azure_function_endpoint"
	// EventHubEndpointID ...
	EventHubEndpointID EventSubscriptionEndpointType = "eventhub_endpoint_id"
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
		Type:         pluginsdk.TypeString,
		Optional:     true,
		ForceNew:     true,
		Default:      string(eventsubscriptions.EventDeliverySchemaEventGridSchema),
		ValidateFunc: validation.StringInSlice(eventsubscriptions.PossibleValuesForEventDeliverySchema(), false),
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
					ValidateFunc: azure.ValidateResourceID, // TODO: validation for a Function App ID
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
		ValidateFunc:  eventhubs.ValidateEventhubID,
	}
}

func eventSubscriptionSchemaHybridConnectionEndpointID(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeString,
		Optional:      true,
		Computed:      true,
		ConflictsWith: conflictsWith,
		ValidateFunc:  hybridconnections.ValidateHybridConnectionID,
	}
}

func eventSubscriptionSchemaServiceBusQueueEndpointID(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeString,
		Optional:      true,
		ConflictsWith: conflictsWith,
		ValidateFunc:  serviceBusQueues.ValidateQueueID,
	}
}

func eventSubscriptionSchemaServiceBusTopicEndpointID(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeString,
		Optional:      true,
		ConflictsWith: conflictsWith,
		ValidateFunc:  serviceBusTopics.ValidateTopicID,
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
					ValidateFunc: commonids.ValidateStorageAccountID,
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
					ValidateFunc: commonids.ValidateStorageAccountID,
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

func eventSubscriptionSchemaIdentity() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(eventsubscriptions.PossibleValuesForEventSubscriptionIdentityType(), false),
				},
				"user_assigned_identity": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}
