// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action

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

			param := edgeactions.EdgeActionUpdate{}

			if metadata.ResourceData.HasChange("sku_name") {
				param.Sku = &edgeactions.SkuTypeUpdate{
					Name: pointer.To(config.SkuName),
					Tier: pointer.To(config.SkuName),
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				// the API rejects a `null` value for `tags`, so an empty map is sent to remove all tags
				if config.Tags == nil {
					config.Tags = map[string]string{}
				}
				param.Tags = pointer.To(config.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
