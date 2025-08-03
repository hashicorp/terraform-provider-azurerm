// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/availabilitysets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/clouds"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachineinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachinetemplates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SystemCenterVirtualMachineManagerVirtualMachineInstanceModel struct {
	ScopedResourceId                                    string             `tfschema:"scoped_resource_id"`
	CustomLocationId                                    string             `tfschema:"custom_location_id"`
	Hardware                                            []Hardware         `tfschema:"hardware"`
	Infrastructure                                      []Infrastructure   `tfschema:"infrastructure"`
	NetworkInterfaces                                   []NetworkInterface `tfschema:"network_interface"`
	OperatingSystem                                     []OperatingSystem  `tfschema:"operating_system"`
	StorageDisks                                        []StorageDisk      `tfschema:"storage_disk"`
	SystemCenterVirtualMachineManagerAvailabilitySetIds []string           `tfschema:"system_center_virtual_machine_manager_availability_set_ids"`
}

type Hardware struct {
	CpuCount                    int64 `tfschema:"cpu_count"`
	DynamicMemoryMaxInMb        int64 `tfschema:"dynamic_memory_max_in_mb"`
	DynamicMemoryMinInMb        int64 `tfschema:"dynamic_memory_min_in_mb"`
	LimitCpuForMigrationEnabled bool  `tfschema:"limit_cpu_for_migration_enabled"`
	MemoryInMb                  int64 `tfschema:"memory_in_mb"`
}

type Infrastructure struct {
	CheckpointType                                          string `tfschema:"checkpoint_type"`
	SystemCenterVirtualMachineManagerCloudId                string `tfschema:"system_center_virtual_machine_manager_cloud_id"`
	SystemCenterVirtualMachineManagerInventoryItemId        string `tfschema:"system_center_virtual_machine_manager_inventory_item_id"`
	SystemCenterVirtualMachineManagerTemplateId             string `tfschema:"system_center_virtual_machine_manager_template_id"`
	SystemCenterVirtualMachineManagerVirtualMachineServerId string `tfschema:"system_center_virtual_machine_manager_virtual_machine_server_id"`
}

type NetworkInterface struct {
	Name             string `tfschema:"name"`
	Ipv4AddressType  string `tfschema:"ipv4_address_type"`
	Ipv6AddressType  string `tfschema:"ipv6_address_type"`
	MacAddressType   string `tfschema:"mac_address_type"`
	VirtualNetworkId string `tfschema:"virtual_network_id"`
}

type OperatingSystem struct {
	ComputerName  string `tfschema:"computer_name"`
	AdminPassword string `tfschema:"admin_password"`
}

type StorageDisk struct {
	Bus                  int64  `tfschema:"bus"`
	BusType              string `tfschema:"bus_type"`
	DiskSizeGB           int64  `tfschema:"disk_size_gb"`
	Lun                  int64  `tfschema:"lun"`
	Name                 string `tfschema:"name"`
	StorageQoSPolicyName string `tfschema:"storage_qos_policy_name"`
	TemplateDiskId       string `tfschema:"template_disk_id"`
	VhdType              string `tfschema:"vhd_type"`
}

var (
	_ sdk.Resource           = SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{}
	_ sdk.ResourceWithUpdate = SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{}
)

