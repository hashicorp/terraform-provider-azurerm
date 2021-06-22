package frontdoor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceFrontDoor() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontDoorCreateUpdate,
		Read:   resourceFrontDoorRead,
		Update: resourceFrontDoorCreateUpdate,
		Delete: resourceFrontDoorDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontDoorID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.FrontDoorUpgradeV0ToV1{},
			1: migration.FrontDoorUpgradeV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(6 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(6 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"cname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"header_frontdoor_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"friendly_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"load_balancer_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			// TODO: In 3.0
			// Move 'enforce_backend_pools_certificate_name_check' and 'backend_pools_send_receive_timeout_seconds'
			// into a 'backend_pool_settings' block
			"enforce_backend_pools_certificate_name_check": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"backend_pools_send_receive_timeout_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntBetween(0, 240),
			},

			// TODO: Remove in 3.0
			"location": {
				Type:       pluginsdk.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Due to the service's API changing 'location' must now always be set to 'Global' for new resources, however if the Front Door service was created prior 2020/03/10 it may continue to exist in a specific current location",
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"routing_rule": {
				Type:     pluginsdk.TypeList,
				MaxItems: 500,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.BackendPoolRoutingRuleName,
						},
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
						"accepted_protocols": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 2,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(frontdoor.HTTP),
									string(frontdoor.HTTPS),
								}, false),
							},
						},
						"patterns_to_match": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 25,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						"frontend_endpoints": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 500,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						"redirect_configuration": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"custom_fragment": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
									"custom_host": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
									"custom_path": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
									"custom_query_string": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
									"redirect_protocol": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.RedirectProtocolHTTPOnly),
											string(frontdoor.RedirectProtocolHTTPSOnly),
											string(frontdoor.RedirectProtocolMatchRequest),
										}, false),
									},
									"redirect_type": {
										Type:     pluginsdk.TypeString,
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
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"backend_pool_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validate.BackendPoolRoutingRuleName,
									},
									"cache_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},
									"cache_use_dynamic_compression": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},
									"cache_query_parameter_strip_directive": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(frontdoor.StripAll),
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.StripAll),
											string(frontdoor.StripNone),
										}, false),
									},
									"custom_forwarding_path": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
									"forwarding_protocol": {
										Type:     pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeList,
				MaxItems: 5000,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.BackendPoolRoutingRuleName,
						},
						"sample_size": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  4,
						},
						"successful_samples_required": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  2,
						},
						"additional_latency_milliseconds": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},

			"backend_pool_health_probe": {
				Type:     pluginsdk.TypeList,
				MaxItems: 5000,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.BackendPoolRoutingRuleName,
						},
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
						"path": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "/",
						},
						"protocol": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(frontdoor.HTTP),
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.HTTP),
								string(frontdoor.HTTPS),
							}, false),
						},
						"probe_method": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(frontdoor.GET),
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.GET),
								string(frontdoor.HEAD),
							}, false),
						},
						"interval_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  120,
						},
					},
				},
			},

			"backend_pool": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"backend": {
							Type:     pluginsdk.TypeList,
							MaxItems: 500,
							Required: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
									},
									"address": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
									"http_port": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"https_port": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"weight": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      50,
										ValidateFunc: validation.IntBetween(1, 1000),
									},
									"priority": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      1,
										ValidateFunc: validation.IntBetween(1, 5),
									},
									"host_header": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
								},
							},
						},
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.BackendPoolRoutingRuleName,
						},
						"health_probe_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"load_balancing_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"frontend_endpoint": {
				Type:     pluginsdk.TypeList,
				MaxItems: 500,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.BackendPoolRoutingRuleName,
						},
						"host_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"session_affinity_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"session_affinity_ttl_seconds": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  0,
						},
						"web_application_firewall_policy_link_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							// TODO: validation that this is a resource id
						},
					},
				},
			},

			// Computed values
			"explicit_resource_order": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"backend_pool_ids": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"frontend_endpoint_ids": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"routing_rule_ids": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"backend_pool_load_balancing_ids": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"backend_pool_health_probe_ids": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"backend_pool_health_probes": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"backend_pool_load_balancing_settings": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"backend_pools": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"frontend_endpoints": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"routing_rules": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(frontDoorCustomizeDiff),
	}
}

func resourceFrontDoorCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	frontDoorId := parse.NewFrontDoorID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for present of existing Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_frontdoor", frontDoorId.ID())
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

	friendlyName := d.Get("friendly_name").(string)
	routingRules := d.Get("routing_rule").([]interface{})
	loadBalancingSettings := d.Get("backend_pool_load_balancing").([]interface{})
	healthProbeSettings := d.Get("backend_pool_health_probe").([]interface{})
	backendPools := d.Get("backend_pool").([]interface{})
	frontendEndpoints := d.Get("frontend_endpoint").([]interface{})
	backendPoolsSettings := d.Get("enforce_backend_pools_certificate_name_check").(bool)
	backendPoolsSendReceiveTimeoutSeconds := int32(d.Get("backend_pools_send_receive_timeout_seconds").(int))
	enabledState := d.Get("load_balancer_enabled").(bool)
	explicitResourceOrder := d.Get("explicit_resource_order").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	// If the explicitResourceOrder is empty and it's not a new resource set the mapping table to the state file and return an error.
	// If the explicitResourceOrder is empty and it is a new resource it will run the CreateOrUpdate as expected
	// If the explicitResourceOrder is NOT empty and it is NOT a new resource it will run the CreateOrUpdate as expected
	if len(explicitResourceOrder) == 0 && !d.IsNewResource() {
		d.Set("explicit_resource_order", flattenExplicitResourceOrder(backendPools, frontendEndpoints, routingRules, loadBalancingSettings, healthProbeSettings, frontDoorId))
	} else {
		frontDoorParameters := frontdoor.FrontDoor{
			Location: utils.String(location),
			Properties: &frontdoor.Properties{
				FriendlyName:          utils.String(friendlyName),
				RoutingRules:          expandFrontDoorRoutingRule(routingRules, frontDoorId),
				BackendPools:          expandFrontDoorBackendPools(backendPools, frontDoorId),
				BackendPoolsSettings:  expandFrontDoorBackendPoolsSettings(backendPoolsSettings, backendPoolsSendReceiveTimeoutSeconds),
				FrontendEndpoints:     expandFrontDoorFrontendEndpoint(frontendEndpoints, frontDoorId),
				HealthProbeSettings:   expandFrontDoorHealthProbeSettingsModel(healthProbeSettings, frontDoorId),
				LoadBalancingSettings: expandFrontDoorLoadBalancingSettingsModel(loadBalancingSettings, frontDoorId),
				EnabledState:          expandFrontDoorEnabledState(enabledState),
			},
			Tags: tags.Expand(t),
		}

		future, err := client.CreateOrUpdate(ctx, frontDoorId.ResourceGroup, frontDoorId.Name, frontDoorParameters)
		if err != nil {
			return fmt.Errorf("creating Front Door %q (Resource Group %q): %+v", frontDoorId.Name, frontDoorId.ResourceGroup, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation of Front Door %q (Resource Group %q): %+v", frontDoorId.Name, frontDoorId.ResourceGroup, err)
		}

		d.SetId(frontDoorId.ID())
		d.Set("explicit_resource_order", flattenExplicitResourceOrder(backendPools, frontendEndpoints, routingRules, loadBalancingSettings, healthProbeSettings, frontDoorId))
	}

	return resourceFrontDoorRead(d, meta)
}

func resourceFrontDoorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorIDInsensitively(d.Id())
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
		explicitResourceOrder := d.Get("explicit_resource_order").([]interface{})
		flattenedBackendPools, err := flattenFrontDoorBackendPools(props.BackendPools, *id, explicitResourceOrder)
		if err != nil {
			return fmt.Errorf("flattening `backend_pool`: %+v", err)
		}
		if err := d.Set("backend_pool", flattenedBackendPools); err != nil {
			return fmt.Errorf("setting `backend_pool`: %+v", err)
		}

		backendPoolSettings := flattenFrontDoorBackendPoolsSettings(props.BackendPoolsSettings)

		d.Set("enforce_backend_pools_certificate_name_check", backendPoolSettings.enforceBackendPoolsCertificateNameCheck)
		d.Set("backend_pools_send_receive_timeout_seconds", backendPoolSettings.backendPoolsSendReceiveTimeoutSeconds)
		d.Set("cname", props.Cname)
		d.Set("header_frontdoor_id", props.FrontdoorID)
		d.Set("load_balancer_enabled", props.EnabledState == frontdoor.EnabledStateEnabled)
		d.Set("friendly_name", props.FriendlyName)

		// Need to call frontEndEndpointClient here to get the frontEndEndpoint information from that client
		// because the information is hidden from the main frontDoorClient "by design"...
		frontEndEndpointsClient := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
		frontEndEndpointInfo, err := retrieveFrontEndEndpointInformation(ctx, frontEndEndpointsClient, *id, props.FrontendEndpoints)
		if err != nil {
			return fmt.Errorf("retrieving FrontEnd Endpoint Information: %+v", err)
		}

		// Force the returned flattenFrontEndEndpoints into the order defined in the explicit_resource_order mapping table
		frontDoorFrontendEndpoints, err := flattenFrontEndEndpoints(frontEndEndpointInfo, *id, explicitResourceOrder)
		if err != nil {
			return fmt.Errorf("flattening `frontend_endpoint`: %+v", err)
		}
		if err := d.Set("frontend_endpoint", frontDoorFrontendEndpoints); err != nil {
			return fmt.Errorf("setting `frontend_endpoint`: %+v", err)
		}

		// Force the returned flattenFrontDoorHealthProbeSettingsModel into the order defined in the explicit_resource_order mapping table
		if err := d.Set("backend_pool_health_probe", flattenFrontDoorHealthProbeSettingsModel(props.HealthProbeSettings, *id, explicitResourceOrder)); err != nil {
			return fmt.Errorf("setting `backend_pool_health_probe`: %+v", err)
		}

		// Force the returned flattenFrontDoorLoadBalancingSettingsModel into the order defined in the explicit_resource_order mapping table
		if err := d.Set("backend_pool_load_balancing", flattenFrontDoorLoadBalancingSettingsModel(props.LoadBalancingSettings, *id, explicitResourceOrder)); err != nil {
			return fmt.Errorf("setting `backend_pool_load_balancing`: %+v", err)
		}

		var flattenedRoutingRules *[]interface{}
		// Force the returned flattenedRoutingRules into the order defined in the explicit_resource_order mapping table
		flattenedRoutingRules, err = flattenFrontDoorRoutingRule(props.RoutingRules, d.Get("routing_rule"), *id, explicitResourceOrder)
		if err != nil {
			return fmt.Errorf("flattening `routing_rules`: %+v", err)
		}
		if err := d.Set("routing_rule", flattenedRoutingRules); err != nil {
			return fmt.Errorf("setting `routing_rules`: %+v", err)
		}

		// Populate computed values
		bpHealthProbeSettings := make(map[string]string)
		if props.HealthProbeSettings != nil {
			for _, v := range *props.HealthProbeSettings {
				if v.Name == nil || v.ID == nil {
					continue
				}
				rid, err := parse.HealthProbeIDInsensitively(*v.ID)
				if err != nil {
					continue
				}
				bpHealthProbeSettings[*v.Name] = rid.ID()
			}
		}
		if err := d.Set("backend_pool_health_probes", bpHealthProbeSettings); err != nil {
			return fmt.Errorf("setting `backend_pool_health_probes`: %+v", err)
		}

		bpLBSettings := make(map[string]string)
		if props.LoadBalancingSettings != nil {
			for _, v := range *props.LoadBalancingSettings {
				if v.Name == nil || v.ID == nil {
					continue
				}
				rid, err := parse.LoadBalancingIDInsensitively(*v.ID)
				if err != nil {
					continue
				}
				bpLBSettings[*v.Name] = rid.ID()
			}
		}
		if err := d.Set("backend_pool_load_balancing_settings", bpLBSettings); err != nil {
			return fmt.Errorf("setting `backend_pool_load_balancing_settings`: %+v", err)
		}

		backendPools := make(map[string]string)
		if props.BackendPools != nil {
			for _, v := range *props.BackendPools {
				if v.Name == nil || v.ID == nil {
					continue
				}
				rid, err := parse.BackendPoolIDInsensitively(*v.ID)
				if err != nil {
					continue
				}
				backendPools[*v.Name] = rid.ID()
			}
		}
		if err := d.Set("backend_pools", backendPools); err != nil {
			return fmt.Errorf("setting `backend_pools`: %+v", err)
		}

		frontendEndpoints := make(map[string]string)
		if props.FrontendEndpoints != nil {
			for _, v := range *props.FrontendEndpoints {
				if v.Name == nil || v.ID == nil {
					continue
				}
				rid, err := parse.FrontendEndpointIDInsensitively(*v.ID)
				if err != nil {
					continue
				}
				frontendEndpoints[*v.Name] = rid.ID()
			}
		}
		if err := d.Set("frontend_endpoints", frontendEndpoints); err != nil {
			return fmt.Errorf("setting `frontend_endpoints`: %+v", err)
		}

		routingRules := make(map[string]string)
		if props.RoutingRules != nil {
			for _, v := range *props.RoutingRules {
				if v.Name == nil || v.ID == nil {
					continue
				}
				rid, err := parse.RoutingRuleIDInsensitively(*v.ID)
				if err != nil {
					continue
				}
				routingRules[*v.Name] = rid.ID()
			}
		}
		if err := d.Set("routing_rules", routingRules); err != nil {
			return fmt.Errorf("setting `routing_rules`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceFrontDoorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if future.Response() != nil {
			if response.WasNotFound(future.Response()) {
				return nil
			}
		}
		return fmt.Errorf("deleting Front Door %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if future.Response() != nil {
			if !response.WasNotFound(future.Response()) {
				return fmt.Errorf("waiting for deleting Front Door %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}
		}
	}

	return nil
}

func expandFrontDoorBackendPools(input []interface{}, frontDoorId parse.FrontDoorId) *[]frontdoor.BackendPool {
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

		backendPoolId := parse.NewBackendPoolID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, backendPoolName).ID()
		healthProbeId := parse.NewHealthProbeID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, backendPoolHealthProbeName).ID()
		loadBalancingId := parse.NewLoadBalancingID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, backendPoolLoadBalancingName).ID()

		result := frontdoor.BackendPool{
			ID:   utils.String(backendPoolId),
			Name: utils.String(backendPoolName),
			BackendPoolProperties: &frontdoor.BackendPoolProperties{
				Backends: expandFrontDoorBackend(backends),
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

func expandFrontDoorBackend(input []interface{}) *[]frontdoor.Backend {
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
			EnabledState:      expandFrontDoorBackendEnabledState(enabled),
			HTTPPort:          utils.Int32(httpPort),
			HTTPSPort:         utils.Int32(httpsPort),
			Priority:          utils.Int32(priority),
			Weight:            utils.Int32(weight),
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorBackendEnabledState(isEnabled bool) frontdoor.BackendEnabledState {
	if isEnabled {
		return frontdoor.Enabled
	}
	return frontdoor.Disabled
}

func expandFrontDoorBackendPoolsSettings(enforceCertificateNameCheck bool, backendPoolsSendReceiveTimeoutSeconds int32) *frontdoor.BackendPoolsSettings {
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

func expandFrontDoorFrontendEndpoint(input []interface{}, frontDoorId parse.FrontDoorId) *[]frontdoor.FrontendEndpoint {
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
		id := parse.NewFrontendEndpointID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, name).ID()
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

func expandFrontDoorHealthProbeSettingsModel(input []interface{}, frontDoorId parse.FrontDoorId) *[]frontdoor.HealthProbeSettingsModel {
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
		healthProbeId := parse.NewHealthProbeID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, name).ID()

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

func expandFrontDoorLoadBalancingSettingsModel(input []interface{}, frontDoorId parse.FrontDoorId) *[]frontdoor.LoadBalancingSettingsModel {
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
		loadBalancingId := parse.NewLoadBalancingID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, name).ID()

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

func expandFrontDoorRoutingRule(input []interface{}, frontDoorId parse.FrontDoorId) *[]frontdoor.RoutingRule {
	if len(input) == 0 {
		return nil
	}

	output := make([]frontdoor.RoutingRule, 0)

	for _, rr := range input {
		routingRule := rr.(map[string]interface{})
		name := routingRule["name"].(string)
		frontendEndpoints := routingRule["frontend_endpoints"].([]interface{})
		acceptedProtocols := routingRule["accepted_protocols"].([]interface{})
		ptm := routingRule["patterns_to_match"].([]interface{})
		enabled := routingRule["enabled"].(bool)

		patternsToMatch := make([]string, 0)
		for _, p := range ptm {
			patternsToMatch = append(patternsToMatch, p.(string))
		}

		var routingConfiguration frontdoor.BasicRouteConfiguration
		if rc := routingRule["redirect_configuration"].([]interface{}); len(rc) != 0 {
			routingConfiguration = expandFrontDoorRedirectConfiguration(rc)
		} else if fc := routingRule["forwarding_configuration"].([]interface{}); len(fc) != 0 {
			routingConfiguration = expandFrontDoorForwardingConfiguration(fc, frontDoorId)
		}
		routingRuleId := parse.NewRoutingRuleID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, name).ID()

		currentRoutingRule := frontdoor.RoutingRule{
			ID:   utils.String(routingRuleId),
			Name: utils.String(name),
			RoutingRuleProperties: &frontdoor.RoutingRuleProperties{
				FrontendEndpoints:  expandFrontDoorFrontEndEndpoints(frontendEndpoints, frontDoorId),
				AcceptedProtocols:  expandFrontDoorAcceptedProtocols(acceptedProtocols),
				PatternsToMatch:    &patternsToMatch,
				EnabledState:       frontdoor.RoutingRuleEnabledState(expandFrontDoorEnabledState(enabled)),
				RouteConfiguration: routingConfiguration,
			},
		}
		output = append(output, currentRoutingRule)
	}

	return &output
}

func expandFrontDoorAcceptedProtocols(input []interface{}) *[]frontdoor.Protocol {
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

func expandFrontDoorFrontEndEndpoints(input []interface{}, frontDoorId parse.FrontDoorId) *[]frontdoor.SubResource {
	if len(input) == 0 {
		return &[]frontdoor.SubResource{}
	}

	output := make([]frontdoor.SubResource, 0)

	for _, name := range input {
		frontendEndpointId := parse.NewFrontendEndpointID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, name.(string)).ID()
		result := frontdoor.SubResource{
			ID: utils.String(frontendEndpointId),
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorEnabledState(enabled bool) frontdoor.EnabledState {
	if enabled {
		return frontdoor.EnabledStateEnabled
	}
	return frontdoor.EnabledStateDisabled
}

func expandFrontDoorRedirectConfiguration(input []interface{}) frontdoor.RedirectConfiguration {
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

func expandFrontDoorForwardingConfiguration(input []interface{}, frontDoorId parse.FrontDoorId) frontdoor.ForwardingConfiguration {
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

	backendPoolId := parse.NewBackendPoolID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, backendPoolName).ID()
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

func flattenExplicitResourceOrder(backendPools, frontendEndpoints, routingRules, loadBalancingSettings, healthProbeSettings []interface{}, frontDoorId parse.FrontDoorId) *[]interface{} {
	output := make([]interface{}, 0)
	var backendPoolOrder []string
	var frontedEndpointOrder []string
	var routingRulesOrder []string
	var backendPoolLoadBalancingOrder []string
	var backendPoolHealthProbeOrder []string
	if len(backendPools) > 0 {
		flattenendBackendPools, err := flattenFrontDoorBackendPools(expandFrontDoorBackendPools(backendPools, frontDoorId), frontDoorId, make([]interface{}, 0))
		if err == nil {
			for _, ids := range *flattenendBackendPools {
				backendPool := ids.(map[string]interface{})
				backendPoolOrder = append(backendPoolOrder, backendPool["id"].(string))
			}
		}
	}
	if len(frontendEndpoints) > 0 {
		flattenendfrontendEndpoints, err := flattenFrontEndEndpoints(expandFrontDoorFrontendEndpoint(frontendEndpoints, frontDoorId), frontDoorId, make([]interface{}, 0))
		if err == nil {
			for _, ids := range *flattenendfrontendEndpoints {
				frontendEndPoint := ids.(map[string]interface{})
				frontedEndpointOrder = append(frontedEndpointOrder, frontendEndPoint["id"].(string))
			}
		}
	}
	if len(routingRules) > 0 {
		var oldBlocks interface{}
		flattenendRoutingRules, err := flattenFrontDoorRoutingRule(expandFrontDoorRoutingRule(routingRules, frontDoorId), oldBlocks, frontDoorId, make([]interface{}, 0))
		if err == nil {
			for _, ids := range *flattenendRoutingRules {
				routingRule := ids.(map[string]interface{})
				routingRulesOrder = append(routingRulesOrder, routingRule["id"].(string))
			}
		}
	}
	if len(loadBalancingSettings) > 0 {
		flattenendLoadBalancingSettings := flattenFrontDoorLoadBalancingSettingsModel(expandFrontDoorLoadBalancingSettingsModel(loadBalancingSettings, frontDoorId), frontDoorId, make([]interface{}, 0))

		if len(flattenendLoadBalancingSettings) > 0 {
			for _, ids := range flattenendLoadBalancingSettings {
				loadBalancingSetting := ids.(map[string]interface{})
				backendPoolLoadBalancingOrder = append(backendPoolLoadBalancingOrder, loadBalancingSetting["id"].(string))
			}
		}
	}
	if len(healthProbeSettings) > 0 {
		flattenendHealthProbeSettings := flattenFrontDoorHealthProbeSettingsModel(expandFrontDoorHealthProbeSettingsModel(healthProbeSettings, frontDoorId), frontDoorId, make([]interface{}, 0))

		if len(flattenendHealthProbeSettings) > 0 {
			for _, ids := range flattenendHealthProbeSettings {
				healthProbeSetting := ids.(map[string]interface{})
				backendPoolHealthProbeOrder = append(backendPoolHealthProbeOrder, healthProbeSetting["id"].(string))
			}
		}
	}

	output = append(output, map[string]interface{}{
		"backend_pool_ids":                backendPoolOrder,
		"frontend_endpoint_ids":           frontedEndpointOrder,
		"routing_rule_ids":                routingRulesOrder,
		"backend_pool_load_balancing_ids": backendPoolLoadBalancingOrder,
		"backend_pool_health_probe_ids":   backendPoolHealthProbeOrder,
	})

	return &output
}

func combineBackendPools(allPools []frontdoor.BackendPool, orderedIds []interface{}, frontDoorId parse.FrontDoorId) ([]interface{}, error) {
	output := make([]interface{}, 0)
	found := false

	// first find all of the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, backend := range allPools {
			if strings.EqualFold(v.(string), *backend.ID) {
				orderedBackendPool, err := flattenSingleFrontDoorBackendPools(&backend, frontDoorId)
				if err == nil {
					output = append(output, orderedBackendPool)
					break
				} else {
					return nil, err
				}
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, backend := range allPools {
		found = false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *backend.ID) {
				found = true
				break
			}
		}

		if !found {
			newBackendPool, err := flattenSingleFrontDoorBackendPools(&backend, frontDoorId)
			if err == nil {
				output = append(output, newBackendPool)
			} else {
				return nil, err
			}
		}
	}

	return output, nil
}

func flattenFrontDoorBackendPools(input *[]frontdoor.BackendPool, frontDoorId parse.FrontDoorId, explicitOrder []interface{}) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	output := make([]interface{}, 0)

	if len(explicitOrder) > 0 {
		orderedBackendPools := explicitOrder[0].(map[string]interface{})
		orderedBackendPoolsIds := orderedBackendPools["backend_pool_ids"].([]interface{})
		combinedBackendPools, err := combineBackendPools(*input, orderedBackendPoolsIds, frontDoorId)
		if err == nil {
			output = combinedBackendPools
		} else {
			return nil, err
		}
	} else {
		for _, backend := range *input {
			backendPool, err := flattenSingleFrontDoorBackendPools(&backend, frontDoorId)
			if err == nil {
				output = append(output, backendPool)
			} else {
				return nil, err
			}
		}
	}

	return &output, nil
}

func flattenSingleFrontDoorBackendPools(input *frontdoor.BackendPool, frontDoorId parse.FrontDoorId) (map[string]interface{}, error) {
	if input == nil {
		return make(map[string]interface{}), nil
	}

	id := ""
	name := ""
	if input.Name != nil {
		name = *input.Name
		// rewrite the ID to ensure it's consistent
		id = parse.NewBackendPoolID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, name).ID()
	}

	backend := make([]interface{}, 0)
	healthProbeName := ""
	loadBalancingName := ""
	if props := input.BackendPoolProperties; props != nil {
		backend = flattenFrontDoorBackend(props.Backends)
		if props.HealthProbeSettings != nil && props.HealthProbeSettings.ID != nil {
			name, err := parse.HealthProbeIDInsensitively(*props.HealthProbeSettings.ID)
			if err != nil {
				return nil, err
			}
			healthProbeName = name.HealthProbeSettingName
		}

		if props.LoadBalancingSettings != nil && props.LoadBalancingSettings.ID != nil {
			name, err := parse.LoadBalancingIDInsensitively(*props.LoadBalancingSettings.ID)
			if err != nil {
				return nil, err
			}
			loadBalancingName = name.LoadBalancingSettingName
		}
	}

	output := map[string]interface{}{
		"backend":             backend,
		"health_probe_name":   healthProbeName,
		"id":                  id,
		"load_balancing_name": loadBalancingName,
		"name":                name,
	}

	return output, nil
}

type flattenedBackendPoolSettings struct {
	enforceBackendPoolsCertificateNameCheck bool
	backendPoolsSendReceiveTimeoutSeconds   int
}

func flattenFrontDoorBackendPoolsSettings(input *frontdoor.BackendPoolsSettings) flattenedBackendPoolSettings {
	if input == nil {
		return flattenedBackendPoolSettings{
			enforceBackendPoolsCertificateNameCheck: true,
			backendPoolsSendReceiveTimeoutSeconds:   60,
		}
	}

	enforceCertificateNameCheck := false
	sendReceiveTimeoutSeconds := 0
	if input.EnforceCertificateNameCheck != "" && input.EnforceCertificateNameCheck == frontdoor.EnforceCertificateNameCheckEnabledStateEnabled {
		enforceCertificateNameCheck = true
	}
	if input.SendRecvTimeoutSeconds != nil {
		sendReceiveTimeoutSeconds = int(*input.SendRecvTimeoutSeconds)
	}

	return flattenedBackendPoolSettings{
		enforceBackendPoolsCertificateNameCheck: enforceCertificateNameCheck,
		backendPoolsSendReceiveTimeoutSeconds:   sendReceiveTimeoutSeconds,
	}
}

func flattenFrontDoorBackend(input *[]frontdoor.Backend) []interface{} {
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

func combineFrontEndEndpoints(allEndpoints []frontdoor.FrontendEndpoint, orderedIds []interface{}, frontDoorId parse.FrontDoorId) ([]interface{}, error) {
	output := make([]interface{}, 0)
	found := false

	// first find all of the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, frontendEndpoint := range allEndpoints {
			if strings.EqualFold(v.(string), *frontendEndpoint.ID) {
				orderedFrontendEndpoint, err := flattenSingleFrontEndEndpoints(frontendEndpoint, frontDoorId)
				if err == nil {
					output = append(output, orderedFrontendEndpoint)
					break
				} else {
					return nil, err
				}
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, frontendEndpoint := range allEndpoints {
		found = false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *frontendEndpoint.ID) {
				found = true
				break
			}
		}

		if !found {
			newFrontendEndpoint, err := flattenSingleFrontEndEndpoints(frontendEndpoint, frontDoorId)
			if err == nil {
				output = append(output, newFrontendEndpoint)
			} else {
				return nil, err
			}
		}
	}

	return output, nil
}

func flattenFrontEndEndpoints(input *[]frontdoor.FrontendEndpoint, frontDoorId parse.FrontDoorId, explicitOrder []interface{}) (*[]interface{}, error) {
	output := make([]interface{}, 0)
	if input == nil {
		return &output, nil
	}

	if len(explicitOrder) > 0 {
		orderedFrontEnd := explicitOrder[0].(map[string]interface{})
		orderedFrontEndIds := orderedFrontEnd["frontend_endpoint_ids"].([]interface{})
		combinedFrontEndEndpoints, err := combineFrontEndEndpoints(*input, orderedFrontEndIds, frontDoorId)
		if err == nil {
			output = combinedFrontEndEndpoints
		} else {
			return nil, err
		}
	} else {
		for _, v := range *input {
			frontendEndpoint, err := flattenSingleFrontEndEndpoints(v, frontDoorId)
			if err == nil {
				output = append(output, frontendEndpoint)
			} else {
				return nil, err
			}
		}
	}

	return &output, nil
}

func flattenSingleFrontEndEndpoints(input frontdoor.FrontendEndpoint, frontDoorId parse.FrontDoorId) (map[string]interface{}, error) {
	id := ""
	name := ""
	if input.Name != nil {
		// rewrite the ID to ensure it's consistent
		id = parse.NewFrontendEndpointID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, *input.Name).ID()
		name = *input.Name
	}
	// TODO: I may have to include the customHTTPSConfiguration as returned from the frontendEndpoint due to an issue in
	// portal. Still investigating this.
	// customHTTPSConfiguration := make([]interface{}, 0)
	// customHttpsProvisioningEnabled := false
	hostName := ""
	sessionAffinityEnabled := false
	sessionAffinityTlsSeconds := 0
	webApplicationFirewallPolicyLinkId := ""
	if props := input.FrontendEndpointProperties; props != nil {
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
			// rewrite the ID to ensure it's consistent
			parsed, err := parse.WebApplicationFirewallPolicyIDInsensitively(*waf.ID)
			if err != nil {
				return nil, err
			}
			webApplicationFirewallPolicyLinkId = parsed.ID()
		}
		// flattenedHttpsConfig := flattenCustomHttpsConfiguration(props)
		// customHTTPSConfiguration = flattenedHttpsConfig.CustomHTTPSConfiguration
		// customHttpsProvisioningEnabled = flattenedHttpsConfig.CustomHTTPSProvisioningEnabled
	}

	output := map[string]interface{}{
		// "custom_https_configuration":        customHTTPSConfiguration,
		// "custom_https_provisioning_enabled": customHttpsProvisioningEnabled,
		"host_name":                    hostName,
		"id":                           id,
		"name":                         name,
		"session_affinity_enabled":     sessionAffinityEnabled,
		"session_affinity_ttl_seconds": sessionAffinityTlsSeconds,
		"web_application_firewall_policy_link_id": webApplicationFirewallPolicyLinkId,
	}

	return output, nil
}

func combineHealthProbeSettingsModel(allHealthProbeSettings []frontdoor.HealthProbeSettingsModel, orderedIds []interface{}, frontDoorId parse.FrontDoorId) []interface{} {
	output := make([]interface{}, 0)
	found := false

	// first find all of the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, healthProbeSetting := range allHealthProbeSettings {
			if strings.EqualFold(v.(string), *healthProbeSetting.ID) {
				orderedHealthProbeSetting := flattenSingleFrontDoorHealthProbeSettingsModel(&healthProbeSetting, frontDoorId)
				output = append(output, orderedHealthProbeSetting)
				break
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, healthProbeSetting := range allHealthProbeSettings {
		found = false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *healthProbeSetting.ID) {
				found = true
				break
			}
		}

		if !found {
			newHealthProbeSetting := flattenSingleFrontDoorHealthProbeSettingsModel(&healthProbeSetting, frontDoorId)
			output = append(output, newHealthProbeSetting)
		}
	}

	return output
}

func flattenFrontDoorHealthProbeSettingsModel(input *[]frontdoor.HealthProbeSettingsModel, frontDoorId parse.FrontDoorId, explicitOrder []interface{}) []interface{} {
	output := make([]interface{}, 0)
	if input == nil {
		return output
	}

	if len(explicitOrder) > 0 {
		orderedHealthProbeSetting := explicitOrder[0].(map[string]interface{})
		orderedHealthProbeSettingIds := orderedHealthProbeSetting["backend_pool_health_probe_ids"].([]interface{})
		output = combineHealthProbeSettingsModel(*input, orderedHealthProbeSettingIds, frontDoorId)
	} else {
		for _, v := range *input {
			healthProbeSetting := flattenSingleFrontDoorHealthProbeSettingsModel(&v, frontDoorId)
			output = append(output, healthProbeSetting)
		}
	}

	return output
}

func flattenSingleFrontDoorHealthProbeSettingsModel(input *frontdoor.HealthProbeSettingsModel, frontDoorId parse.FrontDoorId) map[string]interface{} {
	if input == nil {
		return make(map[string]interface{})
	}

	id := ""
	name := ""
	if input.Name != nil {
		name = *input.Name
		// rewrite the ID to ensure it's consistent
		id = parse.NewHealthProbeID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, name).ID()
	}

	enabled := false
	intervalInSeconds := 0
	path := ""
	probeMethod := ""
	protocol := ""

	if properties := input.HealthProbeSettingsProperties; properties != nil {
		if properties.IntervalInSeconds != nil {
			intervalInSeconds = int(*properties.IntervalInSeconds)
		}
		if properties.Path != nil {
			path = *properties.Path
		}
		if healthProbeMethod := properties.HealthProbeMethod; healthProbeMethod != "" {
			// I have to upper this as the frontdoor.GET and frontdoor.HEAD types are uppercased
			// but Azure stores them in the resource as sentence cased (e.g. "Get" and "Head")
			probeMethod = strings.ToUpper(string(healthProbeMethod))
		}
		if properties.EnabledState != "" {
			enabled = properties.EnabledState == frontdoor.HealthProbeEnabledEnabled
		}
		protocol = string(properties.Protocol)
	}

	output := map[string]interface{}{
		"enabled":             enabled,
		"id":                  id,
		"name":                name,
		"protocol":            protocol,
		"interval_in_seconds": intervalInSeconds,
		"path":                path,
		"probe_method":        probeMethod,
	}

	return output
}

