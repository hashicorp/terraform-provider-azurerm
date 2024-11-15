package azurestackhci

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/marketplacegalleryimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/storagecontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/virtualharddisks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/virtualmachineinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = StackHCIWindowsVirtualMachineResource{}
	_ sdk.ResourceWithUpdate = StackHCIWindowsVirtualMachineResource{}
)

type StackHCIWindowsVirtualMachineResource struct{}

func (StackHCIWindowsVirtualMachineResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.StackHCIVirtualMachineID
}

func (StackHCIWindowsVirtualMachineResource) ResourceType() string {
	return "azurerm_stack_hci_windows_virtual_machine"
}

func (StackHCIWindowsVirtualMachineResource) ModelObject() interface{} {
	return &StackHCIWindowsVirtualMachineResourceModel{}
}

type StackHCIWindowsVirtualMachineResourceModel struct {
	ArcMachineId           string                                         `tfschema:"arc_machine_id"`
	CustomLocationId       string                                         `tfschema:"custom_location_id"`
	HardwareProfile        []StackHCIVirtualMachineHardwareProfile        `tfschema:"hardware_profile"`
	HttpProxyConfiguration []StackHCIVirtualMachineHttpProxyConfiguration `tfschema:"http_proxy_configuration"`
	Identity               []identity.ModelSystemAssigned                 `tfschema:"identity"`
	NetworkProfile         []StackHCIVirtualMachineNetworkProfile         `tfschema:"network_profile"`
	OsProfile              []StackHCIVirtualMachineOsProfileWindows       `tfschema:"os_profile"`
	SecurityProfile        []StackHCIVirtualMachineSecurityProfile        `tfschema:"security_profile"`
	StorageProfile         []StackHCIVirtualMachineStorageProfile         `tfschema:"storage_profile"`
}

type StackHCIVirtualMachineHardwareProfile struct {
	DynamicMemory   []StackHCIVirtualMachineDynamicMemory `tfschema:"dynamic_memory"`
	MemoryMb        int64                                 `tfschema:"memory_mb"`
	ProcessorNumber int64                                 `tfschema:"processor_number"`
	VmSize          string                                `tfschema:"vm_size"`
}

type StackHCIVirtualMachineDynamicMemory struct {
	MaximumMemoryMb    int64 `tfschema:"maximum_memory_mb"`
	MinimumMemoryMb    int64 `tfschema:"minimum_memory_mb"`
	TargetMemoryBuffer int64 `tfschema:"target_memory_buffer"`
}

type StackHCIVirtualMachineHttpProxyConfiguration struct {
	HttpProxy  string   `tfschema:"http_proxy"`
	HttpsProxy string   `tfschema:"https_proxy"`
	NoProxy    []string `tfschema:"no_proxy"`
	TrustedCa  string   `tfschema:"trusted_ca"`
}

type StackHCIVirtualMachineNetworkProfile struct {
	NetworkInterfaceIds []string `tfschema:"network_interface_ids"`
}

type StackHCIVirtualMachineOsProfileWindows struct {
	AdminUsername string `tfschema:"admin_username"`
	AdminPassword string `tfschema:"admin_password"`
	ComputerName  string `tfschema:"computer_name"`

	// windowsConfiguration
	AutomaticUpdateEnabled        bool                                 `tfschema:"automatic_update_enabled"`
	ProvisionVmAgentEnabled       bool                                 `tfschema:"provision_vm_agent_enabled"`
	ProvisionVmConfigAgentEnabled bool                                 `tfschema:"provision_vm_config_agent_enabled"`
	SshPublicKey                  []StackHCIVirtualMachineSshPublicKey `tfschema:"ssh_public_key"`
	TimeZone                      string                               `tfschema:"time_zone"`
}

type StackHCIVirtualMachineSshPublicKey struct {
	KeyData string `tfschema:"key_data"`
	Path    string `tfschema:"path"`
}

