// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group

import (
	"fmt"
	"log"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/management/2020-05-01/managementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	accountClient := meta.(*clients.Client)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseManagementGroupID(d.Id())
	if err != nil {
		return err
	}

	tenantID := accountClient.Account.TenantId
	tenantScopedID := parse.NewTenantScopedManagementGroupID(tenantID, id.GroupId)
	d.Set("tenant_scoped_id", tenantScopedID.TenantScopedID())

	recurse := pointer.To(true)
	resp, err := client.Get(ctx, *id, managementgroups.GetOperationOptions{
		CacheControl: &managementGroupCacheControl,
		Filter:       pointer.To("children.childType eq Subscription"),
		Expand:       pointer.To(managementgroups.ExpandChildren),
		Recurse:      recurse,
	})
	if err != nil {
		if response.WasForbidden(resp.HttpResponse) || response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Management Group %q doesn't exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("unable to read Management Group %q: %+v", d.Id(), err)
	}

	d.Set("name", id.GroupId)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("display_name", props.DisplayName)

			subscriptionIds, err := flattenManagementGroupSubscriptionIds(props.Children)
			if err != nil {
				return fmt.Errorf("unable to flatten `subscription_ids`: %+v", err)
			}
			d.Set("subscription_ids", subscriptionIds)

			parentId := ""
			if details := props.Details; details != nil {
				if parent := details.Parent; parent != nil {
					if pid := parent.Id; pid != nil {
						parentId = *pid
					}
				}
			}
			d.Set("parent_management_group_id", parentId)
		}
	}

	return nil
}

func flattenManagementGroupSubscriptionIds(input *[]managementgroups.ManagementGroupChildInfo) (*pluginsdk.Set, error) {
	subscriptionIds := &pluginsdk.Set{F: pluginsdk.HashString}
	if input == nil {
		return subscriptionIds, nil
	}

	for _, child := range *input {
		if child.Type == nil || *child.Type != managementgroups.ManagementGroupChildTypeSubscriptions || child.Id == nil {
			continue
		}

		id, err := commonids.ParseSubscriptionID(*child.Id)
		if err != nil {
			return nil, fmt.Errorf("unable to parse child Subscription ID %+v", err)
		}

		if id != nil {
			subscriptionIds.Add(id.SubscriptionId)
		}
	}

	return subscriptionIds, nil
}

var managementGroupCacheControl = "no-cache"

func determineManagementGroupSubscriptionsIdsToRemove(existing *[]managementgroups.ManagementGroupChildInfo, updated []string) (*[]string, error) {
	subscriptionIdsToRemove := make([]string, 0)
	if existing == nil {
		return &subscriptionIdsToRemove, nil
	}

	for _, v := range *existing {
		if v.Type == nil || *v.Type != managementgroups.ManagementGroupChildTypeSubscriptions || v.Id == nil {
			continue
		}

		id, err := commonids.ParseSubscriptionID(*v.Id)
		if err != nil {
			return nil, fmt.Errorf("unable to parse Subscription ID %q: %+v", *v.Id, err)
		}

		// not a Subscription - so let's skip it
		if id == nil {
			continue
		}

		found := slices.Contains(updated, id.SubscriptionId)

		if !found {
			subscriptionIdsToRemove = append(subscriptionIdsToRemove, id.SubscriptionId)
		}
	}

	return &subscriptionIdsToRemove, nil
}
