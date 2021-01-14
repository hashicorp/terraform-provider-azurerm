package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	computeValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	logAnalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetworkConnectionMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkConnectionMonitorCreateUpdate,
		Read:   resourceArmNetworkConnectionMonitorRead,
		Update: resourceArmNetworkConnectionMonitorCreateUpdate,
		Delete: resourceArmNetworkConnectionMonitorDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"network_watcher_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.NetworkWatcherID,
			},

			"location": azure.SchemaLocation(),

			"auto_start": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "The field belongs to the v1 network connection monitor, which is now deprecated in favour of v2 by Azure. Please check the document (https://www.terraform.io/docs/providers/azurerm/r/network_connection_monitor.html) for the v2 properties.",
			},

			"interval_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(30),
				Deprecated:   "The field belongs to the v1 network connection monitor, which is now deprecated in favour of v2 by Azure. Please check the document (https://www.terraform.io/docs/providers/azurerm/r/network_connection_monitor.html) for the v2 properties.",
			},

			"source": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				Deprecated: "The field belongs to the v1 network connection monitor, which is now deprecated in favour of v2 by Azure. Please check the document (https://www.terraform.io/docs/providers/azurerm/r/network_connection_monitor.html) for the v2 properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_machine_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: azure.ValidateResourceID,
							Deprecated:   "The field belongs to the v1 network connection monitor, which is now deprecated in favour of v2 by Azure. Please check the document (https://www.terraform.io/docs/providers/azurerm/r/network_connection_monitor.html) for the v2 properties.",
						},

						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.PortNumberOrZero,
							Deprecated:   "The field belongs to the v1 network connection monitor, which is now deprecated in favour of v2 by Azure. Please check the document (https://www.terraform.io/docs/providers/azurerm/r/network_connection_monitor.html) for the v2 properties.",
						},
					},
				},
			},

			"destination": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				Deprecated: "The field belongs to the v1 network connection monitor, which is now deprecated in favour of v2 by Azure. Please check the document (https://www.terraform.io/docs/providers/azurerm/r/network_connection_monitor.html) for the v2 properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_machine_id": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ValidateFunc:  azure.ValidateResourceID,
							ConflictsWith: []string{"destination.0.address"},
							Deprecated:    "The field belongs to the v1 network connection monitor, which is now deprecated in favour of v2 by Azure. Please check the document (https://www.terraform.io/docs/providers/azurerm/r/network_connection_monitor.html) for the v2 properties.",
						},

						"address": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"destination.0.virtual_machine_id"},
							Deprecated:    "The field belongs to the v1 network connection monitor, which is now deprecated in favour of v2 by Azure. Please check the document (https://www.terraform.io/docs/providers/azurerm/r/network_connection_monitor.html) for the v2 properties.",
						},

						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.PortNumber,
							Deprecated:   "The field belongs to the v1 network connection monitor, which is now deprecated in favour of v2 by Azure. Please check the document (https://www.terraform.io/docs/providers/azurerm/r/network_connection_monitor.html) for the v2 properties.",
						},
					},
				},
			},

			"endpoint": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"address": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.Any(
								validation.IsIPv4Address,
								networkValidate.NetworkConnectionMonitorEndpointAddress,
							),
						},

						"filter": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"item": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: azure.ValidateResourceID,
												},

												"type": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  string(network.AgentAddress),
													ValidateFunc: validation.StringInSlice([]string{
														string(network.AgentAddress),
													}, false),
												},
											},
										},
									},

									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(network.Include),
										ValidateFunc: validation.StringInSlice([]string{
											string(network.Include),
										}, false),
									},
								},
							},
						},

						"virtual_machine_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: computeValidate.VirtualMachineID,
						},
					},
				},
			},

			"test_configuration": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ConnectionMonitorTestConfigurationProtocolTCP),
								string(network.ConnectionMonitorTestConfigurationProtocolHTTP),
								string(network.ConnectionMonitorTestConfigurationProtocolIcmp),
							}, false),
						},

						"http_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(network.Get),
										ValidateFunc: validation.StringInSlice([]string{
											string(network.Get),
											string(network.Post),
										}, false),
									},

									"path": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: networkValidate.NetworkConnectionMonitorHttpPath,
									},

									"port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validate.PortNumber,
									},

									"prefer_https": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"request_header": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"value": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},

									"valid_status_code_ranges": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: networkValidate.NetworkConnectionMonitorValidStatusCodeRanges,
										},
									},
								},
							},
						},

						"icmp_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"trace_route_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},

						"preferred_ip_version": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.PreferredIPVersionIPv4),
								string(network.PreferredIPVersionIPv6),
							}, false),
						},

						"success_threshold": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"checks_failed_percent": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 100),
									},

									"round_trip_time_ms": {
										Type:         schema.TypeFloat,
										Optional:     true,
										ValidateFunc: validation.FloatAtLeast(0),
									},
								},
							},
						},

						"tcp_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validate.PortNumber,
									},

									"trace_route_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},

						"test_frequency_in_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      60,
							ValidateFunc: validation.IntBetween(30, 1800),
						},
					},
				},
			},

			"test_group": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"destination_endpoints": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"source_endpoints": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"test_configuration_names": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			// API accepts any value including empty string.
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"output_workspace_resource_ids": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: logAnalyticsValidate.LogAnalyticsWorkspaceID,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNetworkConnectionMonitorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ConnectionMonitorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))

	watcherId := d.Get("network_watcher_id").(string)
	id, err := parse.NetworkWatcherID(watcherId)
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Connection Monitor %q (Watcher %q / Resource Group %q): %s", name, id.Name, id.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_network_connection_monitor", *existing.ID)
		}
	}

	properties := network.ConnectionMonitor{
		Location: utils.String(location),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		ConnectionMonitorParameters: &network.ConnectionMonitorParameters{
			Endpoints:          expandArmNetworkConnectionMonitorEndpoint(d.Get("endpoint").(*schema.Set).List()),
			Outputs:            expandArmNetworkConnectionMonitorOutput(d.Get("output_workspace_resource_ids").(*schema.Set).List()),
			TestConfigurations: expandArmNetworkConnectionMonitorTestConfiguration(d.Get("test_configuration").(*schema.Set).List()),
			TestGroups:         expandArmNetworkConnectionMonitorTestGroup(d.Get("test_group").(*schema.Set).List()),
		},
	}

	if notes, ok := d.GetOk("notes"); ok {
		properties.Notes = utils.String(notes.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, name, properties)
	if err != nil {
		return fmt.Errorf("Error creating Connection Monitor %q (Watcher %q / Resource Group %q): %+v", name, id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Connection Monitor %q (Watcher %q / Resource Group %q): %+v", name, id.Name, id.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Connection Monitor %q (Watcher %q / Resource Group %q): %+v", name, id.Name, id.ResourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Connection Monitor %q (Watcher %q / Resource Group %q) ID", name, id.Name, id.ResourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmNetworkConnectionMonitorRead(d, meta)
}

func resourceArmNetworkConnectionMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ConnectionMonitorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConnectionMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NetworkWatcherName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Connection Monitor %q (Watcher %q / Resource Group %q) %+v", id.Name, id.NetworkWatcherName, id.ResourceGroup, err)
	}

	if resp.ConnectionMonitorType == network.SingleSourceDestination {
		return fmt.Errorf("the resource created via API version 2019-06-01 or before (a.k.a v1) isn't compatible to this version of provider. Please migrate to v2 resource.")
	}

	d.Set("name", id.Name)

	networkWatcherId := parse.NewNetworkWatcherID(id.SubscriptionId, id.ResourceGroup, id.NetworkWatcherName)
	d.Set("network_watcher_id", networkWatcherId.ID())

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ConnectionMonitorResultProperties; props != nil {
		d.Set("notes", props.Notes)

		if err := d.Set("endpoint", flattenArmNetworkConnectionMonitorEndpoint(props.Endpoints)); err != nil {
			return fmt.Errorf("setting `endpoint`: %+v", err)
		}

		if err := d.Set("output_workspace_resource_ids", flattenArmNetworkConnectionMonitorOutput(props.Outputs)); err != nil {
			return fmt.Errorf("setting `output`: %+v", err)
		}

		if err := d.Set("test_configuration", flattenArmNetworkConnectionMonitorTestConfiguration(props.TestConfigurations)); err != nil {
			return fmt.Errorf("setting `test_configuration`: %+v", err)
		}

		if err := d.Set("test_group", flattenArmNetworkConnectionMonitorTestGroup(props.TestGroups)); err != nil {
			return fmt.Errorf("setting `test_group`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmNetworkConnectionMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ConnectionMonitorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConnectionMonitorID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.NetworkWatcherName, id.Name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Connection Monitor %q (Watcher %q / Resource Group %q): %+v", id.Name, id.NetworkWatcherName, id.ResourceGroup, err)
		}
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Connection Monitor %q (Watcher %q / Resource Group %q): %+v", id.Name, id.NetworkWatcherName, id.ResourceGroup, err)
	}

	return nil
}

