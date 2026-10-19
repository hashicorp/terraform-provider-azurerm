// Copyright IBM Corp. 2014, 2025
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
									ValidateFunc: validate.StorageAccountName,
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

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
