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

		"parameter_values": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
