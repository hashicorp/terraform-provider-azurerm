package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/frontdoor/mgmt/2019-04-01/frontdoor"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFrontDoor() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFrontDoorCreateUpdate,
		Read:   resourceArmFrontDoorRead,
		Update: resourceArmFrontDoorCreateUpdate,
		Delete: resourceArmFrontDoorDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group": azure.SchemaResourceGroupNameDiffSuppress(),

			"backend_pools": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backends": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"backend_host_header": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enabled_state": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.Enabled),
											string(frontdoor.Disabled),
										}, false),
										Default: string(frontdoor.Enabled),
									},
									"http_port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"https_port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"health_probe_settings": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"load_balancing_settings": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_state": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.PolicyResourceStateCreating),
								string(frontdoor.PolicyResourceStateEnabling),
								string(frontdoor.PolicyResourceStateEnabled),
								string(frontdoor.PolicyResourceStateDisabling),
								string(frontdoor.PolicyResourceStateDisabled),
								string(frontdoor.PolicyResourceStateDeleting),
							}, false),
							Default: string(frontdoor.PolicyResourceStateCreating),
						},
					},
				},
			},

			"backend_pools_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enforce_certificate_name_check": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.Enabled),
								string(frontdoor.Disabled),
							}, false),
							Default: string(frontdoor.Enabled),
						},
					},
				},
			},

			"enabled_state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(frontdoor.Enabled),
					string(frontdoor.Disabled),
				}, false),
				Default: string(frontdoor.Enabled),
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"frontend_endpoints": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_state": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.PolicyResourceStateCreating),
								string(frontdoor.PolicyResourceStateEnabling),
								string(frontdoor.PolicyResourceStateEnabled),
								string(frontdoor.PolicyResourceStateDisabling),
								string(frontdoor.PolicyResourceStateDisabled),
								string(frontdoor.PolicyResourceStateDeleting),
							}, false),
							Default: string(frontdoor.PolicyResourceStateCreating),
						},
						"session_affinity_enabled_state": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.Enabled),
								string(frontdoor.Disabled),
							}, false),
							Default: string(frontdoor.Enabled),
						},
						"session_affinity_ttl_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"web_application_firewall_policy_link": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"health_probe_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"interval_in_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.HTTP),
								string(frontdoor.HTTPS),
							}, false),
							Default: string(frontdoor.HTTP),
						},
						"resource_state": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.PolicyResourceStateCreating),
								string(frontdoor.PolicyResourceStateEnabling),
								string(frontdoor.PolicyResourceStateEnabled),
								string(frontdoor.PolicyResourceStateDisabling),
								string(frontdoor.PolicyResourceStateDisabled),
								string(frontdoor.PolicyResourceStateDeleting),
							}, false),
							Default: string(frontdoor.PolicyResourceStateCreating),
						},
					},
				},
			},

			"load_balancing_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_latency_milliseconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_state": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.PolicyResourceStateCreating),
								string(frontdoor.PolicyResourceStateEnabling),
								string(frontdoor.PolicyResourceStateEnabled),
								string(frontdoor.PolicyResourceStateDisabling),
								string(frontdoor.PolicyResourceStateDisabled),
								string(frontdoor.PolicyResourceStateDeleting),
							}, false),
							Default: string(frontdoor.PolicyResourceStateCreating),
						},
						"sample_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"successful_samples_required": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"resource_state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(frontdoor.PolicyResourceStateCreating),
					string(frontdoor.PolicyResourceStateEnabling),
					string(frontdoor.PolicyResourceStateEnabled),
					string(frontdoor.PolicyResourceStateDisabling),
					string(frontdoor.PolicyResourceStateDisabled),
					string(frontdoor.PolicyResourceStateDeleting),
				}, false),
				Default: string(frontdoor.PolicyResourceStateCreating),
			},

			"routing_rules": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accepted_protocols": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.HTTP),
								string(frontdoor.HTTPS),
							}, false),
							Default: string(frontdoor.HTTP),
						},
						"enabled_state": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.Enabled),
								string(frontdoor.Disabled),
							}, false),
							Default: string(frontdoor.Enabled),
						},
						"frontend_endpoints": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"patterns_to_match": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_state": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.PolicyResourceStateCreating),
								string(frontdoor.PolicyResourceStateEnabling),
								string(frontdoor.PolicyResourceStateEnabled),
								string(frontdoor.PolicyResourceStateDisabling),
								string(frontdoor.PolicyResourceStateDisabled),
								string(frontdoor.PolicyResourceStateDeleting),
							}, false),
							Default: string(frontdoor.PolicyResourceStateCreating),
						},
					},
				},
			},

			"tags": tagsSchema(),

			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"provisioning_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmFrontDoorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).frontDoorsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group").(string)

	if requireResourcesToBeImported {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_front_door", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	backendPools := d.Get("backend_pools").([]interface{})
	backendPoolsSettings := d.Get("backend_pools_settings").([]interface{})
	enabledState := d.Get("enabled_state").(string)
	friendlyName := d.Get("friendly_name").(string)
	frontendEndpoints := d.Get("frontend_endpoints").([]interface{})
	healthProbeSettings := d.Get("health_probe_settings").([]interface{})
	loadBalancingSettings := d.Get("load_balancing_settings").([]interface{})
	resourceState := d.Get("resource_state").(string)
	routingRules := d.Get("routing_rules").([]interface{})
	tags := d.Get("tags").(map[string]interface{})

	frontDoorParameters := frontdoor.FrontDoor{
		Location: utils.String(location),
		Properties: &frontdoor.Properties{
			BackendPools:          expandArmFrontDoorBackendPool(backendPools),
			BackendPoolsSettings:  expandArmFrontDoorBackendPoolsSettings(backendPoolsSettings),
			EnabledState:          frontdoor.EnabledState(enabledState),
			FriendlyName:          utils.String(friendlyName),
			FrontendEndpoints:     expandArmFrontDoorFrontendEndpoint(frontendEndpoints),
			HealthProbeSettings:   expandArmFrontDoorHealthProbeSettingsModel(healthProbeSettings),
			LoadBalancingSettings: expandArmFrontDoorLoadBalancingSettingsModel(loadBalancingSettings),
			ResourceState:         frontdoor.ResourceState(resourceState),
			RoutingRules:          expandArmFrontDoorRoutingRule(routingRules),
		},
		Tags: expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, frontDoorParameters)
	if err != nil {
		return fmt.Errorf("Error creating Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Front Door %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmFrontDoorRead(d, meta)
}

func resourceArmFrontDoorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).frontDoorsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["frontDoors"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Front Door %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if properties := resp.Properties; properties != nil {
		if err := d.Set("backend_pools", flattenArmFrontDoorBackendPool(properties.BackendPools)); err != nil {
			return fmt.Errorf("Error setting `backend_pools`: %+v", err)
		}
		if err := d.Set("backend_pools_settings", flattenArmFrontDoorBackendPoolsSettings(properties.BackendPoolsSettings)); err != nil {
			return fmt.Errorf("Error setting `backend_pools_settings`: %+v", err)
		}
		d.Set("cname", properties.Cname)
		d.Set("enabled_state", string(properties.EnabledState))
		d.Set("friendly_name", properties.FriendlyName)
		if err := d.Set("frontend_endpoints", flattenArmFrontDoorFrontendEndpoint(properties.FrontendEndpoints)); err != nil {
			return fmt.Errorf("Error setting `frontend_endpoints`: %+v", err)
		}
		if err := d.Set("health_probe_settings", flattenArmFrontDoorHealthProbeSettingsModel(properties.HealthProbeSettings)); err != nil {
			return fmt.Errorf("Error setting `health_probe_settings`: %+v", err)
		}
		if err := d.Set("load_balancing_settings", flattenArmFrontDoorLoadBalancingSettingsModel(properties.LoadBalancingSettings)); err != nil {
			return fmt.Errorf("Error setting `load_balancing_settings`: %+v", err)
		}
		d.Set("provisioning_state", properties.ProvisioningState)
		d.Set("resource_state", string(properties.ResourceState))
		if err := d.Set("routing_rules", flattenArmFrontDoorRoutingRule(properties.RoutingRules)); err != nil {
			return fmt.Errorf("Error setting `routing_rules`: %+v", err)
		}
	}
	d.Set("type", resp.Type)
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmFrontDoorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).frontDoorsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["frontDoors"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmFrontDoorBackendPool(input []interface{}) *[]frontdoor.BackendPool {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
  
	id := v["id"].(string)
	backends := v["backends"].([]interface{})
	loadBalancingSettings := v["load_balancing_settings"].([]interface{})
	healthProbeSettings := v["health_probe_settings"].([]interface{})
	resourceState := v["resource_state"].(string)
	name := v["name"].(string)

	result := frontdoor.BackendPool{
		ID:   utils.String(id),
		Name: utils.String(name),
		BackendPoolProperties: &frontdoor.BackendPoolProperties{
			Backends:              expandArmFrontDoorBackend(backends),
			HealthProbeSettings:   expandArmFrontDoorSubResource(healthProbeSettings),
			LoadBalancingSettings: expandArmFrontDoorSubResource(loadBalancingSettings),
			ResourceState:         frontdoor.ResourceState(resourceState),
		},
	}
	return &[]frontdoor.BackendPool{result}
}

func expandArmFrontDoorBackendPoolsSettings(input []interface{}) *frontdoor.BackendPoolsSettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	enforceCertificateNameCheck := v["enforce_certificate_name_check"].(string)

	result := frontdoor.BackendPoolsSettings{
		EnforceCertificateNameCheck: frontdoor.EnforceCertificateNameCheckEnabledState(enforceCertificateNameCheck),
	}
	return &result
}

func expandArmFrontDoorFrontendEndpoint(input []interface{}) *[]frontdoor.FrontendEndpoint {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	id := v["id"].(string)
	hostName := v["host_name"].(string)
	sessionAffinityEnabledState := v["session_affinity_enabled_state"].(string)
	sessionAffinityTtlSeconds := v["session_affinity_ttl_seconds"].(int32)
	webApplicationFirewallPolicyLink := v["web_application_firewall_policy_link"].([]interface{})
	resourceState := v["resource_state"].(string)
	name := v["name"].(string)

	result := frontdoor.FrontendEndpoint{
		ID:   utils.String(id),
		Name: utils.String(name),
		FrontendEndpointProperties: &frontdoor.FrontendEndpointProperties{
			HostName:                         utils.String(hostName),
			ResourceState:                    frontdoor.ResourceState(resourceState),
			SessionAffinityEnabledState:      frontdoor.SessionAffinityEnabledState(sessionAffinityEnabledState),
			SessionAffinityTTLSeconds:        utils.Int32(sessionAffinityTtlSeconds),
			WebApplicationFirewallPolicyLink: expandArmFrontDoorFrontendEndpointUpdateParameters_webApplicationFirewallPolicyLink(webApplicationFirewallPolicyLink),
		},
	}
	return &[]frontdoor.FrontendEndpoint{result}
}

func expandArmFrontDoorHealthProbeSettingsModel(input []interface{}) *[]frontdoor.HealthProbeSettingsModel {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	id := v["id"].(string)
	path := v["path"].(string)
	protocol := v["protocol"].(string)
	intervalInSeconds := v["interval_in_seconds"].(int32)
	resourceState := v["resource_state"].(string)
	name := v["name"].(string)

	result := frontdoor.HealthProbeSettingsModel{
		ID:   utils.String(id),
		Name: utils.String(name),
		HealthProbeSettingsProperties: &frontdoor.HealthProbeSettingsProperties{
			IntervalInSeconds: utils.Int32(intervalInSeconds),
			Path:              utils.String(path),
			Protocol:          frontdoor.Protocol(protocol),
			ResourceState:     frontdoor.ResourceState(resourceState),
		},
	}
	return &[]frontdoor.HealthProbeSettingsModel{result}
}

func expandArmFrontDoorLoadBalancingSettingsModel(input []interface{}) *[]frontdoor.LoadBalancingSettingsModel {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	id := v["id"].(string)
	sampleSize := v["sample_size"].(int32)
	successfulSamplesRequired := v["successful_samples_required"].(int32)
	additionalLatencyMilliseconds := v["additional_latency_milliseconds"].(int32)
	resourceState := v["resource_state"].(string)
	name := v["name"].(string)

	result := frontdoor.LoadBalancingSettingsModel{
		ID:   utils.String(id),
		Name: utils.String(name),
		LoadBalancingSettingsProperties: &frontdoor.LoadBalancingSettingsProperties{
			AdditionalLatencyMilliseconds: utils.Int32(additionalLatencyMilliseconds),
			ResourceState:                 frontdoor.ResourceState(resourceState),
			SampleSize:                    utils.Int32(sampleSize),
			SuccessfulSamplesRequired:     utils.Int32(successfulSamplesRequired),
		},
	}
	return &[]frontdoor.LoadBalancingSettingsModel{result}
}

func expandArmFrontDoorRoutingRule(input []interface{}) *[]frontdoor.RoutingRule {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	id := v["id"].(string)
	frontendEndpoints := v["frontend_endpoints"].([]interface{})
	acceptedProtocols := v["accepted_protocols"].([]interface{})
	patternsToMatch := v["patterns_to_match"].([]interface{})
	enabledState := v["enabled_state"].(string)
	resourceState := v["resource_state"].(string)
	name := v["name"].(string)

	result := frontdoor.RoutingRule{
		ID:   utils.String(id),
		Name: utils.String(name),
		RoutingRuleProperties: &frontdoor.RoutingRuleProperties{
			AcceptedProtocols: expandArmFrontDoorProtocols(acceptedProtocols),
			EnabledState:      frontdoor.RoutingRuleEnabledState(enabledState),
			FrontendEndpoints: expandArmFrontDoorSubResource(frontendEndpoints),
			PatternsToMatch:   *patternsToMatch,
			ResourceState:     frontdoor.ResourceState(resourceState),
		},
	}
	return &[]frontdoor.RoutingRule{result}
}
func expandArmFrontDoorProtocols(input []interface{}) *[]frontdoor.Protocol {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(string)

	protocol := v[0]
	result := frontdoor.Protocol{protocol}

	return &[]frontdoor.Protocol{result}
}

func expandArmFrontDoorBackend(input []interface{}) *[]frontdoor.Backend {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	address := v["address"].(string)
	httpPort := v["http_port"].(int32)
	httpsPort := v["https_port"].(int32)
	enabledState := v["enabled_state"].(string)
	priority := v["priority"].(int32)
	weight := v["weight"].(int32)
	backendHostHeader := v["backend_host_header"].(string)

	result := frontdoor.Backend{
		Address:           utils.String(address),
		BackendHostHeader: utils.String(backendHostHeader),
		EnabledState:      frontdoor.BackendEnabledState(enabledState),
		HTTPPort:          utils.Int32(httpPort),
		HTTPSPort:         utils.Int32(httpsPort),
		Priority:          utils.Int32(priority),
		Weight:            utils.Int32(weight),
	}
	return &[]frontdoor.Backend{result}
}

func expandArmFrontDoorSubResource(input []interface{}) *frontdoor.SubResource {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	id := v["id"].(string)

	result := frontdoor.SubResource{
		ID: utils.String(id),
	}
	return &frontdoor.SubResource{result}
}

func expandArmFrontDoorFrontendEndpointUpdateParameters_webApplicationFirewallPolicyLink(input []interface{}) *frontdoor.FrontendEndpointUpdateParametersWebApplicationFirewallPolicyLink {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	id := v["id"].(string)

	result := frontdoor.FrontendEndpointUpdateParametersWebApplicationFirewallPolicyLink{
		ID: utils.String(id),
	}
	return &result
}

func flattenArmFrontDoorBackendPool(input *[]frontdoor.BackendPool) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if id := input.ID; id != nil {
		result["id"] = *id
	}
	if properties := input.BackendPoolProperties; properties != nil {
		result["backends"] = flattenArmFrontDoorBackend(properties.Backends)
		result["health_probe_settings"] = flattenArmFrontDoorSubResource(properties.HealthProbeSettings)
		result["load_balancing_settings"] = flattenArmFrontDoorSubResource(properties.LoadBalancingSettings)
		if resourceState := string(properties.ResourceState); resourceState != nil {
			result["resource_state"] = *resourceState
		}
	}

	return []interface{}{result}
}

