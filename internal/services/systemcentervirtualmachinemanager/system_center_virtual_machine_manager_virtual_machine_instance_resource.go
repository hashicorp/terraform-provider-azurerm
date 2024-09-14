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
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachineinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerVirtualMachineInstanceModel struct {
	ScopedResourceId                                    string             `tfschema:"scoped_resource_id"`
	CustomLocationId                                    string             `tfschema:"custom_location_id"`
	Hardware                                            []Hardware         `tfschema:"hardware"`
	Infrastructure                                      []Infrastructure   `tfschema:"infrastructure"`
	NetworkInterfaces                                   []NetworkInterface `tfschema:"network_interface"`
	OS                                                  []OS               `tfschema:"os"`
	StorageDisks                                        []StorageDisk      `tfschema:"storage_disk"`
	SystemCenterVirtualMachineManagerAvailabilitySetIds []string           `tfschema:"system_center_virtual_machine_manager_availability_set_ids"`
}

type Hardware struct {
	CpuCount                    int64 `tfschema:"cpu_count"`
	DynamicMemoryEnabled        bool  `tfschema:"dynamic_memory_enabled"`
	DynamicMemoryMaxInMb        int64 `tfschema:"dynamic_memory_max_in_mb"`
	DynamicMemoryMinInMb        int64 `tfschema:"dynamic_memory_min_in_mb"`
	LimitCpuForMigrationEnabled bool  `tfschema:"limit_cpu_for_migration_enabled"`
	MemoryInMb                  int64 `tfschema:"memory_in_mb"`
}

type Infrastructure struct {
	BiosGuid                                                string `tfschema:"bios_guid"`
	CheckpointType                                          string `tfschema:"checkpoint_type"`
	Generation                                              int64  `tfschema:"generation"`
	SystemCenterVirtualMachineManagerCloudId                string `tfschema:"system_center_virtual_machine_manager_cloud_id"`
	SystemCenterVirtualMachineManagerInventoryItemId        string `tfschema:"system_center_virtual_machine_manager_inventory_item_id"`
	SystemCenterVirtualMachineManagerTemplateId             string `tfschema:"system_center_virtual_machine_manager_template_id"`
	SystemCenterVirtualMachineManagerVirtualMachineName     string `tfschema:"system_center_virtual_machine_manager_virtual_machine_name"`
	SystemCenterVirtualMachineManagerVirtualMachineServerId string `tfschema:"system_center_virtual_machine_manager_virtual_machine_server_id"`
	Uuid                                                    string `tfschema:"uuid"`
}

type NetworkInterface struct {
	Id               string `tfschema:"id"`
	Name             string `tfschema:"name"`
	Ipv4AddressType  string `tfschema:"ipv4_address_type"`
	Ipv6AddressType  string `tfschema:"ipv6_address_type"`
	MacAddress       string `tfschema:"mac_address"`
	MacAddressType   string `tfschema:"mac_address_type"`
	VirtualNetworkId string `tfschema:"virtual_network_id"`
}

type OS struct {
	ComputerName  string `tfschema:"computer_name"`
	AdminPassword string `tfschema:"admin_password"`
}

type StorageDisk struct {
	Bus                   int64              `tfschema:"bus"`
	BusType               string             `tfschema:"bus_type"`
	CreateDiffDiskEnabled bool               `tfschema:"create_diff_disk_enabled"`
	DiskId                string             `tfschema:"disk_id"`
	DiskSizeGB            int64              `tfschema:"disk_size_gb"`
	Lun                   int64              `tfschema:"lun"`
	Name                  string             `tfschema:"name"`
	StorageQoSPolicy      []StorageQoSPolicy `tfschema:"storage_qos_policy"`
	TemplateDiskId        string             `tfschema:"template_disk_id"`
	VhdType               string             `tfschema:"vhd_type"`
}

type StorageQoSPolicy struct {
	Id   string `tfschema:"id"`
	Name string `tfschema:"name"`
}

