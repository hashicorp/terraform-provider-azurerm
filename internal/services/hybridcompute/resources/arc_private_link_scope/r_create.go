// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_private_link_scope

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (a Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateLinkScopeModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.HybridCompute.PrivateLinkScopesClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := privatelinkscopes.NewProviderPrivateLinkScopeID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(a.ResourceType(), id)
				}
			}

			properties := privatelinkscopes.HybridComputePrivateLinkScope{
				Location:   location.Normalize(model.Location),
				Name:       &model.Name,
				Tags:       &model.Tags,
				Properties: &privatelinkscopes.HybridComputePrivateLinkScopeProperties{},
			}

			publicNetwork := privatelinkscopes.PublicNetworkAccessTypeDisabled

			if model.PublicNetworkAccessEnabled {
				publicNetwork = privatelinkscopes.PublicNetworkAccessTypeEnabled
			}

			properties.Properties.PublicNetworkAccess = &publicNetwork

			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
