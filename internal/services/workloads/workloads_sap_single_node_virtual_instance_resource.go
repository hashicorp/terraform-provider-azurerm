package workloads

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
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
	OsSapConfiguration        []OsSapConfiguration         `tfschema:"os_sap_configuration"`
	SapProduct                string                       `tfschema:"sap_product"`
	SingleServerConfiguration []SingleServerConfiguration  `tfschema:"single_server_configuration"`
	Identity                  []identity.ModelUserAssigned `tfschema:"identity"`
	ManagedResourceGroupName  string                       `tfschema:"managed_resource_group_name"`
	Tags                      map[string]string            `tfschema:"tags"`
}

type OsSapConfiguration struct {
	DeployerVmPackages []DeployerVmPackages `tfschema:"deployer_virtual_machine_packages"`
	SapFqdn            string               `tfschema:"sap_fqdn"`
}

type DeployerVmPackages struct {
	StorageAccountId string `tfschema:"storage_account_id"`
	Url              string `tfschema:"url"`
}

type SingleServerConfiguration struct {
	AppResourceGroupName        string                        `tfschema:"app_resource_group_name"`
	DatabaseType                string                        `tfschema:"database_type"`
	DiskVolumeConfigurations    []DiskVolumeConfiguration     `tfschema:"disk_volume_configuration"`
	IsSecondaryIpEnabled        bool                          `tfschema:"secondary_ip_enabled"`
	SubnetId                    string                        `tfschema:"subnet_id"`
	VirtualMachineConfiguration []VirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
	VirtualMachineResourceNames []VirtualMachineResourceNames `tfschema:"virtual_machine_resource_names"`
}

type DiskVolumeConfiguration struct {
	VolumeName    string `tfschema:"volume_name"`
	NumberOfDisks int    `tfschema:"number_of_disks"`
	SizeGb        int    `tfschema:"size_in_gb"`
	SkuName       string `tfschema:"sku_name"`
}

type VirtualMachineConfiguration struct {
	ImageReference []ImageReference `tfschema:"image"`
	OSProfile      []OSProfile      `tfschema:"os_profile"`
	VmSize         string           `tfschema:"virtual_machine_size"`
}

type ImageReference struct {
	Offer     string `tfschema:"offer"`
	Publisher string `tfschema:"publisher"`
	Sku       string `tfschema:"sku"`
	Version   string `tfschema:"version"`
}

type OSProfile struct {
	AdminUsername string `tfschema:"admin_username"`
	SshPrivateKey string `tfschema:"ssh_private_key"`
	SshPublicKey  string `tfschema:"ssh_public_key"`
}

type VirtualMachineResourceNames struct {
	DataDiskNames         map[string]interface{} `tfschema:"data_disk_names"`
	HostName              string                 `tfschema:"host_name"`
	NetworkInterfaceNames []string               `tfschema:"network_interface_names"`
	OSDiskName            string                 `tfschema:"os_disk_name"`
	VMName                string                 `tfschema:"virtual_machine_name"`
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

		"os_sap_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"sap_fqdn": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(2, 34),
					},

					"deployer_virtual_machine_packages": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"storage_account_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: commonids.ValidateStorageAccountID,
								},

								"url": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsURLWithHTTPorHTTPS,
								},
							},
						},
					},
				},
			},
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
									Type:     pluginsdk.TypeInt,
									Required: true,
									ForceNew: true,
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
								"data_disk_names": {
									Type:         pluginsdk.TypeMap,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk_names", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringLenBetween(1, 80),
									},
								},

								"host_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk_names", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
									ValidateFunc: validation.StringLenBetween(1, 13),
								},

								"network_interface_names": {
									Type:         pluginsdk.TypeList,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk_names", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: networkValidate.NetworkInterfaceName,
									},
								},

								"os_disk_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk_names", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
									ValidateFunc: validation.StringLenBetween(1, 80),
								},

								"virtual_machine_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"single_server_configuration.0.virtual_machine_resource_names.0.data_disk_names", "single_server_configuration.0.virtual_machine_resource_names.0.host_name", "single_server_configuration.0.virtual_machine_resource_names.0.network_interface_names", "single_server_configuration.0.virtual_machine_resource_names.0.os_disk_name", "single_server_configuration.0.virtual_machine_resource_names.0.virtual_machine_name"},
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
						OsSapConfiguration:          expandOsSapConfiguration(model.OsSapConfiguration),
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
						state.OsSapConfiguration = flattenOsSapConfiguration(v.OsSapConfiguration)

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

