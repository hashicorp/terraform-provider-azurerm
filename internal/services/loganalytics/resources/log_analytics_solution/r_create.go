// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_solution

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationsmanagement/2015-11-01-preview/solution"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.SolutionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config SolutionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			// The resource requires both .name and .plan.name are set in the format
			// "SolutionName(WorkspaceName)". Feedback will be submitted to the OMS team as IMO this isn't ideal.
			id := solution.NewSolutionID(subscriptionId, config.ResourceGroupName, fmt.Sprintf("%s(%s)", config.SolutionName, config.WorkspaceName))

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %s", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_log_analytics_solution", id.ID())
				}
			}

			workspaceID, err := workspaces.ParseWorkspaceID(config.WorkspaceResourceId)
			if err != nil {
				return err
			}

			parameters := solution.Solution{
				Name:     pointer.To(id.SolutionName),
				Location: pointer.To(location.Normalize(config.Location)),
				Properties: &solution.SolutionProperties{
					WorkspaceResourceId: workspaceID.ID(),
				},
				Tags: pointer.To(config.Tags),
			}

			if len(config.SolutionPlan) > 0 {
				solutionPlan := expandAzureRmLogAnalyticsSolutionPlan(config.SolutionPlan)
				solutionPlan.Name = &id.SolutionName
				parameters.Plan = &solutionPlan
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandAzureRmLogAnalyticsSolutionPlan(plans []SolutionPlanModel) solution.SolutionPlan {
	if len(plans) == 0 {
		return solution.SolutionPlan{}
	}

	plan := plans[0]
	name := plan.Name
	publisher := plan.Publisher
	promotionCode := plan.PromotionCode
	product := plan.Product

	return solution.SolutionPlan{
		Name:          pointer.To(name),
		PromotionCode: pointer.To(promotionCode),
		Publisher:     pointer.To(publisher),
		Product:       pointer.To(product),
	}
}
