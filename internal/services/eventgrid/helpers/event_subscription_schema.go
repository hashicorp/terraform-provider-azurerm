// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/eventsubscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// EventSubscriptionEndpointType enumerates the values for event subscription endpoint types.
type EventSubscriptionEndpointType string

const (
	AzureFunctionEndpoint EventSubscriptionEndpointType = "azure_function_endpoint"
	EventHubID            EventSubscriptionEndpointType = "eventhub_id"
	HybridConnectionID    EventSubscriptionEndpointType = "hybrid_connection_id"
	ServiceBusQueueID     EventSubscriptionEndpointType = "service_bus_queue_id"
	ServiceBusTopicID     EventSubscriptionEndpointType = "service_bus_topic_id"
	StorageQueueEndpoint  EventSubscriptionEndpointType = "storage_queue_endpoint"
	WebHookEndpoint       EventSubscriptionEndpointType = "webhook_endpoint"
)

func EventSubscriptionSchemaEventSubscriptionName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description: "The name of the Event Subscription.",
		Type:        pluginsdk.TypeString,
		Required:    true,
		ForceNew:    true,
		ValidateFunc: validation.All(
			validation.StringIsNotEmpty,
			validation.StringMatch(
				regexp.MustCompile("^[-a-zA-Z0-9]{3,64}$"),
				"EventGrid subscription name must be 3 - 64 characters long, contain only letters, numbers and hyphens.",
			),
		),
	}
}

func EventSubscriptionSchemaEventDeliverySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description:  "The event delivery schema for the event subscription.",
		Type:         pluginsdk.TypeString,
		Optional:     true,
		ForceNew:     true,
		Default:      string(eventsubscriptions.EventDeliverySchemaEventGridSchema),
		ValidateFunc: validation.StringInSlice(eventsubscriptions.PossibleValuesForEventDeliverySchema(), false),
	}
}

func EventSubscriptionSchemaDeliveryProperty() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"header_name": {
					Description:      "The name of the header to forward to the destination.",
					Type:             pluginsdk.TypeString,
					Required:         true,
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"type": {
					Description: "The type of header value.",
					Type:        pluginsdk.TypeString,
					Required:    true,
					ValidateFunc: validation.StringInSlice([]string{
						"Static",
						"Dynamic",
					}, false),
				},

				"value": {
					Description: "The header value. Must be provided when the `type` is `Static`.",
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Sensitive:   true,
				},

				"source_field": {
					Description: "If the `type` is `Dynamic`, then provide the payload field to be used as the value. Valid source fields differ by subscription type.",
					Type:        pluginsdk.TypeString,
					Optional:    true,
				},

				"secret": {
					Description: "Whether the value is secret and should be protected. When `true`, this value won't be returned via Azure API requests.",
					Type:        pluginsdk.TypeBool,
					Optional:    true,
				},
			},
		},
	}
}

func EventSubscriptionSchemaExpirationTimeUTC() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description:  "The expiration time of the event subscription (Datetime Format `RFC 3339`).",
		Type:         pluginsdk.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
}

func EventSubscriptionSchemaAzureFunctionEndpoint(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: conflictsWith,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"function_id": {
					Description:  "The ID of the Function where the Event Subscription will receive events. This must be the functions ID in format {function_app.id}/functions/{name}.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: webapps.ValidateFunctionID,
				},
				"max_events_per_batch": {
					Description: "Maximum number of events per batch.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
				},
				"preferred_batch_size_in_kilobytes": {
					Description: "Preferred batch size in Kilobytes.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
				},
			},
		},
	}
}

func EventSubscriptionSchemaEventHubEndpointID(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description:   "The ID where the Event Hub is located.",
		Type:          pluginsdk.TypeString,
		Optional:      true,
		Computed:      true, // azignore:AZS007 - pre-existing violation
		ValidateFunc:  eventhubs.ValidateEventhubID,
		ConflictsWith: conflictsWith,
	}
}