func combineLoadBalancingSettingsModel(allLoadBalancingSettings []frontdoor.LoadBalancingSettingsModel, orderedIds []interface{}, frontDoorId parse.FrontDoorId) []interface{} {
	output := make([]interface{}, 0)
	found := false

	// first find all of the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, loadBalancingSetting := range allLoadBalancingSettings {
			if strings.EqualFold(v.(string), *loadBalancingSetting.ID) {
				orderedLoadBalanceSetting := flattenSingleFrontDoorLoadBalancingSettingsModel(&loadBalancingSetting, frontDoorId)
				output = append(output, orderedLoadBalanceSetting)
				break
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, loadBalanceSetting := range allLoadBalancingSettings {
		found = false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *loadBalanceSetting.ID) {
				found = true
				break
			}
		}

		if !found {
			newLoadBalanceSetting := flattenSingleFrontDoorLoadBalancingSettingsModel(&loadBalanceSetting, frontDoorId)
			output = append(output, newLoadBalanceSetting)
		}
	}

	return output
}

func flattenFrontDoorLoadBalancingSettingsModel(input *[]frontdoor.LoadBalancingSettingsModel, frontDoorId parse.FrontDoorId, explicitOrder []interface{}) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	output := make([]interface{}, 0)

	if len(explicitOrder) > 0 {
		orderedLoadBalancingSettings := explicitOrder[0].(map[string]interface{})
		orderedLoadBalancingIds := orderedLoadBalancingSettings["backend_pool_load_balancing_ids"].([]interface{})
		output = combineLoadBalancingSettingsModel(*input, orderedLoadBalancingIds, frontDoorId)
	} else {
		for _, v := range *input {
			loadBalanceSetting := flattenSingleFrontDoorLoadBalancingSettingsModel(&v, frontDoorId)
			output = append(output, loadBalanceSetting)
		}
	}

	return output
}