type StackHCIVirtualMachineSecurityProfile struct {
	SecureBootEnabled bool   `tfschema:"secure_boot_enabled"`
	SecurityType      string `tfschema:"security_type"`
	TpmEnabled        bool   `tfschema:"tpm_enabled"`
}

type StackHCIVirtualMachineStorageProfile struct {
	DataDiskIds           []string `tfschema:"data_disk_ids"`
	ImageId               string   `tfschema:"image_id"`
	OsDiskId              string   `tfschema:"os_disk_id"`
	VmConfigStoragePathId string   `tfschema:"vm_config_storage_path_id"`
}

type StackHCIVirtualMachineOsDisk struct{}

func (StackHCIWindowsVirtualMachineResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"arc_machine_id": commonschema.ResourceIDReferenceRequiredForceNew(&machines.MachineId{}),

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}),

		"hardware_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"vm_size": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(virtualmachineinstances.PossibleValuesForVMSizeEnum(), false),
					},

					"processor_number": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"memory_mb": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"dynamic_memory": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"maximum_memory_mb": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntAtLeast(1),
								},

								"minimum_memory_mb": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntAtLeast(1),
								},

								"target_memory_buffer": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntAtLeast(1),
								},
							},
						},
					},
				},
			},
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_interface_ids": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: networkinterfaces.ValidateNetworkInterfaceID,
						},
					},
				},
			},
		},

		"os_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"admin_password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"computer_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[\-a-zA-Z0-9]{0,15}$`),
							"name must begin and end with an alphanumeric character, be between 2 and 64 characters in length and can only contain alphanumeric characters, hyphens, periods or underscores.",
						),
					},

					"automatic_update_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"ssh_public_key": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"path": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"key_data": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"time_zone": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"provision_vm_agent_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"provision_vm_config_agent_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"security_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tpm_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"secure_boot_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  true,
					},

					"security_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(virtualmachineinstances.PossibleValuesForSecurityTypes(), false),
					},
				},
			},
		},

		"storage_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_disk_ids": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: virtualharddisks.ValidateVirtualHardDiskID,
						},
					},

					"image_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: marketplacegalleryimages.ValidateMarketplaceGalleryImageID,
					},

					"os_disk_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: virtualharddisks.ValidateVirtualHardDiskID,
					},

					"vm_config_storage_path_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: storagecontainers.ValidateStorageContainerID,
					},
				},
			},
		},

		"http_proxy_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"http_proxy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"https_proxy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"no_proxy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"trusted_ca": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedIdentityOptional(),
	}
}

func (StackHCIWindowsVirtualMachineResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StackHCIWindowsVirtualMachineResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualMachineInstances

			var config StackHCIWindowsVirtualMachineResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			arcMachineId, err := machines.ParseMachineID(config.ArcMachineId)
			if err != nil {
				return fmt.Errorf("parsing `arc_machine_id`: %+v", err)
			}

			id := parse.NewStackHCIVirtualMachineID(arcMachineId.SubscriptionId, arcMachineId.ResourceGroupName, arcMachineId.MachineName, "default")
			scopeId := commonids.NewScopeID(config.ArcMachineId)

			existing, err := client.Get(ctx, scopeId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := virtualmachineinstances.VirtualMachineInstance{
				ExtendedLocation: &virtualmachineinstances.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(virtualmachineinstances.ExtendedLocationTypesCustomLocation),
				},
				Properties: &virtualmachineinstances.VirtualMachineInstanceProperties{
					HardwareProfile: expandVirtualMachineHardwareProfile(config.HardwareProfile),
					HTTPProxyConfig: expandVirtualMachineHttpProxyConfig(config.HttpProxyConfiguration),
					NetworkProfile:  expandVirtualMachineNetworkProfile(config.NetworkProfile),
					OsProfile:       expandVirtualMachineOsProfileWindows(config.OsProfile),
					SecurityProfile: expandVirtualMachineSecurityProfile(config.SecurityProfile),
					StorageProfile:  expandVirtualMachineStorageProfileWindows(config.StorageProfile),
				},
			}
			if len(config.Identity) > 0 {
				payload.Identity, err = identity.ExpandSystemAssignedFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, scopeId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			time.Sleep(5 * time.Second)

			return nil
		},
	}
}

func (r StackHCIWindowsVirtualMachineResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualMachineInstances

			id, err := parse.StackHCIVirtualMachineID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroup, id.MachineName)
			scopeId := commonids.NewScopeID(arcMachineId.ID())

			resp, err := client.Get(ctx, scopeId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema := StackHCIWindowsVirtualMachineResourceModel{
				ArcMachineId: arcMachineId.ID(),
			}

			if model := resp.Model; model != nil {
				schema.Identity = identity.FlattenSystemAssignedToModel(model.Identity)

				if model.ExtendedLocation != nil && model.ExtendedLocation.Name != nil {
					customLocationId, err := customlocations.ParseCustomLocationIDInsensitively(*model.ExtendedLocation.Name)
					if err != nil {
						return err
					}

					schema.CustomLocationId = customLocationId.ID()
				}

				if props := model.Properties; props != nil {
					schema.HardwareProfile = flattenVirtualMachineHardwareProfile(props.HardwareProfile)
					schema.HttpProxyConfiguration = flattenVirtualMachineHttpProxyConfig(props.HTTPProxyConfig)
					schema.NetworkProfile = flattenVirtualMachineNetworkProfile(props.NetworkProfile)
					schema.OsProfile = flattenVirtualMachineOsProfileWindows(props.OsProfile)
					schema.SecurityProfile = flattenVirtualMachineSecurityProfile(props.SecurityProfile)
					schema.StorageProfile = flattenVirtualMachineStorageProfileWindows(props.StorageProfile)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r StackHCIWindowsVirtualMachineResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualMachineInstances

			id, err := parse.StackHCIVirtualMachineID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroup, id.MachineName)
			scopeId := commonids.NewScopeID(arcMachineId.ID())

			var config StackHCIWindowsVirtualMachineResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, scopeId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil || resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			existing := resp.Model

			if metadata.ResourceData.HasChange("identity") {
				if len(config.Identity) != 0 {
					expandedIdentity, err := identity.ExpandSystemAssignedFromModel(config.Identity)
					if err != nil {
						return fmt.Errorf("expanding `identity`: %+v", err)
					}

					existing.Identity = expandedIdentity
				} else {
					existing.Identity = nil
				}
			}

			if metadata.ResourceData.HasChange("hardware_profile") {
				existing.Properties.HardwareProfile = expandVirtualMachineHardwareProfile(config.HardwareProfile)
			}

			if metadata.ResourceData.HasChange("network_profile") {
				existing.Properties.NetworkProfile = expandVirtualMachineNetworkProfile(config.NetworkProfile)
			}

			if metadata.ResourceData.HasChange("storage_profile") {
				existing.Properties.StorageProfile = expandVirtualMachineStorageProfileWindows(config.StorageProfile)
			}

			if metadata.ResourceData.HasChange("os_profile") {
				existing.Properties.OsProfile = expandVirtualMachineOsProfileWindows(config.OsProfile)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, scopeId, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r StackHCIWindowsVirtualMachineResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualMachineInstances

			id, err := parse.StackHCIVirtualMachineID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroup, id.MachineName)
			scopeId := commonids.NewScopeID(arcMachineId.ID())

			if err := client.DeleteThenPoll(ctx, scopeId); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandVirtualMachineHardwareProfile(input []StackHCIVirtualMachineHardwareProfile) *virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	output := &virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfile{
		DynamicMemoryConfig: expandVirtualMachineDynamicMemory(v.DynamicMemory),
		MemoryMB:            pointer.To(v.MemoryMb),
		Processors:          pointer.To(v.ProcessorNumber),
		VMSize:              pointer.To(virtualmachineinstances.VMSizeEnum(v.VmSize)),
	}

	return output
}

func flattenVirtualMachineHardwareProfile(input *virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfile) []StackHCIVirtualMachineHardwareProfile {
	if input == nil {
		return make([]StackHCIVirtualMachineHardwareProfile, 0)
	}

	return []StackHCIVirtualMachineHardwareProfile{
		{
			DynamicMemory:   nil,
			MemoryMb:        pointer.From(input.MemoryMB),
			ProcessorNumber: pointer.From(input.Processors),
			VmSize:          string(pointer.From(input.VMSize)),
		},
	}
}

func expandVirtualMachineDynamicMemory(input []StackHCIVirtualMachineDynamicMemory) *virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfileDynamicMemoryConfig {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	output := &virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfileDynamicMemoryConfig{
		MaximumMemoryMB: pointer.To(v.MaximumMemoryMb),
		MinimumMemoryMB: pointer.To(v.MinimumMemoryMb),
	}

	return output
}

func flattenVirtualMachineDynamicMemory(input *virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfileDynamicMemoryConfig) []StackHCIVirtualMachineDynamicMemory {
	if input == nil {
		return make([]StackHCIVirtualMachineDynamicMemory, 0)
	}

	return []StackHCIVirtualMachineDynamicMemory{
		{
			MaximumMemoryMb: pointer.From(input.MaximumMemoryMB),
			MinimumMemoryMb: pointer.From(input.MinimumMemoryMB),
		},
	}
}

func expandVirtualMachineHttpProxyConfig(input []StackHCIVirtualMachineHttpProxyConfiguration) *virtualmachineinstances.HTTPProxyConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	output := &virtualmachineinstances.HTTPProxyConfiguration{
		HTTPProxy:  pointer.To(v.HttpProxy),
		HTTPSProxy: pointer.To(v.HttpsProxy),
		NoProxy:    pointer.To(v.NoProxy),
		TrustedCa:  pointer.To(v.TrustedCa),
	}

	return output
}

