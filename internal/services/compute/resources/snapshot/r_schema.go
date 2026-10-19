// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package snapshot

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskaccesses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func snapshotSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SnapshotName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"create_option": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(snapshots.DiskCreateOptionCopy),
				string(snapshots.DiskCreateOptionCopyStart),
				string(snapshots.DiskCreateOptionImport),
			}, false),
		},

		"incremental_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"source_uri": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"network_access_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(snapshots.PossibleValuesForNetworkAccessPolicy(), false),
			Default:      string(snapshots.NetworkAccessPolicyAllowAll),
		},

		"disk_access_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: diskaccesses.ValidateDiskAccessID,
			// TODO:
			// the snapshot API is broken and returns the Resource Group name in UPPERCASE
			// tracked by https://github.com/Azure/azure-rest-api-specs/issues/29187
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"source_resource_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"storage_account_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"disk_size_gb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// Note: O+C because Azure computes disk size when not specified
			Computed: true,
		},

		"encryption_settings": helpers.EncryptionSettingsSchema(),

		"trusted_launch_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
