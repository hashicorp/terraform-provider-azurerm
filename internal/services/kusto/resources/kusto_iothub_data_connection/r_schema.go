// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_iothub_data_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func kustoIothubDataConnectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Kusto IotHub Data Connection to create.",
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

		"iothub_id": {
			Description:  "The resource ID of the IotHub this data connection will use for ingestion.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iothubValidate.IotHubID,
		},

		"consumer_group": {
			Description:  "The IotHub consumer group this data connection will use for ingestion.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iothubValidate.IoTHubConsumerGroupName,
		},

		"shared_access_policy_name": {
			Description:  "The IotHub Shared Access Policy this data connection will use for ingestion, which must have read permission.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iothubValidate.IotHubSharedAccessPolicyName,
		},

		"table_name": {
			Description:  "The target table name used for the message ingestion. Table must exist before resource is created.",
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validate.EntityName,
		},

		"mapping_rule_name": {
			Description:  "The mapping rule used for the message ingestion. Mapping rule must exist before resource is created.",
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validate.EntityName,
		},

		"data_format": {
			Description:  "The data format of the IoTHub messages.",
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForIotHubDataFormat(), false),
		},

		"database_routing_type": {
			Description:  "Indication for database routing information from the data connection, by default only database routing information is allowed.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(dataconnections.DatabaseRoutingSingle),
			ValidateFunc: validation.StringInSlice(dataconnections.PossibleValuesForDatabaseRouting(), false),
		},

		"event_system_properties": {
			Description: "The System Properties that each IoT Hub message should contain.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"retrieval_start_date": {
			Description:  "The date after which data should be retrieved from IoT Hub. When defined, the data connection retrieves existing events created since the specified retrieval start date. It can only retrieve events retained by the IoT Hub, based on its retention period.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IsRFC3339Time,
		},
	}
}
