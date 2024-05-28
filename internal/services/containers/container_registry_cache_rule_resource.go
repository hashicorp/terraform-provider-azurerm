// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

var _ sdk.Resource = ContainerRegistryCacheRule{}

type ContainerRegistryCacheRule struct{}

func (ContainerRegistryCacheRule) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "The name of the cache rule.",
			ValidateFunc: validate.ContainerRegistryCacheRuleName,
		},
		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: registries.ValidateRegistryID,
		},
		"credential_set_id": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Description: "The ARM resource ID of the credential store which is associated with the cache rule.",
		},
		"source_repo": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The full source repository path such as 'docker.io/library/ubuntu'.",
		},

		"target_repo": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The target repository namespace such as 'ubuntu'.",
		},
	}
}

func (ContainerRegistryCacheRule) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (ContainerRegistryCacheRule) ModelObject() interface{} {
	return nil
}

func (ContainerRegistryCacheRule) ResourceType() string {
	return "azurerm_container_registry_cache_rule"
}

func (r ContainerRegistryCacheRule) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			cacheRulesClient := metadata.Client.Containers.ContainerRegistryClient_v2023_07_01.CacheRules
			subscriptionId := metadata.Client.Account.SubscriptionId
			ctx, cancel := timeouts.ForCreate(metadata.Client.StopContext, metadata.ResourceData)

			var config ContainerRegistryCacheRuleModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			defer cancel()
			log.Printf("[INFO] preparing arguments for Container Registry Cache Rule creation.")

			registryId, err := registries.ParseRegistryID(metadata.ResourceData.Get("container_registry_id").(string))
			if err != nil {
				return err
			}

			id := cacherules.NewCacheRuleID(subscriptionId,
				registryId.ResourceGroupName,
				registryId.RegistryName,
				metadata.ResourceData.Get("name").(string),
			)


		existing, err := cacheRulesClient.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_container_registry_cache_rule", id.ID())
		}

		// TODO: make a check that the repo is available in the registry.
		// TODO: validate the source repo.

		parameters := cacherules.CacheRule{
			Name: &id.CacheRuleName,
			Properties: &cacherules.CacheRuleProperties{
				SourceRepository: pointer.To(config.SourceRepo),
				TargetRepository: pointer.To(config.TargetRepo),
			},
		}

		// Conditionally add CredentialSetResourceId if credentialSetId is not empty
		if config.CredentialSetId != "" {
			parameters.Properties.CredentialSetResourceId = pointer.To(config.CredentialSetId)
		}

		if err := cacheRulesClient.CreateThenPoll(ctx, id, parameters); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}


			metadata.SetID(id)

			return nil
		},
	}
}

func (ContainerRegistryCacheRule) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			cacheRulesClient := metadata.Client.Containers.ContainerRegistryClient_v2023_07_01.CacheRules
			ctx, cancel := timeouts.ForRead(metadata.Client.StopContext, metadata.ResourceData)
			defer cancel()

			id, err := cacherules.ParseCacheRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			resourceGroupName := id.ResourceGroupName
			registryName := id.RegistryName

			registryId := registries.NewRegistryID(subscriptionId, resourceGroupName, registryName)

			resp, err := cacheRulesClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] Container Registry Cache Rule %s was not found.", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving Container Registry Cache Rule %s: %+v", *id, err)
			}

			metadata.ResourceData.Set("name", id.CacheRuleName)
			metadata.ResourceData.Set("container_registry_id", registryId.ID())

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					metadata.ResourceData.Set("source_repo", properties.SourceRepository)
					metadata.ResourceData.Set("target_repo", properties.TargetRepository)
					metadata.ResourceData.Set("credential_set_id", properties.CredentialSetResourceId)
				}
			}

			return nil
		},
	}
}

func (r ContainerRegistryCacheRule) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			cacheRulesClient := metadata.Client.Containers.CacheRulesClient
			ctx, cancel := timeouts.ForRead(metadata.Client.StopContext, metadata.ResourceData)

			defer cancel()
			log.Printf("[INFO] preparing arguments for Container Registry Cache Rule update.")

			id, err := cacherules.ParseCacheRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			parameters := cacherules.CacheRuleUpdateParameters{}
			credentialSetId := metadata.ResourceData.Get("credential_set_id").(string)

			if credentialSetId != "" {
				parameters = cacherules.CacheRuleUpdateParameters{
					Properties: &cacherules.CacheRuleUpdateProperties{CredentialSetResourceId: &credentialSetId},
				}
			} else {
				//This is due to an issue with the Azure CacheRule API that prevents removing credentials
				return fmt.Errorf("Error on update: credential_set_id must not be empty: %s", id)
			}

			if err := cacheRulesClient.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (ContainerRegistryCacheRule) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			cacheRulesClient := metadata.Client.Containers.CacheRulesClient
			ctx, cancel := timeouts.ForDelete(metadata.Client.StopContext, metadata.ResourceData)
			defer cancel()

			id, err := cacherules.ParseCacheRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := cacheRulesClient.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting Container Registry Cache Rule %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (ContainerRegistryCacheRule) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cacherules.ValidateCacheRuleID
}