func flattenArmFrontDoorBackendPoolsSettings(input *frontdoor.BackendPoolsSettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if enforceCertificateNameCheck := string(input.EnforceCertificateNameCheck); enforceCertificateNameCheck != "" {
		result["enforce_certificate_name_check"] = enforceCertificateNameCheck
	}

	return []interface{}{result}
}

func flattenArmFrontDoorFrontendEndpoint(input *[]frontdoor.FrontendEndpoint) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if id := input.ID; id != nil {
		result["id"] = *id
	}
	if properties := input.FrontendEndpointProperties; properties != nil {
		if hostName := properties.HostName; hostName != nil {
			result["host_name"] = *hostName
		}
		if resourceState := string(properties.ResourceState); resourceState != nil {
			result["resource_state"] = *resourceState
		}
		if sessionAffinityEnabledState := string(properties.SessionAffinityEnabledState); sessionAffinityEnabledState != nil {
			result["session_affinity_enabled_state"] = *sessionAffinityEnabledState
		}
		if sessionAffinityTtlSeconds := properties.SessionAffinityTtlSeconds; sessionAffinityTtlSeconds != nil {
			result["session_affinity_ttl_seconds"] = *sessionAffinityTtlSeconds
		}
		result["web_application_firewall_policy_link"] = flattenArmFrontDoorFrontendEndpointUpdateParameters_webApplicationFirewallPolicyLink(properties.WebApplicationFirewallPolicyLink)
	}

	return []interface{}{result}
}

