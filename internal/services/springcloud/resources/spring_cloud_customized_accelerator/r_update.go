// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_customized_accelerator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
)

func (s Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := springcloudClient.From(metadata.Client).AppPlatformClient

			id, err := appplatform.ParseCustomizedAcceleratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SpringCloudCustomizedAcceleratorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.CustomizedAcceleratorsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model.Properties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("git_repository") {
				properties.GitRepository = expandSpringCloudCustomizedAcceleratorGitRepository(model.GitRepository)
			}

			if metadata.ResourceData.HasChange("accelerator_tags") {
				properties.AcceleratorTags = &model.AcceleratorTags
			}

			if metadata.ResourceData.HasChange("accelerator_type") {
				properties.AcceleratorType = pointer.ToEnum[appplatform.CustomizedAcceleratorType](model.AcceleratorType)
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Description = &model.Description
			}

			if metadata.ResourceData.HasChange("display_name") {
				properties.DisplayName = &model.DisplayName
			}

			if metadata.ResourceData.HasChange("icon_url") {
				properties.IconURL = &model.IconURL
			}

			CustomizedAcceleratorResource := appplatform.CustomizedAcceleratorResource{
				Properties: properties,
			}
			if err = client.CustomizedAcceleratorsCreateOrUpdateThenPoll(ctx, *id, CustomizedAcceleratorResource); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
