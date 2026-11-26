// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_scope_map

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func containerRegistryScopeMapSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the scope map.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryScopeMapName,
		},

		"description": {
			Description:  "The description of the Container Registry.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 256),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"container_registry_name": {
			Description:  "The name of the Container Registry.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"actions": {
			Description: "A list of actions to attach to the scope map (e.g. `repo/content/read`, `repo2/content/delete`).",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
