package subscription

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2021-01-01/subscriptions"
	subscriptionAlias "github.com/Azure/azure-sdk-for-go/services/subscription/mgmt/2020-09-01/subscription"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var SubscriptionNameResourceName = "azurerm_subscription_name"

func resourceSubscriptionName() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSubscriptionNameCreate,
		Update: resourceSubscriptionNameUpdate,
		Read:   resourceSubscriptionNameRead,
		Delete: resourceSubscriptionNameDelete,

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
			"alias": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The Alias Name of the subscription. If omitted a new UUID will be generated for this property.",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Description:  "The Display Name for the Subscription.",
				ValidateFunc: validate.SubscriptionName,
			},

			"subscription_id": {
				Type:         pluginsdk.TypeString,
				Description:  "The GUID of the Subscription.",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:        pluginsdk.TypeString,
				Description: "The Tenant ID to which the subscription belongs",
				Computed:    true,
			},

			"tags": tags.Schema(),

			"workload": {
				Type:        pluginsdk.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The workload type for the Subscription. Possible values are `Production` (default) and `DevTest`.",
				// Other RP's have updated Constants with contextual prefixes so these are likely to change
				ValidateFunc: validation.StringInSlice([]string{
					string(subscriptionAlias.Production),
					string(subscriptionAlias.DevTest),
				}, false),
				// Workload is not exposed in any way, so must be ignored if the resource is imported.
				DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
					return new == ""
				},
			},
		},
	}
}

func resourceSubscriptionNameCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	aliasClient := meta.(*clients.Client).Subscription.AliasClient
	subscriptionClient := meta.(*clients.Client).Subscription.SubscriptionClient
	client := meta.(*clients.Client).Subscription.Client
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	aliasName := ""
	if aliasNameRaw, ok := d.GetOk("alias"); ok {
		aliasName = aliasNameRaw.(string)
	} else {
		aliasName = uuid.New().String()
		d.Set("alias", aliasName)
	}

	id := parse.NewSubscriptionAliasId(aliasName)

	existing, err := aliasClient.Get(ctx, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existence of Subscription by Alias %q: %+v", id.Name, err)
		}
	}

	if props := existing.Properties; props != nil {
		return tf.ImportAsExistsError("azurerm_subscription_name", id.ID())
	}

	locks.ByName(aliasName, SubscriptionNameResourceName)
	defer locks.UnlockByName(aliasName, SubscriptionNameResourceName)

	workload := subscriptionAlias.Production
	workloadRaw := d.Get("workload").(string)
	if workloadRaw != "" {
		workload = subscriptionAlias.Workload(workloadRaw)
	}

	req := subscriptionAlias.PutAliasRequest{
		Properties: &subscriptionAlias.PutAliasRequestProperties{
			Workload: workload,
		},
	}

	subscriptionId := d.Get("subscription_id").(string)

	locks.ByID(subscriptionId)
	defer locks.UnlockByID(subscriptionId)

	// Terraform assumes a 1:1 mapping between a Subscription and an Alias - first check if there's any existing aliases
	exists, aliasCount, err := checkExistingAliases(ctx, *aliasClient, subscriptionId)
	if err != nil {
		return err
	}
	if exists != nil {
		if aliasCount > 1 {
			return fmt.Errorf("multiple Aliases for Subscription %q already exist - to be managed via Terraform only one Alias can exist and this resource needs to be imported into the State. Please see the resource documentation for %q for more information", subscriptionId, "azurerm_subscription_name")
		}
		return tf.ImportAsExistsError("azurerm_subscription_name", id.ID())
	}

	req.Properties.SubscriptionID = utils.String(subscriptionId)
	existingSub, err := client.Get(ctx, subscriptionId)
	if err != nil {
		return fmt.Errorf("could not read existing Subscription %q", subscriptionId)
	}
	// Disabled and Warned are both "effectively" cancelled states,
	if existingSub.State == subscriptions.StateDisabled || existingSub.State == subscriptions.StateWarned {
		return fmt.Errorf("subscription %q has state: %s, cannot manage. %+v", subscriptionId, existingSub.State, err)
	}

	// create alias
	future, err := aliasClient.Create(ctx, aliasName, req)
	if err != nil {
		return fmt.Errorf("creating new Subscription (Alias %q): %+v", aliasName, err)
	}

	if err := future.WaitForCompletionRef(ctx, aliasClient.Client); err != nil {
		return fmt.Errorf("waiting for creation of Subscription with Alias %q: %+v", id.Name, err)
	}

	alias, err := aliasClient.Get(ctx, id.Name)
	if err != nil || alias.Properties == nil || alias.Properties.SubscriptionID == nil {
		return fmt.Errorf("failed reading subscription details for Alias %q: %+v", id.Name, err)
	}

	// update subscription name
	displayName := subscriptionAlias.Name{
		SubscriptionName: utils.String(d.Get("name").(string)),
	}
	if _, err := subscriptionClient.Rename(ctx, subscriptionId, displayName); err != nil {
		return fmt.Errorf("could not update Display Name of Subscription %q: %+v", subscriptionId, err)
	}

	deadline, _ := ctx.Deadline()
	createDeadline := time.Until(deadline)

	err = waitForNameChange(ctx, meta.(*clients.Client), subscriptionId, d.Get("name").(string), createDeadline)
	if err != nil {
		return err
	}

	// update tags
	tagsClient := meta.(*clients.Client).Resource.TagsClientForSubscription(*alias.Properties.SubscriptionID)
	t := tags.Expand(d.Get("tags").(map[string]interface{}))
	scope := fmt.Sprintf("subscriptions/%s", *alias.Properties.SubscriptionID)
	tagsResource := resources.TagsResource{
		Properties: &resources.Tags{
			Tags: t,
		},
	}
	if _, err = tagsClient.CreateOrUpdateAtScope(ctx, scope, tagsResource); err != nil {
		return fmt.Errorf("setting tags on %s: %+v", id, err)
	}

	// wait for tags to replicate
	deadline, _ = ctx.Deadline()
	createDeadline = time.Until(deadline)
	err = waitForTagsChange(ctx, meta.(*clients.Client), subscriptionId, scope, t, createDeadline)
	if err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceSubscriptionNameRead(d, meta)
}

func resourceSubscriptionNameUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	aliasClient := meta.(*clients.Client).Subscription.AliasClient
	subscriptionClient := meta.(*clients.Client).Subscription.SubscriptionClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionAliasID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, SubscriptionNameResourceName)
	defer locks.UnlockByName(id.Name, SubscriptionNameResourceName)
	resp, err := aliasClient.Get(ctx, id.Name)
	if err != nil || resp.Properties == nil {
		return fmt.Errorf("could not read Subscription Alias for update: %+v", err)
	}

	subscriptionId := resp.Properties.SubscriptionID
	if subscriptionId == nil || *subscriptionId == "" {
		return fmt.Errorf("could not read Subscription ID from Alias")
	}

	if d.HasChange("name") {
		locks.ByID(*subscriptionId)
		defer locks.UnlockByID(*subscriptionId)

		displayName := subscriptionAlias.Name{
			SubscriptionName: utils.String(d.Get("name").(string)),
		}
		if _, err := subscriptionClient.Rename(ctx, *subscriptionId, displayName); err != nil {
			return fmt.Errorf("could not update Display Name of Subscription %q: %+v", *subscriptionId, err)
		}

		deadline, _ := ctx.Deadline()
		updateDeadline := time.Until(deadline)

		err = waitForNameChange(ctx, meta.(*clients.Client), *subscriptionId, d.Get("name").(string), updateDeadline)
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		tagsClient := meta.(*clients.Client).Resource.TagsClientForSubscription(*subscriptionId)
		t := tags.Expand(d.Get("tags").(map[string]interface{}))
		scope := fmt.Sprintf("subscriptions/%s", *subscriptionId)
		tagsResource := resources.TagsResource{
			Properties: &resources.Tags{
				Tags: t,
			},
		}
		if _, err = tagsClient.CreateOrUpdateAtScope(ctx, scope, tagsResource); err != nil {
			return fmt.Errorf("setting tags on %s: %+v", *id, err)
		}
		deadline, _ := ctx.Deadline()
		updateDeadline := time.Until(deadline)
		err = waitForTagsChange(ctx, meta.(*clients.Client), *subscriptionId, scope, t, updateDeadline)
		if err != nil {
			return err
		}
	}

	return resourceSubscriptionNameRead(d, meta)
}