func flattenArmFrontDoorHealthProbeSettingsModel(input *[]frontdoor.HealthProbeSettingsModel) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if id := input.ID; id != nil {
		result["id"] = *id
	}
	if properties := input.HealthProbeSettingsProperties; properties != nil {
		if intervalInSeconds := properties.IntervalInSeconds; intervalInSeconds != nil {
			result["interval_in_seconds"] = *intervalInSeconds
		}
		if path := properties.Path; path != nil {
			result["path"] = *path
		}
		if protocol := string(properties.Protocol); protocol != nil {
			result["protocol"] = *protocol
		}
		if resourceState := string(properties.ResourceState); resourceState != nil {
			result["resource_state"] = *resourceState
		}
	}

	return []interface{}{result}
}

func flattenArmFrontDoorLoadBalancingSettingsModel(input *[]frontdoor.LoadBalancingSettingsModel) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if id := input.ID; id != nil {
		result["id"] = *id
	}
	if properties := input.Properties; properties != nil {
		if additionalLatencyMilliseconds := properties.AdditionalLatencyMilliseconds; additionalLatencyMilliseconds != nil {
			result["additional_latency_milliseconds"] = *additionalLatencyMilliseconds
		}
		if resourceState := string(properties.ResourceState); resourceState != nil {
			result["resource_state"] = *resourceState
		}
		if sampleSize := properties.SampleSize; sampleSize != nil {
			result["sample_size"] = *sampleSize
		}
		if successfulSamplesRequired := properties.SuccessfulSamplesRequired; successfulSamplesRequired != nil {
			result["successful_samples_required"] = *successfulSamplesRequired
		}
	}

	return []interface{}{result}
}

