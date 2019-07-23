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
				ValidateFunc: azure.ValidateFrontDoorName,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
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
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: azure.ValidateBackendPoolRoutingRuleName,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default: true,
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
								Type: schema.TypeString,
								Default: "/*",
							},
						},
						"frontend_endpoints": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 100,
							Elem: &schema.Schema{
								Type:     schema.TypeString,
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
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.RedirectProtocolHTTPOnly),
											string(frontdoor.RedirectProtocolHTTPSOnly),
											string(frontdoor.RedirectProtocolMatchRequest),
										}, false),
									},
									"redirect_type": {
										Type:     schema.TypeString,
										Optional: true,
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
										Type:     schema.TypeString,
										Optional: true,
										Default: false,
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
										Required: true,
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
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: azure.ValidateBackendPoolRoutingRuleName,
						},
						"sample_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Default: 4,
						},
						"successful_samples_required": {
							Type:     schema.TypeInt,
							Optional: true,
							Default: 2,
						},
						"additional_latency_milliseconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Default: 2,
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
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: azure.ValidateBackendPoolRoutingRuleName,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "/",
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
							Default: 120,
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
										Default: true,
									},
									"address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"http_port": {
										Type:     schema.TypeInt,
										Required: true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"https_port": {
										Type:     schema.TypeInt,
										Required: true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"weight": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: validation.IntBetween(1, 1000),
										Default: 50,
									},
									"priority": {
										Type:     schema.TypeInt,
										Required: true,
										ValidateFunc: validation.IntBetween(1, 5),
									},
									"host_header": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
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
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: azure.ValidateBackendPoolRoutingRuleName,
						},
						"host_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"session_affinity_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default: true,
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

	//return resourceArmFrontDoorRead(d, meta)
	return nil
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
		if err := d.Set("frontend_endpoint", flattenArmFrontDoorFrontendEndpoint(properties.FrontendEndpoints)); err != nil {
			return fmt.Errorf("Error setting `frontend_endpoint`: %+v", err)
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
			ID:   utils.String(azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, serviceName, "backendPools", backendPoolName)),
			Name: utils.String(backendPoolName),
			BackendPoolProperties: &frontdoor.BackendPoolProperties{
				// ResourceState
				Backends:              expandArmFrontDoorBackend(backends),
				LoadBalancingSettings: expandArmFrontDoorSubResource(subscriptionId, resourceGroup, serviceName, "loadBalancingSettings", backendPoolLoadBalancingName),
				HealthProbeSettings:   expandArmFrontDoorSubResource(subscriptionId, resourceGroup, serviceName, "healthProbeSettings", backendPoolHealthProbeName),
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
		httpPort := backend["http_port"].(int32)
		httpsPort := backend["https_port"].(int32)
		enabled := backend["enabled"].(bool)
		priority := backend["priority"].(int32)
		weight := backend["weight"].(int32)
		hostHeader := backend["host_header"].(string)

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
		sessionAffinityTtlSeconds := frontendEndpoint["session_affinity_ttl_seconds"].(int32)
		customHttpsConfiguration := frontendEndpoint["custom_https_configuration"].([]interface{})
		name := frontendEndpoint["name"].(string)
		id := utils.String(azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, serviceName, "frontendEndpoints", name))

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
				CustomHTTPSConfiguration:     expandArmFrontDoorCustomHTTPSConfiguration(customHttpsConfiguration),
				HostName:                     utils.String(hostName),
				SessionAffinityEnabledState:  sessionAffinityEnabled,
				SessionAffinityTTLSeconds:    utils.Int32(sessionAffinityTtlSeconds),
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
	
	result := frontdoor.CustomHTTPSConfiguration {
		ProtocolType: frontdoor.ServerNameIndication,
	}

	if certSource == "AzureKeyVault" {
		vaultSecret := v["azure_key_vault_certificate_secret_name"].(string)
		vaultVersion := v["azure_key_vault_certificate_secret_version"].(string)
		vaultId := v["azure_key_vault_certificate_vault_id"].(string)

		result.CertificateSource = frontdoor.CertificateSourceAzureKeyVault
		result.KeyVaultCertificateSourceParameters = &frontdoor.KeyVaultCertificateSourceParameters{
			Vault: &frontdoor.KeyVaultCertificateSourceParametersVault {
				ID: utils.String(vaultId),
			},
			SecretName: utils.String(vaultSecret),
			SecretVersion:utils.String(vaultVersion),
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

	for _,hps := range input {
		v := hps.(map[string]interface{})

		path := v["path"].(string)
		protocol := v["protocol"].(string)
		intervalInSeconds := v["interval_in_seconds"].(int32)
		name := v["name"].(string)
	
		result := frontdoor.HealthProbeSettingsModel{
			ID:   utils.String(azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, serviceName, "healthProbeSettings", name)),
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

	output := make([]frontdoor.LoadBalancingSettingsModel,0)

	for _, lbs := range input {
		loadBalanceSetting := lbs.(map[string]interface{})

		name := loadBalanceSetting["name"].(string)
		sampleSize := loadBalanceSetting["sample_size"].(int32)
		successfulSamplesRequired := loadBalanceSetting["successful_samples_required"].(int32)
		additionalLatencyMilliseconds := loadBalanceSetting["additional_latency_milliseconds"].(int32)
		id := utils.String(azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, serviceName, "loadBalancingSettings", name))

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
		patternsToMatch := routingRule["patterns_to_match"].(*[]string)
		enabled := routingRule["enabled"].(bool)
		name := routingRule["name"].(string)

		var routingConfiguration frontdoor.BasicRouteConfiguration

		if rc := routingRule["redirect_configuration"].([]interface{}); rc != nil {
			redirectConfiguration := expandArmFrontDoorRedirectConfiguration(rc)
			routingConfiguration =  redirectConfiguration
		} else if fc := routingRule["forwarding_configuration"].([]interface{}); fc != nil {
			forwardingConfiguration := expandArmFrontDoorForwardingConfiguration(fc, subscriptionId, resourceGroup, serviceName)
			routingConfiguration = forwardingConfiguration
		}

		currentRoutingRule := frontdoor.RoutingRule{
			ID:   utils.String(id),
			Name: utils.String(name),
			RoutingRuleProperties: &frontdoor.RoutingRuleProperties{
				//ResourceState:
				FrontendEndpoints:  expandArmFrontDoorFrontEndEndpoints(frontendEndpoints, subscriptionId, resourceGroup, serviceName),
				AcceptedProtocols:  expandArmFrontDoorAcceptedProtocols(acceptedProtocols),
				PatternsToMatch:    patternsToMatch,
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

	output := make([]frontdoor.Protocol,0)

	for _,ap := range input {
		result := frontdoor.HTTPS

		if ap.(frontdoor.Protocol) == frontdoor.HTTP {
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

	v := input[0].(map[string]interface{})
	output := make([]frontdoor.SubResource, 0)

	for _, SubResource := range v {
		result := *expandArmFrontDoorSubResource(subscriptionId, resourceGroup, serviceName, "frontendEndpoints", SubResource.(string))
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
	
	redirectType := v["redirect_type"].(frontdoor.RedirectType)
	redirectProtocol := v["redirect_protocol"].(frontdoor.RedirectProtocol)
	customHost := v["custom_host"].(string)
	customPath := v["custom_path"].(string)
	customFragment := v["custom_fragment"].(string)
	customQueryString := v["custom_query_string"].(string)

	redirectConfiguration := frontdoor.RedirectConfiguration {
		RedirectType:      redirectType,
		RedirectProtocol:  redirectProtocol,
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
	forwardingProtocol := v["forwarding_protocol"].(frontdoor.ForwardingProtocol)
	cacheUseDynamicCompression := v["cache_use_dynamic_compression"].(bool)
	cacheQueryParameterStripDirective := v["cache_query_parameter_strip_directive"].(frontdoor.Query)
	backendPoolName := v["backend_pool_name"].(string)

	useDynamicCompression := frontdoor.DynamicCompressionEnabledDisabled

	if cacheUseDynamicCompression {
		useDynamicCompression = frontdoor.DynamicCompressionEnabledEnabled
	}

	cacheConfiguration  := &frontdoor.CacheConfiguration  {
		QueryParameterStripDirective: cacheQueryParameterStripDirective,
		DynamicCompression: useDynamicCompression,
	}

	backend := &frontdoor.SubResource{
		ID: utils.String(azure.GetFrontDoorSubResourceId(subscriptionId, resourceGroup, serviceName, "backendPools", backendPoolName)),
	}

	forwardingConfiguration := frontdoor.ForwardingConfiguration {
		CustomForwardingPath: utils.String(customForwardingPath),
		ForwardingProtocol:   forwardingProtocol,
		CacheConfiguration:   cacheConfiguration,
		BackendPool:          backend,
		OdataType:            frontdoor.OdataTypeMicrosoftAzureFrontDoorModelsFrontdoorForwardingConfiguration,
	}

	return forwardingConfiguration
}






























func flattenArmFrontDoorBackendPool(input *[]frontdoor.BackendPool) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	for _, v := range *input {
		if id := v.ID; id != nil {
			result["id"] = *id
		}
		if properties := v.BackendPoolProperties; properties != nil {
			result["backends"] = flattenArmFrontDoorBackend(properties.Backends)
			result["health_probe_settings"] = flattenArmFrontDoorSubResource(properties.HealthProbeSettings)
			result["load_balancing_settings"] = flattenArmFrontDoorSubResource(properties.LoadBalancingSettings)
			result["resource_state"] = string(properties.ResourceState)
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
	for _, v := range *input {
		if id := v.ID; id != nil {
			result["id"] = *id
		}
		if properties := v.FrontendEndpointProperties; properties != nil {
			if hostName := properties.HostName; hostName != nil {
				result["host_name"] = *hostName
			}
			result["resource_state"] = string(properties.ResourceState)
			result["session_affinity_enabled_state"] = string(properties.SessionAffinityEnabledState)
			if sessionAffinityTtlSeconds := properties.SessionAffinityTTLSeconds; sessionAffinityTtlSeconds != nil {
				result["session_affinity_ttl_seconds"] = *sessionAffinityTtlSeconds
			}
			result["web_application_firewall_policy_link"] = flattenArmFrontDoorFrontendEndpointUpdateParameters_webApplicationFirewallPolicyLink(properties.WebApplicationFirewallPolicyLink)
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
		if properties := v.HealthProbeSettingsProperties; properties != nil {
			if intervalInSeconds := properties.IntervalInSeconds; intervalInSeconds != nil {
				result["interval_in_seconds"] = *intervalInSeconds
			}
			if path := properties.Path; path != nil {
				result["path"] = *path
			}
			result["protocol"] = string(properties.Protocol)
			result["resource_state"] = string(properties.ResourceState)
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
		if properties := v.LoadBalancingSettingsProperties; properties != nil {
			if additionalLatencyMilliseconds := properties.AdditionalLatencyMilliseconds; additionalLatencyMilliseconds != nil {
				result["additional_latency_milliseconds"] = *additionalLatencyMilliseconds
			}
			result["resource_state"] = string(properties.ResourceState)
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

	result := make(map[string]interface{})

	for _, v := range *input {
		if id := v.ID; id != nil {
			result["id"] = *id
		}
		if properties := v.RoutingRuleProperties; properties != nil {
			if acceptedProtocols := properties.AcceptedProtocols; acceptedProtocols != nil {
				for _, ap := range *acceptedProtocols {
					result["accepted_protocols"] = string(ap)
				}
			}
				result["enabled_state"] = string(properties.EnabledState)
				result["frontend_endpoints"] = flattenArmFrontDoorSubResources(properties.FrontendEndpoints)
				if patternsToMatch := properties.PatternsToMatch; patternsToMatch != nil {
					result["patterns_to_match"] = *patternsToMatch
				}
				result["resource_state"] = string(properties.ResourceState)
		}
	}

	return []interface{}{result}
}

func flattenArmFrontDoorBackend(input *[]frontdoor.Backend) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	for _, v := range *input {
		if address := v.Address; address != nil {
			result["address"] = *address
		}
		if backendHostHeader := v.BackendHostHeader; backendHostHeader != nil {
			result["backend_host_header"] = *backendHostHeader
		}
		result["enabled_state"] = string(v.EnabledState)
		if httpPort := v.HTTPPort; httpPort != nil {
			result["http_port"] = *httpPort
		}
		if httpsPort := v.HTTPSPort; httpsPort != nil {
			result["https_port"] = *httpsPort
		}
		if priority := v.Priority; priority != nil {
			result["priority"] = *priority
		}
		if weight := v.Weight; weight != nil {
			result["weight"] = *weight
		}
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

func flattenArmFrontDoorSubResources(input *[]frontdoor.SubResource) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	for _,v := range *input {
		if id := v.ID; id != nil {
			result["id"] = *id
		}
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
