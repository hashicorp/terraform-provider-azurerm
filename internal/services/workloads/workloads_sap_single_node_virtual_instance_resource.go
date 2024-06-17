// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workloads

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPSingleNodeVirtualInstanceModel struct {
	Name                      string                       `tfschema:"name"`
	ResourceGroupName         string                       `tfschema:"resource_group_name"`
	Location                  string                       `tfschema:"location"`
	AppLocation               string                       `tfschema:"app_location"`
	Environment               string                       `tfschema:"environment"`
	SapFqdn                   string                       `tfschema:"sap_fqdn"`
	SapProduct                string                       `tfschema:"sap_product"`
	SingleServerConfiguration []SingleServerConfiguration  `tfschema:"single_server_configuration"`
	Identity                  []identity.ModelUserAssigned `tfschema:"identity"`
	ManagedResourceGroupName  string                       `tfschema:"managed_resource_group_name"`
	Tags                      map[string]string            `tfschema:"tags"`
}

type SingleServerConfiguration struct {
	AppResourceGroupName        string                                    `tfschema:"app_resource_group_name"`
	DatabaseType                string                                    `tfschema:"database_type"`
	DiskVolumeConfigurations    []SingleServerDiskVolumeConfiguration     `tfschema:"disk_volume_configuration"`
	IsSecondaryIpEnabled        bool                                      `tfschema:"secondary_ip_enabled"`
	SubnetId                    string                                    `tfschema:"subnet_id"`
	VirtualMachineConfiguration []SingleServerVirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
	VirtualMachineResourceNames []SingleServerVirtualMachineResourceNames `tfschema:"virtual_machine_resource_names"`
}

type SingleServerDiskVolumeConfiguration struct {
	VolumeName    string `tfschema:"volume_name"`
	NumberOfDisks int64  `tfschema:"number_of_disks"`
	SizeGb        int64  `tfschema:"size_in_gb"`
	SkuName       string `tfschema:"sku_name"`
}

type SingleServerVirtualMachineConfiguration struct {
	ImageReference []SingleServerImageReference `tfschema:"image"`
	OSProfile      []SingleServerOSProfile      `tfschema:"os_profile"`
	VmSize         string                       `tfschema:"virtual_machine_size"`
}

type SingleServerImageReference struct {
	Offer     string `tfschema:"offer"`
	Publisher string `tfschema:"publisher"`
	Sku       string `tfschema:"sku"`
	Version   string `tfschema:"version"`
}

type SingleServerOSProfile struct {
	AdminUsername string `tfschema:"admin_username"`
	SshPrivateKey string `tfschema:"ssh_private_key"`
	SshPublicKey  string `tfschema:"ssh_public_key"`
}

type SingleServerVirtualMachineResourceNames struct {
	DataDisks             []SingleServerDataDisk `tfschema:"data_disk"`
	HostName              string                 `tfschema:"host_name"`
	NetworkInterfaceNames []string               `tfschema:"network_interface_names"`
	OSDiskName            string                 `tfschema:"os_disk_name"`
	VMName                string                 `tfschema:"virtual_machine_name"`
}

type SingleServerDataDisk struct {
	VolumeName string   `tfschema:"volume_name"`
	Names      []string `tfschema:"names"`
}

type WorkloadsSAPSingleNodeVirtualInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPSingleNodeVirtualInstanceResource{}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_single_node_virtual_instance"
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPSingleNodeVirtualInstanceModel{}
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapvirtualinstances.ValidateSapVirtualInstanceID
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SAPVirtualInstanceName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"app_location": commonschema.Location(),

		"sap_fqdn": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(2, 34),
		},

		"single_server_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"app_resource_group_name": commonschema.ResourceGroupName(),

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: networkValidate.SubnetID,
					},

					"virtual_machine_configuration": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"image": {
									Type:     pluginsdk.TypeList,
									Required: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"offer": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"publisher": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ForceNew: true,
												ValidateFunc: validation.StringInSlice([]string{
													"RedHat",
													"SUSE",
												}, false),
											},

											"sku": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"version": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
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
												ValidateFunc: validation.StringLenBetween(1, 64),
											},

											"ssh_private_key": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												Sensitive:    true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"ssh_public_key": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"virtual_machine_size": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"database_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForSAPDatabaseType(), false),
					},

					"disk_volume_configuration": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"volume_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										"backup",
										"hana/data",
										"hana/log",
										"hana/shared",
										"os",
										"usr/sap",
									}, false),
								},

								"number_of_disks": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntAtLeast(1),
								},

								"size_in_gb": {
									Type:     pluginsdk.TypeInt,
									Required: true,
									ForceNew: true,
								},

								"sku_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForDiskSkuName(), false),
								},
							},
						},
					},

					"secondary_ip_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"virtual_machine_resource_names": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_disk": {
									Type:         pluginsdk.TypeSet,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"volume_name": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ForceNew: true,
												ValidateFunc: validation.StringInSlice([]string{
													"default",
												}, false),
											},

											"names": {
												Type:     pluginsdk.TypeList,
												Required: true,
												ForceNew: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringLenBetween(1, 80),
												},
											},
										},
									},
								},

								"host_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
									ValidateFunc: validation.StringLenBetween(1, 13),
								},

								"network_interface_names": {
									Type:         pluginsdk.TypeList,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: networkValidate.NetworkInterfaceName,
									},
								},

								"os_disk_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
									ValidateFunc: validation.StringLenBetween(1, 80),
								},

								"virtual_machine_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
									ValidateFunc: computeValidate.VirtualMachineName,
								},
							},
						},
					},
				},
			},
		},

		"environment": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForSAPEnvironmentType(), false),
		},

		"sap_product": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForSAPProductType(), false),
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"managed_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: resourcegroups.ValidateName,
		},

		"tags": commonschema.Tags(),
	}
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if v := rd.Get("single_server_configuration.0.disk_volume_configuration"); v != nil {
				diskVolumes := v.(*pluginsdk.Set).List()
				if hasDuplicateVolumeNameForSAPSingleNodeVirtualInstance(diskVolumes) {
					return fmt.Errorf("`volume_name` cannot be duplicated")
				}
			}

			return nil
		},
	}
}

