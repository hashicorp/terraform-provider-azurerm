package workloads

import (
	"fmt"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/validate"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func SchemaForSAPVirtualInstanceSingleServerConfiguration(deploymentType string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
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

				"database_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForSAPDatabaseType(), false),
				},

				"disk_volume_configuration": SchemaForSAPVirtualInstanceDiskVolumeConfiguration(),

				"is_secondary_ip_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"virtual_machine_configuration": SchemaForSAPVirtualInstanceVirtualMachineConfiguration(),

				"virtual_machine_full_resource_names": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"data_disk_names": {
								Type:     pluginsdk.TypeMap,
								Optional: true,
								ForceNew: true,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validate.DiskName,
								},
							},

							"host_name": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validate.HostName,
							},

							"network_interface_names": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: networkValidate.NetworkInterfaceName,
								},
							},

							"os_disk_name": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validate.DiskName,
							},

							"vm_name": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: computeValidate.VirtualMachineName,
							},
						},
					},
				},
			},
		},
		AtLeastOneOf: []string{
			fmt.Sprintf("%s.0.single_server_configuration", deploymentType),
			fmt.Sprintf("%s.0.three_tier_configuration", deploymentType),
		},
	}
}

func SchemaForSAPVirtualInstanceVirtualMachineConfiguration() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"image_reference": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"offer": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"publisher": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"sku": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"version": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
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
							"admin_password": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"admin_username": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"linux_configuration": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"password_authentication_enabled": {
											Type:     pluginsdk.TypeBool,
											Optional: true,
											ForceNew: true,
										},

										"ssh_key_pair": {
											Type:     pluginsdk.TypeList,
											Optional: true,
											ForceNew: true,
											MaxItems: 1,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													"private_key": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},

													"public_key": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},

				"vm_size": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaForSAPVirtualInstanceDiskVolumeConfiguration() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"volume_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"count": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"size_gb": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"sku_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaForSAPVirtualInstanceThreeTierConfiguration(deploymentType string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"app_resource_group_name": commonschema.ResourceGroupName(),

				"application_server_configuration": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"instance_count": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
								ForceNew: true,
							},

							"subnet_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: networkValidate.SubnetID,
							},

							"virtual_machine_configuration": SchemaForSAPVirtualInstanceVirtualMachineConfiguration(),
						},
					},
				},

				"central_server_configuration": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"instance_count": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
								ForceNew: true,
							},

							"subnet_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: networkValidate.SubnetID,
							},

							"virtual_machine_configuration": SchemaForSAPVirtualInstanceVirtualMachineConfiguration(),
						},
					},
				},

				"database_server_configuration": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"database_type": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"disk_volume_configuration": SchemaForSAPVirtualInstanceDiskVolumeConfiguration(),

							"instance_count": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
								ForceNew: true,
							},

							"subnet_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: networkValidate.SubnetID,
							},

							"virtual_machine_configuration": SchemaForSAPVirtualInstanceVirtualMachineConfiguration(),
						},
					},
				},

				"full_resource_names": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"application_server": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"availability_set_name": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ForceNew:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},

										"virtual_machine": {
											Type:     pluginsdk.TypeList,
											Optional: true,
											ForceNew: true,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													"data_disk_names": {
														Type:     pluginsdk.TypeMap,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},

													"host_name": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},

													"network_interface_names": {
														Type:     pluginsdk.TypeList,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},

													"os_disk_name": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},

													"vm_name": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},
												},
											},
										},
									},
								},
							},

							"central_server": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"availability_set_name": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ForceNew:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},

										"load_balancer": {
											Type:     pluginsdk.TypeList,
											Optional: true,
											ForceNew: true,
											MaxItems: 1,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													"name": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},

													"backend_pool_names": {
														Type:     pluginsdk.TypeList,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},

													"frontend_ip_configuration_names": {
														Type:     pluginsdk.TypeList,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},

													"health_probe_names": {
														Type:     pluginsdk.TypeList,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},
										},

										"virtual_machine": {
											Type:     pluginsdk.TypeList,
											Optional: true,
											ForceNew: true,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													"data_disk_names": {
														Type:     pluginsdk.TypeMap,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},

													"host_name": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},

													"network_interface_names": {
														Type:     pluginsdk.TypeList,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},

													"os_disk_name": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},

													"vm_name": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},
												},
											},
										},
									},
								},
							},

							"database_server": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"availability_set_name": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ForceNew:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},

										"load_balancer": {
											Type:     pluginsdk.TypeList,
											Optional: true,
											ForceNew: true,
											MaxItems: 1,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													"backend_pool_names": {
														Type:     pluginsdk.TypeList,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},

													"frontend_ip_configuration_names": {
														Type:     pluginsdk.TypeList,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},

													"health_probe_names": {
														Type:     pluginsdk.TypeList,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},

													"name": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},
												},
											},
										},

										"virtual_machine": {
											Type:     pluginsdk.TypeList,
											Optional: true,
											ForceNew: true,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													"data_disk_names": {
														Type:     pluginsdk.TypeMap,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},

													"host_name": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},

													"network_interface_names": {
														Type:     pluginsdk.TypeList,
														Optional: true,
														ForceNew: true,
														Elem: &pluginsdk.Schema{
															Type:         pluginsdk.TypeString,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},

													"os_disk_name": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},

													"vm_name": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},
												},
											},
										},
									},
								},
							},

							"shared_storage": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"account_name": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ForceNew:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},

										"private_endpoint_name": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ForceNew:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},
					},
				},

				"high_availability_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"is_secondary_ip_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"storage_configuration": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"transport_create_and_mount": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"resource_group_name": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ForceNew:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},

										"storage_account_name": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ForceNew:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},

							"transport_mount": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"file_share_id": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ForceNew:     true,
											ValidateFunc: storageValidate.StorageShareID,
										},

										"private_endpoint_id": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ForceNew:     true,
											ValidateFunc: networkValidate.PrivateEndpointID,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		AtLeastOneOf: []string{
			fmt.Sprintf("%s.0.single_server_configuration", deploymentType),
			fmt.Sprintf("%s.0.three_tier_configuration", deploymentType),
		},
	}
}

