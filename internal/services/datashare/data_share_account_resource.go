// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datashare

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	helperTags "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataShareAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataShareAccountCreate,
		Read:   resourceDataShareAccountRead,
		Update: resourceDataShareAccountUpdate,
		Delete: resourceDataShareAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := account.ParseAccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

			// the api will save and return the tag keys in lowercase, so an extra validation of the key is all in lowercase is added
			// issue has been created https://github.com/Azure/azure-rest-api-specs/issues/9280
			"tags": tags.SchemaEnforceLowerCaseKeys(),
		},
	}
}

func resourceDataShareAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.AccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := account.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_data_share_account", id.ID())
	}

	expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	account := account.Account{
		Name:     utils.String(id.AccountName),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: *expandedIdentity,
		Tags:     helperTags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateThenPoll(ctx, id, account); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataShareAccountRead(d, meta)
}

func resourceDataShareAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if err := d.Set("identity", identity.FlattenSystemAssigned(&model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
		return helperTags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceDataShareAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.AccountClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	props := account.AccountUpdateParameters{}

	if d.HasChange("tags") {
		props.Tags = helperTags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err = client.Update(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceDataShareAccountRead(d, meta)
}

func resourceDataShareAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.AccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
