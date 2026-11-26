// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Snapshot resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SnapshotName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"create_option": {
			Description: "How the snapshot is to be created.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(snapshots.DiskCreateOptionCopy),
				string(snapshots.DiskCreateOptionCopyStart),
				string(snapshots.DiskCreateOptionImport),
			}, false),
		},

		"incremental_enabled": {
			Description: "Whether the Snapshot is incremental.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"source_uri": {
			Description: "The URI to a Managed or Unmanaged Disk.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},

		"network_access_policy": {
			Description:  "Policy for accessing the disk via network.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(snapshots.NetworkAccessPolicyAllowAll),
			ValidateFunc: validation.StringInSlice(snapshots.PossibleValuesForNetworkAccessPolicy(), false),
		},

		"disk_access_id": {
			Description:  "The ID of the Disk Access which should be used for this Snapshot. This is used in conjunction with setting `network_access_policy` to `AllowPrivate`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: diskaccesses.ValidateDiskAccessID,
			// TODO:
			// the snapshot API is broken and returns the Resource Group name in UPPERCASE
			// tracked by https://github.com/Azure/azure-rest-api-specs/issues/29187
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"public_network_access_enabled": {
			Description: "Policy for controlling export on the disk.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"source_resource_id": {
			Description: "A reference to an existing snapshot, when `create_option` is `Copy` or `CopyStart`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},

		"storage_account_id": {
			Description: "The ID of a Storage Account. Used with `source_uri` to allow authorization during import of unmanaged blobs from a different subscription.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},

		"disk_size_gb": {
			Description: "The size of the Snapshotted Disk in GB.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// Note: O+C because Azure computes disk size when not specified
			Computed: true,
		},

		"encryption_settings": helpers.EncryptionSettingsSchema(),

		"trusted_launch_enabled": {
			Description: "Whether Trusted Launch is enabled for the Snapshot.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
