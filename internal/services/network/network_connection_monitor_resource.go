// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/connectionmonitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkwatchers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetworkConnectionMonitor() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkConnectionMonitorCreateUpdate,
		Read:   resourceNetworkConnectionMonitorRead,
		Update: resourceNetworkConnectionMonitorCreateUpdate,
		Delete: resourceNetworkConnectionMonitorDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := connectionmonitors.ParseConnectionMonitorID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceNetworkConnectionMonitorSchema(),
	}
}

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
							networkValidate.NetworkConnectionMonitorEndpointAddress,
						),
					},

					"coverage_level": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(connectionmonitors.CoverageLevelAboveAverage),
							string(connectionmonitors.CoverageLevelAverage),
							string(connectionmonitors.CoverageLevelBelowAverage),
							string(connectionmonitors.CoverageLevelDefault),
							string(connectionmonitors.CoverageLevelFull),
							string(connectionmonitors.CoverageLevelLow),
						}, false),
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
												Type:     pluginsdk.TypeString,
												Optional: true,
												Default:  string(connectionmonitors.ConnectionMonitorEndpointFilterItemTypeAgentAddress),
												ValidateFunc: validation.StringInSlice([]string{
													string(connectionmonitors.ConnectionMonitorEndpointFilterItemTypeAgentAddress),
												}, false),
											},
										},
									},
								},

								"type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  string(connectionmonitors.ConnectionMonitorEndpointFilterTypeInclude),
									ValidateFunc: validation.StringInSlice([]string{
										string(connectionmonitors.ConnectionMonitorEndpointFilterTypeInclude),
									}, false),
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
						Computed: true,
						ValidateFunc: validation.Any(
							computeValidate.VirtualMachineID,
							workspaces.ValidateWorkspaceID,
							commonids.ValidateSubnetID,
							commonids.ValidateVirtualNetworkID,
						),
					},

					"target_resource_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
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
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(connectionmonitors.ConnectionMonitorTestConfigurationProtocolTcp),
							string(connectionmonitors.ConnectionMonitorTestConfigurationProtocolHTTP),
							string(connectionmonitors.ConnectionMonitorTestConfigurationProtocolIcmp),
						}, false),
					},

					"http_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"method": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  string(connectionmonitors.HTTPConfigurationMethodGet),
									ValidateFunc: validation.StringInSlice([]string{
										string(connectionmonitors.HTTPConfigurationMethodGet),
										string(connectionmonitors.HTTPConfigurationMethodPost),
									}, false),
								},

								"path": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: networkValidate.NetworkConnectionMonitorHttpPath,
								},

								"port": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validate.PortNumber,
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
										ValidateFunc: networkValidate.NetworkConnectionMonitorValidStatusCodeRanges,
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
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(connectionmonitors.PreferredIPVersionIPvFour),
							string(connectionmonitors.PreferredIPVersionIPvSix),
						}, false),
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
									ValidateFunc: validate.PortNumber,
								},

								"trace_route_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},

								"destination_port_behavior": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(connectionmonitors.DestinationPortBehaviorNone),
										string(connectionmonitors.DestinationPortBehaviorListenIfAvailable),
									}, false),
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
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},
		},

		"tags": commonschema.Tags(),
	}
}

func resourceNetworkConnectionMonitorCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ConnectionMonitors
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location").(string))

	watcherId, err := connectionmonitors.ParseNetworkWatcherID(d.Get("network_watcher_id").(string))
	if err != nil {
		return err
	}

	connectionMonitorId := connectionmonitors.NewConnectionMonitorID(subscriptionId, watcherId.ResourceGroupName, watcherId.NetworkWatcherName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, connectionMonitorId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", connectionMonitorId, err)
			}
		}

		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_network_connection_monitor", connectionMonitorId.ID())
		}
	}

	properties := connectionmonitors.ConnectionMonitor{
		Location: utils.String(location),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: connectionmonitors.ConnectionMonitorParameters{
			Outputs:            expandNetworkConnectionMonitorOutput(d.Get("output_workspace_resource_ids").(*pluginsdk.Set).List()),
			TestConfigurations: expandNetworkConnectionMonitorTestConfiguration(d.Get("test_configuration").(*pluginsdk.Set).List()),
			TestGroups:         expandNetworkConnectionMonitorTestGroup(d.Get("test_group").(*pluginsdk.Set).List()),
		},
	}

	properties.Properties.Endpoints = expandNetworkConnectionMonitorEndpoint(d.Get("endpoint").(*pluginsdk.Set).List())

	if notes, ok := d.GetOk("notes"); ok {
		properties.Properties.Notes = utils.String(notes.(string))
	}

	if err = client.CreateOrUpdateThenPoll(ctx, connectionMonitorId, properties, connectionmonitors.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %+v", connectionMonitorId, err)
	}

	d.SetId(connectionMonitorId.ID())

	return resourceNetworkConnectionMonitorRead(d, meta)
}

func resourceNetworkConnectionMonitorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ConnectionMonitors
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connectionmonitors.ParseConnectionMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectionMonitorName)

	if model := resp.Model; model != nil {
		networkWatcherId := networkwatchers.NewNetworkWatcherID(id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName)
		d.Set("network_watcher_id", networkWatcherId.ID())

		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if props := model.Properties; props != nil {
			if props.ConnectionMonitorType != nil && *props.ConnectionMonitorType == connectionmonitors.ConnectionMonitorTypeSingleSourceDestination {
				return fmt.Errorf("the resource created via API version 2019-06-01 or before (a.k.a v1) isn't compatible to this version of provider. Please migrate to v2 pluginsdk.")
			}
			d.Set("notes", props.Notes)

			if err := d.Set("endpoint", flattenNetworkConnectionMonitorEndpoint(props.Endpoints)); err != nil {
				return fmt.Errorf("setting `endpoint`: %+v", err)
			}

			if err := d.Set("output_workspace_resource_ids", flattenNetworkConnectionMonitorOutput(props.Outputs)); err != nil {
				return fmt.Errorf("setting `output`: %+v", err)
			}

			if err := d.Set("test_configuration", flattenNetworkConnectionMonitorTestConfiguration(props.TestConfigurations)); err != nil {
				return fmt.Errorf("setting `test_configuration`: %+v", err)
			}

			if err := d.Set("test_group", flattenNetworkConnectionMonitorTestGroup(props.TestGroups)); err != nil {
				return fmt.Errorf("setting `test_group`: %+v", err)
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceNetworkConnectionMonitorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ConnectionMonitors
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connectionmonitors.ParseConnectionMonitorID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandNetworkConnectionMonitorEndpoint(input []interface{}) *[]connectionmonitors.ConnectionMonitorEndpoint {
	results := make([]connectionmonitors.ConnectionMonitorEndpoint, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := connectionmonitors.ConnectionMonitorEndpoint{
			Name:   v["name"].(string),
			Filter: expandNetworkConnectionMonitorEndpointFilter(v["filter"].([]interface{})),
		}

		if address := v["address"]; address != "" {
			result.Address = utils.String(address.(string))
		}

		if coverageLevel := v["coverage_level"]; coverageLevel != "" {
			result.CoverageLevel = pointer.To(connectionmonitors.CoverageLevel(coverageLevel.(string)))
		}

		excludedItems := v["excluded_ip_addresses"].(*pluginsdk.Set).List()
		includedItems := v["included_ip_addresses"].(*pluginsdk.Set).List()
		if len(excludedItems) != 0 || len(includedItems) != 0 {
			result.Scope = &connectionmonitors.ConnectionMonitorEndpointScope{}

			if len(excludedItems) != 0 {
				var excludedAddresses []connectionmonitors.ConnectionMonitorEndpointScopeItem
				for _, v := range excludedItems {
					excludedAddresses = append(excludedAddresses, connectionmonitors.ConnectionMonitorEndpointScopeItem{
						Address: utils.String(v.(string)),
					})
				}
				result.Scope.Exclude = &excludedAddresses
			}

			if len(includedItems) != 0 {
				var includedAddresses []connectionmonitors.ConnectionMonitorEndpointScopeItem
				for _, v := range includedItems {
					includedAddresses = append(includedAddresses, connectionmonitors.ConnectionMonitorEndpointScopeItem{
						Address: utils.String(v.(string)),
					})
				}
				result.Scope.Include = &includedAddresses
			}
		}

		if resourceId := v["target_resource_id"]; resourceId != "" {
			result.ResourceId = utils.String(resourceId.(string))
		}

		if endpointType := v["target_resource_type"]; endpointType != "" {
			result.Type = pointer.To(connectionmonitors.EndpointType(endpointType.(string)))
		}

		results = append(results, result)
	}

	return &results
}

func expandNetworkConnectionMonitorEndpointFilter(input []interface{}) *connectionmonitors.ConnectionMonitorEndpointFilter {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &connectionmonitors.ConnectionMonitorEndpointFilter{
		Type:  pointer.To(connectionmonitors.ConnectionMonitorEndpointFilterType(v["type"].(string))),
		Items: expandNetworkConnectionMonitorEndpointFilterItem(v["item"].(*pluginsdk.Set).List()),
	}
}

func expandNetworkConnectionMonitorEndpointFilterItem(input []interface{}) *[]connectionmonitors.ConnectionMonitorEndpointFilterItem {
	if len(input) == 0 {
		return nil
	}

	results := make([]connectionmonitors.ConnectionMonitorEndpointFilterItem, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := connectionmonitors.ConnectionMonitorEndpointFilterItem{
			Type: pointer.To(connectionmonitors.ConnectionMonitorEndpointFilterItemType(v["type"].(string))),
		}

		if address := v["address"]; address != "" {
			result.Address = utils.String(address.(string))
		}

		results = append(results, result)
	}

	return &results
}

func expandNetworkConnectionMonitorTestConfiguration(input []interface{}) *[]connectionmonitors.ConnectionMonitorTestConfiguration {
	results := make([]connectionmonitors.ConnectionMonitorTestConfiguration, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := connectionmonitors.ConnectionMonitorTestConfiguration{
			Name:              v["name"].(string),
			HTTPConfiguration: expandNetworkConnectionMonitorHTTPConfiguration(v["http_configuration"].([]interface{})),
			IcmpConfiguration: expandNetworkConnectionMonitorIcmpConfiguration(v["icmp_configuration"].([]interface{})),
			Protocol:          connectionmonitors.ConnectionMonitorTestConfigurationProtocol(v["protocol"].(string)),
			SuccessThreshold:  expandNetworkConnectionMonitorSuccessThreshold(v["success_threshold"].([]interface{})),
			TcpConfiguration:  expandNetworkConnectionMonitorTCPConfiguration(v["tcp_configuration"].([]interface{})),
			TestFrequencySec:  utils.Int64(int64(v["test_frequency_in_seconds"].(int))),
		}

		if preferredIPVersion := v["preferred_ip_version"]; preferredIPVersion != "" {
			result.PreferredIPVersion = pointer.To(connectionmonitors.PreferredIPVersion(preferredIPVersion.(string)))
		}

		results = append(results, result)
	}

	return &results
}

func expandNetworkConnectionMonitorHTTPConfiguration(input []interface{}) *connectionmonitors.ConnectionMonitorHTTPConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	props := &connectionmonitors.ConnectionMonitorHTTPConfiguration{
		Method:         pointer.To(connectionmonitors.HTTPConfigurationMethod(v["method"].(string))),
		PreferHTTPS:    utils.Bool(v["prefer_https"].(bool)),
		RequestHeaders: expandNetworkConnectionMonitorHTTPHeader(v["request_header"].(*pluginsdk.Set).List()),
	}

	if path := v["path"]; path != "" {
		props.Path = utils.String(path.(string))
	}

	if port := v["port"]; port != 0 {
		props.Port = utils.Int64(int64(port.(int)))
	}

	if ranges := v["valid_status_code_ranges"].(*pluginsdk.Set).List(); len(ranges) != 0 {
		props.ValidStatusCodeRanges = utils.ExpandStringSlice(ranges)
	}

	return props
}

func expandNetworkConnectionMonitorTCPConfiguration(input []interface{}) *connectionmonitors.ConnectionMonitorTcpConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := &connectionmonitors.ConnectionMonitorTcpConfiguration{
		Port:              utils.Int64(int64(v["port"].(int))),
		DisableTraceRoute: utils.Bool(!v["trace_route_enabled"].(bool)),
	}

	if destinationPortBehavior := v["destination_port_behavior"].(string); destinationPortBehavior != "" {
		result.DestinationPortBehavior = pointer.To(connectionmonitors.DestinationPortBehavior(destinationPortBehavior))
	}

	return result
}

func expandNetworkConnectionMonitorIcmpConfiguration(input []interface{}) *connectionmonitors.ConnectionMonitorIcmpConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &connectionmonitors.ConnectionMonitorIcmpConfiguration{
		DisableTraceRoute: utils.Bool(!v["trace_route_enabled"].(bool)),
	}
}

