// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume

import (
	"context"
	"fmt"
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticSanVolumeName,
		},

		"volume_group_id": commonschema.ResourceIDReferenceRequiredForceNew(&volumes.VolumeGroupId{}),

		"size_in_gib": {
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
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.Any(
							commonids.ValidateManagedDiskID,
							restorepoints.ValidateRestorePointID,
							diskSnapshots.ValidateSnapshotID,
							snapshots.ValidateSnapshotID,
						),
					},
					"source_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
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
			Computed: true,
			Type:     pluginsdk.TypeString,
		},

		"target_portal_hostname": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},

		"target_portal_port": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"volume_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldVal, newVal := metadata.ResourceDiff.GetChange("size_in_gib"); newVal.(int) < oldVal.(int) {
				return fmt.Errorf("new size_in_gib should be greater than the existing one")
			}

			return nil
		},
	}
}
