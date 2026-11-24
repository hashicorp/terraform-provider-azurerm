// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_subscription_association

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/management/2020-05-01/managementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementGroupSubscriptionAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	// There is no "read" function on the appropriate client so we need to check if the Subscription is in the Management Group subscription list
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managementgroups.ParseSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	managementGroup, err := client.Get(ctx, commonids.NewManagementGroupID(id.GroupId), managementgroups.GetOperationOptions{
		CacheControl: &managementGroupCacheControl,
		Expand:       pointer.To(managementgroups.ExpandChildren),
		Recurse:      pointer.To(false),
	})
	if err != nil {
		return fmt.Errorf("reading Management Group %q for Subscription Associations: %+v", id.GroupId, err)
	}
	found := false
	if model := managementGroup.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.Children != nil {
				for _, v := range *props.Children {
					if v.Type != nil && *v.Type == managementgroups.ManagementGroupChildTypeSubscriptions && v.Name != nil && strings.EqualFold(*v.Name, id.SubscriptionId) {
						found = true
					}
				}
			}

			if !found {
				log.Printf("[INFO] Subscription %q not found in Management group %q, removing from state", id.SubscriptionId, id.GroupId)
				d.SetId("")
				return nil
			}

			managementGroupId := parse.NewManagementGroupId(id.GroupId)
			d.Set("management_group_id", managementGroupId.ID())
			subscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)
			d.Set("subscription_id", subscriptionId.ID())
		}
	}

	return nil
}

func subscriptionAssociationRefreshFunc(ctx context.Context, client *managementgroups.ManagementGroupsClient, id managementgroups.SubscriptionId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		managementGroup, err := client.Get(ctx, commonids.NewManagementGroupID(id.GroupId), managementgroups.GetOperationOptions{
			CacheControl: &managementGroupCacheControl,
			Expand:       pointer.To(managementgroups.ExpandChildren),
			Recurse:      pointer.To(false),
		})
		if err != nil {
			return nil, "", fmt.Errorf("reading Management Group %q for Subscription Associations: %+v", id.GroupId, err)
		}

		if model := managementGroup.Model; model != nil {
			if props := model.Properties; props != nil && props.Children != nil {
				for _, v := range *props.Children {
					if v.Type != nil && *v.Type == managementgroups.ManagementGroupChildTypeSubscriptions {
						if v.Name != nil && strings.EqualFold(*v.Name, id.SubscriptionId) {
							return managementGroup, "Exists", nil
						}
					}
				}
			}
		}

		return "NotFound", "NotFound", nil
	}
}
