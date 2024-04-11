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
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/availabilitysets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachineinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerVirtualMachineInstanceModel struct {
	ScopedResourceId                                    string                  `tfschema:"scoped_resource_id"`
	CustomLocationId                                    string                  `tfschema:"custom_location_id"`
	HardwareProfile                                     []HardwareProfile       `tfschema:"hardware_profile"`
	InfrastructureProfile                               []InfrastructureProfile `tfschema:"infrastructure_profile"`
	NetworkInterfaces                                   []NetworkInterface      `tfschema:"network_interface"`
	OSProfile                                           []OSProfile             `tfschema:"os_profile"`
	SystemCenterVirtualMachineManagerAvailabilitySetIds []string                `tfschema:"system_center_virtual_machine_manager_availability_set_id"`
}

type HardwareProfile struct {
	CpuCount                    int  `tfschema:"cpu_count"`
	DynamicMemoryEnabled        bool `tfschema:"dynamic_memory_enabled"`
	DynamicMemoryMaxInMb        int  `tfschema:"dynamic_memory_max_in_mb"`
	DynamicMemoryMinInMb        int  `tfschema:"dynamic_memory_min_in_mb"`
	LimitCpuForMigrationEnabled bool `tfschema:"limit_cpu_for_migration_enabled"`
	MemoryInMb                  int  `tfschema:"memory_in_mb"`
}

type InfrastructureProfile struct {
	BiosGuid                                                string       `tfschema:"bios_guid"`
	CheckpointType                                          string       `tfschema:"checkpoint_type"`
	Checkpoints                                             []Checkpoint `tfschema:"checkpoint"`
	Generation                                              int          `tfschema:"generation"`
	SystemCenterVirtualMachineManagerCloudId                string       `tfschema:"system_center_virtual_machine_manager_cloud_id"`
	SystemCenterVirtualMachineManagerInventoryItemId        string       `tfschema:"system_center_virtual_machine_manager_inventory_item_id"`
	SystemCenterVirtualMachineManagerTemplateId             string       `tfschema:"system_center_virtual_machine_manager_template_id"`
	SystemCenterVirtualMachineManagerVirtualMachineName     string       `tfschema:"system_center_virtual_machine_manager_virtual_machine_name"`
	SystemCenterVirtualMachineManagerVirtualMachineServerId string       `tfschema:"system_center_virtual_machine_manager_virtual_machine_server_id"`
	Uuid                                                    string       `tfschema:"uuid"`
}

type Checkpoint struct {
	CheckpointId       string `tfschema:"checkpoint_id"`
	Description        string `tfschema:"description"`
	Name               string `tfschema:"name"`
	ParentCheckpointId string `tfschema:"parent_checkpoint_id"`
}

type NetworkInterface struct {
	id               string `tfschema:"id"`
	name             string `tfschema:"name"`
	Ipv4AddressType  string `tfschema:"ipv4_address_type"`
	Ipv6AddressType  string `tfschema:"ipv6_address_type"`
	MacAddress       string `tfschema:"mac_address"`
	MacAddressType   string `tfschema:"mac_address_type"`
	VirtualNetworkId string `tfschema:"virtual_network_id"`
}

