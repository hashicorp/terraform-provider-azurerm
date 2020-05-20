package network

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

func dataSourceArmFirewall() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmFirewallRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateAzureFirewallName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"ip_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_public_ip_address_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip_address_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.AzureFirewallsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("Firewall %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*read.ID)
	d.Set("name", read.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := read.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := read.AzureFirewallPropertiesFormat; props != nil {
		if err := d.Set("ip_configuration", flattenArmFirewallIPConfigurations(props.IPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, read.Tags)
}
