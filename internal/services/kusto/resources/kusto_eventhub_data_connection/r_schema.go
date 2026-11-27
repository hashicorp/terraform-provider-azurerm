// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_eventhub_data_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func kustoEventhubDataConnectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Kusto EventHub Data Connection to create.",
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

		"compression": {
			Description:  "The compression type for the connection.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      dataconnections.CompressionNone,
			ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForCompression(), false),
		},

		"database_name": {
			Description:  "The name of the Kusto Database this data connection will be added to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatabaseName,
		},

		"eventhub_id": {
			Description:  "The resource ID of the EventHub this data connection will use for ingestion.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: eventhubs.ValidateEventhubID,
		},

		"event_system_properties": {
			Description: "A list of system properties for the Event Hub.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"consumer_group": {
			Description: "The EventHub consumer group this data connection will use for ingestion.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				eventhubValidate.ValidateEventHubConsumerName(),
				validation.StringInSlice([]string{"$Default"}, false),
			),
		},

		"table_name": {
			Description:  "The target table name used for the message ingestion. Table must exist before resource is created.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.EntityName,
		},

		"identity_id": {
			Description: "The resource ID of a Managed Identity (system or user assigned) to be used to authenticate with Event Hub.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.Any(
				commonids.ValidateKustoClusterID,
				commonids.ValidateUserAssignedIdentityID,
			),
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
			ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForEventHubDataFormat(), false),
		},

		"database_routing_type": {
			Description:  "Indication for database routing information from the data connection, by default only database routing information is allowed.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(dataconnections.DatabaseRoutingSingle),
			ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForDatabaseRouting(), false),
		},

		"retrieval_start_date": {
			Description:  "The date after which data should be retrieved from Event Hub. When defined, the data connection retrieves existing events created since the specified retrieval start date. It can only retrieve events retained by the Event Hub, based on its retention period.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IsRFC3339Time,
		},
	}
}
