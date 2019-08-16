package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/frontdoor/mgmt/2019-04-01/frontdoor"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmFrontDoor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmFrontDoorRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"load_balancer_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"routing_rule": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"accepted_protocols": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 2,
							Elem: &schema.Schema{
								Type:     schema.TypeString,
								Computed: true,
							},
						},
						"patterns_to_match": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 25,
							Elem: &schema.Schema{
								Type:     schema.TypeString,
								Computed: true,
							},
						},
						"frontend_endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 100,
							Elem: &schema.Schema{
								Type:     schema.TypeString,
								Computed: true,
							},
						},
						"redirect_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"custom_fragment": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"custom_host": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"custom_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"custom_query_string": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"redirect_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"redirect_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"forwarding_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"backend_pool_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cache_use_dynamic_compression": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"cache_query_parameter_strip_directive": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"custom_forwarding_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"forwarding_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"enforce_backend_pools_certificate_name_check": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"backend_pool_load_balancing": {
				Type:     schema.TypeList,
				MaxItems: 5000,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sample_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"successful_samples_required": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"additional_latency_milliseconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"backend_pool_health_probe": {
				Type:     schema.TypeList,
				MaxItems: 5000,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interval_in_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"backend_pool": {
				Type:     schema.TypeList,
				MaxItems: 50,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend": {
							Type:     schema.TypeList,
							MaxItems: 100,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"http_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"https_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"host_header": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_probe_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancing_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"frontend_endpoint": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"session_affinity_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"session_affinity_ttl_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"custom_https_provisioning_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"custom_https_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_source": {
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
										Computed: true,
									},
									"azure_key_vault_certificate_secret_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"azure_key_vault_certificate_vault_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func dataSourceArmFrontDoorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).frontdoor.FrontDoorsClient
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
			d.Set("load_balancer_enabled", true)
		} else {
			d.Set("load_balancer_enabled", false)
		}
		d.Set("friendly_name", properties.FriendlyName)

		frontDoorFrontendEndpoint, flattenErr := flattenArmFrontDoorFrontendEndpoint(properties.FrontendEndpoints, resourceGroup, *resp.Name, meta)

		if flattenErr == nil {
			if err := d.Set("frontend_endpoint", frontDoorFrontendEndpoint); err != nil {
				return fmt.Errorf("Error setting `frontend_endpoint`: %+v", err)
			}
		} else {
			return fmt.Errorf("Error setting `frontend_endpoint`: %+v", flattenErr)
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
