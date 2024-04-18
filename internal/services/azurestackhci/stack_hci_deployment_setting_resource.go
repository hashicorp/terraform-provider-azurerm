package azurestackhci

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = StackHCIDeploymentSettingResource{}

type StackHCIDeploymentSettingResource struct{}

func (StackHCIDeploymentSettingResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deploymentsettings.ValidateDeploymentSettingID
}

func (StackHCIDeploymentSettingResource) ResourceType() string {
	return "azurerm_stack_hci_deployment_setting"
}

type StackHCIDeploymentSettingModel struct {
	StackHCIClusterId string           `tfschema:"stack_hci_cluster_id"`
	ArcResourceIds    []string         `tfschema:"arc_resource_ids"`
	Version           string           `tfschema:"version"`
	ScaleUnit         []ScaleUnitModel `tfschema:"scale_unit"`
}

type ScaleUnitModel struct {
	AdouPath              string                       `tfschema:"adou_path"`
	Cluster               []ClusterModel               `tfschema:"cluster"`
	DomainFqdn            string                       `tfschema:"domain_fqdn"`
	HostNetwork           []HostNetworkModel           `tfschema:"host_network"`
	InfrastructureNetwork []InfrastructureNetworkModel `tfschema:"infrastructure_network"`
	NamingPrefix          string                       `tfschema:"naming_prefix"`
	Observability         []ObservabilityModel         `tfschema:"observability"`
	OptionalService       []OptionalServiceModel       `tfschema:"optional_service"`
	PhysicalNode          []PhysicalNodeModel          `tfschema:"physical_node"`
	SecretsLocation       string                       `tfschema:"secrets_location"`
	SecuritySetting       []SecuritySettingModel       `tfschema:"security_setting"`
	Storage               []StorageModel               `tfschema:"storage"`
}

type ClusterModel struct {
	AzureServiceEndpoint string `tfschema:"azure_service_endpoint"`
	CloudAccountName     string `tfschema:"cloud_account_name"`
	Name                 string `tfschema:"name"`
	WitnessType          string `tfschema:"witness_type"`
	WitnessPath          string `tfschema:"witness_path"`
}

type HostNetworkModel struct {
	Intent                        []HostNetworkIntentModel         `tfschema:"intent"`
	StorageAutoIpEnabled          bool                             `tfschema:"storage_auto_ip_enabled"`
	StorageConnectivitySwitchless bool                             `tfschema:"storage_connectivity_switchless"`
	StorageNetwork                []HostNetworkStorageNetworkModel `tfschema:"storage_network"`
}

type HostNetworkIntentModel struct {
	Adapter                                   []string                                                   `tfschema:"adapter"`
	AdapterPropertyOverride                   []HostNetworkIntentAdapterPropertyOverrideModel            `tfschema:"adapter_property_override"`
	Name                                      string                                                     `tfschema:"name"`
	OverrideAdapterPropertyEnabled            bool                                                       `tfschema:"override_adapter_property_enabled"`
	OverrideQosPolicyEnabled                  bool                                                       `tfschema:"override_qos_policy_enabled"`
	OverrideVirtualSwitchConfigurationEnabled bool                                                       `tfschema:"override_virtual_switch_configuration_enabled"`
	QosPolicyOverride                         []HostNetworkIntentQosPolicyOverrideModel                  `tfschema:"qos_policy_override"`
	TrafficType                               []string                                                   `tfschema:"traffic_type"`
	VirtualSwitchConfigurationOverride        []HostNetworkIntentVirtualSwitchConfigurationOverrideModel `tfschema:"virtual_switch_configuration_override"`
}

type HostNetworkIntentAdapterPropertyOverrideModel struct {
	JumboPacket             string `tfschema:"jumbo_packet"`
	NetworkDirect           string `tfschema:"network_direct"`
	NetworkDirectTechnology string `tfschema:"network_direct_technology"`
}

type HostNetworkIntentQosPolicyOverrideModel struct {
	BandWidthPercentageSMB         string `tfschema:"bandwidth_percentage_smb"`
	PriorityValue8021ActionCluster string `tfschema:"priority_value8021_action_cluster"`
	PriorityValue8021ActionSMB     string `tfschema:"priority_value8021_action_smb"`
}

