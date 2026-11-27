// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	cdnClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := cdnClient.From(metadata.Client).EdgeActionsClient

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

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(config.Tags)
			}

			// `Create` is the SDK's name for the `PUT` operation, which is used rather than the `PATCH` (`Update`)
			// so that the full model is sent back and fields removed from the configuration are cleared
			if err := client.CreateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
