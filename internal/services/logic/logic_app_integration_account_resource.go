// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogicAppIntegrationAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationAccountCreateUpdate,
		Read:   resourceLogicAppIntegrationAccountRead,
		Update: resourceLogicAppIntegrationAccountCreateUpdate,
		Delete: resourceLogicAppIntegrationAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := integrationaccounts.ParseIntegrationAccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(integrationaccounts.IntegrationAccountSkuNameBasic),
					string(integrationaccounts.IntegrationAccountSkuNameFree),
					string(integrationaccounts.IntegrationAccountSkuNameStandard),
				}, false),
			},

			"integration_service_environment_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: integrationaccounts.ValidateIntegrationAccountID,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceLogicAppIntegrationAccountCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := integrationaccounts.NewIntegrationAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account", id.ID())
		}
	}

	account := integrationaccounts.IntegrationAccount{
		Properties: &integrationaccounts.IntegrationAccountProperties{},
		Location:   utils.String(location.Normalize(d.Get("location").(string))),
		Sku: &integrationaccounts.IntegrationAccountSku{
			Name: integrationaccounts.IntegrationAccountSkuName(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("integration_service_environment_id"); ok {
		account.Properties.IntegrationServiceEnvironment = &integrationaccounts.ResourceReference{
			Id: utils.String(v.(string)),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, account); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceLogicAppIntegrationAccountRead(d, meta)
}

func resourceLogicAppIntegrationAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccounts.ParseIntegrationAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.IntegrationAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if model.Sku != nil {
			d.Set("sku_name", string(model.Sku.Name))
		}

		if props := model.Properties; props != nil {
			iseId := ""
			if props.IntegrationServiceEnvironment != nil && props.IntegrationServiceEnvironment.Id != nil {
				iseId = *props.IntegrationServiceEnvironment.Id
			}
			d.Set("integration_service_environment_id", iseId)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceLogicAppIntegrationAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccounts.ParseIntegrationAccountID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
