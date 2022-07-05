package trafficmanager

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmTrafficManagerProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmTrafficManagerProfileCreate,
		Read:   resourceArmTrafficManagerProfileRead,
		Update: resourceArmTrafficManagerProfileUpdate,
		Delete: resourceArmTrafficManagerProfileDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := profiles.ParseTrafficManagerProfileID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"profile_status": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(profiles.ProfileStatusEnabled),
					string(profiles.ProfileStatusDisabled),
				}, false),
			},

			"traffic_routing_method": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(profiles.TrafficRoutingMethodGeographic),
					string(profiles.TrafficRoutingMethodWeighted),
					string(profiles.TrafficRoutingMethodPerformance),
					string(profiles.TrafficRoutingMethodPriority),
					string(profiles.TrafficRoutingMethodSubnet),
					string(profiles.TrafficRoutingMethodMultiValue),
				}, false),
			},

			"dns_config": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"relative_name": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Required: true,
						},
						"ttl": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 2147483647),
						},
					},
				},
			},

			"monitor_config": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"expected_status_code_ranges": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.StatusCodeRange,
							},
						},

						"custom_header": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"value": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
								},
							},
						},

						"protocol": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(profiles.MonitorProtocolHTTP),
								string(profiles.MonitorProtocolHTTPS),
								string(profiles.MonitorProtocolTCP),
							}, false),
						},

						"port": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"path": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"interval_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{10, 30}),
							Default:      30,
						},

						"timeout_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(5, 10),
							Default:      10,
						},

						"tolerated_number_of_failures": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 9),
							Default:      3,
						},
					},
				},
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"max_return": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 8),
			},

			"traffic_view_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceArmTrafficManagerProfileCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Traffic Manager Profile creation.")

	id := profiles.NewTrafficManagerProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s", id)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_traffic_manager_profile", id.ID())
	}

	trafficRoutingMethod := profiles.TrafficRoutingMethod(d.Get("traffic_routing_method").(string))
	// No existing profile - start from a new struct.
	profile := profiles.Profile{
		Name:     utils.String(id.ProfileName),
		Location: utils.String("global"), // must be provided in request
		Properties: &profiles.ProfileProperties{
			TrafficRoutingMethod: &trafficRoutingMethod,
			DnsConfig:            expandArmTrafficManagerDNSConfig(d),
			MonitorConfig:        expandArmTrafficManagerMonitorConfig(d),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if maxReturn, ok := d.GetOk("max_return"); ok {
		profile.Properties.MaxReturn = utils.Int64(int64(maxReturn.(int)))
	}

	if status, ok := d.GetOk("profile_status"); ok {
		profileStatus := profiles.ProfileStatus(status.(string))
		profile.Properties.ProfileStatus = &profileStatus
	}

	if trafficViewStatus, ok := d.GetOk("traffic_view_enabled"); ok {
		profile.Properties.TrafficViewEnrollmentStatus = expandArmTrafficManagerTrafficView(trafficViewStatus.(bool))
	}

	trafficRoutingMethodPtr := profile.Properties.TrafficRoutingMethod
	if *trafficRoutingMethodPtr == profiles.TrafficRoutingMethodMultiValue &&
		profile.Properties.MaxReturn == nil {
		return fmt.Errorf("`max_return` must be specified when `traffic_routing_method` is set to `MultiValue`")
	}

	if *profile.Properties.MonitorConfig.IntervalInSeconds == int64(10) &&
		*profile.Properties.MonitorConfig.TimeoutInSeconds == int64(10) {
		return fmt.Errorf("`timeout_in_seconds` must be between `5` and `9` when `interval_in_seconds` is set to `10`")
	}

	if _, err := client.CreateOrUpdate(ctx, id, profile); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArmTrafficManagerProfileRead(d, meta)
}

func resourceArmTrafficManagerProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := profiles.ParseTrafficManagerProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ProfileName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if profile := model.Properties; profile != nil {
			profileStatus := ""
			if profile.ProfileStatus != nil {
				profileStatus = string(*profile.ProfileStatus)
			}
			d.Set("profile_status", profileStatus)
			trafficRoutingMethod := ""
			if profile.TrafficRoutingMethod != nil {
				trafficRoutingMethod = string(*profile.TrafficRoutingMethod)
			}
			d.Set("traffic_routing_method", trafficRoutingMethod)
			d.Set("max_return", profile.MaxReturn)

			d.Set("dns_config", flattenAzureRMTrafficManagerProfileDNSConfig(profile.DnsConfig))
			d.Set("monitor_config", flattenAzureRMTrafficManagerProfileMonitorConfig(profile.MonitorConfig))
			trafficViewEnabled := false
			if profile.TrafficViewEnrollmentStatus != nil {
				trafficViewEnabled = *profile.TrafficViewEnrollmentStatus == profiles.TrafficViewEnrollmentStatusEnabled
			}
			d.Set("traffic_view_enabled", trafficViewEnabled)

			// fqdn is actually inside DNSConfig, inlined for simpler reference
			if dns := profile.DnsConfig; dns != nil {
				d.Set("fqdn", dns.Fqdn)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceArmTrafficManagerProfileUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := profiles.ParseTrafficManagerProfileID(d.Id())
	if err != nil {
		return err
	}

	update := profiles.Profile{
		Properties: &profiles.ProfileProperties{},
	}
	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("profile_status") {
		status := profiles.ProfileStatus(d.Get("profile_status").(string))
		update.Properties.ProfileStatus = &status
	}

	if d.HasChange("traffic_routing_method") {
		routingMethod := profiles.TrafficRoutingMethod(d.Get("traffic_routing_method").(string))
		update.Properties.TrafficRoutingMethod = &routingMethod
	}

	if d.HasChange("max_return") {
		if maxReturn, ok := d.GetOk("max_return"); ok {
			update.Properties.MaxReturn = utils.Int64(int64(maxReturn.(int)))
		}
	}

	if d.HasChange("dns_config") {
		update.Properties.DnsConfig = expandArmTrafficManagerDNSConfig(d)
	}

	if d.HasChange("monitor_config") {
		update.Properties.MonitorConfig = expandArmTrafficManagerMonitorConfig(d)
	}

	if d.HasChange("traffic_view_enabled") {
		if trafficViewStatus, ok := d.GetOk("traffic_view_enabled"); ok {
			update.Properties.TrafficViewEnrollmentStatus = expandArmTrafficManagerTrafficView(trafficViewStatus.(bool))
		}
	}

	if _, err := client.Update(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceArmTrafficManagerProfileRead(d, meta)
}

func resourceArmTrafficManagerProfileDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := profiles.ParseTrafficManagerProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return err
		}
	}

	return nil
}

