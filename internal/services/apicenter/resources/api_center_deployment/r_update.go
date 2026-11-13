// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/deployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.DeploymentsClient

			id, err := deployments.ParseDeploymentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ApiCenterDeploymentModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			model := existing.Model
			if metadata.ResourceData.HasChange("custom_properties") {
				model.Properties.CustomProperties, err = helpers.ExpandCustomProperties(state.CustomProperties)
				if err != nil {
					return err
				}
			}

			if metadata.ResourceData.HasChange("api_definition_id") {
				definition, err := apidefinitions.ParseDefinitionID(state.DefinitionId)
				if err != nil {
					return err
				}
				scopedDefinitionId := fmt.Sprintf("/workspaces/%s/apis/%s/versions/%s/definitions/%s", definition.WorkspaceName, definition.ApiName, definition.VersionName, definition.DefinitionName)
				model.Properties.DefinitionId = pointer.To(scopedDefinitionId)
			}

			if metadata.ResourceData.HasChange("description") {
				model.Properties.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("environment_id") {
				environment, err := environments.ParseEnvironmentID(state.EnvironmentId)
				if err != nil {
					return err
				}
				scopedEnvironmentId := fmt.Sprintf("/workspaces/%s/environments/%s", environment.WorkspaceName, environment.EnvironmentName)
				model.Properties.EnvironmentId = pointer.To(scopedEnvironmentId)
			}

			if metadata.ResourceData.HasChange("server_runtime_uri") {
				model.Properties.Server.RuntimeUri = pointer.To(state.ServerRuntimeUri)
			}
			if metadata.ResourceData.HasChange("state") {
				model.Properties.State = pointer.To(deployments.DeploymentState(state.State))
			}
			if metadata.ResourceData.HasChange("display_name") {
				model.Properties.Title = pointer.To(state.Title)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
