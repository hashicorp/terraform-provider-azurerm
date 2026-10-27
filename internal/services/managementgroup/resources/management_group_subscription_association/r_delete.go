// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_subscription_association

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/management/2020-05-01/managementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementGroupSubscriptionAssociationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managementgroups.ParseSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SubscriptionsDelete(ctx, *id, managementgroups.SubscriptionsDeleteOperationOptions{})
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting Management Group Subscription Association between Management Group %q and Subscription %q: %+v", id.GroupId, id.SubscriptionId, err)
		}
	}

	// It's a workaround to solve the replication delay issue: DELETE operation happens in one region, but it needs more time to sync the result to other regions.
	log.Printf("[DEBUG] Waiting for %s to be fully deleted..", d.Id())
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   subscriptionAssociationRefreshFunc(ctx, meta.(*clients.Client).ManagementGroups.GroupsClient, *id),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be fully deleted: %+v", d.Id(), err)
	}

	return nil
}
