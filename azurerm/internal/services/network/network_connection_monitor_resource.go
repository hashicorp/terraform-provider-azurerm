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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/migration"
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

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    migration.NetworkConnectionMonitorV0Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.NetworkConnectionMonitorV0ToV1,
				Version: 0,
			},
		},

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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"network_watcher_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"auto_start": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "This field has been deprecated in new api version 2020-05-01",
			},

			"destination": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"endpoint"},
				Deprecated:    "Deprecated in favor of `endpoint`",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.PortNumber,
						},

						"virtual_machine_id": {
							Type:          schema.TypeString,
							Optional:      true,
							ValidateFunc:  computeValidate.VirtualMachineID,
							ConflictsWith: []string{"destination.0.address"},
						},

						"address": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"destination.0.virtual_machine_id"},
						},
					},
				},
			},

			"endpoint": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"source", "destination"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"address": {
							Type:     schema.TypeString,
							Optional: true,
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
												"type": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  string(network.AgentAddress),
													ValidateFunc: validation.StringInSlice([]string{
														string(network.AgentAddress),
													}, false),
												},

												"address": {
													Type:     schema.TypeString,
													Optional: true,
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

			"interval_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Deprecated:   "Deprecated in favor of `test_frequency_sec`",
				ValidateFunc: validation.IntAtLeast(30),
			},

			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"source": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"endpoint"},
				Deprecated:    "Deprecated in favor of `endpoint`",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_machine_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: computeValidate.VirtualMachineID,
						},

						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validate.PortNumberOrZero,
						},
					},
				},
			},

			"test_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
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

						"test_frequency_sec": {
							Type:          schema.TypeInt,
							Optional:      true,
							Default:       60,
							ValidateFunc:  validation.IntBetween(30, 1800),
							ConflictsWith: []string{"interval_in_seconds"},
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

									"port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validate.PortNumber,
									},

									"path": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: networkValidate.NetworkConnectionMonitorHttpPath,
									},

									"prefer_https": {
										Type:     schema.TypeBool,
										Optional: true,
									},

									"request_header": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},

												"value": {
													Type:     schema.TypeString,
													Required: true,
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
									"disable_trace_route": {
										Type:     schema.TypeBool,
										Optional: true,
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

									"disable_trace_route": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"test_group": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"destinations": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"sources": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"test_configurations": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"output_workspace_resource_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				// Allow to switch workspace from specified one to default one.
				// 1. Set `output = []` in tfconfig to switch workspace from specified one to default one.
				// 2. Remove `output = []` from tfconfig to ensure no diff.
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
	watcherName := d.Get("network_watcher_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, watcherName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Connection Monitor %q (Watcher %q / Resource Group %q): %s", name, watcherName, resourceGroup, err)
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
			Notes:              utils.String(d.Get("notes").(string)),
			Outputs:            expandArmNetworkConnectionMonitorOutput(d.Get("output_workspace_resource_ids").(*schema.Set).List()),
			TestConfigurations: expandArmNetworkConnectionMonitorTestConfiguration(d.Get("test_configuration").(*schema.Set).List()),
			TestGroups:         expandArmNetworkConnectionMonitorTestGroup(d.Get("test_group").(*schema.Set).List()),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, watcherName, name, properties)
	if err != nil {
		return fmt.Errorf("Error creating Connection Monitor %q (Watcher %q / Resource Group %q): %+v", name, watcherName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Connection Monitor %q (Watcher %q / Resource Group %q): %+v", name, watcherName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, watcherName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Connection Monitor %q (Watcher %q / Resource Group %q): %+v", name, watcherName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Connection Monitor %q (Watcher %q / Resource Group %q) ID", name, watcherName, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmNetworkConnectionMonitorRead(d, meta)
}

func resourceArmNetworkConnectionMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ConnectionMonitorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NetworkConnectionMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WatcherName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Connection Monitor %q (Watcher %q / Resource Group %q) %+v", id.Name, id.WatcherName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("network_watcher_name", id.WatcherName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ConnectionMonitorResultProperties; props != nil {
		d.Set("notes", props.Notes)

		if err := d.Set("output_workspace_resource_ids", flattenArmNetworkConnectionMonitorOutput(props.Outputs)); err != nil {
			return fmt.Errorf("setting `output`: %+v", err)
		}

		if err := d.Set("endpoint", flattenArmNetworkConnectionMonitorEndpoint(props.Endpoints)); err != nil {
			return fmt.Errorf("setting `endpoint`: %+v", err)
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

	id, err := parse.NetworkConnectionMonitorID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WatcherName, id.Name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Connection Monitor %q (Watcher %q / Resource Group %q): %+v", id.Name, id.WatcherName, id.ResourceGroup, err)
		}
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Connection Monitor %q (Watcher %q / Resource Group %q): %+v", id.Name, id.WatcherName, id.ResourceGroup, err)
	}

	return nil
}

