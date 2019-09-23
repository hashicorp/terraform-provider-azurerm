package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/managementgroups"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmManagementGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmManagementGroupRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"parent_management_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subscription_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceArmManagementGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).managementGroups.GroupsClient
	ctx := meta.(*ArmClient).StopContext

	groupId := d.Get("group_id").(string)

	recurse := true
	resp, err := client.Get(ctx, groupId, "children", &recurse, "", managementGroupCacheControl)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Management Group %q was not found", groupId)
		}

		return fmt.Errorf("Error reading Management Group %q: %+v", d.Id(), err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Client returned an nil ID for Management Group %q", groupId)
	}

	d.SetId(*resp.ID)
	d.Set("group_id", groupId)

	if props := resp.Properties; props != nil {
		d.Set("display_name", props.DisplayName)

		subscriptionIds, err := flattenArmManagementGroupDataSourceSubscriptionIds(props.Children)
		if err != nil {
			return fmt.Errorf("Error flattening `subscription_ids`: %+v", err)
		}
		d.Set("subscription_ids", subscriptionIds)

		parentId := ""
		if details := props.Details; details != nil {
			if parent := details.Parent; parent != nil {
				if pid := parent.ID; pid != nil {
					parentId = *pid
				}
			}
		}
		d.Set("parent_management_group_id", parentId)

	}

	return nil
}

func flattenArmManagementGroupDataSourceSubscriptionIds(input *[]managementgroups.ChildInfo) (*schema.Set, error) {
	subscriptionIds := &schema.Set{F: schema.HashString}
	if input == nil {
		return subscriptionIds, nil
	}

	for _, child := range *input {
		if child.ID == nil {
			continue
		}

		id, err := parseManagementGroupSubscriptionID(*child.ID)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse child Subscription ID %+v", err)
		}

		if id != nil {
			subscriptionIds.Add(id.subscriptionId)
		}
	}

	return subscriptionIds, nil
}
