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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/privateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	storageParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPDeploymentVirtualInstanceModel struct {
	Name                      string                       `tfschema:"name"`
	ResourceGroupName         string                       `tfschema:"resource_group_name"`
	Location                  string                       `tfschema:"location"`
	AppLocation               string                       `tfschema:"app_location"`
	OsSapConfiguration        []OsSapConfiguration         `tfschema:"os_sap_configuration"`
	SingleServerConfiguration []SingleServerConfiguration  `tfschema:"single_server_configuration"`
	ThreeTierConfiguration    []ThreeTierConfiguration     `tfschema:"three_tier_configuration"`
	Environment               string                       `tfschema:"environment"`
	Identity                  []identity.ModelUserAssigned `tfschema:"identity"`
	ManagedResourceGroupName  string                       `tfschema:"managed_resource_group_name"`
	SapProduct                string                       `tfschema:"sap_product"`
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
	AppResourceGroupName            string                            `tfschema:"app_resource_group_name"`
	DatabaseType                    string                            `tfschema:"database_type"`
	DiskVolumeConfigurations        []DiskVolumeConfiguration         `tfschema:"disk_volume_configuration"`
	IsSecondaryIpEnabled            bool                              `tfschema:"secondary_ip_enabled"`
	SubnetId                        string                            `tfschema:"subnet_id"`
	VirtualMachineConfiguration     []VirtualMachineConfiguration     `tfschema:"virtual_machine_configuration"`
	VirtualMachineFullResourceNames []VirtualMachineFullResourceNames `tfschema:"virtual_machine_full_resource_names"`
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

type VirtualMachineFullResourceNames struct {
	DataDiskNames         map[string]interface{} `tfschema:"data_disk_names"`
	HostName              string                 `tfschema:"host_name"`
	NetworkInterfaceNames []string               `tfschema:"network_interface_names"`
	OSDiskName            string                 `tfschema:"os_disk_name"`
	VMName                string                 `tfschema:"virtual_machine_name"`
}

type ThreeTierConfiguration struct {
	ApplicationServerConfiguration []ApplicationServerConfiguration `tfschema:"application_server_configuration"`
	AppResourceGroupName           string                           `tfschema:"app_resource_group_name"`
	CentralServerConfiguration     []CentralServerConfiguration     `tfschema:"central_server_configuration"`
	DatabaseServerConfiguration    []DatabaseServerConfiguration    `tfschema:"database_server_configuration"`
	FullResourceNames              []FullResourceNames              `tfschema:"full_resource_names"`
	HighAvailabilityType           string                           `tfschema:"high_availability_type"`
	IsSecondaryIpEnabled           bool                             `tfschema:"secondary_ip_enabled"`
	TransportCreateAndMount        []TransportCreateAndMount        `tfschema:"transport_create_and_mount"`
	TransportMount                 []TransportMount                 `tfschema:"transport_mount"`
}

type TransportCreateAndMount struct {
	ResourceGroupId    string `tfschema:"resource_group_id"`
	StorageAccountName string `tfschema:"storage_account_name"`
}

type TransportMount struct {
	FileShareId       string `tfschema:"file_share_id"`
	PrivateEndpointId string `tfschema:"private_endpoint_id"`
}

type ApplicationServerConfiguration struct {
	InstanceCount               int                           `tfschema:"instance_count"`
	SubnetId                    string                        `tfschema:"subnet_id"`
	VirtualMachineConfiguration []VirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
}

type CentralServerConfiguration struct {
	InstanceCount               int                           `tfschema:"instance_count"`
	SubnetId                    string                        `tfschema:"subnet_id"`
	VirtualMachineConfiguration []VirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
}

type DatabaseServerConfiguration struct {
	DatabaseType                string                        `tfschema:"database_type"`
	DiskVolumeConfigurations    []DiskVolumeConfiguration     `tfschema:"disk_volume_configuration"`
	InstanceCount               int                           `tfschema:"instance_count"`
	SubnetId                    string                        `tfschema:"subnet_id"`
	VirtualMachineConfiguration []VirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
}

type FullResourceNames struct {
	ApplicationServer []ApplicationServerFullResourceNames `tfschema:"application_server"`
	CentralServer     []CentralServerFullResourceNames     `tfschema:"central_server"`
	DatabaseServer    []DatabaseServerFullResourceNames    `tfschema:"database_server"`
	SharedStorage     []SharedStorage                      `tfschema:"shared_storage"`
}

type ApplicationServerFullResourceNames struct {
	AvailabilitySetName string                            `tfschema:"availability_set_name"`
	VirtualMachines     []VirtualMachineFullResourceNames `tfschema:"virtual_machine"`
}

type CentralServerFullResourceNames struct {
	AvailabilitySetName string                            `tfschema:"availability_set_name"`
	LoadBalancer        []LoadBalancer                    `tfschema:"load_balancer"`
	VirtualMachines     []VirtualMachineFullResourceNames `tfschema:"virtual_machine"`
}

type DatabaseServerFullResourceNames struct {
	AvailabilitySetName string                            `tfschema:"availability_set_name"`
	LoadBalancer        []LoadBalancer                    `tfschema:"load_balancer"`
	VirtualMachines     []VirtualMachineFullResourceNames `tfschema:"virtual_machine"`
}

type LoadBalancer struct {
	BackendPoolNames             []string `tfschema:"backend_pool_names"`
	FrontendIpConfigurationNames []string `tfschema:"frontend_ip_configuration_names"`
	HealthProbeNames             []string `tfschema:"health_probe_names"`
	Name                         string   `tfschema:"name"`
}

type SharedStorage struct {
	AccountName         string `tfschema:"account_name"`
	PrivateEndpointName string `tfschema:"private_endpoint_name"`
}

type WorkloadsSAPDeploymentVirtualInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPDeploymentVirtualInstanceResource{}

func (r WorkloadsSAPDeploymentVirtualInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_deployment_virtual_instance"
}

func (r WorkloadsSAPDeploymentVirtualInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPDeploymentVirtualInstanceModel{}
}

func (r WorkloadsSAPDeploymentVirtualInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapvirtualinstances.ValidateSapVirtualInstanceID
}

func (r WorkloadsSAPDeploymentVirtualInstanceResource) Arguments() map[string]*pluginsdk.Schema {
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
						ValidateFunc: validate.SAPFQDN,
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
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			AtLeastOneOf: []string{
				"single_server_configuration",
				"three_tier_configuration",
			},
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
												ValidateFunc: validate.AdminUsername,
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

								"virtual_machine_name": {
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
		},

		"three_tier_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			AtLeastOneOf: []string{
				"single_server_configuration",
				"three_tier_configuration",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"app_resource_group_name": commonschema.ResourceGroupName(),

					"application_server_configuration": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"instance_count": {
									Type:     pluginsdk.TypeInt,
									Required: true,
									ForceNew: true,
								},

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
															ValidateFunc: validate.AdminUsername,
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
							},
						},
					},

					"central_server_configuration": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"instance_count": {
									Type:     pluginsdk.TypeInt,
									Required: true,
									ForceNew: true,
								},

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
															ValidateFunc: validate.AdminUsername,
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
							},
						},
					},

					"database_server_configuration": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"instance_count": {
									Type:     pluginsdk.TypeInt,
									Required: true,
									ForceNew: true,
								},

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
															ValidateFunc: validate.AdminUsername,
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
												ValidateFunc: validate.AvailabilitySetName,
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

														"virtual_machine_name": {
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
												ValidateFunc: validate.AvailabilitySetName,
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

														"virtual_machine_name": {
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
												ValidateFunc: validate.AvailabilitySetName,
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

														"virtual_machine_name": {
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
												ValidateFunc: storageValidate.StorageAccountName,
											},

											"private_endpoint_name": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ForceNew:     true,
												ValidateFunc: networkValidate.PrivateLinkName,
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
						ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForSAPHighAvailabilityType(), false),
					},

					"secondary_ip_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"transport_create_and_mount": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"resource_group_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: commonids.ValidateResourceGroupID,
								},

								"storage_account_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: storageValidate.StorageAccountName,
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
									Required:     true,
									ForceNew:     true,
									ValidateFunc: storageValidate.StorageShareResourceManagerID,
								},

								"private_endpoint_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: privateendpoints.ValidatePrivateEndpointID,
								},
							},
						},
					},
				},
			},
		},

		"environment": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sapvirtualinstances.SAPEnvironmentTypeNonProd),
				string(sapvirtualinstances.SAPEnvironmentTypeProd),
			}, false),
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"sap_product": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sapvirtualinstances.SAPProductTypeECC),
				string(sapvirtualinstances.SAPProductTypeOther),
				string(sapvirtualinstances.SAPProductTypeSFourHANA),
			}, false),
		},

		"managed_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: resourcegroups.ValidateName,
		},

		"tags": commonschema.Tags(),
	}
}

