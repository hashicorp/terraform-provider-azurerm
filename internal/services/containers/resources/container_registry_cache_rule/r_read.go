// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_cache_rule

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Read() sdk.ResourceFunc {
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