func hasDuplicateVolumeNameForSAPSingleNodeVirtualInstance(input []interface{}) bool {
	seen := make(map[string]bool)

	for _, v := range input {
		diskVolume := v.(map[string]interface{})
		volumeName := diskVolume["volume_name"].(string)

		if seen[volumeName] {
			return true
		}
		seen[volumeName] = true
	}

	return false
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPSingleNodeVirtualInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.SAPVirtualInstances
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := sapvirtualinstances.NewSapVirtualInstanceID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identity, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			parameters := &sapvirtualinstances.SAPVirtualInstance{
				Identity: identity,
				Location: location.Normalize(model.Location),
				Properties: sapvirtualinstances.SAPVirtualInstanceProperties{
					Configuration: sapvirtualinstances.DeploymentWithOSConfiguration{
						AppLocation:                 utils.String(location.Normalize(model.AppLocation)),
						InfrastructureConfiguration: expandSingleServerConfiguration(model.SingleServerConfiguration),
						OsSapConfiguration: &sapvirtualinstances.OsSapConfiguration{
							SapFqdn: utils.String(model.SapFqdn),
						},
					},
					Environment: sapvirtualinstances.SAPEnvironmentType(model.Environment),
					SapProduct:  sapvirtualinstances.SAPProductType(model.SapProduct),
				},
				Tags: &model.Tags,
			}

			if v := model.ManagedResourceGroupName; v != "" {
				parameters.Properties.ManagedResourceGroupConfiguration = &sapvirtualinstances.ManagedRGConfiguration{
					Name: utils.String(v),
				}
			}

			if err := client.CreateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPSingleNodeVirtualInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := &sapvirtualinstances.UpdateSAPVirtualInstanceRequest{}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				parameters.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = &model.Tags
			}

			if _, err := client.Update(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := WorkloadsSAPSingleNodeVirtualInstanceModel{}
			if model := resp.Model; model != nil {
				state.Name = id.SapVirtualInstanceName
				state.ResourceGroupName = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)

				identity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = pointer.From(identity)

				props := &model.Properties
				state.Environment = string(props.Environment)
				state.SapProduct = string(props.SapProduct)
				state.Tags = pointer.From(model.Tags)

				if config := props.Configuration; config != nil {
					if v, ok := config.(sapvirtualinstances.DeploymentWithOSConfiguration); ok {
						appLocation := ""
						if appLocationVal := v.AppLocation; appLocationVal != nil {
							appLocation = *v.AppLocation
						}
						state.AppLocation = location.Normalize(appLocation)

						sapFqdn := ""
						if osSapConfiguration := v.OsSapConfiguration; osSapConfiguration != nil {
							sapFqdn = pointer.From(osSapConfiguration.SapFqdn)
						}
						state.SapFqdn = sapFqdn

						if configuration := v.InfrastructureConfiguration; configuration != nil {
							if singleServerConfiguration, singleServerConfigurationExists := configuration.(sapvirtualinstances.SingleServerConfiguration); singleServerConfigurationExists {
								state.SingleServerConfiguration = flattenSingleServerConfiguration(singleServerConfiguration, metadata.ResourceData)
							}
						}
					}
				}

				if v := props.ManagedResourceGroupConfiguration; v != nil {
					state.ManagedResourceGroupName = pointer.From(v.Name)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandSAPSingleNodeVirtualInstanceVirtualMachineConfiguration(input []SingleServerVirtualMachineConfiguration) *sapvirtualinstances.VirtualMachineConfiguration {
	if len(input) == 0 {
		return nil
	}

	virtualMachineConfiguration := input[0]

	result := &sapvirtualinstances.VirtualMachineConfiguration{
		ImageReference: pointer.From(expandSAPSingleNodeVirtualInstanceImageReference(virtualMachineConfiguration.ImageReference)),
		OsProfile:      pointer.From(expandSAPSingleNodeVirtualInstanceOsProfile(virtualMachineConfiguration.OSProfile)),
		VMSize:         virtualMachineConfiguration.VmSize,
	}

	return result
}

func expandSAPSingleNodeVirtualInstanceImageReference(input []SingleServerImageReference) *sapvirtualinstances.ImageReference {
	if len(input) == 0 {
		return nil
	}

	imageReference := input[0]

	result := &sapvirtualinstances.ImageReference{
		Offer:     utils.String(imageReference.Offer),
		Publisher: utils.String(imageReference.Publisher),
		Sku:       utils.String(imageReference.Sku),
		Version:   utils.String(imageReference.Version),
	}

	return result
}

func expandSAPSingleNodeVirtualInstanceOsProfile(input []SingleServerOSProfile) *sapvirtualinstances.OSProfile {
	if len(input) == 0 {
		return nil
	}

	osProfile := input[0]

	result := &sapvirtualinstances.OSProfile{
		AdminUsername: utils.String(osProfile.AdminUsername),
		OsConfiguration: &sapvirtualinstances.LinuxConfiguration{
			DisablePasswordAuthentication: utils.Bool(true),
			SshKeyPair: &sapvirtualinstances.SshKeyPair{
				PrivateKey: utils.String(osProfile.SshPrivateKey),
				PublicKey:  utils.String(osProfile.SshPublicKey),
			},
		},
	}

	return result
}

func expandSAPSingleNodeVirtualInstanceVirtualMachineFullResourceNames(input []SingleServerVirtualMachineResourceNames) *sapvirtualinstances.SingleServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	virtualMachineFullResourceNames := input[0]

	result := &sapvirtualinstances.SingleServerFullResourceNames{
		VirtualMachine: &sapvirtualinstances.VirtualMachineResourceNames{
			DataDiskNames:     expandSAPSingleNodeVirtualInstanceDataDisks(virtualMachineFullResourceNames.DataDisks),
			NetworkInterfaces: expandSAPSingleNodeVirtualInstanceNetworkInterfaceNames(virtualMachineFullResourceNames.NetworkInterfaceNames),
		},
	}

	if v := virtualMachineFullResourceNames.HostName; v != "" {
		result.VirtualMachine.HostName = utils.String(v)
	}

	if v := virtualMachineFullResourceNames.OSDiskName; v != "" {
		result.VirtualMachine.OsDiskName = utils.String(v)
	}

	if v := virtualMachineFullResourceNames.VMName; v != "" {
		result.VirtualMachine.VirtualMachineName = utils.String(v)
	}

	return result
}

func expandSAPSingleNodeVirtualInstanceNetworkInterfaceNames(input []string) *[]sapvirtualinstances.NetworkInterfaceResourceNames {
	result := make([]sapvirtualinstances.NetworkInterfaceResourceNames, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		networkInterfaceName := sapvirtualinstances.NetworkInterfaceResourceNames{
			NetworkInterfaceName: utils.String(v),
		}

		result = append(result, networkInterfaceName)
	}

	return &result
}

func expandSAPSingleNodeVirtualInstanceDataDisks(input []SingleServerDataDisk) *map[string][]string {
	result := make(map[string][]string)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		result[v.VolumeName] = v.Names
	}

	return &result
}

