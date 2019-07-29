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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateFrontDoorName,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"routing_rule": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateBackendPoolRoutingRuleName,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"accepted_protocols": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 2,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(frontdoor.HTTP),
									string(frontdoor.HTTPS),
								}, false),
								Default: string(frontdoor.HTTP),
							},
						},
						"patterns_to_match": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 25,
							Elem: &schema.Schema{
								Type:    schema.TypeString,
								Default: "/*",
							},
						},
						"frontend_endpoints": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 100,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"redirect_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"custom_fragment": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"custom_host": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"custom_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"custom_query_string": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"redirect_protocol": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.RedirectProtocolHTTPOnly),
											string(frontdoor.RedirectProtocolHTTPSOnly),
											string(frontdoor.RedirectProtocolMatchRequest),
										}, false),
									},
									"redirect_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.Found),
											string(frontdoor.Moved),
											string(frontdoor.PermanentRedirect),
											string(frontdoor.TemporaryRedirect),
										}, false),
									},
								},
							},
						},
						"forwarding_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"backend_pool_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"cache_use_dynamic_compression": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"cache_query_parameter_strip_directive": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.StripAll),
											string(frontdoor.StripNone),
										}, false),
										Default: string(frontdoor.StripNone),
									},
									"custom_forwarding_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"forwarding_protocol": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.HTTPOnly),
											string(frontdoor.HTTPSOnly),
											string(frontdoor.MatchRequest),
										}, false),
										Default: string(frontdoor.MatchRequest),
									},
								},
							},
						},
					},
				},
			},
			"enforce_backend_pools_certificate_name_check": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"backend_pool_load_balancing": {
				Type:     schema.TypeList,
				MaxItems: 5000,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateBackendPoolRoutingRuleName,
						},
						"sample_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  4,
						},
						"successful_samples_required": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  2,
						},
						"additional_latency_milliseconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  2,
						},
					},
				},
			},
			"backend_pool_health_probe": {
				Type:     schema.TypeList,
				MaxItems: 5000,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateBackendPoolRoutingRuleName,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "/",
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
						"interval_in_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  120,
						},
					},
				},
			},
			"backend_pool": {
				Type:     schema.TypeList,
				MaxItems: 50,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend": {
							Type:     schema.TypeList,
							MaxItems: 100,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"http_port": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"https_port": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"weight": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 1000),
										Default:      50,
									},
									"priority": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 5),
									},
									"host_header": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateBackendPoolRoutingRuleName,
						},
						"health_probe_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"load_balancing_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"frontend_endpoint": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateBackendPoolRoutingRuleName,
						},
						"host_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"session_affinity_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"session_affinity_ttl_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"custom_https_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_source": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.CertificateSourceAzureKeyVault),
											string(frontdoor.CertificateSourceFrontDoor),
										}, false),
										Default: string(frontdoor.CertificateSourceFrontDoor),
									},
									// NOTE: None of these attributes are valid if
									//       certificate_source is set to FrontDoor
									"azure_key_vault_certificate_secret_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"azure_key_vault_certificate_secret_version": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"azure_key_vault_certificate_vault_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmFrontDoorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if err := azure.ValidateFrontdoor(d); err != nil {
		return fmt.Errorf("Error creating Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	client := meta.(*ArmClient).frontDoorsClient
	ctx := meta.(*ArmClient).StopContext
	subscriptionId := meta.(*ArmClient).subscriptionId

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

	// if subId := azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, name, resourceType, resourceName); subId == "" {
	// 	return fmt.Errorf("Error Front Door %q (Resource Group %q): unable to creating ID for sub resource", name, resourceGroup)
	// }

	location := azure.NormalizeLocation(d.Get("location").(string))
	friendlyName := d.Get("friendly_name").(string)
	routingRules := d.Get("routing_rule").([]interface{})
	loadBalancingSettings := d.Get("backend_pool_load_balancing").([]interface{})
	healthProbeSettings := d.Get("backend_pool_health_probe").([]interface{})
	backendPools := d.Get("backend_pool").([]interface{})
	frontendEndpoints := d.Get("frontend_endpoint").([]interface{})
	backendPoolsSettings := d.Get("enforce_backend_pools_certificate_name_check").(bool)
	enabledState := d.Get("enabled").(bool)
	tags := d.Get("tags").(map[string]interface{})

	frontDoorParameters := frontdoor.FrontDoor{
		Location: utils.String(location),
		Properties: &frontdoor.Properties{
			FriendlyName:          utils.String(friendlyName),
			RoutingRules:          expandArmFrontDoorRoutingRule(routingRules, subscriptionId, resourceGroup, name),
			BackendPools:          expandArmFrontDoorBackendPools(backendPools, subscriptionId, resourceGroup, name),
			BackendPoolsSettings:  expandArmFrontDoorBackendPoolsSettings(backendPoolsSettings),
			FrontendEndpoints:     expandArmFrontDoorFrontendEndpoint(frontendEndpoints, subscriptionId, resourceGroup, name),
			HealthProbeSettings:   expandArmFrontDoorHealthProbeSettingsModel(healthProbeSettings, subscriptionId, resourceGroup, name),
			LoadBalancingSettings: expandArmFrontDoorLoadBalancingSettingsModel(loadBalancingSettings, subscriptionId, resourceGroup, name),
			EnabledState:          expandArmFrontDoorEnabledState(enabledState),
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
	name := id.Path["frontdoors"]

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
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if properties := resp.Properties; properties != nil {
		if err := d.Set("backend_pool", flattenArmFrontDoorBackendPools(properties.BackendPools)); err != nil {
			return fmt.Errorf("Error setting `backend_pool`: %+v", err)
		}
		if err := d.Set("enforce_backend_pools_certificate_name_check", flattenArmFrontDoorBackendPoolsSettings(properties.BackendPoolsSettings)); err != nil {
			return fmt.Errorf("Error setting `enforce_backend_pools_certificate_name_check`: %+v", err)
		}
		d.Set("cname", properties.Cname)
		if properties.EnabledState == frontdoor.EnabledStateEnabled {
			d.Set("enabled", true)
		} else {
			d.Set("enabled", false)
		}
		d.Set("friendly_name", properties.FriendlyName)
		if err := d.Set("frontend_endpoint", flattenArmFrontDoorFrontendEndpoint(properties.FrontendEndpoints)); err != nil {
			return fmt.Errorf("Error setting `frontend_endpoint`: %+v", err)
		}
		if err := d.Set("backend_pool_health_probe", flattenArmFrontDoorHealthProbeSettingsModel(properties.HealthProbeSettings)); err != nil {
			return fmt.Errorf("Error setting `backend_pool_health_probe`: %+v", err)
		}
		if err := d.Set("backend_pool_load_balancing", flattenArmFrontDoorLoadBalancingSettingsModel(properties.LoadBalancingSettings)); err != nil {
			return fmt.Errorf("Error setting `backend_pool_load_balancing`: %+v", err)
		}
		if err := d.Set("routing_rule", flattenArmFrontDoorRoutingRule(properties.RoutingRules)); err != nil {
			return fmt.Errorf("Error setting `routing_rules`: %+v", err)
		}
	}

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
	name := id.Path["frontdoors"]

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

func expandArmFrontDoorBackendPools(input []interface{}, subscriptionId string, resourceGroup string, serviceName string) *[]frontdoor.BackendPool {
	if len(input) == 0 {
		return &[]frontdoor.BackendPool{}
	}

	output := make([]frontdoor.BackendPool, 0)

	for _, bp := range input {
		backendPool := bp.(map[string]interface{})

		backendPoolName := backendPool["name"].(string)
		backendPoolLoadBalancingName := backendPool["load_balancing_name"].(string)
		backendPoolHealthProbeName := backendPool["health_probe_name"].(string)

		backends := backendPool["backend"].([]interface{})

		result := frontdoor.BackendPool{
			ID:   utils.String(azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, serviceName, "BackendPools", backendPoolName)),
			Name: utils.String(backendPoolName),
			BackendPoolProperties: &frontdoor.BackendPoolProperties{
				// ResourceState
				Backends:              expandArmFrontDoorBackend(backends),
				LoadBalancingSettings: expandArmFrontDoorSubResource(subscriptionId, resourceGroup, serviceName, "LoadBalancingSettings", backendPoolLoadBalancingName),
				HealthProbeSettings:   expandArmFrontDoorSubResource(subscriptionId, resourceGroup, serviceName, "HealthProbeSettings", backendPoolHealthProbeName),
			},
		}

		output = append(output, result)
	}

	return &output
}

func expandArmFrontDoorBackend(input []interface{}) *[]frontdoor.Backend {
	if len(input) == 0 {
		return &[]frontdoor.Backend{}
	}

	output := make([]frontdoor.Backend, 0)

	for _, be := range input {
		backend := be.(map[string]interface{})

		address := backend["address"].(string)
		hostHeader := backend["host_header"].(string)
		enabled := backend["enabled"].(bool)
		httpPort := int32(backend["http_port"].(int))
		httpsPort := int32(backend["https_port"].(int))
		priority := int32(backend["priority"].(int))
		weight := int32(backend["weight"].(int))

		result := frontdoor.Backend{
			Address:           utils.String(address),
			BackendHostHeader: utils.String(hostHeader),
			EnabledState:      expandArmFrontDoorBackendEnabledState(enabled),
			HTTPPort:          utils.Int32(httpPort),
			HTTPSPort:         utils.Int32(httpsPort),
			Priority:          utils.Int32(priority),
			Weight:            utils.Int32(weight),
		}

		output = append(output, result)
	}

	return &output
}

func expandArmFrontDoorBackendEnabledState(isEnabled bool) frontdoor.BackendEnabledState {
	if isEnabled {
		return frontdoor.Enabled
	}

	return frontdoor.Disabled
}

func expandArmFrontDoorBackendPoolsSettings(enforceCertificateNameCheck bool) *frontdoor.BackendPoolsSettings {
	enforceCheck := frontdoor.EnforceCertificateNameCheckEnabledStateDisabled

	if enforceCertificateNameCheck {
		enforceCheck = frontdoor.EnforceCertificateNameCheckEnabledStateEnabled
	}

	result := frontdoor.BackendPoolsSettings{
		EnforceCertificateNameCheck: enforceCheck,
	}

	return &result
}

func expandArmFrontDoorFrontendEndpoint(input []interface{}, subscriptionId string, resourceGroup string, serviceName string) *[]frontdoor.FrontendEndpoint {
	if len(input) == 0 {
		return &[]frontdoor.FrontendEndpoint{}
	}

	output := make([]frontdoor.FrontendEndpoint, 0)

	for _, frontendEndpoints := range input {
		frontendEndpoint := frontendEndpoints.(map[string]interface{})

		hostName := frontendEndpoint["host_name"].(string)
		isSessionAffinityEnabled := frontendEndpoint["session_affinity_enabled"].(bool)
		sessionAffinityTtlSeconds := int32(frontendEndpoint["session_affinity_ttl_seconds"].(int))
		customHttpsConfiguration := frontendEndpoint["custom_https_configuration"].([]interface{})
		name := frontendEndpoint["name"].(string)
		id := utils.String(azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, serviceName, "FrontendEndpoints", name))

		sessionAffinityEnabled := frontdoor.SessionAffinityEnabledStateDisabled

		if isSessionAffinityEnabled {
			sessionAffinityEnabled = frontdoor.SessionAffinityEnabledStateEnabled
		}

		result := frontdoor.FrontendEndpoint{
			ID:   id,
			Name: utils.String(name),
			FrontendEndpointProperties: &frontdoor.FrontendEndpointProperties{
				// ResourceState:
				// CustomHTTPSProvisioningState:
				// CustomHTTPSProvisioningSubstate:
				CustomHTTPSConfiguration:    expandArmFrontDoorCustomHTTPSConfiguration(customHttpsConfiguration),
				HostName:                    utils.String(hostName),
				SessionAffinityEnabledState: sessionAffinityEnabled,
				SessionAffinityTTLSeconds:   utils.Int32(sessionAffinityTtlSeconds),
				// WebApplicationFirewallPolicyLink:

			},
		}

		output = append(output, result)
	}

	return &output
}

func expandArmFrontDoorCustomHTTPSConfiguration(input []interface{}) *frontdoor.CustomHTTPSConfiguration {
	if len(input) == 0 {
		return &frontdoor.CustomHTTPSConfiguration{}
	}

	v := input[0].(map[string]interface{})
	certSource := v["certificate_source"].(string)

	result := frontdoor.CustomHTTPSConfiguration{
		ProtocolType: frontdoor.ServerNameIndication,
	}

	if certSource == "AzureKeyVault" {
		vaultSecret := v["azure_key_vault_certificate_secret_name"].(string)
		vaultVersion := v["azure_key_vault_certificate_secret_version"].(string)
		vaultId := v["azure_key_vault_certificate_vault_id"].(string)

		result.CertificateSource = frontdoor.CertificateSourceAzureKeyVault
		result.KeyVaultCertificateSourceParameters = &frontdoor.KeyVaultCertificateSourceParameters{
			Vault: &frontdoor.KeyVaultCertificateSourceParametersVault{
				ID: utils.String(vaultId),
			},
			SecretName:    utils.String(vaultSecret),
			SecretVersion: utils.String(vaultVersion),
		}
	} else {
		result.CertificateSource = frontdoor.CertificateSourceFrontDoor
		result.CertificateSourceParameters = &frontdoor.CertificateSourceParameters{
			CertificateType: frontdoor.Dedicated,
		}
	}

	return &result
}

func expandArmFrontDoorHealthProbeSettingsModel(input []interface{}, subscriptionId string, resourceGroup string, serviceName string) *[]frontdoor.HealthProbeSettingsModel {
	if len(input) == 0 {
		return &[]frontdoor.HealthProbeSettingsModel{}
	}

	output := make([]frontdoor.HealthProbeSettingsModel, 0)

	for _, hps := range input {
		v := hps.(map[string]interface{})

		path := v["path"].(string)
		protocol := v["protocol"].(string)
		intervalInSeconds := int32(v["interval_in_seconds"].(int))
		name := v["name"].(string)

		result := frontdoor.HealthProbeSettingsModel{
			ID:   utils.String(azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, serviceName, "HealthProbeSettings", name)),
			Name: utils.String(name),
			HealthProbeSettingsProperties: &frontdoor.HealthProbeSettingsProperties{
				IntervalInSeconds: utils.Int32(intervalInSeconds),
				Path:              utils.String(path),
				Protocol:          frontdoor.Protocol(protocol),
			},
		}

		output = append(output, result)
	}

	return &output
}

func expandArmFrontDoorLoadBalancingSettingsModel(input []interface{}, subscriptionId string, resourceGroup string, serviceName string) *[]frontdoor.LoadBalancingSettingsModel {
	if len(input) == 0 {
		return &[]frontdoor.LoadBalancingSettingsModel{}
	}

	output := make([]frontdoor.LoadBalancingSettingsModel, 0)

	for _, lbs := range input {
		loadBalanceSetting := lbs.(map[string]interface{})

		name := loadBalanceSetting["name"].(string)
		sampleSize := int32(loadBalanceSetting["sample_size"].(int))
		successfulSamplesRequired := int32(loadBalanceSetting["successful_samples_required"].(int))
		additionalLatencyMilliseconds := int32(loadBalanceSetting["additional_latency_milliseconds"].(int))
		id := utils.String(azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, serviceName, "LoadBalancingSettings", name))

		result := frontdoor.LoadBalancingSettingsModel{
			ID:   id,
			Name: utils.String(name),
			LoadBalancingSettingsProperties: &frontdoor.LoadBalancingSettingsProperties{
				SampleSize:                    utils.Int32(sampleSize),
				SuccessfulSamplesRequired:     utils.Int32(successfulSamplesRequired),
				AdditionalLatencyMilliseconds: utils.Int32(additionalLatencyMilliseconds),
			},
		}

		output = append(output, result)
	}

	return &output
}

