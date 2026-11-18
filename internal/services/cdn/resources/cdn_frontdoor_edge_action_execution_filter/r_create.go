// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_execution_filter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionexecutionfilters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.EdgeActionFiltersClient

			var config EdgeActionExecutionFilterModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			edgeActionId, err := edgeactionexecutionfilters.ParseEdgeActionID(config.CdnFrontDoorEdgeActionId)
			if err != nil {
				return err
			}

			id := edgeactionexecutionfilters.NewExecutionFilterID(edgeActionId.SubscriptionId, edgeActionId.ResourceGroupName, edgeActionId.EdgeActionName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			// Edge Action Execution Filters are global resources, the API only accepts `global` as the location
			param := edgeactionexecutionfilters.EdgeActionExecutionFilter{
				Location: "global",
				Properties: &edgeactionexecutionfilters.EdgeActionExecutionFilterProperties{
					VersionId:                            config.CdnFrontDoorEdgeActionVersionId,
					ExecutionFilterIdentifierHeaderName:  config.HeaderName,
					ExecutionFilterIdentifierHeaderValue: config.HeaderValue,
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
