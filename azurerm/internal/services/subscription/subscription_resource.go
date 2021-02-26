package subscription

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-11-01/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/subscription/mgmt/2020-09-01/subscription"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	azValidate "github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var (
	billingScopeEnrollmentFmt = "/providers/Microsoft.Billing/billingAccounts/%s/enrollmentAccounts/%s"
	billingScopeMCAFmt        = "/providers/Microsoft.Billing/billingAccounts/%s/billingProfiles/%s/invoiceSections/%s"
)

func resourceSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubscriptionCreate,
		Update: resourceSubscriptionUpdate,
		Read:   resourceSubscriptionRead,
		Delete: resourceSubscriptionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SubscriptionAliasID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"subscription_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Display Name for the Subscription.",
				ValidateFunc: validate.SubscriptionName,
			},

			"alias": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The Alias Name of the subscription. If omitted a new UUID will be generated for this property.",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Active",
				ValidateFunc: validation.StringInSlice([]string{
					"Active",
					"Cancelled",
				}, false),
			},

			"billing_account": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				Description:  "The name of the billing account under which the Subscription will be created.",
				ValidateFunc: azValidate.NoEmptyStrings, // TODO - Put some actual protection around this?
			},

			"enrollment_account": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The name of the enrollment account in which to create the subscription. Used for EA accounts.",
				ConflictsWith: []string{
					"invoice_section",
					"billing_profile",
					"subscription_id",
				},
				ValidateFunc: azValidate.NoEmptyStrings, // TODO - Put some actual protection around this?
			},

			"billing_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The name of the Billing Profile under which the Subscription should be created. Used for MCA and Partner Agreements.",
				ConflictsWith: []string{
					"enrollment_account",
					"subscription_id",
				},
				RequiredWith: []string{
					"invoice_section",
				},
				ValidateFunc: azValidate.NoEmptyStrings, // TODO - Put some actual protection around this?
			},

			"invoice_section": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The Invoice Section of the Billing Profile which will be used for the subscription. Used for MCA and Partner Agreements.",
				ConflictsWith: []string{
					"enrollment_account",
					"subscription_id",
				},
				RequiredWith: []string{
					"billing_profile",
				},
				ValidateFunc: azValidate.NoEmptyStrings, // TODO - Put some actual protection around this?
			},

			// Optional
			"workload": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The workload type for the Subscription. Possible values are `Production` (default) and `DevTest`.",
				// Other RP's have updated Constants with contextual prefixes so these are likely to change
				ValidateFunc: validation.StringInSlice([]string{
					string(subscription.Production),
					string(subscription.DevTest),
				}, false),
				// Workload is not exposed in any way, so must be ignored if the resource is imported.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == ""
				},
			},

			"subscription_id": {
				Type:        schema.TypeString,
				Description: "The GUID of the Subscription.",
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Description: "The Tenant ID to which the subscription belongs",
				Computed:    true,
			},
		},
	}
}

func resourceSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Subscription.AliasClient
	subscriptionClient := meta.(*clients.Client).Subscription.SubscriptionClient
	subscriptionsClient := meta.(*clients.Client).Subscription.Client
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

	existing, err := client.Get(ctx, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existence of Subscription by Alias %q: %+v", id.Name, err)
		}
	}

	if props := existing.Properties; props != nil {
		return tf.ImportAsExistsError("azurerm_subscription", id.ID())
	}

	workload := subscription.Production
	if workloadRaw := d.Get("workload").(string); workloadRaw != "" {
		workload = subscription.Workload(workloadRaw)
	}

	req := subscription.PutAliasRequest{
		Properties: &subscription.PutAliasRequestProperties{
			Workload: workload,
		},
	}

	targetState := d.Get("state").(string)
	reactivate := false
	subscriptionId := ""

	// Check if we're adding alias management for an existing subscription
	if subscriptionIdRaw, ok := d.GetOk("subscription_id"); ok {
		subscriptionId = subscriptionIdRaw.(string)
		exists, err := checkExistingAliases(ctx, *client, subscriptionId)
		if err != nil {
			return err
		}
		if exists != nil {
			return fmt.Errorf("An Alias for Subscription %q already exists with name %q - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information.", subscriptionId, *exists, "azurerm_subscription")
		}

		req.Properties.SubscriptionID = utils.String(subscriptionId)
		existingSub, err := subscriptionsClient.Get(ctx, subscriptionId)
		if err != nil {
			return fmt.Errorf("could not read existing Subscription %q", subscriptionId)
		}
		// Disabled and Warned are both "effectively" cancelled states,
		if (existingSub.State == subscriptions.Disabled || existingSub.State == subscriptions.Warned) && targetState == "Active" {
			log.Printf("[DEBUG] Existing subscription already in use in Disabled state Terraform will attempt to re-activate it")
			reactivate = true
		}
	} else {
		// If we're not assuming control of an existing Subscription, we need to know where to create it.
		req.Properties.DisplayName = utils.String(d.Get("subscription_name").(string))

		billingAccount := d.Get("billing_account").(string)

		if enrollmentAccount, ok := d.GetOk("enrollment_account"); ok && enrollmentAccount.(string) != "" {
			req.Properties.BillingScope = utils.String(fmt.Sprintf(billingScopeEnrollmentFmt, billingAccount, enrollmentAccount))
		} else {
			billingProfile := d.Get("billing_profile").(string)
			invoiceSection := d.Get("invoice_section").(string)
			req.Properties.BillingScope = utils.String(fmt.Sprintf(billingScopeMCAFmt, billingAccount, billingProfile, invoiceSection))
		}
	}

	future, err := client.Create(ctx, aliasName, req)
	if err != nil {
		return fmt.Errorf("creating new Subscription (Alias %q): %+v", aliasName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Subscription with Alias %q: %+v", id.Name, err)
	}

	if reactivate {
		if _, err := subscriptionClient.Enable(ctx, subscriptionId); err != nil {
			return fmt.Errorf("enabling Subscription %q: %+v", subscriptionId, err)
		}
	}

	alias, err := client.Get(ctx, id.Name)
	if err != nil || alias.Properties == nil || alias.Properties.SubscriptionID == nil {
		return fmt.Errorf("failed reading subscription details for Alias %q: %+v", id.Name, err)
	}

	if err := waitForSubscriptionStateToSettle(ctx, meta.(*clients.Client), *alias.Properties.SubscriptionID, targetState, d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("failed waiting for Subscription %q (Alias %q) to enter %q state: %+v", subscriptionId, id.Name, targetState, err)
	}

	d.SetId(id.ID())

	return resourceSubscriptionRead(d, meta)
}

func resourceSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Subscription.AliasClient
	subscriptionClient := meta.(*clients.Client).Subscription.SubscriptionClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionAliasID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name)
	if err != nil || resp.Properties == nil {
		return fmt.Errorf("could not read Subscription Alias for update: %+v", err)
	}

	subscriptionId := resp.Properties.SubscriptionID
	if subscriptionId == nil || *subscriptionId == "" {
		return fmt.Errorf("could not read Subscription ID from Alias")
	}
	if d.HasChange("subscription_name") {
		displayName := subscription.Name{
			SubscriptionName: utils.String(d.Get("subscription_name").(string)),
		}
		if _, err := subscriptionClient.Rename(ctx, *subscriptionId, displayName); err != nil {
			return fmt.Errorf("could not update Display Name of Subscription %q: %+v", *subscriptionId, err)
		}
	}

	if d.HasChange("state") {
		newState := d.Get("state").(string)
		switch newState {
		case "Active":
			if _, err := subscriptionClient.Enable(ctx, *subscriptionId); err != nil {
				return fmt.Errorf("failed to Enable Subscription %q: %+v", *subscriptionId, err)
			}
		case "Cancelled":
			if _, err := subscriptionClient.Cancel(ctx, *subscriptionId); err != nil {
				return fmt.Errorf("failed to Disable Subscription %q: %+v", *subscriptionId, err)
			}
		default:
			return fmt.Errorf("unsupported Subscription State %q", newState)
		}

		if err := waitForSubscriptionStateToSettle(ctx, meta.(*clients.Client), *subscriptionId, newState, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("failed to set Subscription %q (Alias %q) to %q: %+v", *subscriptionId, id.Name, newState, err)
		}
	}

	return nil
}

func resourceSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Subscription.AliasClient
	subscriptionsClient := meta.(*clients.Client).Subscription.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionAliasID(d.Id())
	if err != nil {
		return err
	}
	d.Set("alias", id.Name)

	alias, err := client.Get(ctx, id.Name)
	if err != nil {
		return fmt.Errorf("reading Subscription Aliss %q: %+v", id.Name, err)
	}

	subscriptionId := ""
	subscriptionName := ""
	tenantId := ""
	if props := alias.Properties; props != nil && props.SubscriptionID != nil {
		subscriptionId = *props.SubscriptionID
		resp, err := subscriptionsClient.Get(ctx, subscriptionId)

		if err != nil {
			return fmt.Errorf("failed to read Subscription %q (Alias %q) for Tenant Information: %+v", subscriptionId, id.Name, err)
		}
		if resp.TenantID != nil {
			tenantId = *resp.TenantID
		}

		if resp.DisplayName != nil {
			subscriptionName = *resp.DisplayName
		}

		state := ""
		if resp.State == subscriptions.Disabled || resp.State == subscriptions.Warned {
			state = "Cancelled"
		}

		if resp.State == subscriptions.Enabled {
			state = "Active"
		}

		d.Set("state", state)
	}

	// (@jackofallops) A subscription's billing scope is not exposed in any way in the API/SDK so we cannot read it back here

	d.Set("subscription_id", subscriptionId)
	d.Set("subscription_name", subscriptionName)
	d.Set("tenant_id", tenantId)

	return nil
}

// (@jackofallops) - Delete here is a misnomer.  The nature of subscriptions is such that they are never truly deleted
// Deleted here means Cancelled. Cancelled subscriptions are held in this state for  90 days before being purged from
// active systems.  However, the backend billing systems _never_ remove this data, so once a Subscription ID has been
// used and purged from active use it can never be recovered nor the UUID reused.
// Note Cancelling a Subscription leaves it in one of several states, `Disabled` for a Subscription with no Resources or
// Alias assignments, `Warned` for Cancelled with "something" associated with it.
func resourceSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Subscription.AliasClient
	subscriptionClient := meta.(*clients.Client).Subscription.SubscriptionClient
	subscriptionsClient := meta.(*clients.Client).Subscription.Client
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionAliasID(d.Id())
	if err != nil {
		return err
	}

	// Get subscription details for later
	alias, err := client.Get(ctx, id.Name)
	if err != nil || alias.Properties == nil {
		return fmt.Errorf("could not read Alias %q for Subscription: %+v", id.Name, err)
	}
	subscriptionId := ""
	if subscriptionIdRaw := alias.Properties.SubscriptionID; subscriptionIdRaw != nil {
		subscriptionId = *subscriptionIdRaw
	}

	sub, err := subscriptionsClient.Get(ctx, subscriptionId)
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
	resp, err := client.Delete(ctx, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("could not delete Alias %q for Subscription %q (ID: %q): %+v", id.Name, subscriptionName, subscriptionId, err)
		}
	}

	// Cancel the Subscription
	if _, err := subscriptionClient.Cancel(ctx, subscriptionId); err != nil {
		return fmt.Errorf("failed to cancel Subscription: %+v", err)
	}

	if err := waitForSubscriptionStateToSettle(ctx, meta.(*clients.Client), subscriptionId, "Cancelled", d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("failed to cancel Subscription %q (Alias %q): %+v", subscriptionId, id.Name, err)
	}

	return nil
}

func waitForSubscriptionStateToSettle(ctx context.Context, clients *clients.Client, subscriptionId string, targetState string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
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

	if actual, err := stateConf.WaitForState(); err != nil {
		sub, ok := actual.(subscriptions.Subscription)
		if !ok {
			return fmt.Errorf("failure in parsing response while waiting for Subscription %q to become %q: %+v", subscriptionId, targetState, err)
		}
		actualState := sub.State
		return fmt.Errorf("waiting for Subscription %q to become %q, currently %q", subscriptionId, targetState, actualState)
	}

	return nil
}

func checkExistingAliases(ctx context.Context, client subscription.AliasClient, subscriptionId string) (*string, error) {
	aliasList, err := client.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not List existing Subscription Aliases")
	}

	for _, v := range *aliasList.Value {
		if v.Properties != nil && v.Properties.SubscriptionID != nil && subscriptionId == *v.Properties.SubscriptionID {
			return v.Name, nil
		}
	}
	return nil, nil
}
