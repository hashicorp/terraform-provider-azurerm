package maintenance

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceMaintenanceConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmMaintenanceConfigurationRead,

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

			"location": azure.SchemaLocationForDataSource(),

			"scope": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"visibility": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"window": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"start_date_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"expiration_date_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"duration": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"time_zone": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"recur_every": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"properties": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmMaintenanceConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("maintenance Configuration %q was not found in Resource Group %q", name, resGroup)
		}
		return fmt.Errorf("retrieving Maintenance Configuration %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.ConfigurationProperties; props != nil {
		d.Set("scope", props.MaintenanceScope)
		d.Set("visibility", props.Visibility)
		d.Set("properties", props.ExtensionProperties)

		window := flattenMaintenanceConfigurationWindow(props.Window)
		if err := d.Set("window", window); err != nil {
			return fmt.Errorf("error setting `window`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
