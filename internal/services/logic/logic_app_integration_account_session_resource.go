// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountsessions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppIntegrationAccountSession() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationAccountSessionCreateUpdate,
		Read:   resourceLogicAppIntegrationAccountSessionRead,
		Update: resourceLogicAppIntegrationAccountSessionCreateUpdate,
		Delete: resourceLogicAppIntegrationAccountSessionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := integrationaccountsessions.ParseSessionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"integration_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountName(),
			},

			"content": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},
		},
	}
}

func resourceLogicAppIntegrationAccountSessionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationAccountSessionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := integrationaccountsessions.NewSessionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account_session", id.ID())
		}
	}

	content := d.Get("content")

	parameters := integrationaccountsessions.IntegrationAccountSession{
		Properties: integrationaccountsessions.IntegrationAccountSessionProperties{
			Content: &content,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountSessionRead(d, meta)
}

func resourceLogicAppIntegrationAccountSessionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountSessionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountsessions.ParseSessionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SessionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if model := resp.Model; model != nil {
		props := model.Properties
		if props.Content != nil {
			d.Set("content", props.Content)
		}
	}

	return nil
}

func resourceLogicAppIntegrationAccountSessionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountSessionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountsessions.ParseSessionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
