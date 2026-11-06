package api_center_deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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

			api, err := apis.ParseApiID(config.ApiId)
			if err != nil {
				return err
			}

			id := deployments.NewDeploymentID(subscriptionId, api.ResourceGroupName, api.ServiceName, api.WorkspaceName, api.ApiName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := deployments.Deployment{
				Name: new(config.Name),
				Properties: &deployments.DeploymentProperties{
					DefinitionId:  new(config.DefinitionId),
					Description:   new(config.Description),
					EnvironmentId: new(config.EnvironmentId),
					Server:        expandDeploymentServer(config.Server),
					State:         new(deployments.DeploymentState(config.State)),
					Title:         new(config.Title),
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
