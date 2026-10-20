// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_system_topic_event_subscription

import (
	"slices"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func eventgridSystemTopicEventSubscriptionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": helpers.EventSubscriptionSchemaEventSubscriptionName(),

		"system_topic": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"event_delivery_schema": helpers.EventSubscriptionSchemaEventDeliverySchema(),

		"expiration_time_utc": helpers.EventSubscriptionSchemaExpirationTimeUTC(),

		"azure_function_endpoint": helpers.EventSubscriptionSchemaAzureFunctionEndpoint(
			slices.DeleteFunc(
				possibleSystemTopicEventSubscriptionEndpointTypes(),
				func(s string) bool { return s == string(helpers.AzureFunctionEndpoint) },
			),
		),

		"eventhub_id": helpers.EventSubscriptionSchemaEventHubEndpointID(
			slices.DeleteFunc(
				possibleSystemTopicEventSubscriptionEndpointTypes(),
				func(s string) bool { return s == string(helpers.EventHubID) },
			),
		),

		"hybrid_connection_id": helpers.EventSubscriptionSchemaHybridConnectionEndpointID(
			slices.DeleteFunc(
				possibleSystemTopicEventSubscriptionEndpointTypes(),
				func(s string) bool { return s == string(helpers.HybridConnectionID) },
			),
		),

		"service_bus_queue_id": helpers.EventSubscriptionSchemaServiceBusQueueEndpointID(
			slices.DeleteFunc(
				possibleSystemTopicEventSubscriptionEndpointTypes(),
				func(s string) bool { return s == string(helpers.ServiceBusQueueID) },
			),
		),

		"service_bus_topic_id": helpers.EventSubscriptionSchemaServiceBusTopicEndpointID(
			slices.DeleteFunc(
				possibleSystemTopicEventSubscriptionEndpointTypes(),
				func(s string) bool { return s == string(helpers.ServiceBusTopicID) },
			),
		),

		"storage_queue_endpoint": helpers.EventSubscriptionSchemaStorageQueueEndpoint(
			slices.DeleteFunc(
				possibleSystemTopicEventSubscriptionEndpointTypes(),
				func(s string) bool { return s == string(helpers.StorageQueueEndpoint) },
			),
		),

		"webhook_endpoint": helpers.EventSubscriptionSchemaWebHookEndpoint(
			slices.DeleteFunc(
				possibleSystemTopicEventSubscriptionEndpointTypes(),
				func(s string) bool { return s == string(helpers.WebHookEndpoint) },
			),
		),

		"included_event_types": helpers.EventSubscriptionSchemaIncludedEventTypes(),

		"subject_filter": helpers.EventSubscriptionSchemaSubjectFilter(),

		"advanced_filter": helpers.EventSubscriptionSchemaAdvancedFilter(),

		"delivery_identity": helpers.EventSubscriptionSchemaIdentity(),

		"dead_letter_identity": helpers.EventSubscriptionSchemaIdentity(),

		"storage_blob_dead_letter_destination": helpers.EventSubscriptionSchemaStorageBlobDeadletterDestination(),

		"retry_policy": helpers.EventSubscriptionSchemaRetryPolicy(),

		"labels": helpers.EventSubscriptionSchemaLabels(),

		"advanced_filtering_on_arrays_enabled": helpers.EventSubscriptionSchemaEnableAdvancedFilteringOnArrays(),

		"delivery_property": helpers.EventSubscriptionSchemaDeliveryProperty(),
	}
}
