// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataSetName(),
		},

		"share_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: share.ValidateShareID,
		},

		"kusto_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKustoClusterID,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kusto_cluster_location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
