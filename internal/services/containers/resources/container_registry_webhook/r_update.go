// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_webhook

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/webhooks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryWebhookUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient.WebHooks
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webhooks.ParseWebHookID(d.Id())
	if err != nil {
		return err
	}

	webhook := webhooks.WebhookUpdateParameters{
		Properties: expandWebhookPropertiesUpdateParameters(d),
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.UpdateThenPoll(ctx, *id, webhook); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceContainerRegistryWebhookRead(d, meta)
}
