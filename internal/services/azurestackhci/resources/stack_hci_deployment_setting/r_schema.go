// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_hci_deployment_setting

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"stack_hci_cluster_id": commonschema.ResourceIDReferenceRequiredForceNew(&deploymentsettings.ClusterId{}),

		"arc_resource_ids": {
			Description: "A list of IDs of Azure ARC machine resource to be part of cluster.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: machines.ValidateMachineID,
			},
		},

		"version": {
			Description: "The deployment template version. The format must be a set of numbers separated by dots such as `10.0.0.0`.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
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
						Description:  "The full name of the Active Directory Organizational Unit container object prepared for the deployment, including the domain components. For example: `OU=HCI01,DC=contoso,DC=com`.",
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
									Description: "The name of the cluster.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ForceNew:    true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile("^[a-zA-Z0-9-]{3,15}$"),
										"the cluster name must be 3-15 characters long and contain only letters, numbers and hyphens",
									),
								},

								"azure_service_endpoint": {
									Description: "The Azure blob service endpoint, for example, `core.windows.net`.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ForceNew:    true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`),
										"`azure_service_endpoint` must be a valid domain name, for example, \"core.windows.net\"",
									),
								},

								"cloud_account_name": {
									Description:  "The Azure Storage Account name of the cloud witness for the Azure Stack HCI cluster.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validate.StorageAccountName,
								},

								"witness_type": {
									Description: "The type of the witness.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ForceNew:    true,
									ValidateFunc: validation.StringInSlice([]string{
										"Cloud",
										"FileShare",
									}, false),
								},

								"witness_path": {
									Description:  "The fileshare path of the local witness for the Azure Stack HCI cluster.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"domain_fqdn": {
						Description: "The FQDN for deploying cluster.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
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
												Description:  "The name of the intent.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"adapter": {
												Description: "A list of ID of network interfaces used for the network intent.",
												Type:        pluginsdk.TypeList,
												Required:    true,
												ForceNew:    true,
												MinItems:    1,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},

											"traffic_type": {
												Description: "A list of network traffic types.",
												Type:        pluginsdk.TypeList,
												Required:    true,
												ForceNew:    true,
												MinItems:    1,
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
															Description:  "The jumbo frame size of the adapter. This parameter should only be modified based on your OEM guidance.",
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"network_direct": {
															Description:  "The network direct of the adapter. This parameter should only be modified based on your OEM guidance.",
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"network_direct_technology": {
															Description:  "The network direct technology of the adapter. This parameter should only be modified based on your OEM guidance.",
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},

											"adapter_property_override_enabled": {
												Description: "Whether to override adapter properties.",
												Type:        pluginsdk.TypeBool,
												Optional:    true,
												ForceNew:    true,
												Default:     false,
											},

											"qos_policy_override": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												ForceNew: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"bandwidth_percentage_smb": {
															Description:  "The percentage of the allocated storage traffic bandwidth. This parameter should only be modified based on your OEM guidance.",
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"priority_value8021_action_cluster": {
															Description:  "The Cluster traffic priority. This parameter should only be modified based on your OEM guidance.",
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"priority_value8021_action_smb": {
															Description:  "The Priority Flow Control where Data Center Bridging (DCB) is used. This parameter should only be modified based on your OEM guidance.",
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},

											"qos_policy_override_enabled": {
												Description: "Whether to override QoS policy.",
												Type:        pluginsdk.TypeBool,
												Optional:    true,
												ForceNew:    true,
												Default:     false,
											},

											"virtual_switch_configuration_override": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												ForceNew: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"enable_iov": {
															Description:  "The IoV enable status for Virtual Switch.",
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"load_balancing_algorithm": {
															Description:  "The load balancing algorithm for Virtual Switch.",
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},

											"virtual_switch_configuration_override_enabled": {
												Description: "Whether to override virtual switch configuration.",
												Type:        pluginsdk.TypeBool,
												Optional:    true,
												ForceNew:    true,
												Default:     false,
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
												Description:  "The name of the storage network.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"network_adapter_name": {
												Description:  "The name of the network adapter.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"vlan_id": {
												Description:  "The ID for the VLAN storage network. This setting is applied to the network interfaces that route the storage and VM migration traffic.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"storage_auto_ip_enabled": {
									Description: "Whether allows users to specify IPs and Mask for Storage NICs when Network ATC is not assigning the IPs for storage automatically. Optional parameter required only for [3 nodes switchless deployments](https://learn.microsoft.com/azure-stack/hci/concepts/physical-network-requirements?tabs=overview%2C23H2reqs#using-switchless).",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									ForceNew:    true,
									Default:     true,
								},

								"storage_connectivity_switchless_enabled": {
									Description: "How the storage adapters between nodes are connected either switch or switch less.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									ForceNew:    true,
									Default:     false,
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
									Description: "A list of IPv4 addresses of the DNS servers in your environment.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									ForceNew:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsIPv4Address,
									},
								},

								"gateway": {
									Description:  "The default gateway that should be used for the provided IP address space.",
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
												Description:  "The ending IP address for the management network. A minimum of six free, contiguous IPv4 addresses (excluding your host IPs) are needed for infrastructure services such as clustering.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.IsIPv4Address,
											},

											"ending_address": {
												Description:  "The starting IP address for the management network. A minimum of six free, contiguous IPv4 addresses (excluding your host IPs) are needed for infrastructure services such as clustering.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.IsIPv4Address,
											},
										},
									},
								},

								"subnet_mask": {
									Description:  "The subnet mask that matches the provided IP address space.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},

								"dhcp_enabled": {
									Description: "Whether DHCP is enabled for hosts and cluster IPs.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									ForceNew:    true,
									Default:     false,
								},
							},
						},
					},

					"name_prefix": {
						Description: "The name prefix to deploy cluster.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
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
									Description:  "The name of custom location. A custom location will be created after the deployment is completed.",
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
									Description:  "The NETBIOS name of each physical server on your Azure Stack HCI cluster.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"ipv4_address": {
									Description:  "The IPv4 address assigned to each physical server on your Azure Stack HCI cluster.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},
							},
						},
					},

					"secrets_location": {
						Description:  "The URI to the Key Vault or secret store.",
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
									Description: "The configuration mode of storage. If set to `Express` and your storage is configured as per best practices based on the number of nodes in the cluster.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ForceNew:    true,
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
						Description: "Whether the telemetry data will be sent to Microsoft.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						ForceNew:    true,
					},

					"eu_location_enabled": {
						Description: "Whether to store data sent to Microsoft in EU. The log and diagnostic data is sent to the appropriate diagnostics servers depending upon where your cluster resides. Setting this to `false` results in all data sent to Microsoft to be stored outside of the EU.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
						ForceNew:    true,
					},

					"episodic_data_upload_enabled": {
						Description: "Whether to collect log data to facilitate quicker issue resolution.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						ForceNew:    true,
					},

					"bitlocker_boot_volume_enabled": {
						Description: "Whether to enable BitLocker for boot volume. When set to `true`, BitLocker XTS_AES 256-bit encryption is enabled for all data-at-rest on the OS volume of your Azure Stack HCI cluster. This setting is TPM-hardware dependent.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						ForceNew:    true,
					},

					"bitlocker_data_volume_enabled": {
						Description: "Whether to enable BitLocker for data volume. When set to `true`, BitLocker XTS-AES 256-bit encryption is enabled for all data-at-rest on your Azure Stack HCI cluster shared volumes.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						ForceNew:    true,
					},

					"credential_guard_enabled": {
						Description: "Whether to enable credential guard.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
						ForceNew:    true,
					},

					"drift_control_enabled": {
						Description: "Whether to enable drift control. When set to `true`, the security baseline is re-applied regularly.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						ForceNew:    true,
					},

					"drtm_protection_enabled": {
						Description: "Whether to enable DRTM protection. When set to `true`, Secure Boot is enabled on your Azure HCI cluster. This setting is hardware dependent.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						ForceNew:    true,
					},

					"hvci_protection_enabled": {
						Description: "Whether to enable HVCI protection. When set to `true`, Hypervisor-protected Code Integrity is enabled on your Azure HCI cluster.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						ForceNew:    true,
					},

					"side_channel_mitigation_enabled": {
						Description: "Whether to enable side channel mitigation. When set to `true`, all side channel mitigations are enabled on your Azure HCI cluster.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						ForceNew:    true,
					},

					"smb_signing_enabled": {
						Description: "Whether to enable SMB signing. When set to `true`, the SMB default instance requires sign in for the client and server services.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						ForceNew:    true,
					},

					"smb_cluster_encryption_enabled": {
						Description: "Whether to enable SMB cluster encryption. When set to `true`, cluster east-west traffic is encrypted.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
						ForceNew:    true,
					},

					"wdac_enabled": {
						Description: "Whether to enable WDAC. When set to `true`, applications and the code that you can run on your Azure Stack HCI cluster are limited.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						ForceNew:    true,
					},
				},
			},
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
