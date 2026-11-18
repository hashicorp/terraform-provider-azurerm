// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_version

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		// deploying the version code can take more than 10 minutes
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.EdgeActionVersionsClient

			id, err := edgeactionversions.ParseVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config EdgeActionVersionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing, err := client.Get(ctx, *id)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}

				if existing.Model == nil {
					return fmt.Errorf("retrieving %s: `model` was nil", id)
				}

				payload := existing.Model
				payload.Tags = tags.Expand(config.Tags)

				// `Create` is the SDK's name for the `PUT` operation, which is used rather than the `PATCH` (`Update`)
				// so that the full model is sent back and fields removed from the configuration are cleared
				if err := client.CreateThenPoll(ctx, *id, *payload); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChanges("code_content_base64", "code_file_name") && config.CodeContentBase64 != "" {
				if err := deployVersionCode(ctx, client, *id, config); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
