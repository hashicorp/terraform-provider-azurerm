package workloads

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
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

type WorkloadsSAPVirtualInstanceModel struct {
	Name                          string                          `tfschema:"name"`
	ResourceGroupName             string                          `tfschema:"resource_group_name"`
	Location                      string                          `tfschema:"location"`
	DeploymentWithOSConfiguration []DeploymentWithOSConfiguration `tfschema:"deployment_with_os_configuration"`
	DiscoveryConfiguration        []DiscoveryConfiguration        `tfschema:"discovery_configuration"`
	Environment                   string                          `tfschema:"environment"`
	Identity                      []identity.ModelUserAssigned    `tfschema:"identity"`
	ManagedResourceGroupName      string                          `tfschema:"managed_resource_group_name"`
	SapProduct                    string                          `tfschema:"sap_product"`
	Tags                          map[string]string               `tfschema:"tags"`
}

type DeploymentWithOSConfiguration struct {
	AppLocation               string                      `tfschema:"app_location"`
	OsSapConfiguration        []OsSapConfiguration        `tfschema:"os_sap_configuration"`
	SingleServerConfiguration []SingleServerConfiguration `tfschema:"single_server_configuration"`
	ThreeTierConfiguration    []ThreeTierConfiguration    `tfschema:"three_tier_configuration"`
}

type DiscoveryConfiguration struct {
	CentralServerVmId         string `tfschema:"central_server_vm_id"`
	ManagedStorageAccountName string `tfschema:"managed_storage_account_name"`
}

type OsSapConfiguration struct {
	DeployerVmPackages []DeployerVmPackages `tfschema:"deployer_vm_packages"`
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
	IsSecondaryIpEnabled            bool                              `tfschema:"is_secondary_ip_enabled"`
	SubnetId                        string                            `tfschema:"subnet_id"`
	VirtualMachineConfiguration     []VirtualMachineConfiguration     `tfschema:"virtual_machine_configuration"`
	VirtualMachineFullResourceNames []VirtualMachineFullResourceNames `tfschema:"virtual_machine_full_resource_names"`
}

type DiskVolumeConfiguration struct {
	VolumeName string `tfschema:"volume_name"`
	Count      int64  `tfschema:"count"`
	SizeGb     int64  `tfschema:"size_gb"`
	SkuName    string `tfschema:"sku_name"`
}

type VirtualMachineConfiguration struct {
	ImageReference []ImageReference `tfschema:"image_reference"`
	OSProfile      []OSProfile      `tfschema:"os_profile"`
	VmSize         string           `tfschema:"vm_size"`
}

type ImageReference struct {
	Offer     string `tfschema:"offer"`
	Publisher string `tfschema:"publisher"`
	Sku       string `tfschema:"sku"`
	Version   string `tfschema:"version"`
}

type OSProfile struct {
	AdminUsername string       `tfschema:"admin_username"`
	SshKeyPair    []SshKeyPair `tfschema:"ssh_key_pair"`
}

type SshKeyPair struct {
	PrivateKey string `tfschema:"private_key"`
	PublicKey  string `tfschema:"public_key"`
}

type VirtualMachineFullResourceNames struct {
	DataDiskNames         map[string]interface{} `tfschema:"data_disk_names"`
	HostName              string                 `tfschema:"host_name"`
	NetworkInterfaceNames []string               `tfschema:"network_interface_names"`
	OSDiskName            string                 `tfschema:"os_disk_name"`
	VMName                string                 `tfschema:"vm_name"`
}

type ThreeTierConfiguration struct {
	ApplicationServerConfiguration []ApplicationServerConfiguration `tfschema:"application_server_configuration"`
	AppResourceGroupName           string                           `tfschema:"app_resource_group_name"`
	CentralServerConfiguration     []CentralServerConfiguration     `tfschema:"central_server_configuration"`
	DatabaseServerConfiguration    []DatabaseServerConfiguration    `tfschema:"database_server_configuration"`
	FullResourceNames              []FullResourceNames              `tfschema:"full_resource_names"`
	HighAvailabilityType           string                           `tfschema:"high_availability_type"`
	IsSecondaryIpEnabled           bool                             `tfschema:"is_secondary_ip_enabled"`
	TransportCreateAndMount        []TransportCreateAndMount        `tfschema:"transport_create_and_mount"`
	TransportMount                 []TransportMount                 `tfschema:"transport_mount"`
}

