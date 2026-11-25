// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_time_series_database_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/timeseriesdatabaseconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/digitaltwins/validate"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	kustoValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (m Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Digital Twins Time Series Database Connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DigitalTwinsTimeSeriesDatabaseConnectionName,
		},

		"digital_twins_id": {
			Description:  "The ID of the Digital Twins.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: timeseriesdatabaseconnections.ValidateDigitalTwinsInstanceID,
		},

		"eventhub_name": {
			Description:  "Name of the Event Hub.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: eventhubValidate.ValidateEventHubName(),
		},

		"eventhub_namespace_id": {
			Description:  "The ID of the Event Hub Namespace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: eventhubs.ValidateNamespaceID,
		},

		"eventhub_namespace_endpoint_uri": {
			Description:  "URI of the Event Hub Namespace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"kusto_cluster_id": {
			Description:  "The ID of the Kusto Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKustoClusterID,
		},

		"kusto_cluster_uri": {
			Description:  "URI of the Kusto Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"kusto_database_name": {
			Description:  "Name of the Kusto Database.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: kustoValidate.DatabaseName,
		},

		"eventhub_consumer_group_name": {
			Description:  "Name of the Event Hub Consumer Group.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "$Default",
			ForceNew:     true,
			ValidateFunc: eventhubValidate.ValidateEventHubConsumerName(),
		},

		"kusto_table_name": {
			Description:  "Name of the Kusto Table.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "AdtPropertyEvents",
			ForceNew:     true,
			ValidateFunc: kustoValidate.EntityName,
		},
	}
}

func (m Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