func expandSingleServerConfiguration(input []SingleServerConfiguration) *sapvirtualinstances.SingleServerConfiguration {
	if len(input) == 0 {
		return nil
	}

	singleServerConfiguration := &input[0]

	result := sapvirtualinstances.SingleServerConfiguration{
		AppResourceGroup: singleServerConfiguration.AppResourceGroupName,
		NetworkConfiguration: &sapvirtualinstances.NetworkConfiguration{
			IsSecondaryIPEnabled: utils.Bool(singleServerConfiguration.IsSecondaryIpEnabled), //maybe need to use GetRawConfig
		},
		SubnetId: singleServerConfiguration.SubnetId,
	}

	if v := singleServerConfiguration.DatabaseType; v != "" {
		dbType := sapvirtualinstances.SAPDatabaseType(v)
		result.DatabaseType = &dbType
	}

	if v := singleServerConfiguration.DiskVolumeConfigurations; v != nil {
		result.DbDiskConfiguration = expandDiskVolumeConfigurations(v)
	}

	if v := singleServerConfiguration.VirtualMachineConfiguration; v != nil {
		result.VirtualMachineConfiguration = expandVirtualMachineConfiguration(v)
	}

	if v := singleServerConfiguration.VirtualMachineFullResourceNames; v != nil {
		result.CustomResourceNames = expandVirtualMachineFullResourceNames(v)
	}

	return &result
}

func expandVirtualMachineConfiguration(input []VirtualMachineConfiguration) sapvirtualinstances.VirtualMachineConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.VirtualMachineConfiguration{}
	}

	virtualMachineConfiguration := &input[0]

	result := sapvirtualinstances.VirtualMachineConfiguration{}

	if v := virtualMachineConfiguration.VmSize; v != "" {
		result.VMSize = v
	}

	if v := virtualMachineConfiguration.ImageReference; v != nil {
		result.ImageReference = expandImageReference(v)
	}

	if v := virtualMachineConfiguration.OSProfile; v != nil {
		result.OsProfile = expandOsProfile(v)
	}

	return result
}

func expandImageReference(input []ImageReference) sapvirtualinstances.ImageReference {
	if len(input) == 0 {
		return sapvirtualinstances.ImageReference{}
	}

	imageReference := &input[0]

	result := sapvirtualinstances.ImageReference{}

	if v := imageReference.Offer; v != "" {
		result.Offer = utils.String(v)
	}

	if v := imageReference.Publisher; v != "" {
		result.Publisher = utils.String(v)
	}

	if v := imageReference.Sku; v != "" {
		result.Sku = utils.String(v)
	}

	if v := imageReference.Version; v != "" {
		result.Version = utils.String(v)
	}

	return result
}

func expandOsProfile(input []OSProfile) sapvirtualinstances.OSProfile {
	if len(input) == 0 {
		return sapvirtualinstances.OSProfile{}
	}

	osProfile := &input[0]

	result := sapvirtualinstances.OSProfile{}

	if v := osProfile.AdminUsername; v != "" {
		result.AdminUsername = utils.String(v)
	}

	if v := osProfile.AdminPassword; v != "" {
		result.AdminPassword = utils.String(v)
	}

	if v := osProfile.LinuxConfiguration; v != nil {
		result.OsConfiguration = expandLinuxConfiguration(v)
	} else {
		result.OsConfiguration = sapvirtualinstances.WindowsConfiguration{}
	}

	return result
}

