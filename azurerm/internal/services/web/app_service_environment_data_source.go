package web

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceAppServiceEnvironment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAppServiceEnvironmentRead,

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
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: App Service Environment %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error retrieving App Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	vipInfo, err := client.GetVipInfo(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(vipInfo.Response) {
			return fmt.Errorf("Error retrieving VIP info: App Service Environment %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error retrieving VIP info App Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
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