func expandArmFrontDoorRoutingRule(input []interface{}, subscriptionId string, resourceGroup, serviceName string) *[]frontdoor.RoutingRule {
	if len(input) == 0 {
		return nil
	}

	output := make([]frontdoor.RoutingRule, 0)

	for _, rr := range input {
		routingRule := rr.(map[string]interface{})

		id := routingRule["id"].(string)
		frontendEndpoints := routingRule["frontend_endpoints"].([]interface{})
		acceptedProtocols := routingRule["accepted_protocols"].([]interface{})
		ptm := routingRule["patterns_to_match"].([]interface{})
		enabled := routingRule["enabled"].(bool)
		name := routingRule["name"].(string)

		patternsToMatch := make([]string, 0)

		for _, p := range ptm {
			patternsToMatch = append(patternsToMatch, p.(string))
		}

		var routingConfiguration frontdoor.BasicRouteConfiguration

		if rc := routingRule["redirect_configuration"].([]interface{}); len(rc) != 0 {
			routingConfiguration = expandArmFrontDoorRedirectConfiguration(rc)
		} else if fc := routingRule["forwarding_configuration"].([]interface{}); len(fc) != 0 {
			routingConfiguration = expandArmFrontDoorForwardingConfiguration(fc, subscriptionId, resourceGroup, serviceName)
		}

		currentRoutingRule := frontdoor.RoutingRule{
			ID:   utils.String(id),
			Name: utils.String(name),
			RoutingRuleProperties: &frontdoor.RoutingRuleProperties{
				//ResourceState:
				FrontendEndpoints:  expandArmFrontDoorFrontEndEndpoints(frontendEndpoints, subscriptionId, resourceGroup, serviceName),
				AcceptedProtocols:  expandArmFrontDoorAcceptedProtocols(acceptedProtocols),
				PatternsToMatch:    &patternsToMatch,
				EnabledState:       frontdoor.RoutingRuleEnabledState(expandArmFrontDoorEnabledState(enabled)),
				RouteConfiguration: routingConfiguration,
			},
		}

		output = append(output, currentRoutingRule)
	}

	return &output
}

