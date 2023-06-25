package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2022-10-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2022-10-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
	ScaleType deployments.DeploymentScaleType `tfschema:"type"`
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
	return map[string]*pluginsdk.Schema{
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

		"scale": {
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
						ValidateFunc: validation.StringInSlice([]string{
							string(deployments.DeploymentScaleTypeStandard),
						}, false),
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

			properties.Properties.ScaleSettings = expandDeploymentScaleSettingsModel(model.ScaleSettings)

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

				state.ScaleSettings = flattenDeploymentScaleSettingsModel(properties.ScaleSettings)

				if v := properties.RaiPolicyName; v != nil {
					state.RaiPolicyName = *v
				}
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

func expandDeploymentScaleSettingsModel(inputList []DeploymentScaleSettingsModel) *deployments.DeploymentScaleSettings {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := deployments.DeploymentScaleSettings{
		ScaleType: &input.ScaleType,
	}

	return &output
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
	var outputList []DeploymentScaleSettingsModel
	if input == nil {
		return outputList
	}

	output := DeploymentScaleSettingsModel{}

	if input.ScaleType != nil {
		output.ScaleType = *input.ScaleType
	}

	return append(outputList, output)
}
