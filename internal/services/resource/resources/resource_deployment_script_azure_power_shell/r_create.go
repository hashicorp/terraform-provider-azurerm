// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_deployment_script_azure_power_shell

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-10-01/deploymentscripts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ResourceDeploymentScriptAzurePowerShellModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Resource.DeploymentScriptsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := deploymentscripts.NewDeploymentScriptID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &deploymentscripts.AzurePowerShellScript{
				Location: location.Normalize(model.Location),
				Properties: deploymentscripts.AzurePowerShellScriptProperties{
					AzPowerShellVersion:    model.Version,
					CleanupPreference:      &model.CleanupPreference,
					RetentionInterval:      model.RetentionInterval,
					SupportingScriptUris:   &model.SupportingScriptUris,
					ContainerSettings:      helpers.ExpandContainerConfigurationModel(model.ContainerSettings),
					EnvironmentVariables:   helpers.ExpandEnvironmentVariableModelArray(model.EnvironmentVariables),
					StorageAccountSettings: helpers.ExpandStorageAccountConfigurationModel(model.StorageAccountSettings),
				},
			}

			identityValue, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]any))
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

			if err := client.CreateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
