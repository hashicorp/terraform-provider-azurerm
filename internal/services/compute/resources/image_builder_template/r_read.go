// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder_template

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/imagebuilder/2024-02-01/virtualmachineimagetemplate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineImageTemplateClient

			state := ImageBuilderTemplateResourceModel{}

			id, err := virtualmachineimagetemplate.ParseImageTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] %s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving image builder template %q (Resource Group %q): %+v", id.ImageTemplateName, id.ResourceGroupName, err)
			}

			state.Name = id.ImageTemplateName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				identityFlattened, err := identity.FlattenUserAssignedMapToModel(&model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = *identityFlattened

				if imageTemplateProperties := resp.Model.Properties; imageTemplateProperties != nil {
					managedImageId, platformImage, sharedImageVersionId := flattenBasicImageTemplateSource(imageTemplateProperties.Source)

					// only one among managedImageId / platformImage / sharedImageVersionId would be returned.
					if managedImageId != "" {
						state.SourceManagedImageId = managedImageId
					}
					if len(platformImage) > 0 {
						state.SourcePlatformImage = platformImage
					}

					if sharedImageVersionId != "" {
						state.SourceSharedImageVersionId = sharedImageVersionId
					}

					flattenedDistributions, err := flattenBasicImageTemplateDistributor(&imageTemplateProperties.Distribute)
					if err != nil {
						return err
					}

					state.Distributions = flattenedDistributions
					state.Customizer = flattenBasicImageTemplateCustomizer(imageTemplateProperties.Customize)
					state.BuildTimeoutMinutes = *imageTemplateProperties.BuildTimeoutInMinutes

					if vmProfile := resp.Model.Properties.VMProfile; vmProfile != nil {
						state.Size = *vmProfile.VMSize
						state.DiskSizeGb = *vmProfile.OsDiskSizeGB

						if vnetConfig := vmProfile.VnetConfig; vnetConfig != nil {
							state.SubnetId = *vnetConfig.SubnetId
						}
					}
				}
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}