func expandNetworkConnectionMonitorSuccessThreshold(input []interface{}) *connectionmonitors.ConnectionMonitorSuccessThreshold {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &connectionmonitors.ConnectionMonitorSuccessThreshold{
		ChecksFailedPercent: utils.Int64(int64(v["checks_failed_percent"].(int))),
		RoundTripTimeMs:     utils.Float(v["round_trip_time_ms"].(float64)),
	}
}

func expandNetworkConnectionMonitorHTTPHeader(input []interface{}) *[]connectionmonitors.HTTPHeader {
	if len(input) == 0 {
		return nil
	}

	results := make([]connectionmonitors.HTTPHeader, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := connectionmonitors.HTTPHeader{
			Name:  utils.String(v["name"].(string)),
			Value: utils.String(v["value"].(string)),
		}

		results = append(results, result)
	}

	return &results
}

func expandNetworkConnectionMonitorTestGroup(input []interface{}) *[]connectionmonitors.ConnectionMonitorTestGroup {
	results := make([]connectionmonitors.ConnectionMonitorTestGroup, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := connectionmonitors.ConnectionMonitorTestGroup{
			Name:               v["name"].(string),
			Destinations:       *utils.ExpandStringSlice(v["destination_endpoints"].(*pluginsdk.Set).List()),
			Disable:            utils.Bool(!v["enabled"].(bool)),
			Sources:            *utils.ExpandStringSlice(v["source_endpoints"].(*pluginsdk.Set).List()),
			TestConfigurations: *utils.ExpandStringSlice(v["test_configuration_names"].(*pluginsdk.Set).List()),
		}

		results = append(results, result)
	}

	return &results
}

func expandNetworkConnectionMonitorOutput(input []interface{}) *[]connectionmonitors.ConnectionMonitorOutput {
	results := make([]connectionmonitors.ConnectionMonitorOutput, 0)

	for _, item := range input {
		result := connectionmonitors.ConnectionMonitorOutput{
			Type: pointer.To(connectionmonitors.OutputTypeWorkspace),
			WorkspaceSettings: &connectionmonitors.ConnectionMonitorWorkspaceSettings{
				WorkspaceResourceId: utils.String(item.(string)),
			},
		}

		results = append(results, result)
	}

	return &results
}

