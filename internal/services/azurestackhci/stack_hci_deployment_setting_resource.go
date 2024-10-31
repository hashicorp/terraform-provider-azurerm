// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/storagecontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
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
	ActiveDirectoryOrganizationalUnitPath string                       `tfschema:"active_directory_organizational_unit_path"`
	Cluster                               []ClusterModel               `tfschema:"cluster"`
	DomainFqdn                            string                       `tfschema:"domain_fqdn"`
	HostNetwork                           []HostNetworkModel           `tfschema:"host_network"`
	InfrastructureNetwork                 []InfrastructureNetworkModel `tfschema:"infrastructure_network"`
	NamePrefix                            string                       `tfschema:"name_prefix"`
	OptionalService                       []OptionalServiceModel       `tfschema:"optional_service"`
	PhysicalNode                          []PhysicalNodeModel          `tfschema:"physical_node"`
	SecretsLocation                       string                       `tfschema:"secrets_location"`
	Storage                               []StorageModel               `tfschema:"storage"`

	// flatten 'observability' block, for API always return them
	StreamingDataClientEnabled bool `tfschema:"streaming_data_client_enabled"`
	EuLocationEnabled          bool `tfschema:"eu_location_enabled"`
	EpisodicDataUploadEnabled  bool `tfschema:"episodic_data_upload_enabled"`

	// flatten 'securitySetting' block, for API always return them
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

type ClusterModel struct {
	AzureServiceEndpoint string `tfschema:"azure_service_endpoint"`
	CloudAccountName     string `tfschema:"cloud_account_name"`
	Name                 string `tfschema:"name"`
	WitnessType          string `tfschema:"witness_type"`
	WitnessPath          string `tfschema:"witness_path"`
}

type HostNetworkModel struct {
	Intent                               []HostNetworkIntentModel         `tfschema:"intent"`
	StorageAutoIpEnabled                 bool                             `tfschema:"storage_auto_ip_enabled"`
	StorageConnectivitySwitchlessEnabled bool                             `tfschema:"storage_connectivity_switchless_enabled"`
	StorageNetwork                       []HostNetworkStorageNetworkModel `tfschema:"storage_network"`
}

