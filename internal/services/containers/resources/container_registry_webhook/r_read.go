// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_webhook

import (
	"fmt"
	"log"
	"maps"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/webhooks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryWebhookRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient.WebHooks
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webhooks.ParseWebHookID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	callbackConfig, err := client.GetCallbackConfig(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Callback Config for %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("registry_name", id.RegistryName)
	d.Set("name", id.WebHookName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			status := ""
			if v := props.Status; v != nil {
				status = string(*v)
			}
			d.Set("status", status)

			d.Set("scope", pointer.From(props.Scope))

			webhookActions := make([]string, len(props.Actions))
			for i, action := range props.Actions {
				webhookActions[i] = string(action)
			}
			d.Set("actions", webhookActions)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	if callbackModel := callbackConfig.Model; callbackModel != nil {
		d.Set("service_uri", callbackModel.ServiceUri)

		customHeaders := make(map[string]string)
		if callbackModel.CustomHeaders != nil {
			maps.Copy(customHeaders, *callbackModel.CustomHeaders)
		}
		d.Set("custom_headers", customHeaders)
	}
	return nil
}
