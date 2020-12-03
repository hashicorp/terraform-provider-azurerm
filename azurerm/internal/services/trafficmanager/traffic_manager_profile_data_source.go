package trafficmanager

import (
	"fmt"
	"time"


	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmTrafficManagerProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmTrafficManagerProfileRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"profile_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"traffic_routing_method": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dns_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"relative_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"monitor_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expected_status_code_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"custom_header": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"interval_in_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"timeout_in_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"tolerated_number_of_failures": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmTrafficManagerProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Traffic Manager Profile %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Traffic Manager Profile %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if loc := resp.Location; loc != nil {
		d.Set("location", location.NormalizeNilable(loc))
	}

	if profile := resp.ProfileProperties; profile != nil {
		d.Set("profile_status", profile.ProfileStatus)
		d.Set("traffic_routing_method", profile.TrafficRoutingMethod)

		d.Set("dns_config", flattenAzureRMTrafficManagerProfileDNSConfig(profile.DNSConfig))
		d.Set("monitor_config", flattenAzureRMTrafficManagerProfileMonitorConfig(profile.MonitorConfig))

		// fqdn is actually inside DNSConfig, inlined for simpler reference
		if dns := profile.DNSConfig; dns != nil {
			d.Set("fqdn", dns.Fqdn)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}