func expandArmNetworkConnectionMonitorEndpoint(input []interface{}) *[]network.ConnectionMonitorEndpoint {
	results := make([]network.ConnectionMonitorEndpoint, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := network.ConnectionMonitorEndpoint{
			Name:   utils.String(v["name"].(string)),
			Filter: expandArmNetworkConnectionMonitorEndpointFilter(v["filter"].([]interface{})),
		}

		if address := v["address"]; address != "" {
			result.Address = utils.String(address.(string))
		}

		if resourceId := v["virtual_machine_id"]; resourceId != "" {
			result.ResourceID = utils.String(resourceId.(string))
		}

		results = append(results, result)
	}

	return &results
}

func expandArmNetworkConnectionMonitorEndpointFilter(input []interface{}) *network.ConnectionMonitorEndpointFilter {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &network.ConnectionMonitorEndpointFilter{
		Type:  network.ConnectionMonitorEndpointFilterType(v["type"].(string)),
		Items: expandArmNetworkConnectionMonitorEndpointFilterItem(v["item"].(*schema.Set).List()),
	}
}

func expandArmNetworkConnectionMonitorEndpointFilterItem(input []interface{}) *[]network.ConnectionMonitorEndpointFilterItem {
	if len(input) == 0 {
		return nil
	}

	results := make([]network.ConnectionMonitorEndpointFilterItem, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := network.ConnectionMonitorEndpointFilterItem{
			Type: network.ConnectionMonitorEndpointFilterItemType(v["type"].(string)),
		}

		if address := v["address"]; address != "" {
			result.Address = utils.String(address.(string))
		}

		results = append(results, result)
	}

	return &results
}

