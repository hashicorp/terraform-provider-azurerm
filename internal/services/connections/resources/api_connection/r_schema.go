// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/managedapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiConnectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"managed_api_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedapis.ValidateManagedApiID,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because Azure sets a default when `display_name` is not defined but the value depends on which managed API is provided.
			// For example:
			//   - Managed API `servicebus` defaults to `Service Bus`
			//   - Managed API `sftpwithssh` defaults to `SFTP - SSH`
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because Azure sets a default `kind` (e.g. `V1`) when it is not specified.
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"parameter_value_set": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"parameter_value_type"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"values": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"parameter_value_type": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"parameter_value_set", "parameter_values"},
		},

		"parameter_values": {
			Type:          pluginsdk.TypeMap,
			Optional:      true,
			ConflictsWith: []string{"parameter_value_type"},
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
