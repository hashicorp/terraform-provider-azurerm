package managementgroup

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/managementgroups"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceManagementGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceManagementGroupRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "group_id", "display_name"},
				Deprecated:   "Deprecated in favour of `name`",
				ValidateFunc: validate.ManagementGroupName,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "group_id", "display_name"},
				ValidateFunc: validate.ManagementGroupName,
			},

			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "group_id", "display_name"},
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

func dataSourceManagementGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	groupName := ""
	if v, ok := d.GetOk("name"); ok {
		groupName = v.(string)
	}
	if v, ok := d.GetOk("group_id"); ok {
		groupName = v.(string)
	}
	displayName := d.Get("display_name").(string)

	// one of displayName and groupName must be non-empty, this is guaranteed by schema
	// if the user is retrieving the mgmt group by display name, use the list api to get the group name first
	var err error
	if displayName != "" {
		groupName, err = getManagementGroupNameByDisplayName(ctx, client, displayName)
		if err != nil {
			return fmt.Errorf("Error reading Management Group (Display Name %q): %+v", displayName, err)
		}
	}
	recurse := true
	resp, err := client.Get(ctx, groupName, "children", &recurse, "", managementGroupCacheControl)
	if err != nil {
		if utils.ResponseWasForbidden(resp.Response) {
			return fmt.Errorf("Management Group %q was not found", groupName)
		}

		return fmt.Errorf("Error reading Management Group %q: %+v", groupName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Client returned an nil ID for Management Group %q", groupName)
	}

	d.SetId(*resp.ID)
	d.Set("name", groupName)
	d.Set("group_id", groupName)

	if props := resp.Properties; props != nil {
		d.Set("display_name", props.DisplayName)

		subscriptionIds, err := flattenManagementGroupDataSourceSubscriptionIds(props.Children)
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

func getManagementGroupNameByDisplayName(ctx context.Context, client *managementgroups.Client, displayName string) (string, error) {
	iterator, err := client.ListComplete(ctx, managementGroupCacheControl, "")
	if err != nil {
		return "", fmt.Errorf("Error listing Management Groups: %+v", err)
	}

	var results []string
	for iterator.NotDone() {
		group := iterator.Value()
		if group.DisplayName != nil && *group.DisplayName == displayName && group.Name != nil && *group.Name != "" {
			results = append(results, *group.Name)
		}

		if err := iterator.NextWithContext(ctx); err != nil {
			return "", fmt.Errorf("Error listing Management Groups: %+v", err)
		}
	}

	// we found none
	if len(results) == 0 {
		return "", fmt.Errorf("Management Group (Display Name %q) was not found", displayName)
	}

	// we found more than one
	if len(results) > 1 {
		return "", fmt.Errorf("expected a single Management Group with the Display Name %q but expected one", displayName)
	}

	return results[0], nil
}

func flattenManagementGroupDataSourceSubscriptionIds(input *[]managementgroups.ChildInfo) (*schema.Set, error) {
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
