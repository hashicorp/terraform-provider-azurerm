// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-10-01/deploymentscripts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ResourceDeploymentScriptAzureCliModel ResourceDeploymentScriptModel

type ResourceDeploymentScriptAzureCliResource struct{}

var _ sdk.ResourceWithUpdate = ResourceDeploymentScriptAzureCliResource{}

func (r ResourceDeploymentScriptAzureCliResource) ResourceType() string {
	return "azurerm_resource_deployment_script_azure_cli"
}

func (r ResourceDeploymentScriptAzureCliResource) ModelObject() interface{} {
	return &ResourceDeploymentScriptAzureCliModel{}
}

func (r ResourceDeploymentScriptAzureCliResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deploymentscripts.ValidateDeploymentScriptID
}

func (r ResourceDeploymentScriptAzureCliResource) Arguments() map[string]*pluginsdk.Schema {
	return getDeploymentScriptArguments(AzureCliKind)
}

func (r ResourceDeploymentScriptAzureCliResource) Attributes() map[string]*pluginsdk.Schema {
	return getDeploymentScriptAttributes()
}

func (r ResourceDeploymentScriptAzureCliResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ResourceDeploymentScriptAzureCliModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Resource.DeploymentScriptsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := deploymentscripts.NewDeploymentScriptID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &deploymentscripts.AzureCliScript{
				Location: location.Normalize(model.Location),
				Properties: deploymentscripts.AzureCliScriptProperties{
					AzCliVersion:           model.Version,
					CleanupPreference:      &model.CleanupPreference,
					RetentionInterval:      model.RetentionInterval,
					SupportingScriptUris:   &model.SupportingScriptUris,
					ContainerSettings:      expandContainerConfigurationModel(model.ContainerSettings),
					EnvironmentVariables:   expandEnvironmentVariableModelArray(model.EnvironmentVariables),
					StorageAccountSettings: expandStorageAccountConfigurationModel(model.StorageAccountSettings),
				},
			}

			identityValue, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return err
			}

			if identityValue != nil && identityValue.Type != identity.TypeNone {
				properties.Identity = identityValue
			}

			if model.Arguments != "" {
				properties.Properties.Arguments = &model.Arguments
			}

			if model.ForceUpdateTag != "" {
				properties.Properties.ForceUpdateTag = &model.ForceUpdateTag
			}

			if model.PrimaryScriptUri != "" {
				properties.Properties.PrimaryScriptUri = &model.PrimaryScriptUri
			}

			if model.ScriptContent != "" {
				properties.Properties.ScriptContent = &model.ScriptContent
			}

			if model.Timeout != "" {
				properties.Properties.Timeout = &model.Timeout
			}

			if model.Tags != nil {
				properties.Tags = &model.Tags
			}

			if err := client.CreateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ResourceDeploymentScriptAzureCliResource) Update() sdk.ResourceFunc {
	return updateDeploymentScript()
}

func (r ResourceDeploymentScriptAzureCliResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentScriptsClient

			id, err := deploymentscripts.ParseDeploymentScriptID(metadata.ResourceData.Id())
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

			model, ok := (*resp.Model).(deploymentscripts.AzureCliScript)
			if !ok {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := ResourceDeploymentScriptAzureCliModel{
				Name:              id.DeploymentScriptName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identityValue, err := identity.FlattenUserAssignedMap(model.Identity)
			if err != nil {
				return err
			}

			if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			properties := &model.Properties
			if properties.Arguments != nil {
				state.Arguments = *properties.Arguments
			}

			state.Version = properties.AzCliVersion

			if properties.CleanupPreference != nil {
				state.CleanupPreference = *properties.CleanupPreference
			}

			state.ContainerSettings = flattenContainerConfigurationModel(properties.ContainerSettings)

			var originalModel ResourceDeploymentScriptAzureCliModel
			if err := metadata.Decode(&originalModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state.EnvironmentVariables = flattenEnvironmentVariableModelArray(properties.EnvironmentVariables, originalModel.EnvironmentVariables)

			if properties.ForceUpdateTag != nil {
				state.ForceUpdateTag = *properties.ForceUpdateTag
			}

			if properties.Outputs != nil && *properties.Outputs != nil {
				outputsValue, err := json.Marshal(*properties.Outputs)
				if err != nil {
					return err
				}

				state.Outputs = string(outputsValue)
			}

			if properties.PrimaryScriptUri != nil {
				state.PrimaryScriptUri = *properties.PrimaryScriptUri
			}

			state.RetentionInterval = properties.RetentionInterval

			if properties.ScriptContent != nil {
				state.ScriptContent = *properties.ScriptContent
			}

			state.StorageAccountSettings = flattenStorageAccountConfigurationModel(properties.StorageAccountSettings, originalModel.StorageAccountSettings)

			if properties.SupportingScriptUris != nil {
				state.SupportingScriptUris = *properties.SupportingScriptUris
			}

			if properties.Timeout != nil {
				state.Timeout = *properties.Timeout
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ResourceDeploymentScriptAzureCliResource) Delete() sdk.ResourceFunc {
	return deleteDeploymentScript()
}
