package trafficmanager

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

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

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourceArmTrafficManagerProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := profiles.NewTrafficManagerProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		if profile := model.Properties; profile != nil {
			profileStatus := ""
			if profile.ProfileStatus != nil {
				profileStatus = string(*profile.ProfileStatus)
			}
			d.Set("profile_status", profileStatus)
			trafficRoutingMethod := ""
			if profile.TrafficRoutingMethod != nil {
				trafficRoutingMethod = string(*profile.TrafficRoutingMethod)
			}
			d.Set("traffic_routing_method", trafficRoutingMethod)

			d.Set("dns_config", flattenAzureRMTrafficManagerProfileDNSConfig(profile.DnsConfig))
			d.Set("monitor_config", flattenAzureRMTrafficManagerProfileMonitorConfig(profile.MonitorConfig))

			trafficViewEnabled := false
			if profile.TrafficViewEnrollmentStatus != nil {
				trafficViewEnabled = *profile.TrafficViewEnrollmentStatus == profiles.TrafficViewEnrollmentStatusEnabled
			}
			d.Set("traffic_view_enabled", trafficViewEnabled)

			// fqdn is actually inside DNSConfig, inlined for simpler reference
			if dns := profile.DnsConfig; dns != nil {
				d.Set("fqdn", dns.Fqdn)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}
