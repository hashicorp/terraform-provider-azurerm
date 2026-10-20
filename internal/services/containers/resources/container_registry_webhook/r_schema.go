// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_webhook

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/webhooks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func containerRegistryWebhookSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryWebhookName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"registry_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"service_uri": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ContainerRegistryWebhookServiceUri,
		},

		"custom_headers": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"status": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      webhooks.WebhookStatusEnabled,
			ValidateFunc: validation.StringInSlice(webhooks.PossibleValuesForWebhookStatus(), false),
		},

		"scope": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "",
		},

		"actions": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(webhooks.PossibleValuesForWebhookAction(), false),
			},
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}
