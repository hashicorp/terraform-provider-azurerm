// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_workspace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Azure Monitor Workspace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"public_network_access_enabled": {
			Description: "Whether public network access is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"query_endpoint": {
			Description: "The query endpoint for the Azure Monitor Workspace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"default_data_collection_endpoint_id": {
			Description: "The ID of the managed default Data Collection Endpoint created with the Azure Monitor Workspace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"default_data_collection_rule_id": {
			Description: "The ID of the managed default Data Collection Rule created with the Azure Monitor Workspace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