func flattenVirtualMachineHttpProxyConfig(input *virtualmachineinstances.HTTPProxyConfiguration) []StackHCIVirtualMachineHttpProxyConfiguration {
	if input == nil {
		return make([]StackHCIVirtualMachineHttpProxyConfiguration, 0)
	}

	return []StackHCIVirtualMachineHttpProxyConfiguration{
		{
			HttpProxy:  pointer.From(input.HTTPProxy),
			HttpsProxy: pointer.From(input.HTTPSProxy),
			NoProxy:    pointer.From(input.NoProxy),
			TrustedCa:  pointer.From(input.TrustedCa),
		},
	}
}

func expandVirtualMachineNetworkProfile(input []StackHCIVirtualMachineNetworkProfile) *virtualmachineinstances.VirtualMachineInstancePropertiesNetworkProfile {
	if len(input) == 0 {
		return nil
	}

	networkInterfaces := make([]virtualmachineinstances.VirtualMachineInstancePropertiesNetworkProfileNetworkInterfacesInlined, 0)
	for _, networkInterfaceId := range input[0].NetworkInterfaceIds {
		networkInterfaces = append(networkInterfaces, virtualmachineinstances.VirtualMachineInstancePropertiesNetworkProfileNetworkInterfacesInlined{
			Id: pointer.To(networkInterfaceId),
		})
	}

	output := &virtualmachineinstances.VirtualMachineInstancePropertiesNetworkProfile{
		NetworkInterfaces: &networkInterfaces,
	}

	return output
}

func flattenVirtualMachineNetworkProfile(input *virtualmachineinstances.VirtualMachineInstancePropertiesNetworkProfile) []StackHCIVirtualMachineNetworkProfile {
	if input == nil || input.NetworkInterfaces == nil {
		return make([]StackHCIVirtualMachineNetworkProfile, 0)
	}

	networkInterfaceIds := make([]string, 0)
	for _, networkInterface := range *input.NetworkInterfaces {
		if networkInterface.Id != nil {
			networkInterfaceIds = append(networkInterfaceIds, *networkInterface.Id)
		}
	}

	return []StackHCIVirtualMachineNetworkProfile{
		{
			NetworkInterfaceIds: networkInterfaceIds,
		},
	}
}

