// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_gallery_application_assignment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient

			id, err := parse.VirtualMachineGalleryApplicationAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.VirtualMachineId, virtualmachines.GetOperationOptions{})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %q: %s", id, err)
			}

			virtualMachine := resp.Model
			if virtualMachine == nil {
				return fmt.Errorf("retrieving model of %q: %s", id, err)
			}

			if virtualMachine.Properties == nil || virtualMachine.Properties.ApplicationProfile == nil || virtualMachine.Properties.ApplicationProfile.GalleryApplications == nil {
				return metadata.MarkAsGone(id)
			}

			for _, application := range pointer.From(virtualMachine.Properties.ApplicationProfile.GalleryApplications) {
				if strings.EqualFold(id.GalleryApplicationVersionId.ID(), application.PackageReferenceId) {
					state := VirtualMachineGalleryApplicationAssignmentResourceResourceModel{
						VirtualMachineId:            id.VirtualMachineId.ID(),
						GalleryApplicationVersionId: id.GalleryApplicationVersionId.ID(),
						ConfigurationBlobUri:        pointer.From(application.ConfigurationReference),
						Order:                       pointer.From(application.Order),
						Tag:                         pointer.From(application.Tags),
					}
					return metadata.Encode(&state)
				}
			}

			return metadata.MarkAsGone(id)
		},
		Timeout: 5 * time.Minute,
	}
}
