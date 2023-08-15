// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceCdnFrontDoorOriginGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCdnFrontDoorOriginGroupRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorOriginGroupName,
			},

			"profile_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			// Computed
			"cdn_frontdoor_profile_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"health_probe": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"interval_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"path": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"request_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"load_balancing": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"additional_latency_in_milliseconds": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"sample_size": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"successful_samples_required": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"restore_traffic_time_to_healed_or_new_endpoint_in_minutes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"session_affinity_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceCdnFrontDoorOriginGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFrontDoorOriginGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.OriginGroupName)
	d.Set("profile_name", id.ProfileName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontDoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.AFDOriginGroupProperties; props != nil {
		if err := d.Set("health_probe", flattenCdnFrontDoorOriginGroupHealthProbeParameters(props.HealthProbeSettings)); err != nil {
			return fmt.Errorf("setting 'health_probe': %+v", err)
		}

		if err := d.Set("load_balancing", flattenCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(props.LoadBalancingSettings)); err != nil {
			return fmt.Errorf("setting 'load_balancing': %+v", err)
		}

		d.Set("session_affinity_enabled", flattenEnabledBool(props.SessionAffinityState))
		d.Set("restore_traffic_time_to_healed_or_new_endpoint_in_minutes", props.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes)
	}

	return nil
}
