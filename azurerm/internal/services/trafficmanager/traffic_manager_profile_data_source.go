package trafficmanager

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-08-01/trafficmanager"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/trafficmanager/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmTrafficManagerProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmTrafficManagerProfileRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"profile_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"traffic_routing_method": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"dns_config": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"relative_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"monitor_config": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"expected_status_code_ranges": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"custom_header": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"value": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},

						"protocol": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"port": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"path": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"interval_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"timeout_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"tolerated_number_of_failures": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"traffic_view_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmTrafficManagerProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewTrafficManagerProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Traffic Manager Profile %q was not found in Resource Group %q", id.Name, id.ResourceGroup)
		}

		return fmt.Errorf("retrieving Traffic Manager Profile %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.SetId(id.ID())

	if profile := resp.ProfileProperties; profile != nil {
		d.Set("profile_status", profile.ProfileStatus)
		d.Set("traffic_routing_method", profile.TrafficRoutingMethod)

		d.Set("dns_config", flattenAzureRMTrafficManagerProfileDNSConfig(profile.DNSConfig))
		d.Set("monitor_config", flattenAzureRMTrafficManagerProfileMonitorConfig(profile.MonitorConfig))
		d.Set("traffic_view_enabled", profile.TrafficViewEnrollmentStatus == trafficmanager.TrafficViewEnrollmentStatusEnabled)

		// fqdn is actually inside DNSConfig, inlined for simpler reference
		if dns := profile.DNSConfig; dns != nil {
			d.Set("fqdn", dns.Fqdn)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}
