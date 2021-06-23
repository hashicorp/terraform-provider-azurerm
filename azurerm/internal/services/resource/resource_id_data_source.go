package resource

import (
	"errors"
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
			"parent_resource_type": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "Resource type",
			},
			"parent_name": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "Resource name",
			},
			"full_resource_type": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "Full resource type (including parent if applicable)",
			},
		},
	}
}

func badIdError(id string) string {
	return fmt.Sprintf("The specified ID %v is not a valid Azure resource ID.", id)
}

func dataSourceResourceIdRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id := d.Get("resource_id").(string)
	splits := strings.Split(strings.Trim(id, "/"), "/")

	count := len(splits)
	err := errors.New(badIdError(id))

	if count%2 == 1 {
		return err
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
		return err
	}

	if !strings.EqualFold(splits[0], "subscriptions") {
		return err
	}

	d.Set("subscription_id", splits[1])

	if count >= 4 {
		if !strings.EqualFold(splits[2], "resourceGroups") {
			return err
		}
		d.Set("resource_group_name", splits[3])
	}

	if count >= 6 {
		if !strings.EqualFold(splits[4], "providers") {
			return err
		}
		d.Set("provider_namespace", splits[5])
	}

	if count >= 8 {
		d.Set("resource_type", splits[count-2])
		d.Set("name", splits[count-1])
		d.Set("full_resource_type", fmt.Sprintf("%v/%v", splits[5], splits[count-2]))
	}

	if count == 10 {
		d.Set("parent_resource_type", splits[6])
		d.Set("parent_name", splits[7])
		d.Set("full_resource_type", fmt.Sprintf("%v/%v/%v", splits[5], splits[6], splits[count-2]))
	}

	d.SetId(id)
	return nil
}