func flattenSingleFrontDoorLoadBalancingSettingsModel(input *frontdoor.LoadBalancingSettingsModel, frontDoorId parse.FrontDoorId) map[string]interface{} {
	if input == nil {
		return make(map[string]interface{})
	}

	id := ""
	name := ""
	if input.Name != nil {
		name = *input.Name
		// rewrite the ID to ensure it's consistent
		id = parse.NewLoadBalancingID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, name).ID()
	}

	additionalLatencyMilliseconds := 0
	sampleSize := 0
	successfulSamplesRequired := 0
	if properties := input.LoadBalancingSettingsProperties; properties != nil {
		if properties.AdditionalLatencyMilliseconds != nil {
			additionalLatencyMilliseconds = int(*properties.AdditionalLatencyMilliseconds)
		}
		if properties.SampleSize != nil {
			sampleSize = int(*properties.SampleSize)
		}
		if properties.SuccessfulSamplesRequired != nil {
			successfulSamplesRequired = int(*properties.SuccessfulSamplesRequired)
		}
	}

	output := map[string]interface{}{
		"additional_latency_milliseconds": additionalLatencyMilliseconds,
		"id":                              id,
		"name":                            name,
		"sample_size":                     sampleSize,
		"successful_samples_required":     successfulSamplesRequired,
	}

	return output
}

