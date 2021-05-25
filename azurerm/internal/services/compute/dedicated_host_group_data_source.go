package compute

import (
	"fmt"
	"regexp"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceDedicatedHostGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDedicatedHostGroupRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^_\W][\w-.]{0,78}[\w]$`), ""),
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"platform_fault_domain_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"automatic_placement_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"zones": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceDedicatedHostGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroupName, name, "")
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
		platformFaultDomainCount := 0
		if props.PlatformFaultDomainCount != nil {
			platformFaultDomainCount = int(*props.PlatformFaultDomainCount)
		}
		d.Set("platform_fault_domain_count", platformFaultDomainCount)

		d.Set("automatic_placement_enabled", props.SupportAutomaticPlacement)
	}

	d.Set("zones", utils.FlattenStringSlice(resp.Zones))

	return tags.FlattenAndSet(d, resp.Tags)
}
