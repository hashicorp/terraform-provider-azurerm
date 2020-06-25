package frontdoor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Hour),
			Delete: schema.DefaultTimeout(6 * time.Hour),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"header_frontdoor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"load_balancer_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			// TODO: In 3.0
			// Move 'enforce_backend_pools_certificate_name_check' and 'backend_pools_send_receive_timeout_seconds'
			// into a 'backend_pool_settings' block
			"enforce_backend_pools_certificate_name_check": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"backend_pools_send_receive_timeout_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntBetween(0, 240),
			},

			// TODO: Remove in 3.0
			"location": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Due to the service's API changing 'location' must now always be set to 'Global' for new resources, however if the Front Door service was created prior 2020/03/10 it may continue to exist in a specific current location",
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"routing_rule": {
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
							Required:     true,
							ValidateFunc: validate.FrontDoorBackendPoolRoutingRuleName,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"accepted_protocols": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 2,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(frontdoor.HTTP),
									string(frontdoor.HTTPS),
								}, false),
							},
						},
						"patterns_to_match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 25,
							Elem: &schema.Schema{
								Type: schema.TypeString,
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
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.FrontDoorBackendPoolRoutingRuleName,
									},
									"cache_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"cache_use_dynamic_compression": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"cache_query_parameter_strip_directive": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(frontdoor.StripAll),
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.StripAll),
											string(frontdoor.StripNone),
										}, false),
									},
									"custom_forwarding_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"forwarding_protocol": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(frontdoor.HTTPSOnly),
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.HTTPOnly),
											string(frontdoor.HTTPSOnly),
											string(frontdoor.MatchRequest),
										}, false),
									},
								},
							},
						},
					},
				},
			},

			"backend_pool_load_balancing": {
				Type:     schema.TypeList,
				MaxItems: 5000,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.FrontDoorBackendPoolRoutingRuleName,
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
							Default:  0,
						},
					},
				},
			},

			"backend_pool_health_probe": {
				Type:     schema.TypeList,
				MaxItems: 5000,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.FrontDoorBackendPoolRoutingRuleName,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "/",
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(frontdoor.HTTP),
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.HTTP),
								string(frontdoor.HTTPS),
							}, false),
						},
						"probe_method": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(frontdoor.GET),
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.GET),
								string(frontdoor.HEAD),
							}, false),
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
										Default:      50,
										ValidateFunc: validation.IntBetween(1, 1000),
									},
									"priority": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      1,
										ValidateFunc: validation.IntBetween(1, 5),
									},
									"host_header": {
										Type:     schema.TypeString,
										Required: true,
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
							ValidateFunc: validate.FrontDoorBackendPoolRoutingRuleName,
						},
						"health_probe_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"load_balancing_name": {
							Type:     schema.TypeString,
							Required: true,
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
							Required:     true,
							ValidateFunc: validate.FrontDoorBackendPoolRoutingRuleName,
						},
						"host_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"session_affinity_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"session_affinity_ttl_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"custom_https_provisioning_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"web_application_firewall_policy_link_id": {
							Type:     schema.TypeString,
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
										Default:  string(frontdoor.CertificateSourceFrontDoor),
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.CertificateSourceAzureKeyVault),
											string(frontdoor.CertificateSourceFrontDoor),
										}, false),
									},
									"minimum_tls_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"provisioning_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"provisioning_substate": {
										Type:     schema.TypeString,
										Computed: true,
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

			"tags": tags.Schema(),
		},

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			if err := validate.FrontdoorSettings(d); err != nil {
				return fmt.Errorf("creating Front Door %q (Resource Group %q): %+v", d.Get("name").(string), d.Get("resource_group_name").(string), err)
			}

			return nil
		},
	}
}

func resourceArmFrontDoorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for present of existing Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_frontdoor", *resp.ID)
		}
	}

	// remove in 3.0
	// due to a change in the RP, if a Frontdoor exists in a location other than 'Global' it may continue to
	// exist in that location, if this is a brand new Frontdoor it must be created in the 'Global' location
	location := "Global"
	preExists := false
	cfgLocation, hasLocation := d.GetOk("location")

	exists, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(exists.Response) {
			return fmt.Errorf("locating Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	} else {
		preExists = true
		location = azure.NormalizeLocation(*exists.Location)
	}

	if hasLocation && preExists {
		if location != azure.NormalizeLocation(cfgLocation) {
			return fmt.Errorf("the Front Door %q (Resource Group %q) already exists in %q and cannot be moved to the %q location", name, resourceGroup, location, cfgLocation)
		}
	}

	frontDoorPath := fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/Microsoft.Network/Frontdoors/%s", subscriptionId, resourceGroup, name)
	friendlyName := d.Get("friendly_name").(string)
	routingRules := d.Get("routing_rule").([]interface{})
	loadBalancingSettings := d.Get("backend_pool_load_balancing").([]interface{})
	healthProbeSettings := d.Get("backend_pool_health_probe").([]interface{})
	backendPools := d.Get("backend_pool").([]interface{})
	frontendEndpoints := d.Get("frontend_endpoint").([]interface{})
	backendPoolsSettings := d.Get("enforce_backend_pools_certificate_name_check").(bool)
	backendPoolsSendReceiveTimeoutSeconds := int32(d.Get("backend_pools_send_receive_timeout_seconds").(int))
	enabledState := d.Get("load_balancer_enabled").(bool)

	t := d.Get("tags").(map[string]interface{})

	frontDoorParameters := frontdoor.FrontDoor{
		Location: utils.String(location),
		Properties: &frontdoor.Properties{
			FriendlyName:          utils.String(friendlyName),
			RoutingRules:          expandArmFrontDoorRoutingRule(routingRules, frontDoorPath),
			BackendPools:          expandArmFrontDoorBackendPools(backendPools, frontDoorPath),
			BackendPoolsSettings:  expandArmFrontDoorBackendPoolsSettings(backendPoolsSettings, backendPoolsSendReceiveTimeoutSeconds),
			FrontendEndpoints:     expandArmFrontDoorFrontendEndpoint(frontendEndpoints, frontDoorPath),
			HealthProbeSettings:   expandArmFrontDoorHealthProbeSettingsModel(healthProbeSettings, frontDoorPath),
			LoadBalancingSettings: expandArmFrontDoorLoadBalancingSettingsModel(loadBalancingSettings, frontDoorPath),
			EnabledState:          expandArmFrontDoorEnabledState(enabledState),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, frontDoorParameters)
	if err != nil {
		return fmt.Errorf("creating Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("cannot read Front Door %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	// Now loop through the FrontendEndpoints and enable/disable Custom Domain HTTPS
	// on each individual Frontend Endpoint if required
	feClient := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
	feCtx, feCancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer feCancel()

	for _, v := range frontendEndpoints {
		frontendEndpoint := v.(map[string]interface{})
		customHttpsProvisioningEnabled := frontendEndpoint["custom_https_provisioning_enabled"].(bool)
		frontendEndpointName := frontendEndpoint["name"].(string)

		// Get current state of endpoint from Azure
		resp, err := feClient.Get(feCtx, resourceGroup, name, frontendEndpointName)
		if err != nil {
			return fmt.Errorf("retrieving Front Door Frontend Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
		}
		if resp.ID == nil {
			return fmt.Errorf("cannot read Front Door Frontend Endpoint %q (Resource Group %q) ID", frontendEndpointName, resourceGroup)
		}

		if properties := resp.FrontendEndpointProperties; properties != nil {
			if provisioningState := properties.CustomHTTPSProvisioningState; provisioningState != "" {
				// Check to see if we are going to change the CustomHTTPSProvisioningState, if so check to
				// see if its current state is configurable, if not return an error...
				if customHttpsProvisioningEnabled != NormalizeCustomHTTPSProvisioningStateToBool(provisioningState) {
					if err := IsFrontDoorFrontendEndpointConfigurable(provisioningState, customHttpsProvisioningEnabled, frontendEndpointName, resourceGroup); err != nil {
						return err
					}
				}

				if customHttpsProvisioningEnabled && provisioningState == frontdoor.CustomHTTPSProvisioningStateDisabled {
					// Build a custom Https configuration based off the config file to send to the enable call
					// NOTE: I do not need to check to see if this exists since I already do that in the validation code
					chc := frontendEndpoint["custom_https_configuration"].([]interface{})
					customHTTPSConfiguration := chc[0].(map[string]interface{})
					minTLSVersion := frontdoor.OneFullStopTwo // Default to TLS 1.2
					if httpsConfig := properties.CustomHTTPSConfiguration; httpsConfig != nil {
						minTLSVersion = httpsConfig.MinimumTLSVersion
					}
					customHTTPSConfigurationUpdate := makeCustomHttpsConfiguration(customHTTPSConfiguration, minTLSVersion)
					// Enable Custom Domain HTTPS for the Frontend Endpoint
					if err := resourceArmFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx, true, name, frontendEndpointName, resourceGroup, customHTTPSConfigurationUpdate, meta); err != nil {
						return fmt.Errorf("Unable enable Custom Domain HTTPS for Frontend Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
					}
				} else if !customHttpsProvisioningEnabled && provisioningState == frontdoor.CustomHTTPSProvisioningStateEnabled {
					// Disable Custom Domain HTTPS for the Frontend Endpoint
					if err := resourceArmFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx, false, name, frontendEndpointName, resourceGroup, frontdoor.CustomHTTPSConfiguration{}, meta); err != nil {
						return fmt.Errorf("Unable to disable Custom Domain HTTPS for Frontend Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
					}
				}
			}
		}
	}

	return resourceArmFrontDoorRead(d, meta)
}

func resourceArmFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx context.Context, enableCustomHttpsProvisioning bool, frontDoorName string, frontendEndpointName string, resourceGroup string, customHTTPSConfiguration frontdoor.CustomHTTPSConfiguration, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient

	if enableCustomHttpsProvisioning {
		future, err := client.EnableHTTPS(ctx, resourceGroup, frontDoorName, frontendEndpointName, customHTTPSConfiguration)

		if err != nil {
			return fmt.Errorf("enabling Custom Domain HTTPS for Frontend Endpoint: %+v", err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting to enable Custom Domain HTTPS for Frontend Endpoint: %+v", err)
		}
	} else {
		future, err := client.DisableHTTPS(ctx, resourceGroup, frontDoorName, frontendEndpointName)

		if err != nil {
			return fmt.Errorf("disabling Custom Domain HTTPS for Frontend Endpoint: %+v", err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting to disable Custom Domain HTTPS for Frontend Endpoint: %+v", err)
		}
	}

	return nil
}

func resourceArmFrontDoorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["frontdoors"]
	// Link to issue: https://github.com/Azure/azure-sdk-for-go/issues/6762
	if name == "" {
		name = id.Path["Frontdoors"]
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Front Door %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", azure.NormalizeLocation(*resp.Location))

	if properties := resp.Properties; properties != nil {
		if err := d.Set("backend_pool", flattenArmFrontDoorBackendPools(properties.BackendPools)); err != nil {
			return fmt.Errorf("setting `backend_pool`: %+v", err)
		}

		backendPoolSettings := flattenArmFrontDoorBackendPoolsSettings(properties.BackendPoolsSettings)

		if err := d.Set("enforce_backend_pools_certificate_name_check", backendPoolSettings["enforce_backend_pools_certificate_name_check"].(bool)); err != nil {
			return fmt.Errorf("setting `enforce_backend_pools_certificate_name_check`: %+v", err)
		}

		if err := d.Set("backend_pools_send_receive_timeout_seconds", backendPoolSettings["backend_pools_send_receive_timeout_seconds"].(int32)); err != nil {
			return fmt.Errorf("setting `backend_pools_send_receive_timeout_seconds`: %+v", err)
		}

		d.Set("cname", properties.Cname)
		d.Set("header_frontdoor_id", properties.FrontdoorID)
		d.Set("load_balancer_enabled", properties.EnabledState == frontdoor.EnabledStateEnabled)
		d.Set("friendly_name", properties.FriendlyName)

		if frontendEndpoints := properties.FrontendEndpoints; frontendEndpoints != nil {
			if resp.Name != nil {
				if frontDoorFrontendEndpoints, err := flattenArmFrontDoorFrontendEndpoint(ctx, frontendEndpoints, resourceGroup, *resp.Name, meta); frontDoorFrontendEndpoints != nil {
					if err := d.Set("frontend_endpoint", frontDoorFrontendEndpoints); err != nil {
						return fmt.Errorf("setting `frontend_endpoint`: %+v", err)
					}
				} else {
					return fmt.Errorf("flattening `frontend_endpoint`: %+v", err)
				}
			} else {
				return fmt.Errorf("flattening `frontend_endpoint`: Unable to read Frontdoor Name")
			}
		}

		if err := d.Set("backend_pool_health_probe", flattenArmFrontDoorHealthProbeSettingsModel(properties.HealthProbeSettings)); err != nil {
			return fmt.Errorf("setting `backend_pool_health_probe`: %+v", err)
		}

		if err := d.Set("backend_pool_load_balancing", flattenArmFrontDoorLoadBalancingSettingsModel(properties.LoadBalancingSettings)); err != nil {
			return fmt.Errorf("setting `backend_pool_load_balancing`: %+v", err)
		}

		if err := d.Set("routing_rule", flattenArmFrontDoorRoutingRule(properties.RoutingRules, d.Get("routing_rule"))); err != nil {
			return fmt.Errorf("setting `routing_rules`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmFrontDoorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["frontdoors"]
	// Link to issue: https://github.com/Azure/azure-sdk-for-go/issues/6762
	if name == "" {
		name = id.Path["Frontdoors"]
	}

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmFrontDoorBackendPools(input []interface{}, frontDoorPath string) *[]frontdoor.BackendPool {
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
			ID:   utils.String(frontDoorPath + "/BackendPools/" + backendPoolName),
			Name: utils.String(backendPoolName),
			BackendPoolProperties: &frontdoor.BackendPoolProperties{
				Backends: expandArmFrontDoorBackend(backends),
				LoadBalancingSettings: &frontdoor.SubResource{
					ID: utils.String(frontDoorPath + "/LoadBalancingSettings/" + backendPoolLoadBalancingName),
				},
				HealthProbeSettings: &frontdoor.SubResource{
					ID: utils.String(frontDoorPath + "/HealthProbeSettings/" + backendPoolHealthProbeName),
				},
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

func expandArmFrontDoorBackendPoolsSettings(enforceCertificateNameCheck bool, backendPoolsSendReceiveTimeoutSeconds int32) *frontdoor.BackendPoolsSettings {
	enforceCheck := frontdoor.EnforceCertificateNameCheckEnabledStateDisabled

	if enforceCertificateNameCheck {
		enforceCheck = frontdoor.EnforceCertificateNameCheckEnabledStateEnabled
	}

	result := frontdoor.BackendPoolsSettings{
		EnforceCertificateNameCheck: enforceCheck,
		SendRecvTimeoutSeconds:      utils.Int32(backendPoolsSendReceiveTimeoutSeconds),
	}

	return &result
}

func expandArmFrontDoorFrontendEndpoint(input []interface{}, frontDoorPath string) *[]frontdoor.FrontendEndpoint {
	if len(input) == 0 {
		return &[]frontdoor.FrontendEndpoint{}
	}

	output := make([]frontdoor.FrontendEndpoint, 0)

	for _, frontendEndpoints := range input {
		frontendEndpoint := frontendEndpoints.(map[string]interface{})

		hostName := frontendEndpoint["host_name"].(string)
		isSessionAffinityEnabled := frontendEndpoint["session_affinity_enabled"].(bool)
		sessionAffinityTtlSeconds := int32(frontendEndpoint["session_affinity_ttl_seconds"].(int))
		waf := frontendEndpoint["web_application_firewall_policy_link_id"].(string)
		name := frontendEndpoint["name"].(string)
		id := utils.String(frontDoorPath + "/FrontendEndpoints/" + name)

		sessionAffinityEnabled := frontdoor.SessionAffinityEnabledStateDisabled
		if isSessionAffinityEnabled {
			sessionAffinityEnabled = frontdoor.SessionAffinityEnabledStateEnabled
		}

		result := frontdoor.FrontendEndpoint{
			ID:   id,
			Name: utils.String(name),
			FrontendEndpointProperties: &frontdoor.FrontendEndpointProperties{
				HostName:                    utils.String(hostName),
				SessionAffinityEnabledState: sessionAffinityEnabled,
				SessionAffinityTTLSeconds:   utils.Int32(sessionAffinityTtlSeconds),
			},
		}

		if waf != "" {
			result.FrontendEndpointProperties.WebApplicationFirewallPolicyLink = &frontdoor.FrontendEndpointUpdateParametersWebApplicationFirewallPolicyLink{
				ID: utils.String(waf),
			}
		}

		output = append(output, result)
	}

	return &output
}

func expandArmFrontDoorHealthProbeSettingsModel(input []interface{}, frontDoorPath string) *[]frontdoor.HealthProbeSettingsModel {
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
		enabled := v["enabled"].(bool)

		healthProbeEnabled := frontdoor.HealthProbeEnabledEnabled
		if !enabled {
			healthProbeEnabled = frontdoor.HealthProbeEnabledDisabled
		}

		result := frontdoor.HealthProbeSettingsModel{
			ID:   utils.String(frontDoorPath + "/HealthProbeSettings/" + name),
			Name: utils.String(name),
			HealthProbeSettingsProperties: &frontdoor.HealthProbeSettingsProperties{
				IntervalInSeconds: utils.Int32(intervalInSeconds),
				Path:              utils.String(path),
				Protocol:          frontdoor.Protocol(protocol),
				HealthProbeMethod: frontdoor.HealthProbeMethod(v["probe_method"].(string)),
				EnabledState:      healthProbeEnabled,
			},
		}

		output = append(output, result)
	}

	return &output
}

func expandArmFrontDoorLoadBalancingSettingsModel(input []interface{}, frontDoorPath string) *[]frontdoor.LoadBalancingSettingsModel {
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
		id := utils.String(frontDoorPath + "/LoadBalancingSettings/" + name)

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

func expandArmFrontDoorRoutingRule(input []interface{}, frontDoorPath string) *[]frontdoor.RoutingRule {
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
			routingConfiguration = expandArmFrontDoorForwardingConfiguration(fc, frontDoorPath)
		}

		currentRoutingRule := frontdoor.RoutingRule{
			ID:   utils.String(id),
			Name: utils.String(name),
			RoutingRuleProperties: &frontdoor.RoutingRuleProperties{
				FrontendEndpoints:  expandArmFrontDoorFrontEndEndpoints(frontendEndpoints, frontDoorPath),
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

		if ap.(string) == string(frontdoor.HTTP) {
			result = frontdoor.HTTP
		}

		output = append(output, result)
	}

	return &output
}

func expandArmFrontDoorFrontEndEndpoints(input []interface{}, frontDoorPath string) *[]frontdoor.SubResource {
	if len(input) == 0 {
		return &[]frontdoor.SubResource{}
	}

	output := make([]frontdoor.SubResource, 0)

	for _, SubResource := range input {
		result := frontdoor.SubResource{
			ID: utils.String(frontDoorPath + "/FrontendEndpoints/" + SubResource.(string)),
		}
		output = append(output, result)
	}

	return &output
}

func expandArmFrontDoorEnabledState(enabled bool) frontdoor.EnabledState {
	if enabled {
		return frontdoor.EnabledStateEnabled
	}

	return frontdoor.EnabledStateDisabled
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
		CustomHost:       utils.String(customHost),
		RedirectType:     frontdoor.RedirectType(redirectType),
		RedirectProtocol: frontdoor.RedirectProtocol(redirectProtocol),
		OdataType:        frontdoor.OdataTypeMicrosoftAzureFrontDoorModelsFrontdoorRedirectConfiguration,
	}

	// The way the API works is if you don't include the attribute in the structure
	// it is treated as Preserve instead of Replace...
	if customHost != "" {
		redirectConfiguration.CustomHost = utils.String(customHost)
	}
	if customPath != "" {
		redirectConfiguration.CustomPath = utils.String(customPath)
	}
	if customFragment != "" {
		redirectConfiguration.CustomFragment = utils.String(customFragment)
	}
	if customQueryString != "" {
		redirectConfiguration.CustomQueryString = utils.String(customQueryString)
	}

	return redirectConfiguration
}

func expandArmFrontDoorForwardingConfiguration(input []interface{}, frontDoorPath string) frontdoor.ForwardingConfiguration {
	if len(input) == 0 {
		return frontdoor.ForwardingConfiguration{}
	}
	v := input[0].(map[string]interface{})

	customForwardingPath := v["custom_forwarding_path"].(string)
	forwardingProtocol := v["forwarding_protocol"].(string)
	backendPoolName := v["backend_pool_name"].(string)
	cacheUseDynamicCompression := v["cache_use_dynamic_compression"].(bool)
	cacheQueryParameterStripDirective := v["cache_query_parameter_strip_directive"].(string)
	cacheEnabled := v["cache_enabled"].(bool)

	backend := &frontdoor.SubResource{
		ID: utils.String(frontDoorPath + "/BackendPools/" + backendPoolName),
	}

	forwardingConfiguration := frontdoor.ForwardingConfiguration{
		ForwardingProtocol: frontdoor.ForwardingProtocol(forwardingProtocol),
		BackendPool:        backend,
		OdataType:          frontdoor.OdataTypeMicrosoftAzureFrontDoorModelsFrontdoorForwardingConfiguration,
	}

	// Per the portal, if you enable the cache the cache_query_parameter_strip_directive
	// is then a required attribute else the CacheConfiguration type is null
	if cacheEnabled {
		// Set the default value for dynamic compression or use the value defined in the config
		dynamicCompression := frontdoor.DynamicCompressionEnabledEnabled
		if !cacheUseDynamicCompression {
			dynamicCompression = frontdoor.DynamicCompressionEnabledDisabled
		}

		if cacheQueryParameterStripDirective == "" {
			// Set Default Value for strip directive is not in the key slice and cache is enabled
			cacheQueryParameterStripDirective = string(frontdoor.StripAll)
		}

		forwardingConfiguration.CacheConfiguration = &frontdoor.CacheConfiguration{
			DynamicCompression:           dynamicCompression,
			QueryParameterStripDirective: frontdoor.Query(cacheQueryParameterStripDirective),
		}
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
			// Link to issue: https://github.com/Azure/azure-sdk-for-go/issues/6762
			if result["health_probe_name"] == "" {
				result["health_probe_name"] = flattenArmFrontDoorSubResource(properties.HealthProbeSettings, "healthProbeSettings")
			}

			result["load_balancing_name"] = flattenArmFrontDoorSubResource(properties.LoadBalancingSettings, "LoadBalancingSettings")
			// Link to issue: https://github.com/Azure/azure-sdk-for-go/issues/6762
			if result["load_balancing_name"] == "" {
				result["load_balancing_name"] = flattenArmFrontDoorSubResource(properties.LoadBalancingSettings, "loadBalancingSettings")
			}
		}
		output = append(output, result)
	}

	return output
}

func flattenArmFrontDoorBackendPoolsSettings(input *frontdoor.BackendPoolsSettings) map[string]interface{} {
	result := make(map[string]interface{})

	// Set default values
	result["enforce_backend_pools_certificate_name_check"] = true
	result["backend_pools_send_receive_timeout_seconds"] = int32(60)

	if input == nil {
		return result
	}

	result["enforce_backend_pools_certificate_name_check"] = false

	if enforceCertificateNameCheck := input.EnforceCertificateNameCheck; enforceCertificateNameCheck != "" && enforceCertificateNameCheck == frontdoor.EnforceCertificateNameCheckEnabledStateEnabled {
		result["enforce_backend_pools_certificate_name_check"] = true
	}

	if sendRecvTimeoutSeconds := input.SendRecvTimeoutSeconds; sendRecvTimeoutSeconds != nil {
		result["backend_pools_send_receive_timeout_seconds"] = *sendRecvTimeoutSeconds
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

		result["enabled"] = v.EnabledState == frontdoor.Enabled

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

func flattenArmFrontDoorFrontendEndpoint(ctx context.Context, input *[]frontdoor.FrontendEndpoint, resourceGroup string, frontDoorName string, meta interface{}) ([]interface{}, error) {
	if input == nil {
		return make([]interface{}, 0), fmt.Errorf("cannot read Front Door Frontend Endpoint (Resource Group %q): slice is empty", resourceGroup)
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		result := make(map[string]interface{})
		customHttpsConfiguration := make([]interface{}, 0)
		chc := make(map[string]interface{})

		if name := v.Name; name != nil {
			result["name"] = *name

			// Need to call frontEndEndpointClient here to get customConfiguration information from that client
			// because the information is hidden from the main frontDoorClient "by design"...
			client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient

			resp, err := client.Get(ctx, resourceGroup, frontDoorName, *name)
			if err != nil {
				return make([]interface{}, 0), fmt.Errorf("retrieving Front Door Frontend Endpoint Custom HTTPS Configuration %q (Resource Group %q): %+v", *name, resourceGroup, err)
			}
			if resp.ID == nil {
				return make([]interface{}, 0), fmt.Errorf("cannot read Front Door Frontend Endpoint Custom HTTPS Configuration %q (Resource Group %q) ID", *name, resourceGroup)
			}

			result["id"] = resp.ID

			if properties := resp.FrontendEndpointProperties; properties != nil {
				if hostName := properties.HostName; hostName != nil {
					result["host_name"] = *hostName
				}

				if sessionAffinityEnabled := properties.SessionAffinityEnabledState; sessionAffinityEnabled != "" {
					result["session_affinity_enabled"] = sessionAffinityEnabled == frontdoor.SessionAffinityEnabledStateEnabled
				}

				if sessionAffinityTtlSeconds := properties.SessionAffinityTTLSeconds; sessionAffinityTtlSeconds != nil {
					result["session_affinity_ttl_seconds"] = *sessionAffinityTtlSeconds
				}

				if waf := properties.WebApplicationFirewallPolicyLink; waf != nil {
					result["web_application_firewall_policy_link_id"] = *waf.ID
				}

				if properties.CustomHTTPSConfiguration != nil {
					customHTTPSConfiguration := properties.CustomHTTPSConfiguration
					if customHTTPSConfiguration.CertificateSource == frontdoor.CertificateSourceAzureKeyVault {
						if kvcsp := customHTTPSConfiguration.KeyVaultCertificateSourceParameters; kvcsp != nil {
							chc["certificate_source"] = string(frontdoor.CertificateSourceAzureKeyVault)
							chc["azure_key_vault_certificate_vault_id"] = *kvcsp.Vault.ID
							chc["azure_key_vault_certificate_secret_name"] = *kvcsp.SecretName
							chc["azure_key_vault_certificate_secret_version"] = *kvcsp.SecretVersion
						}
					} else {
						chc["certificate_source"] = string(frontdoor.CertificateSourceFrontDoor)
					}

					chc["minimum_tls_version"] = string(customHTTPSConfiguration.MinimumTLSVersion)

					if provisioningState := properties.CustomHTTPSProvisioningState; provisioningState != "" {
						chc["provisioning_state"] = provisioningState
						if provisioningState == frontdoor.CustomHTTPSProvisioningStateEnabled || provisioningState == frontdoor.CustomHTTPSProvisioningStateEnabling {
							result["custom_https_provisioning_enabled"] = true

							if provisioningSubstate := properties.CustomHTTPSProvisioningSubstate; provisioningSubstate != "" {
								chc["provisioning_substate"] = provisioningSubstate
							}
						} else {
							result["custom_https_provisioning_enabled"] = false
						}

						customHttpsConfiguration = append(customHttpsConfiguration, chc)
						result["custom_https_configuration"] = customHttpsConfiguration
					}
				}
			}
		}

		output = append(output, result)
	}

	return output, nil
}

func flattenArmFrontDoorHealthProbeSettingsModel(input *[]frontdoor.HealthProbeSettingsModel) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		result := make(map[string]interface{})
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
			if healthProbeMethod := properties.HealthProbeMethod; healthProbeMethod != "" {
				// I have to upper this as the frontdoor.GET and frondoor.HEAD types are uppercased
				// but Azure stores them in the resource as pascal cased (e.g. "Get" and "Head")
				result["probe_method"] = strings.ToUpper(string(healthProbeMethod))
			}
			if enabled := properties.EnabledState; enabled != "" {
				result["enabled"] = enabled == frontdoor.HealthProbeEnabledEnabled
			}
			result["protocol"] = string(properties.Protocol)
		}

		results = append(results, result)
	}

	return results
}

func flattenArmFrontDoorLoadBalancingSettingsModel(input *[]frontdoor.LoadBalancingSettingsModel) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		result := make(map[string]interface{})
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

		results = append(results, result)
	}

	return results
}

func flattenArmFrontDoorRoutingRule(input *[]frontdoor.RoutingRule, oldBlocks interface{}) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	oldByName := map[string]map[string]interface{}{}

	for _, i := range oldBlocks.([]interface{}) {
		v := i.(map[string]interface{})

		oldByName[v["name"].(string)] = v
	}

	output := make([]interface{}, 0)
	for _, v := range *input {
		result := make(map[string]interface{})

		if id := v.ID; id != nil {
			result["id"] = *id
		}

		name := v.Name
		if name != nil {
			result["name"] = *name
		}

		if properties := v.RoutingRuleProperties; properties != nil {
			result["accepted_protocols"] = flattenArmFrontDoorAcceptedProtocol(properties.AcceptedProtocols)
			result["enabled"] = properties.EnabledState == frontdoor.RoutingRuleEnabledStateEnabled
			result["frontend_endpoints"] = flattenArmFrontDoorFrontendEndpointsSubResources(properties.FrontendEndpoints)
			if patternsToMatch := properties.PatternsToMatch; patternsToMatch != nil {
				result["patterns_to_match"] = *patternsToMatch
			}

			brc := properties.RouteConfiguration
			if routeConfigType := GetFrontDoorBasicRouteConfigurationType(brc.(interface{})); routeConfigType != "" {
				rc := make([]interface{}, 0)
				c := make(map[string]interface{})

				// there are only two types of Route Configuration
				if routeConfigType == "ForwardingConfiguration" {
					v := brc.(frontdoor.ForwardingConfiguration)

					c["backend_pool_name"] = flattenArmFrontDoorSubResource(v.BackendPool, "BackendPools")
					// Link to issue: https://github.com/Azure/azure-sdk-for-go/issues/6762
					if c["backend_pool_name"] == "" {
						c["backend_pool_name"] = flattenArmFrontDoorSubResource(v.BackendPool, "backendPools")
					}
					c["custom_forwarding_path"] = v.CustomForwardingPath
					c["forwarding_protocol"] = string(v.ForwardingProtocol)

					if cacheConfiguration := v.CacheConfiguration; cacheConfiguration != nil {
						c["cache_enabled"] = true
						if stripDirective := cacheConfiguration.QueryParameterStripDirective; stripDirective != "" {
							c["cache_query_parameter_strip_directive"] = string(stripDirective)
						} else {
							// Default value if not set
							c["cache_query_parameter_strip_directive"] = string(frontdoor.StripAll)
						}

						if dynamicCompression := cacheConfiguration.DynamicCompression; dynamicCompression != "" {
							c["cache_use_dynamic_compression"] = string(dynamicCompression) == string(frontdoor.DynamicCompressionEnabledEnabled)
						}
					} else {
						// if the cache is disabled, set the default values or revert to what they were in the previous plan
						c["cache_enabled"] = false
						c["cache_query_parameter_strip_directive"] = string(frontdoor.StripAll)

						if name != nil {
							// get `forwarding_configuration`
							if o, ok := oldByName[*name]; ok {
								ofcs := o["forwarding_configuration"].([]interface{})
								if len(ofcs) > 0 {
									ofc := ofcs[0].(map[string]interface{})
									c["cache_query_parameter_strip_directive"] = ofc["cache_query_parameter_strip_directive"]
									c["cache_use_dynamic_compression"] = ofc["cache_use_dynamic_compression"]
								}
							}
						}
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

		output = append(output, result)
	}

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
		aid, err := azure.ParseAzureResourceID(*id)
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
		// Link to issue: https://github.com/Azure/azure-sdk-for-go/issues/6762
		if name == "" {
			name = flattenArmFrontDoorSubResource(&v, "frontendEndpoints")
		}
		output = append(output, name)
	}

	return output
}

func makeCustomHttpsConfiguration(customHttpsConfiguration map[string]interface{}, minTLSVersion frontdoor.MinimumTLSVersion) frontdoor.CustomHTTPSConfiguration {
	// https://github.com/Azure/azure-sdk-for-go/issues/6882
	defaultProtocolType := "ServerNameIndication"

	customHTTPSConfigurationUpdate := frontdoor.CustomHTTPSConfiguration{
		ProtocolType:      &defaultProtocolType,
		MinimumTLSVersion: minTLSVersion,
	}

	if customHttpsConfiguration["certificate_source"].(string) == "AzureKeyVault" {
		vaultSecret := customHttpsConfiguration["azure_key_vault_certificate_secret_name"].(string)
		vaultVersion := customHttpsConfiguration["azure_key_vault_certificate_secret_version"].(string)
		vaultId := customHttpsConfiguration["azure_key_vault_certificate_vault_id"].(string)

		customHTTPSConfigurationUpdate.CertificateSource = frontdoor.CertificateSourceAzureKeyVault
		customHTTPSConfigurationUpdate.KeyVaultCertificateSourceParameters = &frontdoor.KeyVaultCertificateSourceParameters{
			Vault: &frontdoor.KeyVaultCertificateSourceParametersVault{
				ID: utils.String(vaultId),
			},
			SecretName:    utils.String(vaultSecret),
			SecretVersion: utils.String(vaultVersion),
		}
	} else {
		customHTTPSConfigurationUpdate.CertificateSource = frontdoor.CertificateSourceFrontDoor
		customHTTPSConfigurationUpdate.CertificateSourceParameters = &frontdoor.CertificateSourceParameters{
			CertificateType: frontdoor.Dedicated,
		}
	}

	return customHTTPSConfigurationUpdate
}