func expandLinuxConfiguration(input []LinuxConfiguration) sapvirtualinstances.LinuxConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.LinuxConfiguration{}
	}

	linuxConfiguration := &input[0]

	result := sapvirtualinstances.LinuxConfiguration{
		DisablePasswordAuthentication: utils.Bool(!linuxConfiguration.PasswordAuthenticationEnabled), // mabye need GetRawConfig
	}

	if v := linuxConfiguration.SshKeyPair; v != nil {
		result.SshKeyPair = expandSshKeyPair(v)
	}

	return result
}

func expandSshKeyPair(input []SshKeyPair) *sapvirtualinstances.SshKeyPair {
	if len(input) == 0 {
		return nil
	}

	sshKeyPair := &input[0]

	result := sapvirtualinstances.SshKeyPair{}

	if v := sshKeyPair.PrivateKey; v != "" {
		result.PrivateKey = utils.String(v)
	}

	if v := sshKeyPair.PublicKey; v != "" {
		result.PublicKey = utils.String(v)
	}

	return &result
}

func expandVirtualMachineFullResourceNames(input []VirtualMachineFullResourceNames) sapvirtualinstances.SingleServerFullResourceNames {
	if len(input) == 0 {
		return sapvirtualinstances.SingleServerFullResourceNames{}
	}

	virtualMachineFullResourceNames := &input[0]

	result := sapvirtualinstances.SingleServerFullResourceNames{
		VirtualMachine: &sapvirtualinstances.VirtualMachineResourceNames{},
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

	if v := virtualMachineFullResourceNames.NetworkInterfaceNames; v != nil {
		result.VirtualMachine.NetworkInterfaces = expandNetworkInterfaceNames(v)
	}

	if v := virtualMachineFullResourceNames.DataDiskNames; v != nil {
		result.VirtualMachine.DataDiskNames = expandDataDiskNames(v)
	}

	return result
}

func expandNetworkInterfaceNames(input []string) *[]sapvirtualinstances.NetworkInterfaceResourceNames {
	if len(input) == 0 {
		return nil
	}

	result := make([]sapvirtualinstances.NetworkInterfaceResourceNames, 0)

	for _, v := range input {
		networkInterfaceName := sapvirtualinstances.NetworkInterfaceResourceNames{
			NetworkInterfaceName: utils.String(v),
		}

		result = append(result, networkInterfaceName)
	}

	return &result
}

func expandDataDiskNames(input map[string]interface{}) *map[string][]string {
	if len(input) == 0 {
		return nil
	}

	result := make(map[string][]string)

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
			Count:  utils.Int64(v.Count),  // maybe need GetRawConfig
			SizeGB: utils.Int64(v.SizeGb), // maybe need GetRawConfig
			Sku: &sapvirtualinstances.DiskSku{
				Name: &skuName,
			},
		}
	}

	return &sapvirtualinstances.DiskConfiguration{
		DiskVolumeConfigurations: &result,
	}
}

func flattenSingleServerConfiguration(input sapvirtualinstances.SingleServerConfiguration) []SingleServerConfiguration {
	result := SingleServerConfiguration{
		AppResourceGroupName:        input.AppResourceGroup,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration),
	}

	if v := input.DatabaseType; v != nil {
		result.DatabaseType = string(*v)
	}

	if networkConfiguration := input.NetworkConfiguration; networkConfiguration != nil && networkConfiguration.IsSecondaryIPEnabled != nil {
		result.IsSecondaryIpEnabled = *networkConfiguration.IsSecondaryIPEnabled
	}

	if v := input.DbDiskConfiguration; v != nil && v.DiskVolumeConfigurations != nil {
		result.DiskVolumeConfigurations = flattenDiskVolumeConfigurations(v.DiskVolumeConfigurations)
	}

	if customResourceNames := input.CustomResourceNames; customResourceNames != nil {
		if v, ok := customResourceNames.(sapvirtualinstances.SingleServerFullResourceNames); ok {
			result.VirtualMachineFullResourceNames = flattenVirtualMachineFullResourceNames(v)
		}
	}

	return []SingleServerConfiguration{
		result,
	}
}

func flattenDiskVolumeConfigurations(input *map[string]sapvirtualinstances.DiskVolumeConfiguration) []DiskVolumeConfiguration {
	if input == nil {
		return nil
	}

	result := make([]DiskVolumeConfiguration, 0)

	for k, v := range *input {
		diskVolumeConfiguration := DiskVolumeConfiguration{
			VolumeName: k,
		}

		if count := v.Count; count != nil {
			diskVolumeConfiguration.Count = *count
		}

		if sizeGB := v.SizeGB; sizeGB != nil {
			diskVolumeConfiguration.SizeGb = *sizeGB
		}

		if sku := v.Sku; sku != nil && sku.Name != nil {
			diskVolumeConfiguration.SkuName = string(*sku.Name)
		}

		result = append(result, diskVolumeConfiguration)
	}

	return result
}