func expandArmTrafficManagerMonitorConfig(d *pluginsdk.ResourceData) *profiles.MonitorConfig {
	monitorSets := d.Get("monitor_config").([]interface{})
	monitor := monitorSets[0].(map[string]interface{})

	customHeaders := expandArmTrafficManagerCustomHeadersConfig(monitor["custom_header"].([]interface{}))

	protocol := profiles.MonitorProtocol(monitor["protocol"].(string))
	cfg := profiles.MonitorConfig{
		Protocol:                  &protocol,
		CustomHeaders:             customHeaders,
		Port:                      utils.Int64(int64(monitor["port"].(int))),
		Path:                      utils.String(monitor["path"].(string)),
		IntervalInSeconds:         utils.Int64(int64(monitor["interval_in_seconds"].(int))),
		TimeoutInSeconds:          utils.Int64(int64(monitor["timeout_in_seconds"].(int))),
		ToleratedNumberOfFailures: utils.Int64(int64(monitor["tolerated_number_of_failures"].(int))),
	}

	if v, ok := monitor["expected_status_code_ranges"].([]interface{}); ok {
		ranges := make([]profiles.MonitorConfigExpectedStatusCodeRangesInlined, 0)
		for _, r := range v {
			parts := strings.Split(r.(string), "-")
			min, _ := strconv.Atoi(parts[0])
			max, _ := strconv.Atoi(parts[1])
			ranges = append(ranges, profiles.MonitorConfigExpectedStatusCodeRangesInlined{
				Min: utils.Int64(int64(min)),
				Max: utils.Int64(int64(max)),
			})
		}
		cfg.ExpectedStatusCodeRanges = &ranges
	}

	return &cfg
}

func expandArmTrafficManagerCustomHeadersConfig(d []interface{}) *[]profiles.MonitorConfigCustomHeadersInlined {
	if len(d) == 0 || d[0] == nil {
		return nil
	}

	customHeaders := make([]profiles.MonitorConfigCustomHeadersInlined, len(d))

	for i, v := range d {
		ch := v.(map[string]interface{})
		customHeaders[i] = profiles.MonitorConfigCustomHeadersInlined{
			Name:  utils.String(ch["name"].(string)),
			Value: utils.String(ch["value"].(string)),
		}
	}

	return &customHeaders
}

func flattenArmTrafficManagerCustomHeadersConfig(input *[]profiles.MonitorConfigCustomHeadersInlined) []interface{} {
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

func expandArmTrafficManagerDNSConfig(d *pluginsdk.ResourceData) *profiles.DnsConfig {
	dnsSets := d.Get("dns_config").([]interface{})
	dns := dnsSets[0].(map[string]interface{})

	name := dns["relative_name"].(string)
	ttl := int64(dns["ttl"].(int))

	return &profiles.DnsConfig{
		RelativeName: &name,
		Ttl:          &ttl,
	}
}

func expandArmTrafficManagerTrafficView(s bool) *profiles.TrafficViewEnrollmentStatus {
	enabled := profiles.TrafficViewEnrollmentStatusDisabled
	if s {
		enabled = profiles.TrafficViewEnrollmentStatusEnabled
	}
	return &enabled
}

func flattenAzureRMTrafficManagerProfileDNSConfig(dns *profiles.DnsConfig) []interface{} {
	result := make(map[string]interface{})

	result["relative_name"] = *dns.RelativeName
	result["ttl"] = int(*dns.Ttl)

	return []interface{}{result}
}

func flattenAzureRMTrafficManagerProfileMonitorConfig(cfg *profiles.MonitorConfig) []interface{} {
	result := make(map[string]interface{})

	protocol := ""
	if cfg.Protocol != nil {
		protocol = string(*cfg.Protocol)
	}
	result["protocol"] = protocol
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
