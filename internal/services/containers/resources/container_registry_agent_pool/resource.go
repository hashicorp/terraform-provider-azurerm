// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_agent_pool

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceContainerRegistryAgentPoolCreate,
		Read:   resourceContainerRegistryAgentPoolRead,
		Update: resourceContainerRegistryAgentPoolUpdate,
		Delete: resourceContainerRegistryAgentPoolDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := agentpools.ParseAgentPoolID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: containerRegistryAgentPoolSchema(),
	}
}
