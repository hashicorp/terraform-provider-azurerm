// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataConnectionName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"cluster_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterName,
		},

		"database_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatabaseName,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"eventhub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: eventhubs.ValidateEventhubID,
		},

		"eventhub_consumer_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				eventhubValidate.ValidateEventHubConsumerName(),
				validation.StringInSlice([]string{"$Default"}, false),
			),
		},

		"blob_storage_event_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(dataconnections.BlobStorageEventTypeMicrosoftPointStoragePointBlobCreated),
			ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForBlobStorageEventType(), false),
		},

		"skip_first_record": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"table_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.EntityName,
		},

		"mapping_rule_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.EntityName,
		},

		"data_format": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForEventGridDataFormat(), false),
		},

		"database_routing_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(dataconnections.DatabaseRoutingSingle),
			ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForDatabaseRouting(), false),
		},

		"eventgrid_event_subscription_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: eventsubscriptions.ValidateScopedEventSubscriptionID,
		},

		"managed_identity_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.Any(
				commonids.ValidateKustoClusterID,
				commonids.ValidateUserAssignedIdentityID,
			),
		},
	}
}
