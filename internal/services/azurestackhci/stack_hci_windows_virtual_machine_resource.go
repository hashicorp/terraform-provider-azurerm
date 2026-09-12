package azurestackhci

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/galleryimages"
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
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
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
	NetworkProfile         []StackHCIVirtualMachineNetworkProfile         `tfschema:"network_profile"`
	OsProfile              []StackHCIVirtualMachineOsProfileWindows       `tfschema:"os_profile"`
	StorageProfile         []StackHCIVirtualMachineStorageProfile         `tfschema:"storage_profile"`

	// securityProfile is flattened, since it is always returned from API
	SecureBootEnabled bool   `tfschema:"secure_boot_enabled"`
	SecurityType      string `tfschema:"security_type"`
	TpmEnabled        bool   `tfschema:"tpm_enabled"`
}

type StackHCIVirtualMachineHardwareProfile struct {
	DynamicMemory   []StackHCIVirtualMachineDynamicMemory `tfschema:"dynamic_memory"`
	MemoryInMb      int64                                 `tfschema:"memory_in_mb"`
	ProcessorNumber int64                                 `tfschema:"processor_number"`
	VmSize          string                                `tfschema:"vm_size"`
}

type StackHCIVirtualMachineDynamicMemory struct {
	MaximumMemoryInMb            int64 `tfschema:"maximum_memory_in_mb"`
	MinimumMemoryInMb            int64 `tfschema:"minimum_memory_in_mb"`
	TargetMemoryBufferPercentage int64 `tfschema:"target_memory_buffer_percentage"`
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

type StackHCIVirtualMachineStorageProfile struct {
	DataDiskIds           []string `tfschema:"data_disk_ids"`
	ImageId               string   `tfschema:"image_id"`
	VmConfigStoragePathId string   `tfschema:"vm_config_storage_path_id"`
}

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

					"memory_in_mb": {
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
								"maximum_memory_in_mb": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntAtLeast(1),
								},

								"minimum_memory_in_mb": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntAtLeast(1),
								},

								"target_memory_buffer_percentage": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntBetween(5, 2000),
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
						Elem:     commonschema.ResourceIDReferenceElem(&networkinterfaces.NetworkInterfaceId{}),
					},
				},
			},
		},

		"os_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: computeValidate.WindowsAdminUsername,
					},

					"admin_password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Sensitive:    true,
						ValidateFunc: computeValidate.WindowsAdminPassword,
					},

					"computer_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: computeValidate.WindowsComputerNameFull,
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
						ForceNew: true,
						Default:  false,
					},

					"provision_vm_config_agent_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},
				},
			},
		},

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
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualmachineinstances.PossibleValuesForSecurityTypes(), false),
		},

		"storage_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_disk_ids": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem:     commonschema.ResourceIDReferenceElem(&virtualharddisks.VirtualHardDiskId{}),
					},

					"image_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.Any(
							marketplacegalleryimages.ValidateMarketplaceGalleryImageID,
							galleryimages.ValidateGalleryImageID,
						),
					},

					"os_disk_id": commonschema.ResourceIDReferenceOptionalForceNew(&virtualharddisks.VirtualHardDiskId{}),

					"vm_config_storage_path_id": commonschema.ResourceIDReferenceOptionalForceNew(&storagecontainers.StorageContainerId{}),
				},
			},
		},

		"http_proxy_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"http_proxy": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"https_proxy": {
						Type:         pluginsdk.TypeString,
						Required:     true,
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
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
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
					SecurityProfile: expandVirtualMachineSecurityProfile(config),
					StorageProfile:  expandVirtualMachineStorageProfileWindows(config.StorageProfile),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, scopeId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// https://github.com/Azure/azure-rest-api-specs/issues/31876
			if err := resourceVirtualMachineWaitForCreated(ctx, *client, id); err != nil {
				return fmt.Errorf("waiting for %s to be created: %+v", id, err)
			}

			metadata.SetID(id)

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

			var config StackHCIWindowsVirtualMachineResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
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
				if model.ExtendedLocation != nil && model.ExtendedLocation.Name != nil {
					customLocationId, err := customlocations.ParseCustomLocationIDInsensitively(*model.ExtendedLocation.Name)
					if err != nil {
						return err
					}

					schema.CustomLocationId = customLocationId.ID()
				}

				if props := model.Properties; props != nil {
					schema.HardwareProfile = flattenVirtualMachineHardwareProfile(props.HardwareProfile)
					schema.HttpProxyConfiguration = flattenVirtualMachineHttpProxyConfig(props.HTTPProxyConfig, config.HttpProxyConfiguration)
					schema.NetworkProfile = flattenVirtualMachineNetworkProfile(props.NetworkProfile)
					schema.OsProfile = flattenVirtualMachineOsProfileWindows(props.OsProfile, config.OsProfile)
					schema.StorageProfile = flattenVirtualMachineStorageProfileWindows(props.StorageProfile)

					if securityProfile := props.SecurityProfile; securityProfile != nil {
						schema.TpmEnabled = pointer.From(securityProfile.EnableTPM)
						schema.SecurityType = string(pointer.From(securityProfile.SecurityType))

						secureBootEnabled := false
						if securityProfile.UefiSettings != nil {
							secureBootEnabled = pointer.From(securityProfile.UefiSettings.SecureBootEnabled)
						}
						schema.SecureBootEnabled = secureBootEnabled
					}
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

			payload := virtualmachineinstances.VirtualMachineInstanceUpdateRequest{
				Properties: &virtualmachineinstances.VirtualMachineInstanceUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("network_profile") {
				payload.Properties.NetworkProfile = expandVirtualMachineNetworkProfileForUpdate(config.NetworkProfile)
			}

			if metadata.ResourceData.HasChange("storage_profile") {
				payload.Properties.StorageProfile = expandVirtualMachineStorageProfileWindowsForUpdate(config.StorageProfile)
			}

			if err := client.UpdateThenPoll(ctx, scopeId, payload); err != nil {
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
		MemoryMB:            pointer.To(v.MemoryInMb),
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
			DynamicMemory:   flattenVirtualMachineDynamicMemory(input.DynamicMemoryConfig),
			MemoryInMb:      pointer.From(input.MemoryMB),
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
		MaximumMemoryMB:    pointer.To(v.MaximumMemoryInMb),
		MinimumMemoryMB:    pointer.To(v.MinimumMemoryInMb),
		TargetMemoryBuffer: pointer.To(v.TargetMemoryBufferPercentage),
	}

	return output
}

func flattenVirtualMachineDynamicMemory(input *virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfileDynamicMemoryConfig) []StackHCIVirtualMachineDynamicMemory {
	if input == nil {
		return make([]StackHCIVirtualMachineDynamicMemory, 0)
	}

	// The API may return blocks with all values set to 0, in this case we ignore the block
	if pointer.From(input.MaximumMemoryMB) == 0 && pointer.From(input.MinimumMemoryMB) == 0 && pointer.From(input.TargetMemoryBuffer) == 0 {
		return make([]StackHCIVirtualMachineDynamicMemory, 0)
	}

	return []StackHCIVirtualMachineDynamicMemory{
		{
			MaximumMemoryInMb:            pointer.From(input.MaximumMemoryMB),
			MinimumMemoryInMb:            pointer.From(input.MinimumMemoryMB),
			TargetMemoryBufferPercentage: pointer.From(input.TargetMemoryBuffer),
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

func flattenVirtualMachineHttpProxyConfig(input *virtualmachineinstances.HTTPProxyConfiguration, configuredHttpProxyConfig []StackHCIVirtualMachineHttpProxyConfiguration) []StackHCIVirtualMachineHttpProxyConfiguration {
	if input == nil {
		return make([]StackHCIVirtualMachineHttpProxyConfiguration, 0)
	}

	output := StackHCIVirtualMachineHttpProxyConfiguration{
		NoProxy:   pointer.From(input.NoProxy),
		TrustedCa: pointer.From(input.TrustedCa),
	}

	// httpProxy and httpsProxy are not returned from the server
	if len(configuredHttpProxyConfig) > 0 {
		output.HttpProxy = configuredHttpProxyConfig[0].HttpProxy
		output.HttpsProxy = configuredHttpProxyConfig[0].HttpsProxy
	}

	return []StackHCIVirtualMachineHttpProxyConfiguration{
		output,
	}
}

func expandVirtualMachineNetworkProfileForUpdate(input []StackHCIVirtualMachineNetworkProfile) *virtualmachineinstances.NetworkProfileUpdate {
	if len(input) == 0 {
		return &virtualmachineinstances.NetworkProfileUpdate{}
	}

	networkInterfaces := make([]virtualmachineinstances.NetworkProfileUpdateNetworkInterfacesInlined, 0)
	for _, networkInterfaceId := range input[0].NetworkInterfaceIds {
		networkInterfaces = append(networkInterfaces, virtualmachineinstances.NetworkProfileUpdateNetworkInterfacesInlined{
			Id: pointer.To(networkInterfaceId),
		})
	}

	output := &virtualmachineinstances.NetworkProfileUpdate{
		NetworkInterfaces: &networkInterfaces,
	}

	return output
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

func flattenVirtualMachineOsProfileWindows(input *virtualmachineinstances.VirtualMachineInstancePropertiesOsProfile, configuredOsProfile []StackHCIVirtualMachineOsProfileWindows) []StackHCIVirtualMachineOsProfileWindows {
	if input == nil {
		return make([]StackHCIVirtualMachineOsProfileWindows, 0)
	}

	output := StackHCIVirtualMachineOsProfileWindows{
		AdminUsername: pointer.From(input.AdminUsername),
		ComputerName:  pointer.From(input.ComputerName),
	}

	// adminPassword is not returned from the server, so it should be taken from the configured value
	if len(configuredOsProfile) > 0 {
		output.AdminPassword = configuredOsProfile[0].AdminPassword
	}

	if input.WindowsConfiguration != nil {
		output.AutomaticUpdateEnabled = pointer.From(input.WindowsConfiguration.EnableAutomaticUpdates)
		output.ProvisionVmAgentEnabled = pointer.From(input.WindowsConfiguration.ProvisionVMAgent)
		output.ProvisionVmConfigAgentEnabled = pointer.From(input.WindowsConfiguration.ProvisionVMConfigAgent)
		output.SshPublicKey = flattenVirtualMachineOsProfileSsh(input.WindowsConfiguration.Ssh)
		output.TimeZone = pointer.From(input.WindowsConfiguration.TimeZone)
	}

	return []StackHCIVirtualMachineOsProfileWindows{
		output,
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

func expandVirtualMachineSecurityProfile(input StackHCIWindowsVirtualMachineResourceModel) *virtualmachineinstances.VirtualMachineInstancePropertiesSecurityProfile {
	output := &virtualmachineinstances.VirtualMachineInstancePropertiesSecurityProfile{
		EnableTPM:    pointer.To(input.TpmEnabled),
		SecurityType: pointer.To(virtualmachineinstances.SecurityTypes(input.SecurityType)),
		UefiSettings: &virtualmachineinstances.VirtualMachineInstancePropertiesSecurityProfileUefiSettings{
			SecureBootEnabled: pointer.To(input.SecureBootEnabled),
		},
	}

	return output
}

func expandVirtualMachineStorageProfileWindowsForUpdate(input []StackHCIVirtualMachineStorageProfile) *virtualmachineinstances.StorageProfileUpdate {
	if len(input) == 0 {
		return &virtualmachineinstances.StorageProfileUpdate{}
	}

	v := input[0]

	dataDiskIds := make([]virtualmachineinstances.StorageProfileUpdateDataDisksInlined, 0)
	for _, dataDiskId := range v.DataDiskIds {
		dataDiskIds = append(dataDiskIds, virtualmachineinstances.StorageProfileUpdateDataDisksInlined{
			Id: pointer.To(dataDiskId),
		})
	}

	output := &virtualmachineinstances.StorageProfileUpdate{
		DataDisks: pointer.To(dataDiskIds),
	}

	return output
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

	return []StackHCIVirtualMachineStorageProfile{
		{
			DataDiskIds:           dataDiskIds,
			ImageId:               imageId,
			VmConfigStoragePathId: pointer.From(input.VMConfigStoragePathId),
		},
	}
}

func resourceVirtualMachineWaitForCreated(ctx context.Context, client virtualmachineinstances.VirtualMachineInstancesClient, id parse.StackHCIVirtualMachineId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal error: context had no deadline")
	}

	state := &pluginsdk.StateChangeConf{
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 3,
		Pending:                   []string{"NotFound"},
		Target:                    []string{"Found"},
		Refresh:                   resourceVirtualMachineRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be created: %+v", id, err)
	}

	return nil
}

func resourceVirtualMachineRefreshFunc(ctx context.Context, client virtualmachineinstances.VirtualMachineInstancesClient, id parse.StackHCIVirtualMachineId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking status for %s ..", id)

		arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroup, id.MachineName)
		scopeId := commonids.NewScopeID(arcMachineId.ID())

		resp, err := client.Get(ctx, scopeId)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return resp, "Found", nil
	}
}
