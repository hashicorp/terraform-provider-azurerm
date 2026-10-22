// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/management/2020-05-01/managementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementGroupDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseManagementGroupID(d.Id())
	if err != nil {
		return err
	}

	group, err := client.Get(ctx, *id, managementgroups.GetOperationOptions{
		CacheControl: &managementGroupCacheControl,
		Filter:       pointer.To("children.childType eq Subscription"),
		Expand:       pointer.To(managementgroups.ExpandChildren),
		Recurse:      pointer.To(true),
	})
	if err != nil {
		if response.WasNotFound(group.HttpResponse) || response.WasForbidden(group.HttpResponse) {
			log.Printf("[DEBUG] Management Group %q doesn't exist in Azure - nothing to do!", id.GroupId)
			return nil
		}

		return fmt.Errorf("unable to retrieve Management Group %q: %+v", id.GroupId, err)
	}

	// before deleting a management group, return any subscriptions to the root management group
	if model := group.Model; model != nil {
		if props := model.Properties; props != nil {
			if children := props.Children; children != nil {
				for _, v := range *children {
					if v.Id == nil {
						continue
					}

					subscriptionId, err := commonids.ParseSubscriptionID(*v.Id)
					if err != nil {
						return fmt.Errorf("unable to parse child Subscription ID %+v", err)
					}
					managementGroupSubscriptionId := managementgroups.NewSubscriptionID(id.GroupId, subscriptionId.SubscriptionId)

					log.Printf("[DEBUG] De-associating Subscription %q from Management Group %q..", subscriptionId, id.GroupId)
					// NOTE: whilst this says `Delete` it's actually `Deassociate` - which is /really/ helpful
					deleteResp, err := client.SubscriptionsDelete(ctx, managementGroupSubscriptionId, managementgroups.SubscriptionsDeleteOperationOptions{
						CacheControl: &managementGroupCacheControl,
					})
					if err != nil {
						if !response.WasNotFound(deleteResp.HttpResponse) {
							return fmt.Errorf("unable to de-associate Subscription %q from Management Group %q: %+v", subscriptionId.SubscriptionId, id.GroupId, err)
						}
					}
				}
			}
		}
	}

	if err = client.DeleteThenPoll(ctx, *id, managementgroups.DeleteOperationOptions{
		CacheControl: &managementGroupCacheControl,
	}); err != nil {
		return fmt.Errorf("unable to delete Management Group %q: %+v", id.GroupId, err)
	}

	return nil
}