func expandSAPSingleNodeVirtualInstanceDiskVolumeConfigurations(input []SingleServerDiskVolumeConfiguration) *sapvirtualinstances.DiskConfiguration {
	if len(input) == 0 {
		return nil
	}

	result := make(map[string]sapvirtualinstances.DiskVolumeConfiguration, 0)

	for _, v := range input {
		skuName := sapvirtualinstances.DiskSkuName(v.SkuName)

		result[v.VolumeName] = sapvirtualinstances.DiskVolumeConfiguration{
			Count:  utils.Int64(v.NumberOfDisks),
			SizeGB: utils.Int64(v.SizeGb),
			Sku: &sapvirtualinstances.DiskSku{
				Name: &skuName,
			},
		}
	}

	return &sapvirtualinstances.DiskConfiguration{
		DiskVolumeConfigurations: &result,
	}
}

func expandSingleServerConfiguration(input []SingleServerConfiguration) *sapvirtualinstances.SingleServerConfiguration {
	if len(input) == 0 {
		return nil
	}

	singleServerConfiguration := input[0]

	result := &sapvirtualinstances.SingleServerConfiguration{
		AppResourceGroup:    singleServerConfiguration.AppResourceGroupName,
		CustomResourceNames: expandSAPSingleNodeVirtualInstanceVirtualMachineFullResourceNames(singleServerConfiguration.VirtualMachineResourceNames),
		DbDiskConfiguration: expandSAPSingleNodeVirtualInstanceDiskVolumeConfigurations(singleServerConfiguration.DiskVolumeConfigurations),
		NetworkConfiguration: &sapvirtualinstances.NetworkConfiguration{
			IsSecondaryIPEnabled: utils.Bool(singleServerConfiguration.IsSecondaryIpEnabled),
		},
		SubnetId:                    singleServerConfiguration.SubnetId,
		VirtualMachineConfiguration: pointer.From(expandSAPSingleNodeVirtualInstanceVirtualMachineConfiguration(singleServerConfiguration.VirtualMachineConfiguration)),
	}

	if v := singleServerConfiguration.DatabaseType; v != "" {
		dbType := sapvirtualinstances.SAPDatabaseType(v)
		result.DatabaseType = &dbType
	}

	return result
}

func flattenSAPSingleNodeVirtualInstanceDiskVolumeConfigurations(input *sapvirtualinstances.DiskConfiguration) []SingleServerDiskVolumeConfiguration {
	result := make([]SingleServerDiskVolumeConfiguration, 0)
	if input == nil || input.DiskVolumeConfigurations == nil {
		return result
	}

	for k, v := range *input.DiskVolumeConfigurations {
		diskVolumeConfiguration := SingleServerDiskVolumeConfiguration{
			NumberOfDisks: pointer.From(v.Count),
			SizeGb:        pointer.From(v.SizeGB),
			VolumeName:    k,
		}

		if sku := v.Sku; sku != nil {
			diskVolumeConfiguration.SkuName = string(pointer.From(sku.Name))
		}

		result = append(result, diskVolumeConfiguration)
	}

	return result
}

func flattenSAPSingleNodeVirtualInstanceVirtualMachineFullResourceNames(input sapvirtualinstances.SingleServerFullResourceNames) []SingleServerVirtualMachineResourceNames {
	result := make([]SingleServerVirtualMachineResourceNames, 0)

	if vm := input.VirtualMachine; vm != nil {
		vmFullResourceNames := SingleServerVirtualMachineResourceNames{
			HostName:              pointer.From(vm.HostName),
			OSDiskName:            pointer.From(vm.OsDiskName),
			VMName:                pointer.From(vm.VirtualMachineName),
			NetworkInterfaceNames: flattenSAPSingleNodeVirtualInstanceNetworkInterfaceResourceNames(vm.NetworkInterfaces),
			DataDisks:             flattenSAPSingleNodeVirtualInstanceDataDisks(vm.DataDiskNames),
		}

		result = append(result, vmFullResourceNames)
	}

	return result
}