func expandVirtualMachineOsProfileWindows(input []StackHCIVirtualMachineOsProfileWindows) *virtualmachineinstances.VirtualMachineInstancePropertiesOsProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	output := &virtualmachineinstances.VirtualMachineInstancePropertiesOsProfile{
		AdminUsername: pointer.To(v.AdminUsername),
		AdminPassword: pointer.To(v.AdminPassword),
		ComputerName:  pointer.To(v.ComputerName),
		WindowsConfiguration: &virtualmachineinstances.VirtualMachineInstancePropertiesOsProfileWindowsConfiguration{
			EnableAutomaticUpdates: pointer.To(v.AutomaticUpdateEnabled),
			ProvisionVMAgent:       pointer.To(v.ProvisionVmAgentEnabled),
			ProvisionVMConfigAgent: pointer.To(v.ProvisionVmConfigAgentEnabled),
			Ssh:                    expandVirtualMachineOsProfileSsh(v.SshPublicKey),
			TimeZone:               pointer.To(v.TimeZone),
		},
	}

	return output
}

func flattenVirtualMachineOsProfileWindows(input *virtualmachineinstances.VirtualMachineInstancePropertiesOsProfile) []StackHCIVirtualMachineOsProfileWindows {
	if input == nil {
		return make([]StackHCIVirtualMachineOsProfileWindows, 0)
	}

	result := StackHCIVirtualMachineOsProfileWindows{
		AdminUsername: pointer.From(input.AdminUsername),
		AdminPassword: pointer.From(input.AdminPassword),
		ComputerName:  pointer.From(input.ComputerName),
	}

	if input.WindowsConfiguration != nil {
		result.AutomaticUpdateEnabled = pointer.From(input.WindowsConfiguration.EnableAutomaticUpdates)
		result.ProvisionVmAgentEnabled = pointer.From(input.WindowsConfiguration.ProvisionVMAgent)
		result.ProvisionVmConfigAgentEnabled = pointer.From(input.WindowsConfiguration.ProvisionVMConfigAgent)
		result.SshPublicKey = flattenVirtualMachineOsProfileSsh(input.WindowsConfiguration.Ssh)
		result.TimeZone = pointer.From(input.WindowsConfiguration.TimeZone)
	}

	return []StackHCIVirtualMachineOsProfileWindows{
		result,
	}
}

func expandVirtualMachineOsProfileSsh(input []StackHCIVirtualMachineSshPublicKey) *virtualmachineinstances.SshConfiguration {
	if len(input) == 0 {
		return nil
	}

	sshPublicKeys := make([]virtualmachineinstances.SshPublicKey, 0)
	for _, key := range input {
		sshPublicKeys = append(sshPublicKeys, virtualmachineinstances.SshPublicKey{
			KeyData: pointer.To(key.KeyData),
			Path:    pointer.To(key.Path),
		})
	}

	return &virtualmachineinstances.SshConfiguration{
		PublicKeys: &sshPublicKeys,
	}
}

func flattenVirtualMachineOsProfileSsh(input *virtualmachineinstances.SshConfiguration) []StackHCIVirtualMachineSshPublicKey {
	if input == nil || input.PublicKeys == nil {
		return make([]StackHCIVirtualMachineSshPublicKey, 0)
	}

	output := make([]StackHCIVirtualMachineSshPublicKey, 0)
	for _, key := range *input.PublicKeys {
		output = append(output, StackHCIVirtualMachineSshPublicKey{
			KeyData: pointer.From(key.KeyData),
			Path:    pointer.From(key.Path),
		})
	}

	return output
}