type HostNetworkIntentModel struct {
	Adapter                                   []string                                                   `tfschema:"adapter"`
	AdapterPropertyOverride                   []HostNetworkIntentAdapterPropertyOverrideModel            `tfschema:"adapter_property_override"`
	AdapterPropertyOverrideEnabled            bool                                                       `tfschema:"adapter_property_override_enabled"`
	Name                                      string                                                     `tfschema:"name"`
	QosPolicyOverride                         []HostNetworkIntentQosPolicyOverrideModel                  `tfschema:"qos_policy_override"`
	QosPolicyOverrideEnabled                  bool                                                       `tfschema:"qos_policy_override_enabled"`
	TrafficType                               []string                                                   `tfschema:"traffic_type"`
	VirtualSwitchConfigurationOverride        []HostNetworkIntentVirtualSwitchConfigurationOverrideModel `tfschema:"virtual_switch_configuration_override"`
	VirtualSwitchConfigurationOverrideEnabled bool                                                       `tfschema:"virtual_switch_configuration_override_enabled"`
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

type OptionalServiceModel struct {
	CustomLocation string `tfschema:"custom_location"`
}

type PhysicalNodeModel struct {
	Name        string `tfschema:"name"`
	Ipv4Address string `tfschema:"ipv4_address"`
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
				"the version must be a set of numbers separated by dots: `10.0.0.1`",
			),
		},

		"scale_unit": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"active_directory_organizational_unit_path": {
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
										"the cluster name must be 3-15 characters long and contain only letters, numbers and hyphens",
									),
								},

								"azure_service_endpoint": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`),
										"`azure_service_endpoint` must be a valid domain name, for example, \"core.windows.net\"",
									),
								},

								"cloud_account_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: storageValidate.StorageAccountName,
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
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`),
							"`domain_fqdn` must be a valid domain name, for example, \"jumpstart.local\"",
						),
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
											"name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

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

											"adapter_property_override_enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												ForceNew: true,
												Default:  false,
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

											"qos_policy_override_enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												ForceNew: true,
												Default:  false,
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

											"virtual_switch_configuration_override_enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												ForceNew: true,
												Default:  false,
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

								"storage_connectivity_switchless_enabled": {
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
									Default:  false,
								},
							},
						},
					},

					"name_prefix": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^[a-zA-Z0-9-]{1,8}$"),
							"the naming prefix must be 1-8 characters long and contain only letters, numbers and hyphens",
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

					"streaming_data_client_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
					},

					"eu_location_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
						ForceNew: true,
					},

					"episodic_data_upload_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
					},

					"bitlocker_boot_volume_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
					},

					"bitlocker_data_volume_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
					},

					"credential_guard_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
						ForceNew: true,
					},

					"drift_control_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
					},

					"drtm_protection_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
					},

					"hvci_protection_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
					},

					"side_channel_mitigation_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
					},

					"smb_signing_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
					},

					"smb_cluster_encryption_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
						ForceNew: true,
					},

					"wdac_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
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
						ScaleUnits: expandDeploymentSettingScaleUnits(config.ScaleUnit),
					},
				},
			}

			// do validation
			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("validating %s: %+v", id, err)
			}

			// do deployment
			payload.Properties.DeploymentMode = deploymentsettings.DeploymentModeDeploy
			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("deploying %s: %+v", id, err)
			}

			metadata.SetID(id)

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
					schema.ScaleUnit = flattenDeploymentSettingScaleUnits(props.DeploymentConfiguration.ScaleUnits)
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

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			applianceName := fmt.Sprintf("%s-arcbridge", id.ClusterName)
			applianceId := appliances.NewApplianceID(id.SubscriptionId, id.ResourceGroupName, applianceName)

			log.Printf("[DEBUG] Deleting Arc Resource Bridge Appliance generated during deployment: %s", applianceId.ID())

			applianceClient := metadata.Client.ArcResourceBridge.AppliancesClient
			if err := applianceClient.DeleteThenPoll(ctx, applianceId); err != nil {
				return fmt.Errorf("deleting Arc Resource Bridge Appliance generated during deployment: deleting %s: %+v", applianceId, err)
			}

			log.Printf("[DEBUG] Deleting Custom Location and Stack HCI Storage Paths generated during deployment")

			var customLocationName string
			if resp.Model != nil && resp.Model.Properties != nil &&
				len(resp.Model.Properties.DeploymentConfiguration.ScaleUnits) > 0 &&
				resp.Model.Properties.DeploymentConfiguration.ScaleUnits[0].DeploymentData.OptionalServices != nil {
				customLocationName = pointer.From(resp.Model.Properties.DeploymentConfiguration.ScaleUnits[0].DeploymentData.OptionalServices.CustomLocation)
			}

			if customLocationName != "" {
				customLocationId := customlocations.NewCustomLocationID(id.SubscriptionId, id.ResourceGroupName, customLocationName)

				// we need to delete the Storage Paths before the Custom Location, otherwise the Custom Location cannot be deleted if there is any Resource in it
				log.Printf("[DEBUG] Deleting Stack HCI Storage Paths under Custom Location %s", customLocationId.ID())

				storageContainerClient := metadata.Client.AzureStackHCI.StorageContainers
				resourceGroupId := commonids.NewResourceGroupID(id.SubscriptionId, id.ResourceGroupName)
				storageContainers, err := storageContainerClient.ListComplete(ctx, resourceGroupId)
				if err != nil {
					return fmt.Errorf("deleting Stack HCI Storage Paths generated during deployment: retrieving Storage Path under %s: %+v", resourceGroupId, err)
				}

				// match Storage Paths under the Custom Location, the generated Storage Path name should match below pattern
				storageContainerNamePattern := regexp.MustCompile(`UserStorage[0-9]+-[a-z0-9]{32}`)
				for _, v := range storageContainers.Items {
					if v.Id != nil && v.ExtendedLocation != nil && v.ExtendedLocation.Name != nil && strings.EqualFold(*v.ExtendedLocation.Name, customLocationId.ID()) && v.Name != nil && storageContainerNamePattern.Match([]byte(*v.Name)) {
						storageContainerId, err := storagecontainers.ParseStorageContainerIDInsensitively(*v.Id)
						if err != nil {
							return fmt.Errorf("parsing the Stack HCI Storage Path ID generated during deployment: %+v", err)
						}

						if err := storageContainerClient.DeleteThenPoll(ctx, *storageContainerId); err != nil {
							return fmt.Errorf("deleting the Stack HCI Storage Path generated during deployment: deleting %s: %+v", storageContainerId, err)
						}
					}
				}

				customLocationsClient := metadata.Client.ExtendedLocation.CustomLocationsClient
				if err := customLocationsClient.DeleteThenPoll(ctx, customLocationId); err != nil {
					return fmt.Errorf("deleting the Custom Location generated during deployment: deleting %s: %+v", customLocationId, err)
				}
			}

			return nil
		},
	}
}