func flattenNetworkConnectionMonitorEndpoint(input *[]connectionmonitors.ConnectionMonitorEndpoint) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var address string
		if item.Address != nil {
			address = *item.Address
		}

		var coverageLevel string
		if item.CoverageLevel != nil && string(*item.CoverageLevel) != "" {
			coverageLevel = string(*item.CoverageLevel)
		}

		var endpointType string
		if item.Type != nil && string(*item.Type) != "" {
			endpointType = string(*item.Type)
		}

		var resourceId string
		if item.ResourceId != nil {
			resourceId = *item.ResourceId
		}

		v := map[string]interface{}{
			"name":                 item.Name,
			"address":              address,
			"coverage_level":       coverageLevel,
			"target_resource_id":   resourceId,
			"target_resource_type": endpointType,
			"filter":               flattenNetworkConnectionMonitorEndpointFilter(item.Filter),
		}

		if scope := item.Scope; scope != nil {
			if includeScope := scope.Include; includeScope != nil {
				includedAddresses := make([]interface{}, 0)

				for _, includedItem := range *includeScope {
					if includedAddress := includedItem.Address; includedAddress != nil {
						includedAddresses = append(includedAddresses, includedAddress)
					}
				}

				v["included_ip_addresses"] = includedAddresses
			}

			if excludeScope := scope.Exclude; excludeScope != nil {
				excludedAddresses := make([]interface{}, 0)

				for _, excludedItem := range *excludeScope {
					if excludedAddress := excludedItem.Address; excludedAddress != nil {
						excludedAddresses = append(excludedAddresses, excludedAddress)
					}
				}

				v["excluded_ip_addresses"] = excludedAddresses
			}
		}

		results = append(results, v)
	}
	return results
}

func flattenNetworkConnectionMonitorEndpointFilter(input *connectionmonitors.ConnectionMonitorEndpointFilter) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var t connectionmonitors.ConnectionMonitorEndpointFilterType
	if input.Type != nil && string(*input.Type) != "" {
		t = *input.Type
	}
	return []interface{}{
		map[string]interface{}{
			"item": flattenNetworkConnectionMonitorEndpointFilterItem(input.Items),
			"type": t,
		},
	}
}

func flattenNetworkConnectionMonitorEndpointFilterItem(input *[]connectionmonitors.ConnectionMonitorEndpointFilterItem) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var address string
		if item.Address != nil {
			address = *item.Address
		}

		var t connectionmonitors.ConnectionMonitorEndpointFilterItemType
		if item.Type != nil && string(*item.Type) != "" {
			t = *item.Type
		}

		v := map[string]interface{}{
			"address": address,
			"type":    t,
		}

		results = append(results, v)
	}

	return results
}

func flattenNetworkConnectionMonitorTestConfiguration(input *[]connectionmonitors.ConnectionMonitorTestConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var protocol connectionmonitors.ConnectionMonitorTestConfigurationProtocol
		if item.Protocol != "" {
			protocol = item.Protocol
		}

		var preferredIpVersion connectionmonitors.PreferredIPVersion
		if item.PreferredIPVersion != nil && string(*item.PreferredIPVersion) != "" {
			preferredIpVersion = *item.PreferredIPVersion
		}

		var testFrequencySec int64
		if item.TestFrequencySec != nil {
			testFrequencySec = *item.TestFrequencySec
		}

		v := map[string]interface{}{
			"name":                      item.Name,
			"protocol":                  protocol,
			"http_configuration":        flattenNetworkConnectionMonitorHTTPConfiguration(item.HTTPConfiguration),
			"icmp_configuration":        flattenNetworkConnectionMonitorIcmpConfiguration(item.IcmpConfiguration),
			"preferred_ip_version":      preferredIpVersion,
			"success_threshold":         flattenNetworkConnectionMonitorSuccessThreshold(item.SuccessThreshold),
			"tcp_configuration":         flattenNetworkConnectionMonitorTCPConfiguration(item.TcpConfiguration),
			"test_frequency_in_seconds": testFrequencySec,
		}

		results = append(results, v)
	}

	return results
}

