// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
	tagsSdk "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/tags"
	subscriptionAlias "github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	resourceClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	subscriptionClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubscriptionCreate(d *pluginsdk.ResourceData, meta any) error {
	aliasClient := subscriptionClient.From(meta).AliasClient
	client := subscriptionClient.From(meta).SubscriptionsClient
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	var aliasName string
	if aliasNameRaw, ok := d.GetOk("alias"); ok {
		aliasName = aliasNameRaw.(string)
	} else {
		aliasName = uuid.New().String()
		d.Set("alias", aliasName)
	}

	id := subscriptionAlias.NewAliasID(aliasName)

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := aliasClient.AliasGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existence of Subscription by Alias %q: %+v", id.AliasName, err)
			}
		}

		if model := existing.Model; model != nil && model.Properties != nil {
			return tf.ImportAsExistsError("azurerm_subscription", id.ID())
		}
	}

	locks.ByName(aliasName, SubscriptionResourceName)
	defer locks.UnlockByName(aliasName, SubscriptionResourceName)

	workload := subscriptionAlias.WorkloadProduction
	workloadRaw := d.Get("workload").(string)
	if workloadRaw != "" {
		workload = subscriptionAlias.Workload(workloadRaw)
	}

	req := subscriptionAlias.PutAliasRequest{
		Properties: &subscriptionAlias.PutAliasRequestProperties{
			Workload: &workload,
		},
	}

	// Check if we're adding alias management for an existing subscription
	if subscriptionIdRaw, ok := d.GetOk("subscription_id"); ok {
		subscriptionId := subscriptionIdRaw.(string)
		subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)

		locks.ByID(subscriptionId)
		defer locks.UnlockByID(subscriptionId)

		// Terraform assumes a 1:1 mapping between a Subscription and an Alias - first check if there's any existing aliases
		exists, aliasCount, err := checkExistingAliases(ctx, *aliasClient, subscriptionId)
		if err != nil {
			return err
		}
		if exists != nil {
			if aliasCount > 1 {
				return fmt.Errorf("multiple Aliases for Subscription %q already exist - to be managed via Terraform only one Alias can exist and this resource needs to be imported into the State. Please see the resource documentation for %q for more information", subscriptionId, "azurerm_subscription")
			}
			return fmt.Errorf("an Alias for Subscription %q already exists with name %q - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information", subscriptionId, *exists, "azurerm_subscription")
		}

		req.Properties.SubscriptionId = pointer.To(subscriptionId)
		existingSub, err := client.Get(ctx, subscriptionResourceId)
		if err != nil {
			return fmt.Errorf("retrieving existing %s: %+v", subscriptionResourceId, err)
		}
		if existingSub.Model == nil {
			return fmt.Errorf("retrieving existing %s: `model` was nil", subscriptionResourceId)
		}
		if existingSub.Model.State == nil {
			return fmt.Errorf("retrieving existing %s: `model.State` was nil", subscriptionResourceId)
		}

		// Disabled and Warned are both "effectively" cancelled states,
		if *existingSub.Model.State == subscriptions.SubscriptionStateDisabled || *existingSub.Model.State == subscriptions.SubscriptionStateWarned {
			log.Printf("[DEBUG] Existing subscription in Disabled/Cancelled state Terraform will attempt to re-activate it")
			if _, err := aliasClient.SubscriptionEnable(ctx, subscriptionResourceId); err != nil {
				return fmt.Errorf("enabling Subscription %q: %+v", subscriptionId, err)
			}
			deadline, _ := ctx.Deadline()
			createDeadline := time.Until(deadline)
			if err := waitForSubscriptionStateToSettle(ctx, client, subscriptionResourceId, "Active", createDeadline); err != nil {
				return fmt.Errorf("failed waiting for Subscription %q (Alias %q) to enter %q state: %+v", subscriptionId, id.AliasName, "Active", err)
			}
		}
	} else {
		// If we're not assuming control of an existing Subscription, we need to know where to create it.
		req.Properties.DisplayName = pointer.To(d.Get("subscription_name").(string))
		req.Properties.BillingScope = pointer.To(d.Get("billing_scope_id").(string))
	}

	if err := aliasClient.AliasCreateCallbackThenPoll(ctx, id, req, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating new Subscription (Alias %q): %+v", aliasName, err)
	}
	d.SetId(id.ID())

	alias, err := aliasClient.AliasGet(ctx, id)
	if err != nil || alias.Model == nil || alias.Model.Properties == nil || alias.Model.Properties.SubscriptionId == nil {
		return fmt.Errorf("failed reading subscription details for Alias %q: %+v", id.AliasName, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return errors.New("internal-error: context with no deadline")
	}
	createDeadline := time.Until(deadline)

	subscriptionResourceId := commonids.NewSubscriptionID(*alias.Model.Properties.SubscriptionId)
	if err := waitForSubscriptionStateToSettle(ctx, client, subscriptionResourceId, "Active", createDeadline); err != nil {
		return fmt.Errorf("failed waiting for Subscription %q (Alias %q) to enter %q state: %+v", *alias.Model.Properties.SubscriptionId, id.AliasName, "Active", err)
	}

	if d.HasChange("tags") {
		tagsClient := resourceClient.From(meta).TagsClient
		t := tags.Expand(d.Get("tags").(map[string]any))
		scope := commonids.NewScopeID(commonids.NewSubscriptionID(*alias.Model.Properties.SubscriptionId).ID())
		tagsResource := tagsSdk.TagsResource{
			Properties: tagsSdk.Tags{
				Tags: t,
			},
		}
		if _, err := tagsClient.CreateOrUpdateAtScope(ctx, scope, tagsResource); err != nil {
			return fmt.Errorf("setting tags on %s: %+v", id, err)
		}
	}

	return resourceSubscriptionRead(d, meta)
}
