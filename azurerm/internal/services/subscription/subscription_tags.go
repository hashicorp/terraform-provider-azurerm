package subscription

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-11-01/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func subscriptionTags() *schema.Resource {
	return &schema.Resource{
		Create: subscriptionTagCreateUpdate,
		Read:   subscriptionTagRead,
		Update: subscriptionTagCreateUpdate,
		Delete: subscriptionTagDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"tags": tags.Schema(),
			"subscription_id": {
				Type:         schema.TypeString,
				Description:  "The GUID of the Subscription.",
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func subscriptionTagCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	tagsClient := meta.(*clients.Client).Subscription.TagsClient
	client := meta.(*clients.Client).Subscription.Client
	subscriptionClient := meta.(*clients.Client).Subscription.SubscriptionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId := ""

	//verify existing subscription
	if subscriptionIdRaw, ok := d.GetOk("subscription_id"); ok {
		subscriptionId = subscriptionIdRaw.(string)

		locks.ByID(subscriptionId)
		defer locks.UnlockByID(subscriptionId)
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
			if err := waitForSubscriptionStateToSettleSub(ctx, meta.(*clients.Client), subscriptionId, "Active", createDeadline); err != nil {
				return fmt.Errorf("failed waiting for Subscription %q  to enter %q state: %+v", subscriptionId, "Active", err)
			}
		}
	}

	d.Set("subscription_id", subscriptionId)
	t := d.Get("tags").(map[string]interface{})
	resource_tags := resources.Tags{
		Tags: tags.Expand(t),
	}
	tagPatchParamter := resources.TagsPatchResource{Operation: "Merge", Properties: &resource_tags}
	uptags, urerr := tagsClient.UpdateAtScope(context.Background(), "subscriptions/"+subscriptionId, tagPatchParamter)
	if urerr != nil {
		if !utils.ResponseWasNotFound(uptags.Response) {
			return fmt.Errorf("Adding tag value %q for subscription %q: %+v", tags.Flatten(resource_tags.Tags), subscriptionId, urerr)
		}

	}
	d.SetId(*uptags.ID)

	return subscriptionTagRead(d, meta)
}

func subscriptionTagRead(d *schema.ResourceData, meta interface{}) error {
	tagsClient := meta.(*clients.Client).Subscription.TagsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := d.Get("subscription_id").(string)
	resp, err := tagsClient.GetAtScope(ctx, "subscriptions/"+subscriptionId)
	fmt.Println(resp)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Failed to Retrieve tags for subscription %q: %+v", subscriptionId, err)
		}

	}
	t := make(map[string]*string)
	t = *&resp.Properties.Tags
	d.Set("subscription_id", subscriptionId)

	return tags.FlattenAndSet(d, t)
}

func subscriptionTagDelete(d *schema.ResourceData, meta interface{}) error {
	tagsClient := meta.(*clients.Client).Subscription.TagsClient
	client := meta.(*clients.Client).Subscription.Client
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId := d.Get("subscription_id").(string)
	locks.ByID(subscriptionId)
	defer locks.UnlockByID(subscriptionId)
	_, serr := client.Get(ctx, subscriptionId)

	if serr != nil {
		return fmt.Errorf("Failed to read Subscription %q: %+v", subscriptionId, serr)
	}
	t := d.Get("tags").(map[string]interface{})
	resource_tags := resources.Tags{
		Tags: tags.Expand(t),
	}
	tagPatchParamter := resources.TagsPatchResource{Operation: "Delete", Properties: &resource_tags}
	_, delerr := tagsClient.UpdateAtScope(context.Background(), "subscriptions/"+subscriptionId, tagPatchParamter)
	if delerr != nil {
		return fmt.Errorf("Failed to Remove tags %q from subscription %q: %+v", tags.Flatten(resource_tags.Tags), subscriptionId, delerr)

	}

	return nil
}

func waitForSubscriptionStateToSettleSub(ctx context.Context, clients *clients.Client, subscriptionId string, targetState string, timeout time.Duration) error {
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
