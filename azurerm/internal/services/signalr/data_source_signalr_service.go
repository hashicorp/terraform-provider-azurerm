package signalr

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSignalRService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSignalRServiceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"server_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"primary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmSignalRServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("SignalR %q was not found in Resource Group %q", name, resourceGroup)
		}
		return fmt.Errorf("Error getting SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	keys, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error getting keys of SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if properties := resp.Properties; properties != nil {
		d.Set("hostname", properties.HostName)
		d.Set("ip_address", properties.ExternalIP)
		d.Set("public_port", properties.PublicPort)
		d.Set("server_port", properties.ServerPort)
	}

	d.Set("primary_access_key", keys.PrimaryKey)
	d.Set("primary_connection_string", keys.PrimaryConnectionString)
	d.Set("secondary_access_key", keys.SecondaryKey)
	d.Set("secondary_connection_string", keys.SecondaryConnectionString)

	return tags.FlattenAndSet(d, resp.Tags)
}
