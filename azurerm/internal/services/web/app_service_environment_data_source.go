package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceAppServiceEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppServiceEnvironmentRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"front_end_scale_factor": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"internal_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"outbound_ip_addresses": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},

			"pricing_tier": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceAppServiceEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
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
	}

	d.Set("internal_ip_address", vipInfo.InternalIPAddress)
	d.Set("service_ip_address", vipInfo.ServiceIPAddress)
	d.Set("outbound_ip_addresses", vipInfo.OutboundIPAddresses)

	return tags.FlattenAndSet(d, resp.Tags)
}