func resourceSubscriptionNameRead(d *pluginsdk.ResourceData, meta interface{}) error {
	aliasClient := meta.(*clients.Client).Subscription.AliasClient
	client := meta.(*clients.Client).Subscription.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionAliasID(d.Id())
	if err != nil {
		return err
	}
	d.Set("alias", id.Name)

	alias, err := aliasClient.Get(ctx, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(alias.Response) {
			log.Printf("[INFO] Error reading Subscription %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Subscription Alias %q: %+v", id.Name, err)
	}

	subscriptionId := ""
	subscriptionName := ""
	tenantId := ""
	t := make(map[string]*string)
	if props := alias.Properties; props != nil && props.SubscriptionID != nil {
		subscriptionId = *props.SubscriptionID
		resp, err := client.Get(ctx, subscriptionId)

		if err != nil {
			return fmt.Errorf("failed to read Subscription %q (Alias %q) for Tenant Information: %+v", subscriptionId, id.Name, err)
		}
		if resp.TenantID != nil {
			tenantId = *resp.TenantID
		}

		if resp.DisplayName != nil {
			subscriptionName = *resp.DisplayName
		}
	}

	scope := fmt.Sprintf("subscriptions/%s", *alias.Properties.SubscriptionID)
	tagsObject, err := meta.(*clients.Client).Resource.TagsClientForSubscription(subscriptionId).GetAtScope(ctx, scope)
	if err != nil {
		if !utils.ResponseWasNotFound(tagsObject.Response) {
			return fmt.Errorf("failed to read Subscription %q tags: %+v", subscriptionId, err)
		}
	} else {
		t = tagsObject.Properties.Tags
	}

	d.Set("subscription_id", subscriptionId)
	d.Set("name", subscriptionName)
	d.Set("tenant_id", tenantId)
	if err := tags.FlattenAndSet(d, t); err != nil {
		return err
	}

	return nil
}

func waitForNameChange(ctx context.Context, clients *clients.Client, subscriptionId string, subscriptionName string, timeout time.Duration) error {
	stateConf := &pluginsdk.StateChangeConf{
		Refresh: func() (result interface{}, state string, err error) {
			status, err := clients.Subscription.Client.Get(ctx, subscriptionId)
			return status, *status.DisplayName, err
		},
		PollInterval:              30 * time.Second,
		Timeout:                   timeout,
		ContinuousTargetOccurence: 2,
		Delay:                     60 * time.Second,
	}
	stateConf.Target = []string{subscriptionName}

	if actual, err := stateConf.WaitForStateContext(ctx); err != nil {
		sub, ok := actual.(subscriptions.Subscription)
		if !ok {
			return fmt.Errorf("failure in parsing response while waiting for Subscription %s to become %s: %+v", subscriptionId, subscriptionName, err)
		}
		actualState := sub.DisplayName
		return fmt.Errorf("waiting for Subscription %s to become %s, currently %s", subscriptionId, subscriptionName, *actualState)
	}

	return nil
}

func waitForTagsChange(ctx context.Context, clients *clients.Client, subscriptionId string, scope string, expectedTags map[string]*string, timeout time.Duration) error {
	stateConf := &pluginsdk.StateChangeConf{
		Refresh: func() (result interface{}, state string, err error) {
			tagsObject, err := clients.Resource.TagsClientForSubscription(subscriptionId).GetAtScope(ctx, scope)

			stateValue := "Incorrect"
			if tagsObject.Properties != nil {
				if reflect.DeepEqual(tagsObject.Properties.Tags, expectedTags) {
					stateValue = "Correct"
				}
			}
			return tagsObject, stateValue, err
		},
		PollInterval: 10 * time.Second,
		Timeout:      timeout,
		Delay:        10 * time.Second,
	}
	stateConf.Target = []string{"Correct"}

	if actual, err := stateConf.WaitForStateContext(ctx); err != nil {
		sub, ok := actual.(resources.TagsResource)
		if !ok {
			return fmt.Errorf("failure in parsing response while waiting for Subscription %s tags to become %v: %+v", subscriptionId, expectedTags, err)
		}
		actualState := sub.Properties.Tags
		return fmt.Errorf("waiting for Subscription %s tags to become %v, currently %v", subscriptionId, expectedTags, actualState)
	}

	return nil
}

func resourceSubscriptionNameDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	aliasClient := meta.(*clients.Client).Subscription.AliasClient
	client := meta.(*clients.Client).Subscription.Client
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionAliasID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, SubscriptionNameResourceName)
	defer locks.UnlockByName(id.Name, SubscriptionNameResourceName)

	// Get subscription details
	alias, err := aliasClient.Get(ctx, id.Name)
	if err != nil || alias.Properties == nil {
		return fmt.Errorf("could not read Alias %q for Subscription: %+v", id.Name, err)
	}
	subscriptionId := ""
	if subscriptionIdRaw := alias.Properties.SubscriptionID; subscriptionIdRaw != nil {
		subscriptionId = *subscriptionIdRaw
	}

	locks.ByID(subscriptionId)
	defer locks.UnlockByID(subscriptionId)

	sub, err := client.Get(ctx, subscriptionId)
	if err != nil {
		return fmt.Errorf("could not read Subscription details for %q: %+v", subscriptionId, err)
	}

	subscriptionName := ""
	if subscriptionNameRaw := sub.DisplayName; subscriptionNameRaw != nil {
		subscriptionName = *sub.DisplayName
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

	resp, err := aliasClient.Delete(ctx, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("could not delete Alias %q for Subscription %q (ID: %q): %+v", id.Name, subscriptionName, subscriptionId, err)
		}
	}

	return nil
}
