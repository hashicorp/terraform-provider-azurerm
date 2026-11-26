// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/management/2020-05-01/managementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementGroupCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	accountClient := meta.(*clients.Client)
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	armTenantID := meta.(*clients.Client).Account.TenantId

	var groupName string
	if v := d.Get("name"); v != "" {
		groupName = v.(string)
	}

	if groupName == "" {
		groupName = uuid.New().String()
	}

	id := commonids.NewManagementGroupID(groupName)

	tenantID := accountClient.Account.TenantId
	tenantScopedID := parse.NewTenantScopedManagementGroupID(tenantID, id.GroupId)
	d.Set("tenant_scoped_id", tenantScopedID.TenantScopedID())

	parentManagementGroupId := d.Get("parent_management_group_id").(string)
	if parentManagementGroupId == "" {
		parentManagementGroupId = fmt.Sprintf("/providers/Microsoft.Management/managementGroups/%s", armTenantID)
	}

	recurse := false
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, managementgroups.GetOperationOptions{
				CacheControl: &managementGroupCacheControl,
				Expand:       pointer.To(managementgroups.ExpandChildren),
				Recurse:      &recurse,
			})
			if err != nil {
				// 403 is returned if group does not exist, bug tracked at: https://github.com/Azure/azure-rest-api-specs/issues/9549
				if !response.WasNotFound(existing.HttpResponse) && !response.WasForbidden(existing.HttpResponse) {
					return fmt.Errorf("unable to check for presence of existing Management Group %q: %s", groupName, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) && !response.WasForbidden(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_management_group", id.ID())
			}
		}
	}

	properties := managementgroups.CreateManagementGroupRequest{
		Name: pointer.To(groupName),
		Properties: &managementgroups.CreateManagementGroupProperties{
			TenantId: pointer.To(armTenantID),
			Details: &managementgroups.CreateManagementGroupDetails{
				Parent: &managementgroups.CreateParentGroupInfo{
					Id: pointer.To(parentManagementGroupId),
				},
			},
		},
	}

	if v := d.Get("display_name"); v != "" {
		properties.Properties.DisplayName = pointer.To(v.(string))
	}

	opts := managementgroups.CreateOrUpdateOperationOptions{
		CacheControl: &managementGroupCacheControl,
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, opts, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, properties, opts); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	// We have a potential race condition / consistency issue whereby the implicit role assignment for the SP may not be
	// completed before the read-back here or an eventually consistent read is creating a temporary 403 error.
	poller := custompollers.NewEventualConsistencyPoller(5, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, id, managementgroups.GetOperationOptions{
			CacheControl: &managementGroupCacheControl,
			Expand:       pointer.To(managementgroups.ExpandChildren),
			Recurse:      pointer.To(true),
		})
		return resp.HttpResponse, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:              10 * time.Second,
		RetryErrorStatusCodes: []int{http.StatusForbidden},
	})
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("failed waiting for read on Managementgroup %q", groupName)
	}

	resp, err := client.Get(ctx, id, managementgroups.GetOperationOptions{
		CacheControl: &managementGroupCacheControl,
		Expand:       pointer.To(managementgroups.ExpandChildren),
		Filter:       pointer.To("children.childType eq Subscription"),
		Recurse:      &recurse,
	})
	if err != nil {
		return fmt.Errorf("unable to retrieve Management Group %q: %+v", groupName, err)
	}

	subscriptionIds := expandManagementGroupSubscriptionIds(d.Get("subscription_ids").(*pluginsdk.Set))

	// first remove any which need to be removed
	if !d.IsNewResource() {
		log.Printf("[DEBUG] Determine which Subscriptions should be removed from Management Group %q", groupName)
		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				subscriptionIdsToRemove, err := determineManagementGroupSubscriptionsIdsToRemove(props.Children, subscriptionIds)
				if err != nil {
					return fmt.Errorf("unable to determine which subscriptions should be removed from Management Group %q: %+v", groupName, err)
				}

				for _, subscriptionId := range *subscriptionIdsToRemove {
					log.Printf("[DEBUG] De-associating Subscription ID %q from Management Group %q", subscriptionId, groupName)
					deleteResp, err := client.SubscriptionsDelete(ctx, managementgroups.NewSubscriptionID(groupName, subscriptionId), managementgroups.SubscriptionsDeleteOperationOptions{
						CacheControl: &managementGroupCacheControl,
					})
					if err != nil {
						if !response.WasNotFound(deleteResp.HttpResponse) {
							return fmt.Errorf("unable to de-associate Subscription %q from Management Group %q: %+v", subscriptionId, groupName, err)
						}
					}
				}
			}
		}
	}

	// then add the new ones
	for _, subscriptionId := range subscriptionIds {
		log.Printf("[DEBUG] Assigning Subscription ID %q to management group %q", subscriptionId, groupName)
		if _, err := client.SubscriptionsCreate(ctx, managementgroups.NewSubscriptionID(groupName, subscriptionId), managementgroups.SubscriptionsCreateOperationOptions{
			CacheControl: &managementGroupCacheControl,
		}); err != nil {
			return fmt.Errorf("[DEBUG] Error assigning Subscription ID %q to Management Group %q: %+v", subscriptionId, groupName, err)
		}
	}

	return resourceManagementGroupRead(d, meta)
}

func expandManagementGroupSubscriptionIds(input *pluginsdk.Set) []string {
	output := make([]string, 0)

	if input != nil {
		for _, v := range input.List() {
			output = append(output, v.(string))
		}
	}

	return output
}
