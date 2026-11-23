// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_deployment_script_azure_power_shell

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
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

			model, ok := resp.Model.(deploymentscripts.AzurePowerShellScript)
			if !ok {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := ResourceDeploymentScriptAzurePowerShellModel{
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

			state.Version = properties.AzPowerShellVersion

			if properties.CleanupPreference != nil {
				state.CleanupPreference = *properties.CleanupPreference
			}

			state.ContainerSettings = helpers.FlattenContainerConfigurationModel(properties.ContainerSettings)

			var originalModel ResourceDeploymentScriptAzurePowerShellModel
			if err := metadata.Decode(&originalModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state.EnvironmentVariables = helpers.FlattenEnvironmentVariableModelArray(properties.EnvironmentVariables, originalModel.EnvironmentVariables)

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

			state.StorageAccountSettings = helpers.FlattenStorageAccountConfigurationModel(properties.StorageAccountSettings, originalModel.StorageAccountSettings)

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
