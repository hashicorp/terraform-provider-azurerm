package api_center_deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/deployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.DeploymentsClient

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ApiCenterDeploymentModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			if config.Title == "" {
				config.Title = config.Name
			}

			definition, err := apidefinitions.ParseDefinitionID(config.DefinitionId)
			if err != nil {
				return err
			}

			id := deployments.NewDeploymentID(subscriptionId, definition.ResourceGroupName, definition.ServiceName, definition.WorkspaceName, definition.ApiName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			scopedDefinitionId := fmt.Sprintf("/workspaces/%s/apis/%s/versions/%s/definitions/%s", definition.WorkspaceName, definition.ApiName, definition.VersionName, definition.DefinitionName)

			environment, err := environments.ParseEnvironmentID(config.EnvironmentId)
			if err != nil {
				return err
			}
			scopedEnvironmentId := fmt.Sprintf("/workspaces/%s/environments/%s", environment.WorkspaceName, environment.EnvironmentName)

			customProperties, err := helpers.ExpandCustomProperties(config.CustomProperties)
			if err != nil {
				return err
			}

			payload := deployments.Deployment{
				Name: new(config.Name),
				Properties: &deployments.DeploymentProperties{
					CustomProperties: customProperties,
					DefinitionId:     new(scopedDefinitionId),
					Description:      new(config.Description),
					EnvironmentId:    new(scopedEnvironmentId),
					Server:           &deployments.DeploymentServer{RuntimeUri: pointer.To(config.ServerRuntimeUri)},
					State:            new(deployments.DeploymentState(config.State)),
					Title:            new(config.Title),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
