// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type cognitiveDeploymentModel struct {
	Name                     string                 `tfschema:"name"`
	CognitiveAccountId       string                 `tfschema:"cognitive_account_id"`
	DynamicThrottlingEnabled bool                   `tfschema:"dynamic_throttling_enabled"`
	Model                    []DeploymentModelModel `tfschema:"model"`
	RaiPolicyName            string                 `tfschema:"rai_policy_name"`
	Sku                      []DeploymentSkuModel   `tfschema:"sku"`
	VersionUpgradeOption     string                 `tfschema:"version_upgrade_option"`
}

type DeploymentModelModel struct {
	Format  string `tfschema:"format"`
	Name    string `tfschema:"name"`
	Version string `tfschema:"version"`
}

type DeploymentSkuModel struct {
	Name     string `tfschema:"name"`
	Tier     string `tfschema:"tier"`
	Size     string `tfschema:"size"`
	Family   string `tfschema:"family"`
	Capacity int64  `tfschema:"capacity"`
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

		"dynamic_throttling_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
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
							"Cohere",
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
						Optional: true,
					},
				},
			},
		},

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Standard",
							"DataZoneBatch",
							"DataZoneProvisionedManaged",
							"DataZoneStandard",
							"GlobalBatch",
							"GlobalProvisionedManaged",
							"GlobalStandard",
							"ProvisionedManaged",
						}, false),
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
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"rai_policy_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_upgrade_option": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(deployments.DeploymentModelVersionUpgradeOptionOnceNewDefaultVersionAvailable),
			ValidateFunc: validation.StringInSlice([]string{
				string(deployments.DeploymentModelVersionUpgradeOptionOnceCurrentVersionExpired),
				string(deployments.DeploymentModelVersionUpgradeOptionOnceNewDefaultVersionAvailable),
				string(deployments.DeploymentModelVersionUpgradeOptionNoAutoUpgrade),
			}, false),
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

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

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

			if model.DynamicThrottlingEnabled {
				properties.Properties.DynamicThrottlingEnabled = &model.DynamicThrottlingEnabled
			}

			if model.VersionUpgradeOption != "" {
				option := deployments.DeploymentModelVersionUpgradeOption(model.VersionUpgradeOption)
				properties.Properties.VersionUpgradeOption = &option
			}

			properties.Sku = expandDeploymentSkuModel(model.Sku)

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CognitiveDeploymentResource) Update() sdk.ResourceFunc {
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

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id, err := deployments.ParseDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return err
			}

			properties := resp.Model

			if metadata.ResourceData.HasChange("dynamic_throttling_enabled") {
				properties.Properties.DynamicThrottlingEnabled = pointer.To(model.DynamicThrottlingEnabled)
			}

			if metadata.ResourceData.HasChange("sku.0.capacity") {
				properties.Sku.Capacity = pointer.To(model.Sku[0].Capacity)
			}

			if metadata.ResourceData.HasChange("rai_policy_name") {
				properties.Properties.RaiPolicyName = pointer.To(model.RaiPolicyName)
			}

			if metadata.ResourceData.HasChange("model.0.version") {
				properties.Properties.Model.Version = pointer.To(model.Model[0].Version)
			}

			properties.Properties.VersionUpgradeOption = pointer.To(deployments.DeploymentModelVersionUpgradeOption(model.VersionUpgradeOption))

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
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

				state.DynamicThrottlingEnabled = pointer.From(properties.DynamicThrottlingEnabled)
				state.RaiPolicyName = pointer.From(properties.RaiPolicyName)
				state.VersionUpgradeOption = string(pointer.From(properties.VersionUpgradeOption))
			}
			if sku := flattenDeploymentSkuModel(model.Sku); sku != nil {
				state.Sku = sku
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

func expandDeploymentSkuModel(inputList []DeploymentSkuModel) *deployments.Sku {
	if len(inputList) == 0 {
		return nil
	}
	input := inputList[0]
	s := &deployments.Sku{
		Name: input.Name,
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

func flattenDeploymentSkuModel(input *deployments.Sku) []DeploymentSkuModel {
	if input == nil {
		return nil
	}
	output := DeploymentSkuModel{
		Name: input.Name,
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
	return []DeploymentSkuModel{output}
}
