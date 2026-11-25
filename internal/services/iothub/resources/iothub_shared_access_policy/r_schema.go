// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_shared_access_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func iothubSharedAccessPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the IotHub Shared Access Policy resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubSharedAccessPolicyName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"iothub_name": {
			Description:  "The name of the IoTHub to which this Shared Access Policy belongs.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IoTHubName,
		},

		"registry_read": {
			Description: "Adds `RegistryRead` permission to this Shared Access Account. It allows read access to the identity registry.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"registry_write": {
			Description: "Adds `RegistryWrite` permission to this Shared Access Account. It allows write access to the identity registry.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"service_connect": {
			Description: "Adds `ServiceConnect` permission to this Shared Access Account. It allows sending and receiving on the cloud-side endpoints.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"device_connect": {
			Description: "Adds `DeviceConnect` permission to this Shared Access Account. It allows sending and receiving on the device-side endpoints.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"primary_key": {
			Description: "The primary key used to create the authentication token.",
			Type:        pluginsdk.TypeString,
			Sensitive:   true,
			Computed:    true,
		},

		"primary_connection_string": {
			Description: "The primary connection string of the Shared Access Policy.",
			Type:        pluginsdk.TypeString,
			Sensitive:   true,
			Computed:    true,
		},

		"secondary_key": {
			Description: "The secondary key used to create the authentication token.",
			Type:        pluginsdk.TypeString,
			Sensitive:   true,
			Computed:    true,
		},

		"secondary_connection_string": {
			Description: "The secondary connection string of the Shared Access Policy.",
			Type:        pluginsdk.TypeString,
			Sensitive:   true,
			Computed:    true,
		},
	}
}