var _ sdk.Resource = SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{}
var _ sdk.ResourceWithUpdate = SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{}

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
					"system_center_virtual_machine_manager_cloud_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						RequiredWith: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id", "infrastructure.0.system_center_virtual_machine_manager_template_id"},
					},

					"system_center_virtual_machine_manager_template_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						RequiredWith: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id", "infrastructure.0.system_center_virtual_machine_manager_template_id"},
					},

					"bios_guid": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"checkpoint_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"generation": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						ForceNew: true,
					},

					"system_center_virtual_machine_manager_inventory_item_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"system_center_virtual_machine_manager_virtual_machine_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"system_center_virtual_machine_manager_virtual_machine_server_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"uuid": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
				},
			},
		},

		"hardware": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cpu_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"dynamic_memory_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"dynamic_memory_max_in_mb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"dynamic_memory_min_in_mb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"limit_cpu_for_migration_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"memory_in_mb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"virtual_network_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"ipv4_address_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"ipv6_address_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"mac_address": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"mac_address_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"os": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"computer_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"admin_password": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						ForceNew:  true,
						Sensitive: true,
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
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"bus_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"create_diff_disk_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"disk_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"disk_size_gb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"lun": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"storage_qos_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"template_disk_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"vhd_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"system_center_virtual_machine_manager_availability_set_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
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
		Timeout: 30 * time.Minute,
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
					Type: utils.String("customLocation"),
					Name: utils.String(model.CustomLocationId),
				},
				Properties: &virtualmachineinstances.VirtualMachineInstanceProperties{
					HardwareProfile:       expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForCreate(model.Hardware),
					InfrastructureProfile: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForCreate(model.Infrastructure),
					NetworkProfile: &virtualmachineinstances.NetworkProfile{
						NetworkInterfaces: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForCreate(model.NetworkInterfaces),
					},
					StorageProfile: &virtualmachineinstances.StorageProfile{
						Disks: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForCreate(model.StorageDisks),
					},
					OsProfile: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(model.OS),
				},
			}

			availabilitySets, err := expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(model.SystemCenterVirtualMachineManagerAvailabilitySetIds)
			if err != nil {
				return err
			}
			parameters.Properties.AvailabilitySets = availabilitySets

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
				state.Hardware = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfile(model.Properties.HardwareProfile)
				state.Infrastructure = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfile(model.Properties.InfrastructureProfile)
				state.OS = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(model.Properties.OsProfile)
				state.SystemCenterVirtualMachineManagerAvailabilitySetIds = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(model.Properties.AvailabilitySets)

				if v := model.Properties.NetworkProfile; v != nil {
					state.NetworkInterfaces = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfaces(v.NetworkInterfaces)
				}

				if v := model.Properties.StorageProfile; v != nil {
					state.StorageDisks = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisks(v.Disks)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
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

			if metadata.ResourceData.HasChange("system_center_virtual_machine_manager_availability_set_ids") {
				availabilitySets, err := expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(model.SystemCenterVirtualMachineManagerAvailabilitySetIds)
				if err != nil {
					return err
				}
				parameters.Properties.AvailabilitySets = availabilitySets
			}

			if metadata.ResourceData.HasChange("hardware") {
				parameters.Properties.HardwareProfile = expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForUpdate(model.Hardware)
			}

			if metadata.ResourceData.HasChange("infrastructure") {
				parameters.Properties.InfrastructureProfile = expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForUpdate(model.Infrastructure)
			}

			if metadata.ResourceData.HasChange("network_interface") {
				parameters.Properties.NetworkProfile = &virtualmachineinstances.NetworkProfileUpdate{
					NetworkInterfaces: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForUpdate(model.NetworkInterfaces),
				}
			}

			if metadata.ResourceData.HasChange("storage_disk") {
				parameters.Properties.StorageProfile = &virtualmachineinstances.StorageProfileUpdate{
					Disks: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForUpdate(model.StorageDisks),
				}
			}

			if err := client.UpdateThenPoll(ctx, commonids.NewScopeID(id.Scope), parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
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

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForCreate(input []Hardware) *virtualmachineinstances.HardwareProfile {
	if len(input) == 0 {
		return nil
	}

	hardwareProfile := input[0]

	return &virtualmachineinstances.HardwareProfile{
		CpuCount:             pointer.To(hardwareProfile.CpuCount),
		DynamicMemoryEnabled: pointer.To(virtualmachineinstances.DynamicMemoryEnabled(strconv.FormatBool(hardwareProfile.DynamicMemoryEnabled))),
		DynamicMemoryMaxMB:   pointer.To(hardwareProfile.DynamicMemoryMaxInMb),
		DynamicMemoryMinMB:   pointer.To(hardwareProfile.DynamicMemoryMinInMb),
		LimitCPUForMigration: pointer.To(virtualmachineinstances.LimitCPUForMigration(strconv.FormatBool(hardwareProfile.LimitCpuForMigrationEnabled))),
		MemoryMB:             pointer.To(hardwareProfile.MemoryInMb),
	}
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForCreate(input []Infrastructure) *virtualmachineinstances.InfrastructureProfile {
	if len(input) == 0 {
		return nil
	}

	infrastructureProfile := input[0]

	result := virtualmachineinstances.InfrastructureProfile{
		Generation: pointer.To(infrastructureProfile.Generation),
	}

	if v := infrastructureProfile.BiosGuid; v != "" {
		result.BiosGuid = pointer.To(v)
	}

	if v := infrastructureProfile.CheckpointType; v != "" {
		result.CheckpointType = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerCloudId; v != "" {
		result.CloudId = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerTemplateId; v != "" {
		result.TemplateId = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerVirtualMachineName; v != "" {
		result.VirtualMachineName = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerVirtualMachineServerId; v != "" {
		result.VMmServerId = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerInventoryItemId; v != "" {
		result.InventoryItemId = pointer.To(v)
	}

	if v := infrastructureProfile.Uuid; v != "" {
		result.Uuid = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForUpdate(input []Hardware) *virtualmachineinstances.HardwareProfileUpdate {
	if len(input) == 0 {
		return nil
	}

	hardwareProfile := input[0]

	return &virtualmachineinstances.HardwareProfileUpdate{
		CpuCount:             pointer.To(hardwareProfile.CpuCount),
		DynamicMemoryEnabled: pointer.To(virtualmachineinstances.DynamicMemoryEnabled(strconv.FormatBool(hardwareProfile.DynamicMemoryEnabled))),
		DynamicMemoryMaxMB:   pointer.To(hardwareProfile.DynamicMemoryMaxInMb),
		DynamicMemoryMinMB:   pointer.To(hardwareProfile.DynamicMemoryMinInMb),
		LimitCPUForMigration: pointer.To(virtualmachineinstances.LimitCPUForMigration(strconv.FormatBool(hardwareProfile.LimitCpuForMigrationEnabled))),
		MemoryMB:             pointer.To(hardwareProfile.MemoryInMb),
	}
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

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForCreate(input []NetworkInterface) *[]virtualmachineinstances.NetworkInterface {
	result := make([]virtualmachineinstances.NetworkInterface, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		networkInterface := virtualmachineinstances.NetworkInterface{}

		if nicId := v.Id; nicId != "" {
			networkInterface.NicId = pointer.To(nicId)
		}

		if name := v.Name; name != "" {
			networkInterface.Name = pointer.To(name)
		}

		if vnetId := v.VirtualNetworkId; vnetId != "" {
			networkInterface.VirtualNetworkId = pointer.To(vnetId)
		}

		if ipv4AddressType := v.Ipv4AddressType; ipv4AddressType != "" {
			networkInterface.IPv4AddressType = pointer.To(virtualmachineinstances.AllocationMethod(ipv4AddressType))
		}

		if ipv6AddressType := v.Ipv6AddressType; ipv6AddressType != "" {
			networkInterface.IPv6AddressType = pointer.To(virtualmachineinstances.AllocationMethod(ipv6AddressType))
		}

		if macAddress := v.MacAddress; macAddress != "" {
			networkInterface.MacAddress = pointer.To(macAddress)
		}

		if macAddressType := v.MacAddressType; macAddressType != "" {
			networkInterface.MacAddressType = pointer.To(virtualmachineinstances.AllocationMethod(macAddressType))
		}

		result = append(result, networkInterface)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForCreate(input []StorageDisk) *[]virtualmachineinstances.VirtualDisk {
	result := make([]virtualmachineinstances.VirtualDisk, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		createDiffDisk := virtualmachineinstances.CreateDiffDiskFalse
		if v.CreateDiffDiskEnabled {
			createDiffDisk = virtualmachineinstances.CreateDiffDiskTrue
		}

		virtualDisk := virtualmachineinstances.VirtualDisk{
			Bus:              pointer.To(v.Bus),
			CreateDiffDisk:   pointer.To(createDiffDisk),
			DiskSizeGB:       pointer.To(v.DiskSizeGB),
			Lun:              pointer.To(v.Lun),
			StorageQoSPolicy: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageQoSPolicy(v.StorageQoSPolicy),
		}

		if busType := v.BusType; busType != "" {
			virtualDisk.BusType = pointer.To(busType)
		}

		if diskId := v.DiskId; diskId != "" {
			virtualDisk.DiskId = pointer.To(diskId)
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

		result = append(result, virtualDisk)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageQoSPolicy(input []StorageQoSPolicy) *virtualmachineinstances.StorageQosPolicyDetails {
	if len(input) == 0 {
		return nil
	}

	storageQoSPolicy := input[0]

	result := virtualmachineinstances.StorageQosPolicyDetails{}

	if v := storageQoSPolicy.Id; v != "" {
		result.Id = pointer.To(v)
	}

	if v := storageQoSPolicy.Name; v != "" {
		result.Name = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForUpdate(input []NetworkInterface) *[]virtualmachineinstances.NetworkInterfaceUpdate {
	result := make([]virtualmachineinstances.NetworkInterfaceUpdate, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		networkInterface := virtualmachineinstances.NetworkInterfaceUpdate{}

		if nicId := v.Id; nicId != "" {
			networkInterface.NicId = pointer.To(nicId)
		}

		if name := v.Name; name != "" {
			networkInterface.Name = pointer.To(name)
		}

		if vnetId := v.VirtualNetworkId; vnetId != "" {
			networkInterface.VirtualNetworkId = pointer.To(vnetId)
		}

		if ipv4AddressType := v.Ipv4AddressType; ipv4AddressType != "" {
			networkInterface.IPv4AddressType = pointer.To(virtualmachineinstances.AllocationMethod(ipv4AddressType))
		}

		if ipv6AddressType := v.Ipv6AddressType; ipv6AddressType != "" {
			networkInterface.IPv6AddressType = pointer.To(virtualmachineinstances.AllocationMethod(ipv6AddressType))
		}

		if macAddress := v.MacAddress; macAddress != "" {
			networkInterface.MacAddress = pointer.To(macAddress)
		}

		if macAddressType := v.MacAddressType; macAddressType != "" {
			networkInterface.MacAddressType = pointer.To(virtualmachineinstances.AllocationMethod(macAddressType))
		}

		result = append(result, networkInterface)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForUpdate(input []StorageDisk) *[]virtualmachineinstances.VirtualDiskUpdate {
	result := make([]virtualmachineinstances.VirtualDiskUpdate, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		virtualDisk := virtualmachineinstances.VirtualDiskUpdate{
			Bus:              pointer.To(v.Bus),
			DiskSizeGB:       pointer.To(v.DiskSizeGB),
			Lun:              pointer.To(v.Lun),
			StorageQoSPolicy: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageQoSPolicy(v.StorageQoSPolicy),
		}

		if busType := v.BusType; busType != "" {
			virtualDisk.BusType = pointer.To(busType)
		}

		if diskId := v.DiskId; diskId != "" {
			virtualDisk.DiskId = pointer.To(diskId)
		}

		if name := v.Name; name != "" {
			virtualDisk.Name = pointer.To(name)
		}

		if vhdType := v.VhdType; vhdType != "" {
			virtualDisk.VhdType = pointer.To(vhdType)
		}

		result = append(result, virtualDisk)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(input []OS) *virtualmachineinstances.OsProfileForVMInstance {
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

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfile(input *virtualmachineinstances.HardwareProfile) []Hardware {
	result := make([]Hardware, 0)
	if input == nil {
		return result
	}

	return append(result, Hardware{
		CpuCount:                    pointer.From(input.CpuCount),
		DynamicMemoryEnabled:        pointer.From(input.DynamicMemoryEnabled) == virtualmachineinstances.DynamicMemoryEnabledTrue,
		DynamicMemoryMaxInMb:        pointer.From(input.DynamicMemoryMaxMB),
		DynamicMemoryMinInMb:        pointer.From(input.DynamicMemoryMinMB),
		LimitCpuForMigrationEnabled: pointer.From(input.LimitCPUForMigration) == virtualmachineinstances.LimitCPUForMigrationTrue,
		MemoryInMb:                  pointer.From(input.MemoryMB),
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfile(input *virtualmachineinstances.InfrastructureProfile) []Infrastructure {
	result := make([]Infrastructure, 0)
	if input == nil {
		return result
	}

	return append(result, Infrastructure{
		BiosGuid:                                 pointer.From(input.BiosGuid),
		CheckpointType:                           pointer.From(input.CheckpointType),
		Generation:                               pointer.From(input.Generation),
		SystemCenterVirtualMachineManagerCloudId: pointer.From(input.CloudId),
		SystemCenterVirtualMachineManagerInventoryItemId:        pointer.From(input.InventoryItemId),
		SystemCenterVirtualMachineManagerTemplateId:             pointer.From(input.TemplateId),
		SystemCenterVirtualMachineManagerVirtualMachineName:     pointer.From(input.VirtualMachineName),
		SystemCenterVirtualMachineManagerVirtualMachineServerId: pointer.From(input.VMmServerId),
		Uuid: pointer.From(input.Uuid),
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfaces(input *[]virtualmachineinstances.NetworkInterface) []NetworkInterface {
	result := make([]NetworkInterface, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, NetworkInterface{
			Id:               pointer.From(v.NicId),
			Name:             pointer.From(v.Name),
			VirtualNetworkId: pointer.From(v.VirtualNetworkId),
			MacAddress:       pointer.From(v.MacAddress),
			Ipv4AddressType:  string(pointer.From(v.IPv4AddressType)),
			Ipv6AddressType:  string(pointer.From(v.IPv6AddressType)),
			MacAddressType:   string(pointer.From(v.MacAddressType)),
		})
	}

	return result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisks(input *[]virtualmachineinstances.VirtualDisk) []StorageDisk {
	result := make([]StorageDisk, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, StorageDisk{
			Bus:                   pointer.From(v.Bus),
			BusType:               pointer.From(v.BusType),
			CreateDiffDiskEnabled: pointer.From(v.CreateDiffDisk) == virtualmachineinstances.CreateDiffDiskTrue,
			DiskId:                pointer.From(v.DiskId),
			DiskSizeGB:            pointer.From(v.DiskSizeGB),
			Lun:                   pointer.From(v.Lun),
			Name:                  pointer.From(v.Name),
			StorageQoSPolicy:      flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageQoSPolicy(v.StorageQoSPolicy),
			TemplateDiskId:        pointer.From(v.TemplateDiskId),
			VhdType:               pointer.From(v.VhdType),
		})
	}

	return result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageQoSPolicy(input *virtualmachineinstances.StorageQosPolicyDetails) []StorageQoSPolicy {
	result := make([]StorageQoSPolicy, 0)
	if input == nil {
		return result
	}

	return append(result, StorageQoSPolicy{
		Id:   pointer.From(input.Id),
		Name: pointer.From(input.Name),
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(input *virtualmachineinstances.OsProfileForVMInstance) []OS {
	result := make([]OS, 0)
	if input == nil {
		return result
	}

	return append(result, OS{
		ComputerName:  pointer.From(input.ComputerName),
		AdminPassword: pointer.From(input.AdminPassword),
	})
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