func EventSubscriptionSchemaHybridConnectionEndpointID(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description:   "The ID where the Hybrid Connection is located.",
		Type:          pluginsdk.TypeString,
		Optional:      true,
		Computed:      true, // azignore:AZS007 - pre-existing violation
		ValidateFunc:  hybridconnections.ValidateHybridConnectionID,
		ConflictsWith: conflictsWith,
	}
}

func EventSubscriptionSchemaServiceBusQueueEndpointID(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description:   "The ID where the Service Bus Queue is located.",
		Type:          pluginsdk.TypeString,
		Optional:      true,
		ValidateFunc:  queues.ValidateQueueID,
		ConflictsWith: conflictsWith,
	}
}

func EventSubscriptionSchemaServiceBusTopicEndpointID(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description:   "The ID where the Service Bus Topic is located.",
		Type:          pluginsdk.TypeString,
		Optional:      true,
		ValidateFunc:  topics.ValidateTopicID,
		ConflictsWith: conflictsWith,
	}
}

func EventSubscriptionSchemaStorageQueueEndpoint(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: conflictsWith,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"storage_account_id": {
					Description:  "The ID of the Storage Account ID where the Storage Queue is located.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: commonids.ValidateStorageAccountID,
				},
				"queue_name": {
					Description:  "The name of the Storage Queue where the Event Subscription will receive events.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"queue_message_time_to_live_in_seconds": {
					Description: "Storage Queue message time to live in seconds.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
				},
			},
		},
	}
}

func EventSubscriptionSchemaWebHookEndpoint(conflictsWith []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: conflictsWith,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"url": {
					Description:  "The URL of the webhook where the Event Subscription will receive events.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsURLWithHTTPS,
				},
				"base_url": {
					Description: "The base URL of the webhook where the Event Subscription will receive events.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},
				"max_events_per_batch": {
					Description:  "Maximum number of events per batch.",
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 5000),
				},
				"preferred_batch_size_in_kilobytes": {
					Description:  "Preferred batch size in Kilobytes.",
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 1024),
				},
				"active_directory_tenant_id": {
					Description: "The Azure Active Directory Tenant ID to get the access token that will be included as the bearer token in delivery requests.",
					Type:        pluginsdk.TypeString,
					Optional:    true,
				},
				"active_directory_app_id_or_uri": {
					Description: "The Azure Active Directory Application ID or URI to get the access token that will be included as the bearer token in delivery requests.",
					Type:        pluginsdk.TypeString,
					Optional:    true,
				},
			},
		},
	}
}

func EventSubscriptionSchemaIncludedEventTypes() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description: "A list of applicable event types that need to be part of the event subscription.",
		Type:        pluginsdk.TypeList,
		Optional:    true,
		Computed:    true, // azignore:AZS007 - pre-existing violation
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func EventSubscriptionSchemaEnableAdvancedFilteringOnArrays() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description: "Whether advanced filters should be evaluated against an array of values instead of expecting a singular value.",
		Type:        pluginsdk.TypeBool,
		Optional:    true,
		Default:     false,
	}
}

func EventSubscriptionSchemaSubjectFilter() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"subject_begins_with": {
					Description:  "A string to filter events for an event subscription based on a resource path prefix.",
					Type:         pluginsdk.TypeString,
					Optional:     true,
					AtLeastOneOf: []string{"subject_filter.0.subject_begins_with", "subject_filter.0.subject_ends_with", "subject_filter.0.case_sensitive"},
				},
				"subject_ends_with": {
					Description:  "A string to filter events for an event subscription based on a resource path suffix.",
					Type:         pluginsdk.TypeString,
					Optional:     true,
					AtLeastOneOf: []string{"subject_filter.0.subject_begins_with", "subject_filter.0.subject_ends_with", "subject_filter.0.case_sensitive"},
				},
				"case_sensitive": {
					Description:  "Whether `subject_begins_with` and `subject_ends_with` are case sensitive.",
					Type:         pluginsdk.TypeBool,
					Optional:     true,
					AtLeastOneOf: []string{"subject_filter.0.subject_begins_with", "subject_filter.0.subject_ends_with", "subject_filter.0.case_sensitive"},
				},
			},
		},
	}
}

