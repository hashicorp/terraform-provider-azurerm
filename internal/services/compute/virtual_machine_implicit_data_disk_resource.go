package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource                   = VirtualMachineImplicitDataDiskResource{}
	_ sdk.ResourceWithUpdate         = VirtualMachineImplicitDataDiskResource{}
	_ sdk.ResourceWithCustomImporter = VirtualMachineImplicitDataDiskResource{}
)

type VirtualMachineImplicitDataDiskResource struct{}

func (r VirtualMachineImplicitDataDiskResource) ModelObject() interface{} {
	return &VirtualMachineImplicitDataDiskResourceModel{}
}

type VirtualMachineImplicitDataDiskResourceModel struct {
	Name                    string `tfschema:"name"`
	VirtualMachineId        string `tfschema:"virtual_machine_id"`
	Lun                     int64  `tfschema:"lun"`
	Caching                 string `tfschema:"caching"`
	CreateOption            string `tfschema:"create_option"`
	DeleteOption            string `tfschema:"delete_option"`
	DiskSizeGb              int64  `tfschema:"disk_size_gb"`
	SourceResourceId        string `tfschema:"source_resource_id"`
	WriteAcceleratorEnabled bool   `tfschema:"write_accelerator_enabled"`
}

func (r VirtualMachineImplicitDataDiskResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataDiskID
}

func (r VirtualMachineImplicitDataDiskResource) ResourceType() string {
	return "azurerm_virtual_machine_implicit_data_disk"
}

func (r VirtualMachineImplicitDataDiskResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"virtual_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"lun": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"caching": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachines.CachingTypesNone),
				string(virtualmachines.CachingTypesReadOnly),
				string(virtualmachines.CachingTypesReadWrite),
			}, false),
		},

		"create_option": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachines.DiskCreateOptionTypesCopy),
			}, false),
		},

		"delete_option": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachines.DiskDeleteOptionTypesDetach),
				string(virtualmachines.DiskDeleteOptionTypesDelete),
			}, false),
		},

		"disk_size_gb": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 1023),
		},

		"source_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.Any(snapshots.ValidateSnapshotID, commonids.ValidateManagedDiskID),
		},

		"write_accelerator_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (r VirtualMachineImplicitDataDiskResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r VirtualMachineImplicitDataDiskResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient

			var config VirtualMachineImplicitDataDiskResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			virtualMachineId, err := virtualmachines.ParseVirtualMachineID(config.VirtualMachineId)
			if err != nil {
				return err
			}

			locks.ByName(virtualMachineId.VirtualMachineName, VirtualMachineResourceName)
			defer locks.UnlockByName(virtualMachineId.VirtualMachineName, VirtualMachineResourceName)

			id := parse.NewDataDiskID(subscriptionId, virtualMachineId.ResourceGroupName, virtualMachineId.VirtualMachineName, config.Name)

			resp, err := client.Get(ctx, *virtualMachineId, virtualmachines.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", *virtualMachineId, err)
			}

			expandedDisk := virtualmachines.DataDisk{
				Name:         pointer.To(config.Name),
				Caching:      pointer.To(virtualmachines.CachingTypes(config.Caching)),
				CreateOption: virtualmachines.DiskCreateOptionTypes(config.CreateOption),
				DeleteOption: pointer.To(virtualmachines.DiskDeleteOptionTypes(config.DeleteOption)),
				DiskSizeGB:   pointer.To(config.DiskSizeGb),
				Lun:          config.Lun,
				SourceResource: &virtualmachines.ApiEntityReference{
					Id: pointer.To(config.SourceResourceId),
				},
				WriteAcceleratorEnabled: pointer.To(config.WriteAcceleratorEnabled),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if profile := props.StorageProfile; profile != nil {
						if dataDisks := profile.DataDisks; dataDisks != nil {
							existingIndex := -1
							disks := *dataDisks
							for i, disk := range disks {
								if pointer.From(disk.Name) == config.Name {
									existingIndex = i
									break
								}
							}

							if existingIndex != -1 {
								return metadata.ResourceRequiresImport(r.ResourceType(), id)
							}

							disks = append(disks, expandedDisk)
							profile.DataDisks = &disks
							// fixes #24145
							props.ApplicationProfile = nil
							// fixes #2485
							model.Identity = nil
							// fixes #1600
							model.Resources = nil
							err = client.CreateOrUpdateThenPoll(ctx, *virtualMachineId, *model, virtualmachines.DefaultCreateOrUpdateOperationOptions())
							if err != nil {
								return fmt.Errorf("creating %s: %+v", id, err)
							}
						}
					}
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r VirtualMachineImplicitDataDiskResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient

			id, err := parse.DataDiskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
			resp, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
			}

			schema := VirtualMachineImplicitDataDiskResourceModel{
				Name:             id.Name,
				VirtualMachineId: virtualMachineId.ID(),
			}

			var disk *virtualmachines.DataDisk
			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if profile := props.StorageProfile; profile != nil {
						if dataDisks := profile.DataDisks; dataDisks != nil {
							for _, dataDisk := range *dataDisks {
								if pointer.From(dataDisk.Name) == id.Name {
									disk = &dataDisk
									break
								}
							}
						}
					}
				}
			}

			if disk == nil {
				return metadata.MarkAsGone(*id)
			}

			schema.Lun = disk.Lun
			schema.Caching = string(pointer.From(disk.Caching))
			schema.CreateOption = string(disk.CreateOption)
			schema.DeleteOption = string(pointer.From(disk.DeleteOption))
			schema.DiskSizeGb = pointer.From(disk.DiskSizeGB)
			if disk.SourceResource != nil {
				schema.SourceResourceId = pointer.From(disk.SourceResource.Id)
			}

			schema.WriteAcceleratorEnabled = pointer.From(disk.WriteAcceleratorEnabled)

			return metadata.Encode(&schema)
		},
	}
}

func (r VirtualMachineImplicitDataDiskResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient

			id, err := parse.DataDiskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.VirtualMachineName, VirtualMachineResourceName)
			defer locks.UnlockByName(id.VirtualMachineName, VirtualMachineResourceName)

			virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
			resp, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if profile := props.StorageProfile; profile != nil {
						newDisks := make([]virtualmachines.DataDisk, 0)
						var toBeDeletedDisk *virtualmachines.DataDisk
						if dataDisks := profile.DataDisks; dataDisks != nil {
							for _, dataDisk := range *dataDisks {
								if pointer.From(dataDisk.Name) != id.Name {
									newDisks = append(newDisks, dataDisk)
								} else {
									toBeDeletedDisk = &dataDisk
								}
							}
						}

						profile.DataDisks = &newDisks

						// fixes #24145
						model.Properties.ApplicationProfile = nil

						// fixes #2485
						model.Identity = nil
						// fixes #1600
						model.Resources = nil

						err = client.CreateOrUpdateThenPoll(ctx, virtualMachineId, *model, virtualmachines.DefaultCreateOrUpdateOperationOptions())
						if err != nil {
							return fmt.Errorf("deleting %s: %+v", id, err)
						}

						// delete data disk if delete_option is set to delete
						if toBeDeletedDisk != nil && pointer.From(toBeDeletedDisk.DeleteOption) == virtualmachines.DiskDeleteOptionTypesDelete &&
							toBeDeletedDisk.ManagedDisk != nil && toBeDeletedDisk.ManagedDisk.Id != nil {
							diskClient := metadata.Client.Compute.DisksClient
							diskId, err := commonids.ParseManagedDiskID(*toBeDeletedDisk.ManagedDisk.Id)
							if err != nil {
								return err
							}

							err = diskClient.DeleteThenPoll(ctx, *diskId)
							if err != nil {
								return fmt.Errorf("deleting Managed Disk %s: %+v", *diskId, err)
							}
						}
					}
				}
			}

			return nil
		},
	}
}

func (r VirtualMachineImplicitDataDiskResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient

			var config VirtualMachineImplicitDataDiskResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := parse.DataDiskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.VirtualMachineName, VirtualMachineResourceName)
			defer locks.UnlockByName(id.VirtualMachineName, VirtualMachineResourceName)

			virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
			resp, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if profile := props.StorageProfile; profile != nil {
						if dataDisks := profile.DataDisks; dataDisks != nil {
							existingIndex := -1
							disks := *dataDisks
							for i, disk := range disks {
								if pointer.From(disk.Name) == config.Name {
									existingIndex = i
									break
								}
							}

							if existingIndex == -1 {
								return fmt.Errorf("unable to retrieve the data disk %s ", *id)
							}

							expandedDisk := &disks[existingIndex]
							if metadata.ResourceData.HasChange("caching") {
								expandedDisk.Caching = pointer.To(virtualmachines.CachingTypes(config.Caching))
							}

							if metadata.ResourceData.HasChange("delete_option") {
								expandedDisk.DeleteOption = pointer.To(virtualmachines.DiskDeleteOptionTypes(config.DeleteOption))
							}

							if metadata.ResourceData.HasChange("disk_size_gb") {
								expandedDisk.DiskSizeGB = pointer.To(config.DiskSizeGb)
							}

							if metadata.ResourceData.HasChange("write_accelerator_enabled") {
								expandedDisk.WriteAcceleratorEnabled = pointer.To(config.WriteAcceleratorEnabled)
							}

							profile.DataDisks = &disks
							// fixes #24145
							model.Properties.ApplicationProfile = nil
							// fixes #2485
							model.Identity = nil
							// fixes #1600
							model.Resources = nil

							err = client.CreateOrUpdateThenPoll(ctx, virtualMachineId, *model, virtualmachines.DefaultCreateOrUpdateOperationOptions())
							if err != nil {
								return fmt.Errorf("updating %s: %+v", id, err)
							}
						}
					}
				}
			}

			return nil
		},
	}
}

func (r VirtualMachineImplicitDataDiskResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Compute.VirtualMachinesClient

		id, err := parse.DataDiskID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
		resp, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				if profile := props.StorageProfile; profile != nil {
					if dataDisks := profile.DataDisks; dataDisks != nil {
						var disk *virtualmachines.DataDisk
						for _, dataDisk := range *dataDisks {
							if pointer.From(dataDisk.Name) == id.Name {
								disk = &dataDisk
								break
							}
						}

						if disk == nil {
							return fmt.Errorf("unable to retrieve an existing data disk %s", *id)
						}

						if disk.CreateOption != virtualmachines.DiskCreateOptionTypesCopy {
							return fmt.Errorf("the value of `create_option` for the imported `azurerm_virtual_machine_implicit_data_disk` instance must be `Copy`, whereas now is %s", disk.CreateOption)
						}
					}
				}
			}
		}

		return nil
	}
}
