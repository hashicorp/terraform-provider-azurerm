// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workloads

import (
	"context"
	"fmt"
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
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPThreeTierVirtualInstanceModel struct {
	Name                     string                       `tfschema:"name"`
	ResourceGroupName        string                       `tfschema:"resource_group_name"`
	Location                 string                       `tfschema:"location"`
	AppLocation              string                       `tfschema:"app_location"`
	Environment              string                       `tfschema:"environment"`
	Identity                 []identity.ModelUserAssigned `tfschema:"identity"`
	ManagedResourceGroupName string                       `tfschema:"managed_resource_group_name"`
	SapFqdn                  string                       `tfschema:"sap_fqdn"`
	SapProduct               string                       `tfschema:"sap_product"`
	ThreeTierConfiguration   []ThreeTierConfiguration     `tfschema:"three_tier_configuration"`
	Tags                     map[string]string            `tfschema:"tags"`
}

type DiskVolumeConfiguration struct {
	VolumeName    string `tfschema:"volume_name"`
	NumberOfDisks int64  `tfschema:"number_of_disks"`
	SizeGb        int64  `tfschema:"size_in_gb"`
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
	DataDisks             []DataDisk `tfschema:"data_disk"`
	HostName              string     `tfschema:"host_name"`
	NetworkInterfaceNames []string   `tfschema:"network_interface_names"`
	OSDiskName            string     `tfschema:"os_disk_name"`
	VMName                string     `tfschema:"virtual_machine_name"`
}

type DataDisk struct {
	VolumeName string   `tfschema:"volume_name"`
	Names      []string `tfschema:"names"`
}

type ThreeTierConfiguration struct {
	ApplicationServerConfiguration []ApplicationServerConfiguration `tfschema:"application_server_configuration"`
	AppResourceGroupName           string                           `tfschema:"app_resource_group_name"`
	CentralServerConfiguration     []CentralServerConfiguration     `tfschema:"central_server_configuration"`
	DatabaseServerConfiguration    []DatabaseServerConfiguration    `tfschema:"database_server_configuration"`
	ResourceNames                  []ResourceNames                  `tfschema:"resource_names"`
	HighAvailabilityType           string                           `tfschema:"high_availability_type"`
	IsSecondaryIpEnabled           bool                             `tfschema:"secondary_ip_enabled"`
	TransportCreateAndMount        []TransportCreateAndMount        `tfschema:"transport_create_and_mount"`
}

type TransportCreateAndMount struct {
	ResourceGroupId    string `tfschema:"resource_group_id"`
	StorageAccountName string `tfschema:"storage_account_name"`
}

type ApplicationServerConfiguration struct {
	InstanceCount               int64                         `tfschema:"instance_count"`
	SubnetId                    string                        `tfschema:"subnet_id"`
	VirtualMachineConfiguration []VirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
}

type CentralServerConfiguration struct {
	InstanceCount               int64                         `tfschema:"instance_count"`
	SubnetId                    string                        `tfschema:"subnet_id"`
	VirtualMachineConfiguration []VirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
}

type DatabaseServerConfiguration struct {
	DatabaseType                string                        `tfschema:"database_type"`
	DiskVolumeConfigurations    []DiskVolumeConfiguration     `tfschema:"disk_volume_configuration"`
	InstanceCount               int64                         `tfschema:"instance_count"`
	SubnetId                    string                        `tfschema:"subnet_id"`
	VirtualMachineConfiguration []VirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
}

type ResourceNames struct {
	ApplicationServer []ApplicationServerResourceNames `tfschema:"application_server"`
	CentralServer     []CentralServerResourceNames     `tfschema:"central_server"`
	DatabaseServer    []DatabaseServerResourceNames    `tfschema:"database_server"`
	SharedStorage     []SharedStorage                  `tfschema:"shared_storage"`
}

type ApplicationServerResourceNames struct {
	AvailabilitySetName string                        `tfschema:"availability_set_name"`
	VirtualMachines     []VirtualMachineResourceNames `tfschema:"virtual_machine"`
}

type CentralServerResourceNames struct {
	AvailabilitySetName string                        `tfschema:"availability_set_name"`
	LoadBalancer        []LoadBalancer                `tfschema:"load_balancer"`
	VirtualMachines     []VirtualMachineResourceNames `tfschema:"virtual_machine"`
}

type DatabaseServerResourceNames struct {
	AvailabilitySetName string                        `tfschema:"availability_set_name"`
	LoadBalancer        []LoadBalancer                `tfschema:"load_balancer"`
	VirtualMachines     []VirtualMachineResourceNames `tfschema:"virtual_machine"`
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

type WorkloadsSAPThreeTierVirtualInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPThreeTierVirtualInstanceResource{}

func (r WorkloadsSAPThreeTierVirtualInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_three_tier_virtual_instance"
}

func (r WorkloadsSAPThreeTierVirtualInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPThreeTierVirtualInstanceModel{}
}

func (r WorkloadsSAPThreeTierVirtualInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapvirtualinstances.ValidateSapVirtualInstanceID
}

func (r WorkloadsSAPThreeTierVirtualInstanceResource) Arguments() map[string]*pluginsdk.Schema {
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

		"environment": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForSAPEnvironmentType(), false),
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"sap_fqdn": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(2, 34),
		},

		"sap_product": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForSAPProductType(), false),
		},

		"three_tier_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
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
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntAtLeast(1),
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
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntAtLeast(1),
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
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntAtLeast(1),
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
							},
						},
					},

					"resource_names": {
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
												ValidateFunc: validation.StringLenBetween(1, 80),
											},

											"virtual_machine": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												ForceNew: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"data_disk": {
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
															ValidateFunc: validation.StringLenBetween(1, 13),
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
															ValidateFunc: validation.StringLenBetween(1, 80),
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
												ValidateFunc: validation.StringLenBetween(1, 80),
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
														"data_disk": {
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
															ValidateFunc: validation.StringLenBetween(1, 13),
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
															ValidateFunc: validation.StringLenBetween(1, 80),
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
												ValidateFunc: validation.StringLenBetween(1, 80),
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
														"data_disk": {
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
																			"hanaData",
																			"hanaLog",
																			"hanaShared",
																			"usrSap",
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
															ValidateFunc: validation.StringLenBetween(1, 13),
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
															ValidateFunc: validation.StringLenBetween(1, 80),
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
				},
			},
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

func (r WorkloadsSAPThreeTierVirtualInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSAPThreeTierVirtualInstanceResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if v := rd.Get("three_tier_configuration.0.database_server_configuration.0.disk_volume_configuration"); v != nil {
				diskVolumes := v.(*pluginsdk.Set).List()
				if hasDuplicateVolumeName(diskVolumes) {
					return fmt.Errorf("`volume_name` cannot be duplicated")
				}
			}

			return nil
		},
	}
}

