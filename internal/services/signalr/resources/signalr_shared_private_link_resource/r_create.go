// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_shared_private_link_resource

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSignalRSharedPrivateLinkCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
	}

	parameters := signalr.SharedPrivateLinkResource{
		Properties: &signalr.SharedPrivateLinkResourceProperties{
			GroupId:               d.Get("sub_resource_name").(string),
			PrivateLinkResourceId: d.Get("target_resource_id").(string),
		},
	}

	requestMessage := d.Get("request_message").(string)
	if requestMessage != "" {
		parameters.Properties.RequestMessage = pointer.To(requestMessage)
	}

	if d.IsNewResource() {
		if err := client.SharedPrivateLinkResourcesCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.SharedPrivateLinkResourcesCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceSignalRSharedPrivateLinkRead(d, meta)
}