func expandArmFrontDoorAcceptedProtocols(input []interface{}) *[]frontdoor.Protocol {
	if len(input) == 0 {
		return &[]frontdoor.Protocol{}
	}

	output := make([]frontdoor.Protocol, 0)

	for _, ap := range input {
		result := frontdoor.HTTPS

		if ap.(string) == fmt.Sprintf("%s", frontdoor.HTTP) {
			result = frontdoor.HTTP
		}

		output = append(output, result)
	}

	return &output
}

func expandArmFrontDoorSubResource(subscriptionId string, resourceGroup string, serviceName string, resourceType string, resourceName string) *frontdoor.SubResource {
	result := frontdoor.SubResource{
		ID: utils.String(azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, serviceName, resourceType, resourceName)),
	}

	return &result
}

func expandArmFrontDoorFrontEndEndpoints(input []interface{}, subscriptionId string, resourceGroup string, serviceName string) *[]frontdoor.SubResource {
	if len(input) == 0 {
		return &[]frontdoor.SubResource{}
	}

	//v := input.(map[string]interface{})
	output := make([]frontdoor.SubResource, 0)

	for _, SubResource := range input {
		result := *expandArmFrontDoorSubResource(subscriptionId, resourceGroup, serviceName, "FrontendEndpoints", SubResource.(string))
		output = append(output, result)
	}

	return &output
}

