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
	CentralServerVmId         string `tfschema:"central_server_virtual_machine_id"`
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
	IsSecondaryIpEnabled            bool                              `tfschema:"secondary_ip_enabled"`
	SubnetId                        string                            `tfschema:"subnet_id"`
	VirtualMachineConfiguration     []VirtualMachineConfiguration     `tfschema:"virtual_machine_configuration"`
	VirtualMachineFullResourceNames []VirtualMachineFullResourceNames `tfschema:"virtual_machine_full_resource_names"`
}

type DiskVolumeConfiguration struct {
	VolumeName string `tfschema:"volume_name"`
	Count      int64  `tfschema:"count"`
	SizeGb     int64  `tfschema:"size_in_gb"`
	SkuName    string `tfschema:"sku_name"`
}

type VirtualMachineConfiguration struct {
	ImageReference []ImageReference `tfschema:"image"`
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
	AdminUsername string `tfschema:"admin_username"`
	SshPrivateKey string `tfschema:"ssh_private_key"`
	SshPublicKey  string `tfschema:"ssh_public_key"`
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

								"secondary_ip_enabled": {
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
					"central_server_virtual_machine_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateVirtualMachineID,
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

			state := WorkloadsSAPVirtualInstanceModel{}
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
						state.DeploymentWithOSConfiguration = flattenDeploymentWithOSConfiguration(v, metadata.ResourceData, id.SubscriptionId)
					}

					if v, ok := config.(sapvirtualinstances.DiscoveryConfiguration); ok {
						state.DiscoveryConfiguration = flattenDiscoveryConfiguration(v)
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

func expandDiscoveryConfiguration(input []DiscoveryConfiguration) sapvirtualinstances.DiscoveryConfiguration {
	configuration := &input[0]

	result := sapvirtualinstances.DiscoveryConfiguration{
		CentralServerVMId: utils.String(configuration.CentralServerVmId),
	}

	if v := configuration.ManagedStorageAccountName; v != "" {
		result.ManagedRgStorageAccountName = utils.String(v)
	}

	return result
}

func flattenDiscoveryConfiguration(input sapvirtualinstances.DiscoveryConfiguration) []DiscoveryConfiguration {
	result := make([]DiscoveryConfiguration, 0)

	return append(result, DiscoveryConfiguration{
		CentralServerVmId:         pointer.From(input.CentralServerVMId),
		ManagedStorageAccountName: pointer.From(input.ManagedRgStorageAccountName),
	})
}

func expandDeploymentWithOSConfiguration(input []DeploymentWithOSConfiguration) (sapvirtualinstances.DeploymentWithOSConfiguration, error) {
	configuration := &input[0]

	result := sapvirtualinstances.DeploymentWithOSConfiguration{
		AppLocation:        utils.String(configuration.AppLocation),
		OsSapConfiguration: expandOsSapConfiguration(configuration.OsSapConfiguration),
	}

	if len(configuration.SingleServerConfiguration) != 0 {
		singleServerConfiguration, err := expandSingleServerConfiguration(configuration.SingleServerConfiguration)
		if err != nil {
			return sapvirtualinstances.DeploymentWithOSConfiguration{}, err
		}
		result.InfrastructureConfiguration = singleServerConfiguration
	}

	if len(configuration.ThreeTierConfiguration) != 0 {
		threeTierConfiguration, err := expandThreeTierConfiguration(configuration.ThreeTierConfiguration)
		if err != nil {
			return sapvirtualinstances.DeploymentWithOSConfiguration{}, err
		}
		result.InfrastructureConfiguration = threeTierConfiguration
	}

	return result, nil
}

func expandSingleServerConfiguration(input []SingleServerConfiguration) (sapvirtualinstances.SingleServerConfiguration, error) {
	singleServerConfiguration := &input[0]

	virtualMachineFullResourceNames, err := expandVirtualMachineFullResourceNames(singleServerConfiguration.VirtualMachineFullResourceNames)
	if err != nil {
		return sapvirtualinstances.SingleServerConfiguration{}, err
	}

	result := sapvirtualinstances.SingleServerConfiguration{
		AppResourceGroup:    singleServerConfiguration.AppResourceGroupName,
		CustomResourceNames: virtualMachineFullResourceNames,
		DbDiskConfiguration: expandDiskVolumeConfigurations(singleServerConfiguration.DiskVolumeConfigurations),
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

	return result, nil
}

func expandThreeTierConfiguration(input []ThreeTierConfiguration) (sapvirtualinstances.ThreeTierConfiguration, error) {
	threeTierConfiguration := &input[0]

	result := sapvirtualinstances.ThreeTierConfiguration{
		ApplicationServer:   expandApplicationServer(threeTierConfiguration.ApplicationServerConfiguration),
		AppResourceGroup:    threeTierConfiguration.AppResourceGroupName,
		CentralServer:       expandCentralServer(threeTierConfiguration.CentralServerConfiguration),
		CustomResourceNames: expandFullResourceNames(threeTierConfiguration.FullResourceNames),
		DatabaseServer:      expandDatabaseServer(threeTierConfiguration.DatabaseServerConfiguration),
		NetworkConfiguration: &sapvirtualinstances.NetworkConfiguration{
			IsSecondaryIPEnabled: utils.Bool(threeTierConfiguration.IsSecondaryIpEnabled),
		},
	}

	storageConfiguration, err := expandStorageConfiguration(threeTierConfiguration)
	if err != nil {
		return sapvirtualinstances.ThreeTierConfiguration{}, err
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

	osSapConfiguration := &input[0]

	result := sapvirtualinstances.OsSapConfiguration{
		DeployerVMPackages: expandDeployerVmPackages(osSapConfiguration.DeployerVmPackages),
		SapFqdn:            utils.String(osSapConfiguration.SapFqdn),
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

func flattenDeploymentWithOSConfiguration(input sapvirtualinstances.DeploymentWithOSConfiguration, d *pluginsdk.ResourceData, subscriptionId string) []DeploymentWithOSConfiguration {
	result := make([]DeploymentWithOSConfiguration, 0)

	deploymentWithOSConfiguration := DeploymentWithOSConfiguration{
		AppLocation:        pointer.From(input.AppLocation),
		OsSapConfiguration: flattenOsSapConfiguration(input.OsSapConfiguration),
	}

	if configuration := input.InfrastructureConfiguration; configuration != nil {
		if v, ok := configuration.(sapvirtualinstances.SingleServerConfiguration); ok {
			deploymentWithOSConfiguration.SingleServerConfiguration = flattenSingleServerConfiguration(v, d, "deployment_with_os_configuration")
		}

		if v, ok := configuration.(sapvirtualinstances.ThreeTierConfiguration); ok {
			deploymentWithOSConfiguration.ThreeTierConfiguration = flattenThreeTierConfiguration(v, d, "deployment_with_os_configuration", subscriptionId)
		}
	}

	return append(result, deploymentWithOSConfiguration)
}

func flattenSingleServerConfiguration(input sapvirtualinstances.SingleServerConfiguration, d *pluginsdk.ResourceData, basePath string) []SingleServerConfiguration {
	result := make([]SingleServerConfiguration, 0)

	singleServerConfig := SingleServerConfiguration{
		AppResourceGroupName:        input.AppResourceGroup,
		DatabaseType:                string(pointer.From(input.DatabaseType)),
		DiskVolumeConfigurations:    flattenDiskVolumeConfigurations(input.DbDiskConfiguration),
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.single_server_configuration", basePath)),
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

func flattenThreeTierConfiguration(input sapvirtualinstances.ThreeTierConfiguration, d *pluginsdk.ResourceData, basePath string, subscriptionId string) []ThreeTierConfiguration {
	result := make([]ThreeTierConfiguration, 0)

	threeTierConfig := ThreeTierConfiguration{
		AppResourceGroupName:           input.AppResourceGroup,
		ApplicationServerConfiguration: flattenApplicationServer(input.ApplicationServer, d, fmt.Sprintf("%s.0.three_tier_configuration", basePath)),
		CentralServerConfiguration:     flattenCentralServer(input.CentralServer, d, fmt.Sprintf("%s.0.three_tier_configuration", basePath)),
		DatabaseServerConfiguration:    flattenDatabaseServer(input.DatabaseServer, d, fmt.Sprintf("%s.0.three_tier_configuration", basePath)),
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

			if mountFileShareConfiguration, ok := transportFileShareConfiguration.(sapvirtualinstances.MountFileShareConfiguration); ok {
				// Currently, the last segment of the Storage File Share resource manager ID in Swagger is defined as `/shares/` but it's unexpected.
				// The last segment of the Storage File Share resource manager ID should be `/fileshares/` not `/shares/` in Swagger since the backend service is using `/fileshares/`.
				// See more details from https://github.com/Azure/azure-rest-api-specs/issues/25209
				threeTierConfig.TransportMount = []TransportMount{
					{
						FileShareId:       strings.Replace(mountFileShareConfiguration.Id, "/shares/", "/fileshares/", 1),
						PrivateEndpointId: mountFileShareConfiguration.PrivateEndpointId,
					},
				}
			}
		}
	}

	return append(result, threeTierConfig)
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

func sapVirtualInstanceStateRefreshFunc(ctx context.Context, client *sapvirtualinstances.SAPVirtualInstancesClient, id sapvirtualinstances.SapVirtualInstanceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if model := res.Model; model != nil {
			if provisioningState := model.Properties.ProvisioningState; provisioningState != nil {
				if *provisioningState == sapvirtualinstances.SapVirtualInstanceProvisioningStateFailed {
					errorMessage := "the provisioning state is in a failed state"

					if model.Properties.Errors != nil && model.Properties.Errors.Properties != nil && model.Properties.Errors.Properties.Message != nil {
						errorMessage = fmt.Sprintf("%s due to %s", errorMessage, *model.Properties.Errors.Properties.Message)
					}

					return res, string(*provisioningState), fmt.Errorf(errorMessage)
				}

				return res, string(*provisioningState), nil
			}
		}
		return nil, "", fmt.Errorf("unable to read state")
	}
}