func expandArmNetworkConnectionMonitorTestConfiguration(input []interface{}) *[]network.ConnectionMonitorTestConfiguration {
	results := make([]network.ConnectionMonitorTestConfiguration, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := network.ConnectionMonitorTestConfiguration{
			Name:              utils.String(v["name"].(string)),
			HTTPConfiguration: expandArmNetworkConnectionMonitorHTTPConfiguration(v["http_configuration"].([]interface{})),
			IcmpConfiguration: expandArmNetworkConnectionMonitorIcmpConfiguration(v["icmp_configuration"].([]interface{})),
			Protocol:          network.ConnectionMonitorTestConfigurationProtocol(v["protocol"].(string)),
			SuccessThreshold:  expandArmNetworkConnectionMonitorSuccessThreshold(v["success_threshold"].([]interface{})),
			TCPConfiguration:  expandArmNetworkConnectionMonitorTCPConfiguration(v["tcp_configuration"].([]interface{})),
			TestFrequencySec:  utils.Int32(int32(v["test_frequency_in_seconds"].(int))),
		}

		if preferredIPVersion := v["preferred_ip_version"]; preferredIPVersion != "" {
			result.PreferredIPVersion = network.PreferredIPVersion(preferredIPVersion.(string))
		}

		results = append(results, result)
	}

	return &results
}