func flattenArmFrontDoorRoutingRule(input *[]frontdoor.RoutingRule) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if id := input.ID; id != nil {
		result["id"] = *id
	}
	if properties := input.Properties; properties != nil {
		if acceptedProtocols := string(properties.AcceptedProtocols); acceptedProtocols != nil {
			result["accepted_protocols"] = *acceptedProtocols
		}
		if enabledState := string(properties.EnabledState); enabledState != nil {
			result["enabled_state"] = *enabledState
		}
		result["frontend_endpoints"] = flattenArmFrontDoorSubResource(properties.FrontendEndpoints)
		if patternsToMatch := properties.PatternsToMatch; patternsToMatch != nil {
			result["patterns_to_match"] = *patternsToMatch
		}
		if resourceState := string(properties.ResourceState); resourceState != nil {
			result["resource_state"] = *resourceState
		}
	}

	return []interface{}{result}
}

func flattenArmFrontDoorBackend(input *frontdoor.Backend) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if address := input.Address; address != nil {
		result["address"] = *address
	}
	if backendHostHeader := input.BackendHostHeader; backendHostHeader != nil {
		result["backend_host_header"] = *backendHostHeader
	}
	if enabledState := string(input.EnabledState); enabledState != nil {
		result["enabled_state"] = *enabledState
	}
	if httpPort := input.HttpPort; httpPort != nil {
		result["http_port"] = *httpPort
	}
	if httpsPort := input.HttpsPort; httpsPort != nil {
		result["https_port"] = *httpsPort
	}
	if priority := input.Priority; priority != nil {
		result["priority"] = *priority
	}
	if weight := input.Weight; weight != nil {
		result["weight"] = *weight
	}

	return []interface{}{result}
}

func flattenArmFrontDoorSubResource(input *frontdoor.SubResource) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if id := input.ID; id != nil {
		result["id"] = *id
	}

	return []interface{}{result}
}

func flattenArmFrontDoorFrontendEndpointUpdateParameters_webApplicationFirewallPolicyLink(input *frontdoor.FrontendEndpointUpdateParametersWebApplicationFirewallPolicyLink) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if id := input.ID; id != nil {
		result["id"] = *id
	}

	return []interface{}{result}
}
