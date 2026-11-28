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
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudCustomizedAcceleratorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			springAcceleratorId, err := appplatform.ParseApplicationAcceleratorID(model.SpringCloudAcceleratorId)
			if err != nil {
				return err
			}
			id := appplatform.NewCustomizedAcceleratorID(springAcceleratorId.SubscriptionId, springAcceleratorId.ResourceGroupName, springAcceleratorId.SpringName, springAcceleratorId.ApplicationAcceleratorName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.CustomizedAcceleratorsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(s.ResourceType(), id)
				}
			}

			CustomizedAcceleratorResource := appplatform.CustomizedAcceleratorResource{
				Properties: &appplatform.CustomizedAcceleratorProperties{
					AcceleratorType: pointer.ToEnum[appplatform.CustomizedAcceleratorType](model.AcceleratorType),
					DisplayName:     pointer.To(model.DisplayName),
					Description:     pointer.To(model.Description),
					IconURL:         pointer.To(model.IconURL),
					AcceleratorTags: pointer.To(model.AcceleratorTags),
					GitRepository:   expandSpringCloudCustomizedAcceleratorGitRepository(model.GitRepository),
				},
			}

			if err := client.CustomizedAcceleratorsCreateOrUpdateCallbackThenPoll(ctx, id, CustomizedAcceleratorResource, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func expandSpringCloudCustomizedAcceleratorGitRepository(repository []GitRepositoryModel) appplatform.AcceleratorGitRepository {
	if len(repository) == 0 {
		return appplatform.AcceleratorGitRepository{}
	}
	repo := repository[0]
	var authSetting appplatform.AcceleratorAuthSetting
	var caCertResourceID *string
	if repo.CaCertificateId != "" {
		caCertResourceID = pointer.To(repo.CaCertificateId)
	}
	authSetting = appplatform.AcceleratorPublicSetting{
		CaCertResourceId: caCertResourceID,
	}
	if len(repo.BasicAuth) != 0 {
		basicAuth := repo.BasicAuth[0]
		authSetting = appplatform.AcceleratorBasicAuthSetting{
			Username:         basicAuth.Username,
			Password:         pointer.To(basicAuth.Password),
			CaCertResourceId: caCertResourceID,
		}
	}
	if len(repo.SshAuth) != 0 {
		sshAuth := repo.SshAuth[0]
		authSetting = appplatform.AcceleratorSshSetting{
			HostKey:          pointer.To(sshAuth.HostKey),
			HostKeyAlgorithm: pointer.To(sshAuth.PrivateKeyAlgorithm),
			PrivateKey:       pointer.To(sshAuth.PrivateKey),
		}
	}
	res := appplatform.AcceleratorGitRepository{
		Url:         repo.Url,
		Branch:      pointer.To(repo.Branch),
		Commit:      pointer.To(repo.Commit),
		GitTag:      pointer.To(repo.GitTag),
		AuthSetting: authSetting,
		SubPath:     pointer.To(repo.Path),
	}
	if repo.IntervalInSeconds != 0 {
		res.IntervalInSeconds = pointer.To(repo.IntervalInSeconds)
	}
	return res
}
