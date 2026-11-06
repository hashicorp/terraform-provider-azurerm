// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder_template

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/imagebuilder/2024-02-01/virtualmachineimagetemplate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineImageTemplateClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ImageBuilderTemplateResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := virtualmachineimagetemplate.NewImageTemplateID(subscriptionId, model.ResourceGroupName, model.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for the presence of existing Image Builder Template %q (Resource Group %q): %s", model.Name, model.ResourceGroupName, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			distribution, err := expandBasicImageTemplateDistributor(model.Distributions, subscriptionId)
			if err != nil {
				return err
			}

			customizer, err := expandBasicImageTemplateCustomizer(metadata.ResourceData, model.Customizer)
			if err != nil {
				return err
			}

			identityExpanded, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			parameters := virtualmachineimagetemplate.ImageTemplate{
				Location: model.Location,
				Identity: *identityExpanded,
				Properties: &virtualmachineimagetemplate.ImageTemplateProperties{
					VMProfile: &virtualmachineimagetemplate.ImageTemplateVMProfile{
						VMSize:       &model.Size,
						OsDiskSizeGB: &model.DiskSizeGb,
					},

					Source:                expandBasicImageTemplateSource(model.SourceManagedImageId, model.SourcePlatformImage, model.SourceSharedImageVersionId),
					Distribute:            *distribution,
					Customize:             customizer,
					BuildTimeoutInMinutes: &model.BuildTimeoutMinutes,
				},

				Tags: &model.Tags,
			}

			if model.SubnetId != "" {
				parameters.Properties.VMProfile.VnetConfig = &virtualmachineimagetemplate.VirtualNetworkConfig{
					SubnetId: &model.SubnetId,
				}
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating image builder template %q (Resource Group %q): %+v", model.Name, model.ResourceGroupName, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
