// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
	tagsSdk "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/tags"
	subscriptionAlias "github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	billingValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var SubscriptionResourceName = "azurerm_subscription"

func resourceSubscription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSubscriptionCreate,
		Update: resourceSubscriptionUpdate,
		Read:   resourceSubscriptionRead,
		Delete: resourceSubscriptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.SubscriptionAliasID(id)
			return err
		}, importSubscriptionByAlias()),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"subscription_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Description:  "The Display Name for the Subscription.",
				ValidateFunc: validate.SubscriptionName,
			},

			"alias": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true, // O+C - This value is supplied by the provider if omitted so must remain `Computed`
				ForceNew:     true,
				Description:  "The Alias Name of the subscription. If omitted a new UUID will be generated for this property.",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"billing_scope_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"subscription_id",
					"billing_scope_id",
				},
				ValidateFunc: validation.Any(
					billingValidate.MicrosoftCustomerAccountBillingScopeID,
					billingValidate.EnrollmentBillingScopeID,
					billingValidate.MicrosoftPartnerAccountBillingScopeID,
				),
			},

			"workload": {
				Type:        pluginsdk.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The workload type for the Subscription. Possible values are `Production` (default) and `DevTest`.",
				// Other RP's have updated Constants with contextual prefixes so these are likely to change
				ValidateFunc: validation.StringInSlice([]string{
					string(subscriptionAlias.WorkloadProduction),
					string(subscriptionAlias.WorkloadDevTest),
				}, false),
				// Workload is not exposed in any way, so must be ignored if the resource is imported.
				DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
					return new == ""
				},
			},

			"subscription_id": {
				Type:        pluginsdk.TypeString,
				Description: "The GUID of the Subscription.",
				ForceNew:    true,
				Optional:    true,
				Computed:    true, // O+C This must remain computed due to the unique nature of this resource - See resource documentation for notes.
				ExactlyOneOf: []string{
					"subscription_id",
					"billing_scope_id",
				},
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:        pluginsdk.TypeString,
				Description: "The Tenant ID to which the subscription belongs",
				Computed:    true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceSubscriptionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	aliasClient := meta.(*clients.Client).Subscription.AliasClient
	client := meta.(*clients.Client).Subscription.SubscriptionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	aliasName := ""
	if aliasNameRaw, ok := d.GetOk("alias"); ok {
		aliasName = aliasNameRaw.(string)
	} else {
		aliasName = uuid.New().String()
		d.Set("alias", aliasName)
	}

	id := subscriptionAlias.NewAliasID(aliasName)
	existing, err := aliasClient.AliasGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existence of Subscription by Alias %q: %+v", id.AliasName, err)
		}
	}

	if model := existing.Model; model != nil && model.Properties != nil {
		return tf.ImportAsExistsError("azurerm_subscription", id.ID())
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

	subscriptionId := ""

	// Check if we're adding alias management for an existing subscription
	if subscriptionIdRaw, ok := d.GetOk("subscription_id"); ok {
		subscriptionId = subscriptionIdRaw.(string)
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

		req.Properties.SubscriptionId = utils.String(subscriptionId)
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
		req.Properties.DisplayName = utils.String(d.Get("subscription_name").(string))
		req.Properties.BillingScope = utils.String(d.Get("billing_scope_id").(string))
	}

	if err := aliasClient.AliasCreateThenPoll(ctx, id, req); err != nil {
		return fmt.Errorf("creating new Subscription (Alias %q): %+v", aliasName, err)
	}

	alias, err := aliasClient.AliasGet(ctx, id)
	if err != nil || alias.Model == nil || alias.Model.Properties == nil || alias.Model.Properties.SubscriptionId == nil {
		return fmt.Errorf("failed reading subscription details for Alias %q: %+v", id.AliasName, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context with no deadline")
	}
	createDeadline := time.Until(deadline)

	subscriptionResourceId := commonids.NewSubscriptionID(*alias.Model.Properties.SubscriptionId)
	if err := waitForSubscriptionStateToSettle(ctx, client, subscriptionResourceId, "Active", createDeadline); err != nil {
		return fmt.Errorf("failed waiting for Subscription %q (Alias %q) to enter %q state: %+v", *alias.Model.Properties.SubscriptionId, id.AliasName, "Active", err)
	}

	if d.HasChange("tags") {
		tagsClient := meta.(*clients.Client).Resource.TagsClient
		t := tags.Expand(d.Get("tags").(map[string]interface{}))
		scope := commonids.NewScopeID(commonids.NewSubscriptionID(*alias.Model.Properties.SubscriptionId).ID())
		tagsResource := tagsSdk.TagsResource{
			Properties: tagsSdk.Tags{
				Tags: t,
			},
		}
		if _, err = tagsClient.CreateOrUpdateAtScope(ctx, scope, tagsResource); err != nil {
			return fmt.Errorf("setting tags on %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceSubscriptionRead(d, meta)
}

func resourceSubscriptionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	aliasClient := meta.(*clients.Client).Subscription.AliasClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := subscriptionAlias.ParseAliasID(d.Id())
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

		displayName := subscriptionAlias.SubscriptionName{
			SubscriptionName: utils.String(d.Get("subscription_name").(string)),
		}
		if _, err := aliasClient.SubscriptionRename(ctx, subscriptionId, displayName); err != nil {
			return fmt.Errorf("could not update Display Name of Subscription %q: %+v", subscriptionId, err)
		}
	}

	if d.HasChange("tags") {
		tagsClient := meta.(*clients.Client).Resource.TagsClient
		t := tags.Expand(d.Get("tags").(map[string]interface{}))
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

func resourceSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

// (@jackofallops) - Delete here is a misnomer.  The nature of subscriptions is such that they are never truly deleted
// Deleted here means Cancelled. Cancelled subscriptions are held in this state for  90 days before being purged from
// active systems.  However, the backend billing systems _never_ remove this data, so once a Subscription ID has been
// used and purged from active use it can never be recovered nor the UUID reused.
// Note Cancelling a Subscription leaves it in one of several states, `Disabled` for a Subscription with no Resources or
// Alias assignments, `Warned` for Cancelled with "something" associated with it.
func resourceSubscriptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
	subscriptionId := ""
	if subscriptionIdRaw := alias.Model.Properties.SubscriptionId; subscriptionIdRaw != nil {
		subscriptionId = *subscriptionIdRaw
	}
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

func waitForSubscriptionStateToSettle(ctx context.Context, client *subscriptions.SubscriptionsClient, subscriptionId commonids.SubscriptionId, targetState string, timeout time.Duration) error {
	stateConf := &pluginsdk.StateChangeConf{
		Refresh: func() (result interface{}, state string, err error) {
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
		actualState := string(pointer.From(sub.State))
		return fmt.Errorf("waiting for Subscription %q to become %q, currently %q", subscriptionId, targetState, actualState)
	}

	return nil
}

func checkExistingAliases(ctx context.Context, client subscriptionAlias.SubscriptionsClient, subscriptionId string) (*string, int, error) {
	aliasList, err := client.AliasListComplete(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("could not List existing Subscription Aliases")
	}

	for _, v := range aliasList.Items {
		if v.Properties != nil && v.Properties.SubscriptionId != nil && subscriptionId == *v.Properties.SubscriptionId {
			return v.Name, len(aliasList.Items), nil
		}
	}

	return nil, len(aliasList.Items), nil
}
