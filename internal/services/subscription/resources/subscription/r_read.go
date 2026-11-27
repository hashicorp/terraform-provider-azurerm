// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
	subscriptionAlias "github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubscriptionRead(d *pluginsdk.ResourceData, meta any) error {
	aliasClient := meta.(*clients.Client).Subscription.AliasClient
	client := meta.(*clients.Client).Subscription.SubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := subscriptionAlias.ParseAliasID(d.Id())
	if err != nil {
		return err
	}
	d.Set("alias", id.AliasName)

	alias, err := aliasClient.AliasGet(ctx, *id)
	if err != nil || alias.Model == nil {
		if response.WasNotFound(alias.HttpResponse) {
			log.Printf("[INFO] Error reading Subscription %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Subscription Alias %q: %+v", id.AliasName, err)
	}

	subscriptionId := ""
	subscriptionName := ""
	tenantId := ""
	var t *map[string]string
	if props := alias.Model.Properties; props != nil && props.SubscriptionId != nil {
		subscriptionId = *props.SubscriptionId
		subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
		resp, err := client.Get(ctx, subscriptionResourceId)
		if err != nil {
			return fmt.Errorf("retrieving %s (Alias %q) to obtain the Tenant Information: %+v", subscriptionResourceId, id.AliasName, err)
		}
		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil", subscriptionResourceId)
		}

		if model := resp.Model; model != nil {
			subscriptionName = pointer.From(model.DisplayName)
			tenantId = pointer.From(model.TenantId)
			t = model.Tags
		}
	}

	// (@jackofallops) A subscription's billing scope is not exposed in any way in the API/SDK so we cannot read it back here

	d.Set("subscription_id", subscriptionId)
	d.Set("subscription_name", subscriptionName)
	d.Set("tenant_id", tenantId)
	if err := tags.FlattenAndSet(d, t); err != nil {
		return err
	}

	return nil
}

func waitForSubscriptionStateToSettle(ctx context.Context, client *subscriptions.SubscriptionsClient, subscriptionId commonids.SubscriptionId, targetState string, timeout time.Duration) error {
	stateConf := &pluginsdk.StateChangeConf{
		Refresh: func() (result any, state string, err error) {
			status, err := client.Get(ctx, subscriptionId)
			if err != nil {
				return status, "Failed", err
			}
			if status.Model == nil || status.Model.State == nil {
				return status, "Unknown", err
			}

			return status, string(*status.Model.State), err
		},
		PollInterval:              10 * time.Second,
		Timeout:                   timeout,
		ContinuousTargetOccurence: 4,
		Delay:                     60 * time.Second,
	}
	switch targetState {
	case "Cancelled":
		stateConf.Target = []string{
			string(subscriptions.SubscriptionStateDisabled),
			string(subscriptions.SubscriptionStateWarned),
		}
		stateConf.Pending = []string{
			string(subscriptions.SubscriptionStateEnabled),
			"", // The `State` field can be empty whilst being updated
		}

	case "Active":
		stateConf.Target = []string{
			string(subscriptions.SubscriptionStateEnabled),
		}
		stateConf.Pending = []string{
			string(subscriptions.SubscriptionStateDisabled),
			string(subscriptions.SubscriptionStateWarned),
			"", // The `State` field can be empty whilst being updated
		}
	default:
		return fmt.Errorf("unsupported target state %q for Subscription %q", targetState, subscriptionId)
	}

	if actual, err := stateConf.WaitForStateContext(ctx); err != nil {
		sub, ok := actual.(subscriptions.Subscription)
		if !ok {
			return fmt.Errorf("failure in parsing response while waiting for Subscription %q to become %q: %+v", subscriptionId, targetState, err)
		}
		actualState := pointer.FromEnum(sub.State)
		return fmt.Errorf("waiting for Subscription %q to become %q, currently %q", subscriptionId, targetState, actualState)
	}

	return nil
}

func checkExistingAliases(ctx context.Context, client subscriptionAlias.SubscriptionsClient, subscriptionId string) (*string, int, error) {
	aliasList, err := client.AliasListComplete(ctx)
	if err != nil {
		return nil, 0, errors.New("could not List existing Subscription Aliases")
	}

	for _, v := range aliasList.Items {
		if v.Properties != nil && v.Properties.SubscriptionId != nil && subscriptionId == *v.Properties.SubscriptionId {
			return v.Name, len(aliasList.Items), nil
		}
	}

	return nil, len(aliasList.Items), nil
}