func (r WorkloadsSAPDeploymentVirtualInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSAPDeploymentVirtualInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPDeploymentVirtualInstanceModel
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
					Environment: sapvirtualinstances.SAPEnvironmentType(model.Environment),
					SapProduct:  sapvirtualinstances.SAPProductType(model.SapProduct),
				},
				Tags: &model.Tags,
			}

			deploymentWithOSConfiguration := &sapvirtualinstances.DeploymentWithOSConfiguration{
				AppLocation:        utils.String(model.AppLocation),
				OsSapConfiguration: expandOsSapConfiguration(model.OsSapConfiguration),
			}

			if len(model.SingleServerConfiguration) != 0 {
				deploymentWithOSConfiguration.InfrastructureConfiguration = expandSingleServerConfiguration(model.SingleServerConfiguration)
			}

			if len(model.ThreeTierConfiguration) != 0 {
				threeTierConfiguration, err := expandThreeTierConfiguration(model.ThreeTierConfiguration)
				if err != nil {
					return err
				}
				deploymentWithOSConfiguration.InfrastructureConfiguration = threeTierConfiguration
			}

			parameters.Properties.Configuration = deploymentWithOSConfiguration

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

func (r WorkloadsSAPDeploymentVirtualInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPDeploymentVirtualInstanceModel
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

func (r WorkloadsSAPDeploymentVirtualInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances
			subscriptionId := metadata.Client.Account.SubscriptionId

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

			state := WorkloadsSAPDeploymentVirtualInstanceModel{}
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
						state.AppLocation = pointer.From(v.AppLocation)
						state.OsSapConfiguration = flattenOsSapConfiguration(v.OsSapConfiguration)

						if configuration := v.InfrastructureConfiguration; configuration != nil {
							if singleServerConfiguration, singleServerConfigurationExists := configuration.(sapvirtualinstances.SingleServerConfiguration); singleServerConfigurationExists {
								state.SingleServerConfiguration = flattenSingleServerConfiguration(singleServerConfiguration, metadata.ResourceData)
							}

							if threeTierConfiguration, threeTierConfigurationExists := configuration.(sapvirtualinstances.ThreeTierConfiguration); threeTierConfigurationExists {
								threeTierConfig, err := flattenThreeTierConfiguration(threeTierConfiguration, metadata.ResourceData, subscriptionId)
								if err != nil {
									return err
								}
								state.ThreeTierConfiguration = threeTierConfig
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

func (r WorkloadsSAPDeploymentVirtualInstanceResource) Delete() sdk.ResourceFunc {
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

func expandVirtualMachineFullResourceNames(input []VirtualMachineFullResourceNames) *sapvirtualinstances.SingleServerFullResourceNames {
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

func expandApplicationServer(input []ApplicationServerConfiguration) *sapvirtualinstances.ApplicationServerConfiguration {
	if len(input) == 0 {
		return nil
	}

	applicationServer := input[0]

	result := &sapvirtualinstances.ApplicationServerConfiguration{
		InstanceCount:               int64(applicationServer.InstanceCount),
		SubnetId:                    applicationServer.SubnetId,
		VirtualMachineConfiguration: pointer.From(expandVirtualMachineConfiguration(applicationServer.VirtualMachineConfiguration)),
	}

	return result
}

func expandCentralServer(input []CentralServerConfiguration) *sapvirtualinstances.CentralServerConfiguration {
	if len(input) == 0 {
		return nil
	}

	centralServer := input[0]

	result := &sapvirtualinstances.CentralServerConfiguration{
		InstanceCount:               int64(centralServer.InstanceCount),
		SubnetId:                    centralServer.SubnetId,
		VirtualMachineConfiguration: pointer.From(expandVirtualMachineConfiguration(centralServer.VirtualMachineConfiguration)),
	}

	return result
}

func expandDatabaseServer(input []DatabaseServerConfiguration) *sapvirtualinstances.DatabaseConfiguration {
	if len(input) == 0 {
		return nil
	}

	databaseServer := input[0]

	result := &sapvirtualinstances.DatabaseConfiguration{
		DiskConfiguration:           expandDiskVolumeConfigurations(databaseServer.DiskVolumeConfigurations),
		InstanceCount:               int64(databaseServer.InstanceCount),
		SubnetId:                    databaseServer.SubnetId,
		VirtualMachineConfiguration: pointer.From(expandVirtualMachineConfiguration(databaseServer.VirtualMachineConfiguration)),
	}

	if v := databaseServer.DatabaseType; v != "" {
		dbType := sapvirtualinstances.SAPDatabaseType(v)
		result.DatabaseType = &dbType
	}

	return result
}

func expandStorageConfiguration(input ThreeTierConfiguration) (*sapvirtualinstances.StorageConfiguration, error) {
	if len(input.TransportCreateAndMount) == 0 && len(input.TransportMount) == 0 {
		return &sapvirtualinstances.StorageConfiguration{
			TransportFileShareConfiguration: sapvirtualinstances.SkipFileShareConfiguration{},
		}, nil
	}

	result := &sapvirtualinstances.StorageConfiguration{}

	if len(input.TransportMount) != 0 {
		transportMount, err := expandTransportMount(input.TransportMount)
		if err != nil {
			return nil, err
		}
		result.TransportFileShareConfiguration = transportMount
	}

	if len(input.TransportCreateAndMount) != 0 {
		transportCreateAndMount, err := expandTransportCreateAndMount(input.TransportCreateAndMount)
		if err != nil {
			return nil, err
		}
		result.TransportFileShareConfiguration = transportCreateAndMount
	}

	return result, nil
}

func expandTransportCreateAndMount(input []TransportCreateAndMount) (*sapvirtualinstances.CreateAndMountFileShareConfiguration, error) {
	if len(input) == 0 {
		return nil, nil
	}

	transportCreateAndMount := input[0]

	result := &sapvirtualinstances.CreateAndMountFileShareConfiguration{}

	if v := transportCreateAndMount.ResourceGroupId; v != "" {
		resourceGroupId, err := commonids.ParseResourceGroupID(v)
		if err != nil {
			return nil, err
		}
		result.ResourceGroup = utils.String(resourceGroupId.ResourceGroupName)
	}

	if v := transportCreateAndMount.StorageAccountName; v != "" {
		result.StorageAccountName = utils.String(v)
	}

	return result, nil
}

func expandTransportMount(input []TransportMount) (*sapvirtualinstances.MountFileShareConfiguration, error) {
	if len(input) == 0 {
		return nil, nil
	}

	transportMount := input[0]

	// Currently, the last segment of the Storage File Share resource manager ID in Swagger is defined as `/shares/` but it's unexpected.
	// The last segment of the Storage File Share resource manager ID should be `/fileshares/` not `/shares/` in Swagger since the backend service is using `/fileshares/`.
	// See more details from https://github.com/Azure/azure-rest-api-specs/issues/25209
	storageShareResourceManagerId, err := storageParse.StorageShareResourceManagerID(transportMount.FileShareId)
	if err != nil {
		return nil, err
	}
	result := &sapvirtualinstances.MountFileShareConfiguration{
		Id:                storageParse.NewLegacyStorageShareResourceManagerID(storageShareResourceManagerId.SubscriptionId, storageShareResourceManagerId.ResourceGroup, storageShareResourceManagerId.StorageAccountName, storageShareResourceManagerId.FileServiceName, storageShareResourceManagerId.FileshareName).ID(),
		PrivateEndpointId: transportMount.PrivateEndpointId,
	}

	return result, nil
}

func expandFullResourceNames(input []FullResourceNames) *sapvirtualinstances.ThreeTierFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	fullResourceNames := input[0]

	result := &sapvirtualinstances.ThreeTierFullResourceNames{
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

	applicationServerFullResourceNames := input[0]

	result := &sapvirtualinstances.ApplicationServerFullResourceNames{
		VirtualMachines: expandVirtualMachinesFullResourceNames(applicationServerFullResourceNames.VirtualMachines),
	}

	if v := applicationServerFullResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = utils.String(v)
	}

	return result
}

func expandVirtualMachinesFullResourceNames(input []VirtualMachineFullResourceNames) *[]sapvirtualinstances.VirtualMachineResourceNames {
	result := make([]sapvirtualinstances.VirtualMachineResourceNames, 0)
	if len(input) == 0 {
		return &result
	}

	for _, item := range input {
		vmResourceNames := sapvirtualinstances.VirtualMachineResourceNames{
			DataDiskNames:     expandDataDiskNames(item.DataDiskNames),
			NetworkInterfaces: expandNetworkInterfaceNames(item.NetworkInterfaceNames),
		}

		if v := item.HostName; v != "" {
			vmResourceNames.HostName = utils.String(v)
		}

		if v := item.OSDiskName; v != "" {
			vmResourceNames.OsDiskName = utils.String(v)
		}

		if v := item.VMName; v != "" {
			vmResourceNames.VirtualMachineName = utils.String(v)
		}

		result = append(result, vmResourceNames)
	}

	return &result
}

func expandCentralServerFullResourceNames(input []CentralServerFullResourceNames) *sapvirtualinstances.CentralServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	centralServerFullResourceNames := input[0]

	result := &sapvirtualinstances.CentralServerFullResourceNames{
		LoadBalancer:    expandLoadBalancerFullResourceNames(centralServerFullResourceNames.LoadBalancer),
		VirtualMachines: expandVirtualMachinesFullResourceNames(centralServerFullResourceNames.VirtualMachines),
	}

	if v := centralServerFullResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = utils.String(v)
	}

	return result
}

func expandLoadBalancerFullResourceNames(input []LoadBalancer) *sapvirtualinstances.LoadBalancerResourceNames {
	result := &sapvirtualinstances.LoadBalancerResourceNames{}
	if len(input) == 0 {
		return result
	}

	loadBalancerFullResourceNames := input[0]

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

	return result
}

func expandDatabaseServerFullResourceNames(input []DatabaseServerFullResourceNames) *sapvirtualinstances.DatabaseServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	databaseServerFullResourceNames := input[0]

	result := &sapvirtualinstances.DatabaseServerFullResourceNames{
		LoadBalancer:    expandLoadBalancerFullResourceNames(databaseServerFullResourceNames.LoadBalancer),
		VirtualMachines: expandVirtualMachinesFullResourceNames(databaseServerFullResourceNames.VirtualMachines),
	}

	if v := databaseServerFullResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = utils.String(v)
	}

	return result
}

func expandSharedStorage(input []SharedStorage) *sapvirtualinstances.SharedStorageResourceNames {
	result := &sapvirtualinstances.SharedStorageResourceNames{}
	if len(input) == 0 {
		return result
	}

	sharedStorage := input[0]

	if v := sharedStorage.AccountName; v != "" {
		result.SharedStorageAccountName = utils.String(v)
	}

	if v := sharedStorage.PrivateEndpointName; v != "" {
		result.SharedStorageAccountPrivateEndPointName = utils.String(v)
	}

	return result
}

func expandSingleServerConfiguration(input []SingleServerConfiguration) *sapvirtualinstances.SingleServerConfiguration {
	singleServerConfiguration := input[0]

	result := &sapvirtualinstances.SingleServerConfiguration{
		AppResourceGroup:    singleServerConfiguration.AppResourceGroupName,
		CustomResourceNames: expandVirtualMachineFullResourceNames(singleServerConfiguration.VirtualMachineFullResourceNames),
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

func expandThreeTierConfiguration(input []ThreeTierConfiguration) (*sapvirtualinstances.ThreeTierConfiguration, error) {
	threeTierConfiguration := input[0]

	result := &sapvirtualinstances.ThreeTierConfiguration{
		ApplicationServer:   pointer.From(expandApplicationServer(threeTierConfiguration.ApplicationServerConfiguration)),
		AppResourceGroup:    threeTierConfiguration.AppResourceGroupName,
		CentralServer:       pointer.From(expandCentralServer(threeTierConfiguration.CentralServerConfiguration)),
		CustomResourceNames: expandFullResourceNames(threeTierConfiguration.FullResourceNames),
		DatabaseServer:      pointer.From(expandDatabaseServer(threeTierConfiguration.DatabaseServerConfiguration)),
		NetworkConfiguration: &sapvirtualinstances.NetworkConfiguration{
			IsSecondaryIPEnabled: utils.Bool(threeTierConfiguration.IsSecondaryIpEnabled),
		},
	}

	storageConfiguration, err := expandStorageConfiguration(threeTierConfiguration)
	if err != nil {
		return nil, err
	}
	result.StorageConfiguration = storageConfiguration

	if v := threeTierConfiguration.HighAvailabilityType; v != "" {
		result.HighAvailabilityConfig = &sapvirtualinstances.HighAvailabilityConfiguration{
			HighAvailabilityType: sapvirtualinstances.SAPHighAvailabilityType(v),
		}
	}

	return result, nil
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

func flattenApplicationServer(input sapvirtualinstances.ApplicationServerConfiguration, d *pluginsdk.ResourceData, basePath string) []ApplicationServerConfiguration {
	result := make([]ApplicationServerConfiguration, 0)

	applicationServerConfig := ApplicationServerConfiguration{
		InstanceCount:               int(input.InstanceCount),
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.application_server_configuration", basePath)),
	}

	return append(result, applicationServerConfig)
}

func flattenCentralServer(input sapvirtualinstances.CentralServerConfiguration, d *pluginsdk.ResourceData, basePath string) []CentralServerConfiguration {
	result := make([]CentralServerConfiguration, 0)

	centralServerConfig := CentralServerConfiguration{
		InstanceCount:               int(input.InstanceCount),
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.central_server_configuration", basePath)),
	}

	return append(result, centralServerConfig)
}

func flattenDatabaseServer(input sapvirtualinstances.DatabaseConfiguration, d *pluginsdk.ResourceData, basePath string) []DatabaseServerConfiguration {
	result := make([]DatabaseServerConfiguration, 0)

	databaseServerConfig := DatabaseServerConfiguration{
		DiskVolumeConfigurations:    flattenDiskVolumeConfigurations(input.DiskConfiguration),
		InstanceCount:               int(input.InstanceCount),
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.database_server_configuration", basePath)),
	}

	if v := input.DatabaseType; v != nil {
		databaseServerConfig.DatabaseType = string(*v)
	}

	return append(result, databaseServerConfig)
}

func flattenFullResourceNames(input sapvirtualinstances.ThreeTierFullResourceNames) []FullResourceNames {
	result := make([]FullResourceNames, 0)

	fullResourceNames := FullResourceNames{
		ApplicationServer: flattenApplicationServerFullResourceNames(input.ApplicationServer),
		CentralServer:     flattenCentralServerFullResourceNames(input.CentralServer),
		DatabaseServer:    flattenDatabaseServerFullResourceNames(input.DatabaseServer),
		SharedStorage:     flattenSharedStorage(input.SharedStorage),
	}

	return append(result, fullResourceNames)
}

func flattenApplicationServerFullResourceNames(input *sapvirtualinstances.ApplicationServerFullResourceNames) []ApplicationServerFullResourceNames {
	result := make([]ApplicationServerFullResourceNames, 0)
	if input == nil {
		return result
	}

	return append(result, ApplicationServerFullResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		VirtualMachines:     flattenVirtualMachinesFullResourceNames(input.VirtualMachines),
	})
}

func flattenVirtualMachinesFullResourceNames(input *[]sapvirtualinstances.VirtualMachineResourceNames) []VirtualMachineFullResourceNames {
	result := make([]VirtualMachineFullResourceNames, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		result = append(result, VirtualMachineFullResourceNames{
			HostName:              pointer.From(item.HostName),
			OSDiskName:            pointer.From(item.OsDiskName),
			VMName:                pointer.From(item.VirtualMachineName),
			DataDiskNames:         flattenDataDiskNames(item.DataDiskNames),
			NetworkInterfaceNames: flattenNetworkInterfaceResourceNames(item.NetworkInterfaces),
		})
	}

	return result
}

func flattenCentralServerFullResourceNames(input *sapvirtualinstances.CentralServerFullResourceNames) []CentralServerFullResourceNames {
	result := make([]CentralServerFullResourceNames, 0)
	if input == nil {
		return result
	}

	centralServerFullResourceNames := CentralServerFullResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		LoadBalancer:        flattenLoadBalancerFullResourceNames(input.LoadBalancer),
		VirtualMachines:     flattenVirtualMachinesFullResourceNames(input.VirtualMachines),
	}

	return append(result, centralServerFullResourceNames)
}

func flattenLoadBalancerFullResourceNames(input *sapvirtualinstances.LoadBalancerResourceNames) []LoadBalancer {
	result := make([]LoadBalancer, 0)
	if input == nil {
		return result
	}

	return append(result, LoadBalancer{
		Name:                         pointer.From(input.LoadBalancerName),
		BackendPoolNames:             pointer.From(input.BackendPoolNames),
		FrontendIpConfigurationNames: pointer.From(input.FrontendIPConfigurationNames),
		HealthProbeNames:             pointer.From(input.HealthProbeNames),
	})
}

func flattenDatabaseServerFullResourceNames(input *sapvirtualinstances.DatabaseServerFullResourceNames) []DatabaseServerFullResourceNames {
	result := make([]DatabaseServerFullResourceNames, 0)
	if input == nil {
		return result
	}

	return append(result, DatabaseServerFullResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		LoadBalancer:        flattenLoadBalancerFullResourceNames(input.LoadBalancer),
		VirtualMachines:     flattenVirtualMachinesFullResourceNames(input.VirtualMachines),
	})
}

func flattenSharedStorage(input *sapvirtualinstances.SharedStorageResourceNames) []SharedStorage {
	result := make([]SharedStorage, 0)
	if input == nil {
		return result
	}

	return append(result, SharedStorage{
		AccountName:         pointer.From(input.SharedStorageAccountName),
		PrivateEndpointName: pointer.From(input.SharedStorageAccountPrivateEndPointName),
	})
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

func flattenVirtualMachineFullResourceNames(input sapvirtualinstances.SingleServerFullResourceNames) []VirtualMachineFullResourceNames {
	result := make([]VirtualMachineFullResourceNames, 0)
	vmFullResourceNames := VirtualMachineFullResourceNames{}

	if vm := input.VirtualMachine; vm != nil {
		vmFullResourceNames.HostName = pointer.From(vm.HostName)
		vmFullResourceNames.OSDiskName = pointer.From(vm.OsDiskName)
		vmFullResourceNames.VMName = pointer.From(vm.VirtualMachineName)
		vmFullResourceNames.NetworkInterfaceNames = flattenNetworkInterfaceResourceNames(vm.NetworkInterfaces)
		vmFullResourceNames.DataDiskNames = flattenDataDiskNames(vm.DataDiskNames)
	}

	return append(result, vmFullResourceNames)
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

func flattenVirtualMachineConfiguration(input sapvirtualinstances.VirtualMachineConfiguration, d *pluginsdk.ResourceData, basePath string) []VirtualMachineConfiguration {
	result := make([]VirtualMachineConfiguration, 0)

	return append(result, VirtualMachineConfiguration{
		ImageReference: flattenImageReference(input.ImageReference),
		OSProfile:      flattenOSProfile(input.OsProfile, d, fmt.Sprintf("%s.0.virtual_machine_configuration", basePath)),
		VmSize:         input.VMSize,
	})
}

func flattenImageReference(input sapvirtualinstances.ImageReference) []ImageReference {
	result := make([]ImageReference, 0)

	imageReference := ImageReference{
		Offer:     pointer.From(input.Offer),
		Publisher: pointer.From(input.Publisher),
		Sku:       pointer.From(input.Sku),
		Version:   pointer.From(input.Version),
	}

	return append(result, imageReference)
}

func flattenOSProfile(input sapvirtualinstances.OSProfile, d *pluginsdk.ResourceData, basePath string) []OSProfile {
	result := make([]OSProfile, 0)

	osProfile := OSProfile{
		AdminUsername: pointer.From(input.AdminUsername),
	}

	if osConfiguration := input.OsConfiguration; osConfiguration != nil {
		if v, ok := osConfiguration.(sapvirtualinstances.LinuxConfiguration); ok {
			if sshKeyPair := v.SshKeyPair; sshKeyPair != nil {
				osProfile.SshPrivateKey = d.Get(fmt.Sprintf("%s.0.os_profile.0.ssh_private_key", basePath)).(string)
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
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, "single_server_configuration"),
	}

	if networkConfiguration := input.NetworkConfiguration; networkConfiguration != nil {
		singleServerConfig.IsSecondaryIpEnabled = pointer.From(networkConfiguration.IsSecondaryIPEnabled)
	}

	if customResourceNames := input.CustomResourceNames; customResourceNames != nil {
		if v, ok := customResourceNames.(sapvirtualinstances.SingleServerFullResourceNames); ok {
			singleServerConfig.VirtualMachineFullResourceNames = flattenVirtualMachineFullResourceNames(v)
		}
	}

	return append(result, singleServerConfig)
}

func flattenThreeTierConfiguration(input sapvirtualinstances.ThreeTierConfiguration, d *pluginsdk.ResourceData, subscriptionId string) ([]ThreeTierConfiguration, error) {
	result := make([]ThreeTierConfiguration, 0)

	threeTierConfig := ThreeTierConfiguration{
		AppResourceGroupName:           input.AppResourceGroup,
		ApplicationServerConfiguration: flattenApplicationServer(input.ApplicationServer, d, "three_tier_configuration"),
		CentralServerConfiguration:     flattenCentralServer(input.CentralServer, d, "three_tier_configuration"),
		DatabaseServerConfiguration:    flattenDatabaseServer(input.DatabaseServer, d, "three_tier_configuration"),
	}

	if customResourceNames := input.CustomResourceNames; customResourceNames != nil {
		if v, ok := customResourceNames.(sapvirtualinstances.ThreeTierFullResourceNames); ok {
			threeTierConfig.FullResourceNames = flattenFullResourceNames(v)
		}
	}

	if v := input.HighAvailabilityConfig; v != nil && v.HighAvailabilityType != "" {
		threeTierConfig.HighAvailabilityType = string(v.HighAvailabilityType)
	}

	if v := input.NetworkConfiguration; v != nil && v.IsSecondaryIPEnabled != nil {
		threeTierConfig.IsSecondaryIpEnabled = *v.IsSecondaryIPEnabled
	}

	if storageConfiguration := input.StorageConfiguration; storageConfiguration != nil {
		if transportFileShareConfiguration := storageConfiguration.TransportFileShareConfiguration; transportFileShareConfiguration != nil {
			if createAndMountFileShareConfiguration, ok := transportFileShareConfiguration.(sapvirtualinstances.CreateAndMountFileShareConfiguration); ok {
				transportCreateAndMount := TransportCreateAndMount{
					StorageAccountName: pointer.From(createAndMountFileShareConfiguration.StorageAccountName),
				}

				var resourceGroupId string
				if v := createAndMountFileShareConfiguration.ResourceGroup; v != nil {
					resourceGroupId = commonids.NewResourceGroupID(subscriptionId, *createAndMountFileShareConfiguration.ResourceGroup).ID()
				}
				transportCreateAndMount.ResourceGroupId = resourceGroupId

				threeTierConfig.TransportCreateAndMount = []TransportCreateAndMount{
					transportCreateAndMount,
				}
			}

			if v, ok := transportFileShareConfiguration.(sapvirtualinstances.MountFileShareConfiguration); ok {
				transportMount, err := flattenTransportMount(v)
				if err != nil {
					return nil, err
				}
				threeTierConfig.TransportMount = transportMount
			}
		}
	}

	return append(result, threeTierConfig), nil
}

func flattenTransportMount(input sapvirtualinstances.MountFileShareConfiguration) ([]TransportMount, error) {
	result := make([]TransportMount, 0)

	// Currently, the last segment of the Storage File Share resource manager ID in Swagger is defined as `/shares/` but it's unexpected.
	// The last segment of the Storage File Share resource manager ID should be `/fileshares/` not `/shares/` in Swagger since the backend service is using `/fileshares/`.
	// See more details from https://github.com/Azure/azure-rest-api-specs/issues/25209
	legacyStorageShareResourceManagerId, err := storageParse.LegacyStorageShareResourceManagerID(input.Id)
	if err != nil {
		return nil, err
	}
	return append(result, TransportMount{
		FileShareId:       storageParse.NewStorageShareResourceManagerID(legacyStorageShareResourceManagerId.SubscriptionId, legacyStorageShareResourceManagerId.ResourceGroup, legacyStorageShareResourceManagerId.StorageAccountName, legacyStorageShareResourceManagerId.FileServiceName, legacyStorageShareResourceManagerId.ShareName).ID(),
		PrivateEndpointId: input.PrivateEndpointId,
	}), nil
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
