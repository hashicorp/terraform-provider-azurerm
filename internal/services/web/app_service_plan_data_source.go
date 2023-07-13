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

func dataSourceAppServicePlan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: AppServicePlanDataSourceRead,

		DeprecationMessage: "The `azurerm_app_service_plan` data source has been superseded by the `azurerm_service_plan` data source. Whilst this resource will continue to be available in the 2.x and 3.x releases it is feature-frozen for compatibility purposes, will no longer receive any updates and will be removed in a future major release of the Azure Provider.",

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

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"tier": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"size": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"capacity": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"app_service_environment_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"reserved": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"per_site_scaling": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"maximum_number_of_workers": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"maximum_elastic_worker_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"is_xenon": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"zone_redundant": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func AppServicePlanDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicePlansClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAppServicePlanID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerfarmName)
	if err != nil {
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	if utils.ResponseWasNotFound(resp.Response) {
		return fmt.Errorf("%s was not found", id)
	}

	d.SetId(id.ID())

	d.Set("name", id.ServerfarmName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("kind", resp.Kind)

	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.AppServicePlanProperties; props != nil {
		if profile := props.HostingEnvironmentProfile; profile != nil {
			d.Set("app_service_environment_id", profile.ID)
		}
		d.Set("per_site_scaling", props.PerSiteScaling)
		d.Set("reserved", props.Reserved)

		if props.MaximumNumberOfWorkers != nil {
			d.Set("maximum_number_of_workers", int(*props.MaximumNumberOfWorkers))
		}

		if props.MaximumElasticWorkerCount != nil {
			d.Set("maximum_elastic_worker_count", int(*props.MaximumElasticWorkerCount))
		}

		d.Set("is_xenon", props.IsXenon)
		d.Set("zone_redundant", props.ZoneRedundant)
	}

	if err := d.Set("sku", flattenAppServicePlanSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
