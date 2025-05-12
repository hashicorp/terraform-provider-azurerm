// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSignalRSharedPrivateLinkResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSignalRSharedPrivateLinkCreateUpdate,
		Read:   resourceSignalRSharedPrivateLinkRead,
		Update: resourceSignalRSharedPrivateLinkCreateUpdate,
		Delete: resourceSignalrSharedPrivateLinkDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := signalr.ParseSharedPrivateLinkResourceIDInsensitively(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"signalr_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: signalr.ValidateSignalRID,
			},

			"sub_resource_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.PrivateLinkSubResourceName,
			},

			"target_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"request_message": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSignalRSharedPrivateLinkCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	signalrID, err := signalr.ParseSignalRID(d.Get("signalr_service_id").(string))
	if err != nil {
		return fmt.Errorf("parsing ID of %s: %+v", signalrID, err)
	}

	id := signalr.NewSharedPrivateLinkResourceID(subscriptionId, signalrID.ResourceGroupName, signalrID.SignalRName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.SharedPrivateLinkResourcesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_signalr_shared_private_link_resource", id.ID())
		}
	}

	parameters := signalr.SharedPrivateLinkResource{
		Properties: &signalr.SharedPrivateLinkResourceProperties{
			GroupId:               d.Get("sub_resource_name").(string),
			PrivateLinkResourceId: d.Get("target_resource_id").(string),
		},
	}

	requestMessage := d.Get("request_message").(string)
	if requestMessage != "" {
		parameters.Properties.RequestMessage = utils.String(requestMessage)
	}

	if err := client.SharedPrivateLinkResourcesCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating the shared private link for signalr %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSignalRSharedPrivateLinkRead(d, meta)
}

func resourceSignalRSharedPrivateLinkRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSharedPrivateLinkResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SharedPrivateLinkResourcesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving shared private link %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		d.Set("name", model.Name)
		d.Set("signalr_service_id", signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName).ID())

		if props := model.Properties; props != nil {
			d.Set("sub_resource_name", props.GroupId)
			d.Set("target_resource_id", props.PrivateLinkResourceId)

			if props.RequestMessage != nil {
				d.Set("request_message", props.RequestMessage)
			}

			status := string(*props.Status)
			d.Set("status", status)
		}
	}
	return nil
}

func resourceSignalrSharedPrivateLinkDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSharedPrivateLinkResourceIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	if err := client.SharedPrivateLinkResourcesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
