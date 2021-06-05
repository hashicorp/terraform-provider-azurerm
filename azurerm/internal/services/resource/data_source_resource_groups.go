package resource

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
						"subscription_id_filter": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
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

	subscription_id_filter := []string(nil)
	if v, ok := d.GetOk("subscription_id_filter"); ok {
		subscription_id_filter = v.([]string)
	}

	// iterate across each resource groups and append them to slice
	resource_groups := make([]map[string]interface{}, 0)

	for results.NotDone() {
		val := results.Value()

		rg := make(map[string]interface{})

		if v := val.ID; v != nil {
			rg["id"] = *v
			rgStruct, err := parse.ResourceGroupID(*v)
			if err != nil {
				return fmt.Errorf("parsing Resource Group ID")
			}
			rg["subscription_id"] = rgStruct.SubscriptionId
		}

		if subscription_id_filter == nil || contains(subscription_id_filter, rg["subscription_id"].(string)) {
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
	}

	d.SetId("resource_groups-" + armClient.Account.TenantId)
	if err = d.Set("resource_groups", resource_groups); err != nil {
		return fmt.Errorf("setting `resource_groups`: %+v", err)
	}

	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
