// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_log_analytics_workspace_onboarding

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/sentinelonboardingstates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SecurityInsightsSentinelOnboardingStateModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Sentinel.OnboardingStatesClient
			// the service only support `default` state
			var id sentinelonboardingstates.OnboardingStateId
			parsedWorkspaceId, err := workspaces.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing `log_analytics_workspace_id`: %+v", err)
			}
			id = sentinelonboardingstates.NewOnboardingStateID(parsedWorkspaceId.SubscriptionId, parsedWorkspaceId.ResourceGroupName, parsedWorkspaceId.WorkspaceName, "default")

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &sentinelonboardingstates.SentinelOnboardingState{
				Properties: &sentinelonboardingstates.SentinelOnboardingStateProperties{
					CustomerManagedKey: &model.CustomerManagedKeyEnabled,
				},
			}

			if _, err := client.Create(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			poller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
				resp, err := client.Get(pollerCtx, id)
				return resp.HttpResponse, err
			}, &custompollers.EventualConsistencyPollerOptions{
				Interval:              10 * time.Second,
				TargetStatusCode:      pointer.To(http.StatusOK),
				RetryErrorStatusCodes: []int{http.StatusNotFound},
			})
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be fully onboarded: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
