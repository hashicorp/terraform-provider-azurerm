// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/sentinelonboardingstates"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SecurityInsightsSentinelOnboardingStateModel struct {
	ResourceGroupName         string `tfschema:"resource_group_name"`
	WorkspaceName             string `tfschema:"workspace_name"`
	CustomerManagedKeyEnabled bool   `tfschema:"customer_managed_key_enabled"`
	WorkspaceId               string `tfschema:"workspace_id"`
}

type LogAnalyticsWorkspaceOnboardResource struct{}

var _ sdk.Resource = LogAnalyticsWorkspaceOnboardResource{}

func (r LogAnalyticsWorkspaceOnboardResource) ResourceType() string {
	return "azurerm_sentinel_log_analytics_workspace_onboarding"
}

func (r LogAnalyticsWorkspaceOnboardResource) ModelObject() interface{} {
	return &SecurityInsightsSentinelOnboardingStateModel{}
}

func (r LogAnalyticsWorkspaceOnboardResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sentinelonboardingstates.ValidateOnboardingStateID
}

func (r LogAnalyticsWorkspaceOnboardResource) Arguments() map[string]*pluginsdk.Schema {
	out := map[string]*pluginsdk.Schema{
		"resource_group_name": {
			Deprecated:    "this property has been deprecated in favour of `workspace_id`",
			Type:          schema.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"workspace_id"},
			ValidateFunc:  resourcegroups.ValidateName,
		},

		"workspace_name": {
			Deprecated:    "this property will be removed in favour of `workspace_id` in version 4.0 of the AzureRM Provider",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"workspace_id"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		// lintignore:S013
		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     features.FourPointOhBeta(),
			Optional:     !features.FourPointOhBeta(),
			Computed:     !features.FourPointOhBeta(),
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"customer_managed_key_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},
	}

	if features.FourPointOhBeta() {
		delete(out, "resource_group_name")
		delete(out, "workspace_name")
	} else {
		out["workspace_id"].ConflictsWith = []string{"resource_group_name", "workspace_name"}
	}

	return out
}

func (r LogAnalyticsWorkspaceOnboardResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LogAnalyticsWorkspaceOnboardResource) Create() sdk.ResourceFunc {
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
			if model.WorkspaceId != "" {
				parsedWorkspaceId, err := workspaces.ParseWorkspaceID(model.WorkspaceId)
				if err != nil {
					return fmt.Errorf("parsing `log_analytics_workspace_id`: %+v", err)
				}
				id = sentinelonboardingstates.NewOnboardingStateID(parsedWorkspaceId.SubscriptionId, parsedWorkspaceId.ResourceGroupName, parsedWorkspaceId.WorkspaceName, "default")
			} else { // TODO: remove in 4.0
				subscriptionId := metadata.Client.Account.SubscriptionId
				id = sentinelonboardingstates.NewOnboardingStateID(subscriptionId, model.ResourceGroupName, model.WorkspaceName, "default")
			}

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &sentinelonboardingstates.SentinelOnboardingState{
				Properties: &sentinelonboardingstates.SentinelOnboardingStateProperties{
					CustomerManagedKey: &model.CustomerManagedKeyEnabled,
				},
			}

			if _, err := client.Create(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("context has no deadline")
			}

			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{"404"},
				Target:  []string{"200"},
				Refresh: func() (interface{}, string, error) {
					resp, err := client.Get(ctx, id)
					statusCode := "dropped connection"
					if resp.HttpResponse != nil {
						statusCode = strconv.Itoa(resp.HttpResponse.StatusCode)
					}

					if err != nil {
						if response.WasNotFound(resp.HttpResponse) {
							return resp, statusCode, nil
						}
						return resp, "", err
					}

					return resp, statusCode, nil
				},
				Timeout: time.Until(deadline),
				Delay:   15 * time.Second,
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be fully onboarded: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LogAnalyticsWorkspaceOnboardResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.OnboardingStatesClient

			id, err := sentinelonboardingstates.ParseOnboardingStateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID()

			state := SecurityInsightsSentinelOnboardingStateModel{
				ResourceGroupName: id.ResourceGroupName,
				WorkspaceName:     id.WorkspaceName,
				WorkspaceId:       workspaceId,
			}

			if properties := model.Properties; properties != nil {
				if properties.CustomerManagedKey != nil {
					state.CustomerManagedKeyEnabled = *properties.CustomerManagedKey
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LogAnalyticsWorkspaceOnboardResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.OnboardingStatesClient

			id, err := sentinelonboardingstates.ParseOnboardingStateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