func hasDuplicateVolumeName(input []interface{}) bool {
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

func (r WorkloadsSAPThreeTierVirtualInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPThreeTierVirtualInstanceModel
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
				AppLocation: utils.String(location.Normalize(model.AppLocation)),
				OsSapConfiguration: &sapvirtualinstances.OsSapConfiguration{
					SapFqdn: utils.String(model.SapFqdn),
				},
			}

			threeTierConfiguration, err := expandThreeTierConfiguration(model.ThreeTierConfiguration)
			if err != nil {
				return err
			}
			deploymentWithOSConfiguration.InfrastructureConfiguration = threeTierConfiguration

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

func (r WorkloadsSAPThreeTierVirtualInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPThreeTierVirtualInstanceModel
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

func (r WorkloadsSAPThreeTierVirtualInstanceResource) Read() sdk.ResourceFunc {
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

			state := WorkloadsSAPThreeTierVirtualInstanceModel{}
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
						state.AppLocation = location.Normalize(pointer.From(v.AppLocation))

						if osSapConfiguration := v.OsSapConfiguration; osSapConfiguration != nil {
							state.SapFqdn = pointer.From(osSapConfiguration.SapFqdn)
						}

						if configuration := v.InfrastructureConfiguration; configuration != nil {
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

func (r WorkloadsSAPThreeTierVirtualInstanceResource) Delete() sdk.ResourceFunc {
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

func expandDataDisks(input []DataDisk) *map[string][]string {
	result := make(map[string][]string)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		result[v.VolumeName] = v.Names
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

func expandApplicationServer(input []ApplicationServerConfiguration) *sapvirtualinstances.ApplicationServerConfiguration {
	if len(input) == 0 {
		return nil
	}

	applicationServer := input[0]

	result := &sapvirtualinstances.ApplicationServerConfiguration{
		InstanceCount:               applicationServer.InstanceCount,
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
		InstanceCount:               centralServer.InstanceCount,
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
		InstanceCount:               databaseServer.InstanceCount,
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
	if len(input.TransportCreateAndMount) == 0 {
		return &sapvirtualinstances.StorageConfiguration{
			TransportFileShareConfiguration: sapvirtualinstances.SkipFileShareConfiguration{},
		}, nil
	}

	result := &sapvirtualinstances.StorageConfiguration{}

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

func expandResourceNames(input []ResourceNames) *sapvirtualinstances.ThreeTierFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	resourceNames := input[0]

	result := &sapvirtualinstances.ThreeTierFullResourceNames{
		ApplicationServer: expandApplicationServerResourceNames(resourceNames.ApplicationServer),
		CentralServer:     expandCentralServerResourceNames(resourceNames.CentralServer),
		DatabaseServer:    expandDatabaseServerResourceNames(resourceNames.DatabaseServer),
		SharedStorage:     expandSharedStorage(resourceNames.SharedStorage),
	}

	return result
}

func expandApplicationServerResourceNames(input []ApplicationServerResourceNames) *sapvirtualinstances.ApplicationServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	applicationServerResourceNames := input[0]

	result := &sapvirtualinstances.ApplicationServerFullResourceNames{
		VirtualMachines: expandVirtualMachinesResourceNames(applicationServerResourceNames.VirtualMachines),
	}

	if v := applicationServerResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = utils.String(v)
	}

	return result
}

func expandVirtualMachinesResourceNames(input []VirtualMachineResourceNames) *[]sapvirtualinstances.VirtualMachineResourceNames {
	result := make([]sapvirtualinstances.VirtualMachineResourceNames, 0)
	if len(input) == 0 {
		return &result
	}

	for _, item := range input {
		vmResourceNames := sapvirtualinstances.VirtualMachineResourceNames{
			DataDiskNames:     expandDataDisks(item.DataDisks),
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

func expandCentralServerResourceNames(input []CentralServerResourceNames) *sapvirtualinstances.CentralServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	centralServerResourceNames := input[0]

	result := &sapvirtualinstances.CentralServerFullResourceNames{
		LoadBalancer:    expandLoadBalancerResourceNames(centralServerResourceNames.LoadBalancer),
		VirtualMachines: expandVirtualMachinesResourceNames(centralServerResourceNames.VirtualMachines),
	}

	if v := centralServerResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = utils.String(v)
	}

	return result
}

func expandLoadBalancerResourceNames(input []LoadBalancer) *sapvirtualinstances.LoadBalancerResourceNames {
	result := &sapvirtualinstances.LoadBalancerResourceNames{}
	if len(input) == 0 {
		return result
	}

	loadBalancerResourceNames := input[0]

	if v := loadBalancerResourceNames.Name; v != "" {
		result.LoadBalancerName = utils.String(v)
	}

	if v := loadBalancerResourceNames.BackendPoolNames; v != nil {
		result.BackendPoolNames = &v
	}

	if v := loadBalancerResourceNames.FrontendIpConfigurationNames; v != nil {
		result.FrontendIPConfigurationNames = &v
	}

	if v := loadBalancerResourceNames.HealthProbeNames; v != nil {
		result.HealthProbeNames = &v
	}

	return result
}

func expandDatabaseServerResourceNames(input []DatabaseServerResourceNames) *sapvirtualinstances.DatabaseServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	databaseServerResourceNames := input[0]

	result := &sapvirtualinstances.DatabaseServerFullResourceNames{
		LoadBalancer:    expandLoadBalancerResourceNames(databaseServerResourceNames.LoadBalancer),
		VirtualMachines: expandVirtualMachinesResourceNames(databaseServerResourceNames.VirtualMachines),
	}

	if v := databaseServerResourceNames.AvailabilitySetName; v != "" {
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

func expandThreeTierConfiguration(input []ThreeTierConfiguration) (*sapvirtualinstances.ThreeTierConfiguration, error) {
	threeTierConfiguration := input[0]

	result := &sapvirtualinstances.ThreeTierConfiguration{
		ApplicationServer:   pointer.From(expandApplicationServer(threeTierConfiguration.ApplicationServerConfiguration)),
		AppResourceGroup:    threeTierConfiguration.AppResourceGroupName,
		CentralServer:       pointer.From(expandCentralServer(threeTierConfiguration.CentralServerConfiguration)),
		CustomResourceNames: expandResourceNames(threeTierConfiguration.ResourceNames),
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

func flattenApplicationServer(input sapvirtualinstances.ApplicationServerConfiguration, d *pluginsdk.ResourceData, basePath string) []ApplicationServerConfiguration {
	result := make([]ApplicationServerConfiguration, 0)

	applicationServerConfig := ApplicationServerConfiguration{
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.application_server_configuration", basePath)),
	}

	return append(result, applicationServerConfig)
}

func flattenCentralServer(input sapvirtualinstances.CentralServerConfiguration, d *pluginsdk.ResourceData, basePath string) []CentralServerConfiguration {
	result := make([]CentralServerConfiguration, 0)

	centralServerConfig := CentralServerConfiguration{
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.central_server_configuration", basePath)),
	}

	return append(result, centralServerConfig)
}

func flattenDatabaseServer(input sapvirtualinstances.DatabaseConfiguration, d *pluginsdk.ResourceData, basePath string) []DatabaseServerConfiguration {
	result := make([]DatabaseServerConfiguration, 0)

	databaseServerConfig := DatabaseServerConfiguration{
		DiskVolumeConfigurations:    flattenDiskVolumeConfigurations(input.DiskConfiguration),
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.database_server_configuration", basePath)),
	}

	if v := input.DatabaseType; v != nil {
		databaseServerConfig.DatabaseType = string(*v)
	}

	return append(result, databaseServerConfig)
}

func flattenResourceNames(input sapvirtualinstances.ThreeTierFullResourceNames) []ResourceNames {
	result := make([]ResourceNames, 0)

	resourceNames := ResourceNames{
		ApplicationServer: flattenApplicationServerResourceNames(input.ApplicationServer),
		CentralServer:     flattenCentralServerResourceNames(input.CentralServer),
		DatabaseServer:    flattenDatabaseServerResourceNames(input.DatabaseServer),
		SharedStorage:     flattenSharedStorage(input.SharedStorage),
	}

	return append(result, resourceNames)
}

func flattenApplicationServerResourceNames(input *sapvirtualinstances.ApplicationServerFullResourceNames) []ApplicationServerResourceNames {
	result := make([]ApplicationServerResourceNames, 0)
	if input == nil {
		return result
	}

	return append(result, ApplicationServerResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		VirtualMachines:     flattenVirtualMachinesResourceNames(input.VirtualMachines),
	})
}

func flattenVirtualMachinesResourceNames(input *[]sapvirtualinstances.VirtualMachineResourceNames) []VirtualMachineResourceNames {
	result := make([]VirtualMachineResourceNames, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		result = append(result, VirtualMachineResourceNames{
			HostName:              pointer.From(item.HostName),
			OSDiskName:            pointer.From(item.OsDiskName),
			VMName:                pointer.From(item.VirtualMachineName),
			DataDisks:             flattenDataDisks(item.DataDiskNames),
			NetworkInterfaceNames: flattenNetworkInterfaceResourceNames(item.NetworkInterfaces),
		})
	}

	return result
}

func flattenCentralServerResourceNames(input *sapvirtualinstances.CentralServerFullResourceNames) []CentralServerResourceNames {
	result := make([]CentralServerResourceNames, 0)
	if input == nil {
		return result
	}

	centralServerResourceNames := CentralServerResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		LoadBalancer:        flattenLoadBalancerResourceNames(input.LoadBalancer),
		VirtualMachines:     flattenVirtualMachinesResourceNames(input.VirtualMachines),
	}

	return append(result, centralServerResourceNames)
}

func flattenLoadBalancerResourceNames(input *sapvirtualinstances.LoadBalancerResourceNames) []LoadBalancer {
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

func flattenDatabaseServerResourceNames(input *sapvirtualinstances.DatabaseServerFullResourceNames) []DatabaseServerResourceNames {
	result := make([]DatabaseServerResourceNames, 0)
	if input == nil {
		return result
	}

	return append(result, DatabaseServerResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		LoadBalancer:        flattenLoadBalancerResourceNames(input.LoadBalancer),
		VirtualMachines:     flattenVirtualMachinesResourceNames(input.VirtualMachines),
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

func flattenDataDisks(input *map[string][]string) []DataDisk {
	results := make([]DataDisk, 0)
	if input == nil {
		return results
	}

	for k, v := range *input {
		dataDisk := DataDisk{
			VolumeName: k,
			Names:      v,
		}

		results = append(results, dataDisk)
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

	return append(result, ImageReference{
		Offer:     pointer.From(input.Offer),
		Publisher: pointer.From(input.Publisher),
		Sku:       pointer.From(input.Sku),
		Version:   pointer.From(input.Version),
	})
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
			threeTierConfig.ResourceNames = flattenResourceNames(v)
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

			if _, ok := transportFileShareConfiguration.(sapvirtualinstances.MountFileShareConfiguration); ok {
				return nil, fmt.Errorf("currently, the last segment of the Storage File Share resource manager ID in Swagger is defined as `/shares/` but it's unexpected. The last segment of the Storage File Share resource manager ID should be `/fileshares/` not `/shares/` in Swagger since the backend service is using `/fileshares/`. See more details from https://github.com/Azure/azure-rest-api-specs/issues/25209. So the feature of `TransportMount` isn't supported by TF for now due to this service API bug")
			}
		}
	}

	return append(result, threeTierConfig), nil
}
