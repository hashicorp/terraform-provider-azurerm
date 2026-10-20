// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_webhook

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/webhooks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandWebhookPropertiesCreateParameters(d *pluginsdk.ResourceData) *webhooks.WebhookPropertiesCreateParameters {
	customHeaders := make(map[string]string)
	for k, v := range d.Get("custom_headers").(map[string]any) {
		customHeaders[k] = v.(string)
	}

	webhookProperties := webhooks.WebhookPropertiesCreateParameters{
		ServiceUri:    d.Get("service_uri").(string),
		CustomHeaders: &customHeaders,
		Actions:       expandWebhookActions(d),
		Scope:         pointer.To(d.Get("scope").(string)),
		Status:        pointer.ToEnum[webhooks.WebhookStatus](d.Get("status").(string)),
	}

	return &webhookProperties
}

func expandWebhookPropertiesUpdateParameters(d *pluginsdk.ResourceData) *webhooks.WebhookPropertiesUpdateParameters {
	customHeaders := make(map[string]string)
	for k, v := range d.Get("custom_headers").(map[string]any) {
		customHeaders[k] = v.(string)
	}

	webhookProperties := webhooks.WebhookPropertiesUpdateParameters{
		ServiceUri:    pointer.To(d.Get("service_uri").(string)),
		CustomHeaders: &customHeaders,
		Actions:       pointer.To(expandWebhookActions(d)),
		Scope:         pointer.To(d.Get("scope").(string)),
		Status:        pointer.ToEnum[webhooks.WebhookStatus](d.Get("status").(string)),
	}

	return &webhookProperties
}

func expandWebhookActions(d *pluginsdk.ResourceData) []webhooks.WebhookAction {
	actions := make([]webhooks.WebhookAction, 0)
	for _, action := range d.Get("actions").(*pluginsdk.Set).List() {
		actions = append(actions, webhooks.WebhookAction(action.(string)))
	}

	return actions
}
