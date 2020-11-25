package trafficmanager

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-04-01/trafficmanager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmTrafficManagerProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmTrafficManagerProfileCreate,
		Read:   resourceArmTrafficManagerProfileRead,
		Update: resourceArmTrafficManagerProfileUpdate,
		Delete: resourceArmTrafficManagerProfileDelete,
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

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"profile_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(trafficmanager.ProfileStatusEnabled),
					string(trafficmanager.ProfileStatusDisabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"traffic_routing_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(trafficmanager.Geographic),
					string(trafficmanager.Weighted),
					string(trafficmanager.Performance),
					string(trafficmanager.Priority),
					string(trafficmanager.Subnet),
					string(trafficmanager.MultiValue),
				}, false),
			},

			"dns_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"relative_name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"ttl": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 999999),
						},
					},
				},
			},

			"monitor_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expected_status_code_ranges": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateTrafficManagerProfileStatusCodeRange,
							},
						},

						"custom_header": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(trafficmanager.HTTP),
								string(trafficmanager.HTTPS),
								string(trafficmanager.TCP),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"interval_in_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{10, 30}),
							Default:      30,
						},

						"timeout_in_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(5, 10),
							Default:      10,
						},

						"tolerated_number_of_failures": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 9),
							Default:      3,
						},
					},
				},
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"max_return": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 8),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmTrafficManagerProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for TrafficManager Profile creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing TrafficManager profile %s (resource group %s) ID", name, resourceGroup)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_traffic_manager_profile", *existing.ID)
	}

	// No existing profile - start from a new struct.
	profile := trafficmanager.Profile{
		Name:     &name,
		Location: utils.String("global"), // must be provided in request
		ProfileProperties: &trafficmanager.ProfileProperties{
			TrafficRoutingMethod: trafficmanager.TrafficRoutingMethod(d.Get("traffic_routing_method").(string)),
			DNSConfig:            expandArmTrafficManagerDNSConfig(d),
			MonitorConfig:        expandArmTrafficManagerMonitorConfig(d),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if maxReturn, ok := d.GetOk("max_return"); ok {
		maxReturnInt := int64(maxReturn.(int))
		profile.ProfileProperties.MaxReturn = &maxReturnInt
	}

	if status, ok := d.GetOk("profile_status"); ok {
		profile.ProfileStatus = trafficmanager.ProfileStatus(status.(string))
	}

	if profile.ProfileProperties.TrafficRoutingMethod == trafficmanager.MultiValue &&
		profile.ProfileProperties.MaxReturn == nil {
		return fmt.Errorf("`max_return` must be specified when `traffic_routing_method` is set to `MultiValue`")
	}

	if *profile.ProfileProperties.MonitorConfig.IntervalInSeconds == int64(10) &&
		*profile.ProfileProperties.MonitorConfig.TimeoutInSeconds == int64(10) {
		return fmt.Errorf("`timeout_in_seconds` must be between `5` and `9` when `interval_in_seconds` is set to `10`")
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, profile); err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error reading TrafficManager profile %s (resource group %s): %v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read TrafficManager profile %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmTrafficManagerProfileRead(d, meta)
}

func resourceArmTrafficManagerProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["trafficManagerProfiles"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Traffic Manager Profile %s: %+v", name, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("name", resp.Name)

	if profile := resp.ProfileProperties; profile != nil {
		d.Set("profile_status", profile.ProfileStatus)
		d.Set("traffic_routing_method", profile.TrafficRoutingMethod)
		d.Set("max_return", profile.MaxReturn)

		d.Set("dns_config", flattenAzureRMTrafficManagerProfileDNSConfig(profile.DNSConfig))
		d.Set("monitor_config", flattenAzureRMTrafficManagerProfileMonitorConfig(profile.MonitorConfig))

		// fqdn is actually inside DNSConfig, inlined for simpler reference
		if dns := profile.DNSConfig; dns != nil {
			d.Set("fqdn", dns.Fqdn)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmTrafficManagerProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["trafficManagerProfiles"]

	update := trafficmanager.Profile{
		ProfileProperties: &trafficmanager.ProfileProperties{},
	}
	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("profile_status") {
		update.ProfileProperties.ProfileStatus = trafficmanager.ProfileStatus(d.Get("profile_status").(string))
	}

	if d.HasChange("traffic_routing_method") {
		update.ProfileProperties.TrafficRoutingMethod = trafficmanager.TrafficRoutingMethod(d.Get("traffic_routing_method").(string))
	}

	if d.HasChange("max_return") {
		if maxReturn, ok := d.GetOk("max_return"); ok {
			maxReturnInt := int64(maxReturn.(int))
			update.ProfileProperties.MaxReturn = &maxReturnInt
		} else {
			update.ProfileProperties.MaxReturn = nil
		}
	}

	if d.HasChange("dns_config") {
		update.ProfileProperties.DNSConfig = expandArmTrafficManagerDNSConfig(d)
	}

	if d.HasChange("monitor_config") {
		update.ProfileProperties.MonitorConfig = expandArmTrafficManagerMonitorConfig(d)
	}

	if _, err := client.Update(ctx, resourceGroup, name, update); err != nil {
		return fmt.Errorf("Error updating TrafficManager profile %q (resource group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmTrafficManagerProfileRead(d, meta)
}

func resourceArmTrafficManagerProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["trafficManagerProfiles"]

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return err
		}
	}

	return nil
}

func expandArmTrafficManagerMonitorConfig(d *schema.ResourceData) *trafficmanager.MonitorConfig {
	monitorSets := d.Get("monitor_config").([]interface{})
	monitor := monitorSets[0].(map[string]interface{})

	customHeaders := expandArmTrafficManagerCustomHeadersConfig(monitor["custom_header"].([]interface{}))

	cfg := trafficmanager.MonitorConfig{
		Protocol:                  trafficmanager.MonitorProtocol(monitor["protocol"].(string)),
		CustomHeaders:             customHeaders,
		Port:                      utils.Int64(int64(monitor["port"].(int))),
		Path:                      utils.String(monitor["path"].(string)),
		IntervalInSeconds:         utils.Int64(int64(monitor["interval_in_seconds"].(int))),
		TimeoutInSeconds:          utils.Int64(int64(monitor["timeout_in_seconds"].(int))),
		ToleratedNumberOfFailures: utils.Int64(int64(monitor["tolerated_number_of_failures"].(int))),
	}

	if v, ok := monitor["expected_status_code_ranges"].([]interface{}); ok {
		ranges := make([]trafficmanager.MonitorConfigExpectedStatusCodeRangesItem, 0)
		for _, r := range v {
			parts := strings.Split(r.(string), "-")
			min, _ := strconv.Atoi(parts[0])
			max, _ := strconv.Atoi(parts[1])
			ranges = append(ranges, trafficmanager.MonitorConfigExpectedStatusCodeRangesItem{
				Min: utils.Int32(int32(min)),
				Max: utils.Int32(int32(max)),
			})
		}
		cfg.ExpectedStatusCodeRanges = &ranges
	}

	return &cfg
}

func expandArmTrafficManagerCustomHeadersConfig(d []interface{}) *[]trafficmanager.MonitorConfigCustomHeadersItem {
	if len(d) == 0 || d[0] == nil {
		return nil
	}

	customHeaders := make([]trafficmanager.MonitorConfigCustomHeadersItem, len(d))

	for i, v := range d {
		ch := v.(map[string]interface{})
		customHeaders[i] = trafficmanager.MonitorConfigCustomHeadersItem{
			Name:  utils.String(ch["name"].(string)),
			Value: utils.String(ch["value"].(string)),
		}
	}

	return &customHeaders
}

func flattenArmTrafficManagerCustomHeadersConfig(input *[]trafficmanager.MonitorConfigCustomHeadersItem) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	headers := *input
	if len(headers) == 0 {
		return result
	}

	for _, v := range headers {
		header := make(map[string]string, 2)
		header["name"] = *v.Name
		header["value"] = *v.Value
		result = append(result, header)
	}

	return result
}

func expandArmTrafficManagerDNSConfig(d *schema.ResourceData) *trafficmanager.DNSConfig {
	dnsSets := d.Get("dns_config").([]interface{})
	dns := dnsSets[0].(map[string]interface{})

	name := dns["relative_name"].(string)
	ttl := int64(dns["ttl"].(int))

	return &trafficmanager.DNSConfig{
		RelativeName: &name,
		TTL:          &ttl,
	}
}

func flattenAzureRMTrafficManagerProfileDNSConfig(dns *trafficmanager.DNSConfig) []interface{} {
	result := make(map[string]interface{})

	result["relative_name"] = *dns.RelativeName
	result["ttl"] = int(*dns.TTL)

	return []interface{}{result}
}

func flattenAzureRMTrafficManagerProfileMonitorConfig(cfg *trafficmanager.MonitorConfig) []interface{} {
	result := make(map[string]interface{})

	result["protocol"] = string(cfg.Protocol)
	result["port"] = int(*cfg.Port)
	result["custom_header"] = flattenArmTrafficManagerCustomHeadersConfig(cfg.CustomHeaders)

	if cfg.Path != nil {
		result["path"] = *cfg.Path
	}

	result["interval_in_seconds"] = int(*cfg.IntervalInSeconds)
	result["timeout_in_seconds"] = int(*cfg.TimeoutInSeconds)
	result["tolerated_number_of_failures"] = int(*cfg.ToleratedNumberOfFailures)

	if v := cfg.ExpectedStatusCodeRanges; v != nil {
		ranges := make([]string, 0)
		for _, r := range *v {
			if r.Min == nil || r.Max == nil {
				continue
			}

			ranges = append(ranges, fmt.Sprintf("%d-%d", *r.Min, *r.Max))
		}
		result["expected_status_code_ranges"] = ranges
	}

	return []interface{}{result}
}

func validateTrafficManagerProfileStatusCodeRange(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	parts := strings.Split(v, "-")
	if len(parts) != 2 {
		errors = append(errors, fmt.Errorf("expected %s to contain a single '-', got %v", k, i))
		return warnings, errors
	}

	_, err := strconv.Atoi(parts[0])
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %s on the left of - to be an integer, got %v: %v", k, i, err))
		return warnings, errors
	}

	_, err = strconv.Atoi(parts[1])
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %s on the right of - to be an integer, got %v: %v", k, i, err))
		return warnings, errors
	}

	return warnings, errors
}
