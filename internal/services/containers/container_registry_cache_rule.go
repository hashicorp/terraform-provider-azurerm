// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/registries"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceContainerRegistryCacheRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := registries.ParseRegistryID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceContainerRegistrySchema(),
	}
}

func resourceContainerRegistryCacheRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			Description: "The name of the cache rule.",
		},

		"registry": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Description:  "The name of the container registry.",
			ValidateFunc: containerValidate.ContainerRegistryName,
		},

		"source_repo": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			Description: "The full source repository path such as 'docker.io/library/ubuntu'.",
		},

		"target_repo": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			Description: "The target repository namespace such as 'ubuntu'.",
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}
