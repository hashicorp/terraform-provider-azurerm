// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_cache_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			cacheRulesClient := metadata.Client.Containers.CacheRulesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ContainerRegistryCacheRuleModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			registryId, err := registries.ParseRegistryID(metadata.ResourceData.Get("container_registry_id").(string))
			if err != nil {
				return err
			}

			id := cacherules.NewCacheRuleID(
				subscriptionId,
				registryId.ResourceGroupName,
				registryId.RegistryName,
				metadata.ResourceData.Get("name").(string),
			)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := cacheRulesClient.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_container_registry_cache_rule", id.ID())
				}
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

			if err := cacheRulesClient.CreateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
