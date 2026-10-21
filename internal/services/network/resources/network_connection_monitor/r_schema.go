// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_connection_monitor

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectionmonitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkwatchers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceNetworkConnectionMonitorSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_watcher_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkwatchers.ValidateNetworkWatcherID,
		},

		"location": commonschema.Location(),

		"endpoint": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.Any(
							validation.IsIPv4Address,
							validate.NetworkConnectionMonitorEndpointAddress,
						),
					},

					"coverage_level": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForCoverageLevel(), false),
					},

					"excluded_ip_addresses": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								validation.IsIPv4Address,
								validation.IsIPv6Address,
								validation.IsCIDR,
							),
						},
					},

					"filter": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"item": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"address": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: azure.ValidateResourceID,
											},

											"type": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												Default:      string(connectionmonitors.ConnectionMonitorEndpointFilterItemTypeAgentAddress),
												ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForConnectionMonitorEndpointFilterItemType(), false),
											},
										},
									},
								},

								"type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(connectionmonitors.ConnectionMonitorEndpointFilterTypeInclude),
									ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForConnectionMonitorEndpointFilterType(), false),
								},
							},
						},
					},

					"included_ip_addresses": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								validation.IsIPv4Address,
								validation.IsIPv6Address,
								validation.IsCIDR,
							),
						},
					},

					"target_resource_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.Any(
							commonids.ValidateVirtualMachineID,
							workspaces.ValidateWorkspaceID,
							commonids.ValidateSubnetID,
							commonids.ValidateVirtualNetworkID,
							machines.ValidateMachineID,
						),
					},

					"target_resource_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(connectionmonitors.EndpointTypeAzureArcVM),
							string(connectionmonitors.EndpointTypeAzureSubnet),
							string(connectionmonitors.EndpointTypeAzureVM),
							string(connectionmonitors.EndpointTypeAzureVNet),
							string(connectionmonitors.EndpointTypeExternalAddress),
							string(connectionmonitors.EndpointTypeMMAWorkspaceMachine),
							string(connectionmonitors.EndpointTypeMMAWorkspaceNetwork),
						}, false),
					},
				},
			},
		},

		"test_configuration": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"protocol": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForConnectionMonitorTestConfigurationProtocol(), false),
					},

					"http_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"method": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(connectionmonitors.HTTPConfigurationMethodGet),
									ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForHTTPConfigurationMethod(), false),
								},

								"path": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validate.NetworkConnectionMonitorHttpPath,
								},

								"port": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IsPortNumber,
								},

								"prefer_https": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"request_header": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"value": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"valid_status_code_ranges": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validate.NetworkConnectionMonitorValidStatusCodeRanges,
									},
								},
							},
						},
					},

					"icmp_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"trace_route_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},
							},
						},
					},

					"preferred_ip_version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForPreferredIPVersion(), false),
					},

					// lintignore:XS003
					"success_threshold": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"checks_failed_percent": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(0, 100),
								},

								"round_trip_time_ms": {
									Type:         pluginsdk.TypeFloat,
									Optional:     true,
									ValidateFunc: validation.FloatAtLeast(0),
								},
							},
						},
					},

					"tcp_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"port": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IsPortNumber,
								},

								"trace_route_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},

								"destination_port_behavior": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForDestinationPortBehavior(), false),
								},
							},
						},
					},

					"test_frequency_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      60,
						ValidateFunc: validation.IntBetween(30, 1800),
					},
				},
			},
		},

		"test_group": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"destination_endpoints": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"source_endpoints": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"test_configuration_names": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		// API accepts any value including empty string.
		"notes": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"output_workspace_resource_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},
		},

		"tags": commonschema.Tags(),
	}
}