func expandArmFrontDoorEnabledState(enabled bool) frontdoor.EnabledState {
	result := frontdoor.EnabledStateDisabled

	if enabled {
		result = frontdoor.EnabledStateEnabled
	}

	return result
}

func expandArmFrontDoorRedirectConfiguration(input []interface{}) frontdoor.RedirectConfiguration {
	if len(input) == 0 {
		return frontdoor.RedirectConfiguration{}
	}
	v := input[0].(map[string]interface{})

	redirectType := v["redirect_type"].(string)
	redirectProtocol := v["redirect_protocol"].(string)
	customHost := v["custom_host"].(string)
	customPath := v["custom_path"].(string)
	customFragment := v["custom_fragment"].(string)
	customQueryString := v["custom_query_string"].(string)

	redirectConfiguration := frontdoor.RedirectConfiguration{
		RedirectType:      frontdoor.RedirectType(redirectType),
		RedirectProtocol:  frontdoor.RedirectProtocol(redirectProtocol),
		CustomHost:        utils.String(customHost),
		CustomPath:        utils.String(customPath),
		CustomFragment:    utils.String(customFragment),
		CustomQueryString: utils.String(customQueryString),
		OdataType:         frontdoor.OdataTypeMicrosoftAzureFrontDoorModelsFrontdoorRedirectConfiguration,
	}

	return redirectConfiguration
}