func expandDeploymentSettingScaleUnits(input []ScaleUnitModel) []deploymentsettings.ScaleUnits {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.ScaleUnits, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.ScaleUnits{
			DeploymentData: deploymentsettings.DeploymentData{
				AdouPath:              pointer.To(item.ActiveDirectoryOrganizationalUnitPath),
				Cluster:               expandDeploymentSettingCluster(item.Cluster),
				DomainFqdn:            pointer.To(item.DomainFqdn),
				HostNetwork:           expandDeploymentSettingHostNetwork(item.HostNetwork),
				InfrastructureNetwork: expandDeploymentSettingInfrastructureNetwork(item.InfrastructureNetwork),
				NamingPrefix:          pointer.To(item.NamePrefix),
				Observability:         expandDeploymentSettingObservability(item),
				OptionalServices:      expandDeploymentSettingOptionalService(item.OptionalService),
				PhysicalNodes:         expandDeploymentSettingPhysicalNode(item.PhysicalNode),
				SecretsLocation:       pointer.To(item.SecretsLocation),
				SecuritySettings:      expandDeploymentSettingSecuritySetting(item),
				Storage:               expandDeploymentSettingStorage(item.Storage),
			},
		})
	}

	return results
}

func flattenDeploymentSettingScaleUnits(input []deploymentsettings.ScaleUnits) []ScaleUnitModel {
	if len(input) == 0 {
		return make([]ScaleUnitModel, 0)
	}

	results := make([]ScaleUnitModel, 0, len(input))
	for _, item := range input {
		result := ScaleUnitModel{
			ActiveDirectoryOrganizationalUnitPath: pointer.From(item.DeploymentData.AdouPath),
			Cluster:                               flattenDeploymentSettingCluster(item.DeploymentData.Cluster),
			DomainFqdn:                            pointer.From(item.DeploymentData.DomainFqdn),
			HostNetwork:                           flattenDeploymentSettingHostNetwork(item.DeploymentData.HostNetwork),
			InfrastructureNetwork:                 flattenDeploymentSettingInfrastructureNetwork(item.DeploymentData.InfrastructureNetwork),
			NamePrefix:                            pointer.From(item.DeploymentData.NamingPrefix),
			OptionalService:                       flattenDeploymentSettingOptionalService(item.DeploymentData.OptionalServices),
			PhysicalNode:                          flattenDeploymentSettingPhysicalNode(item.DeploymentData.PhysicalNodes),
			SecretsLocation:                       pointer.From(item.DeploymentData.SecretsLocation),
			Storage:                               flattenDeploymentSettingStorage(item.DeploymentData.Storage),
		}

		if observability := item.DeploymentData.Observability; observability != nil {
			result.EpisodicDataUploadEnabled = pointer.From(observability.EpisodicDataUpload)
			result.EuLocationEnabled = pointer.From(observability.EuLocation)
			result.StreamingDataClientEnabled = pointer.From(observability.StreamingDataClient)
		}

		if securitySettings := item.DeploymentData.SecuritySettings; securitySettings != nil {
			result.BitlockerBootVolumeEnabled = pointer.From(securitySettings.BitlockerBootVolume)
			result.BitlockerDataVolumeEnabled = pointer.From(securitySettings.BitlockerDataVolumes)
			result.CredentialGuardEnabled = pointer.From(securitySettings.CredentialGuardEnforced)
			result.DriftControlEnabled = pointer.From(securitySettings.DriftControlEnforced)
			result.DrtmProtectionEnabled = pointer.From(securitySettings.DrtmProtection)
			result.HvciProtectionEnabled = pointer.From(securitySettings.HvciProtection)
			result.SideChannelMitigationEnabled = pointer.From(securitySettings.SideChannelMitigationEnforced)
			result.SmbClusterEncryptionEnabled = pointer.From(securitySettings.SmbClusterEncryption)
			result.SmbSigningEnabled = pointer.From(securitySettings.SmbSigningEnforced)
			result.WdacEnabled = pointer.From(securitySettings.WdacEnforced)
		}

		results = append(results, result)
	}

	return results
}

