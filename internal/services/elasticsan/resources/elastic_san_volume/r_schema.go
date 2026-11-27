// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	diskSnapshots "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/restorepoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Elastic SAN Volume.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticSanVolumeName,
		},

		"volume_group_id": commonschema.ResourceIDReferenceRequiredForceNew(&volumes.VolumeGroupId{}, "Volume Group"),

		"size_in_gib": {
			Description:  "The size of the Elastic SAN Volume in GiB. The size should be within the remaining capacity of the parent Elastic SAN.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65536),
		},

		"create_source": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source_id": {
						Description: "The ID of the source to create the Elastic SAN Volume from.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
						ValidateFunc: validation.Any(
							commonids.ValidateManagedDiskID,
							restorepoints.ValidateRestorePointID,
							diskSnapshots.ValidateSnapshotID,
							snapshots.ValidateSnapshotID,
						),
					},
					"source_type": {
						Description: "The type of the source to create the Elastic SAN Volume from.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
						ValidateFunc: validation.StringInSlice([]string{
							// None is not exposed
							string(volumes.VolumeCreateOptionDisk),
							string(volumes.VolumeCreateOptionDiskRestorePoint),
							string(volumes.VolumeCreateOptionDiskSnapshot),
							string(volumes.VolumeCreateOptionVolumeSnapshot),
						},
							false),
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"target_iqn": {
			Description: "The iSCSI Target IQN of the Elastic SAN Volume.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"target_portal_hostname": {
			Description: "The iSCSI Target Portal Host Name of the Elastic SAN Volume.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"target_portal_port": {
			Description: "The iSCSI Target Portal Port of the Elastic SAN Volume.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"volume_id": {
			Description: "The UUID of the Elastic SAN Volume.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldVal, newVal := metadata.ResourceDiff.GetChange("size_in_gib"); newVal.(int) < oldVal.(int) {
				return errors.New("new size_in_gib should be greater than the existing one")
			}

			return nil
		},
	}
}
