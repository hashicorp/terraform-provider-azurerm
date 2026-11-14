// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_device_update_instance

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubDeviceUpdateInstanceName,
		},

		"device_update_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: deviceupdates.ValidateAccountID,
		},

		"iothub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubID,
		},

		"diagnostic_storage_account": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},

					"connection_string": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"diagnostic_storage_account.0.connection_string_wo"},
					},

					"connection_string_wo": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						WriteOnly:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"diagnostic_storage_account.0.connection_string"},
						RequiredWith: []string{"diagnostic_storage_account.0.connection_string_wo_version"},
					},

					"connection_string_wo_version": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
						RequiredWith: []string{"diagnostic_storage_account.0.connection_string_wo"},
					},
				},
			},
		},

		"diagnostic_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
