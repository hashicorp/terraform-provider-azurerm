package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApplicationGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApplicationGatewayCreateUpdate,
		Read:   resourceArmApplicationGatewayRead,
		Update: resourceArmApplicationGatewayCreateUpdate,
		Delete: resourceArmApplicationGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"zones": zonesSchema(),

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Required
			"backend_address_pool": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"fqdns": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"ip_addresses": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.IPv4Address,
							},
						},

						// TODO: remove in 2.0
						"fqdn_list": {
							Type:       schema.TypeList,
							Optional:   true,
							Computed:   true,
							Deprecated: "`fqdn_list` has been deprecated in favour of the `fqdns` field",
							MinItems:   1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						// TODO: remove in 2.0
						"ip_address_list": {
							Type:       schema.TypeList,
							Optional:   true,
							Computed:   true,
							Deprecated: "`ip_address_list` has been deprecated in favour of the `ip_addresses` field",
							MinItems:   1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.IPv4Address,
							},
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"backend_http_settings": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.PortNumber,
						},

						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.HTTP),
								string(network.HTTPS),
							}, true),
						},

						"cookie_based_affinity": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Enabled),
								string(network.Disabled),
							}, true),
						},

						"host_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"pick_host_name_from_backend_address": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"request_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 86400),
						},

						"authentication_certificate": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"connection_draining": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},

									"drain_timeout_sec": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 3600),
									},
								},
							},
						},

						"probe_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"probe_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"frontend_ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"private_ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"public_ip_address_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"private_ip_address_allocation": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Dynamic),
								string(network.Static),
							}, true),
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"frontend_port": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.PortNumber,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"gateway_ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"http_listener": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"frontend_ip_configuration_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"frontend_port_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.HTTP),
								string(network.HTTPS),
							}, true),
						},

						"host_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"ssl_certificate_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"require_sni": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"frontend_ip_configuration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"frontend_port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ssl_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"custom_error_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status_code": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(network.HTTPStatus403),
											string(network.HTTPStatus502),
										}, false),
									},

									"custom_error_page_url": {
										Type:     schema.TypeString,
										Required: true,
									},

									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"request_routing_rule": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"rule_type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Basic),
								string(network.PathBasedRouting),
							}, true),
						},

						"http_listener_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"backend_address_pool_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"backend_http_settings_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"url_path_map_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"redirect_configuration_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"backend_address_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backend_http_settings_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"http_listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"url_path_map_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"redirect_configuration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"redirect_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"redirect_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Permanent),
								string(network.Temporary),
								string(network.Found),
								string(network.SeeOther),
							}, false),
						},

						"target_listener_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"target_url": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"include_path": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"include_query_string": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"target_listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.StandardSmall),
								string(network.StandardMedium),
								string(network.StandardLarge),
								string(network.StandardV2),
								string(network.WAFLarge),
								string(network.WAFMedium),
								string(network.WAFV2),
							}, true),
						},

						"tier": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ApplicationGatewayTierStandard),
								string(network.ApplicationGatewayTierStandardV2),
								string(network.ApplicationGatewayTierWAF),
								string(network.ApplicationGatewayTierWAFV2),
							}, true),
						},

						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},
					},
				},
			},

			// Optional
			"authentication_certificate": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"data": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			// TODO: @tombuildsstuff deprecate this in favour of a full `ssl_protocol` block in the future
			"disabled_ssl_protocols": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: suppress.CaseDifference,
					ValidateFunc: validation.StringInSlice([]string{
						string(network.TLSv10),
						string(network.TLSv11),
						string(network.TLSv12),
					}, true),
				},
			},

			"enable_http2": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"probe": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.HTTP),
								string(network.HTTPS),
							}, true),
						},

						"path": {
							Type:     schema.TypeString,
							Required: true,
						},

						"host": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"interval": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"timeout": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"pick_host_name_from_backend_http_settings": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"minimum_servers": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},

						"match": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"body": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "*",
									},

									"status_code": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ssl_certificate": {
				// TODO: should this become a Set?
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"data": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							StateFunc: base64EncodedStateFunc,
						},

						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"public_cert_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"url_path_map": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"default_backend_address_pool_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"default_backend_http_settings_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"default_redirect_configuration_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"path_rule": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"paths": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"backend_address_pool_name": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"backend_http_settings_name": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"redirect_configuration_name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},

									"backend_address_pool_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"backend_http_settings_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"redirect_configuration_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"default_backend_address_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"default_backend_http_settings_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"default_redirect_configuration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"waf_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"firewall_mode": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Detection),
								string(network.Prevention),
							}, true),
						},

						"rule_set_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "OWASP",
						},

						"rule_set_version": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"2.2.9",
								"3.0",
							}, false),
						},
						"file_upload_limit_mb": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 500),
							Default:      100,
						},
						"request_body_check": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"max_request_body_size_kb": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 128),
							Default:      128,
						},
					},
				},
			},

			"custom_error_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status_code": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.HTTPStatus403),
								string(network.HTTPStatus502),
							}, false),
						},

						"custom_error_page_url": {
							Type:     schema.TypeString,
							Required: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmApplicationGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	client := armClient.applicationGatewayClient
	ctx := armClient.StopContext

	log.Printf("[INFO] preparing arguments for Application Gateway creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Application Gateway %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_application_gateway", *existing.ID)
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	enablehttp2 := d.Get("enable_http2").(bool)
	tags := d.Get("tags").(map[string]interface{})

	// Gateway ID is needed to link sub-resources together in expand functions
	gatewayIDFmt := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s"
	gatewayID := fmt.Sprintf(gatewayIDFmt, armClient.subscriptionId, resGroup, name)

	authenticationCertificates := expandApplicationGatewayAuthenticationCertificates(d)
	backendAddressPools := expandApplicationGatewayBackendAddressPools(d)
	backendHTTPSettingsCollection := expandApplicationGatewayBackendHTTPSettings(d, gatewayID)
	frontendIPConfigurations := expandApplicationGatewayFrontendIPConfigurations(d)
	frontendPorts := expandApplicationGatewayFrontendPorts(d)
	gatewayIPConfigurations := expandApplicationGatewayIPConfigurations(d)
	httpListeners := expandApplicationGatewayHTTPListeners(d, gatewayID)
	probes := expandApplicationGatewayProbes(d)
	requestRoutingRules := expandApplicationGatewayRequestRoutingRules(d, gatewayID)
	redirectConfigurations := expandApplicationGatewayRedirectConfigurations(d, gatewayID)
	sku := expandApplicationGatewaySku(d)
	sslCertificates := expandApplicationGatewaySslCertificates(d)
	sslPolicy := expandApplicationGatewaySslPolicy(d)
	customErrorConfigurations := expandApplicationGatewayCustomErrorConfigurations(d.Get("custom_error_configuration").([]interface{}))
	urlPathMaps := expandApplicationGatewayURLPathMaps(d, gatewayID)
	zones := expandZones(d.Get("zones").([]interface{}))

	gateway := network.ApplicationGateway{
		Location: utils.String(location),
		Zones:    zones,

		Tags: expandTags(tags),
		ApplicationGatewayPropertiesFormat: &network.ApplicationGatewayPropertiesFormat{
			AuthenticationCertificates:    authenticationCertificates,
			BackendAddressPools:           backendAddressPools,
			BackendHTTPSettingsCollection: backendHTTPSettingsCollection,
			EnableHTTP2:                   utils.Bool(enablehttp2),
			FrontendIPConfigurations:      frontendIPConfigurations,
			FrontendPorts:                 frontendPorts,
			GatewayIPConfigurations:       gatewayIPConfigurations,
			HTTPListeners:                 httpListeners,
			Probes:                        probes,
			RequestRoutingRules:           requestRoutingRules,
			RedirectConfigurations:        redirectConfigurations,
			Sku:                           sku,
			SslCertificates:               sslCertificates,
			SslPolicy:                     sslPolicy,
			CustomErrorConfigurations:     customErrorConfigurations,
			URLPathMaps:                   urlPathMaps,
		},
	}

	for _, backendHttpSettings := range *backendHTTPSettingsCollection {
		backendHttpSettingsProperties := *backendHttpSettings.ApplicationGatewayBackendHTTPSettingsPropertiesFormat
		if backendHttpSettingsProperties.HostName != nil {
			hostName := *backendHttpSettingsProperties.HostName
			pick := *backendHttpSettingsProperties.PickHostNameFromBackendAddress

			if hostName != "" && pick {
				return fmt.Errorf("Only one of `host_name` or `pick_host_name_from_backend_address` can be set")
			}
		}
	}

	for _, probe := range *probes {
		probeProperties := *probe.ApplicationGatewayProbePropertiesFormat
		host := *probeProperties.Host
		pick := *probeProperties.PickHostNameFromBackendHTTPSettings

		if host == "" && !pick {
			return fmt.Errorf("One of `host` or `pick_host_name_from_backend_http_settings` must be set")
		}

		if host != "" && pick {
			return fmt.Errorf("Only one of `host` or `pick_host_name_from_backend_http_settings` can be set")
		}
	}

	if _, ok := d.GetOk("waf_configuration"); ok {
		gateway.ApplicationGatewayPropertiesFormat.WebApplicationFirewallConfiguration = expandApplicationGatewayWafConfig(d)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, gateway)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Application Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the create/update of Application Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Application Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID of Application Gateway %q (Resource Group %q)", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmApplicationGatewayRead(d, meta)
}

func resourceArmApplicationGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationGatewayClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["applicationGateways"]

	applicationGateway, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(applicationGateway.Response) {
			log.Printf("[DEBUG] Application Gateway %q was not found in Resource Group %q - removing from state", name, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Application Gateway %s: %+v", name, err)
	}

	d.Set("name", applicationGateway.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := applicationGateway.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	d.Set("zones", applicationGateway.Zones)

	if props := applicationGateway.ApplicationGatewayPropertiesFormat; props != nil {
		flattenedCerts := flattenApplicationGatewayAuthenticationCertificates(props.AuthenticationCertificates, d)
		if setErr := d.Set("authentication_certificate", flattenedCerts); setErr != nil {
			return fmt.Errorf("Error setting `authentication_certificate`: %+v", setErr)
		}

		if setErr := d.Set("backend_address_pool", flattenApplicationGatewayBackendAddressPools(props.BackendAddressPools)); setErr != nil {
			return fmt.Errorf("Error setting `backend_address_pool`: %+v", setErr)
		}

		backendHttpSettings, err := flattenApplicationGatewayBackendHTTPSettings(props.BackendHTTPSettingsCollection)
		if err != nil {
			return fmt.Errorf("Error flattening `backend_http_settings`: %+v", err)
		}
		if setErr := d.Set("backend_http_settings", backendHttpSettings); setErr != nil {
			return fmt.Errorf("Error setting `backend_http_settings`: %+v", setErr)
		}

		if setErr := d.Set("disabled_ssl_protocols", flattenApplicationGatewayDisabledSSLProtocols(props.SslPolicy)); setErr != nil {
			return fmt.Errorf("Error setting `disabled_ssl_protocols`: %+v", setErr)
		}

		d.Set("enable_http2", props.EnableHTTP2)

		httpListeners, err := flattenApplicationGatewayHTTPListeners(props.HTTPListeners)
		if err != nil {
			return fmt.Errorf("Error flattening `http_listener`: %+v", err)
		}
		if setErr := d.Set("http_listener", httpListeners); setErr != nil {
			return fmt.Errorf("Error setting `http_listener`: %+v", setErr)
		}

		if setErr := d.Set("frontend_port", flattenApplicationGatewayFrontendPorts(props.FrontendPorts)); setErr != nil {
			return fmt.Errorf("Error setting `frontend_port`: %+v", setErr)
		}

		if setErr := d.Set("frontend_ip_configuration", flattenApplicationGatewayFrontendIPConfigurations(props.FrontendIPConfigurations)); setErr != nil {
			return fmt.Errorf("Error setting `frontend_ip_configuration`: %+v", setErr)
		}

		if setErr := d.Set("gateway_ip_configuration", flattenApplicationGatewayIPConfigurations(props.GatewayIPConfigurations)); setErr != nil {
			return fmt.Errorf("Error setting `gateway_ip_configuration`: %+v", setErr)
		}

		if setErr := d.Set("probe", flattenApplicationGatewayProbes(props.Probes)); setErr != nil {
			return fmt.Errorf("Error setting `probe`: %+v", setErr)
		}

		requestRoutingRules, err := flattenApplicationGatewayRequestRoutingRules(props.RequestRoutingRules)
		if err != nil {
			return fmt.Errorf("Error flattening `request_routing_rule`: %+v", err)
		}
		if setErr := d.Set("request_routing_rule", requestRoutingRules); setErr != nil {
			return fmt.Errorf("Error setting `request_routing_rule`: %+v", setErr)
		}

		redirectConfigurations, err := flattenApplicationGatewayRedirectConfigurations(props.RedirectConfigurations)
		if err != nil {
			return fmt.Errorf("Error flattening `redirect configuration`: %+v", err)
		}
		if setErr := d.Set("redirect_configuration", redirectConfigurations); setErr != nil {
			return fmt.Errorf("Error setting `redirect configuration`: %+v", setErr)
		}

		if setErr := d.Set("sku", flattenApplicationGatewaySku(props.Sku)); setErr != nil {
			return fmt.Errorf("Error setting `sku`: %+v", setErr)
		}

		if setErr := d.Set("ssl_certificate", flattenApplicationGatewaySslCertificates(props.SslCertificates, d)); setErr != nil {
			return fmt.Errorf("Error setting `ssl_certificate`: %+v", setErr)
		}

		if setErr := d.Set("custom_error_configuration", flattenApplicationGatewayCustomErrorConfigurations(props.CustomErrorConfigurations)); setErr != nil {
			return fmt.Errorf("Error setting `custom_error_configuration`: %+v", setErr)
		}

		urlPathMaps, err := flattenApplicationGatewayURLPathMaps(props.URLPathMaps)
		if err != nil {
			return fmt.Errorf("Error flattening `url_path_map`: %+v", err)
		}
		if setErr := d.Set("url_path_map", urlPathMaps); setErr != nil {
			return fmt.Errorf("Error setting `url_path_map`: %+v", setErr)
		}

		if setErr := d.Set("waf_configuration", flattenApplicationGatewayWafConfig(props.WebApplicationFirewallConfiguration)); setErr != nil {
			return fmt.Errorf("Error setting `waf_configuration`: %+v", setErr)
		}
	}

	flattenAndSetTags(d, applicationGateway.Tags)

	return nil
}

func resourceArmApplicationGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationGatewayClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["applicationGateways"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting for Application Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Application Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandApplicationGatewayAuthenticationCertificates(d *schema.ResourceData) *[]network.ApplicationGatewayAuthenticationCertificate {
	vs := d.Get("authentication_certificate").([]interface{})
	results := make([]network.ApplicationGatewayAuthenticationCertificate, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		data := v["data"].(string)

		// data must be base64 encoded
		encodedData := base64Encode(data)

		output := network.ApplicationGatewayAuthenticationCertificate{
			Name: utils.String(name),
			ApplicationGatewayAuthenticationCertificatePropertiesFormat: &network.ApplicationGatewayAuthenticationCertificatePropertiesFormat{
				Data: utils.String(encodedData),
			},
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayAuthenticationCertificates(input *[]network.ApplicationGatewayAuthenticationCertificate, d *schema.ResourceData) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for i, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		// since the certificate data isn't returned we have to load it from the same index
		if existing, ok := d.GetOk("authentication_certificate"); ok && existing != nil {
			existingVals := existing.([]interface{})
			if len(existingVals) >= i {
				existingCerts := existingVals[i].(map[string]interface{})
				if data := existingCerts["data"]; data != nil {
					output["data"] = data.(string)
				}
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayBackendAddressPools(d *schema.ResourceData) *[]network.ApplicationGatewayBackendAddressPool {
	vs := d.Get("backend_address_pool").([]interface{})
	results := make([]network.ApplicationGatewayBackendAddressPool, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})
		backendAddresses := make([]network.ApplicationGatewayBackendAddress, 0)

		for _, ip := range v["fqdns"].([]interface{}) {
			backendAddresses = append(backendAddresses, network.ApplicationGatewayBackendAddress{
				Fqdn: utils.String(ip.(string)),
			})
		}
		for _, ip := range v["ip_addresses"].([]interface{}) {
			backendAddresses = append(backendAddresses, network.ApplicationGatewayBackendAddress{
				IPAddress: utils.String(ip.(string)),
			})
		}

		if len(backendAddresses) == 0 {
			// TODO: remove in 2.0
			for _, ip := range v["ip_address_list"].([]interface{}) {
				backendAddresses = append(backendAddresses, network.ApplicationGatewayBackendAddress{
					IPAddress: utils.String(ip.(string)),
				})
			}
			// TODO: remove in 2.0
			for _, ip := range v["fqdn_list"].([]interface{}) {
				backendAddresses = append(backendAddresses, network.ApplicationGatewayBackendAddress{
					Fqdn: utils.String(ip.(string)),
				})
			}
		}

		name := v["name"].(string)
		output := network.ApplicationGatewayBackendAddressPool{
			Name: utils.String(name),
			ApplicationGatewayBackendAddressPoolPropertiesFormat: &network.ApplicationGatewayBackendAddressPoolPropertiesFormat{
				BackendAddresses: &backendAddresses,
			},
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayBackendAddressPools(input *[]network.ApplicationGatewayBackendAddressPool) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, config := range *input {
		ipAddressList := make([]interface{}, 0)
		fqdnList := make([]interface{}, 0)

		if props := config.ApplicationGatewayBackendAddressPoolPropertiesFormat; props != nil {
			if props.BackendAddresses != nil {
				for _, address := range *props.BackendAddresses {
					if address.IPAddress != nil {
						ipAddressList = append(ipAddressList, *address.IPAddress)
					} else if address.Fqdn != nil {
						fqdnList = append(fqdnList, *address.Fqdn)
					}
				}
			}
		}

		output := map[string]interface{}{
			"fqdns":        fqdnList,
			"ip_addresses": ipAddressList,

			// TODO: deprecated - remove in 2.0
			"ip_address_list": ipAddressList,
			"fqdn_list":       fqdnList,
		}

		if config.ID != nil {
			output["id"] = *config.ID
		}

		if config.Name != nil {
			output["name"] = *config.Name
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayBackendHTTPSettings(d *schema.ResourceData, gatewayID string) *[]network.ApplicationGatewayBackendHTTPSettings {
	results := make([]network.ApplicationGatewayBackendHTTPSettings, 0)
	vs := d.Get("backend_http_settings").([]interface{})

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		path := v["path"].(string)
		port := int32(v["port"].(int))
		protocol := v["protocol"].(string)
		cookieBasedAffinity := v["cookie_based_affinity"].(string)
		pickHostNameFromBackendAddress := v["pick_host_name_from_backend_address"].(bool)
		requestTimeout := int32(v["request_timeout"].(int))

		setting := network.ApplicationGatewayBackendHTTPSettings{
			Name: &name,
			ApplicationGatewayBackendHTTPSettingsPropertiesFormat: &network.ApplicationGatewayBackendHTTPSettingsPropertiesFormat{
				CookieBasedAffinity:            network.ApplicationGatewayCookieBasedAffinity(cookieBasedAffinity),
				Path:                           utils.String(path),
				PickHostNameFromBackendAddress: utils.Bool(pickHostNameFromBackendAddress),
				Port:                           utils.Int32(port),
				Protocol:                       network.ApplicationGatewayProtocol(protocol),
				RequestTimeout:                 utils.Int32(requestTimeout),
				ConnectionDraining:             expandApplicationGatewayConnectionDraining(v),
			},
		}

		hostName := v["host_name"].(string)
		if hostName != "" {
			setting.ApplicationGatewayBackendHTTPSettingsPropertiesFormat.HostName = utils.String(hostName)
		}

		if v["authentication_certificate"] != nil {
			authCerts := v["authentication_certificate"].([]interface{})
			authCertSubResources := make([]network.SubResource, 0)

			for _, rawAuthCert := range authCerts {
				authCert := rawAuthCert.(map[string]interface{})
				authCertName := authCert["name"].(string)
				authCertID := fmt.Sprintf("%s/authenticationCertificates/%s", gatewayID, authCertName)
				authCertSubResource := network.SubResource{
					ID: utils.String(authCertID),
				}

				authCertSubResources = append(authCertSubResources, authCertSubResource)
			}

			setting.ApplicationGatewayBackendHTTPSettingsPropertiesFormat.AuthenticationCertificates = &authCertSubResources
		}

		probeName := v["probe_name"].(string)
		if probeName != "" {
			probeID := fmt.Sprintf("%s/probes/%s", gatewayID, probeName)
			setting.ApplicationGatewayBackendHTTPSettingsPropertiesFormat.Probe = &network.SubResource{
				ID: utils.String(probeID),
			}
		}

		results = append(results, setting)
	}

	return &results
}

func flattenApplicationGatewayBackendHTTPSettings(input *[]network.ApplicationGatewayBackendHTTPSettings) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayBackendHTTPSettingsPropertiesFormat; props != nil {
			output["cookie_based_affinity"] = string(props.CookieBasedAffinity)

			if path := props.Path; path != nil {
				output["path"] = *path
			}
			output["connection_draining"] = flattenApplicationGatewayConnectionDraining(props.ConnectionDraining)

			if port := props.Port; port != nil {
				output["port"] = int(*port)
			}

			if hostName := props.HostName; hostName != nil {
				output["host_name"] = *hostName
			}

			if pickHostNameFromBackendAddress := props.PickHostNameFromBackendAddress; pickHostNameFromBackendAddress != nil {
				output["pick_host_name_from_backend_address"] = *pickHostNameFromBackendAddress
			}

			output["protocol"] = string(props.Protocol)

			if timeout := props.RequestTimeout; timeout != nil {
				output["request_timeout"] = int(*timeout)
			}

			authenticationCertificates := make([]interface{}, 0)
			if certs := props.AuthenticationCertificates; certs != nil {
				for _, cert := range *certs {
					if cert.ID == nil {
						continue
					}

					certId, err := parseAzureResourceID(*cert.ID)
					if err != nil {
						return nil, err
					}

					name := certId.Path["authenticationCertificates"]
					certificate := map[string]interface{}{
						"id":   *cert.ID,
						"name": name,
					}
					authenticationCertificates = append(authenticationCertificates, certificate)
				}
			}
			output["authentication_certificate"] = authenticationCertificates

			if probe := props.Probe; probe != nil {
				if probe.ID != nil {
					id, err := parseAzureResourceID(*probe.ID)
					if err != nil {
						return results, err
					}

					output["probe_name"] = id.Path["probes"]
					output["probe_id"] = *probe.ID
				}
			}
		}

		results = append(results, output)
	}

	return results, nil
}

func expandApplicationGatewayConnectionDraining(d map[string]interface{}) *network.ApplicationGatewayConnectionDraining {
	connectionsRaw := d["connection_draining"].([]interface{})

	if len(connectionsRaw) <= 0 {
		return nil
	}

	connectionRaw := connectionsRaw[0].(map[string]interface{})

	return &network.ApplicationGatewayConnectionDraining{
		Enabled:           utils.Bool(connectionRaw["enabled"].(bool)),
		DrainTimeoutInSec: utils.Int32(int32(connectionRaw["drain_timeout_sec"].(int))),
	}
}

func flattenApplicationGatewayConnectionDraining(input *network.ApplicationGatewayConnectionDraining) []interface{} {
	result := map[string]interface{}{}
	if input == nil {
		return []interface{}{}
	}

	if v := input.Enabled; v != nil {
		result["enabled"] = *v
	}
	if v := input.DrainTimeoutInSec; v != nil {
		result["drain_timeout_sec"] = *v
	}

	return []interface{}{result}
}

func expandApplicationGatewaySslPolicy(d *schema.ResourceData) *network.ApplicationGatewaySslPolicy {
	vs := d.Get("disabled_ssl_protocols").([]interface{})
	results := make([]network.ApplicationGatewaySslProtocol, 0)

	for _, v := range vs {
		results = append(results, network.ApplicationGatewaySslProtocol(v.(string)))
	}

	return &network.ApplicationGatewaySslPolicy{
		DisabledSslProtocols: &results,
	}
}

func flattenApplicationGatewayDisabledSSLProtocols(input *network.ApplicationGatewaySslPolicy) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.DisabledSslProtocols == nil {
		return results
	}

	for _, v := range *input.DisabledSslProtocols {
		results = append(results, string(v))
	}

	return results
}

func expandApplicationGatewayHTTPListeners(d *schema.ResourceData, gatewayID string) *[]network.ApplicationGatewayHTTPListener {
	vs := d.Get("http_listener").([]interface{})
	results := make([]network.ApplicationGatewayHTTPListener, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		frontendIPConfigName := v["frontend_ip_configuration_name"].(string)
		frontendPortName := v["frontend_port_name"].(string)
		protocol := v["protocol"].(string)
		requireSNI := v["require_sni"].(bool)

		frontendIPConfigID := fmt.Sprintf("%s/frontendIPConfigurations/%s", gatewayID, frontendIPConfigName)
		frontendPortID := fmt.Sprintf("%s/frontendPorts/%s", gatewayID, frontendPortName)

		customErrorConfigurations := expandApplicationGatewayCustomErrorConfigurations(v["custom_error_configuration"].([]interface{}))

		listener := network.ApplicationGatewayHTTPListener{
			Name: utils.String(name),
			ApplicationGatewayHTTPListenerPropertiesFormat: &network.ApplicationGatewayHTTPListenerPropertiesFormat{
				FrontendIPConfiguration: &network.SubResource{
					ID: utils.String(frontendIPConfigID),
				},
				FrontendPort: &network.SubResource{
					ID: utils.String(frontendPortID),
				},
				Protocol:                    network.ApplicationGatewayProtocol(protocol),
				RequireServerNameIndication: utils.Bool(requireSNI),
				CustomErrorConfigurations:   customErrorConfigurations,
			},
		}

		if host := v["host_name"].(string); host != "" {
			listener.ApplicationGatewayHTTPListenerPropertiesFormat.HostName = &host
		}

		if sslCertName := v["ssl_certificate_name"].(string); sslCertName != "" {
			certID := fmt.Sprintf("%s/sslCertificates/%s", gatewayID, sslCertName)
			listener.ApplicationGatewayHTTPListenerPropertiesFormat.SslCertificate = &network.SubResource{
				ID: utils.String(certID),
			}
		}

		results = append(results, listener)
	}

	return &results
}

func flattenApplicationGatewayHTTPListeners(input *[]network.ApplicationGatewayHTTPListener) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayHTTPListenerPropertiesFormat; props != nil {
			if port := props.FrontendPort; port != nil {
				if port.ID != nil {
					portId, err := parseAzureResourceID(*port.ID)
					if err != nil {
						return nil, err
					}
					portName := portId.Path["frontendPorts"]
					output["frontend_port_name"] = portName
					output["frontend_port_id"] = *port.ID
				}
			}

			if feConfig := props.FrontendIPConfiguration; feConfig != nil {
				if feConfig.ID != nil {
					feConfigId, err := parseAzureResourceID(*feConfig.ID)
					if err != nil {
						return nil, err
					}
					frontendName := feConfigId.Path["frontendIPConfigurations"]
					output["frontend_ip_configuration_name"] = frontendName
					output["frontend_ip_configuration_id"] = *feConfig.ID
				}
			}

			if hostname := props.HostName; hostname != nil {
				output["host_name"] = *hostname
			}

			output["protocol"] = string(props.Protocol)

			if cert := props.SslCertificate; cert != nil {
				if cert.ID != nil {
					certId, err := parseAzureResourceID(*cert.ID)
					if err != nil {
						return nil, err
					}
					sslCertName := certId.Path["sslCertificates"]

					output["ssl_certificate_name"] = sslCertName
					output["ssl_certificate_id"] = *cert.ID
				}
			}

			if sni := props.RequireServerNameIndication; sni != nil {
				output["require_sni"] = *sni
			}

			output["custom_error_configuration"] = flattenApplicationGatewayCustomErrorConfigurations(props.CustomErrorConfigurations)
		}

		results = append(results, output)
	}

	return results, nil
}

func expandApplicationGatewayIPConfigurations(d *schema.ResourceData) *[]network.ApplicationGatewayIPConfiguration {
	vs := d.Get("gateway_ip_configuration").([]interface{})
	results := make([]network.ApplicationGatewayIPConfiguration, 0)

	for _, configRaw := range vs {
		data := configRaw.(map[string]interface{})

		name := data["name"].(string)
		subnetID := data["subnet_id"].(string)

		output := network.ApplicationGatewayIPConfiguration{
			Name: utils.String(name),
			ApplicationGatewayIPConfigurationPropertiesFormat: &network.ApplicationGatewayIPConfigurationPropertiesFormat{
				Subnet: &network.SubResource{
					ID: utils.String(subnetID),
				},
			},
		}
		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayIPConfigurations(input *[]network.ApplicationGatewayIPConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayIPConfigurationPropertiesFormat; props != nil {
			if subnet := props.Subnet; subnet != nil {
				if subnet.ID != nil {
					output["subnet_id"] = *subnet.ID
				}
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayFrontendPorts(d *schema.ResourceData) *[]network.ApplicationGatewayFrontendPort {
	vs := d.Get("frontend_port").([]interface{})
	results := make([]network.ApplicationGatewayFrontendPort, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		port := int32(v["port"].(int))

		output := network.ApplicationGatewayFrontendPort{
			Name: utils.String(name),
			ApplicationGatewayFrontendPortPropertiesFormat: &network.ApplicationGatewayFrontendPortPropertiesFormat{
				Port: utils.Int32(port),
			},
		}
		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayFrontendPorts(input *[]network.ApplicationGatewayFrontendPort) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayFrontendPortPropertiesFormat; props != nil {
			if props.Port != nil {
				output["port"] = int(*props.Port)
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayFrontendIPConfigurations(d *schema.ResourceData) *[]network.ApplicationGatewayFrontendIPConfiguration {
	vs := d.Get("frontend_ip_configuration").([]interface{})
	results := make([]network.ApplicationGatewayFrontendIPConfiguration, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		properties := network.ApplicationGatewayFrontendIPConfigurationPropertiesFormat{}

		if val := v["subnet_id"].(string); val != "" {
			properties.Subnet = &network.SubResource{
				ID: utils.String(val),
			}
		}

		if val := v["private_ip_address_allocation"].(string); val != "" {
			properties.PrivateIPAllocationMethod = network.IPAllocationMethod(val)
		}

		if val := v["private_ip_address"].(string); val != "" {
			properties.PrivateIPAddress = utils.String(val)
		}

		if val := v["public_ip_address_id"].(string); val != "" {
			properties.PublicIPAddress = &network.SubResource{
				ID: utils.String(val),
			}
		}

		name := v["name"].(string)
		output := network.ApplicationGatewayFrontendIPConfiguration{
			Name: utils.String(name),
			ApplicationGatewayFrontendIPConfigurationPropertiesFormat: &properties,
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayFrontendIPConfigurations(input *[]network.ApplicationGatewayFrontendIPConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, config := range *input {
		output := make(map[string]interface{})
		if config.ID != nil {
			output["id"] = *config.ID
		}

		if config.Name != nil {
			output["name"] = *config.Name
		}

		if props := config.ApplicationGatewayFrontendIPConfigurationPropertiesFormat; props != nil {
			output["private_ip_address_allocation"] = string(props.PrivateIPAllocationMethod)

			if props.Subnet != nil && props.Subnet.ID != nil {
				output["subnet_id"] = *props.Subnet.ID
			}

			if props.PrivateIPAddress != nil {
				output["private_ip_address"] = *props.PrivateIPAddress
			}

			if props.PublicIPAddress != nil && props.PublicIPAddress.ID != nil {
				output["public_ip_address_id"] = *props.PublicIPAddress.ID
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayProbes(d *schema.ResourceData) *[]network.ApplicationGatewayProbe {
	vs := d.Get("probe").([]interface{})
	results := make([]network.ApplicationGatewayProbe, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		host := v["host"].(string)
		interval := int32(v["interval"].(int))
		minServers := int32(v["minimum_servers"].(int))
		name := v["name"].(string)
		probePath := v["path"].(string)
		protocol := v["protocol"].(string)
		timeout := int32(v["timeout"].(int))
		unhealthyThreshold := int32(v["unhealthy_threshold"].(int))
		pickHostNameFromBackendHTTPSettings := v["pick_host_name_from_backend_http_settings"].(bool)

		output := network.ApplicationGatewayProbe{
			Name: utils.String(name),
			ApplicationGatewayProbePropertiesFormat: &network.ApplicationGatewayProbePropertiesFormat{
				Host:                                utils.String(host),
				Interval:                            utils.Int32(interval),
				MinServers:                          utils.Int32(minServers),
				Path:                                utils.String(probePath),
				Protocol:                            network.ApplicationGatewayProtocol(protocol),
				Timeout:                             utils.Int32(timeout),
				UnhealthyThreshold:                  utils.Int32(unhealthyThreshold),
				PickHostNameFromBackendHTTPSettings: utils.Bool(pickHostNameFromBackendHTTPSettings),
			},
		}

		matchConfigs := v["match"].([]interface{})
		if len(matchConfigs) > 0 {
			match := matchConfigs[0].(map[string]interface{})
			matchBody := match["body"].(string)

			statusCodes := make([]string, 0)
			for _, statusCode := range match["status_code"].([]interface{}) {
				statusCodes = append(statusCodes, statusCode.(string))
			}

			output.ApplicationGatewayProbePropertiesFormat.Match = &network.ApplicationGatewayProbeHealthResponseMatch{
				Body:        utils.String(matchBody),
				StatusCodes: &statusCodes,
			}
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayProbes(input *[]network.ApplicationGatewayProbe) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayProbePropertiesFormat; props != nil {
			output["protocol"] = string(props.Protocol)

			if host := props.Host; host != nil {
				output["host"] = *host
			}

			if path := props.Path; path != nil {
				output["path"] = *path
			}

			if interval := props.Interval; interval != nil {
				output["interval"] = int(*interval)
			}

			if timeout := props.Timeout; timeout != nil {
				output["timeout"] = int(*timeout)
			}

			if threshold := props.UnhealthyThreshold; threshold != nil {
				output["unhealthy_threshold"] = int(*threshold)
			}

			if pickHostNameFromBackendHTTPSettings := props.PickHostNameFromBackendHTTPSettings; pickHostNameFromBackendHTTPSettings != nil {
				output["pick_host_name_from_backend_http_settings"] = *pickHostNameFromBackendHTTPSettings
			}

			if minServers := props.MinServers; minServers != nil {
				output["minimum_servers"] = int(*minServers)
			}

			matches := make([]interface{}, 0)
			if match := props.Match; match != nil {
				matchConfig := map[string]interface{}{}
				if body := match.Body; body != nil {
					matchConfig["body"] = *body
				}

				statusCodes := make([]interface{}, 0)
				if match.StatusCodes != nil {
					for _, status := range *match.StatusCodes {
						statusCodes = append(statusCodes, status)
					}
				}
				matchConfig["status_code"] = statusCodes
				matches = append(matches, matchConfig)
			}
			output["match"] = matches
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayRequestRoutingRules(d *schema.ResourceData, gatewayID string) *[]network.ApplicationGatewayRequestRoutingRule {
	vs := d.Get("request_routing_rule").([]interface{})
	results := make([]network.ApplicationGatewayRequestRoutingRule, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		ruleType := v["rule_type"].(string)
		httpListenerName := v["http_listener_name"].(string)
		httpListenerID := fmt.Sprintf("%s/httpListeners/%s", gatewayID, httpListenerName)

		rule := network.ApplicationGatewayRequestRoutingRule{
			Name: utils.String(name),
			ApplicationGatewayRequestRoutingRulePropertiesFormat: &network.ApplicationGatewayRequestRoutingRulePropertiesFormat{
				RuleType: network.ApplicationGatewayRequestRoutingRuleType(ruleType),
				HTTPListener: &network.SubResource{
					ID: utils.String(httpListenerID),
				},
			},
		}

		if backendAddressPoolName := v["backend_address_pool_name"].(string); backendAddressPoolName != "" {
			backendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, backendAddressPoolName)
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.BackendAddressPool = &network.SubResource{
				ID: utils.String(backendAddressPoolID),
			}
		}

		if backendHTTPSettingsName := v["backend_http_settings_name"].(string); backendHTTPSettingsName != "" {
			backendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, backendHTTPSettingsName)
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.BackendHTTPSettings = &network.SubResource{
				ID: utils.String(backendHTTPSettingsID),
			}
		}

		if urlPathMapName := v["url_path_map_name"].(string); urlPathMapName != "" {
			urlPathMapID := fmt.Sprintf("%s/urlPathMaps/%s", gatewayID, urlPathMapName)
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.URLPathMap = &network.SubResource{
				ID: utils.String(urlPathMapID),
			}
		}

		if redirectConfigName := v["redirect_configuration_name"].(string); redirectConfigName != "" {
			redirectConfigID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, redirectConfigName)
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.RedirectConfiguration = &network.SubResource{
				ID: utils.String(redirectConfigID),
			}
		}

		results = append(results, rule)
	}

	return &results
}

func flattenApplicationGatewayRequestRoutingRules(input *[]network.ApplicationGatewayRequestRoutingRule) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		if props := config.ApplicationGatewayRequestRoutingRulePropertiesFormat; props != nil {

			if applicationGatewayHasSubResource(props.BackendAddressPool) && applicationGatewayHasSubResource(props.RedirectConfiguration) {
				return nil, fmt.Errorf("[ERROR] Conflict between `backend_address_pool_name` and `redirect_configuration_name` (back-end pool not applicable when redirection specified)")
			}

			if applicationGatewayHasSubResource(props.BackendHTTPSettings) && applicationGatewayHasSubResource(props.RedirectConfiguration) {
				return nil, fmt.Errorf("[ERROR] Conflict between `backend_http_settings_name` and `redirect_configuration_name` (back-end settings not applicable when redirection specified)")
			}

			output := map[string]interface{}{
				"rule_type": string(props.RuleType),
			}

			if config.ID != nil {
				output["id"] = *config.ID
			}

			if config.Name != nil {
				output["name"] = *config.Name
			}

			if pool := props.BackendAddressPool; pool != nil {
				if pool.ID != nil {
					poolId, err := parseAzureResourceID(*pool.ID)
					if err != nil {
						return nil, err
					}
					backendAddressPoolName := poolId.Path["backendAddressPools"]
					output["backend_address_pool_name"] = backendAddressPoolName
					output["backend_address_pool_id"] = *pool.ID
				}
			}

			if settings := props.BackendHTTPSettings; settings != nil {
				if settings.ID != nil {
					settingsId, err := parseAzureResourceID(*settings.ID)
					if err != nil {
						return nil, err
					}
					backendHTTPSettingsName := settingsId.Path["backendHttpSettingsCollection"]
					output["backend_http_settings_name"] = backendHTTPSettingsName
					output["backend_http_settings_id"] = *settings.ID
				}
			}

			if listener := props.HTTPListener; listener != nil {
				if listener.ID != nil {
					listenerId, err := parseAzureResourceID(*listener.ID)
					if err != nil {
						return nil, err
					}
					httpListenerName := listenerId.Path["httpListeners"]
					output["http_listener_id"] = *listener.ID
					output["http_listener_name"] = httpListenerName
				}
			}

			if pathMap := props.URLPathMap; pathMap != nil {
				if pathMap.ID != nil {
					pathMapId, err := parseAzureResourceID(*pathMap.ID)
					if err != nil {
						return nil, err
					}
					urlPathMapName := pathMapId.Path["urlPathMaps"]
					output["url_path_map_name"] = urlPathMapName
					output["url_path_map_id"] = *pathMap.ID
				}
			}

			if redirect := props.RedirectConfiguration; redirect != nil {
				if redirect.ID != nil {
					redirectId, err := parseAzureResourceID(*redirect.ID)
					if err != nil {
						return nil, err
					}
					redirectName := redirectId.Path["redirectConfigurations"]
					output["redirect_configuration_name"] = redirectName
					output["redirect_configuration_id"] = *redirect.ID
				}
			}

			results = append(results, output)
		}
	}

	return results, nil
}

func applicationGatewayHasSubResource(subResource *network.SubResource) bool {
	return subResource != nil && subResource.ID != nil && *subResource.ID != ""
}

func expandApplicationGatewayRedirectConfigurations(d *schema.ResourceData, gatewayID string) *[]network.ApplicationGatewayRedirectConfiguration {

	vs := d.Get("redirect_configuration").([]interface{})
	results := make([]network.ApplicationGatewayRedirectConfiguration, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		redirectType := v["redirect_type"].(string)

		output := network.ApplicationGatewayRedirectConfiguration{
			Name: utils.String(name),
			ApplicationGatewayRedirectConfigurationPropertiesFormat: &network.ApplicationGatewayRedirectConfigurationPropertiesFormat{
				RedirectType: network.ApplicationGatewayRedirectType(redirectType),
			},
		}

		if includePath := v["include_path"].(bool); includePath {
			output.ApplicationGatewayRedirectConfigurationPropertiesFormat.IncludePath = utils.Bool(includePath)
		}

		if includeQueryString := v["include_query_string"].(bool); includeQueryString {
			output.ApplicationGatewayRedirectConfigurationPropertiesFormat.IncludeQueryString = utils.Bool(includeQueryString)
		}

		if targetListenerName := v["target_listener_name"].(string); targetListenerName != "" {
			targetListenerID := fmt.Sprintf("%s/httpListeners/%s", gatewayID, targetListenerName)
			output.ApplicationGatewayRedirectConfigurationPropertiesFormat.TargetListener = &network.SubResource{
				ID: utils.String(targetListenerID),
			}
		}

		if targetUrl := v["target_url"].(string); targetUrl != "" {
			output.ApplicationGatewayRedirectConfigurationPropertiesFormat.TargetURL = utils.String(targetUrl)
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayRedirectConfigurations(input *[]network.ApplicationGatewayRedirectConfiguration) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		if props := config.ApplicationGatewayRedirectConfigurationPropertiesFormat; props != nil {

			if !applicationGatewayHasSubResource(config.TargetListener) && (config.TargetURL == nil || *config.TargetURL == "") {
				return nil, fmt.Errorf("[ERROR] Set either `target_listener_name` or `target_url`")
			}
			if applicationGatewayHasSubResource(config.TargetListener) && config.TargetURL != nil && *config.TargetURL != "" {
				return nil, fmt.Errorf("[ERROR] Conflict between `target_listener_name` and `target_url` (redirection is either to URL or target listener)")
			}

			if config.TargetURL != nil && *config.TargetURL != "" && config.IncludePath != nil {
				return nil, fmt.Errorf("[ERROR] `include_path` is not a valid option when `target_url` is set")
			}

			output := map[string]interface{}{
				"redirect_type": string(props.RedirectType),
			}

			if config.ID != nil {
				output["id"] = *config.ID
			}

			if config.Name != nil {
				output["name"] = *config.Name
			}

			if listener := props.TargetListener; listener != nil {
				if listener.ID != nil {
					listenerId, err := parseAzureResourceID(*listener.ID)
					if err != nil {
						return nil, err
					}
					targetListenerName := listenerId.Path["httpListeners"]
					output["target_listener_name"] = targetListenerName
					output["target_listener_id"] = *listener.ID
				}
			}

			if config.TargetURL != nil {
				output["target_url"] = *config.TargetURL
			}

			if config.IncludePath != nil {
				output["include_path"] = *config.IncludePath
			}

			if config.IncludeQueryString != nil {
				output["include_query_string"] = *config.IncludeQueryString
			}

			results = append(results, output)
		}
	}

	return results, nil
}

func expandApplicationGatewaySku(d *schema.ResourceData) *network.ApplicationGatewaySku {
	vs := d.Get("sku").([]interface{})
	v := vs[0].(map[string]interface{})

	name := v["name"].(string)
	tier := v["tier"].(string)
	capacity := int32(v["capacity"].(int))

	return &network.ApplicationGatewaySku{
		Name:     network.ApplicationGatewaySkuName(name),
		Tier:     network.ApplicationGatewayTier(tier),
		Capacity: utils.Int32(capacity),
	}
}

func flattenApplicationGatewaySku(input *network.ApplicationGatewaySku) []interface{} {
	result := make(map[string]interface{})

	result["name"] = string(input.Name)
	result["tier"] = string(input.Tier)
	if input.Capacity != nil {
		result["capacity"] = int(*input.Capacity)
	}

	return []interface{}{result}
}

func expandApplicationGatewaySslCertificates(d *schema.ResourceData) *[]network.ApplicationGatewaySslCertificate {
	vs := d.Get("ssl_certificate").([]interface{})
	results := make([]network.ApplicationGatewaySslCertificate, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		data := v["data"].(string)
		password := v["password"].(string)

		// data must be base64 encoded
		data = base64Encode(data)

		output := network.ApplicationGatewaySslCertificate{
			Name: utils.String(name),
			ApplicationGatewaySslCertificatePropertiesFormat: &network.ApplicationGatewaySslCertificatePropertiesFormat{
				Data:     utils.String(data),
				Password: utils.String(password),
			},
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewaySslCertificates(input *[]network.ApplicationGatewaySslCertificate, d *schema.ResourceData) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}
		if v.Name == nil {
			continue
		}

		name := *v.Name

		if v.ID != nil {
			output["id"] = *v.ID
		}

		output["name"] = name

		if props := v.ApplicationGatewaySslCertificatePropertiesFormat; props != nil {
			if data := props.PublicCertData; data != nil {
				output["public_cert_data"] = *data
			}
		}

		// since the certificate data isn't returned we have to load it from the same index
		if existing, ok := d.GetOk("ssl_certificate"); ok && existing != nil {
			existingVals := existing.([]interface{})
			for _, existingVal := range existingVals {
				existingCerts := existingVal.(map[string]interface{})
				existingName := existingCerts["name"].(string)

				if name == existingName {
					if data := existingCerts["data"]; data != nil {
						v := base64Encode(data.(string))
						output["data"] = v
					}

					if password := existingCerts["password"]; password != nil {
						output["password"] = password.(string)
					}
				}
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayURLPathMaps(d *schema.ResourceData, gatewayID string) *[]network.ApplicationGatewayURLPathMap {
	vs := d.Get("url_path_map").([]interface{})
	results := make([]network.ApplicationGatewayURLPathMap, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)

		pathRules := make([]network.ApplicationGatewayPathRule, 0)
		for _, ruleConfig := range v["path_rule"].([]interface{}) {
			ruleConfigMap := ruleConfig.(map[string]interface{})

			ruleName := ruleConfigMap["name"].(string)

			rulePaths := make([]string, 0)
			for _, rulePath := range ruleConfigMap["paths"].([]interface{}) {
				rulePaths = append(rulePaths, rulePath.(string))
			}

			rule := network.ApplicationGatewayPathRule{
				Name: utils.String(ruleName),
				ApplicationGatewayPathRulePropertiesFormat: &network.ApplicationGatewayPathRulePropertiesFormat{
					Paths: &rulePaths,
				},
			}

			if backendAddressPoolName := ruleConfigMap["backend_address_pool_name"].(string); backendAddressPoolName != "" {
				backendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, backendAddressPoolName)
				rule.ApplicationGatewayPathRulePropertiesFormat.BackendAddressPool = &network.SubResource{
					ID: utils.String(backendAddressPoolID),
				}
			}

			if backendHTTPSettingsName := ruleConfigMap["backend_http_settings_name"].(string); backendHTTPSettingsName != "" {
				backendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, backendHTTPSettingsName)
				rule.ApplicationGatewayPathRulePropertiesFormat.BackendHTTPSettings = &network.SubResource{
					ID: utils.String(backendHTTPSettingsID),
				}
			}

			if redirectConfigurationName := ruleConfigMap["redirect_configuration_name"].(string); redirectConfigurationName != "" {
				redirectConfigurationID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, redirectConfigurationName)
				rule.ApplicationGatewayPathRulePropertiesFormat.RedirectConfiguration = &network.SubResource{
					ID: utils.String(redirectConfigurationID),
				}
			}

			pathRules = append(pathRules, rule)
		}

		output := network.ApplicationGatewayURLPathMap{
			Name: utils.String(name),
			ApplicationGatewayURLPathMapPropertiesFormat: &network.ApplicationGatewayURLPathMapPropertiesFormat{
				PathRules: &pathRules,
			},
		}

		// treating these three as optional as seems necessary when redirection is also an alternative. Not explicit in the documentation, though
		// see https://docs.microsoft.com/en-us/rest/api/application-gateway/applicationgateways/createorupdate#applicationgatewayurlpathmap
		// see also az docs https://docs.microsoft.com/en-us/cli/azure/network/application-gateway/url-path-map?view=azure-cli-latest#az-network-application-gateway-url-path-map-create

		if defaultBackendAddressPoolName := v["default_backend_address_pool_name"].(string); defaultBackendAddressPoolName != "" {
			defaultBackendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, defaultBackendAddressPoolName)
			output.ApplicationGatewayURLPathMapPropertiesFormat.DefaultBackendAddressPool = &network.SubResource{
				ID: utils.String(defaultBackendAddressPoolID),
			}
		}

		if defaultBackendHTTPSettingsName := v["default_backend_http_settings_name"].(string); defaultBackendHTTPSettingsName != "" {
			defaultBackendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, defaultBackendHTTPSettingsName)
			output.ApplicationGatewayURLPathMapPropertiesFormat.DefaultBackendHTTPSettings = &network.SubResource{
				ID: utils.String(defaultBackendHTTPSettingsID),
			}
		}

		if defaultRedirectConfigurationName := v["default_redirect_configuration_name"].(string); defaultRedirectConfigurationName != "" {
			defaultRedirectConfigurationID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, defaultRedirectConfigurationName)
			output.ApplicationGatewayURLPathMapPropertiesFormat.DefaultRedirectConfiguration = &network.SubResource{
				ID: utils.String(defaultRedirectConfigurationID),
			}
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayURLPathMaps(input *[]network.ApplicationGatewayURLPathMap) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayURLPathMapPropertiesFormat; props != nil {
			if backendPool := props.DefaultBackendAddressPool; backendPool != nil && backendPool.ID != nil {
				poolId, err := parseAzureResourceID(*backendPool.ID)
				if err != nil {
					return nil, err
				}
				backendAddressPoolName := poolId.Path["backendAddressPools"]
				output["default_backend_address_pool_name"] = backendAddressPoolName
				output["default_backend_address_pool_id"] = *backendPool.ID
			}

			if settings := props.DefaultBackendHTTPSettings; settings != nil && settings.ID != nil {
				settingsId, err := parseAzureResourceID(*settings.ID)
				if err != nil {
					return nil, err
				}
				backendHTTPSettingsName := settingsId.Path["backendHttpSettingsCollection"]
				output["default_backend_http_settings_name"] = backendHTTPSettingsName
				output["default_backend_http_settings_id"] = *settings.ID
			}

			if redirect := props.DefaultRedirectConfiguration; redirect != nil && redirect.ID != nil {
				settingsId, err := parseAzureResourceID(*redirect.ID)
				if err != nil {
					return nil, err
				}
				redirectConfigurationName := settingsId.Path["redirectConfigurations"]
				output["default_redirect_configuration_name"] = redirectConfigurationName
				output["default_redirect_configuration_id"] = *redirect.ID
			}

			pathRules := make([]interface{}, 0)
			if rules := props.PathRules; rules != nil {
				for _, rule := range *rules {
					ruleOutput := map[string]interface{}{}

					if rule.ID != nil {
						ruleOutput["id"] = *rule.ID
					}

					if rule.Name != nil {
						ruleOutput["name"] = *rule.Name
					}

					if ruleProps := rule.ApplicationGatewayPathRulePropertiesFormat; ruleProps != nil {
						if applicationGatewayHasSubResource(props.DefaultBackendAddressPool) && applicationGatewayHasSubResource(ruleProps.RedirectConfiguration) {
							return nil, fmt.Errorf("[ERROR] Conflict between `default_backend_address_pool_name` and `redirect_configuration_name` (default back-end pool not applicable when redirection specified)")
						}

						if applicationGatewayHasSubResource(ruleProps.BackendAddressPool) && applicationGatewayHasSubResource(ruleProps.RedirectConfiguration) {
							return nil, fmt.Errorf("[ERROR] Conflict between `backend_address_pool_name` and `redirect_configuration_name` (back-end pool not applicable when redirection specified)")
						}

						if applicationGatewayHasSubResource(props.DefaultBackendHTTPSettings) && applicationGatewayHasSubResource(ruleProps.RedirectConfiguration) {
							return nil, fmt.Errorf("[ERROR] Conflict between `default_backend_http_settings_name` and `redirect_configuration_name` (default back-end settings not applicable when redirection specified)")
						}

						if applicationGatewayHasSubResource(ruleProps.BackendHTTPSettings) && applicationGatewayHasSubResource(ruleProps.RedirectConfiguration) {
							return nil, fmt.Errorf("[ERROR] Conflict between `backend_http_settings_name` and `redirect_configuration_name` (back-end settings not applicable when redirection specified)")
						}

						if pool := ruleProps.BackendAddressPool; pool != nil && pool.ID != nil {
							poolId, err := parseAzureResourceID(*pool.ID)
							if err != nil {
								return nil, err
							}
							backendAddressPoolName2 := poolId.Path["backendAddressPools"]
							ruleOutput["backend_address_pool_name"] = backendAddressPoolName2
							ruleOutput["backend_address_pool_id"] = *pool.ID
						}

						if backend := ruleProps.BackendHTTPSettings; backend != nil && backend.ID != nil {
							backendId, err := parseAzureResourceID(*backend.ID)
							if err != nil {
								return nil, err
							}
							backendHTTPSettingsName2 := backendId.Path["backendHttpSettingsCollection"]
							ruleOutput["backend_http_settings_name"] = backendHTTPSettingsName2
							ruleOutput["backend_http_settings_id"] = *backend.ID
						}

						if redirect := ruleProps.RedirectConfiguration; redirect != nil && redirect.ID != nil {
							redirectId, err := parseAzureResourceID(*redirect.ID)
							if err != nil {
								return nil, err
							}
							redirectConfigurationName2 := redirectId.Path["redirectConfigurations"]
							ruleOutput["redirect_configuration_name"] = redirectConfigurationName2
							ruleOutput["redirect_configuration_id"] = *redirect.ID
						}

						pathOutputs := make([]interface{}, 0)
						if paths := ruleProps.Paths; paths != nil {
							for _, rulePath := range *paths {
								pathOutputs = append(pathOutputs, rulePath)
							}
						}
						ruleOutput["paths"] = pathOutputs
					}

					pathRules = append(pathRules, ruleOutput)
				}
			}
			output["path_rule"] = pathRules
		}

		results = append(results, output)
	}

	return results, nil
}

func expandApplicationGatewayWafConfig(d *schema.ResourceData) *network.ApplicationGatewayWebApplicationFirewallConfiguration {
	vs := d.Get("waf_configuration").([]interface{})
	v := vs[0].(map[string]interface{})

	enabled := v["enabled"].(bool)
	mode := v["firewall_mode"].(string)
	ruleSetType := v["rule_set_type"].(string)
	ruleSetVersion := v["rule_set_version"].(string)
	fileUploadLimitInMb := v["file_upload_limit_mb"].(int)
	requestBodyCheck := v["request_body_check"].(bool)
	maxRequestBodySizeInKb := v["max_request_body_size_kb"].(int)

	return &network.ApplicationGatewayWebApplicationFirewallConfiguration{
		Enabled:                utils.Bool(enabled),
		FirewallMode:           network.ApplicationGatewayFirewallMode(mode),
		RuleSetType:            utils.String(ruleSetType),
		RuleSetVersion:         utils.String(ruleSetVersion),
		FileUploadLimitInMb:    utils.Int32(int32(fileUploadLimitInMb)),
		RequestBodyCheck:       utils.Bool(requestBodyCheck),
		MaxRequestBodySizeInKb: utils.Int32(int32(maxRequestBodySizeInKb)),
	}
}

func flattenApplicationGatewayWafConfig(input *network.ApplicationGatewayWebApplicationFirewallConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	output := make(map[string]interface{})

	if input.Enabled != nil {
		output["enabled"] = *input.Enabled
	}

	output["firewall_mode"] = string(input.FirewallMode)

	if input.RuleSetType != nil {
		output["rule_set_type"] = *input.RuleSetType
	}

	if input.RuleSetVersion != nil {
		output["rule_set_version"] = *input.RuleSetVersion
	}

	if input.FileUploadLimitInMb != nil {
		output["file_upload_limit_mb"] = int(*input.FileUploadLimitInMb)
	}

	if input.RequestBodyCheck != nil {
		output["request_body_check"] = *input.RequestBodyCheck
	}

	if input.MaxRequestBodySizeInKb != nil {
		output["max_request_body_size_kb"] = int(*input.MaxRequestBodySizeInKb)
	}

	results = append(results, output)

	return results
}

func expandApplicationGatewayCustomErrorConfigurations(vs []interface{}) *[]network.ApplicationGatewayCustomError {
	results := make([]network.ApplicationGatewayCustomError, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})
		statusCode := v["status_code"].(string)
		customErrorPageUrl := v["custom_error_page_url"].(string)

		output := network.ApplicationGatewayCustomError{
			StatusCode:         network.ApplicationGatewayCustomErrorStatusCode(statusCode),
			CustomErrorPageURL: utils.String(customErrorPageUrl),
		}
		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayCustomErrorConfigurations(input *[]network.ApplicationGatewayCustomError) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		output["status_code"] = string(v.StatusCode)

		if v.CustomErrorPageURL != nil {
			output["custom_error_page_url"] = *v.CustomErrorPageURL
		}

		results = append(results, output)
	}

	return results
}
