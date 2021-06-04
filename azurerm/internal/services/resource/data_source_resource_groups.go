package resource

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
)

func dataSourceResourceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceResourceGroupsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tags.SchemaDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceResourceGroupsRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*clients.Client)
	rgClient := armClient.Resource.GroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// ListComplete returns an iterator struct
	results, err := rgClient.ListComplete(ctx, "", nil)
	if err != nil {
		return fmt.Errorf("listing resource groups: %+v", err)
	}

	// iterate across each resource groups and append them to slice
	resource_groups := make([]map[string]interface{}, 0)
	for results.NotDone() {
		val := results.Value()

		rg := make(map[string]interface{})

		if v := val.ID; v != nil {
			rg["id"] = *v
			rg["subscription_id"] = parse.ResourceGroupID.SubscriptionId
		}
		if v := val.Name; v != nil {
			rg["name"] = *v
		}
		if v := val.Type; v != nil {
			rg["type"] = *v
		}
		if v := val.Location; v != nil {
			rg["location"] = *v
		}

		if err = results.Next(); err != nil {
			return fmt.Errorf("going to next resource groups value: %+v", err)
		}

		rg["tags"] = tags.Flatten(val.Tags)

		resource_groups = append(resource_groups, rg)
	}

	d.SetId("resource_groups-" + armClient.Account.TenantId)
	if err = d.Set("resource_groups", resource_groups); err != nil {
		return fmt.Errorf("setting `resource_groups`: %+v", err)
	}

	return nil
}