func flattenNetworkConnectionMonitorHTTPConfiguration(input *connectionmonitors.ConnectionMonitorHTTPConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var method connectionmonitors.HTTPConfigurationMethod
	if input.Method != nil && string(*input.Method) != "" {
		method = *input.Method
	}

	var p string
	if input.Path != nil {
		p = *input.Path
	}

	var port int64
	if input.Port != nil {
		port = *input.Port
	}

	var preferHttps bool
	if input.PreferHTTPS != nil {
		preferHttps = *input.PreferHTTPS
	}

	return []interface{}{
		map[string]interface{}{
			"method":                   method,
			"path":                     p,
			"port":                     port,
			"prefer_https":             preferHttps,
			"request_header":           flattenNetworkConnectionMonitorHTTPHeader(input.RequestHeaders),
			"valid_status_code_ranges": utils.FlattenStringSlice(input.ValidStatusCodeRanges),
		},
	}
}

func flattenNetworkConnectionMonitorIcmpConfiguration(input *connectionmonitors.ConnectionMonitorIcmpConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var enableTraceRoute bool
	if input.DisableTraceRoute != nil {
		enableTraceRoute = !*input.DisableTraceRoute
	}

	return []interface{}{
		map[string]interface{}{
			"trace_route_enabled": enableTraceRoute,
		},
	}
}

func flattenNetworkConnectionMonitorSuccessThreshold(input *connectionmonitors.ConnectionMonitorSuccessThreshold) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var checksFailedPercent int64
	if input.ChecksFailedPercent != nil {
		checksFailedPercent = *input.ChecksFailedPercent
	}

	var roundTripTimeMs float64
	if input.RoundTripTimeMs != nil {
		roundTripTimeMs = *input.RoundTripTimeMs
	}

	return []interface{}{
		map[string]interface{}{
			"checks_failed_percent": checksFailedPercent,
			"round_trip_time_ms":    roundTripTimeMs,
		},
	}
}

func flattenNetworkConnectionMonitorTCPConfiguration(input *connectionmonitors.ConnectionMonitorTcpConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var enableTraceRoute bool
	if input.DisableTraceRoute != nil {
		enableTraceRoute = !*input.DisableTraceRoute
	}

	var port int64
	if input.Port != nil {
		port = *input.Port
	}

	var destinationPortBehavior connectionmonitors.DestinationPortBehavior
	if input.DestinationPortBehavior != nil && string(*input.DestinationPortBehavior) != "" {
		destinationPortBehavior = *input.DestinationPortBehavior
	}

	return []interface{}{
		map[string]interface{}{
			"trace_route_enabled":       enableTraceRoute,
			"port":                      port,
			"destination_port_behavior": string(destinationPortBehavior),
		},
	}
}

func flattenNetworkConnectionMonitorHTTPHeader(input *[]connectionmonitors.HTTPHeader) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		var value string
		if item.Value != nil {
			value = *item.Value
		}

		v := map[string]interface{}{
			"name":  name,
			"value": value,
		}

		results = append(results, v)
	}

	return results
}

func flattenNetworkConnectionMonitorTestGroup(input *[]connectionmonitors.ConnectionMonitorTestGroup) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var disable bool
		if item.Disable != nil {
			disable = *item.Disable
		}

		v := map[string]interface{}{
			"name":                     item.Name,
			"destination_endpoints":    item.Destinations,
			"source_endpoints":         item.Sources,
			"test_configuration_names": item.TestConfigurations,
			"enabled":                  !disable,
		}

		results = append(results, v)
	}
	return results
}

func flattenNetworkConnectionMonitorOutput(input *[]connectionmonitors.ConnectionMonitorOutput) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var workspaceResourceId string
		if item.WorkspaceSettings != nil && item.WorkspaceSettings.WorkspaceResourceId != nil {
			workspaceResourceId = *item.WorkspaceSettings.WorkspaceResourceId
		}

		results = append(results, workspaceResourceId)
	}

	return results
}
