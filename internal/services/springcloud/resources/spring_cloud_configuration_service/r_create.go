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
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudConfigurationServiceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := springcloudClient.From(metadata.Client).AppPlatformClient
			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return err
			}
			id := appplatform.NewConfigurationServiceID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ConfigurationServicesGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(s.ResourceType(), id)
				}
			}

			configurationServiceResource := appplatform.ConfigurationServiceResource{
				Properties: &appplatform.ConfigurationServiceProperties{
					Generation: pointer.ToEnum[appplatform.ConfigurationServiceGeneration](model.Generation),
					Settings: &appplatform.ConfigurationServiceSettings{
						GitProperty: &appplatform.ConfigurationServiceGitProperty{
							Repositories: expandConfigurationServiceConfigurationServiceGitRepositoryArray(model.Repository),
						},
						RefreshIntervalInSeconds: pointer.To(model.RefreshInterval),
					},
				},
			}
			if err := client.ConfigurationServicesCreateOrUpdateCallbackThenPoll(ctx, id, configurationServiceResource, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func expandConfigurationServiceConfigurationServiceGitRepositoryArray(input []SpringCloudRepositoryModel) *[]appplatform.ConfigurationServiceGitRepository {
	if len(input) == 0 {
		return nil
	}
	results := make([]appplatform.ConfigurationServiceGitRepository, 0)
	for _, v := range input {
		repo := appplatform.ConfigurationServiceGitRepository{
			Name:                  v.Name,
			Patterns:              v.Patterns,
			Uri:                   v.Uri,
			Label:                 v.Label,
			SearchPaths:           pointer.To(v.SearchPaths),
			Username:              pointer.To(v.Username),
			Password:              pointer.To(v.Password),
			HostKey:               pointer.To(v.HostKey),
			HostKeyAlgorithm:      pointer.To(v.HostKeyAlgorithm),
			PrivateKey:            pointer.To(v.PrivateKey),
			StrictHostKeyChecking: pointer.To(v.StrictHostKeyChecking),
		}
		if v.CaCertificateId != "" {
			repo.CaCertResourceId = pointer.To(v.CaCertificateId)
		}
		results = append(results, repo)
	}
	return &results
}
