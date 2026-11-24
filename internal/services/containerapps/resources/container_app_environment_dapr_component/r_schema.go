// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_dapr_component

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/daprcomponents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DaprComponentName,
			Description:  "The name for this Dapr component.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: daprcomponents.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Managed Environment for this Dapr Component.",
		},

		"component_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The Dapr Component Type. For example `state.azure.blobstorage`.",
		},

		"version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The version of the component.",
		},

		"init_timeout": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "5s",
			ValidateFunc: validate.InitTimeout,
			Description:  "The timeout for component initialisation as a `ISO8601` formatted string. e.g. `5s`, `2h`, `1m`.",
		},

		"ignore_errors": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Should the Dapr sidecar to continue initialisation if the component fails to load?",
		},

		"secret": helpers.SecretsSchema(),

		"metadata": helpers.ContainerAppEnvironmentDaprMetadataSchema(),

		"scopes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A list of scopes to which this component applies.",
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