func flattenVirtualMachineFullResourceNames(input sapvirtualinstances.SingleServerFullResourceNames) []VirtualMachineFullResourceNames {
	result := VirtualMachineFullResourceNames{}

	if vm := input.VirtualMachine; vm != nil {
		if v := vm.HostName; v != nil {
			result.HostName = *v
		}

		if v := vm.OsDiskName; v != nil {
			result.OSDiskName = *v
		}

		if v := vm.VirtualMachineName; v != nil {
			result.VMName = *v
		}

		if v := vm.NetworkInterfaces; v != nil {
			result.NetworkInterfaceNames = flattenNetworkInterfaceResourceNames(v)
		}

		if v := vm.DataDiskNames; v != nil {
			result.DataDiskNames = flattenDataDiskNames(v)
		}
	}

	return []VirtualMachineFullResourceNames{
		result,
	}
}

func flattenNetworkInterfaceResourceNames(input *[]sapvirtualinstances.NetworkInterfaceResourceNames) []string {
	if input == nil {
		return nil
	}

	result := make([]string, 0)

	for _, v := range *input {
		result = append(result, *v.NetworkInterfaceName)
	}

	return result
}

func flattenDataDiskNames(input *map[string][]string) map[string]interface{} {
	if input == nil {
		return nil
	}

	results := make(map[string]interface{})

	for k, v := range *input {
		results[k] = strings.Join(v, ",")
	}

	return results
}

func flattenVirtualMachineConfiguration(input sapvirtualinstances.VirtualMachineConfiguration) []VirtualMachineConfiguration {
	result := VirtualMachineConfiguration{
		ImageReference: flattenImageReference(input.ImageReference),
		OSProfile:      flattenOSProfile(input.OsProfile),
	}

	if v := input.VMSize; v != "" {
		result.VmSize = v
	}

	return []VirtualMachineConfiguration{
		result,
	}
}

func flattenImageReference(input sapvirtualinstances.ImageReference) []ImageReference {
	result := ImageReference{}

	if v := input.Offer; v != nil {
		result.Offer = *v
	}

	if v := input.Publisher; v != nil {
		result.Publisher = *v
	}

	if v := input.Sku; v != nil {
		result.Sku = *v
	}

	if v := input.Version; v != nil {
		result.Version = *v
	}

	return []ImageReference{
		result,
	}
}

func flattenOSProfile(input sapvirtualinstances.OSProfile) []OSProfile {
	result := OSProfile{}

	if v := input.AdminPassword; v != nil {
		result.AdminPassword = *v
	}

	if v := input.AdminUsername; v != nil {
		result.AdminUsername = *v
	}

	if osConfiguration := input.OsConfiguration; osConfiguration != nil {
		if v, ok := osConfiguration.(sapvirtualinstances.LinuxConfiguration); ok {
			result.LinuxConfiguration = flattenLinuxConfiguration(v)
		}
	}

	return []OSProfile{
		result,
	}
}

func flattenLinuxConfiguration(input sapvirtualinstances.LinuxConfiguration) []LinuxConfiguration {
	result := LinuxConfiguration{}

	if v := input.DisablePasswordAuthentication; v != nil {
		result.PasswordAuthenticationEnabled = !*v
	}

	if v := input.SshKeyPair; v != nil {
		result.SshKeyPair = flattenSshKeyPair(v)
	}

	return []LinuxConfiguration{
		result,
	}
}

func flattenSshKeyPair(input *sapvirtualinstances.SshKeyPair) []SshKeyPair {
	if input == nil {
		return nil
	}

	result := SshKeyPair{}

	if v := input.PrivateKey; v != nil {
		result.PrivateKey = *v
	}

	if v := input.PublicKey; v != nil {
		result.PublicKey = *v
	}

	return []SshKeyPair{
		result,
	}
}

