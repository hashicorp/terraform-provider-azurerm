// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_scope_map

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/scopemaps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an Azure Container Registry scope map. For more information on scope maps see the [product documentation](https://learn.microsoft.com/azure/container-registry/container-registry-repository-scoped-permissions).",
		Create:      resourceContainerRegistryScopeMapCreate,
		Read:        resourceContainerRegistryScopeMapRead,
		Update:      resourceContainerRegistryScopeMapUpdate,
		Delete:      resourceContainerRegistryScopeMapDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := scopemaps.ParseScopeMapID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: containerRegistryScopeMapSchema(),
	}
}