func expandArmFrontDoorForwardingConfiguration(input []interface{}, subscriptionId string, resourceGroup, serviceName string) frontdoor.ForwardingConfiguration {
	if len(input) == 0 {
		return frontdoor.ForwardingConfiguration{}
	}
	v := input[0].(map[string]interface{})

	customForwardingPath := v["custom_forwarding_path"].(string)
	forwardingProtocol := v["forwarding_protocol"].(string)
	cacheUseDynamicCompression := v["cache_use_dynamic_compression"].(bool)
	cacheQueryParameterStripDirective := v["cache_query_parameter_strip_directive"].(string)
	backendPoolName := v["backend_pool_name"].(string)

	useDynamicCompression := frontdoor.DynamicCompressionEnabledDisabled

	if cacheUseDynamicCompression {
		useDynamicCompression = frontdoor.DynamicCompressionEnabledEnabled
	}

	cacheConfiguration := &frontdoor.CacheConfiguration{
		QueryParameterStripDirective: frontdoor.Query(cacheQueryParameterStripDirective),
		DynamicCompression:           useDynamicCompression,
	}

	backend := &frontdoor.SubResource{
		ID: utils.String(azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, serviceName, "BackendPools", backendPoolName)),
	}

	forwardingConfiguration := frontdoor.ForwardingConfiguration{
		ForwardingProtocol: frontdoor.ForwardingProtocol(forwardingProtocol),
		CacheConfiguration: cacheConfiguration,
		BackendPool:        backend,
		OdataType:          frontdoor.OdataTypeMicrosoftAzureFrontDoorModelsFrontdoorForwardingConfiguration,
	}

	if customForwardingPath != "" {
		forwardingConfiguration.CustomForwardingPath = utils.String(customForwardingPath)
	}

	return forwardingConfiguration
}

