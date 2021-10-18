package resource

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"

	resource "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceResourceGroups() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceResourceGroupsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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
	subClient := armClient.Subscription.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// ListComplete returns an iterator struct
	var results resource.GroupListResultIterator
	var err error
	// iterate across each resource groups and append them to slice
	resourceGroups := make([]map[string]interface{}, 0)
	for _, subId := range d.Get("subscription_ids").([]string){
		rgSubGroupsClient := resource.NewGroupsClient(subId)
		results, err = rgSubGroupsClient.ListComplete(ctx, "", nil)
		if err != nil {
			return fmt.Errorf("listing resource groups: %+v", err)
		}

	for results.NotDone() {
		val := results.Value()

		rgId := ""
		if v := val.ID; v != nil {
			rgId = *v
		}

		rgSubId := ""
		rgStruct, err := parse.ResourceGroupID(*val.ID)
		if err != nil {
			return fmt.Errorf("parsing Resource Group ID")
		}
		rgSubId = rgStruct.SubscriptionId
		

		rgTenantId := ""
		resp, err := subClient.Get(ctx, rgSubId)
		if err != nil {
			return fmt.Errorf("reading subscription: %+v", err)
		} else {
			rgTenantId = *resp.TenantID
		}

		rgName := ""
		if v := val.Name; v != nil {
			rgName = *v
		}

		rgType := ""
		if v := val.Type; v != nil {
			rgType = *v
		}

		rgLocation := ""
		if v := val.Location; v != nil {
			rgLocation = *v
		}

		rgTags := make(map[string]interface{})
		if val.Tags != nil {
			rgTags = make(map[string]interface{}, len(val.Tags))
			for key, value := range val.Tags {
				if value != nil {
					rgTags[key] = *value
				}
			}
		}

		if err = results.Next(); err != nil {
			return fmt.Errorf("going to next resource groups value: %+v", err)
		}

		resourceGroups = append(resourceGroups, map[string]interface{}{
			"id": rgId,
			"name": rgName,
			"type":     rgType,
			"location": rgLocation,
			"subscriptionId": rgSubId,
			"tenantId": rgTenantId,
			"tags":     rgTags,
		})
	}
}

d.SetId("resource_groups-" + uuid.New().String())
if err = d.Set("resource_groups", resourceGroups); err != nil {
	return fmt.Errorf("setting `resource_groups`: %+v", err)
}

return nil
}