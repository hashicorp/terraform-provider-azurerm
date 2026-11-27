// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	tagsSdk "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	resourceClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubscriptionUpdate(d *pluginsdk.ResourceData, meta any) error {
	aliasClient := client.From(meta).AliasClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := subscriptions.ParseAliasID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.AliasName, SubscriptionResourceName)
	defer locks.UnlockByName(id.AliasName, SubscriptionResourceName)
	resp, err := aliasClient.AliasGet(ctx, *id)
	if err != nil || resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.SubscriptionId == nil {
		return fmt.Errorf("could not read Subscription Alias for update: %+v", err)
	}

	subscriptionId := commonids.NewSubscriptionID(*resp.Model.Properties.SubscriptionId)

	if d.HasChange("subscription_name") {
		locks.ByID(subscriptionId.ID())
		defer locks.UnlockByID(subscriptionId.ID())

		displayName := subscriptions.SubscriptionName{
			SubscriptionName: pointer.To(d.Get("subscription_name").(string)),
		}
		if _, err := aliasClient.SubscriptionRename(ctx, subscriptionId, displayName); err != nil {
			return fmt.Errorf("could not update Display Name of Subscription %q: %+v", subscriptionId, err)
		}
	}

	if d.HasChange("tags") {
		tagsClient := resourceClient.From(meta).TagsClient
		t := tags.Expand(d.Get("tags").(map[string]any))
		scope := commonids.NewScopeID(subscriptionId.ID())
		tagsResource := tagsSdk.TagsResource{
			Properties: tagsSdk.Tags{
				Tags: t,
			},
		}
		if _, err = tagsClient.CreateOrUpdateAtScope(ctx, scope, tagsResource); err != nil {
			return fmt.Errorf("setting tags on %s: %+v", *id, err)
		}
	}

	return nil
}