func flattenArmFrontDoorBackendPools(input *[]frontdoor.BackendPool) []map[string]interface{} {
	if input == nil {
		return make([]map[string]interface{}, 0)
	}

	output := make([]map[string]interface{}, 0)

	for _, v := range *input {
		result := make(map[string]interface{})

		if id := v.ID; id != nil {
			result["id"] = *id
		}

		if name := v.Name; name != nil {
			result["name"] = *name
		}

		if properties := v.BackendPoolProperties; properties != nil {
			result["backend"] = flattenArmFrontDoorBackend(properties.Backends)
			result["health_probe_name"] = flattenArmFrontDoorSubResource(properties.HealthProbeSettings, "HealthProbeSettings")
			result["load_balancing_name"] = flattenArmFrontDoorSubResource(properties.LoadBalancingSettings, "LoadBalancingSettings")
		}

		output = append(output, result)
	}

	return output
}

func flattenArmFrontDoorBackendPoolsSettings(input *frontdoor.BackendPoolsSettings) bool {
	if input == nil {
		return true
	}

	result := false

	if enforceCertificateNameCheck := frontdoor.EnforceCertificateNameCheckEnabledState(input.EnforceCertificateNameCheck); enforceCertificateNameCheck != "" {
		if enforceCertificateNameCheck == frontdoor.EnforceCertificateNameCheckEnabledStateEnabled {
			result = true
		}
	}

	return result
}

