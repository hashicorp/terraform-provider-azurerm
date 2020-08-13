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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// TODO: a state migration to patch the ID's

func resourceArmFrontDoor() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFrontDoorCreateUpdate,
		Read:   resourceArmFrontDoorRead,
		Update: resourceArmFrontDoorCreateUpdate,
		Delete: resourceArmFrontDoorDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.FrontDoorID(id)
			return err
		}),

		SchemaVersion: 1,

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
							Type:       schema.TypeBool,
							Optional:   true,
							Computed:   true,
							Deprecated: "Deprecated in favour of `azurerm_frontdoor_custom_https_configuration` resource",
						},
						"web_application_firewall_policy_link_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"custom_https_configuration": {
							Type:       schema.TypeList,
							Optional:   true,
							Computed:   true,
							MaxItems:   1,
							Deprecated: "Deprecated in favour of `azurerm_frontdoor_custom_https_configuration` resource",
							Elem: &schema.Resource{
								Schema: schemaCustomHttpsConfiguration(),
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

	frontDoorId := parse.NewFrontDoorID(resourceGroup, name)

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
			RoutingRules:          expandArmFrontDoorRoutingRule(routingRules, frontDoorId, subscriptionId),
			BackendPools:          expandArmFrontDoorBackendPools(backendPools, frontDoorId, subscriptionId),
			BackendPoolsSettings:  expandArmFrontDoorBackendPoolsSettings(backendPoolsSettings, backendPoolsSendReceiveTimeoutSeconds),
			FrontendEndpoints:     expandArmFrontDoorFrontendEndpoint(frontendEndpoints, frontDoorId, subscriptionId),
			HealthProbeSettings:   expandArmFrontDoorHealthProbeSettingsModel(healthProbeSettings, frontDoorId, subscriptionId),
			LoadBalancingSettings: expandArmFrontDoorLoadBalancingSettingsModel(loadBalancingSettings, frontDoorId, subscriptionId),
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

	d.SetId(frontDoorId.ID(subscriptionId))

	// Now loop through the FrontendEndpoints and enable/disable Custom Domain HTTPS
	// on each individual Frontend Endpoint if required
	feClient := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient

	for _, v := range frontendEndpoints {
		frontendEndpoint := v.(map[string]interface{})
		customHttpsProvisioningEnabled := frontendEndpoint["custom_https_provisioning_enabled"].(bool)
		endpointName := frontendEndpoint["name"].(string)

		// Get current state of endpoint from Azure
		resp, err := feClient.Get(ctx, resourceGroup, name, endpointName)
		if err != nil {
			return fmt.Errorf("retrieving Front Door Frontend Endpoint %q (Resource Group %q): %+v", endpointName, resourceGroup, err)
		}

		if properties := resp.FrontendEndpointProperties; properties != nil {
			frontendClient := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
			customHttpsConfigurationNew := frontendEndpoint["custom_https_configuration"].([]interface{})
			frontendInputId := parse.NewFrontendEndpointID(frontDoorId, endpointName)
			input := customHttpsConfigurationUpdateInput{
				customHttpsConfigurationCurrent: properties.CustomHTTPSConfiguration,
				customHttpsConfigurationNew:     customHttpsConfigurationNew,
				customHttpsProvisioningEnabled:  customHttpsProvisioningEnabled,
				frontendEndpointId:              frontendInputId,
				provisioningState:               properties.CustomHTTPSProvisioningState,
				subscriptionId:                  subscriptionId,
			}
			if err := updateCustomHttpsConfiguration(ctx, frontendClient, input); err != nil {
				return fmt.Errorf("updating Custom HTTPS configuration for Frontend Endpoint %q (Front Door %q / Resource Group %q): %+v", endpointName, name, resourceGroup, err)
			}
		}
	}

	return resourceArmFrontDoorRead(d, meta)
}

func resourceArmFrontDoorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Front Door %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Front Door %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", azure.NormalizeLocation(*resp.Location))

	if props := resp.Properties; props != nil {
		flattenedBackendPools, err := flattenArmFrontDoorBackendPools(props.BackendPools)
		if err != nil {
			return fmt.Errorf("flattening `backend_pool`: %+v", err)
		}
		if err := d.Set("backend_pool", flattenedBackendPools); err != nil {
			return fmt.Errorf("setting `backend_pool`: %+v", err)
		}

		backendPoolSettings := flattenArmFrontDoorBackendPoolsSettings(props.BackendPoolsSettings)

		if err := d.Set("enforce_backend_pools_certificate_name_check", backendPoolSettings["enforce_backend_pools_certificate_name_check"].(bool)); err != nil {
			return fmt.Errorf("setting `enforce_backend_pools_certificate_name_check`: %+v", err)
		}

		if err := d.Set("backend_pools_send_receive_timeout_seconds", backendPoolSettings["backend_pools_send_receive_timeout_seconds"].(int32)); err != nil {
			return fmt.Errorf("setting `backend_pools_send_receive_timeout_seconds`: %+v", err)
		}

		d.Set("cname", props.Cname)
		d.Set("header_frontdoor_id", props.FrontdoorID)
		d.Set("load_balancer_enabled", props.EnabledState == frontdoor.EnabledStateEnabled)
		d.Set("friendly_name", props.FriendlyName)

		// Need to call frontEndEndpointClient here to get custom(HTTPS)Configuration information from that client
		// because the information is hidden from the main frontDoorClient "by design"...
		frontEndEndpointsClient := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
		frontEndEndpointInfo, err := retrieveFrontEndEndpointInformation(ctx, frontEndEndpointsClient, *id, props.FrontendEndpoints)
		if err != nil {
			return fmt.Errorf("retrieving FrontEnd Endpoint Custom HTTPS Information: %+v", err)
		}
		frontDoorFrontendEndpoints, err := flattenFrontEndEndpoints(frontEndEndpointInfo)
		if err != nil {
			return fmt.Errorf("flattening `frontend_endpoint`: %+v", err)
		}
		if err := d.Set("frontend_endpoint", frontDoorFrontendEndpoints); err != nil {
			return fmt.Errorf("setting `frontend_endpoint`: %+v", err)
		}

		if err := d.Set("backend_pool_health_probe", flattenArmFrontDoorHealthProbeSettingsModel(props.HealthProbeSettings)); err != nil {
			return fmt.Errorf("setting `backend_pool_health_probe`: %+v", err)
		}

		if err := d.Set("backend_pool_load_balancing", flattenArmFrontDoorLoadBalancingSettingsModel(props.LoadBalancingSettings)); err != nil {
			return fmt.Errorf("setting `backend_pool_load_balancing`: %+v", err)
		}

		flattenedRoutingRules, err := flattenArmFrontDoorRoutingRule(props.RoutingRules, d.Get("routing_rule"))
		if err != nil {
			return fmt.Errorf("flattening `routing_rules`: %+v", err)
		}
		if err := d.Set("routing_rule", flattenedRoutingRules); err != nil {
			return fmt.Errorf("setting `routing_rules`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmFrontDoorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting Front Door %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting Front Door %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandArmFrontDoorBackendPools(input []interface{}, frontDoorId parse.FrontDoorId, subscriptionId string) *[]frontdoor.BackendPool {
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

		backendPoolId := parse.NewBackendPoolID(frontDoorId, backendPoolName).ID(subscriptionId)
		healthProbeId := parse.NewHealthProbeID(frontDoorId, backendPoolHealthProbeName).ID(subscriptionId)
		loadBalancingId := parse.NewLoadBalancingID(frontDoorId, backendPoolLoadBalancingName).ID(subscriptionId)

		result := frontdoor.BackendPool{
			ID:   utils.String(backendPoolId),
			Name: utils.String(backendPoolName),
			BackendPoolProperties: &frontdoor.BackendPoolProperties{
				Backends: expandArmFrontDoorBackend(backends),
				HealthProbeSettings: &frontdoor.SubResource{
					ID: utils.String(healthProbeId),
				},
				LoadBalancingSettings: &frontdoor.SubResource{
					ID: utils.String(loadBalancingId),
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

func expandArmFrontDoorFrontendEndpoint(input []interface{}, frontDoorId parse.FrontDoorId, subscriptionId string) *[]frontdoor.FrontendEndpoint {
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
		id := parse.NewFrontendEndpointID(frontDoorId, name).ID(subscriptionId)

		sessionAffinityEnabled := frontdoor.SessionAffinityEnabledStateDisabled
		if isSessionAffinityEnabled {
			sessionAffinityEnabled = frontdoor.SessionAffinityEnabledStateEnabled
		}

		result := frontdoor.FrontendEndpoint{
			ID:   utils.String(id),
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

func expandArmFrontDoorHealthProbeSettingsModel(input []interface{}, frontDoorId parse.FrontDoorId, subscriptionId string) *[]frontdoor.HealthProbeSettingsModel {
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

		healthProbeId := parse.NewHealthProbeID(frontDoorId, name).ID(subscriptionId)

		result := frontdoor.HealthProbeSettingsModel{
			ID:   utils.String(healthProbeId),
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

func expandArmFrontDoorLoadBalancingSettingsModel(input []interface{}, frontDoorId parse.FrontDoorId, subscriptionId string) *[]frontdoor.LoadBalancingSettingsModel {
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
		loadBalancingId := parse.NewLoadBalancingID(frontDoorId, name).ID(subscriptionId)

		result := frontdoor.LoadBalancingSettingsModel{
			ID:   utils.String(loadBalancingId),
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

func expandArmFrontDoorRoutingRule(input []interface{}, frontDoorId parse.FrontDoorId, subscriptionId string) *[]frontdoor.RoutingRule {
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
			routingConfiguration = expandArmFrontDoorForwardingConfiguration(fc, frontDoorId, subscriptionId)
		}

		currentRoutingRule := frontdoor.RoutingRule{
			ID:   utils.String(id),
			Name: utils.String(name),
			RoutingRuleProperties: &frontdoor.RoutingRuleProperties{
				FrontendEndpoints:  expandArmFrontDoorFrontEndEndpoints(frontendEndpoints, frontDoorId, subscriptionId),
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

func expandArmFrontDoorFrontEndEndpoints(input []interface{}, frontDoorId parse.FrontDoorId, subscriptionId string) *[]frontdoor.SubResource {
	if len(input) == 0 {
		return &[]frontdoor.SubResource{}
	}

	output := make([]frontdoor.SubResource, 0)
	for _, name := range input {
		frontendEndpointId := parse.NewFrontendEndpointID(frontDoorId, name.(string)).ID(subscriptionId)
		result := frontdoor.SubResource{
			ID: utils.String(frontendEndpointId),
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

func expandArmFrontDoorForwardingConfiguration(input []interface{}, frontDoorId parse.FrontDoorId, subscriptionId string) frontdoor.ForwardingConfiguration {
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

	backendPoolId := parse.NewBackendPoolID(frontDoorId, backendPoolName).ID(subscriptionId)
	backend := &frontdoor.SubResource{
		ID: utils.String(backendPoolId),
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

func flattenArmFrontDoorBackendPools(input *[]frontdoor.BackendPool) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	output := make([]interface{}, 0)
	for _, v := range *input {
		id := ""
		if v.ID != nil {
			id = *v.ID
		}

		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		backend := make([]interface{}, 0)
		healthProbeName := ""
		loadBalancingName := ""

		if props := v.BackendPoolProperties; props != nil {
			backend = flattenArmFrontDoorBackend(props.Backends)

			if props.HealthProbeSettings != nil && props.HealthProbeSettings.ID != nil {
				name, err := parse.HealthProbeID(*props.HealthProbeSettings.ID)
				if err != nil {
					return nil, err
				}

				healthProbeName = name.Name
			}

			if props.LoadBalancingSettings != nil && props.LoadBalancingSettings.ID != nil {
				name, err := parse.LoadBalancingID(*props.LoadBalancingSettings.ID)
				if err != nil {
					return nil, err
				}

				loadBalancingName = name.Name
			}
		}
		output = append(output, map[string]interface{}{
			"backend":             backend,
			"health_probe_name":   healthProbeName,
			"id":                  id,
			"load_balancing_name": loadBalancingName,
			"name":                name,
		})
	}

	return &output, nil
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

func retrieveFrontEndEndpointInformation(ctx context.Context, client *frontdoor.FrontendEndpointsClient, frontDoorId parse.FrontDoorId, endpoints *[]frontdoor.FrontendEndpoint) (*[]frontdoor.FrontendEndpoint, error) {
	output := make([]frontdoor.FrontendEndpoint, 0)
	if endpoints == nil {
		return &output, nil
	}

	for _, endpoint := range *endpoints {
		if endpoint.Name == nil {
			continue
		}

		name := *endpoint.Name
		resp, err := client.Get(ctx, frontDoorId.ResourceGroup, frontDoorId.Name, name)
		if err != nil {
			return nil, fmt.Errorf("retrieving Custom HTTPS Configuration for Frontend Endpoint %q (FrontDoor %q / Resource Group %q): %+v", name, frontDoorId.Name, frontDoorId.ResourceGroup, err)
		}

		output = append(output, resp)
	}

	return &output, nil
}

func flattenFrontEndEndpoints(input *[]frontdoor.FrontendEndpoint) (*[]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return &results, nil
	}

	for _, item := range *input {
		id := ""
		if item.ID != nil {
			id = *item.ID
		}
		name := ""
		if item.Name != nil {
			name = *item.Name
		}

		customHTTPSConfiguration := make([]interface{}, 0)
		customHttpsProvisioningEnabled := false
		hostName := ""
		sessionAffinityEnabled := false
		sessionAffinityTlsSeconds := 0
		webApplicationFirewallPolicyLinkId := ""
		if props := item.FrontendEndpointProperties; props != nil {
			if props.HostName != nil {
				hostName = *props.HostName
			}

			if props.SessionAffinityEnabledState != "" {
				sessionAffinityEnabled = props.SessionAffinityEnabledState == frontdoor.SessionAffinityEnabledStateEnabled
			}

			if props.SessionAffinityTTLSeconds != nil {
				sessionAffinityTlsSeconds = int(*props.SessionAffinityTTLSeconds)
			}

			if waf := props.WebApplicationFirewallPolicyLink; waf != nil && waf.ID != nil {
				webApplicationFirewallPolicyLinkId = *waf.ID
			}

			flattenedHttpsConfig := flattenCustomHttpsConfiguration(props)
			customHTTPSConfiguration = flattenedHttpsConfig.CustomHTTPSConfiguration
			customHttpsProvisioningEnabled = flattenedHttpsConfig.CustomHTTPSProvisioningEnabled
		}

		results = append(results, map[string]interface{}{
			"custom_https_configuration":        customHTTPSConfiguration,
			"custom_https_provisioning_enabled": customHttpsProvisioningEnabled,
			"host_name":                         hostName,
			"id":                                id,
			"name":                              name,
			"session_affinity_enabled":          sessionAffinityEnabled,
			"session_affinity_ttl_seconds":      sessionAffinityTlsSeconds,
			"web_application_firewall_policy_link_id": webApplicationFirewallPolicyLinkId,
		})
	}

	return &results, nil
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

func flattenArmFrontDoorRoutingRule(input *[]frontdoor.RoutingRule, oldBlocks interface{}) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	output := make([]interface{}, 0)
	for _, v := range *input {
		id := ""
		if v.ID != nil {
			id = *v.ID
		}

		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		acceptedProtocols := make([]string, 0)
		enabled := false
		forwardingConfiguration := make([]interface{}, 0)
		frontEndEndpoints := make([]string, 0)
		patternsToMatch := make([]string, 0)
		redirectConfiguration := make([]interface{}, 0)
		if props := v.RoutingRuleProperties; props != nil {
			acceptedProtocols = flattenArmFrontDoorAcceptedProtocol(props.AcceptedProtocols)
			enabled = props.EnabledState == frontdoor.RoutingRuleEnabledStateEnabled

			forwardConfiguration, err := flattenRoutingRuleForwardingConfiguration(props.RouteConfiguration, oldBlocks)
			if err != nil {
				return nil, err
			}
			forwardingConfiguration = *forwardConfiguration

			frontendEndpoints, err := flattenArmFrontDoorFrontendEndpointsSubResources(props.FrontendEndpoints)
			if err != nil {
				return nil, err
			}
			frontEndEndpoints = *frontendEndpoints

			if props.PatternsToMatch != nil {
				patternsToMatch = *props.PatternsToMatch
			}
			redirectConfiguration = flattenRoutingRuleRedirectConfiguration(props.RouteConfiguration)
		}

		output = append(output, map[string]interface{}{
			"accepted_protocols":       acceptedProtocols,
			"enabled":                  enabled,
			"forwarding_configuration": forwardingConfiguration,
			"frontend_endpoints":       frontEndEndpoints,
			"id":                       id,
			"name":                     name,
			"patterns_to_match":        patternsToMatch,
			"redirect_configuration":   redirectConfiguration,
		})
	}

	return &output, nil
}

func flattenRoutingRuleForwardingConfiguration(config frontdoor.BasicRouteConfiguration, oldConfig interface{}) (*[]interface{}, error) {
	v, ok := config.(frontdoor.ForwardingConfiguration)
	if !ok {
		return &[]interface{}{}, nil
	}

	name := ""
	if v.BackendPool != nil && v.BackendPool.ID != nil {
		backendPoolId, err := parse.BackendPoolID(*v.BackendPool.ID)
		if err != nil {
			return nil, err
		}

		name = backendPoolId.Name
	}

	customForwardingPath := ""
	if v.CustomForwardingPath != nil {
		customForwardingPath = *v.CustomForwardingPath
	}

	cacheEnabled := false
	cacheQueryParameterStripDirective := string(frontdoor.StripAll)
	cacheUseDynamicCompression := false

	if cacheConfiguration := v.CacheConfiguration; cacheConfiguration != nil {
		cacheEnabled = true
		if stripDirective := cacheConfiguration.QueryParameterStripDirective; stripDirective != "" {
			cacheQueryParameterStripDirective = string(stripDirective)
		}

		if dynamicCompression := cacheConfiguration.DynamicCompression; dynamicCompression != "" {
			cacheUseDynamicCompression = string(dynamicCompression) == string(frontdoor.DynamicCompressionEnabledEnabled)
		}
	} else {
		// if the cache is disabled, use the default values or revert to what they were in the previous plan
		old, ok := oldConfig.([]interface{})
		if ok {
			for _, oldValue := range old {
				oldVal, ok := oldValue.(map[string]interface{})
				if ok {
					thisName := oldVal["name"].(string)
					if name == thisName {
						oldConfigs := oldVal["forwarding_configuration"].([]interface{})
						if len(oldConfigs) > 0 {
							ofc := oldConfigs[0].(map[string]interface{})
							cacheQueryParameterStripDirective = ofc["cache_query_parameter_strip_directive"].(string)
							cacheUseDynamicCompression = ofc["cache_use_dynamic_compression"].(bool)
						}
					}
				}
			}
		}
	}

	return &[]interface{}{
		map[string]interface{}{
			"backend_pool_name":                     "name",
			"custom_forwarding_path":                customForwardingPath,
			"forwarding_protocol":                   string(v.ForwardingProtocol),
			"cache_enabled":                         cacheEnabled,
			"cache_query_parameter_strip_directive": cacheQueryParameterStripDirective,
			"cache_use_dynamic_compression":         cacheUseDynamicCompression,
		},
	}, nil
}

func flattenRoutingRuleRedirectConfiguration(config frontdoor.BasicRouteConfiguration) []interface{} {
	v, ok := config.(frontdoor.RedirectConfiguration)
	if !ok {
		return []interface{}{}
	}

	customFragment := ""
	if v.CustomFragment != nil {
		customFragment = *v.CustomFragment
	}

	customHost := ""
	if v.CustomHost != nil {
		customHost = *v.CustomHost
	}

	customQueryString := ""
	if v.CustomQueryString != nil {
		customQueryString = *v.CustomQueryString
	}

	customPath := ""
	if v.CustomPath != nil {
		customPath = *v.CustomPath
	}

	return []interface{}{
		map[string]interface{}{
			"custom_host":         customHost,
			"custom_fragment":     customFragment,
			"custom_query_string": customQueryString,
			"custom_path":         customPath,
			"redirect_protocol":   string(v.RedirectProtocol),
			"redirect_type":       string(v.RedirectType),
		},
	}
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

func flattenArmFrontDoorFrontendEndpointsSubResources(input *[]frontdoor.SubResource) (*[]string, error) {
	output := make([]string, 0)

	if input == nil {
		return &output, nil
	}

	for _, v := range *input {
		if v.ID == nil {
			continue
		}

		id, err := parse.FrontendEndpointID(*v.ID)
		if err != nil {
			return nil, err
		}

		output = append(output, id.Name)
	}

	return &output, nil
}
