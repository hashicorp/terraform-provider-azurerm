// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/credentialsets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Description:  "The ARM resource ID of the credential store which is associated with the cache rule.",
			ValidateFunc: credentialsets.ValidateCredentialSetID,
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

type ContainerRegistryCacheRuleModel struct {
	Name                string `tfschema:"name"`
	ContainerRegistryId string `tfschema:"container_registry_id"`
	CredentialSetId     string `tfschema:"credential_set_id"`
	SourceRepo          string `tfschema:"source_repo"`
	TargetRepo          string `tfschema:"target_repo"`
}

func (ContainerRegistryCacheRule) ModelObject() interface{} {
	return &ContainerRegistryCacheRuleModel{}
}

func (ContainerRegistryCacheRule) ResourceType() string {
	return "azurerm_container_registry_cache_rule"
}

func (r ContainerRegistryCacheRule) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			cacheRulesClient := metadata.Client.Containers.CacheRulesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ContainerRegistryCacheRuleModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

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
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_container_registry_cache_rule", id.ID())
			}

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
			cacheRulesClient := metadata.Client.Containers.CacheRulesClient

			var config ContainerRegistryCacheRuleModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			id, err := cacherules.ParseCacheRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			registryId := registries.NewRegistryID(id.SubscriptionId, id.ResourceGroupName, id.RegistryName)

			resp, err := cacheRulesClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("%s was not found.", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			config.Name = id.CacheRuleName
			config.ContainerRegistryId = registryId.ID()

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					config.SourceRepo = pointer.From(properties.SourceRepository)
					config.TargetRepo = pointer.From(properties.TargetRepository)
					config.CredentialSetId = pointer.From(properties.CredentialSetResourceId)
				}
			}

			return metadata.Encode(&config)
		},
	}
}

func (r ContainerRegistryCacheRule) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			cacheRulesClient := metadata.Client.Containers.CacheRulesClient

			var config ContainerRegistryCacheRuleModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			log.Printf("[INFO] preparing arguments for Container Registry Cache Rule update.")

			id, err := cacherules.ParseCacheRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			credentialSetId := pointer.To(config.CredentialSetId)

			parameters := cacherules.CacheRuleUpdateParameters{
				Properties: &cacherules.CacheRuleUpdateProperties{CredentialSetResourceId: credentialSetId},
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

			id, err := cacherules.ParseCacheRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := cacheRulesClient.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (ContainerRegistryCacheRule) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cacherules.ValidateCacheRuleID
}

func (ContainerRegistryCacheRule) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldVal, newVal := metadata.ResourceDiff.GetChange("credential_set_id"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("credential_set_id"); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
