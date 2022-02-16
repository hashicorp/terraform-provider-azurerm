package compute

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceAvailabilitySet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAvailabilitySetRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"platform_update_domain_count": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"platform_fault_domain_count": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"managed": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceAvailabilitySetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.AvailabilitySetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAvailabilitySetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if resp.Sku != nil && resp.Sku.Name != nil {
		d.Set("managed", strings.EqualFold(*resp.Sku.Name, "Aligned"))
	}
	if props := resp.AvailabilitySetProperties; props != nil {
		if v := props.PlatformUpdateDomainCount; v != nil {
			d.Set("platform_update_domain_count", strconv.Itoa(int(*v)))
		}
		if v := props.PlatformFaultDomainCount; v != nil {
			d.Set("platform_fault_domain_count", strconv.Itoa(int(*v)))
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}
