// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name which should be used for this Network Connection Monitor.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_watcher_id": {
			Description:  "The ID of the Network Watcher.",
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
						Description:  "The name of the endpoint for the Network Connection Monitor.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address": {
						Description: "The IP address or domain name of the Network Connection Monitor endpoint.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.Any(
							validation.IsIPv4Address,
							validate.NetworkConnectionMonitorEndpointAddress,
						),
					},

					"coverage_level": {
						Description:  "The test coverage for the Network Connection Monitor endpoint.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForCoverageLevel(), false),
					},

					"excluded_ip_addresses": {
						Description: "A list of IPv4/IPv6 subnet masks or IPv4/IPv6 IP addresses to be excluded to the Network Connection Monitor endpoint.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
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
												Description:  "The address of the filter item.",
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: azure.ValidateResourceID,
											},

											"type": {
												Description:  "The type of items included in the filter.",
												Type:         pluginsdk.TypeString,
												Optional:     true,
												Default:      string(connectionmonitors.ConnectionMonitorEndpointFilterItemTypeAgentAddress),
												ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForConnectionMonitorEndpointFilterItemType(), false),
											},
										},
									},
								},

								"type": {
									Description:  "The behaviour type of this endpoint filter. Currently the only.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(connectionmonitors.ConnectionMonitorEndpointFilterTypeInclude),
									ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForConnectionMonitorEndpointFilterType(), false),
								},
							},
						},
					},

					"included_ip_addresses": {
						Description: "A list of IPv4/IPv6 subnet masks or IPv4/IPv6 IP addresses to be included to the Network Connection Monitor endpoint.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
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
						Description: "The resource ID which is used as the endpoint by the Network Connection Monitor.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.Any(
							commonids.ValidateVirtualMachineID,
							workspaces.ValidateWorkspaceID,
							commonids.ValidateSubnetID,
							commonids.ValidateVirtualNetworkID,
							machines.ValidateMachineID,
						),
					},

					"target_resource_type": {
						Description: "The endpoint type of the Network Connection Monitor.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
						Description:  "The name of test configuration for the Network Connection Monitor.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"protocol": {
						Description:  "The protocol used to evaluate tests.",
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
									Description:  "The HTTP method for the HTTP request.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(connectionmonitors.HTTPConfigurationMethodGet),
									ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForHTTPConfigurationMethod(), false),
								},

								"path": {
									Description:  "The path component of the URI. It only accepts the absolute path.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validate.NetworkConnectionMonitorHttpPath,
								},

								"port": {
									Description:  "The port for the HTTP connection.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IsPortNumber,
								},

								"prefer_https": {
									Description: "Whether HTTPS should be preferred over HTTP in cases where the choice is not explicit.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},

								"request_header": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Description:  "The name of the HTTP header.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"value": {
												Description:  "The value of the HTTP header.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"valid_status_code_ranges": {
									Description: "The HTTP status codes to consider successful. For instance, `2xx`, `301-304` and `418`.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
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
									Description: "Whether path evaluation with trace route should be enabled.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     true,
								},
							},
						},
					},

					"preferred_ip_version": {
						Description:  "The preferred IP version which is used in the test evaluation.",
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
									Description:  "The maximum percentage of failed checks permitted for a test to be successful.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(0, 100),
								},

								"round_trip_time_ms": {
									Description:  "The maximum round-trip time in milliseconds permitted for a test to be successful.",
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
									Description:  "The port for the TCP connection.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IsPortNumber,
								},

								"trace_route_enabled": {
									Description: "Whether path evaluation with trace route should be enabled.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     true,
								},

								"destination_port_behavior": {
									Description:  "The destination port behaviour for the TCP connection.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(connectionmonitors.PossibleValuesForDestinationPortBehavior(), false),
								},
							},
						},
					},

					"test_frequency_in_seconds": {
						Description:  "The time interval in seconds at which the test evaluation will happen.",
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
						Description:  "The name of the test group for the Network Connection Monitor.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"destination_endpoints": {
						Description: "A list of destination endpoint names.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"source_endpoints": {
						Description: "A list of source endpoint names.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"test_configuration_names": {
						Description: "A list of test configuration names.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"enabled": {
						Description: "Whether the test group should be enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
				},
			},
		},

		// API accepts any value including empty string.
		"notes": {
			Description: "The description of the Network Connection Monitor.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"output_workspace_resource_ids": {
			Description: "A list of IDs of the Log Analytics Workspace which will accept the output from the Network Connection Monitor.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},
		},

		"tags": commonschema.Tags(),
	}
}