func expandThreeTierConfiguration(input []ThreeTierConfiguration) *sapvirtualinstances.ThreeTierConfiguration {
	if len(input) == 0 {
		return nil
	}

	threeTierConfiguration := &input[0]

	result := sapvirtualinstances.ThreeTierConfiguration{
		AppResourceGroup: threeTierConfiguration.AppResourceGroupName,
		NetworkConfiguration: &sapvirtualinstances.NetworkConfiguration{
			IsSecondaryIPEnabled: utils.Bool(threeTierConfiguration.IsSecondaryIpEnabled), //maybe need to use GetRawConfig
		},
	}

	if v := threeTierConfiguration.HighAvailabilityType; v != "" {
		result.HighAvailabilityConfig = &sapvirtualinstances.HighAvailabilityConfiguration{
			HighAvailabilityType: sapvirtualinstances.SAPHighAvailabilityType(v),
		}
	}

	if v := threeTierConfiguration.ApplicationServerConfiguration; v != nil {
		result.ApplicationServer = expandApplicationServer(v)
	}

	if v := threeTierConfiguration.CentralServerConfiguration; v != nil {
		result.CentralServer = expandCentralServer(v)
	}

	if v := threeTierConfiguration.DatabaseServerConfiguration; v != nil {
		result.DatabaseServer = expandDatabaseServer(v)
	}

	if v := threeTierConfiguration.StorageConfiguration; v != nil {
		result.StorageConfiguration = expandStorageConfiguration(v)
	}

	if v := threeTierConfiguration.FullResourceNames; v != nil {
		result.CustomResourceNames = expandFullResourceNames(v)
	}

	return &result
}

func expandApplicationServer(input []ApplicationServerConfiguration) sapvirtualinstances.ApplicationServerConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.ApplicationServerConfiguration{}
	}

	applicationServer := &input[0]

	result := sapvirtualinstances.ApplicationServerConfiguration{
		InstanceCount: applicationServer.InstanceCount, // Maybe it needs Getrawconfig
	}

	if v := applicationServer.SubnetId; v != "" {
		result.SubnetId = v
	}

	if v := applicationServer.VirtualMachineConfiguration; v != nil {
		result.VirtualMachineConfiguration = expandVirtualMachineConfiguration(v)
	}

	return result
}

func expandCentralServer(input []CentralServerConfiguration) sapvirtualinstances.CentralServerConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.CentralServerConfiguration{}
	}

	centralServer := &input[0]

	result := sapvirtualinstances.CentralServerConfiguration{
		InstanceCount: centralServer.InstanceCount, // Maybe it needs Getrawconfig
	}

	if v := centralServer.SubnetId; v != "" {
		result.SubnetId = v
	}

	if v := centralServer.VirtualMachineConfiguration; v != nil {
		result.VirtualMachineConfiguration = expandVirtualMachineConfiguration(v)
	}

	return result
}

func expandDatabaseServer(input []DatabaseServerConfiguration) sapvirtualinstances.DatabaseConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.DatabaseConfiguration{}
	}

	databaseServer := &input[0]

	result := sapvirtualinstances.DatabaseConfiguration{
		InstanceCount: databaseServer.InstanceCount, // Maybe it needs Getrawconfig
	}

	if v := databaseServer.SubnetId; v != "" {
		result.SubnetId = v
	}

	if v := databaseServer.DatabaseType; v != "" {
		dbType := sapvirtualinstances.SAPDatabaseType(v)
		result.DatabaseType = &dbType
	}

	if v := databaseServer.DiskVolumeConfigurations; v != nil {
		result.DiskConfiguration = expandDiskVolumeConfigurations(v)
	}

	if v := databaseServer.VirtualMachineConfiguration; v != nil {
		result.VirtualMachineConfiguration = expandVirtualMachineConfiguration(v)
	}

	return result
}

func expandStorageConfiguration(input []StorageConfiguration) *sapvirtualinstances.StorageConfiguration {
	if len(input) == 0 {
		return nil
	}

	storageConfiguration := &input[0]

	result := sapvirtualinstances.StorageConfiguration{}

	if v := storageConfiguration.TransportCreateAndMount; v != nil {
		result.TransportFileShareConfiguration = expandTransportCreateAndMount(v)
	}

	if v := storageConfiguration.TransportMount; v != nil {
		result.TransportFileShareConfiguration = expandTransportMount(v)
	}

	return &result
}

func expandTransportCreateAndMount(input []TransportCreateAndMount) sapvirtualinstances.CreateAndMountFileShareConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.CreateAndMountFileShareConfiguration{}
	}

	transportCreateAndMount := &input[0]

	result := sapvirtualinstances.CreateAndMountFileShareConfiguration{}

	if v := transportCreateAndMount.ResourceGroupName; v != "" {
		result.ResourceGroup = utils.String(v)
	}

	if v := transportCreateAndMount.StorageAccountName; v != "" {
		result.StorageAccountName = utils.String(v)
	}

	return result
}

func expandTransportMount(input []TransportMount) sapvirtualinstances.MountFileShareConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.MountFileShareConfiguration{}
	}

	transportMount := &input[0]

	result := sapvirtualinstances.MountFileShareConfiguration{}

	if v := transportMount.FileShareId; v != "" {
		result.Id = v
	}

	if v := transportMount.PrivateEndpointId; v != "" {
		result.PrivateEndpointId = v
	}

	return result
}