func EventSubscriptionSchemaAdvancedFilter() *pluginsdk.Schema {
	atLeastOneOf := []string{
		"advanced_filter.0.bool_equals", "advanced_filter.0.number_greater_than", "advanced_filter.0.number_greater_than_or_equals", "advanced_filter.0.number_less_than",
		"advanced_filter.0.number_less_than_or_equals", "advanced_filter.0.number_in", "advanced_filter.0.number_not_in", "advanced_filter.0.string_begins_with", "advanced_filter.0.string_not_begins_with",
		"advanced_filter.0.string_ends_with", "advanced_filter.0.string_not_ends_with", "advanced_filter.0.string_contains", "advanced_filter.0.string_not_contains", "advanced_filter.0.string_in",
		"advanced_filter.0.string_not_in", "advanced_filter.0.is_not_null", "advanced_filter.0.is_null_or_undefined", "advanced_filter.0.number_in_range", "advanced_filter.0.number_not_in_range",
	}
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"bool_equals": {
					Description:  "Compares a value of an event using a single boolean value.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"number_greater_than": {
					Description:  "Compares a value of an event using a single floating point number.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"number_greater_than_or_equals": {
					Description:  "Compares a value of an event using a single floating point number.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"number_less_than": {
					Description:  "Compares a value of an event using a single floating point number.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"number_less_than_or_equals": {
					Description:  "Compares a value of an event using a single floating point number.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"number_in": {
					Description:  "Compares a value of an event using multiple floating point numbers.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"number_not_in": {
					Description:  "Compares a value of an event using multiple floating point numbers.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"string_begins_with": {
					Description:  "Compares a value of an event using multiple string values.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"string_not_begins_with": {
					Description:  "Compares a value of an event using multiple string values.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"string_ends_with": {
					Description:  "Compares a value of an event using multiple string values.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"string_not_ends_with": {
					Description:  "Compares a value of an event using multiple string values.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"string_contains": {
					Description:  "Compares a value of an event using multiple string values.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"string_not_contains": {
					Description:  "Compares a value of an event using multiple string values.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"string_in": {
					Description:  "Compares a value of an event using multiple string values.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"string_not_in": {
					Description:  "Compares a value of an event using multiple string values.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"is_not_null": {
					Description:  "Evaluates if a value of an event isn't NULL or undefined.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"is_null_or_undefined": {
					Description:  "Evaluates if a value of an event is NULL or undefined.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"number_in_range": {
					Description:  "Compares a value of an event using multiple floating point number ranges.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
				"number_not_in_range": {
					Description:  "Compares a value of an event using multiple floating point number ranges.",
					Type:         pluginsdk.TypeList,
					Optional:     true,
					AtLeastOneOf: atLeastOneOf,
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
				},
			},
		},
	}
}

func EventSubscriptionSchemaStorageBlobDeadletterDestination() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"storage_account_id": {
					Description:  "The ID of the Storage Account ID where the storage blob is located.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: commonids.ValidateStorageAccountID,
				},
				"storage_blob_container_name": {
					Description:  "The name of the Storage blob container that is the destination of the deadletter events.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func EventSubscriptionSchemaRetryPolicy() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true, // azignore:AZS007 - pre-existing violation
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_delivery_attempts": {
					Description:  "The maximum number of delivery retry attempts for events.",
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 30),
				},
				"event_time_to_live": {
					Description:  "The time to live (in minutes) for events. See [official documentation](https://learn.microsoft.com/azure/event-grid/manage-event-delivery#set-retry-policy) for more details.",
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 1440),
				},
			},
		},
	}
}

func EventSubscriptionSchemaLabels() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description: "A list of labels to assign to the event subscription.",
		Type:        pluginsdk.TypeList,
		Optional:    true,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

func EventSubscriptionSchemaIdentity() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Description:  "The type of Managed Service Identity to configure on this resource.",
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(eventsubscriptions.PossibleValuesForEventSubscriptionIdentityType(), false),
				},
				"user_assigned_identity": {
					Description: "The user identity associated with the resource.",
					Type:        schema.TypeString,
					Optional:    true,
				},
			},
		},
	}
}
