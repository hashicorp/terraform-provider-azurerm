// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_cache_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
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
