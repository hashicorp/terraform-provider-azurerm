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
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/operation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryCacheRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceContainerRegistryCacheRuleCreate,
		Read:   resourceContainerRegistryCacheRuleRead,

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
	operationClient := meta.(*clients.Client).Containers.ContainerRegistryClient_v2023_07_01.Operation

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)

	defer cancel()
	log.Printf("[INFO] preparing arguments for Container Registry Cache Rule creation.")

	id := cacherules.NewCacheRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("registry").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := cacheRulesClient.Get(ctx, id)

		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_container_registry_cache_rule", id.ID())
		}

		sId := commonids.NewSubscriptionID(subscriptionId)
		availabilityRequest := operation.RegistryNameCheckRequest{
			Name: id.CacheRuleName,
			Type: "Microsoft.ContainerRegistry/registries/cacheRules",
		}

		availabilityResponse, err := operationClient.RegistriesCheckNameAvailability(ctx, sId, availabilityRequest)
		if err != nil {
			return fmt.Errorf("checking availability of %s: %s", id.CacheRuleName, err)
		}

		if availabilityResponse.Model == nil && availabilityResponse.Model.NameAvailable == nil {
			return fmt.Errorf("checking availability of %s: unexpected response from API", id.CacheRuleName)
		}

		if available := *availabilityResponse.Model.NameAvailable; !available {
			return fmt.Errorf("the name %q used for the Container Registry Cache Rule needs to be unique within Container Registry %s and isn't available: %s", id.CacheRuleName, id.RegistryName, *availabilityResponse.Model.Message)
		}

		// TODO: make a check that the repo is available in the registry.
		targetRepo := d.Get("target_repo").(string)

		// TODO: validate the source repo.
		sourceRepo := d.Get("source_repo").(string)

		parameters := cacherules.CacheRule{
			Name: &id.CacheRuleName,
			Properties: &cacherules.CacheRuleProperties{
				SourceRepository: &sourceRepo,
				TargetRepository: &targetRepo,
			},
		}

		if err := cacheRulesClient.CreateThenPoll(ctx, id, parameters); err != nil {
			return fmt.Errorf("creating Container Registry Cache Rule %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceContainerRegistryCacheRuleRead(d, meta)
}

func resourceContainerRegistryCacheRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	cacheRulesClient := meta.(*clients.Client).Containers.ContainerRegistryClient_v2023_07_01.CacheRules
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cacherules.ParseCacheRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := cacheRulesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Container Registry Cache Rule %s was not found.", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Container Registry Cache Rule %s: %+v", *id, err)
	}

	d.Set("name", id.CacheRuleName)
	d.Set("registry", id.RegistryName)

	if model := resp.Model; model != nil {
		if properties := model.Properties; properties != nil {
			d.Set("source_repo", properties.SourceRepository)
			d.Set("target_repo", properties.TargetRepository)
		}
	}

	return nil
}