type SystemCenterVirtualMachineManagerVirtualMachineInstanceResource struct{}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerVirtualMachineInstanceModel{}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SystemCenterVirtualMachineManagerVirtualMachineInstanceID
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_virtual_machine_instance"
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scoped_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: machines.ValidateMachineID,
		},

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}),

		"infrastructure": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"checkpoint_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Disabled",
							"Production",
							"ProductionOnly",
							"Standard",
						}, false),
					},

					"system_center_virtual_machine_manager_cloud_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: clouds.ValidateCloudID,
						RequiredWith: []string{"infrastructure.0.system_center_virtual_machine_manager_template_id"},
						AtLeastOneOf: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id", "infrastructure.0.system_center_virtual_machine_manager_inventory_item_id", "infrastructure.0.system_center_virtual_machine_manager_template_id", "infrastructure.0.system_center_virtual_machine_manager_virtual_machine_server_id"},
					},

					"system_center_virtual_machine_manager_inventory_item_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: inventoryitems.ValidateInventoryItemID,
						RequiredWith: []string{"infrastructure.0.system_center_virtual_machine_manager_virtual_machine_server_id"},
						AtLeastOneOf: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id", "infrastructure.0.system_center_virtual_machine_manager_inventory_item_id", "infrastructure.0.system_center_virtual_machine_manager_template_id", "infrastructure.0.system_center_virtual_machine_manager_virtual_machine_server_id"},
					},

					"system_center_virtual_machine_manager_template_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: virtualmachinetemplates.ValidateVirtualMachineTemplateID,
						RequiredWith: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id"},
						AtLeastOneOf: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id", "infrastructure.0.system_center_virtual_machine_manager_inventory_item_id", "infrastructure.0.system_center_virtual_machine_manager_template_id", "infrastructure.0.system_center_virtual_machine_manager_virtual_machine_server_id"},
					},

					"system_center_virtual_machine_manager_virtual_machine_server_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: vmmservers.ValidateVMmServerID,
						AtLeastOneOf: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id", "infrastructure.0.system_center_virtual_machine_manager_inventory_item_id", "infrastructure.0.system_center_virtual_machine_manager_template_id", "infrastructure.0.system_center_virtual_machine_manager_virtual_machine_server_id"},
					},
				},
			},
		},

		"hardware": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cpu_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 64),
						AtLeastOneOf: []string{"hardware.0.cpu_count", "hardware.0.dynamic_memory_max_in_mb", "hardware.0.dynamic_memory_min_in_mb", "hardware.0.limit_cpu_for_migration_enabled", "hardware.0.memory_in_mb"},
					},

					"dynamic_memory_max_in_mb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(32, 1048576),
						AtLeastOneOf: []string{"hardware.0.cpu_count", "hardware.0.dynamic_memory_max_in_mb", "hardware.0.dynamic_memory_min_in_mb", "hardware.0.limit_cpu_for_migration_enabled", "hardware.0.memory_in_mb"},
					},

					"dynamic_memory_min_in_mb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(32, 1048576),
						AtLeastOneOf: []string{"hardware.0.cpu_count", "hardware.0.dynamic_memory_max_in_mb", "hardware.0.dynamic_memory_min_in_mb", "hardware.0.limit_cpu_for_migration_enabled", "hardware.0.memory_in_mb"},
					},

					"limit_cpu_for_migration_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						AtLeastOneOf: []string{"hardware.0.cpu_count", "hardware.0.dynamic_memory_max_in_mb", "hardware.0.dynamic_memory_min_in_mb", "hardware.0.limit_cpu_for_migration_enabled", "hardware.0.memory_in_mb"},
					},

					"memory_in_mb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(32, 1048576),
						AtLeastOneOf: []string{"hardware.0.cpu_count", "hardware.0.dynamic_memory_max_in_mb", "hardware.0.dynamic_memory_min_in_mb", "hardware.0.limit_cpu_for_migration_enabled", "hardware.0.memory_in_mb"},
					},
				},
			},
		},

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: networkValidate.NetworkInterfaceName,
					},

					"virtual_network_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: virtualnetworks.ValidateVirtualNetworkID,
					},

					"ipv4_address_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(virtualmachineinstances.PossibleValuesForAllocationMethod(), false),
					},

					"ipv6_address_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(virtualmachineinstances.PossibleValuesForAllocationMethod(), false),
					},

					"mac_address_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(virtualmachineinstances.PossibleValuesForAllocationMethod(), false),
					},
				},
			},
		},

		"operating_system": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"computer_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validate.SystemCenterVirtualMachineManagerVirtualMachineInstanceComputerName,
						AtLeastOneOf: []string{"operating_system.0.computer_name", "operating_system.0.admin_password"},
					},

					"admin_password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"operating_system.0.computer_name", "operating_system.0.admin_password"},
					},
				},
			},
		},

		"storage_disk": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"bus": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 3),
					},

					"bus_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"IDE",
							"SCSI",
						}, false),
					},

					"disk_size_gb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"lun": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 63),
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.SystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDiskName,
					},

					"storage_qos_policy_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"template_disk_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsUUID,
					},

					"vhd_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Dynamic",
							"Fixed",
						}, false),
					},
				},
			},
		},

		"system_center_virtual_machine_manager_availability_set_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: availabilitysets.ValidateAvailabilitySetID,
			},
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineInstances

			var model SystemCenterVirtualMachineManagerVirtualMachineInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := parse.NewSystemCenterVirtualMachineManagerVirtualMachineInstanceID(model.ScopedResourceId)

			existing, err := client.Get(ctx, commonids.NewScopeID(id.Scope))
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := virtualmachineinstances.VirtualMachineInstance{
				ExtendedLocation: virtualmachineinstances.ExtendedLocation{
					Type: pointer.To("customLocation"),
					Name: pointer.To(model.CustomLocationId),
				},
				Properties: &virtualmachineinstances.VirtualMachineInstanceProperties{
					InfrastructureProfile: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForCreate(model.Infrastructure),
					HardwareProfile:       expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForCreate(model.Hardware, metadata.ResourceData),
					OsProfile:             expandSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(model.OperatingSystem),
				},
			}

			if v := model.NetworkInterfaces; v != nil {
				parameters.Properties.NetworkProfile = &virtualmachineinstances.NetworkProfile{
					NetworkInterfaces: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForCreate(v),
				}
			}

			if v := model.StorageDisks; v != nil {
				parameters.Properties.StorageProfile = &virtualmachineinstances.StorageProfile{
					Disks: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForCreate(v, metadata.ResourceData),
				}
			}

			if v := model.SystemCenterVirtualMachineManagerAvailabilitySetIds; v != nil {
				availabilitySets, err := expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(v)
				if err != nil {
					return err
				}
				parameters.Properties.AvailabilitySets = availabilitySets
			}

			if err := client.CreateOrUpdateThenPoll(ctx, commonids.NewScopeID(id.Scope), parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineInstances

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, commonids.NewScopeID(id.Scope))
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := SystemCenterVirtualMachineManagerVirtualMachineInstanceModel{
				ScopedResourceId: id.Scope,
			}

			if model := resp.Model; model != nil {
				state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)

				if props := model.Properties; props != nil {
					state.Hardware = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfile(props.HardwareProfile)
					state.Infrastructure = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfile(props.InfrastructureProfile)
					state.OperatingSystem = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(props.OsProfile, metadata.ResourceData.Get("operating_system.0.admin_password").(string))
					state.SystemCenterVirtualMachineManagerAvailabilitySetIds = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(props.AvailabilitySets)

					if v := props.NetworkProfile; v != nil {
						state.NetworkInterfaces = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfaces(v.NetworkInterfaces)
					}

					if v := props.StorageProfile; v != nil {
						state.StorageDisks = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisks(v.Disks)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineInstances

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerVirtualMachineInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := virtualmachineinstances.VirtualMachineInstanceUpdate{
				Properties: &virtualmachineinstances.VirtualMachineInstanceUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("infrastructure") {
				parameters.Properties.InfrastructureProfile = expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForUpdate(model.Infrastructure)
			}

			if metadata.ResourceData.HasChange("hardware") {
				parameters.Properties.HardwareProfile = expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForUpdate(model.Hardware, metadata.ResourceData)
			}

			if metadata.ResourceData.HasChange("network_interface") {
				parameters.Properties.NetworkProfile = &virtualmachineinstances.NetworkProfileUpdate{
					NetworkInterfaces: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForUpdate(model.NetworkInterfaces),
				}
			}

			if metadata.ResourceData.HasChange("storage_disk") {
				parameters.Properties.StorageProfile = &virtualmachineinstances.StorageProfileUpdate{
					Disks: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForUpdate(model.StorageDisks, metadata.ResourceData),
				}
			}

			if metadata.ResourceData.HasChange("system_center_virtual_machine_manager_availability_set_ids") {
				availabilitySets, err := expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(model.SystemCenterVirtualMachineManagerAvailabilitySetIds)
				if err != nil {
					return err
				}
				parameters.Properties.AvailabilitySets = availabilitySets
			}

			needToRestart := metadata.ResourceData.HasChange("hardware") || metadata.ResourceData.HasChange("network_interface") || metadata.ResourceData.HasChange("storage_disk")

			if needToRestart {
				if err := client.StopThenPoll(ctx, commonids.NewScopeID(id.Scope), virtualmachineinstances.StopVirtualMachineOptions{
					SkipShutdown: pointer.To(virtualmachineinstances.SkipShutdownFalse),
				}); err != nil {
					return fmt.Errorf("stopping %s: %+v", *id, err)
				}
			}

			if err := client.UpdateThenPoll(ctx, commonids.NewScopeID(id.Scope), parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if needToRestart {
				if err := client.StartThenPoll(ctx, commonids.NewScopeID(id.Scope)); err != nil {
					return fmt.Errorf("starting %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineInstances

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, commonids.NewScopeID(id.Scope), virtualmachineinstances.DeleteOperationOptions{
				DeleteFromHost: pointer.To(virtualmachineinstances.DeleteFromHostTrue),
				Force:          pointer.To(virtualmachineinstances.ForceDeleteTrue),
			}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForCreate(input []Infrastructure) *virtualmachineinstances.InfrastructureProfile {
	if len(input) == 0 {
		return nil
	}

	infrastructureProfile := input[0]

	result := virtualmachineinstances.InfrastructureProfile{}

	if v := infrastructureProfile.CheckpointType; v != "" {
		result.CheckpointType = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerCloudId; v != "" {
		result.CloudId = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerInventoryItemId; v != "" {
		result.InventoryItemId = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerTemplateId; v != "" {
		result.TemplateId = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerVirtualMachineServerId; v != "" {
		result.VMmServerId = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForCreate(input []Hardware, d *pluginsdk.ResourceData) *virtualmachineinstances.HardwareProfile {
	if len(input) == 0 {
		return nil
	}

	hardwareProfile := input[0]

	result := virtualmachineinstances.HardwareProfile{}

	// As TF always sets bool value to false when it isn't set, so it has to use d.GetRawConfig() to determine whether it is set in the tf config
	if v := d.GetRawConfig().AsValueMap()["hardware"].AsValueSlice()[0].AsValueMap()["limit_cpu_for_migration_enabled"]; !v.IsNull() {
		result.LimitCPUForMigration = pointer.To(virtualmachineinstances.LimitCPUForMigration(strconv.FormatBool(hardwareProfile.LimitCpuForMigrationEnabled)))
	}

	if v := hardwareProfile.CpuCount; v != 0 {
		result.CpuCount = pointer.To(v)
	}

	dynamicMemoryEnabled := false
	if hardwareProfile.DynamicMemoryMaxInMb != 0 || hardwareProfile.DynamicMemoryMinInMb != 0 {
		dynamicMemoryEnabled = true
	}
	result.DynamicMemoryEnabled = pointer.To(virtualmachineinstances.DynamicMemoryEnabled(strconv.FormatBool(dynamicMemoryEnabled)))

	if v := hardwareProfile.DynamicMemoryMaxInMb; v != 0 {
		result.DynamicMemoryMaxMB = pointer.To(v)
	}

	if v := hardwareProfile.DynamicMemoryMinInMb; v != 0 {
		result.DynamicMemoryMinMB = pointer.To(v)
	}

	if v := hardwareProfile.MemoryInMb; v != 0 {
		result.MemoryMB = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForCreate(input []NetworkInterface) *[]virtualmachineinstances.NetworkInterface {
	result := make([]virtualmachineinstances.NetworkInterface, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		networkInterface := virtualmachineinstances.NetworkInterface{
			Name: pointer.To(v.Name),
		}

		if ipv4AddressType := v.Ipv4AddressType; ipv4AddressType != "" {
			networkInterface.IPv4AddressType = pointer.To(virtualmachineinstances.AllocationMethod(ipv4AddressType))
		}

		if ipv6AddressType := v.Ipv6AddressType; ipv6AddressType != "" {
			networkInterface.IPv6AddressType = pointer.To(virtualmachineinstances.AllocationMethod(ipv6AddressType))
		}

		if macAddressType := v.MacAddressType; macAddressType != "" {
			networkInterface.MacAddressType = pointer.To(virtualmachineinstances.AllocationMethod(macAddressType))
		}

		if vnetId := v.VirtualNetworkId; vnetId != "" {
			networkInterface.VirtualNetworkId = pointer.To(vnetId)
		}

		result = append(result, networkInterface)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(input []OperatingSystem) *virtualmachineinstances.OsProfileForVMInstance {
	if len(input) == 0 {
		return nil
	}

	osProfile := input[0]

	result := virtualmachineinstances.OsProfileForVMInstance{}

	if v := osProfile.ComputerName; v != "" {
		result.ComputerName = pointer.To(v)
	}

	if v := osProfile.AdminPassword; v != "" {
		result.AdminPassword = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForCreate(input []StorageDisk, d *pluginsdk.ResourceData) *[]virtualmachineinstances.VirtualDisk {
	result := make([]virtualmachineinstances.VirtualDisk, 0)
	if len(input) == 0 {
		return &result
	}

	for k, v := range input {
		virtualDisk := virtualmachineinstances.VirtualDisk{}

		// As API allows zero value for this property, so TF has to use d.GetRawConfig() to determine whether it's set in tf config
		if bus := d.GetRawConfig().AsValueMap()["storage_disk"].AsValueSlice()[k].AsValueMap()["bus"]; !bus.IsNull() {
			virtualDisk.Bus = pointer.To(v.Bus)
		}

		if lun := d.GetRawConfig().AsValueMap()["storage_disk"].AsValueSlice()[k].AsValueMap()["lun"]; !lun.IsNull() {
			virtualDisk.Lun = pointer.To(v.Lun)
		}

		if busType := v.BusType; busType != "" {
			virtualDisk.BusType = pointer.To(busType)
		}

		if diskSizeGB := v.DiskSizeGB; diskSizeGB != 0 {
			virtualDisk.DiskSizeGB = pointer.To(diskSizeGB)
		}

		if name := v.Name; name != "" {
			virtualDisk.Name = pointer.To(name)
		}

		if templateDiskId := v.TemplateDiskId; templateDiskId != "" {
			virtualDisk.TemplateDiskId = pointer.To(templateDiskId)
		}

		if vhdType := v.VhdType; vhdType != "" {
			virtualDisk.VhdType = pointer.To(vhdType)
		}

		if storageQosPolicyName := v.StorageQoSPolicyName; storageQosPolicyName != "" {
			virtualDisk.StorageQoSPolicy = &virtualmachineinstances.StorageQosPolicyDetails{
				Name: pointer.To(storageQosPolicyName),
			}
		}

		result = append(result, virtualDisk)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(input []string) (*[]virtualmachineinstances.AvailabilitySetListItem, error) {
	result := make([]virtualmachineinstances.AvailabilitySetListItem, 0)
	if len(input) == 0 {
		return &result, nil
	}

	for _, v := range input {
		availabilitySetId, err := availabilitysets.ParseAvailabilitySetID(v)
		if err != nil {
			return nil, err
		}

		result = append(result, virtualmachineinstances.AvailabilitySetListItem{
			Id:   pointer.To(availabilitySetId.ID()),
			Name: pointer.To(availabilitySetId.AvailabilitySetName),
		})
	}

	return &result, nil
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForUpdate(input []Infrastructure) *virtualmachineinstances.InfrastructureProfileUpdate {
	if len(input) == 0 {
		return nil
	}

	infrastructureProfile := input[0]

	result := virtualmachineinstances.InfrastructureProfileUpdate{}

	if v := infrastructureProfile.CheckpointType; v != "" {
		result.CheckpointType = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForUpdate(input []Hardware, d *pluginsdk.ResourceData) *virtualmachineinstances.HardwareProfileUpdate {
	result := virtualmachineinstances.HardwareProfileUpdate{}

	if len(input) == 0 {
		return &result
	}

	hardwareProfile := input[0]

	// As TF always sets bool value to false when it isn't set, so it has to use d.GetRawConfig() to determine whether it is set in the tf config
	if v := d.GetRawConfig().AsValueMap()["hardware"].AsValueSlice()[0].AsValueMap()["limit_cpu_for_migration_enabled"]; !v.IsNull() {
		result.LimitCPUForMigration = pointer.To(virtualmachineinstances.LimitCPUForMigration(strconv.FormatBool(hardwareProfile.LimitCpuForMigrationEnabled)))
	}

	if v := hardwareProfile.CpuCount; v != 0 {
		result.CpuCount = pointer.To(v)
	}

	dynamicMemoryEnabled := false
	if hardwareProfile.DynamicMemoryMaxInMb != 0 || hardwareProfile.DynamicMemoryMinInMb != 0 {
		dynamicMemoryEnabled = true
	}
	result.DynamicMemoryEnabled = pointer.To(virtualmachineinstances.DynamicMemoryEnabled(strconv.FormatBool(dynamicMemoryEnabled)))

	if v := hardwareProfile.DynamicMemoryMaxInMb; v != 0 {
		result.DynamicMemoryMaxMB = pointer.To(v)
	}

	if v := hardwareProfile.DynamicMemoryMinInMb; v != 0 {
		result.DynamicMemoryMinMB = pointer.To(v)
	}

	if v := hardwareProfile.MemoryInMb; v != 0 {
		result.MemoryMB = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForUpdate(input []NetworkInterface) *[]virtualmachineinstances.NetworkInterfaceUpdate {
	result := make([]virtualmachineinstances.NetworkInterfaceUpdate, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		networkInterface := virtualmachineinstances.NetworkInterfaceUpdate{
			Name: pointer.To(v.Name),
		}

		if ipv4AddressType := v.Ipv4AddressType; ipv4AddressType != "" {
			networkInterface.IPv4AddressType = pointer.To(virtualmachineinstances.AllocationMethod(ipv4AddressType))
		}

		if ipv6AddressType := v.Ipv6AddressType; ipv6AddressType != "" {
			networkInterface.IPv6AddressType = pointer.To(virtualmachineinstances.AllocationMethod(ipv6AddressType))
		}

		if macAddressType := v.MacAddressType; macAddressType != "" {
			networkInterface.MacAddressType = pointer.To(virtualmachineinstances.AllocationMethod(macAddressType))
		}

		if vnetId := v.VirtualNetworkId; vnetId != "" {
			networkInterface.VirtualNetworkId = pointer.To(vnetId)
		}

		result = append(result, networkInterface)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForUpdate(input []StorageDisk, d *pluginsdk.ResourceData) *[]virtualmachineinstances.VirtualDiskUpdate {
	result := make([]virtualmachineinstances.VirtualDiskUpdate, 0)
	if len(input) == 0 {
		return &result
	}

	for k, v := range input {
		virtualDisk := virtualmachineinstances.VirtualDiskUpdate{}

		// As API allows zero value for this property, so TF has to use d.GetRawConfig() to determine whether it's set in tf config
		if bus := d.GetRawConfig().AsValueMap()["storage_disk"].AsValueSlice()[k].AsValueMap()["bus"]; !bus.IsNull() {
			virtualDisk.Bus = pointer.To(v.Bus)
		}

		if lun := d.GetRawConfig().AsValueMap()["storage_disk"].AsValueSlice()[k].AsValueMap()["lun"]; !lun.IsNull() {
			virtualDisk.Lun = pointer.To(v.Lun)
		}

		if busType := v.BusType; busType != "" {
			virtualDisk.BusType = pointer.To(busType)
		}

		if diskSizeGB := v.DiskSizeGB; diskSizeGB != 0 {
			virtualDisk.DiskSizeGB = pointer.To(diskSizeGB)
		}

		if name := v.Name; name != "" {
			virtualDisk.Name = pointer.To(name)
		}

		if vhdType := v.VhdType; vhdType != "" {
			virtualDisk.VhdType = pointer.To(vhdType)
		}

		if storageQosPolicyName := v.StorageQoSPolicyName; storageQosPolicyName != "" {
			virtualDisk.StorageQoSPolicy = &virtualmachineinstances.StorageQosPolicyDetails{
				Name: pointer.To(storageQosPolicyName),
			}
		}

		result = append(result, virtualDisk)
	}

	return &result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfile(input *virtualmachineinstances.InfrastructureProfile) []Infrastructure {
	result := make([]Infrastructure, 0)
	if input == nil {
		return result
	}

	return append(result, Infrastructure{
		CheckpointType:                                          pointer.From(input.CheckpointType),
		SystemCenterVirtualMachineManagerCloudId:                pointer.From(input.CloudId),
		SystemCenterVirtualMachineManagerInventoryItemId:        pointer.From(input.InventoryItemId),
		SystemCenterVirtualMachineManagerTemplateId:             pointer.From(input.TemplateId),
		SystemCenterVirtualMachineManagerVirtualMachineServerId: pointer.From(input.VMmServerId),
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfile(input *virtualmachineinstances.HardwareProfile) []Hardware {
	result := make([]Hardware, 0)
	if input == nil {
		return result
	}

	return append(result, Hardware{
		CpuCount:                    pointer.From(input.CpuCount),
		DynamicMemoryMaxInMb:        pointer.From(input.DynamicMemoryMaxMB),
		DynamicMemoryMinInMb:        pointer.From(input.DynamicMemoryMinMB),
		LimitCpuForMigrationEnabled: pointer.From(input.LimitCPUForMigration) == virtualmachineinstances.LimitCPUForMigrationTrue,
		MemoryInMb:                  pointer.From(input.MemoryMB),
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfaces(input *[]virtualmachineinstances.NetworkInterface) []NetworkInterface {
	result := make([]NetworkInterface, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, NetworkInterface{
			Name:             pointer.From(v.Name),
			VirtualNetworkId: pointer.From(v.VirtualNetworkId),
			Ipv4AddressType:  string(pointer.From(v.IPv4AddressType)),
			Ipv6AddressType:  string(pointer.From(v.IPv6AddressType)),
			MacAddressType:   string(pointer.From(v.MacAddressType)),
		})
	}

	return result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(input *virtualmachineinstances.OsProfileForVMInstance, adminPassword string) []OperatingSystem {
	result := make([]OperatingSystem, 0)
	if input == nil {
		return result
	}

	return append(result, OperatingSystem{
		ComputerName:  pointer.From(input.ComputerName),
		AdminPassword: adminPassword,
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisks(input *[]virtualmachineinstances.VirtualDisk) []StorageDisk {
	result := make([]StorageDisk, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		// service team confirmed that diskSizeGB in the request payload and maxDiskSizeGB in the API response both represent the maximum size of the disk but diskSizeGB in the API response represents the actual disk size
		storageDisk := StorageDisk{
			Bus:            pointer.From(v.Bus),
			BusType:        pointer.From(v.BusType),
			DiskSizeGB:     pointer.From(v.MaxDiskSizeGB),
			Lun:            pointer.From(v.Lun),
			Name:           pointer.From(v.Name),
			TemplateDiskId: pointer.From(v.TemplateDiskId),
			VhdType:        pointer.From(v.VhdType),
		}

		if storageQoSPolicy := v.StorageQoSPolicy; storageQoSPolicy != nil {
			storageDisk.StorageQoSPolicyName = pointer.From(storageQoSPolicy.Name)
		}

		result = append(result, storageDisk)
	}

	return result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(input *[]virtualmachineinstances.AvailabilitySetListItem) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, pointer.From(v.Id))
	}

	return result
}
