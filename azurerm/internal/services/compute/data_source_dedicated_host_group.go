package compute

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDedicatedHostGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDedicatedHostGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^_\W][\w-.]{0,78}[\w]$`), ""),
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"host_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"platform_fault_domain_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmDedicatedHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Dedicated Host Group %q (Resource Group %q) was not found", name, resourceGroupName)
		}
		return fmt.Errorf("Error reading Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.DedicatedHostGroupProperties; props != nil {
		if props.Hosts != nil {
			var hostIds []string
			for _, item := range *props.Hosts {
				hostIds = append(hostIds, *item.ID)
			}
			if err := d.Set("host_ids", hostIds); err != nil {
				return fmt.Errorf("Error setting `hosts`: %+v", err)
			}
		}
		d.Set("platform_fault_domain_count", int(*props.PlatformFaultDomainCount))
	}

	d.Set("zones", utils.FlattenStringSlice(resp.Zones))

	return nil
}
