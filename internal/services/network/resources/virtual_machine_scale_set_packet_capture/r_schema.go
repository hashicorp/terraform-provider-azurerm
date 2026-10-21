// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_packet_capture

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkwatchers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/packetcaptures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualMachineScaleSetPacketCaptureSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_watcher_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkwatchers.ValidateNetworkWatcherID,
		},

		"virtual_machine_scale_set_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				commonids.ValidateVirtualMachineScaleSetID,
			),
		},

		"maximum_bytes_per_packet": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  0,
		},

		"maximum_bytes_per_session": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  1073741824,
		},

		"maximum_capture_duration_in_seconds": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      18000,
			ValidateFunc: validation.IntBetween(1, 18000),
		},

		"storage_location": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"file_path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.FilePath,
						AtLeastOneOf: []string{"storage_location.0.file_path", "storage_location.0.storage_account_id"},
					},
					"storage_account_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
						AtLeastOneOf: []string{"storage_location.0.file_path", "storage_location.0.storage_account_id"},
					},
					"storage_path": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"local_ip_address": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"local_port": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"protocol": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(packetcaptures.PossibleValuesForPcProtocol(), false),
					},
					"remote_ip_address": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"remote_port": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
				},
			},
		},

		"machine_scope": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"exclude_instance_ids": {
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ForceNew:      true,
						ConflictsWith: []string{"machine_scope.0.include_instance_ids"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"include_instance_ids": {
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ForceNew:      true,
						ConflictsWith: []string{"machine_scope.0.exclude_instance_ids"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}