type OSProfile struct {
	ComputerName  string `tfschema:"computer_name"`
	AdminPassword string `tfschema:"admin_password"`
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
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}),

		"hardware_profile ": {
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

		"infrastructure_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"bios_guid": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"checkpoint_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"checkpoint": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"checkpoint_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},

								"description": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},

								"parent_checkpoint_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},
							},
						},
					},

					"generation": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						ForceNew: true,
					},

					"system_center_virtual_machine_manager_cloud_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"system_center_virtual_machine_manager_inventory_item_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"system_center_virtual_machine_manager_template_id": {
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

		"os_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"computer_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"admin_password": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
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
				Properties: virtualmachineinstances.VirtualMachineInstanceProperties{
					HardwareProfile:       expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForCreate(model.HardwareProfile),
					InfrastructureProfile: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForCreate(model.InfrastructureProfile),
					NetworkProfile: &virtualmachineinstances.NetworkProfile{
						NetworkInterfaces: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForCreate(model.NetworkInterfaces),
					},
					OsProfile: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(model.OSProfile),
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

			state := SystemCenterVirtualMachineManagerVirtualMachineInstanceModel{}
			state.ScopedResourceId = id.Scope
			if model := resp.Model; model != nil {
				state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)
				state.HardwareProfile = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfile(model.Properties.HardwareProfile)
				state.InfrastructureProfile = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfile(model.Properties.InfrastructureProfile)
				state.OSProfile = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(model.Properties.OsProfile)
				state.SystemCenterVirtualMachineManagerAvailabilitySetIds = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(model.Properties.AvailabilitySets)

				if v := model.Properties.NetworkProfile; v != nil {
					state.NetworkInterfaces = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfaces(v.NetworkInterfaces)
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

			if metadata.ResourceData.HasChange("hardware_profile") {
				parameters.Properties.HardwareProfile = expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForUpdate(model.HardwareProfile)
			}

			if metadata.ResourceData.HasChange("infrastructure_profile") {
				parameters.Properties.InfrastructureProfile = expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForUpdate(model.InfrastructureProfile)
			}

			if metadata.ResourceData.HasChange("network_interface") {
				parameters.Properties.NetworkProfile = &virtualmachineinstances.NetworkProfileUpdate{
					NetworkInterfaces: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForUpdate(model.NetworkInterfaces),
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
				Force:          pointer.To(virtualmachineinstances.ForceTrue),
			}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForCreate(input []HardwareProfile) *virtualmachineinstances.HardwareProfile {
	if len(input) == 0 {
		return nil
	}

	hardwareProfile := input[0]

	return &virtualmachineinstances.HardwareProfile{
		CpuCount:             pointer.To(int64(hardwareProfile.CpuCount)),
		DynamicMemoryEnabled: pointer.To(virtualmachineinstances.DynamicMemoryEnabled(strconv.FormatBool(hardwareProfile.DynamicMemoryEnabled))),
		DynamicMemoryMaxMB:   pointer.To(int64(hardwareProfile.DynamicMemoryMaxInMb)),
		DynamicMemoryMinMB:   pointer.To(int64(hardwareProfile.DynamicMemoryMinInMb)),
		LimitCPUForMigration: pointer.To(virtualmachineinstances.LimitCPUForMigration(strconv.FormatBool(hardwareProfile.LimitCpuForMigrationEnabled))),
		MemoryMB:             pointer.To(int64(hardwareProfile.MemoryInMb)),
	}
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForCreate(input []InfrastructureProfile) *virtualmachineinstances.InfrastructureProfile {
	if len(input) == 0 {
		return nil
	}

	infrastructureProfile := input[0]

	return &virtualmachineinstances.InfrastructureProfile{
		BiosGuid:           pointer.To(infrastructureProfile.BiosGuid),
		CheckpointType:     pointer.To(infrastructureProfile.CheckpointType),
		Checkpoints:        expandSystemCenterVirtualMachineManagerVirtualMachineInstanceCheckpoints(infrastructureProfile.Checkpoints),
		Generation:         pointer.To(int64(infrastructureProfile.Generation)),
		CloudId:            pointer.To(infrastructureProfile.SystemCenterVirtualMachineManagerCloudId),
		InventoryItemId:    pointer.To(infrastructureProfile.SystemCenterVirtualMachineManagerInventoryItemId),
		TemplateId:         pointer.To(infrastructureProfile.SystemCenterVirtualMachineManagerTemplateId),
		VirtualMachineName: pointer.To(infrastructureProfile.SystemCenterVirtualMachineManagerVirtualMachineName),
		VMmServerId:        pointer.To(infrastructureProfile.SystemCenterVirtualMachineManagerVirtualMachineServerId),
		Uuid:               pointer.To(infrastructureProfile.Uuid),
	}
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceCheckpoints(input []Checkpoint) *[]virtualmachineinstances.Checkpoint {
	result := make([]virtualmachineinstances.Checkpoint, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		checkpoint := virtualmachineinstances.Checkpoint{
			CheckpointID:       pointer.To(v.CheckpointId),
			Name:               pointer.To(v.Name),
			Description:        pointer.To(v.Description),
			ParentCheckpointID: pointer.To(v.ParentCheckpointId),
		}

		result = append(result, checkpoint)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForUpdate(input []HardwareProfile) *virtualmachineinstances.HardwareProfileUpdate {
	if len(input) == 0 {
		return nil
	}

	hardwareProfile := input[0]

	return &virtualmachineinstances.HardwareProfileUpdate{
		CpuCount:             pointer.To(int64(hardwareProfile.CpuCount)),
		DynamicMemoryEnabled: pointer.To(virtualmachineinstances.DynamicMemoryEnabled(strconv.FormatBool(hardwareProfile.DynamicMemoryEnabled))),
		DynamicMemoryMaxMB:   pointer.To(int64(hardwareProfile.DynamicMemoryMaxInMb)),
		DynamicMemoryMinMB:   pointer.To(int64(hardwareProfile.DynamicMemoryMinInMb)),
		LimitCPUForMigration: pointer.To(virtualmachineinstances.LimitCPUForMigration(strconv.FormatBool(hardwareProfile.LimitCpuForMigrationEnabled))),
		MemoryMB:             pointer.To(int64(hardwareProfile.MemoryInMb)),
	}
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForUpdate(input []InfrastructureProfile) *virtualmachineinstances.InfrastructureProfileUpdate {
	if len(input) == 0 {
		return nil
	}

	infrastructureProfile := input[0]

	return &virtualmachineinstances.InfrastructureProfileUpdate{
		CheckpointType: pointer.To(infrastructureProfile.CheckpointType),
	}
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForCreate(input []NetworkInterface) *[]virtualmachineinstances.NetworkInterface {
	result := make([]virtualmachineinstances.NetworkInterface, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		networkInterface := virtualmachineinstances.NetworkInterface{
			NicId:            pointer.To(v.id),
			Name:             pointer.To(v.name),
			VirtualNetworkId: pointer.To(v.VirtualNetworkId),
			IPv4AddressType:  pointer.To(virtualmachineinstances.AllocationMethod(v.Ipv4AddressType)),
			IPv6AddressType:  pointer.To(virtualmachineinstances.AllocationMethod(v.Ipv6AddressType)),
			MacAddress:       pointer.To(v.MacAddress),
			MacAddressType:   pointer.To(virtualmachineinstances.AllocationMethod(v.MacAddressType)),
		}

		result = append(result, networkInterface)
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
			NicId:            pointer.To(v.id),
			Name:             pointer.To(v.name),
			VirtualNetworkId: pointer.To(v.VirtualNetworkId),
			IPv4AddressType:  pointer.To(virtualmachineinstances.AllocationMethod(v.Ipv4AddressType)),
			IPv6AddressType:  pointer.To(virtualmachineinstances.AllocationMethod(v.Ipv6AddressType)),
			MacAddress:       pointer.To(v.MacAddress),
			MacAddressType:   pointer.To(virtualmachineinstances.AllocationMethod(v.MacAddressType)),
		}

		result = append(result, networkInterface)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(input []OSProfile) *virtualmachineinstances.OsProfileForVMInstance {
	if len(input) == 0 {
		return nil
	}

	osProfile := input[0]

	return &virtualmachineinstances.OsProfileForVMInstance{
		ComputerName:  pointer.To(osProfile.ComputerName),
		AdminPassword: pointer.To(osProfile.AdminPassword),
	}
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(input []string) (*[]virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined, error) {
	result := make([]virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined, 0)
	if len(input) == 0 {
		return &result, nil
	}

	for _, v := range input {
		availabilitySetId, err := availabilitysets.ParseAvailabilitySetID(v)
		if err != nil {
			return nil, err
		}

		result = append(result, virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined{
			Id:   pointer.To(availabilitySetId.ID()),
			Name: pointer.To(availabilitySetId.AvailabilitySetName),
		})
	}

	return &result, nil
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfile(input *virtualmachineinstances.HardwareProfile) []HardwareProfile {
	result := make([]HardwareProfile, 0)
	if input == nil {
		return result
	}

	return append(result, HardwareProfile{
		CpuCount:                    int(pointer.From(input.CpuCount)),
		DynamicMemoryEnabled:        pointer.From(input.DynamicMemoryEnabled) == virtualmachineinstances.DynamicMemoryEnabledTrue,
		DynamicMemoryMaxInMb:        int(pointer.From(input.DynamicMemoryMaxMB)),
		DynamicMemoryMinInMb:        int(pointer.From(input.DynamicMemoryMinMB)),
		LimitCpuForMigrationEnabled: pointer.From(input.LimitCPUForMigration) == virtualmachineinstances.LimitCPUForMigrationTrue,
		MemoryInMb:                  int(pointer.From(input.MemoryMB)),
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfile(input *virtualmachineinstances.InfrastructureProfile) []InfrastructureProfile {
	result := make([]InfrastructureProfile, 0)
	if input == nil {
		return result
	}

	return append(result, InfrastructureProfile{
		BiosGuid:                                 pointer.From(input.BiosGuid),
		CheckpointType:                           pointer.From(input.CheckpointType),
		Checkpoints:                              flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceCheckpoints(input.Checkpoints),
		Generation:                               int(pointer.From(input.Generation)),
		SystemCenterVirtualMachineManagerCloudId: pointer.From(input.CloudId),
		SystemCenterVirtualMachineManagerInventoryItemId:        pointer.From(input.InventoryItemId),
		SystemCenterVirtualMachineManagerTemplateId:             pointer.From(input.TemplateId),
		SystemCenterVirtualMachineManagerVirtualMachineName:     pointer.From(input.VirtualMachineName),
		SystemCenterVirtualMachineManagerVirtualMachineServerId: pointer.From(input.VMmServerId),
		Uuid: pointer.From(input.Uuid),
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceCheckpoints(input *[]virtualmachineinstances.Checkpoint) []Checkpoint {
	result := make([]Checkpoint, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, Checkpoint{
			CheckpointId:       pointer.From(v.CheckpointID),
			Name:               pointer.From(v.Name),
			Description:        pointer.From(v.Description),
			ParentCheckpointId: pointer.From(v.ParentCheckpointID),
		})
	}

	return result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfaces(input *[]virtualmachineinstances.NetworkInterface) []NetworkInterface {
	result := make([]NetworkInterface, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, NetworkInterface{
			id:               pointer.From(v.NicId),
			name:             pointer.From(v.Name),
			VirtualNetworkId: pointer.From(v.VirtualNetworkId),
			MacAddress:       pointer.From(v.MacAddress),
			Ipv4AddressType:  string(pointer.From(v.IPv4AddressType)),
			Ipv6AddressType:  string(pointer.From(v.IPv6AddressType)),
			MacAddressType:   string(pointer.From(v.MacAddressType)),
		})
	}

	return result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(input *virtualmachineinstances.OsProfileForVMInstance) []OSProfile {
	result := make([]OSProfile, 0)
	if input == nil {
		return result
	}

	return append(result, OSProfile{
		ComputerName:  pointer.From(input.ComputerName),
		AdminPassword: pointer.From(input.AdminPassword),
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(input *[]virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, pointer.From(v.Id))
	}

	return result
}