func expandArmNetworkConnectionMonitorEndpoint(input []interface{}) *[]network.ConnectionMonitorEndpoint {
	results := make([]network.ConnectionMonitorEndpoint, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := network.ConnectionMonitorEndpoint{
			Name:    utils.String(v["name"].(string)),
			Address: utils.String(v["address"].(string)),
			Filter:  expandArmNetworkConnectionMonitorEndpointFilter(v["filter"].([]interface{})),
		}

		resourceId := v["virtual_machine_id"].(string)
		if resourceId != "" {
			result.ResourceID = utils.String(resourceId)
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
	results := make([]network.ConnectionMonitorEndpointFilterItem, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := network.ConnectionMonitorEndpointFilterItem{
			Type:    network.ConnectionMonitorEndpointFilterItemType(v["type"].(string)),
			Address: utils.String(v["address"].(string)),
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
			Name:               utils.String(v["name"].(string)),
			Protocol:           network.ConnectionMonitorTestConfigurationProtocol(v["protocol"].(string)),
			PreferredIPVersion: network.PreferredIPVersion(v["preferred_ip_version"].(string)),
			TestFrequencySec:   utils.Int32(int32(v["test_frequency_sec"].(int))),
			HTTPConfiguration:  expandArmNetworkConnectionMonitorHTTPConfiguration(v["http_configuration"].([]interface{})),
			TCPConfiguration:   expandArmNetworkConnectionMonitorTCPConfiguration(v["tcp_configuration"].([]interface{})),
			IcmpConfiguration:  expandArmNetworkConnectionMonitorIcmpConfiguration(v["icmp_configuration"].([]interface{})),
			SuccessThreshold:   expandArmNetworkConnectionMonitorSuccessThreshold(v["success_threshold"].([]interface{})),
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
		Method:                network.HTTPConfigurationMethod(v["method"].(string)),
		Path:                  utils.String(v["path"].(string)),
		RequestHeaders:        expandArmNetworkConnectionMonitorHTTPHeader(v["request_header"].(*schema.Set).List()),
		ValidStatusCodeRanges: utils.ExpandStringSlice(v["valid_status_code_ranges"].(*schema.Set).List()),
		PreferHTTPS:           utils.Bool(v["prefer_https"].(bool)),
	}

	if port := v["port"].(int); port != 0 {
		props.Port = utils.Int32(int32(port))
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
		DisableTraceRoute: utils.Bool(v["disable_trace_route"].(bool)),
	}
}

func expandArmNetworkConnectionMonitorIcmpConfiguration(input []interface{}) *network.ConnectionMonitorIcmpConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &network.ConnectionMonitorIcmpConfiguration{
		DisableTraceRoute: utils.Bool(v["disable_trace_route"].(bool)),
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
			Disable:            utils.Bool(v["disable"].(bool)),
			TestConfigurations: utils.ExpandStringSlice(v["test_configurations"].(*schema.Set).List()),
			Sources:            utils.ExpandStringSlice(v["sources"].(*schema.Set).List()),
			Destinations:       utils.ExpandStringSlice(v["destinations"].(*schema.Set).List()),
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
			"name":                 name,
			"protocol":             protocol,
			"http_configuration":   flattenArmNetworkConnectionMonitorHTTPConfiguration(item.HTTPConfiguration),
			"icmp_configuration":   flattenArmNetworkConnectionMonitorIcmpConfiguration(item.IcmpConfiguration),
			"preferred_ip_version": preferredIpVersion,
			"success_threshold":    flattenArmNetworkConnectionMonitorSuccessThreshold(item.SuccessThreshold),
			"tcp_configuration":    flattenArmNetworkConnectionMonitorTCPConfiguration(item.TCPConfiguration),
			"test_frequency_sec":   testFrequencySec,
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

	var disableTraceRoute bool
	if input.DisableTraceRoute != nil {
		disableTraceRoute = *input.DisableTraceRoute
	}

	return []interface{}{
		map[string]interface{}{
			"disable_trace_route": disableTraceRoute,
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

	var disableTraceRoute bool
	if input.DisableTraceRoute != nil {
		disableTraceRoute = *input.DisableTraceRoute
	}

	var port int32
	if input.Port != nil {
		port = *input.Port
	}

	return []interface{}{
		map[string]interface{}{
			"disable_trace_route": disableTraceRoute,
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
			"name":                name,
			"destinations":        utils.FlattenStringSlice(item.Destinations),
			"sources":             utils.FlattenStringSlice(item.Sources),
			"test_configurations": utils.FlattenStringSlice(item.TestConfigurations),
			"disable":             disable,
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
