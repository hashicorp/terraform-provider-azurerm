package subscription

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/subscription/mgmt/2020-09-01/subscription"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSubscriptionAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSubscriptionAliasCreateUpdate,
		Read:   resourceArmSubscriptionAliasRead,
		Update: resourceArmSubscriptionAliasCreateUpdate,
		Delete: resourceArmSubscriptionAliasDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SubscriptionAliasID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SubscriptionAliasName,
			},

			"subscription_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceArmSubscriptionAliasCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Subscription.AliasClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Subscription Alias %q : %+v", name, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_subscription_alias", *existing.ID)
		}
	}

	alias := subscription.PutAliasRequest{
		Properties: &subscription.PutAliasRequestProperties{
			SubscriptionID: utils.String(d.Get("subscription_id").(string)),
		},
	}

	// despite this function only says `Create`, but it also works during updating
	future, err := client.Create(ctx, name, alias)
	if err != nil {
		return fmt.Errorf("creating/updating Subscription Alias %q : %+v", name, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating future for Subscription Alias %q : %+v", name, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("retrieving Subscription Alias %q : %+v", name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Subscription Alias %q  ID", name)
	}

	d.SetId(*resp.ID)
	return resourceArmSubscriptionAliasRead(d, meta)
}

func resourceArmSubscriptionAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Subscription.AliasClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionAliasID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] subscription %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Subscription Alias %q : %+v", id.Name, err)
	}

	d.Set("name", id.Name)
	if props := resp.Properties; props != nil {
		d.Set("subscription_id", props.SubscriptionID)
	}

	return nil
}

func resourceArmSubscriptionAliasDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Subscription.AliasClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionAliasID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.Name); err != nil {
		return fmt.Errorf("deleting Subscription Alias %q : %+v", id.Name, err)
	}

	return nil
}