func expandDeploymentSettingCluster(input []ClusterModel) *deploymentsettings.DeploymentCluster {
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

func flattenDeploymentSettingCluster(input *deploymentsettings.DeploymentCluster) []ClusterModel {
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

func expandDeploymentSettingHostNetwork(input []HostNetworkModel) *deploymentsettings.HostNetwork {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.HostNetwork{
		Intents:                       expandDeploymentSettingHostNetworkIntent(v.Intent),
		EnableStorageAutoIP:           pointer.To(v.StorageAutoIpEnabled),
		StorageConnectivitySwitchless: pointer.To(v.StorageConnectivitySwitchlessEnabled),
		StorageNetworks:               expandDeploymentSettingHostNetworkStorageNetwork(v.StorageNetwork),
	}
}

func flattenDeploymentSettingHostNetwork(input *deploymentsettings.HostNetwork) []HostNetworkModel {
	if input == nil {
		return make([]HostNetworkModel, 0)
	}

	return []HostNetworkModel{{
		Intent:                               flattenDeploymentSettingHostNetworkIntent(input.Intents),
		StorageAutoIpEnabled:                 pointer.From(input.EnableStorageAutoIP),
		StorageConnectivitySwitchlessEnabled: pointer.From(input.StorageConnectivitySwitchless),
		StorageNetwork:                       flattenDeploymentSettingHostNetworkStorageNetwork(input.StorageNetworks),
	}}
}

func expandDeploymentSettingHostNetworkIntent(input []HostNetworkIntentModel) *[]deploymentsettings.Intents {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.Intents, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.Intents{
			Adapter:                             pointer.To(item.Adapter),
			AdapterPropertyOverrides:            expandHostNetworkIntentAdapterPropertyOverride(item.AdapterPropertyOverride),
			Name:                                pointer.To(item.Name),
			OverrideAdapterProperty:             pointer.To(item.AdapterPropertyOverrideEnabled),
			OverrideQosPolicy:                   pointer.To(item.QosPolicyOverrideEnabled),
			OverrideVirtualSwitchConfiguration:  pointer.To(item.VirtualSwitchConfigurationOverrideEnabled),
			QosPolicyOverrides:                  expandHostNetworkIntentQosPolicyOverride(item.QosPolicyOverride),
			TrafficType:                         pointer.To(item.TrafficType),
			VirtualSwitchConfigurationOverrides: expandHostNetworkIntentVirtualSwitchConfigurationOverride(item.VirtualSwitchConfigurationOverride),
		})
	}

	return &results
}

func flattenDeploymentSettingHostNetworkIntent(input *[]deploymentsettings.Intents) []HostNetworkIntentModel {
	if input == nil {
		return make([]HostNetworkIntentModel, 0)
	}

	results := make([]HostNetworkIntentModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, HostNetworkIntentModel{
			Adapter:                        pointer.From(item.Adapter),
			AdapterPropertyOverride:        flattenHostNetworkIntentAdapterPropertyOverride(item.AdapterPropertyOverrides),
			Name:                           pointer.From(item.Name),
			AdapterPropertyOverrideEnabled: pointer.From(item.OverrideAdapterProperty),
			QosPolicyOverrideEnabled:       pointer.From(item.OverrideQosPolicy),
			VirtualSwitchConfigurationOverrideEnabled: pointer.From(item.OverrideVirtualSwitchConfiguration),
			QosPolicyOverride:                         flattenHostNetworkIntentQosPolicyOverride(item.QosPolicyOverrides),
			TrafficType:                               pointer.From(item.TrafficType),
			VirtualSwitchConfigurationOverride:        flattenHostNetworkIntentVirtualSwitchConfigurationOverride(item.VirtualSwitchConfigurationOverrides),
		})
	}

	return results
}

