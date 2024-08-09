// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managementgroup

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managementgroups/2020-05-01/managementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var managementGroupCacheControl = "no-cache"

func resourceManagementGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagementGroupCreateUpdate,
		Update: resourceManagementGroupCreateUpdate,
		Read:   resourceManagementGroupRead,
		Delete: resourceManagementGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagementGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.ManagementGroupName,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenant_scoped_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"parent_management_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: commonids.ValidateManagementGroupID,
			},

			"subscription_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
				Set: pluginsdk.HashString,
			},
		},
	}
}

func resourceManagementGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	log.Printf("[INFO] Creating Management Group %q", groupName)

	properties := managementgroups.CreateManagementGroupRequest{
		Name: utils.String(groupName),
		Properties: &managementgroups.CreateManagementGroupProperties{
			TenantId: utils.String(armTenantID),
			Details: &managementgroups.CreateManagementGroupDetails{
				Parent: &managementgroups.CreateParentGroupInfo{
					Id: utils.String(parentManagementGroupId),
				},
			},
		},
	}

	if v := d.Get("display_name"); v != "" {
		properties.Properties.DisplayName = utils.String(v.(string))
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, properties, managementgroups.CreateOrUpdateOperationOptions{
		CacheControl: &managementGroupCacheControl,
	})
	if err != nil {
		return fmt.Errorf("unable to create Management Group %q: %+v", groupName, err)
	}

	// We have a potential race condition / consistency issue whereby the implicit role assignment for the SP may not be
	// completed before the read-back here or an eventually consistent read is creating a temporary 403 error.
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			"pending",
		},
		Target: []string{
			"succeeded",
		},
		Refresh:                   managementGroupCreateStateRefreshFunc(ctx, client, id),
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
		ContinuousTargetOccurence: 5,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
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

	d.SetId(id.ID())

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
	log.Printf("[DEBUG] Preparing to assign Subscriptions to Management Group %q", groupName)
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

func resourceManagementGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	recurse := pointer.FromBool(true)
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

func resourceManagementGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseManagementGroupID(d.Id())
	if err != nil {
		return err
	}

	recurse := true
	group, err := client.Get(ctx, *id, managementgroups.GetOperationOptions{
		CacheControl: &managementGroupCacheControl,
		Filter:       pointer.To("children.childType eq Subscription"),
		Expand:       pointer.To(managementgroups.ExpandChildren),
		Recurse:      &recurse,
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

					subscriptionId, err := managementgroups.ParseSubscriptionID(*v.Id)
					if err != nil {
						return fmt.Errorf("unable to parse child Subscription ID %+v", err)
					}
					if subscriptionId == nil {
						continue
					}
					log.Printf("[DEBUG] De-associating Subscription %q from Management Group %q..", subscriptionId, id.GroupId)
					// NOTE: whilst this says `Delete` it's actually `Deassociate` - which is /really/ helpful
					deleteResp, err := client.SubscriptionsDelete(ctx, *subscriptionId, managementgroups.SubscriptionsDeleteOperationOptions{
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

	err = client.DeleteThenPoll(ctx, *id, managementgroups.DeleteOperationOptions{
		CacheControl: &managementGroupCacheControl,
	})
	if err != nil {
		return fmt.Errorf("unable to delete Management Group %q: %+v", id.GroupId, err)
	}

	return nil
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

		found := false
		for _, subId := range updated {
			if id.SubscriptionId == subId {
				found = true
				break
			}
		}

		if !found {
			subscriptionIdsToRemove = append(subscriptionIdsToRemove, id.SubscriptionId)
		}
	}

	return &subscriptionIdsToRemove, nil
}

func managementGroupCreateStateRefreshFunc(ctx context.Context, client *managementgroups.ManagementGroupsClient, id commonids.ManagementGroupId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id, managementgroups.GetOperationOptions{
			CacheControl: &managementGroupCacheControl,
			Expand:       pointer.To(managementgroups.ExpandChildren),
			Recurse:      pointer.FromBool(true),
		})
		if err != nil {
			if response.WasForbidden(resp.HttpResponse) {
				return resp, "pending", nil
			}
			return resp, "failed", err
		}

		return resp, "succeeded", nil
	}
}
