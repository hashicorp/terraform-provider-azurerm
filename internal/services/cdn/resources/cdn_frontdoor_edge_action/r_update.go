// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.EdgeActionsClient

			id, err := edgeactions.ParseEdgeActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config EdgeActionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("sku") {
				payload.Sku = edgeactions.SkuType{
					Name: config.Sku,
					Tier: config.Sku,
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				// the API rejects a `null` value for `tags`, so an empty map is sent to remove all tags
				if config.Tags == nil {
					config.Tags = map[string]string{}
				}
				payload.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
