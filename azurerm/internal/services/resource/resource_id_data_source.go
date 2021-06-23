package resource

import (
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func dataSourceResourceId() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceResourceIdRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Second),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_id": {
				Type:        pluginsdk.TypeString,
				Required:    true,
				Description: "Resource id to parse",
			},
			"subscription_id": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "Resource subscription id",
			},
			"resource_group_name": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "Resource group name",
			},
			"provider_namespace": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "Resource namespace",
			},
			"resource_type": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "Resource type",
			},
			"name": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "Resource name",
			},
			"parent_resources": {
				Type:        pluginsdk.TypeMap,
				Computed:    true,
				Description: "A map of parent resource types and names",
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"full_resource_type": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "Full resource type (including parents if applicable)",
			},
		},
	}
}

func dataSourceResourceIdRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id := d.Get("resource_id").(string)
	splits := strings.Split(strings.Trim(id, "/"), "/")

	count := len(splits)
	formatErr := fmt.Errorf("The specified ID %v is not a valid Azure resource ID.", id)

	if count%2 == 1 {
		return formatErr
	}

	//   Format of id:
	// 	   /
	//   0 subscriptions/
	//   1 subscriptionId/
	//   2 resourceGroups/
	//   3 resourceGroupName/
	//   4 providers/
	//   5 providerNamespace/
	//  (6 parentResourceType/)*
	//  (7 parentName/)*
	//  ^1 resourceType/
	//  ^0 name

	if count < 2 {
		return formatErr
	}

	if !strings.EqualFold(splits[0], "subscriptions") {
		return formatErr
	}

	d.Set("subscription_id", splits[1])

	if count >= 4 {
		if !strings.EqualFold(splits[2], "resourceGroups") {
			return formatErr
		}
		d.Set("resource_group_name", splits[3])
	}

	fullResourceType := ""
	if count >= 6 {
		if !strings.EqualFold(splits[4], "providers") {
			return formatErr
		}
		d.Set("provider_namespace", splits[5])
		fullResourceType = splits[5]
	}

	parentResources := make(map[string]string)

	for i := 8; i <= count; i += 2 {
		resType := splits[i-2]
		resName := splits[i-1]
		fullResourceType += fmt.Sprintf("/%v", resType)
		if i == count {
			d.Set("resource_type", resType)
			d.Set("name", resName)
		} else {
			parentResources[resType] = resName
		}
	}

	d.Set("full_resource_type", fullResourceType)
	d.Set("parent_resources", parentResources)

	d.SetId(id)
	return nil
}
