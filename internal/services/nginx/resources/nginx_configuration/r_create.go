// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	nginxClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx/client"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := nginxClient.From(metadata.Client).NginxConfiguration

			var model ConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			deployID, err := nginxdeployment.ParseNginxDeploymentID(model.NginxDeploymentId)
			if err != nil {
				return err
			}

			subscriptionID := account.From(metadata.Client).SubscriptionId
			id := nginxconfiguration.NewConfigurationID(subscriptionID, deployID.ResourceGroupName, deployID.NginxDeploymentName, defaultConfigurationName)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ConfigurationsGet(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			req := model.ToSDKModel()

			if err := client.ConfigurationsCreateOrUpdateCallbackThenPoll(ctx, id, req, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
