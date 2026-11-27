// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_webhook

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/webhooks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func containerRegistryWebhookSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Container Registry Webhook. Only Alphanumeric characters allowed.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryWebhookName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"registry_name": {
			Description:  "The Name of Container registry this Webhook belongs to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"service_uri": {
			Description:  "The service URI for the Webhook to post notifications.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ContainerRegistryWebhookServiceUri,
		},

		"custom_headers": {
			Description: "Custom headers that will be added to the webhook notifications request.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"enabled": {
			Description: "Whether this Webhook triggers notifications.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"scope": {
			Description: "The scope of repositories that can trigger an event. For example, `foo:*` means events for all tags under repository `foo`. `foo:bar` means events for `foo:bar` only. `foo` is equivalent to `foo:latest`. Empty means all events.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "",
		},

		"actions": {
			Description: "A list of actions that trigger the Webhook to post notifications. At least one action needs to be specified.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(webhooks.PossibleValuesForWebhookAction(), false),
			},
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}

	if !features.SixPointOh() {
		s["status"] = &pluginsdk.Schema{
			Description:   "Whether this Webhook triggers notifications.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ValidateFunc:  validation.StringInSlice(webhooks.PossibleValuesForWebhookStatus(), false),
			ConflictsWith: []string{"enabled"},
			Deprecated:    "`status` has been deprecated in favour of `enabled` and will be removed in v6.0 of the AzureRM Provider",
		}

		s["enabled"] = &pluginsdk.Schema{
			Description:   "Whether this Webhook triggers notifications.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"status"},
		}
	}

	return s
}