type HostNetworkIntentVirtualSwitchConfigurationOverrideModel struct {
	EnableIov              string `tfschema:"enable_iov"`
	LoadBalancingAlgorithm string `tfschema:"load_balancing_algorithm"`
}

type HostNetworkStorageNetworkModel struct {
	Name               string `tfschema:"name"`
	NetworkAdapterName string `tfschema:"network_adapter_name"`
	VlanId             string `tfschema:"vlan_id"`
}

type InfrastructureNetworkModel struct {
	DhcpEnabled bool          `tfschema:"dhcp_enabled"`
	SubnetMask  string        `tfschema:"subnet_mask"`
	Gateway     string        `tfschema:"gateway"`
	IpPool      []IpPoolModel `tfschema:"ip_pool"`
	DnsServer   []string      `tfschema:"dns_server"`
}

type IpPoolModel struct {
	StartingAddress string `tfschema:"starting_address"`
	EndingAddress   string `tfschema:"ending_address"`
}

type ObservabilityModel struct {
	StreamingDataClientEnabled bool `tfschema:"streaming_data_client_enabled"`
	EuLocationEnabled          bool `tfschema:"eu_location_enabled"`
	EpisodicDataUploadEnabled  bool `tfschema:"episodic_data_upload_enabled"`
}

type OptionalServiceModel struct {
	CustomLocation string `tfschema:"custom_location"`
}

type PhysicalNodeModel struct {
	Name        string `tfschema:"name"`
	Ipv4Address string `tfschema:"ipv4_address"`
}

type SecuritySettingModel struct {
	BitlockerBootVolumeEnabled   bool `tfschema:"bitlocker_boot_volume_enabled"`
	BitlockerDataVolumeEnabled   bool `tfschema:"bitlocker_data_volume_enabled"`
	CredentialGuardEnabled       bool `tfschema:"credential_guard_enabled"`
	DriftControlEnabled          bool `tfschema:"drift_control_enabled"`
	DrtmProtectionEnabled        bool `tfschema:"drtm_protection_enabled"`
	HvciProtectionEnabled        bool `tfschema:"hvci_protection_enabled"`
	SideChannelMitigationEnabled bool `tfschema:"side_channel_mitigation_enabled"`
	SmbSigningEnabled            bool `tfschema:"smb_signing_enabled"`
	SmbClusterEncryptionEnabled  bool `tfschema:"smb_cluster_encryption_enabled"`
	WdacEnabled                  bool `tfschema:"wdac_enabled"`
}

type StorageModel struct {
	ConfigurationMode string `tfschema:"configuration_mode"`
}

func (StackHCIDeploymentSettingResource) ModelObject() interface{} {
	return &StackHCIDeploymentSettingModel{}
}

func (StackHCIDeploymentSettingResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"stack_hci_cluster_id": commonschema.ResourceIDReferenceRequiredForceNew(&deploymentsettings.ClusterId{}),

		"arc_resource_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: machines.ValidateMachineID,
			},
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$`),
				"must be in format `10.0.0.1`",
			),
		},

		"scale_unit": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"adou_path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"cluster": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile("^[a-zA-Z0-9-]{3,15}$"),
										"must be 3-15 characters long and contain only letters, numbers and hyphens",
									),
								},

								"azure_service_endpoint": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"cloud_account_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"witness_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Cloud",
										"FileShare",
									}, false),
								},

								"witness_path": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"domain_fqdn": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"host_network": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"intent": {
									Type:     pluginsdk.TypeList,
									Required: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"adapter": {
												Type:     pluginsdk.TypeList,
												Required: true,
												ForceNew: true,
												MinItems: 1,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},

											"name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"traffic_type": {
												Type:     pluginsdk.TypeList,
												Required: true,
												ForceNew: true,
												MinItems: 1,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
													ValidateFunc: validation.StringInSlice([]string{
														"Compute",
														"Storage",
														"Management",
													}, false),
												},
											},

											"adapter_property_override": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												ForceNew: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"jumbo_packet": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"network_direct": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"network_direct_technology": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},

											"override_adapter_property_enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												ForceNew: true,
											},

											"override_qos_policy_enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												ForceNew: true,
											},

											"override_virtual_switch_configuration_enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												ForceNew: true,
											},

											"qos_policy_override": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												ForceNew: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"bandwidth_percentage_smb": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"priority_value8021_action_cluster": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"priority_value8021_action_smb": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},

											"virtual_switch_configuration_override": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												ForceNew: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"enable_iov": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"load_balancing_algorithm": {
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

								"storage_network": {
									Type:     pluginsdk.TypeList,
									Required: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"network_adapter_name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"vlan_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"storage_auto_ip_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
									Default:  true,
								},

								"storage_connectivity_switchless": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
									Default:  false,
								},
							},
						},
					},

					"infrastructure_network": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"dns_server": {
									Type:     pluginsdk.TypeList,
									Required: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsIPv4Address,
									},
								},

								"gateway": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},

								"ip_pool": {
									Type:     pluginsdk.TypeList,
									Required: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"starting_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.IsIPv4Address,
											},

											"ending_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.IsIPv4Address,
											},
										},
									},
								},

								"subnet_mask": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},

								"dhcp_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},
							},
						},
					},

					"naming_prefix": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^[a-zA-Z0-9-]{1,8}$"),
							"must be 1-8 characters long and contain only letters, numbers and hyphens",
						),
					},

					"optional_service": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"custom_location": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"physical_node": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"ipv4_address": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},
							},
						},
					},

					"secrets_location": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},

					"storage": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"configuration_mode": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Express",
										"InfraOnly",
										"KeepStorage",
									}, false),
								},
							},
						},
					},

					"observability": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"streaming_data_client_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},

								"eu_location_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},

								"episodic_data_upload_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},
							},
						},
					},

					"security_setting": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"bitlocker_boot_volume_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},

								"bitlocker_data_volume_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},

								"credential_guard_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},

								"drift_control_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},

								"drtm_protection_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},

								"hvci_protection_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},

								"side_channel_mitigation_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},

								"smb_signing_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},

								"smb_cluster_encryption_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},

								"wdac_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
									ForceNew: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (StackHCIDeploymentSettingResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StackHCIDeploymentSettingResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.DeploymentSettings

			var config StackHCIDeploymentSettingModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			stackHCIClusterId, err := deploymentsettings.ParseClusterID(config.StackHCIClusterId)
			if err != nil {
				return err
			}
			id := deploymentsettings.NewDeploymentSettingID(stackHCIClusterId.SubscriptionId, stackHCIClusterId.ResourceGroupName, stackHCIClusterId.ClusterName, "default")

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := deploymentsettings.DeploymentSetting{
				Properties: &deploymentsettings.DeploymentSettingsProperties{
					ArcNodeResourceIds: config.ArcResourceIds,
					DeploymentMode:     deploymentsettings.DeploymentModeValidate,
					DeploymentConfiguration: deploymentsettings.DeploymentConfiguration{
						Version:    pointer.To(config.Version),
						ScaleUnits: ExpandDeploymentSettingScaleUnits(config.ScaleUnit),
					},
				},
			}

			// the resource exists even validation error
			metadata.SetID(id)

			// do validation
			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// do deployment
			payload.Properties.DeploymentMode = deploymentsettings.DeploymentModeDeploy
			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("deploying %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (StackHCIDeploymentSettingResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.DeploymentSettings

			id, err := deploymentsettings.ParseDeploymentSettingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			schema := StackHCIDeploymentSettingModel{
				StackHCIClusterId: deploymentsettings.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					schema.ArcResourceIds = props.ArcNodeResourceIds
					schema.Version = pointer.From(props.DeploymentConfiguration.Version)
					schema.ScaleUnit = FlattenDeploymentSettingScaleUnits(props.DeploymentConfiguration.ScaleUnits)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (StackHCIDeploymentSettingResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 1 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.DeploymentSettings

			id, err := deploymentsettings.ParseDeploymentSettingID(metadata.ResourceData.Id())
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

func ExpandDeploymentSettingScaleUnits(input []ScaleUnitModel) []deploymentsettings.ScaleUnits {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.ScaleUnits, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.ScaleUnits{
			DeploymentData: deploymentsettings.DeploymentData{
				AdouPath:              pointer.To(item.AdouPath),
				Cluster:               ExpandDeploymentSettingCluster(item.Cluster),
				DomainFqdn:            pointer.To(item.DomainFqdn),
				HostNetwork:           ExpandDeploymentSettingHostNetwork(item.HostNetwork),
				InfrastructureNetwork: ExpandDeploymentSettingInfrastructureNetwork(item.InfrastructureNetwork),
				NamingPrefix:          pointer.To(item.NamingPrefix),
				Observability:         ExpandDeploymentSettingObservability(item.Observability),
				OptionalServices:      ExpandDeploymentSettingOptionalService(item.OptionalService),
				PhysicalNodes:         ExpandDeploymentSettingPhysicalNode(item.PhysicalNode),
				SecretsLocation:       pointer.To(item.SecretsLocation),
				SecuritySettings:      ExpandDeploymentSettingSecuritySetting(item.SecuritySetting),
				Storage:               ExpandDeploymentSettingStorage(item.Storage),
			},
		})
	}

	return results
}

func FlattenDeploymentSettingScaleUnits(input []deploymentsettings.ScaleUnits) []ScaleUnitModel {
	if len(input) == 0 {
		return make([]ScaleUnitModel, 0)
	}

	results := make([]ScaleUnitModel, 0, len(input))
	for _, item := range input {
		results = append(results, ScaleUnitModel{
			AdouPath:              pointer.From(item.DeploymentData.AdouPath),
			Cluster:               FlattenDeploymentSettingCluster(item.DeploymentData.Cluster),
			DomainFqdn:            pointer.From(item.DeploymentData.DomainFqdn),
			HostNetwork:           FlattenDeploymentSettingHostNetwork(item.DeploymentData.HostNetwork),
			InfrastructureNetwork: FlattenDeploymentSettingInfrastructureNetwork(item.DeploymentData.InfrastructureNetwork),
			NamingPrefix:          pointer.From(item.DeploymentData.NamingPrefix),
			Observability:         FlattenDeploymentSettingObservability(item.DeploymentData.Observability),
			OptionalService:       FlattenDeploymentSettingOptionalService(item.DeploymentData.OptionalServices),
			PhysicalNode:          FlattenDeploymentSettingPhysicalNode(item.DeploymentData.PhysicalNodes),
			SecretsLocation:       pointer.From(item.DeploymentData.SecretsLocation),
			SecuritySetting:       FlattenDeploymentSettingSecuritySetting(item.DeploymentData.SecuritySettings),
			Storage:               FlattenDeploymentSettingStorage(item.DeploymentData.Storage),
		})
	}

	return results
}

func ExpandDeploymentSettingCluster(input []ClusterModel) *deploymentsettings.DeploymentCluster {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.DeploymentCluster{
		AzureServiceEndpoint: pointer.To(v.AzureServiceEndpoint),
		CloudAccountName:     pointer.To(v.CloudAccountName),
		Name:                 pointer.To(v.Name),
		WitnessType:          pointer.To(v.WitnessType),
		WitnessPath:          pointer.To(v.WitnessPath),
	}
}

func FlattenDeploymentSettingCluster(input *deploymentsettings.DeploymentCluster) []ClusterModel {
	if input == nil {
		return make([]ClusterModel, 0)
	}

	return []ClusterModel{{
		AzureServiceEndpoint: pointer.From(input.AzureServiceEndpoint),
		CloudAccountName:     pointer.From(input.CloudAccountName),
		Name:                 pointer.From(input.Name),
		WitnessType:          pointer.From(input.WitnessType),
		WitnessPath:          pointer.From(input.WitnessPath),
	}}
}

func ExpandDeploymentSettingHostNetwork(input []HostNetworkModel) *deploymentsettings.HostNetwork {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.HostNetwork{
		Intents:                       ExpandDeploymentSettingHostNetworkIntent(v.Intent),
		EnableStorageAutoIP:           pointer.To(v.StorageAutoIpEnabled),
		StorageConnectivitySwitchless: pointer.To(v.StorageConnectivitySwitchless),
		StorageNetworks:               ExpandDeploymentSettingHostNetworkStorageNetwork(v.StorageNetwork),
	}
}

func FlattenDeploymentSettingHostNetwork(input *deploymentsettings.HostNetwork) []HostNetworkModel {
	if input == nil {
		return make([]HostNetworkModel, 0)
	}

	return []HostNetworkModel{{
		Intent:                        FlattenDeploymentSettingHostNetworkIntent(input.Intents),
		StorageAutoIpEnabled:          pointer.From(input.EnableStorageAutoIP),
		StorageConnectivitySwitchless: pointer.From(input.StorageConnectivitySwitchless),
		StorageNetwork:                FlattenDeploymentSettingHostNetworkStorageNetwork(input.StorageNetworks),
	}}
}

func ExpandDeploymentSettingHostNetworkIntent(input []HostNetworkIntentModel) *[]deploymentsettings.Intents {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.Intents, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.Intents{
			Adapter:                             pointer.To(item.Adapter),
			AdapterPropertyOverrides:            ExpandHostNetworkIntentAdapterPropertyOverride(item.AdapterPropertyOverride),
			Name:                                pointer.To(item.Name),
			OverrideAdapterProperty:             pointer.To(item.OverrideAdapterPropertyEnabled),
			OverrideQosPolicy:                   pointer.To(item.OverrideQosPolicyEnabled),
			OverrideVirtualSwitchConfiguration:  pointer.To(item.OverrideVirtualSwitchConfigurationEnabled),
			QosPolicyOverrides:                  ExpandHostNetworkIntentQosPolicyOverride(item.QosPolicyOverride),
			TrafficType:                         pointer.To(item.TrafficType),
			VirtualSwitchConfigurationOverrides: ExpandHostNetworkIntentVirtualSwitchConfigurationOverride(item.VirtualSwitchConfigurationOverride),
		})
	}

	return &results
}

func FlattenDeploymentSettingHostNetworkIntent(input *[]deploymentsettings.Intents) []HostNetworkIntentModel {
	if input == nil {
		return make([]HostNetworkIntentModel, 0)
	}

	results := make([]HostNetworkIntentModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, HostNetworkIntentModel{
			Adapter:                        pointer.From(item.Adapter),
			AdapterPropertyOverride:        FlattenHostNetworkIntentAdapterPropertyOverride(item.AdapterPropertyOverrides),
			Name:                           pointer.From(item.Name),
			OverrideAdapterPropertyEnabled: pointer.From(item.OverrideAdapterProperty),
			OverrideQosPolicyEnabled:       pointer.From(item.OverrideQosPolicy),
			OverrideVirtualSwitchConfigurationEnabled: pointer.From(item.OverrideVirtualSwitchConfiguration),
			QosPolicyOverride:                         FlattenHostNetworkIntentQosPolicyOverride(item.QosPolicyOverrides),
			TrafficType:                               pointer.From(item.TrafficType),
			VirtualSwitchConfigurationOverride:        FlattenHostNetworkIntentVirtualSwitchConfigurationOverride(item.VirtualSwitchConfigurationOverrides),
		})
	}

	return results
}

func ExpandHostNetworkIntentAdapterPropertyOverride(input []HostNetworkIntentAdapterPropertyOverrideModel) *deploymentsettings.AdapterPropertyOverrides {
	if len(input) == 0 {
		return &deploymentsettings.AdapterPropertyOverrides{
			JumboPacket:             pointer.To(""),
			NetworkDirect:           pointer.To(""),
			NetworkDirectTechnology: pointer.To(""),
		}
	}

	v := input[0]

	return &deploymentsettings.AdapterPropertyOverrides{
		JumboPacket:             pointer.To(v.JumboPacket),
		NetworkDirect:           pointer.To(v.NetworkDirect),
		NetworkDirectTechnology: pointer.To(v.NetworkDirectTechnology),
	}
}

func FlattenHostNetworkIntentAdapterPropertyOverride(input *deploymentsettings.AdapterPropertyOverrides) []HostNetworkIntentAdapterPropertyOverrideModel {
	if input == nil {
		return make([]HostNetworkIntentAdapterPropertyOverrideModel, 0)
	}

	jumboPacket := pointer.From(input.JumboPacket)
	networkDirect := pointer.From(input.NetworkDirect)
	networkDirectTechnology := pointer.From(input.NetworkDirectTechnology)

	// server will return the block with empty string in all fields by default
	if jumboPacket == "" && networkDirect == "" && networkDirectTechnology == "" {
		return make([]HostNetworkIntentAdapterPropertyOverrideModel, 0)
	}

	return []HostNetworkIntentAdapterPropertyOverrideModel{{
		JumboPacket:             jumboPacket,
		NetworkDirect:           networkDirect,
		NetworkDirectTechnology: networkDirectTechnology,
	}}
}

func ExpandHostNetworkIntentQosPolicyOverride(input []HostNetworkIntentQosPolicyOverrideModel) *deploymentsettings.QosPolicyOverrides {
	if len(input) == 0 {
		return &deploymentsettings.QosPolicyOverrides{
			BandwidthPercentageSMB:         pointer.To(""),
			PriorityValue8021ActionCluster: pointer.To(""),
			PriorityValue8021ActionSMB:     pointer.To(""),
		}
	}

	v := input[0]

	return &deploymentsettings.QosPolicyOverrides{
		BandwidthPercentageSMB:         pointer.To(v.BandWidthPercentageSMB),
		PriorityValue8021ActionCluster: pointer.To(v.PriorityValue8021ActionCluster),
		PriorityValue8021ActionSMB:     pointer.To(v.PriorityValue8021ActionSMB),
	}
}

func FlattenHostNetworkIntentQosPolicyOverride(input *deploymentsettings.QosPolicyOverrides) []HostNetworkIntentQosPolicyOverrideModel {
	if input == nil {
		return make([]HostNetworkIntentQosPolicyOverrideModel, 0)
	}

	bandwidthPercentageSMB := pointer.From(input.BandwidthPercentageSMB)
	priorityValue8021ActionCluster := pointer.From(input.PriorityValue8021ActionCluster)
	priorityValue8021ActionSMB := pointer.From(input.PriorityValue8021ActionSMB)

	// server will return the block with empty string in all fields by default
	if bandwidthPercentageSMB == "" && priorityValue8021ActionCluster == "" && priorityValue8021ActionSMB == "" {
		return make([]HostNetworkIntentQosPolicyOverrideModel, 0)
	}

	return []HostNetworkIntentQosPolicyOverrideModel{{
		BandWidthPercentageSMB:         bandwidthPercentageSMB,
		PriorityValue8021ActionCluster: priorityValue8021ActionCluster,
		PriorityValue8021ActionSMB:     priorityValue8021ActionSMB,
	}}
}

func ExpandHostNetworkIntentVirtualSwitchConfigurationOverride(input []HostNetworkIntentVirtualSwitchConfigurationOverrideModel) *deploymentsettings.VirtualSwitchConfigurationOverrides {
	if len(input) == 0 {
		return &deploymentsettings.VirtualSwitchConfigurationOverrides{
			EnableIov:              pointer.To(""),
			LoadBalancingAlgorithm: pointer.To(""),
		}
	}

	v := input[0]

	return &deploymentsettings.VirtualSwitchConfigurationOverrides{
		EnableIov:              pointer.To(v.EnableIov),
		LoadBalancingAlgorithm: pointer.To(v.LoadBalancingAlgorithm),
	}
}

func FlattenHostNetworkIntentVirtualSwitchConfigurationOverride(input *deploymentsettings.VirtualSwitchConfigurationOverrides) []HostNetworkIntentVirtualSwitchConfigurationOverrideModel {
	if input == nil {
		return make([]HostNetworkIntentVirtualSwitchConfigurationOverrideModel, 0)
	}

	enableIov := pointer.From(input.EnableIov)
	loadBalancingAlgorithm := pointer.From(input.LoadBalancingAlgorithm)

	// server will return the block with empty string in all fields by default
	if enableIov == "" && loadBalancingAlgorithm == "" {
		return make([]HostNetworkIntentVirtualSwitchConfigurationOverrideModel, 0)
	}

	return []HostNetworkIntentVirtualSwitchConfigurationOverrideModel{{
		EnableIov:              enableIov,
		LoadBalancingAlgorithm: loadBalancingAlgorithm,
	}}
}

func ExpandDeploymentSettingHostNetworkStorageNetwork(input []HostNetworkStorageNetworkModel) *[]deploymentsettings.StorageNetworks {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.StorageNetworks, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.StorageNetworks{
			Name:               pointer.To(item.Name),
			NetworkAdapterName: pointer.To(item.NetworkAdapterName),
			VlanId:             pointer.To(item.VlanId),
		})
	}

	return &results
}

func FlattenDeploymentSettingHostNetworkStorageNetwork(input *[]deploymentsettings.StorageNetworks) []HostNetworkStorageNetworkModel {
	if input == nil {
		return make([]HostNetworkStorageNetworkModel, 0)
	}

	results := make([]HostNetworkStorageNetworkModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, HostNetworkStorageNetworkModel{
			Name:               pointer.From(item.Name),
			NetworkAdapterName: pointer.From(item.NetworkAdapterName),
			VlanId:             pointer.From(item.VlanId),
		})
	}

	return results
}

func ExpandDeploymentSettingInfrastructureNetwork(input []InfrastructureNetworkModel) *[]deploymentsettings.InfrastructureNetwork {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.InfrastructureNetwork, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.InfrastructureNetwork{
			DnsServers: pointer.To(item.DnsServer),
			Gateway:    pointer.To(item.Gateway),
			IPPools:    ExpandDeploymentSettingInfrastructureNetworkIpPool(item.IpPool),
			SubnetMask: pointer.To(item.SubnetMask),
			UseDhcp:    pointer.To(item.DhcpEnabled),
		})
	}

	return &results
}

func FlattenDeploymentSettingInfrastructureNetwork(input *[]deploymentsettings.InfrastructureNetwork) []InfrastructureNetworkModel {
	if input == nil {
		return make([]InfrastructureNetworkModel, 0)
	}

	results := make([]InfrastructureNetworkModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, InfrastructureNetworkModel{
			DhcpEnabled: pointer.From(item.UseDhcp),
			DnsServer:   pointer.From(item.DnsServers),
			Gateway:     pointer.From(item.Gateway),
			IpPool:      FlattenDeploymentSettingInfrastructureNetworkIpPool(item.IPPools),
			SubnetMask:  pointer.From(item.SubnetMask),
		})
	}

	return results
}

func ExpandDeploymentSettingInfrastructureNetworkIpPool(input []IpPoolModel) *[]deploymentsettings.IPPools {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.IPPools, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.IPPools{
			EndingAddress:   pointer.To(item.EndingAddress),
			StartingAddress: pointer.To(item.StartingAddress),
		})
	}

	return &results
}

func FlattenDeploymentSettingInfrastructureNetworkIpPool(input *[]deploymentsettings.IPPools) []IpPoolModel {
	if input == nil {
		return make([]IpPoolModel, 0)
	}

	results := make([]IpPoolModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, IpPoolModel{
			EndingAddress:   pointer.From(item.EndingAddress),
			StartingAddress: pointer.From(item.StartingAddress),
		})
	}

	return results
}

func ExpandDeploymentSettingObservability(input []ObservabilityModel) *deploymentsettings.Observability {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.Observability{
		EpisodicDataUpload:  pointer.To(v.EpisodicDataUploadEnabled),
		EuLocation:          pointer.To(v.EuLocationEnabled),
		StreamingDataClient: pointer.To(v.StreamingDataClientEnabled),
	}
}

func FlattenDeploymentSettingObservability(input *deploymentsettings.Observability) []ObservabilityModel {
	if input == nil {
		return make([]ObservabilityModel, 0)
	}

	return []ObservabilityModel{{
		EpisodicDataUploadEnabled:  pointer.From(input.EpisodicDataUpload),
		EuLocationEnabled:          pointer.From(input.EuLocation),
		StreamingDataClientEnabled: pointer.From(input.StreamingDataClient),
	}}
}

func ExpandDeploymentSettingOptionalService(input []OptionalServiceModel) *deploymentsettings.OptionalServices {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.OptionalServices{
		CustomLocation: pointer.To(v.CustomLocation),
	}
}

func FlattenDeploymentSettingOptionalService(input *deploymentsettings.OptionalServices) []OptionalServiceModel {
	if input == nil {
		return make([]OptionalServiceModel, 0)
	}

	return []OptionalServiceModel{{
		CustomLocation: pointer.From(input.CustomLocation),
	}}
}

func ExpandDeploymentSettingPhysicalNode(input []PhysicalNodeModel) *[]deploymentsettings.PhysicalNodes {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.PhysicalNodes, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.PhysicalNodes{
			IPv4Address: pointer.To(item.Ipv4Address),
			Name:        pointer.To(item.Name),
		})
	}

	return &results
}

func FlattenDeploymentSettingPhysicalNode(input *[]deploymentsettings.PhysicalNodes) []PhysicalNodeModel {
	if input == nil {
		return make([]PhysicalNodeModel, 0)
	}

	results := make([]PhysicalNodeModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, PhysicalNodeModel{
			Ipv4Address: pointer.From(item.IPv4Address),
			Name:        pointer.From(item.Name),
		})
	}

	return results
}

func ExpandDeploymentSettingSecuritySetting(input []SecuritySettingModel) *deploymentsettings.DeploymentSecuritySettings {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.DeploymentSecuritySettings{
		BitlockerBootVolume:           pointer.To(v.BitlockerBootVolumeEnabled),
		BitlockerDataVolumes:          pointer.To(v.BitlockerDataVolumeEnabled),
		CredentialGuardEnforced:       pointer.To(v.CredentialGuardEnabled),
		DriftControlEnforced:          pointer.To(v.DriftControlEnabled),
		DrtmProtection:                pointer.To(v.DrtmProtectionEnabled),
		HvciProtection:                pointer.To(v.HvciProtectionEnabled),
		SideChannelMitigationEnforced: pointer.To(v.SideChannelMitigationEnabled),
		SmbClusterEncryption:          pointer.To(v.SmbClusterEncryptionEnabled),
		SmbSigningEnforced:            pointer.To(v.SmbSigningEnabled),
		WdacEnforced:                  pointer.To(v.WdacEnabled),
	}
}

func FlattenDeploymentSettingSecuritySetting(input *deploymentsettings.DeploymentSecuritySettings) []SecuritySettingModel {
	if input == nil {
		return make([]SecuritySettingModel, 0)
	}

	return []SecuritySettingModel{{
		BitlockerBootVolumeEnabled:   pointer.From(input.BitlockerBootVolume),
		BitlockerDataVolumeEnabled:   pointer.From(input.BitlockerDataVolumes),
		CredentialGuardEnabled:       pointer.From(input.CredentialGuardEnforced),
		DriftControlEnabled:          pointer.From(input.DriftControlEnforced),
		DrtmProtectionEnabled:        pointer.From(input.DrtmProtection),
		HvciProtectionEnabled:        pointer.From(input.HvciProtection),
		SideChannelMitigationEnabled: pointer.From(input.SideChannelMitigationEnforced),
		SmbClusterEncryptionEnabled:  pointer.From(input.SmbClusterEncryption),
		SmbSigningEnabled:            pointer.From(input.SmbSigningEnforced),
		WdacEnabled:                  pointer.From(input.WdacEnforced),
	}}
}

func ExpandDeploymentSettingStorage(input []StorageModel) *deploymentsettings.Storage {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.Storage{
		ConfigurationMode: pointer.To(v.ConfigurationMode),
	}
}

func FlattenDeploymentSettingStorage(input *deploymentsettings.Storage) []StorageModel {
	if input == nil {
		return make([]StorageModel, 0)
	}

	return []StorageModel{{
		ConfigurationMode: pointer.From(input.ConfigurationMode),
	}}
}
