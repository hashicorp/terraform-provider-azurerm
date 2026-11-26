// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name which should be used for this API Connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"managed_api_id": {
			Description:  "The ID of the Managed API which this API Connection is linked to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedapis.ValidateManagedApiID,
		},

		"display_name": {
			Description: "A display name for this API Connection.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because Azure sets a default when `display_name` is not defined but the value depends on which managed API is provided.
			// For example:
			//   - Managed API `servicebus` defaults to `Service Bus`
			//   - Managed API `sftpwithssh` defaults to `SFTP - SSH`
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"parameter_values": {
			Description: "A map of parameter values associated with this API Connection.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
