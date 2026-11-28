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
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseConfigurationServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ConfigurationServicesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			springId := commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName)

			var model SpringCloudConfigurationServiceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudConfigurationServiceModel{
				Name:                 id.ConfigurationServiceName,
				SpringCloudServiceId: springId.ID(),
			}

			if resp.Model != nil {
				if props := resp.Model.Properties; props != nil {
					state.Generation = pointer.FromEnum(props.Generation)
					if props.Settings != nil && props.Settings.GitProperty != nil {
						state.Repository = flattenConfigurationServiceConfigurationServiceGitRepositoryArray(props.Settings.GitProperty.Repositories, model.Repository)
						state.RefreshInterval = pointer.From(props.Settings.RefreshIntervalInSeconds)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenConfigurationServiceConfigurationServiceGitRepositoryArray(input *[]appplatform.ConfigurationServiceGitRepository, old []SpringCloudRepositoryModel) []SpringCloudRepositoryModel {
	results := make([]SpringCloudRepositoryModel, 0)
	if input == nil {
		return results
	}

	oldItems := make(map[string]SpringCloudRepositoryModel)
	for _, v := range old {
		oldItems[v.Name] = v
	}

	for _, item := range *input {
		strictHostKeyChecking := pointer.From(item.StrictHostKeyChecking)

		var hostKey string
		var hostKeyAlgorithm string
		var privateKey string
		var username string
		var password string
		if oldItem, ok := oldItems[item.Name]; ok {
			hostKey = oldItem.HostKey
			hostKeyAlgorithm = oldItem.HostKeyAlgorithm
			privateKey = oldItem.PrivateKey
			username = oldItem.Username
			password = oldItem.Password
		}

		var caCertificateId string
		if item.CaCertResourceId != nil {
			certificatedId, err := appplatform.ParseCertificateIDInsensitively(*item.CaCertResourceId)
			if err == nil {
				caCertificateId = certificatedId.ID()
			}
		}
		results = append(results, SpringCloudRepositoryModel{
			Name:                  item.Name,
			Label:                 item.Label,
			Patterns:              item.Patterns,
			Uri:                   item.Uri,
			CaCertificateId:       caCertificateId,
			HostKey:               hostKey,
			HostKeyAlgorithm:      hostKeyAlgorithm,
			Password:              password,
			PrivateKey:            privateKey,
			SearchPaths:           pointer.From(item.SearchPaths),
			StrictHostKeyChecking: strictHostKeyChecking,
			Username:              username,
		})
	}
	return results
}