func flattenArmFrontDoorBackend(input *[]frontdoor.Backend) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		result := make(map[string]interface{})

		if address := v.Address; address != nil {
			result["address"] = *address
		}
		if backendHostHeader := v.BackendHostHeader; backendHostHeader != nil {
			result["host_header"] = *backendHostHeader
		}

		if v.EnabledState == frontdoor.Enabled {
			result["enabled"] = true
		} else {
			result["enabled"] = false
		}

		if httpPort := v.HTTPPort; httpPort != nil {
			result["http_port"] = int(*httpPort)
		}

		if httpsPort := v.HTTPSPort; httpsPort != nil {
			result["https_port"] = int(*httpsPort)
		}

		if priority := v.Priority; priority != nil {
			result["priority"] = int(*priority)
		}

		if weight := v.Weight; weight != nil {
			result["weight"] = int(*weight)
		}

		output = append(output, result)
	}

	return output
}

func flattenArmFrontDoorFrontendEndpoint(input *[]frontdoor.FrontendEndpoint) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	customHttpsConfiguration := make([]interface{}, 0)
	chc := make(map[string]interface{}, 0)

	for _, v := range *input {
		if id := v.ID; id != nil {
			result["id"] = *id
		}
		if name := v.Name; name != nil {
			result["name"] = *name
		}

		if properties := v.FrontendEndpointProperties; properties != nil {
			if hostName := properties.HostName; hostName != nil {
				result["host_name"] = *hostName
			}
			if sessionAffinityEnabled := properties.SessionAffinityEnabledState; sessionAffinityEnabled != "" {
				if sessionAffinityEnabled == frontdoor.SessionAffinityEnabledStateEnabled {
					result["session_affinity_enabled"] = true
				} else {
					result["session_affinity_enabled"] = false
				}
			}

			if sessionAffinityTtlSeconds := properties.SessionAffinityTTLSeconds; sessionAffinityTtlSeconds != nil {
				result["session_affinity_ttl_seconds"] = *sessionAffinityTtlSeconds
			}

			if properties.CustomHTTPSConfiguration != nil {
				chc["certificate_source"] = string(properties.CustomHTTPSConfiguration.CertificateSource)

				if properties.CustomHTTPSConfiguration.CertificateSource == frontdoor.CertificateSourceAzureKeyVault {
					kvcsp := properties.CustomHTTPSConfiguration.KeyVaultCertificateSourceParameters
					chc["azure_key_vault_certificate_vault_id"] = *kvcsp.Vault.ID
					chc["azure_key_vault_certificate_secret_name"] = *kvcsp.SecretName
					chc["azure_key_vault_certificate_secret_version"] = *kvcsp.SecretVersion
				}
			} else {
				// since FrontDoor is the default the API does not set this value (e.g. null) in Azure,
				// Set default value for state file
				chc["certificate_source"] = string(frontdoor.CertificateSourceFrontDoor)
			}

			customHttpsConfiguration = append(customHttpsConfiguration, chc)
			result["custom_https_configuration"] = customHttpsConfiguration

			//result["web_application_firewall_policy_link"] = flattenArmFrontDoorFrontendEndpointUpdateParameters_webApplicationFirewallPolicyLink(properties.WebApplicationFirewallPolicyLink)
		}
	}
	return []interface{}{result}
}

func flattenArmFrontDoorHealthProbeSettingsModel(input *[]frontdoor.HealthProbeSettingsModel) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	result := make(map[string]interface{})

	for _, v := range *input {
		if id := v.ID; id != nil {
			result["id"] = *id
		}
		if name := v.Name; name != nil {
			result["name"] = *name
		}
		if properties := v.HealthProbeSettingsProperties; properties != nil {
			if intervalInSeconds := properties.IntervalInSeconds; intervalInSeconds != nil {
				result["interval_in_seconds"] = *intervalInSeconds
			}
			if path := properties.Path; path != nil {
				result["path"] = *path
			}
			result["protocol"] = string(properties.Protocol)
		}
	}

	return []interface{}{result}
}

