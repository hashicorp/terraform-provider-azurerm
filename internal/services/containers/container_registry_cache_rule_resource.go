// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryCacheRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cacherules.ParseCacheRuleID(id)
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

func resourceContainerRegistryCacheRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	cacheRulesClient := meta.(*clients.Client).Containers.ContainerRegistryClient_v2023_07_01.CacheRules
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)

	defer cancel()
	log.Printf("[INFO] preparing arguments for Container Registry Cache Rule creation.")

	id := cacherules.NewCacheRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("registry").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)

		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_container_registry_cache_rule", id.ID())
		}

		sId := commonids.NewSubscriptionID(subscriptionId)

	}
}
