// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps_shared_access_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func iothubDpsSharedAccessPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubSharedAccessPolicyName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"iothub_dps_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IoTHubName,
		},

		"enrollment_read": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"enrollment_write": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"registration_read": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"registration_write": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"service_config": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"primary_key": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"secondary_key": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},
	}
}
