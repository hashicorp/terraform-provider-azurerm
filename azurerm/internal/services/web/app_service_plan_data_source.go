package web

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceAppServicePlan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: AppServicePlanDataSourceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

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

			"tags": tags.SchemaDataSource(),
		},
	}
}

func AppServicePlanDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicePlansClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on App Service Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if utils.ResponseWasNotFound(resp.Response) {
		return fmt.Errorf("Error: App Service Plan %q (Resource Group %q) was not found", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("kind", resp.Kind)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

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
	}

	if err := d.Set("sku", flattenAppServicePlanSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