func expandFullResourceNames(input []FullResourceNames) sapvirtualinstances.ThreeTierFullResourceNames {
	if len(input) == 0 {
		return sapvirtualinstances.ThreeTierFullResourceNames{}
	}

	fullResourceNames := &input[0]

	result := sapvirtualinstances.ThreeTierFullResourceNames{
		ApplicationServer: expandApplicationServerFullResourceNames(fullResourceNames.ApplicationServer),
		CentralServer:     expandCentralServerFullResourceNames(fullResourceNames.CentralServer),
		DatabaseServer:    expandDatabaseServerFullResourceNames(fullResourceNames.DatabaseServer),
		SharedStorage:     expandSharedStorage(fullResourceNames.SharedStorage),
	}

	return result
}

func expandApplicationServerFullResourceNames(input []ApplicationServerFullResourceNames) *sapvirtualinstances.ApplicationServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	applicationServerFullResourceNames := &input[0]

	result := sapvirtualinstances.ApplicationServerFullResourceNames{
		VirtualMachines: expandVirtualMachinesFullResourceNames(applicationServerFullResourceNames.VirtualMachines),
	}

	if v := applicationServerFullResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = utils.String(v)
	}

	return &result
}

func expandVirtualMachinesFullResourceNames(input []VirtualMachineFullResourceNames) *[]sapvirtualinstances.VirtualMachineResourceNames {
	if len(input) == 0 {
		return nil
	}

	result := make([]sapvirtualinstances.VirtualMachineResourceNames, 0)

	for _, item := range input {
		vmResourceNames := sapvirtualinstances.VirtualMachineResourceNames{}

		if v := item.HostName; v != "" {
			vmResourceNames.HostName = utils.String(v)
		}

		if v := item.OSDiskName; v != "" {
			vmResourceNames.OsDiskName = utils.String(v)
		}

		if v := item.VMName; v != "" {
			vmResourceNames.VirtualMachineName = utils.String(v)
		}

		if v := item.NetworkInterfaceNames; v != nil {
			vmResourceNames.NetworkInterfaces = expandNetworkInterfaceNames(v)
		}

		if v := item.DataDiskNames; v != nil {
			vmResourceNames.DataDiskNames = expandDataDiskNames(v)
		}

		result = append(result, vmResourceNames)
	}

	return &result
}

func expandCentralServerFullResourceNames(input []CentralServerFullResourceNames) *sapvirtualinstances.CentralServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	centralServerFullResourceNames := &input[0]

	result := sapvirtualinstances.CentralServerFullResourceNames{
		LoadBalancer:    expandLoadBalancerFullResourceNames(centralServerFullResourceNames.LoadBalancer),
		VirtualMachines: expandVirtualMachinesFullResourceNames(centralServerFullResourceNames.VirtualMachines),
	}

	if v := centralServerFullResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = utils.String(v)
	}

	return &result
}

func expandLoadBalancerFullResourceNames(input []LoadBalancer) *sapvirtualinstances.LoadBalancerResourceNames {
	if len(input) == 0 {
		return nil
	}

	loadBalancerFullResourceNames := &input[0]

	result := sapvirtualinstances.LoadBalancerResourceNames{}

	if v := loadBalancerFullResourceNames.Name; v != "" {
		result.LoadBalancerName = utils.String(v)
	}

	if v := loadBalancerFullResourceNames.BackendPoolNames; v != nil {
		result.BackendPoolNames = &v
	}

	if v := loadBalancerFullResourceNames.FrontendIpConfigurationNames; v != nil {
		result.FrontendIPConfigurationNames = &v
	}

	if v := loadBalancerFullResourceNames.HealthProbeNames; v != nil {
		result.HealthProbeNames = &v
	}

	return &result
}

func expandDatabaseServerFullResourceNames(input []DatabaseServerFullResourceNames) *sapvirtualinstances.DatabaseServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	databaseServerFullResourceNames := &input[0]

	result := sapvirtualinstances.DatabaseServerFullResourceNames{
		LoadBalancer:    expandLoadBalancerFullResourceNames(databaseServerFullResourceNames.LoadBalancer),
		VirtualMachines: expandVirtualMachinesFullResourceNames(databaseServerFullResourceNames.VirtualMachines),
	}

	if v := databaseServerFullResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = utils.String(v)
	}

	return &result
}

func expandSharedStorage(input []SharedStorage) *sapvirtualinstances.SharedStorageResourceNames {
	if len(input) == 0 {
		return nil
	}

	sharedStorage := &input[0]

	result := sapvirtualinstances.SharedStorageResourceNames{}

	if v := sharedStorage.AccountName; v != "" {
		result.SharedStorageAccountName = utils.String(v)
	}

	if v := sharedStorage.PrivateEndpointName; v != "" {
		result.SharedStorageAccountPrivateEndPointName = utils.String(v)
	}

	return &result
}

