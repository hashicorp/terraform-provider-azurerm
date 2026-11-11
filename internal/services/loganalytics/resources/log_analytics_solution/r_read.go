// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_solution

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationsmanagement/2015-11-01-preview/solution"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.SolutionsClient

			id, err := solution.ParseSolutionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", metadata.ResourceData.Id(), err)
			}

			state := SolutionResourceModel{
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if loc := model.Location; loc != nil {
					state.Location = location.Normalize(*loc)
				}

				// Reversing the mapping used to get .solution_name
				// expecting resp.Name to be in format "SolutionName(WorkspaceName)".
				if v := model.Name; v != nil {
					val := pointer.From(v)
					segments := strings.Split(val, "(")
					if len(segments) != 2 {
						return fmt.Errorf("expected %q to match 'Solution(WorkspaceName)'", val)
					}

					state.SolutionName = segments[0]
					state.WorkspaceName = strings.TrimSuffix(segments[1], ")")
				}

				if props := model.Properties; props != nil {
					var workspaceId string
					if props.WorkspaceResourceId != "" {
						id, err := workspaces.ParseWorkspaceIDInsensitively(props.WorkspaceResourceId)
						if err != nil {
							return err
						}
						workspaceId = id.ID()
					}
					state.WorkspaceResourceId = workspaceId
				}

				if plan := model.Plan; plan != nil {
					state.SolutionPlan = flattenAzureRmLogAnalyticsSolutionPlan(plan)
				}

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenAzureRmLogAnalyticsSolutionPlan(input *solution.SolutionPlan) []SolutionPlanModel {
	output := make([]SolutionPlanModel, 0)
	if input == nil {
		return output
	}

	plan := SolutionPlanModel{}

	plan.Name = pointer.From(input.Name)
	plan.Product = pointer.From(input.Product)
	plan.PromotionCode = pointer.From(input.PromotionCode)
	plan.Publisher = pointer.From(input.Publisher)

	return append(output, plan)
}
