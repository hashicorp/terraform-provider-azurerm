// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_token

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/tokens"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an Azure Container Registry token associated to a scope map. For more information on scope maps and their tokens see the [product documentation](https://learn.microsoft.com/azure/container-registry/container-registry-repository-scoped-permissions).",
		Create:      resourceContainerRegistryTokenCreate,
		Read:        resourceContainerRegistryTokenRead,
		Update:      resourceContainerRegistryTokenUpdate,
		Delete:      resourceContainerRegistryTokenDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := tokens.ParseTokenID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: containerRegistryTokenSchema(),
	}
}
