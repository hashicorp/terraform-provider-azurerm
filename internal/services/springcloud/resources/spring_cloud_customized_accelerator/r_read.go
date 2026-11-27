// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_customized_accelerator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := springcloudClient.From(metadata.Client).AppPlatformClient

			id, err := appplatform.ParseCustomizedAcceleratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.CustomizedAcceleratorsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			state := SpringCloudCustomizedAcceleratorModel{
				Name:                     id.CustomizedAcceleratorName,
				SpringCloudAcceleratorId: appplatform.NewApplicationAcceleratorID(id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ApplicationAcceleratorName).ID(),
			}

			if props := resp.Model.Properties; props != nil {
				if props.AcceleratorTags != nil {
					state.AcceleratorTags = *props.AcceleratorTags
				}
				if props.AcceleratorType != nil {
					state.AcceleratorType = string(*props.AcceleratorType)
				}
				if props.Description != nil {
					state.Description = *props.Description
				}
				if props.DisplayName != nil {
					state.DisplayName = *props.DisplayName
				}

				var model SpringCloudCustomizedAcceleratorModel
				if err := metadata.Decode(&model); err != nil {
					return fmt.Errorf("decoding: %+v", err)
				}
				state.GitRepository = flattenSpringCloudCustomizedAcceleratorGitRepository(model.GitRepository, props.GitRepository)

				if props.IconURL != nil {
					state.IconURL = *props.IconURL
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenSpringCloudCustomizedAcceleratorGitRepository(state []GitRepositoryModel, input appplatform.AcceleratorGitRepository) []GitRepositoryModel {
	basicAuth := make([]BasicAuthModel, 0)

	caCertificateId := ""
	if publicAuthSetting, ok := input.AuthSetting.(appplatform.AcceleratorPublicSetting); ok && publicAuthSetting.CaCertResourceId != nil {
		certificatedId, err := appplatform.ParseCertificateIDInsensitively(*publicAuthSetting.CaCertResourceId)
		if err == nil {
			caCertificateId = certificatedId.ID()
		}
	}
	if basicAuthSetting, ok := input.AuthSetting.(appplatform.AcceleratorBasicAuthSetting); ok {
		if basicAuthSetting.CaCertResourceId != nil {
			certificatedId, err := appplatform.ParseCertificateIDInsensitively(*basicAuthSetting.CaCertResourceId)
			if err == nil {
				caCertificateId = certificatedId.ID()
			}
		}
		var basicAuthState BasicAuthModel
		if len(state) != 0 && len(state[0].BasicAuth) != 0 {
			basicAuthState = state[0].BasicAuth[0]
		}
		basicAuth = append(basicAuth, BasicAuthModel{
			Username: basicAuthSetting.Username,
			Password: basicAuthState.Password,
		})
	}

	sshAuth := make([]SshAuthModel, 0)
	if _, ok := input.AuthSetting.(appplatform.AcceleratorSshSetting); ok {
		var sshAuthState SshAuthModel
		if len(state) != 0 && len(state[0].SshAuth) != 0 {
			sshAuthState = state[0].SshAuth[0]
		}
		sshAuth = append(sshAuth, sshAuthState)
	}

	url := input.Url

	return []GitRepositoryModel{
		{
			BasicAuth:         basicAuth,
			SshAuth:           sshAuth,
			Branch:            pointer.From(input.Branch),
			CaCertificateId:   caCertificateId,
			Commit:            pointer.From(input.Commit),
			GitTag:            pointer.From(input.GitTag),
			IntervalInSeconds: pointer.From(input.IntervalInSeconds),
			Url:               url,
			Path:              pointer.From(input.SubPath),
		},
	}
}