type TransportCreateAndMount struct {
	ResourceGroupName  string `tfschema:"resource_group_name"`
	StorageAccountName string `tfschema:"storage_account_name"`
}

type TransportMount struct {
	ShareFileId       string `tfschema:"share_file_id"`
	PrivateEndpointId string `tfschema:"private_endpoint_id"`
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

type WorkloadsSAPVirtualInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPVirtualInstanceResource{}

func (r WorkloadsSAPVirtualInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_virtual_instance"
}

func (r WorkloadsSAPVirtualInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPVirtualInstanceModel{}
}

func (r WorkloadsSAPVirtualInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapvirtualinstances.ValidateSapVirtualInstanceID
}

func (r WorkloadsSAPVirtualInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SAPVirtualInstanceName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"environment": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sapvirtualinstances.SAPEnvironmentTypeNonProd),
				string(sapvirtualinstances.SAPEnvironmentTypeProd),
			}, false),
		},

		"identity": commonschema.UserAssignedIdentityRequired(),

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

		"deployment_with_os_configuration": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			AtLeastOneOf: []string{"deployment_with_os_configuration", "discovery_configuration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
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

								"deployer_vm_packages": {
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
												ValidateFunc: storageValidate.StorageAccountID,
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
							"deployment_with_os_configuration.0.single_server_configuration",
							"deployment_with_os_configuration.0.three_tier_configuration",
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

								"virtual_machine_configuration": SchemaForSAPVirtualInstanceVirtualMachineConfiguration(),

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
									Default:  false,
								},

								"virtual_machine_full_resource_names": SchemaForSAPVirtualInstanceVirtualMachineFullResourceNames(),
							},
						},
					},

					"three_tier_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						AtLeastOneOf: []string{
							"deployment_with_os_configuration.0.single_server_configuration",
							"deployment_with_os_configuration.0.three_tier_configuration",
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

											"virtual_machine_configuration": SchemaForSAPVirtualInstanceVirtualMachineConfiguration(),
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
											"virtual_machine_configuration": SchemaForSAPVirtualInstanceVirtualMachineConfiguration(),

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

											"virtual_machine_configuration": SchemaForSAPVirtualInstanceVirtualMachineConfiguration(),

											"database_type": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForSAPDatabaseType(), false),
											},

											"disk_volume_configuration": SchemaForSAPVirtualInstanceDiskVolumeConfiguration(),
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

														"virtual_machine": SchemaForSAPVirtualInstanceVirtualMachineFullResourceNames(),
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

														"load_balancer": SchemaForSAPVirtualInstanceLoadBalancerFullResourceNames(),

														"virtual_machine": SchemaForSAPVirtualInstanceVirtualMachineFullResourceNames(),
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

														"load_balancer": SchemaForSAPVirtualInstanceLoadBalancerFullResourceNames(),

														"virtual_machine": SchemaForSAPVirtualInstanceVirtualMachineFullResourceNames(),
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

								"is_secondary_ip_enabled": {
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
											"resource_group_name": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ForceNew:     true,
												ValidateFunc: resourcegroups.ValidateName,
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
											"share_file_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: storageValidate.StorageShareResourceManagerID,
											},

											"private_endpoint_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
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
		},

		"discovery_configuration": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			AtLeastOneOf: []string{"deployment_with_os_configuration", "discovery_configuration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"central_server_vm_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: computeValidate.VirtualMachineID,
					},

					"managed_storage_account_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: storageValidate.StorageAccountName,
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

func (r WorkloadsSAPVirtualInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSAPVirtualInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPVirtualInstanceModel
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

			identity, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
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

			if v := model.DeploymentWithOSConfiguration; v != nil {
				deploymentWithOSConfiguration, err := expandDeploymentWithOSConfiguration(v)
				if err != nil {
					return err
				}
				parameters.Properties.Configuration = deploymentWithOSConfiguration
			}

			if v := model.DiscoveryConfiguration; v != nil {
				parameters.Properties.Configuration = expandDiscoveryConfiguration(v)
			}

			if v := model.ManagedResourceGroupName; v != "" {
				parameters.Properties.ManagedResourceGroupConfiguration = &sapvirtualinstances.ManagedRGConfiguration{
					Name: utils.String(v),
				}
			}

			if _, err := client.Create(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{
					"Accepted",
					"Creating",
					"Preparing System Configuration",
				},
				Target:       []string{string(sapvirtualinstances.SapVirtualInstanceProvisioningStateSucceeded)},
				Refresh:      sapVirtualInstanceStateRefreshFunc(ctx, client, id),
				Timeout:      60 * time.Minute,
				PollInterval: 10 * time.Second,
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsSAPVirtualInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPVirtualInstanceModel
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

func (r WorkloadsSAPVirtualInstanceResource) Read() sdk.ResourceFunc {
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

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := WorkloadsSAPVirtualInstanceModel{
				Name:              id.SapVirtualInstanceName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			state.Identity = *identity

			properties := &model.Properties
			state.Environment = string(properties.Environment)
			state.SapProduct = string(properties.SapProduct)

			if properties.Configuration != nil {
				if v, ok := properties.Configuration.(sapvirtualinstances.DeploymentWithOSConfiguration); ok {
					state.DeploymentWithOSConfiguration = flattenDeploymentWithOSConfiguration(&v, metadata.ResourceData)
				}

				if v, ok := properties.Configuration.(sapvirtualinstances.DiscoveryConfiguration); ok {
					state.DiscoveryConfiguration = flattenDiscoveryConfiguration(&v)
				}
			}

			if v := properties.ManagedResourceGroupConfiguration; v != nil && v.Name != nil {
				state.ManagedResourceGroupName = *v.Name
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPVirtualInstanceResource) Delete() sdk.ResourceFunc {
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

func expandDiscoveryConfiguration(input []DiscoveryConfiguration) *sapvirtualinstances.DiscoveryConfiguration {
	configuration := &input[0]

	result := sapvirtualinstances.DiscoveryConfiguration{
		CentralServerVMId: utils.String(configuration.CentralServerVmId),
	}

	if v := configuration.ManagedStorageAccountName; v != "" {
		result.ManagedRgStorageAccountName = utils.String(v)
	}

	return &result
}

func flattenDiscoveryConfiguration(input *sapvirtualinstances.DiscoveryConfiguration) []DiscoveryConfiguration {
	if input == nil {
		return nil
	}

	result := DiscoveryConfiguration{
		CentralServerVmId: *input.CentralServerVMId,
	}

	if v := input.ManagedRgStorageAccountName; v != nil {
		result.ManagedStorageAccountName = *v
	}

	return []DiscoveryConfiguration{
		result,
	}
}

func expandDeploymentWithOSConfiguration(input []DeploymentWithOSConfiguration) (*sapvirtualinstances.DeploymentWithOSConfiguration, error) {
	configuration := &input[0]

	result := sapvirtualinstances.DeploymentWithOSConfiguration{
		AppLocation:        utils.String(configuration.AppLocation),
		OsSapConfiguration: expandOsSapConfiguration(configuration.OsSapConfiguration),
	}

	if len(configuration.SingleServerConfiguration) != 0 {
		singleServerConfiguration, err := expandSingleServerConfiguration(configuration.SingleServerConfiguration)
		if err != nil {
			return nil, err
		}
		result.InfrastructureConfiguration = singleServerConfiguration
	}

	if len(configuration.ThreeTierConfiguration) != 0 {
		result.InfrastructureConfiguration = expandThreeTierConfiguration(configuration.ThreeTierConfiguration)
	}

	return &result, nil
}

func expandSingleServerConfiguration(input []SingleServerConfiguration) (*sapvirtualinstances.SingleServerConfiguration, error) {
	singleServerConfiguration := &input[0]

	result := sapvirtualinstances.SingleServerConfiguration{
		AppResourceGroup: singleServerConfiguration.AppResourceGroupName,
		NetworkConfiguration: &sapvirtualinstances.NetworkConfiguration{
			IsSecondaryIPEnabled: utils.Bool(singleServerConfiguration.IsSecondaryIpEnabled),
		},
		SubnetId:                    singleServerConfiguration.SubnetId,
		VirtualMachineConfiguration: expandVirtualMachineConfiguration(singleServerConfiguration.VirtualMachineConfiguration),
	}

	if v := singleServerConfiguration.DatabaseType; v != "" {
		dbType := sapvirtualinstances.SAPDatabaseType(v)
		result.DatabaseType = &dbType
	}

	if v := singleServerConfiguration.DiskVolumeConfigurations; v != nil {
		result.DbDiskConfiguration = expandDiskVolumeConfigurations(v)
	}

	virtualMachineFullResourceNames, err := expandVirtualMachineFullResourceNames(singleServerConfiguration.VirtualMachineFullResourceNames)
	if err != nil {
		return nil, err
	}
	result.CustomResourceNames = virtualMachineFullResourceNames

	return &result, nil
}

func expandThreeTierConfiguration(input []ThreeTierConfiguration) *sapvirtualinstances.ThreeTierConfiguration {
	threeTierConfiguration := &input[0]

	result := sapvirtualinstances.ThreeTierConfiguration{
		ApplicationServer: expandApplicationServer(threeTierConfiguration.ApplicationServerConfiguration),
		AppResourceGroup:  threeTierConfiguration.AppResourceGroupName,
		CentralServer:     expandCentralServer(threeTierConfiguration.CentralServerConfiguration),
		DatabaseServer:    expandDatabaseServer(threeTierConfiguration.DatabaseServerConfiguration),
		NetworkConfiguration: &sapvirtualinstances.NetworkConfiguration{
			IsSecondaryIPEnabled: utils.Bool(threeTierConfiguration.IsSecondaryIpEnabled),
		},
		StorageConfiguration: expandStorageConfiguration(threeTierConfiguration),
	}

	if v := threeTierConfiguration.HighAvailabilityType; v != "" {
		result.HighAvailabilityConfig = &sapvirtualinstances.HighAvailabilityConfiguration{
			HighAvailabilityType: sapvirtualinstances.SAPHighAvailabilityType(v),
		}
	}

	if v := threeTierConfiguration.FullResourceNames; v != nil {
		result.CustomResourceNames = expandFullResourceNames(v)
	}

	return &result
}

func expandOsSapConfiguration(input []OsSapConfiguration) *sapvirtualinstances.OsSapConfiguration {
	if len(input) == 0 {
		return nil
	}

	osSapConfiguration := &input[0]

	result := sapvirtualinstances.OsSapConfiguration{
		SapFqdn: utils.String(osSapConfiguration.SapFqdn),
	}

	if v := osSapConfiguration.DeployerVmPackages; v != nil {
		result.DeployerVMPackages = expandDeployerVmPackages(v)
	}

	return &result
}

func expandDeployerVmPackages(input []DeployerVmPackages) *sapvirtualinstances.DeployerVMPackages {
	if len(input) == 0 {
		return nil
	}

	deployerVmPackages := &input[0]

	result := sapvirtualinstances.DeployerVMPackages{
		StorageAccountId: utils.String(deployerVmPackages.StorageAccountId),
		Url:              utils.String(deployerVmPackages.Url),
	}

	return &result
}

func flattenDeploymentWithOSConfiguration(input *sapvirtualinstances.DeploymentWithOSConfiguration, d *pluginsdk.ResourceData) []DeploymentWithOSConfiguration {
	if input == nil {
		return nil
	}

	result := DeploymentWithOSConfiguration{
		AppLocation:        *input.AppLocation,
		OsSapConfiguration: flattenOsSapConfiguration(input.OsSapConfiguration),
	}

	if configuration := input.InfrastructureConfiguration; configuration != nil {
		if v, ok := configuration.(sapvirtualinstances.SingleServerConfiguration); ok {
			result.SingleServerConfiguration = flattenSingleServerConfiguration(v, d, "deployment_with_os_configuration")
		}

		if v, ok := configuration.(sapvirtualinstances.ThreeTierConfiguration); ok {
			result.ThreeTierConfiguration = flattenThreeTierConfiguration(v, d, "deployment_with_os_configuration")
		}
	}

	return []DeploymentWithOSConfiguration{
		result,
	}
}

func flattenSingleServerConfiguration(input sapvirtualinstances.SingleServerConfiguration, d *pluginsdk.ResourceData, basePath string) []SingleServerConfiguration {
	result := SingleServerConfiguration{
		AppResourceGroupName:        input.AppResourceGroup,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.single_server_configuration", basePath)),
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

func flattenThreeTierConfiguration(input sapvirtualinstances.ThreeTierConfiguration, d *pluginsdk.ResourceData, basePath string) []ThreeTierConfiguration {
	result := ThreeTierConfiguration{
		AppResourceGroupName:           input.AppResourceGroup,
		ApplicationServerConfiguration: flattenApplicationServer(input.ApplicationServer, d, fmt.Sprintf("%s.0.three_tier_configuration", basePath)),
		CentralServerConfiguration:     flattenCentralServer(input.CentralServer, d, fmt.Sprintf("%s.0.three_tier_configuration", basePath)),
		DatabaseServerConfiguration:    flattenDatabaseServer(input.DatabaseServer, d, fmt.Sprintf("%s.0.three_tier_configuration", basePath)),
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

	if storageConfiguration := input.StorageConfiguration; storageConfiguration != nil {
		if transportFileShareConfiguration := storageConfiguration.TransportFileShareConfiguration; transportFileShareConfiguration != nil {
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
				transportMount := TransportMount{
					ShareFileId:       mountFileShareConfiguration.Id,
					PrivateEndpointId: mountFileShareConfiguration.PrivateEndpointId,
				}

				result.TransportMount = []TransportMount{
					transportMount,
				}
			}
		}
	}

	return []ThreeTierConfiguration{
		result,
	}
}

func flattenOsSapConfiguration(input *sapvirtualinstances.OsSapConfiguration) []OsSapConfiguration {
	if input == nil {
		return nil
	}

	result := OsSapConfiguration{
		SapFqdn: *input.SapFqdn,
	}

	if v := input.DeployerVMPackages; v != nil {
		result.DeployerVmPackages = flattenDeployerVMPackages(v)
	}

	return []OsSapConfiguration{
		result,
	}
}

func flattenDeployerVMPackages(input *sapvirtualinstances.DeployerVMPackages) []DeployerVmPackages {
	if input == nil {
		return nil
	}

	result := DeployerVmPackages{
		StorageAccountId: *input.StorageAccountId,
		Url:              *input.Url,
	}

	return []DeployerVmPackages{
		result,
	}
}

func sapVirtualInstanceStateRefreshFunc(ctx context.Context, client *sapvirtualinstances.SAPVirtualInstancesClient, id sapvirtualinstances.SapVirtualInstanceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if model := res.Model; model != nil {
			if model.Properties.ProvisioningState != nil {
				if *model.Properties.ProvisioningState == sapvirtualinstances.SapVirtualInstanceProvisioningStateFailed {
					return res, string(*model.Properties.ProvisioningState), fmt.Errorf("the provisioning state is in a failed state due to %s", *model.Properties.Errors.Properties.Message)
				}

				return res, string(*model.Properties.ProvisioningState), nil
			}
		}
		return nil, "", fmt.Errorf("unable to read state")
	}
}
