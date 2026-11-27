// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_configuration_service

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
)

func (s Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudConfigurationServiceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return err
			}
			id := appplatform.NewConfigurationServiceID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			client := springcloudClient.From(metadata.Client).AppPlatformClient
			existing, err := client.ConfigurationServicesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}
			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			properties := existing.Model.Properties
			if metadata.ResourceData.HasChange("generation") {
				properties.Generation = pointer.ToEnum[appplatform.ConfigurationServiceGeneration](model.Generation)
			}

			if metadata.ResourceData.HasChange("repository") {
				properties.Settings.GitProperty.Repositories = expandConfigurationServiceConfigurationServiceGitRepositoryArray(model.Repository)
			}

			if metadata.ResourceData.HasChange("refresh_interval_in_seconds") {
				properties.Settings.RefreshIntervalInSeconds = pointer.To(model.RefreshInterval)
			}

			configurationServiceResource := appplatform.ConfigurationServiceResource{
				Properties: properties,
			}
			if err = client.ConfigurationServicesCreateOrUpdateThenPoll(ctx, id, configurationServiceResource); err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
