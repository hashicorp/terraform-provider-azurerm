package azurestackhci

import (
	"context"
	"fmt"
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = StackHCIVirtualMachineResource{}
	_ sdk.ResourceWithUpdate = StackHCIVirtualMachineResource{}
)

type StackHCIVirtualMachineResource struct{}

func (StackHCIVirtualMachineResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.StackHCIVirtualMachineID
}

func (StackHCIVirtualMachineResource) ResourceType() string {
	return "azurerm_stack_hci_virtual_machine"
}

func (StackHCIVirtualMachineResource) ModelObject() interface{} {
	return &StackHCIVirtualMachineResourceModel{}
}

type StackHCIVirtualMachineResourceModel struct {
	ArcMachineId           string                                         `tfschema:"arc_machine_id"`
	CustomLocationId       string                                         `tfschema:"custom_location_id"`
	HardwareProfile        []StackHCIVirtualMachineHardwareProfile        `tfschema:"hardware_profile"`
	HttpProxyConfiguration []StackHCIVirtualMachineHttpProxyConfiguration `tfschema:"http_proxy_configuration"`
	Identity               []identity.ModelSystemAssigned                 `tfschema:"identity"`
	NetworkProfile         []StackHCIVirtualMachineNetworkProfile         `tfschema:"network_profile"`
	OsProfile              []StackHCIVirtualMachineOsProfile              `tfschema:"os_profile"`
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

type StackHCIVirtualMachineNetworkProfile struct {
	NetworkInterfaceIds []string `tfschema:"network_interface_ids"`
}

type StackHCIVirtualMachineOsProfile struct {
	AdminUsername        string                                       `tfschema:"admin_username"`
	AdminPassword        string                                       `tfschema:"admin_password"`
	ComputerName         string                                       `tfschema:"computer_name"`
	LinuxConfiguration   []StackHCIVirtualMachineLinuxConfiguration   `tfschema:"linux_configuration"`
	WindowsConfiguration []StackHCIVirtualMachineWindowsConfiguration `tfschema:"windows_configuration"`
}

type StackHCIVirtualMachineLinuxConfiguration struct {
	PasswordAuthenticationEnabled bool                                 `tfschema:"password_authentication_enabled"`
	ProvisionVmAgentEnabled       bool                                 `tfschema:"provision_vm_agent_enabled"`
	ProvisionVmConfigAgentEnabled bool                                 `tfschema:"provision_vm_config_agent_enabled"`
	SshPublicKey                  []StackHCIVirtualMachineSshPublicKey `tfschema:"ssh_public_key"`
}

type StackHCIVirtualMachineWindowsConfiguration struct {
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
	DataDiskIds           []string                       `tfschema:"data_disk_ids"`
	ImageId               string                         `tfschema:"image_id"`
	OsDisk                []StackHCIVirtualMachineOsDisk `tfschema:"os_disk"`
	VmConfigStoragePathId string                         `tfschema:"vm_config_storage_path_id"`
}

type StackHCIVirtualMachineOsDisk struct {
	DiskId string `tfschema:"disk_id"`
	OsType string `tfschema:"os_type"`
}

type StackHCIVirtualMachineHttpProxyConfiguration struct {
	HttpProxy  string   `tfschema:"http_proxy"`
	HttpsProxy string   `tfschema:"https_proxy"`
	NoProxy    []string `tfschema:"no_proxy"`
	TrustedCa  string   `tfschema:"trusted_ca"`
}

func (StackHCIVirtualMachineResource) Arguments() map[string]*pluginsdk.Schema {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"linux_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"password_authentication_enabled": {
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

					"windows_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
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
									Required:     true,
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
						Default:  false,
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

					"os_disk": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"disk_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: virtualharddisks.ValidateVirtualHardDiskID,
								},

								"os_type": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(virtualmachineinstances.PossibleValuesForOperatingSystemTypes(), false),
								},
							},
						},
					},

					"vm_config_storage_path_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
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
				Schema: map[string]*pluginsdk.Schema{},
			},
		},

		"identity": commonschema.SystemAssignedIdentityOptional(),
	}
}

func (StackHCIVirtualMachineResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StackHCIVirtualMachineResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualMachineInstances

			var config StackHCIVirtualMachineResourceModel
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
					HardwareProfile: expandVirtualMachineInstanceHardwareProfile(config.HardwareProfile),
					HTTPProxyConfig: nil,
					NetworkProfile:  nil,
					OsProfile:       nil,
					SecurityProfile: nil,
					StorageProfile:  nil,
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, scopeId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r StackHCIVirtualMachineResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualMachineInstances

			id, err := parse.StackHCIVirtualMachineID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroup, id.MachineName)
			scopeId := commonids.NewScopeID(arcMachineId.String())

			resp, err := client.Get(ctx, scopeId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			schema := StackHCIVirtualMachineResourceModel{
				ArcMachineId: arcMachineId.String(),
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
					schema.HardwareProfile = flattenVirtualMachineInstanceHardwareProfile(props.HardwareProfile)
					schema.HttpProxyConfiguration = nil
					schema.NetworkProfile = nil
					schema.OsProfile = nil
					schema.SecurityProfile = nil
					schema.StorageProfile = nil
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r StackHCIVirtualMachineResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualMachineInstances

			id, err := parse.StackHCIVirtualMachineID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroup, id.MachineName)
			scopeId := commonids.NewScopeID(arcMachineId.String())

			var model StackHCIVirtualMachineResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := virtualmachineinstances.VirtualMachineInstanceUpdateRequest{}

			if metadata.ResourceData.HasChange("") {
			}

			if err := client.UpdateThenPoll(ctx, scopeId, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r StackHCIVirtualMachineResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualMachineInstances

			id, err := parse.StackHCIVirtualMachineID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroup, id.MachineName)
			scopeId := commonids.NewScopeID(arcMachineId.String())

			if err := client.DeleteThenPoll(ctx, scopeId); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandVirtualMachineInstanceHardwareProfile(input []StackHCIVirtualMachineHardwareProfile) *virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	output := &virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfile{
		DynamicMemoryConfig: expandVirtualMachineInstanceDynamicMemory(v.DynamicMemory),
		MemoryMB:            pointer.To(v.MemoryMb),
		Processors:          pointer.To(v.ProcessorNumber),
		VMSize:              pointer.To(virtualmachineinstances.VMSizeEnum(v.VmSize)),
	}

	return output
}

func flattenVirtualMachineInstanceHardwareProfile(input *virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfile) []StackHCIVirtualMachineHardwareProfile {
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

func expandVirtualMachineInstanceDynamicMemory(input []StackHCIVirtualMachineDynamicMemory) *virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfileDynamicMemoryConfig {
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

func flattenVirtualMachineInstanceDynamicMemory(input *virtualmachineinstances.VirtualMachineInstancePropertiesHardwareProfileDynamicMemoryConfig) []StackHCIVirtualMachineDynamicMemory {
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