func flattenSAPSingleNodeVirtualInstanceNetworkInterfaceResourceNames(input *[]sapvirtualinstances.NetworkInterfaceResourceNames) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, pointer.From(v.NetworkInterfaceName))
	}

	return result
}

func flattenSAPSingleNodeVirtualInstanceDataDisks(input *map[string][]string) []SingleServerDataDisk {
	results := make([]SingleServerDataDisk, 0)
	if input == nil {
		return results
	}

	for k, v := range *input {
		dataDisk := SingleServerDataDisk{
			VolumeName: k,
			Names:      v,
		}

		results = append(results, dataDisk)
	}

	return results
}

func flattenSAPSingleNodeVirtualInstanceVirtualMachineConfiguration(input sapvirtualinstances.VirtualMachineConfiguration, d *pluginsdk.ResourceData) []SingleServerVirtualMachineConfiguration {
	result := make([]SingleServerVirtualMachineConfiguration, 0)

	return append(result, SingleServerVirtualMachineConfiguration{
		ImageReference: flattenSAPSingleNodeVirtualInstanceImageReference(input.ImageReference),
		OSProfile:      flattenSAPSingleNodeVirtualInstanceOSProfile(input.OsProfile, d),
		VmSize:         input.VMSize,
	})
}

func flattenSAPSingleNodeVirtualInstanceImageReference(input sapvirtualinstances.ImageReference) []SingleServerImageReference {
	result := make([]SingleServerImageReference, 0)

	return append(result, SingleServerImageReference{
		Offer:     pointer.From(input.Offer),
		Publisher: pointer.From(input.Publisher),
		Sku:       pointer.From(input.Sku),
		Version:   pointer.From(input.Version),
	})
}

func flattenSAPSingleNodeVirtualInstanceOSProfile(input sapvirtualinstances.OSProfile, d *pluginsdk.ResourceData) []SingleServerOSProfile {
	result := make([]SingleServerOSProfile, 0)

	osProfile := SingleServerOSProfile{
		AdminUsername: pointer.From(input.AdminUsername),
	}

	if osConfiguration := input.OsConfiguration; osConfiguration != nil {
		if v, ok := osConfiguration.(sapvirtualinstances.LinuxConfiguration); ok {
			if sshKeyPair := v.SshKeyPair; sshKeyPair != nil {
				osProfile.SshPrivateKey = d.Get("single_server_configuration.0.virtual_machine_configuration.0.os_profile.0.ssh_private_key").(string)
				osProfile.SshPublicKey = pointer.From(sshKeyPair.PublicKey)
			}
		}
	}

	return append(result, osProfile)
}

func flattenSingleServerConfiguration(input sapvirtualinstances.SingleServerConfiguration, d *pluginsdk.ResourceData) []SingleServerConfiguration {
	result := make([]SingleServerConfiguration, 0)

	singleServerConfig := SingleServerConfiguration{
		AppResourceGroupName:        input.AppResourceGroup,
		DatabaseType:                string(pointer.From(input.DatabaseType)),
		DiskVolumeConfigurations:    flattenSAPSingleNodeVirtualInstanceDiskVolumeConfigurations(input.DbDiskConfiguration),
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenSAPSingleNodeVirtualInstanceVirtualMachineConfiguration(input.VirtualMachineConfiguration, d),
	}

	if networkConfiguration := input.NetworkConfiguration; networkConfiguration != nil {
		singleServerConfig.IsSecondaryIpEnabled = pointer.From(networkConfiguration.IsSecondaryIPEnabled)
	}

	if customResourceNames := input.CustomResourceNames; customResourceNames != nil {
		if v, ok := customResourceNames.(sapvirtualinstances.SingleServerFullResourceNames); ok {
			singleServerConfig.VirtualMachineResourceNames = flattenSAPSingleNodeVirtualInstanceVirtualMachineFullResourceNames(v)
		}
	}

	return append(result, singleServerConfig)
}