func flattenArmFrontDoorLoadBalancingSettingsModel(input *[]frontdoor.LoadBalancingSettingsModel) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	for _, v := range *input {
		if id := v.ID; id != nil {
			result["id"] = *id
		}
		if name := v.Name; name != nil {
			result["name"] = *name
		}
		if properties := v.LoadBalancingSettingsProperties; properties != nil {
			if additionalLatencyMilliseconds := properties.AdditionalLatencyMilliseconds; additionalLatencyMilliseconds != nil {
				result["additional_latency_milliseconds"] = *additionalLatencyMilliseconds
			}
			if sampleSize := properties.SampleSize; sampleSize != nil {
				result["sample_size"] = *sampleSize
			}
			if successfulSamplesRequired := properties.SuccessfulSamplesRequired; successfulSamplesRequired != nil {
				result["successful_samples_required"] = *successfulSamplesRequired
			}
		}
	}
	return []interface{}{result}
}

func flattenArmFrontDoorRoutingRule(input *[]frontdoor.RoutingRule) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	output := make([]interface{}, 0)
	result := make(map[string]interface{})

	for _, v := range *input {
		if id := v.ID; id != nil {
			result["id"] = *id
		}
		if name := v.Name; name != nil {
			result["name"] = *name
		}

		if properties := v.RoutingRuleProperties; properties != nil {
			result["accepted_protocols"] = flattenArmFrontDoorAcceptedProtocol(properties.AcceptedProtocols)

			if properties.EnabledState == frontdoor.RoutingRuleEnabledStateEnabled {
				result["enabled"] = true
			} else {
				result["enabled"] = false
			}
			result["frontend_endpoints"] = flattenArmFrontDoorFrontendEndpointsSubResources(properties.FrontendEndpoints)
			if patternsToMatch := properties.PatternsToMatch; patternsToMatch != nil {
				result["patterns_to_match"] = *patternsToMatch
			}

			brc := properties.RouteConfiguration
			if routeConfigType := azure.GetFrontDoorBasicRouteConfigurationType(brc.(interface{})); routeConfigType != "" {
				rc := make([]interface{}, 0)
				c := make(map[string]interface{})

				// there are only two types of Route Configuration
				if routeConfigType == "ForwardingConfiguration" {
					v := brc.(frontdoor.ForwardingConfiguration)

					c["backend_pool_name"] = flattenArmFrontDoorSubResource(v.BackendPool, "BackendPools")
					c["custom_forwarding_path"] = v.CustomForwardingPath
					c["forwarding_protocol"] = string(v.ForwardingProtocol)

					cc := v.CacheConfiguration
					c["cache_query_parameter_strip_directive"] = string(cc.QueryParameterStripDirective)
					if cc.DynamicCompression == frontdoor.DynamicCompressionEnabledEnabled {
						c["cache_use_dynamic_compression"] = true
					} else {
						c["cache_use_dynamic_compression"] = false
					}

					rc = append(rc, c)
					result["forwarding_configuration"] = rc
				} else {
					v := brc.(frontdoor.RedirectConfiguration)
					c["custom_fragment"] = v.CustomFragment
					c["custom_host"] = v.CustomHost
					c["custom_path"] = v.CustomPath
					c["custom_query_string"] = v.CustomQueryString
					c["redirect_protocol"] = string(v.RedirectProtocol)
					c["redirect_type"] = string(v.RedirectType)

					rc = append(rc, c)
					result["redirect_configuration"] = rc
				}
			}
		}
	}

	output = append(output, result)

	return output
}

func flattenArmFrontDoorAcceptedProtocol(input *[]frontdoor.Protocol) []string {
	if input == nil {
		return make([]string, 0)
	}

	output := make([]string, 0)
	for _, p := range *input {
		output = append(output, string(p))
	}

	return output
}

func flattenArmFrontDoorSubResource(input *frontdoor.SubResource, resourceType string) string {
	if input == nil {
		return ""
	}

	name := ""

	if id := input.ID; id != nil {
		aid, err := parseAzureResourceID(*id)
		if err != nil {
			return ""
		}
		name = aid.Path[resourceType]
	}

	return name
}

func flattenArmFrontDoorFrontendEndpointsSubResources(input *[]frontdoor.SubResource) []string {
	if input == nil {
		return make([]string, 0)
	}

	output := make([]string, 0)

	for _, v := range *input {
		name := flattenArmFrontDoorSubResource(&v, "FrontendEndpoints")
		output = append(output, name)
	}

	return output
}