func combineRoutingRules(allRoutingRules []frontdoor.RoutingRule, oldBlocks interface{}, orderedIds []interface{}, frontDoorId parse.FrontDoorId) ([]interface{}, error) {
	output := make([]interface{}, 0)
	found := false

	// first find all of the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, routingRule := range allRoutingRules {
			if strings.EqualFold(v.(string), *routingRule.ID) {
				orderedRoutingRule, err := flattenSingleFrontDoorRoutingRule(routingRule, oldBlocks, frontDoorId)
				if err == nil {
					output = append(output, orderedRoutingRule)
					break
				} else {
					return nil, err
				}
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, routingRule := range allRoutingRules {
		found = false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *routingRule.ID) {
				found = true
				break
			}
		}

		if !found {
			newRoutingRule, err := flattenSingleFrontDoorRoutingRule(routingRule, oldBlocks, frontDoorId)
			if err == nil {
				output = append(output, newRoutingRule)
			} else {
				return nil, err
			}
		}
	}

	return output, nil
}

func flattenFrontDoorRoutingRule(input *[]frontdoor.RoutingRule, oldBlocks interface{}, frontDoorId parse.FrontDoorId, explicitOrder []interface{}) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	output := make([]interface{}, 0)

	if len(explicitOrder) > 0 {
		orderedRule := explicitOrder[0].(map[string]interface{})
		orderedRountingRuleIds := orderedRule["routing_rule_ids"].([]interface{})
		combinedRoutingRules, err := combineRoutingRules(*input, oldBlocks, orderedRountingRuleIds, frontDoorId)
		if err != nil {
			return nil, err
		}
		output = combinedRoutingRules
	} else {
		for _, v := range *input {
			routingRule, err := flattenSingleFrontDoorRoutingRule(v, oldBlocks, frontDoorId)
			if err == nil {
				output = append(output, routingRule)
			} else {
				return nil, err
			}
		}
	}

	return &output, nil
}

