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
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachineinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerVirtualMachineInstanceModel struct {
	ScopedResourceId                                    string             `tfschema:"scoped_resource_id"`
	CustomLocationId                                    string             `tfschema:"custom_location_id"`
	SystemCenterVirtualMachineManagerAvailabilitySetIds []string           `tfschema:"system_center_virtual_machine_manager_availability_set_id"`
	HardwareProfile                                     []HardwareProfile  `tfschema:"hardware_profile"`
	NetworkInterfaces                                   []NetworkInterface `tfschema:"network_interface"`
	OSProfile                                           []OSProfile        `tfschema:"os_profile"`
}

type HardwareProfile struct {
	CpuCount                    int  `tfschema:"cpu_count"`
	DynamicMemoryEnabled        bool `tfschema:"dynamic_memory_enabled"`
	DynamicMemoryMaxInMb        int  `tfschema:"dynamic_memory_max_in_mb"`
	DynamicMemoryMinInMb        int  `tfschema:"dynamic_memory_min_in_mb"`
	LimitCpuForMigrationEnabled bool `tfschema:"limit_cpu_for_migration_enabled"`
	MemoryInMb                  int  `tfschema:"memory_in_mb"`
}

type NetworkInterface struct {
	id               string `tfschema:"id"`
	name             string `tfschema:"name"`
	VirtualNetworkId string `tfschema:"virtual_network_id"`
	Ipv4AddressType  string `tfschema:"ipv4_address_type"`
	Ipv6AddressType  string `tfschema:"ipv6_address_type"`
	MacAddress       string `tfschema:"mac_address"`
	MacAddressType   string `tfschema:"mac_address_type"`
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
				ValidateFunc: validation.StringIsNotEmpty,
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
					AvailabilitySets: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(model.SystemCenterVirtualMachineManagerAvailabilitySetIds),
					HardwareProfile:  expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForCreate(model.HardwareProfile),
					NetworkProfile: &virtualmachineinstances.NetworkProfile{
						NetworkInterfaces: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForCreate(model.NetworkInterfaces),
					},
					OsProfile: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(model.OSProfile),
				},
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

			state := SystemCenterVirtualMachineManagerVirtualMachineInstanceModel{}
			if model := resp.Model; model != nil {
				state.ScopedResourceId = id.Scope
				state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)
				state.HardwareProfile = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfile(model.Properties.HardwareProfile)
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
				parameters.Properties.AvailabilitySets = expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(model.SystemCenterVirtualMachineManagerAvailabilitySetIds)
			}

			if metadata.ResourceData.HasChange("hardware_profile") {
				parameters.Properties.HardwareProfile = expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForUpdate(model.HardwareProfile)
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

	hardwareProfile := &input[0]

	result := &virtualmachineinstances.HardwareProfile{
		CpuCount:             pointer.To(int64(hardwareProfile.CpuCount)),
		DynamicMemoryEnabled: pointer.To(virtualmachineinstances.DynamicMemoryEnabled(strconv.FormatBool(hardwareProfile.DynamicMemoryEnabled))),
		DynamicMemoryMaxMB:   pointer.To(int64(hardwareProfile.DynamicMemoryMaxInMb)),
		DynamicMemoryMinMB:   pointer.To(int64(hardwareProfile.DynamicMemoryMinInMb)),
		LimitCPUForMigration: pointer.To(virtualmachineinstances.LimitCPUForMigration(strconv.FormatBool(hardwareProfile.LimitCpuForMigrationEnabled))),
		MemoryMB:             pointer.To(int64(hardwareProfile.MemoryInMb)),
	}

	return result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForUpdate(input []HardwareProfile) *virtualmachineinstances.HardwareProfileUpdate {
	if len(input) == 0 {
		return nil
	}

	hardwareProfile := &input[0]

	result := &virtualmachineinstances.HardwareProfileUpdate{
		CpuCount:             pointer.To(int64(hardwareProfile.CpuCount)),
		DynamicMemoryEnabled: pointer.To(virtualmachineinstances.DynamicMemoryEnabled(strconv.FormatBool(hardwareProfile.DynamicMemoryEnabled))),
		DynamicMemoryMaxMB:   pointer.To(int64(hardwareProfile.DynamicMemoryMaxInMb)),
		DynamicMemoryMinMB:   pointer.To(int64(hardwareProfile.DynamicMemoryMinInMb)),
		LimitCPUForMigration: pointer.To(virtualmachineinstances.LimitCPUForMigration(strconv.FormatBool(hardwareProfile.LimitCpuForMigrationEnabled))),
		MemoryMB:             pointer.To(int64(hardwareProfile.MemoryInMb)),
	}

	return result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfile(input *virtualmachineinstances.HardwareProfile) []HardwareProfile {
	result := make([]HardwareProfile, 0)
	if input == nil {
		return result
	}

	hardwareProfile := HardwareProfile{}

	if v := input.CpuCount; v != nil {
		hardwareProfile.CpuCount = int(*v)
	}

	if v := input.DynamicMemoryEnabled; v != nil {
		hardwareProfile.DynamicMemoryEnabled = *v == virtualmachineinstances.DynamicMemoryEnabledTrue
	}

	if v := input.DynamicMemoryMaxMB; v != nil {
		hardwareProfile.DynamicMemoryMaxInMb = int(*v)
	}

	if v := input.DynamicMemoryMinMB; v != nil {
		hardwareProfile.DynamicMemoryMinInMb = int(*v)
	}

	if v := input.LimitCPUForMigration; v != nil {
		hardwareProfile.LimitCpuForMigrationEnabled = *v == virtualmachineinstances.LimitCPUForMigrationTrue
	}

	if v := input.MemoryMB; v != nil {
		hardwareProfile.MemoryInMb = int(*v)
	}

	return append(result, hardwareProfile)
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

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfaces(input *[]virtualmachineinstances.NetworkInterface) []NetworkInterface {
	result := make([]NetworkInterface, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		networkInterface := NetworkInterface{
			id:               pointer.From(item.NicId),
			name:             pointer.From(item.Name),
			VirtualNetworkId: pointer.From(item.VirtualNetworkId),
			MacAddress:       pointer.From(item.MacAddress),
		}

		if v := item.IPv4AddressType; v != nil {
			networkInterface.Ipv4AddressType = string(*v)
		}

		if v := item.IPv6AddressType; v != nil {
			networkInterface.Ipv6AddressType = string(*v)
		}

		if v := item.MacAddressType; v != nil {
			networkInterface.MacAddressType = string(*v)
		}

		result = append(result, networkInterface)
	}

	return result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(input []OSProfile) *virtualmachineinstances.OsProfileForVMInstance {
	if len(input) == 0 {
		return nil
	}

	osProfile := &input[0]

	result := &virtualmachineinstances.OsProfileForVMInstance{
		ComputerName:  pointer.To(osProfile.ComputerName),
		AdminPassword: pointer.To(osProfile.AdminPassword),
	}

	return result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(input *virtualmachineinstances.OsProfileForVMInstance) []OSProfile {
	result := make([]OSProfile, 0)
	if input == nil {
		return result
	}

	osProfile := OSProfile{
		ComputerName:  pointer.From(input.ComputerName),
		AdminPassword: pointer.From(input.AdminPassword),
	}

	return append(result, osProfile)
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(input []string) *[]virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined {
	result := make([]virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		availabilitySet := virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined{
			Id: pointer.To(v),
		}
		result = append(result, availabilitySet)
	}

	return &result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(input *[]virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		if avsId := v.Id; avsId != nil {
			result = append(result, *avsId)
		}
	}

	return result
}