func flattenThreeTierConfiguration(input sapvirtualinstances.ThreeTierConfiguration) []ThreeTierConfiguration {
	result := ThreeTierConfiguration{
		AppResourceGroupName:           input.AppResourceGroup,
		ApplicationServerConfiguration: flattenApplicationServer(input.ApplicationServer),
		CentralServerConfiguration:     flattenCentralServer(input.CentralServer),
		DatabaseServerConfiguration:    flattenDatabaseServer(input.DatabaseServer),
	}

	if customResourceNames := input.CustomResourceNames; customResourceNames != nil {
		if v, ok := customResourceNames.(sapvirtualinstances.ThreeTierFullResourceNames); ok {
			result.FullResourceNames = flattenFullResourceNames(v)
		}
	}

	if v := input.HighAvailabilityConfig; v != nil && v.HighAvailabilityType != "" {
		result.HighAvailabilityType = string(v.HighAvailabilityType)
	}

	if networkConfiguration := input.NetworkConfiguration; networkConfiguration != nil && networkConfiguration.IsSecondaryIPEnabled != nil {
		result.IsSecondaryIpEnabled = *networkConfiguration.IsSecondaryIPEnabled
	}

	if v := input.StorageConfiguration; v != nil {
		result.StorageConfiguration = flattenStorageConfiguration(v)
	}

	return []ThreeTierConfiguration{
		result,
	}
}

func flattenStorageConfiguration(input *sapvirtualinstances.StorageConfiguration) []StorageConfiguration {
	if input == nil {
		return nil
	}

	result := StorageConfiguration{}

	if transportFileShareConfiguration := input.TransportFileShareConfiguration; transportFileShareConfiguration != nil {
		if createAndMountFileShareConfiguration, ok := transportFileShareConfiguration.(sapvirtualinstances.CreateAndMountFileShareConfiguration); ok {
			transportCreateAndMount := TransportCreateAndMount{}

			if v := createAndMountFileShareConfiguration.ResourceGroup; v != nil {
				transportCreateAndMount.ResourceGroupName = *v
			}

			if v := createAndMountFileShareConfiguration.StorageAccountName; v != nil {
				transportCreateAndMount.StorageAccountName = *v
			}

			result.TransportCreateAndMount = []TransportCreateAndMount{
				transportCreateAndMount,
			}
		}

		if mountFileShareConfiguration, ok := transportFileShareConfiguration.(sapvirtualinstances.MountFileShareConfiguration); ok {
			transportMount := TransportMount{}

			if v := mountFileShareConfiguration.Id; v != "" {
				transportMount.FileShareId = v
			}

			if v := mountFileShareConfiguration.PrivateEndpointId; v != "" {
				transportMount.PrivateEndpointId = v
			}

			result.TransportMount = []TransportMount{
				transportMount,
			}
		}
	}

	return []StorageConfiguration{
		result,
	}
}

func flattenApplicationServer(input sapvirtualinstances.ApplicationServerConfiguration) []ApplicationServerConfiguration {
	result := ApplicationServerConfiguration{
		InstanceCount:               input.InstanceCount,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration),
	}

	if v := input.SubnetId; v != "" {
		result.SubnetId = v
	}

	return []ApplicationServerConfiguration{
		result,
	}
}

func flattenCentralServer(input sapvirtualinstances.CentralServerConfiguration) []CentralServerConfiguration {
	result := CentralServerConfiguration{
		InstanceCount:               input.InstanceCount,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration),
	}

	if v := input.SubnetId; v != "" {
		result.SubnetId = v
	}

	return []CentralServerConfiguration{
		result,
	}
}

func flattenDatabaseServer(input sapvirtualinstances.DatabaseConfiguration) []DatabaseServerConfiguration {
	result := DatabaseServerConfiguration{
		InstanceCount:               input.InstanceCount,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration),
	}

	if v := input.DatabaseType; v != nil {
		result.DatabaseType = string(*v)
	}

	if v := input.DiskConfiguration; v != nil && v.DiskVolumeConfigurations != nil {
		result.DiskVolumeConfigurations = flattenDiskVolumeConfigurations(v.DiskVolumeConfigurations)
	}

	if v := input.SubnetId; v != "" {
		result.SubnetId = v
	}

	return []DatabaseServerConfiguration{
		result,
	}
}

