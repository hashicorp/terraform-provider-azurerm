package resource

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceResourceGroups() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceResourceGroupsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"filter_by_subscription_id": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"resource_groups": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"location": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"subscription_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tags": tags.SchemaDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceResourceGroupsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	armClient := meta.(*clients.Client)
	rgClient := armClient.Resource.GroupsClient
	subClient := armClient.Subscription.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// ListComplete returns an iterator struct
	results, err := rgClient.ListComplete(ctx, "", nil)
	if err != nil {
		return fmt.Errorf("listing resource groups: %+v", err)
	}

	filterBySubscriptionId := []string(nil)
	if v, ok := d.GetOk("filter_by_subscription_id"); ok {
		for _, p := range v.([]interface{}) {
			filterBySubscriptionId = append(filterBySubscriptionId, p.(string))
		}
	}

	// iterate across each resource groups and append them to slice
	resourceGroups := make([]map[string]interface{}, 0)
	subTenantIdMap := make(map[string]string)
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

		if val, ok := subTenantIdMap[rg["subscription_id"].(string)]; ok {
			rg["tenant_id"] = val
		} else {
			resp, err := subClient.Get(ctx, rg["subscription_id"].(string))

			if err != nil {
				return fmt.Errorf("reading subscription: %+v", err)
			} else {
				rg["tenant_id"] = *resp.TenantID
				subTenantIdMap[rg["subscription_id"].(string)] = rg["tenant_id"].(string)
			}
		}

		filterCondition := contains(filterBySubscriptionId, rg["subscription_id"].(string))
		if filterBySubscriptionId == nil || filterCondition {
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

			resourceGroups = append(resourceGroups, rg)
		}
	}

	d.SetId("resource_groups-" + armClient.Account.TenantId)
	if err = d.Set("resource_groups", resourceGroups); err != nil {
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
