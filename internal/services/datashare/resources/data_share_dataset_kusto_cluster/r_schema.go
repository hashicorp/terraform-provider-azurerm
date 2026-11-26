// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share_dataset_kusto_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func dataShareDatasetKustoClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Data Share Kusto Cluster Dataset.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataSetName(),
		},

		"share_id": {
			Description:  "The resource ID of the Data Share where this Data Share Kusto Cluster Dataset should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: share.ValidateShareID,
		},

		"kusto_cluster_id": {
			Description:  "The resource ID of the Kusto Cluster to be shared with the receiver.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKustoClusterID,
		},

		"display_name": {
			Description: "The name of the Data Share Dataset.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"kusto_cluster_location": {
			Description: "The location of the Kusto Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
