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
			Description:  "The name of the IotHub Shared Access Policy resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubSharedAccessPolicyName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"iothub_dps_name": {
			Description:  "The name of the IoT Hub Device Provisioning service to which this Shared Access Policy belongs.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IoTHubName,
		},

		"enrollment_read": {
			Description: "Adds `EnrollmentRead` permission to this Shared Access Account. It allows read access to enrollment data.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"enrollment_write": {
			Description: "Adds `EnrollmentWrite` permission to this Shared Access Account. It allows write access to enrollment data.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"registration_read": {
			Description: "Adds `RegistrationStatusRead` permission to this Shared Access Account. It allows read access to device registrations.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"registration_write": {
			Description: "Adds `RegistrationStatusWrite` permission to this Shared Access Account. It allows write access to device registrations.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"service_config": {
			Description: "Adds `ServiceConfig` permission to this Shared Access Account. It allows configuration of the Device Provisioning Service.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"primary_key": {
			Description: "The primary key used to create the authentication token.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_connection_string": {
			Description: "The primary connection string of the Shared Access Policy.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_key": {
			Description: "The secondary key used to create the authentication token.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_connection_string": {
			Description: "The secondary connection string of the Shared Access Policy.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}
