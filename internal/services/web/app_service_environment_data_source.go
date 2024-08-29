// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceAppServiceEnvironment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read:               dataSourceAppServiceEnvironmentRead,
		DeprecationMessage: "This data source is deprecated due to the [retirement of v1 and v2 App Service Environments](https://azure.microsoft.com/en-gb/updates/app-service-environment-v1-and-v2-retirement-announcement/) and will be removed inv4.0 of the provider. Please use `azurerm_app_service_environment_v3` instead.",
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"cluster_setting": {
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

			"front_end_scale_factor": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"internal_ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"outbound_ip_addresses": {
				Type: pluginsdk.TypeList,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
				Computed: true,
			},

			"pricing_tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceAppServiceEnvironmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAppServiceEnvironmentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.HostingEnvironmentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	vipInfo, err := client.GetVipInfo(ctx, id.ResourceGroup, id.HostingEnvironmentName)
	if err != nil {
		if utils.ResponseWasNotFound(vipInfo.Response) {
			return fmt.Errorf("retrieving VIP info: %s was not found", id)
		}
		return fmt.Errorf("retrieving VIP info %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.HostingEnvironmentName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.AppServiceEnvironment; props != nil {
		frontendScaleFactor := 0
		if props.FrontEndScaleFactor != nil {
			frontendScaleFactor = int(*props.FrontEndScaleFactor)
		}
		d.Set("front_end_scale_factor", frontendScaleFactor)

		pricingTier := ""
		if props.MultiSize != nil {
			pricingTier = convertToIsolatedSKU(*props.MultiSize)
		}
		d.Set("pricing_tier", pricingTier)
		d.Set("cluster_setting", flattenClusterSettings(props.ClusterSettings))
	}

	d.Set("internal_ip_address", vipInfo.InternalIPAddress)
	d.Set("service_ip_address", vipInfo.ServiceIPAddress)
	d.Set("outbound_ip_addresses", vipInfo.OutboundIPAddresses)

	return tags.FlattenAndSet(d, resp.Tags)
}