func expandVirtualMachineSecurityProfile(input []StackHCIVirtualMachineSecurityProfile) *virtualmachineinstances.VirtualMachineInstancePropertiesSecurityProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	output := &virtualmachineinstances.VirtualMachineInstancePropertiesSecurityProfile{
		EnableTPM:    pointer.To(v.TpmEnabled),
		SecurityType: pointer.To(virtualmachineinstances.SecurityTypes(v.SecurityType)),
		UefiSettings: &virtualmachineinstances.VirtualMachineInstancePropertiesSecurityProfileUefiSettings{
			SecureBootEnabled: pointer.To(v.SecureBootEnabled),
		},
	}

	return output
}

func flattenVirtualMachineSecurityProfile(input *virtualmachineinstances.VirtualMachineInstancePropertiesSecurityProfile) []StackHCIVirtualMachineSecurityProfile {
	if input == nil {
		return make([]StackHCIVirtualMachineSecurityProfile, 0)
	}

	secureBootEnabled := false
	if input.UefiSettings != nil {
		secureBootEnabled = pointer.From(input.UefiSettings.SecureBootEnabled)
	}

	return []StackHCIVirtualMachineSecurityProfile{
		{
			TpmEnabled:        pointer.From(input.EnableTPM),
			SecurityType:      string(pointer.From(input.SecurityType)),
			SecureBootEnabled: secureBootEnabled,
		},
	}
}

func expandVirtualMachineStorageProfileWindows(input []StackHCIVirtualMachineStorageProfile) *virtualmachineinstances.VirtualMachineInstancePropertiesStorageProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	dataDiskIds := make([]virtualmachineinstances.VirtualMachineInstancePropertiesStorageProfileDataDisksInlined, 0)
	for _, dataDiskId := range v.DataDiskIds {
		dataDiskIds = append(dataDiskIds, virtualmachineinstances.VirtualMachineInstancePropertiesStorageProfileDataDisksInlined{
			Id: pointer.To(dataDiskId),
		})
	}

	output := &virtualmachineinstances.VirtualMachineInstancePropertiesStorageProfile{
		DataDisks: pointer.To(dataDiskIds),
		OsDisk: &virtualmachineinstances.VirtualMachineInstancePropertiesStorageProfileOsDisk{
			OsType: pointer.To(virtualmachineinstances.OperatingSystemTypesWindows),
		},
		ImageReference: &virtualmachineinstances.VirtualMachineInstancePropertiesStorageProfileImageReference{
			Id: pointer.To(v.ImageId),
		},
	}

	if v.OsDiskId != "" {
		output.OsDisk.Id = pointer.To(v.OsDiskId)
	}

	if v.VmConfigStoragePathId != "" {
		output.VMConfigStoragePathId = pointer.To(v.VmConfigStoragePathId)
	}

	return output
}

func flattenVirtualMachineStorageProfileWindows(input *virtualmachineinstances.VirtualMachineInstancePropertiesStorageProfile) []StackHCIVirtualMachineStorageProfile {
	if input == nil {
		return make([]StackHCIVirtualMachineStorageProfile, 0)
	}

	dataDiskIds := make([]string, 0)
	if input.DataDisks != nil {
		for _, dataDisk := range *input.DataDisks {
			if dataDisk.Id != nil {
				dataDiskIds = append(dataDiskIds, *dataDisk.Id)
			}
		}
	}

	var imageId string
	if input.ImageReference != nil {
		imageId = pointer.From(input.ImageReference.Id)
	}

	var osDiskId string
	if input.OsDisk != nil {
		osDiskId = pointer.From(input.OsDisk.Id)
	}

	return []StackHCIVirtualMachineStorageProfile{
		{
			DataDiskIds:           dataDiskIds,
			ImageId:               imageId,
			OsDiskId:              osDiskId,
			VmConfigStoragePathId: pointer.From(input.VMConfigStoragePathId),
		},
	}
}
