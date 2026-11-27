// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
	subscriptionAliasPandora "github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions"
	subscriptionClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importSubscriptionByAlias() pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta any) (data []*pluginsdk.ResourceData, err error) {
		aliasClient := subscriptionClient.From(meta).AliasClient
		client := subscriptionClient.From(meta).SubscriptionsClient
		aliasId, err := subscriptionAliasPandora.ParseAliasID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, errors.New("failed parsing Subscription Alias ID for import")
		}
		alias, err := aliasClient.AliasGet(ctx, *aliasId)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("failed reading Subscription Alias: %+v", err)
		}
		if alias.Model == nil || alias.Model.Properties == nil || alias.Model.Properties.SubscriptionId == nil {
			return []*pluginsdk.ResourceData{}, errors.New("failed reading Subscription Alias Properties, empty response or missing Subscription ID")
		}
		subscriptionResourceId := commonids.NewSubscriptionID(*alias.Model.Properties.SubscriptionId)
		subscription, err := client.Get(ctx, subscriptionResourceId)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", subscriptionResourceId, err)
		}
		if subscription.Model == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `model` was nil", subscriptionResourceId)
		}
		if subscription.Model.SubscriptionId == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `model.SubscriptionId` was nil", subscriptionResourceId)
		}
		if *subscription.Model.State != subscriptions.SubscriptionStateEnabled {
			return []*pluginsdk.ResourceData{}, errors.New("cannot import a cancelled Subscription by Alias ID, please enable the subscription prior to import")
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}
