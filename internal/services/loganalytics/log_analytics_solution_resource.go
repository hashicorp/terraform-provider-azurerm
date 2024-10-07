// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationsmanagement/2015-11-01-preview/solution"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogAnalyticsSolutionResource struct{}

func (s LogAnalyticsSolutionResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.SolutionV0ToV1{},
		},
	}
}

type SolutionResourceModel struct {
	SolutionName        string              `tfschema:"solution_name"`
	WorkspaceName       string              `tfschema:"workspace_name"`
	WorkspaceResourceId string              `tfschema:"workspace_resource_id"`
	Location            string              `tfschema:"location"`
	ResourceGroupName   string              `tfschema:"resource_group_name"`
	SolutionPlan        []SolutionPlanModel `tfschema:"plan"`
	Tags                map[string]string   `tfschema:"tags"`
}

type SolutionPlanModel struct {
	Name          string `tfschema:"name"`
	Publisher     string `tfschema:"publisher"`
	PromotionCode string `tfschema:"promotion_code"`
	Product       string `tfschema:"product"`
}

var _ sdk.ResourceWithUpdate = LogAnalyticsSolutionResource{}
var _ sdk.ResourceWithStateMigration = LogAnalyticsSolutionResource{}

func (s LogAnalyticsSolutionResource) ModelObject() interface{} {
	return &SolutionResourceModel{}
}

func (s LogAnalyticsSolutionResource) ResourceType() string {
	return "azurerm_log_analytics_solution"
}

func (s LogAnalyticsSolutionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return solution.ValidateSolutionID
}

func (s LogAnalyticsSolutionResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"solution_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"workspace_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LogAnalyticsWorkspaceName,
		},

		"workspace_resource_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"location": commonschema.Location(),

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"plan": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"publisher": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
					"promotion_code": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"product": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}

func (s LogAnalyticsSolutionResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}
func (s LogAnalyticsSolutionResource) Create() sdk.ResourceFunc {
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

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_log_analytics_solution", id.ID())
			}

			workspaceID, err := workspaces.ParseWorkspaceID(config.WorkspaceResourceId)
			if err != nil {
				return err
			}

			parameters := solution.Solution{
				Name:     pointer.To(id.SolutionName),
				Location: pointer.To(azure.NormalizeLocation(config.Location)),
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

			err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (s LogAnalyticsSolutionResource) Read() sdk.ResourceFunc {
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
				if location := model.Location; location != nil {
					state.Location = azure.NormalizeLocation(*location)
				}

				// Reversing the mapping used to get .solution_name
				// expecting resp.Name to be in format "SolutionName(WorkspaceName)".
				if v := model.Name; v != nil {
					val := pointer.From(v)
					segments := strings.Split(val, "(")
					if len(segments) != 2 {
						return fmt.Errorf("expected %q to match 'Solution(WorkspaceName)'", val)
					}

					solutionName := segments[0]
					workspaceName := strings.TrimSuffix(segments[1], ")")
					state.SolutionName = solutionName
					state.WorkspaceName = workspaceName
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

func (s LogAnalyticsSolutionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.SolutionsClient

			id, err := solution.ParseSolutionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (s LogAnalyticsSolutionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.SolutionsClient

			id, err := solution.ParseSolutionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config SolutionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			payload := solution.SolutionPatch{}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

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

	expandedPlan := solution.SolutionPlan{
		Name:          utils.String(name),
		PromotionCode: utils.String(promotionCode),
		Publisher:     utils.String(publisher),
		Product:       utils.String(product),
	}

	return expandedPlan
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
