// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_eventgrid_data_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/eventsubscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func kustoEventgridDataConnectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Kusto Event Grid Data Connection to create.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataConnectionName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"cluster_name": {
			Description:  "The name of the Kusto Cluster this data connection will be added to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterName,
		},

		"database_name": {
			Description:  "The name of the Kusto Database this data connection will be added to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatabaseName,
		},

		"storage_account_id": {
			Description:  "The resource ID of the Storage Account this data connection will use for ingestion.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"eventhub_id": {
			Description:  "The resource ID of the Event Hub this data connection will use for ingestion.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: eventhubs.ValidateEventhubID,
		},

		"eventhub_consumer_group_name": {
			Description: "The Event Hub consumer group this data connection will use for ingestion.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				eventhubValidate.ValidateEventHubConsumerName(),
				validation.StringInSlice([]string{"$Default"}, false),
			),
		},

		"blob_storage_event_type": {
			Description:  "The blob storage event type that needs to be processed.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(dataconnections.BlobStorageEventTypeMicrosoftPointStoragePointBlobCreated),
			ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForBlobStorageEventType(), false),
		},

		"skip_first_record": {
			Description: "Whether the first record of every file is ignored.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"table_name": {
			Description:  "The target table name used for the message ingestion. Table must exist before resource is created.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.EntityName,
		},

		"mapping_rule_name": {
			Description:  "The mapping rule used for the message ingestion. Mapping rule must exist before resource is created.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.EntityName,
		},

		"data_format": {
			Description:  "The data format of the EventHub messages.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForEventGridDataFormat(), false),
		},

		"database_routing_type": {
			Description:  "Indication for database routing information from the data connection, by default only database routing information is allowed.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(dataconnections.DatabaseRoutingSingle),
			ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForDatabaseRouting(), false),
		},

		"eventgrid_event_subscription_id": {
			Description:  "The resource ID of the Event Grid that is subscribed to the Storage Account events.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: eventsubscriptions.ValidateScopedEventSubscriptionID,
		},

		"managed_identity_id": {
			Description: "Empty for non-managed identity based data connection. For System Assigned Identity, provide cluster resource ID. For User Assigned Identity (UAI) provide the UAI resource ID.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.Any(
				commonids.ValidateKustoClusterID,
				commonids.ValidateUserAssignedIdentityID,
			),
		},
	}
}
