// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_webhook

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/webhooks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryWebhookCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient.WebHooks
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webhooks.NewWebHookID(subscriptionId, d.Get("resource_group_name").(string), d.Get("registry_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_container_registry_webhook", id.ID())
		}
	}

	webhook := webhooks.WebhookCreateParameters{
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: expandWebhookPropertiesCreateParameters(d),
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateCallbackThenPoll(ctx, id, webhook, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryWebhookRead(d, meta)
}
