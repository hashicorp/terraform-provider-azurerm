// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_cluster

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func vmwareClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"vmware_cloud_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PrivateCloudID,
		},

		"cluster_node_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(3, 16),
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"av20",
				"av36",
				"av36t",
				"av36p",
				"av36pt",
				"av48",
				"av48t",
				"av52",
				"av52t",
				"av64",
			}, false),
		},

		"cluster_number": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"hosts": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