func expandArmNetworkConnectionMonitorHTTPConfiguration(input []interface{}) *network.ConnectionMonitorHTTPConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	props := &network.ConnectionMonitorHTTPConfiguration{
		Method:         network.HTTPConfigurationMethod(v["method"].(string)),
		PreferHTTPS:    utils.Bool(v["prefer_https"].(bool)),
		RequestHeaders: expandArmNetworkConnectionMonitorHTTPHeader(v["request_header"].(*schema.Set).List()),
	}

	if path := v["path"]; path != "" {
		props.Path = utils.String(path.(string))
	}

	if port := v["port"]; port != 0 {
		props.Port = utils.Int32(int32(port.(int)))
	}

	if ranges := v["valid_status_code_ranges"].(*schema.Set).List(); len(ranges) != 0 {
		props.ValidStatusCodeRanges = utils.ExpandStringSlice(ranges)
	}

	return props
}

func expandArmNetworkConnectionMonitorTCPConfiguration(input []interface{}) *network.ConnectionMonitorTCPConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &network.ConnectionMonitorTCPConfiguration{
		Port:              utils.Int32(int32(v["port"].(int))),
		DisableTraceRoute: utils.Bool(!v["trace_route_enabled"].(bool)),
	}
}

func expandArmNetworkConnectionMonitorIcmpConfiguration(input []interface{}) *network.ConnectionMonitorIcmpConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &network.ConnectionMonitorIcmpConfiguration{
		DisableTraceRoute: utils.Bool(!v["trace_route_enabled"].(bool)),
	}
}

func expandArmNetworkConnectionMonitorSuccessThreshold(input []interface{}) *network.ConnectionMonitorSuccessThreshold {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &network.ConnectionMonitorSuccessThreshold{
		ChecksFailedPercent: utils.Int32(int32(v["checks_failed_percent"].(int))),
		RoundTripTimeMs:     utils.Float(v["round_trip_time_ms"].(float64)),
	}
}

func expandArmNetworkConnectionMonitorHTTPHeader(input []interface{}) *[]network.HTTPHeader {
	if len(input) == 0 {
		return nil
	}

	results := make([]network.HTTPHeader, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := network.HTTPHeader{
			Name:  utils.String(v["name"].(string)),
			Value: utils.String(v["value"].(string)),
		}

		results = append(results, result)
	}

	return &results
}

func expandArmNetworkConnectionMonitorTestGroup(input []interface{}) *[]network.ConnectionMonitorTestGroup {
	results := make([]network.ConnectionMonitorTestGroup, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := network.ConnectionMonitorTestGroup{
			Name:               utils.String(v["name"].(string)),
			Destinations:       utils.ExpandStringSlice(v["destination_endpoints"].(*schema.Set).List()),
			Disable:            utils.Bool(!v["enabled"].(bool)),
			Sources:            utils.ExpandStringSlice(v["source_endpoints"].(*schema.Set).List()),
			TestConfigurations: utils.ExpandStringSlice(v["test_configuration_names"].(*schema.Set).List()),
		}

		results = append(results, result)
	}

	return &results
}

func expandArmNetworkConnectionMonitorOutput(input []interface{}) *[]network.ConnectionMonitorOutput {
	results := make([]network.ConnectionMonitorOutput, 0)

	for _, item := range input {
		result := network.ConnectionMonitorOutput{
			Type: network.Workspace,
			WorkspaceSettings: &network.ConnectionMonitorWorkspaceSettings{
				WorkspaceResourceID: utils.String(item.(string)),
			},
		}

		results = append(results, result)
	}

	return &results
}

func flattenArmNetworkConnectionMonitorEndpoint(input *[]network.ConnectionMonitorEndpoint) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		var address string
		if item.Address != nil {
			address = *item.Address
		}

		var resourceId string
		if item.ResourceID != nil {
			resourceId = *item.ResourceID
		}

		v := map[string]interface{}{
			"name":               name,
			"address":            address,
			"filter":             flattenArmNetworkConnectionMonitorEndpointFilter(item.Filter),
			"virtual_machine_id": resourceId,
		}

		results = append(results, v)
	}
	return results
}

func flattenArmNetworkConnectionMonitorEndpointFilter(input *network.ConnectionMonitorEndpointFilter) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var t network.ConnectionMonitorEndpointFilterType
	if input.Type != "" {
		t = input.Type
	}
	return []interface{}{
		map[string]interface{}{
			"item": flattenArmNetworkConnectionMonitorEndpointFilterItem(input.Items),
			"type": t,
		},
	}
}

