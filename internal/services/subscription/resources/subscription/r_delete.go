// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	subscriptionAlias "github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

// (@jackofallops) - Delete here is a misnomer.  The nature of subscriptions is such that they are never truly deleted
// Deleted here means Cancelled. Cancelled subscriptions are held in this state for  90 days before being purged from
// active systems.  However, the backend billing systems _never_ remove this data, so once a Subscription ID has been
// used and purged from active use it can never be recovered nor the UUID reused.
// Note Cancelling a Subscription leaves it in one of several states, `Disabled` for a Subscription with no Resources or
// Alias assignments, `Warned` for Cancelled with "something" associated with it.
func resourceSubscriptionDelete(d *pluginsdk.ResourceData, meta any) error {
	aliasClient := meta.(*clients.Client).Subscription.AliasClient
	client := meta.(*clients.Client).Subscription.SubscriptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := subscriptionAlias.ParseAliasID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.AliasName, SubscriptionResourceName)
	defer locks.UnlockByName(id.AliasName, SubscriptionResourceName)

	// Get subscription details for later
	alias, err := aliasClient.AliasGet(ctx, *id)
	if err != nil || alias.Model == nil || alias.Model.Properties == nil {
		return fmt.Errorf("could not read Alias %q for Subscription: %+v", id.AliasName, err)
	}
	subscriptionId := pointer.From(alias.Model.Properties.SubscriptionId)
	locks.ByID(subscriptionId)
	defer locks.UnlockByID(subscriptionId)

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	sub, err := client.Get(ctx, subscriptionResourceId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", subscriptionResourceId, err)
	}
	if sub.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", subscriptionResourceId)
	}

	subscriptionName := ""
	if subscriptionNameRaw := sub.Model.DisplayName; subscriptionNameRaw != nil {
		subscriptionName = *sub.Model.DisplayName
	}

	if subscriptionName == "" || subscriptionId == "" {
		return fmt.Errorf("one or both of Subscription Name (%q) and Subscription ID (%q) could not be determined", subscriptionName, subscriptionId)
	}
	// remove the alias
	if _, count, err := checkExistingAliases(ctx, *aliasClient, subscriptionId); err != nil {
		if count > 1 {
			return fmt.Errorf("multiple Aliases found for Subscription %q, cannot remove", subscriptionId)
		}
	}

	resp, err := aliasClient.AliasDelete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("could not delete Alias %q for Subscription %q (ID: %q): %+v", id.AliasName, subscriptionName, subscriptionId, err)
		}
	}

	// Cancel the Subscription
	if !meta.(*clients.Client).Features.Subscription.PreventCancellationOnDestroy {
		log.Printf("[DEBUG] Cancelling subscription %s", subscriptionId)

		opts := subscriptionAlias.DefaultSubscriptionCancelOperationOptions()
		// TODO: support a Provider `features` flag to enable deleting a Subscription containing Resources
		// This is a dangerous operation, and likely wants a similar default value as to that for Resource Groups
		if _, err := aliasClient.SubscriptionCancel(ctx, subscriptionResourceId, opts); err != nil {
			return fmt.Errorf("failed to cancel Subscription: %+v", err)
		}

		deadline, _ := ctx.Deadline()
		deleteDeadline := time.Until(deadline)

		if err := waitForSubscriptionStateToSettle(ctx, client, subscriptionResourceId, "Cancelled", deleteDeadline); err != nil {
			return fmt.Errorf("failed to cancel Subscription %q (Alias %q): %+v", subscriptionId, id.AliasName, err)
		}
	} else {
		log.Printf("[DEBUG] Skipping subscription %s cancellation due to feature flag.", *id)
	}

	return nil
}
