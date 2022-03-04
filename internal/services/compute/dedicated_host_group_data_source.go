package compute

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"platform_fault_domain_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"automatic_placement_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),

			"zones": commonschema.ZonesMultipleComputed(),
		},
	}
}

func dataSourceDedicatedHostGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDedicatedHostGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.HostGroupName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.HostGroupName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("zones", zones.Flatten(resp.Zones))

	if props := resp.DedicatedHostGroupProperties; props != nil {
		d.Set("automatic_placement_enabled", props.SupportAutomaticPlacement)

		platformFaultDomainCount := 0
		if props.PlatformFaultDomainCount != nil {
			platformFaultDomainCount = int(*props.PlatformFaultDomainCount)
		}
		d.Set("platform_fault_domain_count", platformFaultDomainCount)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
