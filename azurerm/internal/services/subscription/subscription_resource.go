package subscription

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-11-01/subscriptions"
	subscriptionAlias "github.com/Azure/azure-sdk-for-go/services/subscription/mgmt/2020-09-01/subscription"
	"github.com/google/uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	billingValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/billing/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
				Computed:     true,
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
				),
			},

			// Optional
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

			"subscription_id": {
				Type:        pluginsdk.TypeString,
				Description: "The GUID of the Subscription.",
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
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

			"tags": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceSubscriptionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return tf.ImportAsExistsError("azurerm_subscription", id.ID())
	}

	locks.ByName(aliasName, SubscriptionResourceName)
	defer locks.UnlockByName(aliasName, SubscriptionResourceName)

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

	subscriptionId := ""

	// Check if we're adding alias management for an existing subscription
	if subscriptionIdRaw, ok := d.GetOk("subscription_id"); ok {
		subscriptionId = subscriptionIdRaw.(string)

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

		req.Properties.SubscriptionID = utils.String(subscriptionId)
		existingSub, err := client.Get(ctx, subscriptionId)
		if err != nil {
			return fmt.Errorf("could not read existing Subscription %q", subscriptionId)
		}
		// Disabled and Warned are both "effectively" cancelled states,
		if existingSub.State == subscriptions.Disabled || existingSub.State == subscriptions.Warned {
			log.Printf("[DEBUG] Existing subscription in Disabled/Cancelled state Terraform will attempt to re-activate it")
			if _, err := subscriptionClient.Enable(ctx, subscriptionId); err != nil {
				return fmt.Errorf("enabling Subscription %q: %+v", subscriptionId, err)
			}
			deadline, _ := ctx.Deadline()
			createDeadline := time.Until(deadline)
			if err := waitForSubscriptionStateToSettle(ctx, meta.(*clients.Client), subscriptionId, "Active", createDeadline); err != nil {
				return fmt.Errorf("failed waiting for Subscription %q (Alias %q) to enter %q state: %+v", subscriptionId, id.Name, "Active", err)
			}
		}
	} else {
		// If we're not assuming control of an existing Subscription, we need to know where to create it.
		req.Properties.DisplayName = utils.String(d.Get("subscription_name").(string))
		req.Properties.BillingScope = utils.String(d.Get("billing_scope_id").(string))
	}

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

	deadline, _ := ctx.Deadline()
	createDeadline := time.Until(deadline)

	if err := waitForSubscriptionStateToSettle(ctx, meta.(*clients.Client), *alias.Properties.SubscriptionID, "Active", createDeadline); err != nil {
		return fmt.Errorf("failed waiting for Subscription %q (Alias %q) to enter %q state: %+v", subscriptionId, id.Name, "Active", err)
	}

	d.SetId(id.ID())

	return resourceSubscriptionRead(d, meta)
}

func resourceSubscriptionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	aliasClient := meta.(*clients.Client).Subscription.AliasClient
	subscriptionClient := meta.(*clients.Client).Subscription.SubscriptionClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionAliasID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, SubscriptionResourceName)
	defer locks.UnlockByName(id.Name, SubscriptionResourceName)
	resp, err := aliasClient.Get(ctx, id.Name)
	if err != nil || resp.Properties == nil {
		return fmt.Errorf("could not read Subscription Alias for update: %+v", err)
	}

	subscriptionId := resp.Properties.SubscriptionID
	if subscriptionId == nil || *subscriptionId == "" {
		return fmt.Errorf("could not read Subscription ID from Alias")
	}

	if d.HasChange("subscription_name") {
		locks.ByID(*subscriptionId)
		defer locks.UnlockByID(*subscriptionId)

		displayName := subscriptionAlias.Name{
			SubscriptionName: utils.String(d.Get("subscription_name").(string)),
		}
		if _, err := subscriptionClient.Rename(ctx, *subscriptionId, displayName); err != nil {
			return fmt.Errorf("could not update Display Name of Subscription %q: %+v", *subscriptionId, err)
		}
	}

	return nil
}

func resourceSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

		t = resp.Tags
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
	subscriptionClient := meta.(*clients.Client).Subscription.SubscriptionClient
	client := meta.(*clients.Client).Subscription.Client
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionAliasID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, SubscriptionResourceName)
	defer locks.UnlockByName(id.Name, SubscriptionResourceName)

	// Get subscription details for later
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

	// Cancel the Subscription
	if _, err := subscriptionClient.Cancel(ctx, subscriptionId); err != nil {
		return fmt.Errorf("failed to cancel Subscription: %+v", err)
	}

	deadline, _ := ctx.Deadline()
	deleteDeadline := time.Until(deadline)

	if err := waitForSubscriptionStateToSettle(ctx, meta.(*clients.Client), subscriptionId, "Cancelled", deleteDeadline); err != nil {
		return fmt.Errorf("failed to cancel Subscription %q (Alias %q): %+v", subscriptionId, id.Name, err)
	}

	return nil
}

func waitForSubscriptionStateToSettle(ctx context.Context, clients *clients.Client, subscriptionId string, targetState string, timeout time.Duration) error {
	stateConf := &pluginsdk.StateChangeConf{
		Refresh: func() (result interface{}, state string, err error) {
			status, err := clients.Subscription.Client.Get(ctx, subscriptionId)
			return status, string(status.State), err
		},
		PollInterval:              20 * time.Second,
		Timeout:                   timeout,
		ContinuousTargetOccurence: 2,
		Delay:                     60 * time.Second,
	}
	switch targetState {
	case "Cancelled":
		stateConf.Target = []string{
			string(subscriptions.Disabled),
			string(subscriptions.Warned),
		}
		stateConf.Pending = []string{
			string(subscriptions.Enabled),
		}

	case "Active":
		stateConf.Target = []string{
			string(subscriptions.Enabled),
		}
		stateConf.Pending = []string{
			string(subscriptions.Disabled),
			string(subscriptions.Warned),
		}
	default:
		return fmt.Errorf("unsupported target state %q for Subscription %q", targetState, subscriptionId)
	}

	if actual, err := stateConf.WaitForStateContext(ctx); err != nil {
		sub, ok := actual.(subscriptions.Subscription)
		if !ok {
			return fmt.Errorf("failure in parsing response while waiting for Subscription %q to become %q: %+v", subscriptionId, targetState, err)
		}
		actualState := sub.State
		return fmt.Errorf("waiting for Subscription %q to become %q, currently %q", subscriptionId, targetState, actualState)
	}

	return nil
}

func checkExistingAliases(ctx context.Context, client subscriptionAlias.AliasClient, subscriptionId string) (*string, int, error) {
	aliasList, err := client.List(ctx)
	if err != nil {
		return nil, len(*aliasList.Value), fmt.Errorf("could not List existing Subscription Aliases")
	}

	if aliasList.Value == nil {
		return nil, len(*aliasList.Value), fmt.Errorf("failed reading Subscription Alias list")
	}

	for _, v := range *aliasList.Value {
		if v.Properties != nil && v.Properties.SubscriptionID != nil && subscriptionId == *v.Properties.SubscriptionID {
			return v.Name, len(*aliasList.Value), nil
		}
	}
	return nil, len(*aliasList.Value), nil
}