func flattenFullResourceNames(input sapvirtualinstances.ThreeTierFullResourceNames) []FullResourceNames {
	result := FullResourceNames{}

	if v := input.ApplicationServer; v != nil {
		result.ApplicationServer = flattenApplicationServerFullResourceNames(v)
	}

	if v := input.CentralServer; v != nil {
		result.CentralServer = flattenCentralServerFullResourceNames(v)
	}

	if v := input.DatabaseServer; v != nil {
		result.DatabaseServer = flattenDatabaseServerFullResourceNames(v)
	}

	if v := input.SharedStorage; v != nil {
		result.SharedStorage = flattenSharedStorage(v)
	}

	return []FullResourceNames{
		result,
	}
}

func flattenApplicationServerFullResourceNames(input *sapvirtualinstances.ApplicationServerFullResourceNames) []ApplicationServerFullResourceNames {
	if input == nil {
		return nil
	}

	result := ApplicationServerFullResourceNames{}

	if v := input.AvailabilitySetName; v != nil {
		result.AvailabilitySetName = *v
	}

	if v := input.VirtualMachines; v != nil {
		result.VirtualMachines = flattenVirtualMachinesFullResourceNames(v)
	}

	return []ApplicationServerFullResourceNames{
		result,
	}
}

func flattenVirtualMachinesFullResourceNames(input *[]sapvirtualinstances.VirtualMachineResourceNames) []VirtualMachineFullResourceNames {
	if input == nil {
		return nil
	}

	result := make([]VirtualMachineFullResourceNames, 0)

	for _, item := range *input {
		virtualMachineFullResourceNames := VirtualMachineFullResourceNames{}

		if v := item.HostName; v != nil {
			virtualMachineFullResourceNames.HostName = *v
		}

		if v := item.OsDiskName; v != nil {
			virtualMachineFullResourceNames.OSDiskName = *v
		}

		if v := item.VirtualMachineName; v != nil {
			virtualMachineFullResourceNames.VMName = *v
		}

		if v := item.NetworkInterfaces; v != nil {
			virtualMachineFullResourceNames.NetworkInterfaceNames = flattenNetworkInterfaceResourceNames(v)
		}

		if v := item.DataDiskNames; v != nil {
			virtualMachineFullResourceNames.DataDiskNames = flattenDataDiskNames(v)
		}

		result = append(result, virtualMachineFullResourceNames)
	}

	return result
}

func flattenCentralServerFullResourceNames(input *sapvirtualinstances.CentralServerFullResourceNames) []CentralServerFullResourceNames {
	if input == nil {
		return nil
	}

	result := CentralServerFullResourceNames{}

	if v := input.AvailabilitySetName; v != nil {
		result.AvailabilitySetName = *v
	}

	if v := input.LoadBalancer; v != nil {
		result.LoadBalancer = flattenLoadBalancerFullResourceNames(v)
	}

	if v := input.VirtualMachines; v != nil {
		result.VirtualMachines = flattenVirtualMachinesFullResourceNames(v)
	}

	return []CentralServerFullResourceNames{
		result,
	}
}

func flattenLoadBalancerFullResourceNames(input *sapvirtualinstances.LoadBalancerResourceNames) []LoadBalancer {
	if input == nil {
		return nil
	}

	result := LoadBalancer{}

	if v := input.LoadBalancerName; v != nil {
		result.Name = *v
	}

	if v := input.BackendPoolNames; v != nil {
		result.BackendPoolNames = *v
	}

	if v := input.FrontendIPConfigurationNames; v != nil {
		result.FrontendIpConfigurationNames = *v
	}

	if v := input.HealthProbeNames; v != nil {
		result.HealthProbeNames = *v
	}

	return []LoadBalancer{
		result,
	}
}

func flattenDatabaseServerFullResourceNames(input *sapvirtualinstances.DatabaseServerFullResourceNames) []DatabaseServerFullResourceNames {
	if input == nil {
		return nil
	}

	result := DatabaseServerFullResourceNames{}

	if v := input.AvailabilitySetName; v != nil {
		result.AvailabilitySetName = *v
	}

	if v := input.LoadBalancer; v != nil {
		result.LoadBalancer = flattenLoadBalancerFullResourceNames(v)
	}

	if v := input.VirtualMachines; v != nil {
		result.VirtualMachines = flattenVirtualMachinesFullResourceNames(v)
	}

	return []DatabaseServerFullResourceNames{
		result,
	}
}

func flattenSharedStorage(input *sapvirtualinstances.SharedStorageResourceNames) []SharedStorage {
	if input == nil {
		return nil
	}

	result := SharedStorage{}

	if v := input.SharedStorageAccountName; v != nil {
		result.AccountName = *v
	}

	if v := input.SharedStorageAccountPrivateEndPointName; v != nil {
		result.PrivateEndpointName = *v
	}

	return []SharedStorage{
		result,
	}
}