func expandHostNetworkIntentAdapterPropertyOverride(input []HostNetworkIntentAdapterPropertyOverrideModel) *deploymentsettings.AdapterPropertyOverrides {
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

func flattenHostNetworkIntentAdapterPropertyOverride(input *deploymentsettings.AdapterPropertyOverrides) []HostNetworkIntentAdapterPropertyOverrideModel {
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

func expandHostNetworkIntentQosPolicyOverride(input []HostNetworkIntentQosPolicyOverrideModel) *deploymentsettings.QosPolicyOverrides {
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

func flattenHostNetworkIntentQosPolicyOverride(input *deploymentsettings.QosPolicyOverrides) []HostNetworkIntentQosPolicyOverrideModel {
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

func expandHostNetworkIntentVirtualSwitchConfigurationOverride(input []HostNetworkIntentVirtualSwitchConfigurationOverrideModel) *deploymentsettings.VirtualSwitchConfigurationOverrides {
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

func flattenHostNetworkIntentVirtualSwitchConfigurationOverride(input *deploymentsettings.VirtualSwitchConfigurationOverrides) []HostNetworkIntentVirtualSwitchConfigurationOverrideModel {
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

func expandDeploymentSettingHostNetworkStorageNetwork(input []HostNetworkStorageNetworkModel) *[]deploymentsettings.StorageNetworks {
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

func flattenDeploymentSettingHostNetworkStorageNetwork(input *[]deploymentsettings.StorageNetworks) []HostNetworkStorageNetworkModel {
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

func expandDeploymentSettingInfrastructureNetwork(input []InfrastructureNetworkModel) *[]deploymentsettings.InfrastructureNetwork {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.InfrastructureNetwork, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.InfrastructureNetwork{
			DnsServers: pointer.To(item.DnsServer),
			Gateway:    pointer.To(item.Gateway),
			IPPools:    expandDeploymentSettingInfrastructureNetworkIpPool(item.IpPool),
			SubnetMask: pointer.To(item.SubnetMask),
			UseDhcp:    pointer.To(item.DhcpEnabled),
		})
	}

	return &results
}

func flattenDeploymentSettingInfrastructureNetwork(input *[]deploymentsettings.InfrastructureNetwork) []InfrastructureNetworkModel {
	if input == nil {
		return make([]InfrastructureNetworkModel, 0)
	}

	results := make([]InfrastructureNetworkModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, InfrastructureNetworkModel{
			DhcpEnabled: pointer.From(item.UseDhcp),
			DnsServer:   pointer.From(item.DnsServers),
			Gateway:     pointer.From(item.Gateway),
			IpPool:      flattenDeploymentSettingInfrastructureNetworkIpPool(item.IPPools),
			SubnetMask:  pointer.From(item.SubnetMask),
		})
	}

	return results
}

func expandDeploymentSettingInfrastructureNetworkIpPool(input []IpPoolModel) *[]deploymentsettings.IPPools {
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

func flattenDeploymentSettingInfrastructureNetworkIpPool(input *[]deploymentsettings.IPPools) []IpPoolModel {
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

func expandDeploymentSettingObservability(input ScaleUnitModel) *deploymentsettings.Observability {
	return &deploymentsettings.Observability{
		EpisodicDataUpload:  pointer.To(input.EpisodicDataUploadEnabled),
		EuLocation:          pointer.To(input.EuLocationEnabled),
		StreamingDataClient: pointer.To(input.StreamingDataClientEnabled),
	}
}

func expandDeploymentSettingOptionalService(input []OptionalServiceModel) *deploymentsettings.OptionalServices {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.OptionalServices{
		CustomLocation: pointer.To(v.CustomLocation),
	}
}

func flattenDeploymentSettingOptionalService(input *deploymentsettings.OptionalServices) []OptionalServiceModel {
	if input == nil {
		return make([]OptionalServiceModel, 0)
	}

	return []OptionalServiceModel{{
		CustomLocation: pointer.From(input.CustomLocation),
	}}
}

func expandDeploymentSettingPhysicalNode(input []PhysicalNodeModel) *[]deploymentsettings.PhysicalNodes {
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

func flattenDeploymentSettingPhysicalNode(input *[]deploymentsettings.PhysicalNodes) []PhysicalNodeModel {
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

func expandDeploymentSettingSecuritySetting(input ScaleUnitModel) *deploymentsettings.DeploymentSecuritySettings {
	return &deploymentsettings.DeploymentSecuritySettings{
		BitlockerBootVolume:           pointer.To(input.BitlockerBootVolumeEnabled),
		BitlockerDataVolumes:          pointer.To(input.BitlockerDataVolumeEnabled),
		CredentialGuardEnforced:       pointer.To(input.CredentialGuardEnabled),
		DriftControlEnforced:          pointer.To(input.DriftControlEnabled),
		DrtmProtection:                pointer.To(input.DrtmProtectionEnabled),
		HvciProtection:                pointer.To(input.HvciProtectionEnabled),
		SideChannelMitigationEnforced: pointer.To(input.SideChannelMitigationEnabled),
		SmbClusterEncryption:          pointer.To(input.SmbClusterEncryptionEnabled),
		SmbSigningEnforced:            pointer.To(input.SmbSigningEnabled),
		WdacEnforced:                  pointer.To(input.WdacEnabled),
	}
}

func expandDeploymentSettingStorage(input []StorageModel) *deploymentsettings.Storage {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.Storage{
		ConfigurationMode: pointer.To(v.ConfigurationMode),
	}
}

func flattenDeploymentSettingStorage(input *deploymentsettings.Storage) []StorageModel {
	if input == nil {
		return make([]StorageModel, 0)
	}

	return []StorageModel{{
		ConfigurationMode: pointer.From(input.ConfigurationMode),
	}}
}
