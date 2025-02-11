// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepointcollections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualMachineRestorePointResource struct{}

var _ sdk.Resource = VirtualMachineRestorePointResource{}

func (r VirtualMachineRestorePointResource) ModelObject() interface{} {
	return &VirtualMachineRestorePointResourceModel{}
}

type VirtualMachineRestorePointResourceModel struct {
	Name                                   string   `tfschema:"name"`
	VirtualMachineRestorePointCollectionId string   `tfschema:"virtual_machine_restore_point_collection_id"`
	CrashConsistencyModeEnabled            bool     `tfschema:"crash_consistency_mode_enabled"`
	ExcludedDisks                          []string `tfschema:"excluded_disks"`
}

func (r VirtualMachineRestorePointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return restorepoints.ValidateRestorePointID
}

func (r VirtualMachineRestorePointResource) ResourceType() string {
	return "azurerm_virtual_machine_restore_point"
}

func (r VirtualMachineRestorePointResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"virtual_machine_restore_point_collection_id": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: restorepointcollections.ValidateRestorePointCollectionID,
		},

		"crash_consistency_mode_enabled": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeBool,
			Default:  false,
		},

		"excluded_disks": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeSet,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateManagedDiskID,
			},
		},
	}
}

func (r VirtualMachineRestorePointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r VirtualMachineRestorePointResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.RestorePointsClient

			var config VirtualMachineRestorePointResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			collectionId, err := restorepointcollections.ParseRestorePointCollectionID(config.VirtualMachineRestorePointCollectionId)
			if err != nil {
				return err
			}

			id := restorepoints.NewRestorePointID(collectionId.SubscriptionId, collectionId.ResourceGroupName, collectionId.RestorePointCollectionName, config.Name)

			existing, err := client.Get(ctx, id, restorepoints.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := restorepoints.RestorePoint{
				Properties: &restorepoints.RestorePointProperties{},
			}

			if config.CrashConsistencyModeEnabled {
				parameters.Properties.ConsistencyMode = pointer.To(restorepoints.ConsistencyModeTypesCrashConsistent)
			}

			if len(config.ExcludedDisks) > 0 {
				excludedDisks := make([]restorepoints.ApiEntityReference, 0)
				for _, diskId := range config.ExcludedDisks {
					excludedDisks = append(excludedDisks, restorepoints.ApiEntityReference{
						Id: pointer.To(diskId),
					})
				}

				parameters.Properties.ExcludeDisks = pointer.To(excludedDisks)
			}

			if err = client.CreateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r VirtualMachineRestorePointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.RestorePointsClient

			schema := VirtualMachineRestorePointResourceModel{}

			id, err := restorepoints.ParseRestorePointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, restorepoints.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.RestorePointName
				schema.VirtualMachineRestorePointCollectionId = restorepointcollections.NewRestorePointCollectionID(id.SubscriptionId, id.ResourceGroupName, id.RestorePointCollectionName).ID()

				if props := model.Properties; props != nil {
					schema.CrashConsistencyModeEnabled = strings.EqualFold(string(pointer.From(props.ConsistencyMode)), string(restorepoints.ConsistencyModeTypesCrashConsistent))

					excludedDisksConfig := make([]string, 0)
					if excludedDisks := props.ExcludeDisks; excludedDisks != nil {
						for _, excludedDisk := range *excludedDisks {
							excludedDisksConfig = append(excludedDisksConfig, pointer.From(excludedDisk.Id))
						}
					}
					schema.ExcludedDisks = excludedDisksConfig
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r VirtualMachineRestorePointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.RestorePointsClient

			id, err := restorepoints.ParseRestorePointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