func flattenArmNetworkConnectionMonitorEndpointFilterItem(input *[]network.ConnectionMonitorEndpointFilterItem) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var address string
		if item.Address != nil {
			address = *item.Address
		}

		var t network.ConnectionMonitorEndpointFilterItemType
		if item.Type != "" {
			t = item.Type
		}

		v := map[string]interface{}{
			"address": address,
			"type":    t,
		}

		results = append(results, v)
	}

	return results
}

func flattenArmNetworkConnectionMonitorTestConfiguration(input *[]network.ConnectionMonitorTestConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		var protocol network.ConnectionMonitorTestConfigurationProtocol
		if item.Protocol != "" {
			protocol = item.Protocol
		}

		var preferredIpVersion network.PreferredIPVersion
		if item.PreferredIPVersion != "" {
			preferredIpVersion = item.PreferredIPVersion
		}

		var testFrequencySec int32
		if item.TestFrequencySec != nil {
			testFrequencySec = *item.TestFrequencySec
		}

		v := map[string]interface{}{
			"name":                      name,
			"protocol":                  protocol,
			"http_configuration":        flattenArmNetworkConnectionMonitorHTTPConfiguration(item.HTTPConfiguration),
			"icmp_configuration":        flattenArmNetworkConnectionMonitorIcmpConfiguration(item.IcmpConfiguration),
			"preferred_ip_version":      preferredIpVersion,
			"success_threshold":         flattenArmNetworkConnectionMonitorSuccessThreshold(item.SuccessThreshold),
			"tcp_configuration":         flattenArmNetworkConnectionMonitorTCPConfiguration(item.TCPConfiguration),
			"test_frequency_in_seconds": testFrequencySec,
		}

		results = append(results, v)
	}

	return results
}

func flattenArmNetworkConnectionMonitorHTTPConfiguration(input *network.ConnectionMonitorHTTPConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var method network.HTTPConfigurationMethod
	if input.Method != "" {
		method = input.Method
	}

	var p string
	if input.Path != nil {
		p = *input.Path
	}

	var port int32
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
			"request_header":           flattenArmNetworkConnectionMonitorHTTPHeader(input.RequestHeaders),
			"valid_status_code_ranges": utils.FlattenStringSlice(input.ValidStatusCodeRanges),
		},
	}
}

func flattenArmNetworkConnectionMonitorIcmpConfiguration(input *network.ConnectionMonitorIcmpConfiguration) []interface{} {
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

func flattenArmNetworkConnectionMonitorSuccessThreshold(input *network.ConnectionMonitorSuccessThreshold) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var checksFailedPercent int32
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

func flattenArmNetworkConnectionMonitorTCPConfiguration(input *network.ConnectionMonitorTCPConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var enableTraceRoute bool
	if input.DisableTraceRoute != nil {
		enableTraceRoute = !*input.DisableTraceRoute
	}

	var port int32
	if input.Port != nil {
		port = *input.Port
	}

	return []interface{}{
		map[string]interface{}{
			"trace_route_enabled": enableTraceRoute,
			"port":                port,
		},
	}
}

func flattenArmNetworkConnectionMonitorHTTPHeader(input *[]network.HTTPHeader) []interface{} {
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

func flattenArmNetworkConnectionMonitorTestGroup(input *[]network.ConnectionMonitorTestGroup) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		var disable bool
		if item.Disable != nil {
			disable = *item.Disable
		}

		v := map[string]interface{}{
			"name":                     name,
			"destination_endpoints":    utils.FlattenStringSlice(item.Destinations),
			"source_endpoints":         utils.FlattenStringSlice(item.Sources),
			"test_configuration_names": utils.FlattenStringSlice(item.TestConfigurations),
			"enabled":                  !disable,
		}

		results = append(results, v)
	}
	return results
}

func flattenArmNetworkConnectionMonitorOutput(input *[]network.ConnectionMonitorOutput) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var workspaceResourceId string
		if item.WorkspaceSettings != nil && item.WorkspaceSettings.WorkspaceResourceID != nil {
			workspaceResourceId = *item.WorkspaceSettings.WorkspaceResourceID
		}

		results = append(results, workspaceResourceId)
	}

	return results
}