func flattenSingleFrontDoorRoutingRule(input frontdoor.RoutingRule, oldBlocks interface{}, frontDoorId parse.FrontDoorId) (map[string]interface{}, error) {
	id := ""
	name := ""
	if input.Name != nil {
		// rewrite the ID to ensure it's consistent
		id = parse.NewRoutingRuleID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroup, frontDoorId.Name, *input.Name).ID()
		name = *input.Name
	}

	acceptedProtocols := make([]string, 0)
	enabled := false
	forwardingConfiguration := make([]interface{}, 0)
	frontEndEndpoints := make([]string, 0)
	patternsToMatch := make([]string, 0)
	redirectConfiguration := make([]interface{}, 0)

	if props := input.RoutingRuleProperties; props != nil {
		acceptedProtocols = flattenFrontDoorAcceptedProtocol(props.AcceptedProtocols)
		enabled = props.EnabledState == frontdoor.RoutingRuleEnabledStateEnabled

		forwardConfiguration, err := flattenRoutingRuleForwardingConfiguration(props.RouteConfiguration, oldBlocks)
		if err != nil {
			return nil, fmt.Errorf("flattening `forward_configuration`: %+v", err)
		}

		forwardingConfiguration = *forwardConfiguration
		frontendEndpoints, err := flattenFrontDoorFrontendEndpointsSubResources(props.FrontendEndpoints)
		if err != nil {
			return nil, fmt.Errorf("flattening `frontend_endpoints`: %+v", err)
		}

		frontEndEndpoints = *frontendEndpoints
		if props.PatternsToMatch != nil {
			patternsToMatch = *props.PatternsToMatch
		}
		redirectConfiguration = flattenRoutingRuleRedirectConfiguration(props.RouteConfiguration)
	}

	output := map[string]interface{}{
		"accepted_protocols":       acceptedProtocols,
		"enabled":                  enabled,
		"forwarding_configuration": forwardingConfiguration,
		"frontend_endpoints":       frontEndEndpoints,
		"id":                       id,
		"name":                     name,
		"patterns_to_match":        patternsToMatch,
		"redirect_configuration":   redirectConfiguration,
	}

	return output, nil
}

func flattenRoutingRuleForwardingConfiguration(config frontdoor.BasicRouteConfiguration, oldConfig interface{}) (*[]interface{}, error) {
	v, ok := config.(frontdoor.ForwardingConfiguration)
	if !ok {
		return &[]interface{}{}, nil
	}

	name := ""
	if v.BackendPool != nil && v.BackendPool.ID != nil {
		backendPoolId, err := parse.BackendPoolIDInsensitively(*v.BackendPool.ID)
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
			"backend_pool_name":                     name,
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

func flattenFrontDoorAcceptedProtocol(input *[]frontdoor.Protocol) []string {
	if input == nil {
		return make([]string, 0)
	}

	output := make([]string, 0)

	for _, p := range *input {
		output = append(output, string(p))
	}

	return output
}

func flattenFrontDoorFrontendEndpointsSubResources(input *[]frontdoor.SubResource) (*[]string, error) {
	output := make([]string, 0)
	if input == nil {
		return &output, nil
	}

	for _, v := range *input {
		if v.ID == nil {
			continue
		}

		id, err := parse.FrontendEndpointIDInsensitively(*v.ID)
		if err != nil {
			return nil, err
		}
		output = append(output, id.Name)
	}

	return &output, nil
}
