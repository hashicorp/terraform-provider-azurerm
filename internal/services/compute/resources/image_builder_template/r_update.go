// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder_template

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/imagebuilder/2024-02-01/virtualmachineimagetemplate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineImageTemplateClient

			var model ImageBuilderTemplateResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id, err := virtualmachineimagetemplate.ParseImageTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			parameters := virtualmachineimagetemplate.ImageTemplateUpdateParameters{}

			if metadata.ResourceData.HasChange("identity") {
				identityExpanded, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				parameters.Identity = identityExpanded
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = &model.Tags
			}

			if err = client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating image builder template %q (Resource Group %q): %+v", id.ImageTemplateName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}
