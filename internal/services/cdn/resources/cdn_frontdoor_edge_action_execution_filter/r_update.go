// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_execution_filter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionexecutionfilters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		// applying an execution filter to the Front Door edge can take more than 10 minutes
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.EdgeActionFiltersClient

			id, err := edgeactionexecutionfilters.ParseExecutionFilterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config EdgeActionExecutionFilterModel
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

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("cdn_frontdoor_edge_action_version_id") {
				payload.Properties.VersionId = config.CdnFrontDoorEdgeActionVersionId
			}

			if metadata.ResourceData.HasChange("header_name") {
				payload.Properties.ExecutionFilterIdentifierHeaderName = config.HeaderName
			}

			if metadata.ResourceData.HasChange("header_value") {
				payload.Properties.ExecutionFilterIdentifierHeaderValue = config.HeaderValue
			}

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
