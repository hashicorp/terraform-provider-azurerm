// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.EdgeActionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config EdgeActionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := edgeactions.NewEdgeActionID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			// Edge Actions are global resources, the API only accepts `global` as the location
			param := edgeactions.EdgeAction{
				Location: "global",
				// `Standard` is currently the only SKU the API accepts, and `name` and `tier` must match
				Sku: edgeactions.SkuType{
					Name: "Standard",
					Tier: "Standard",
				},
				Tags: tags.Expand(config.Tags),
			}

			if err := client.CreateCallbackThenPoll(ctx, id, param, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