func expandVirtualMachineConfiguration(input []VirtualMachineConfiguration) *sapvirtualinstances.VirtualMachineConfiguration {
	if len(input) == 0 {
		return nil
	}

	virtualMachineConfiguration := input[0]

	result := &sapvirtualinstances.VirtualMachineConfiguration{
		ImageReference: pointer.From(expandImageReference(virtualMachineConfiguration.ImageReference)),
		OsProfile:      pointer.From(expandOsProfile(virtualMachineConfiguration.OSProfile)),
		VMSize:         virtualMachineConfiguration.VmSize,
	}

	return result
}

func expandImageReference(input []ImageReference) *sapvirtualinstances.ImageReference {
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

func expandOsProfile(input []OSProfile) *sapvirtualinstances.OSProfile {
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

func expandVirtualMachineFullResourceNames(input []VirtualMachineResourceNames) *sapvirtualinstances.SingleServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	virtualMachineFullResourceNames := input[0]

	result := &sapvirtualinstances.SingleServerFullResourceNames{
		VirtualMachine: &sapvirtualinstances.VirtualMachineResourceNames{
			DataDiskNames:     expandDataDiskNames(virtualMachineFullResourceNames.DataDiskNames),
			NetworkInterfaces: expandNetworkInterfaceNames(virtualMachineFullResourceNames.NetworkInterfaceNames),
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

func expandNetworkInterfaceNames(input []string) *[]sapvirtualinstances.NetworkInterfaceResourceNames {
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

func expandDataDiskNames(input map[string]interface{}) *map[string][]string {
	result := make(map[string][]string)
	if len(input) == 0 {
		return &result
	}

	for k, v := range input {
		result[k] = strings.Split(v.(string), ",")
	}

	return &result
}

func expandDiskVolumeConfigurations(input []DiskVolumeConfiguration) *sapvirtualinstances.DiskConfiguration {
	if len(input) == 0 {
		return nil
	}

	result := make(map[string]sapvirtualinstances.DiskVolumeConfiguration, 0)

	for _, v := range input {
		skuName := sapvirtualinstances.DiskSkuName(v.SkuName)

		result[v.VolumeName] = sapvirtualinstances.DiskVolumeConfiguration{
			Count:  utils.Int64(int64(v.NumberOfDisks)),
			SizeGB: utils.Int64(int64(v.SizeGb)),
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
		CustomResourceNames: expandVirtualMachineFullResourceNames(singleServerConfiguration.VirtualMachineResourceNames),
		DbDiskConfiguration: expandDiskVolumeConfigurations(singleServerConfiguration.DiskVolumeConfigurations),
		NetworkConfiguration: &sapvirtualinstances.NetworkConfiguration{
			IsSecondaryIPEnabled: utils.Bool(singleServerConfiguration.IsSecondaryIpEnabled),
		},
		SubnetId:                    singleServerConfiguration.SubnetId,
		VirtualMachineConfiguration: pointer.From(expandVirtualMachineConfiguration(singleServerConfiguration.VirtualMachineConfiguration)),
	}

	if v := singleServerConfiguration.DatabaseType; v != "" {
		dbType := sapvirtualinstances.SAPDatabaseType(v)
		result.DatabaseType = &dbType
	}

	return result
}

func expandOsSapConfiguration(input []OsSapConfiguration) *sapvirtualinstances.OsSapConfiguration {
	if len(input) == 0 {
		return nil
	}

	osSapConfiguration := input[0]

	result := &sapvirtualinstances.OsSapConfiguration{
		DeployerVMPackages: expandDeployerVmPackages(osSapConfiguration.DeployerVmPackages),
		SapFqdn:            utils.String(osSapConfiguration.SapFqdn),
	}

	return result
}

func expandDeployerVmPackages(input []DeployerVmPackages) *sapvirtualinstances.DeployerVMPackages {
	if len(input) == 0 {
		return nil
	}

	deployerVmPackages := input[0]

	result := &sapvirtualinstances.DeployerVMPackages{
		StorageAccountId: utils.String(deployerVmPackages.StorageAccountId),
		Url:              utils.String(deployerVmPackages.Url),
	}

	return result
}

func flattenDiskVolumeConfigurations(input *sapvirtualinstances.DiskConfiguration) []DiskVolumeConfiguration {
	result := make([]DiskVolumeConfiguration, 0)
	if input == nil || input.DiskVolumeConfigurations == nil {
		return result
	}

	for k, v := range *input.DiskVolumeConfigurations {
		diskVolumeConfiguration := DiskVolumeConfiguration{
			VolumeName: k,
		}

		if count := v.Count; count != nil {
			diskVolumeConfiguration.NumberOfDisks = int(*count)
		}

		if sizeGb := v.SizeGB; sizeGb != nil {
			diskVolumeConfiguration.SizeGb = int(*sizeGb)
		}

		if sku := v.Sku; sku != nil && sku.Name != nil {
			diskVolumeConfiguration.SkuName = string(*sku.Name)
		}

		result = append(result, diskVolumeConfiguration)
	}

	return result
}

func flattenVirtualMachineFullResourceNames(input sapvirtualinstances.SingleServerFullResourceNames) []VirtualMachineResourceNames {
	result := make([]VirtualMachineResourceNames, 0)

	if vm := input.VirtualMachine; vm != nil {
		vmFullResourceNames := VirtualMachineResourceNames{
			HostName:              pointer.From(vm.HostName),
			OSDiskName:            pointer.From(vm.OsDiskName),
			VMName:                pointer.From(vm.VirtualMachineName),
			NetworkInterfaceNames: flattenNetworkInterfaceResourceNames(vm.NetworkInterfaces),
			DataDiskNames:         flattenDataDiskNames(vm.DataDiskNames),
		}

		result = append(result, vmFullResourceNames)
	}

	return result
}

func flattenNetworkInterfaceResourceNames(input *[]sapvirtualinstances.NetworkInterfaceResourceNames) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, pointer.From(v.NetworkInterfaceName))
	}

	return result
}

func flattenDataDiskNames(input *map[string][]string) map[string]interface{} {
	results := make(map[string]interface{})
	if input == nil {
		return results
	}

	for k, v := range *input {
		results[k] = strings.Join(v, ",")
	}

	return results
}

func flattenVirtualMachineConfiguration(input sapvirtualinstances.VirtualMachineConfiguration, d *pluginsdk.ResourceData) []VirtualMachineConfiguration {
	result := make([]VirtualMachineConfiguration, 0)

	return append(result, VirtualMachineConfiguration{
		ImageReference: flattenImageReference(input.ImageReference),
		OSProfile:      flattenOSProfile(input.OsProfile, d),
		VmSize:         input.VMSize,
	})
}

func flattenImageReference(input sapvirtualinstances.ImageReference) []ImageReference {
	result := make([]ImageReference, 0)

	return append(result, ImageReference{
		Offer:     pointer.From(input.Offer),
		Publisher: pointer.From(input.Publisher),
		Sku:       pointer.From(input.Sku),
		Version:   pointer.From(input.Version),
	})
}

func flattenOSProfile(input sapvirtualinstances.OSProfile, d *pluginsdk.ResourceData) []OSProfile {
	result := make([]OSProfile, 0)

	osProfile := OSProfile{
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
		DiskVolumeConfigurations:    flattenDiskVolumeConfigurations(input.DbDiskConfiguration),
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d),
	}

	if networkConfiguration := input.NetworkConfiguration; networkConfiguration != nil {
		singleServerConfig.IsSecondaryIpEnabled = pointer.From(networkConfiguration.IsSecondaryIPEnabled)
	}

	if customResourceNames := input.CustomResourceNames; customResourceNames != nil {
		if v, ok := customResourceNames.(sapvirtualinstances.SingleServerFullResourceNames); ok {
			singleServerConfig.VirtualMachineResourceNames = flattenVirtualMachineFullResourceNames(v)
		}
	}

	return append(result, singleServerConfig)
}

func flattenOsSapConfiguration(input *sapvirtualinstances.OsSapConfiguration) []OsSapConfiguration {
	result := make([]OsSapConfiguration, 0)
	if input == nil {
		return result
	}

	return append(result, OsSapConfiguration{
		DeployerVmPackages: flattenDeployerVMPackages(input.DeployerVMPackages),
		SapFqdn:            pointer.From(input.SapFqdn),
	})
}

func flattenDeployerVMPackages(input *sapvirtualinstances.DeployerVMPackages) []DeployerVmPackages {
	result := make([]DeployerVmPackages, 0)
	if input == nil {
		return result
	}

	return append(result, DeployerVmPackages{
		StorageAccountId: pointer.From(input.StorageAccountId),
		Url:              pointer.From(input.Url),
	})
}
