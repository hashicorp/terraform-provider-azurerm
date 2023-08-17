// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2023-05-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2023-05-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type cognitiveDeploymentModel struct {
	Name               string                         `tfschema:"name"`
	CognitiveAccountId string                         `tfschema:"cognitive_account_id"`
	Model              []DeploymentModelModel         `tfschema:"model"`
	RaiPolicyName      string                         `tfschema:"rai_policy_name"`
	ScaleSettings      []DeploymentScaleSettingsModel `tfschema:"scale"`
}

type DeploymentModelModel struct {
	Format  string `tfschema:"format"`
	Name    string `tfschema:"name"`
	Version string `tfschema:"version"`
}

type DeploymentScaleSettingsModel struct {
	ScaleType string `tfschema:"type"`
	Tier      string `tfschema:"tier"`
	Size      string `tfschema:"size"`
	Family    string `tfschema:"family"`
	Capacity  int64  `tfschema:"capacity"`
}

type CognitiveDeploymentResource struct{}

var _ sdk.Resource = CognitiveDeploymentResource{}

func (r CognitiveDeploymentResource) ResourceType() string {
	return "azurerm_cognitive_deployment"
}

func (r CognitiveDeploymentResource) ModelObject() interface{} {
	return &cognitiveDeploymentModel{}
}

func (r CognitiveDeploymentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deployments.ValidateDeploymentID
}

func (r CognitiveDeploymentResource) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cognitive_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
		},

		"model": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"format": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"OpenAI",
						}, false),
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"version": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"rai_policy_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
	if !features.FourPointOh() {
		arguments["scale"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
					"tier": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(deployments.SkuTierFree),
							string(deployments.SkuTierBasic),
							string(deployments.SkuTierStandard),
							string(deployments.SkuTierPremium),
							string(deployments.SkuTierEnterprise),
						}, false),
					},
					"size": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"family": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		}
	} else {
		//TODO: 4.0 - add corresponding field in cognitiveDeploymentModel struct
		arguments["sku"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
					"tier": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(deployments.SkuTierFree),
							string(deployments.SkuTierBasic),
							string(deployments.SkuTierStandard),
							string(deployments.SkuTierPremium),
							string(deployments.SkuTierEnterprise),
						}, false),
					},
					"size": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"family": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		}
	}
	return arguments
}

func (r CognitiveDeploymentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveDeploymentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model cognitiveDeploymentModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cognitive.DeploymentsClient
			accountId, err := cognitiveservicesaccounts.ParseAccountID(model.CognitiveAccountId)

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			if err != nil {
				return err
			}

			id := deployments.NewDeploymentID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &deployments.Deployment{
				Properties: &deployments.DeploymentProperties{},
			}

			properties.Properties.Model = expandDeploymentModelModel(model.Model)

			if model.RaiPolicyName != "" {
				properties.Properties.RaiPolicyName = &model.RaiPolicyName
			}

			properties.Sku = expandDeploymentSkuModel(model.ScaleSettings)

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CognitiveDeploymentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.DeploymentsClient

			id, err := deployments.ParseDeploymentID(metadata.ResourceData.Id())
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

			state := cognitiveDeploymentModel{
				Name:               id.DeploymentName,
				CognitiveAccountId: cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
			}

			if properties := model.Properties; properties != nil {

				state.Model = flattenDeploymentModelModel(properties.Model)

				if v := properties.RaiPolicyName; v != nil {
					state.RaiPolicyName = *v
				}
				state.ScaleSettings = flattenDeploymentScaleSettingsModel(properties.ScaleSettings)
			}
			if scale := flattenDeploymentSkuModel(model.Sku); scale != nil {
				state.ScaleSettings = scale
			}
			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveDeploymentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.DeploymentsClient

			id, err := deployments.ParseDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			accountId := cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandDeploymentModelModel(inputList []DeploymentModelModel) *deployments.DeploymentModel {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := deployments.DeploymentModel{}

	if input.Format != "" {
		output.Format = &input.Format
	}

	if input.Name != "" {
		output.Name = &input.Name
	}

	if input.Version != "" {
		output.Version = &input.Version
	}

	return &output
}

func expandDeploymentSkuModel(inputList []DeploymentScaleSettingsModel) *deployments.Sku {
	if len(inputList) == 0 {
		return nil
	}
	input := inputList[0]
	s := &deployments.Sku{
		Name: input.ScaleType,
	}
	if input.Capacity != 0 {
		s.Capacity = utils.Int64(input.Capacity)
	}
	if input.Family != "" {
		s.Family = utils.String(input.Family)
	}
	if input.Size != "" {
		s.Size = utils.String(input.Size)
	}
	if input.Tier != "" {
		tier := deployments.SkuTier(input.Tier)
		s.Tier = &tier
	}
	return s
}

func flattenDeploymentModelModel(input *deployments.DeploymentModel) []DeploymentModelModel {
	var outputList []DeploymentModelModel
	if input == nil {
		return outputList
	}

	output := DeploymentModelModel{}
	format := ""
	if input.Format != nil {
		format = *input.Format
	}
	output.Format = format

	name := ""
	if input.Name != nil {
		name = *input.Name
	}
	output.Name = name

	version := ""
	if input.Version != nil {
		version = *input.Version
	}
	output.Version = version

	return append(outputList, output)
}

func flattenDeploymentScaleSettingsModel(input *deployments.DeploymentScaleSettings) []DeploymentScaleSettingsModel {
	if input == nil || input.ScaleType == nil {
		return nil
	}

	output := DeploymentScaleSettingsModel{
		ScaleType: string(*input.ScaleType),
	}
	return []DeploymentScaleSettingsModel{output}
}

func flattenDeploymentSkuModel(input *deployments.Sku) []DeploymentScaleSettingsModel {
	if input == nil {
		return nil
	}
	output := DeploymentScaleSettingsModel{
		ScaleType: input.Name,
	}
	if input.Capacity != nil {
		output.Capacity = *input.Capacity
	}
	if input.Tier != nil {
		output.Tier = string(*input.Tier)
	}
	if input.Size != nil {
		output.Size = *input.Size
	}
	if input.Family != nil {
		output.Family = *input.Family
	}
	return []DeploymentScaleSettingsModel{output}
}
