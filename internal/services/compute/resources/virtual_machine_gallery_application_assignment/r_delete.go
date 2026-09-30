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
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient

			id, err := parse.VirtualMachineGalleryApplicationAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.VirtualMachineId.ID())
			defer locks.UnlockByID(id.VirtualMachineId.ID())

			resp, err := client.Get(ctx, id.VirtualMachineId, virtualmachines.GetOperationOptions{})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("checking for presence of existing  %q: %+v", id, err)
			}

			virtualMachine := resp.Model
			if virtualMachine == nil {
				return fmt.Errorf("retrieving model of %q: %s", id, err)
			}

			if virtualMachine.Properties != nil && virtualMachine.Properties.ApplicationProfile != nil && virtualMachine.Properties.ApplicationProfile.GalleryApplications != nil {
				galleryApplications := virtualMachine.Properties.ApplicationProfile.GalleryApplications
				updatedApplications := make([]virtualmachines.VMGalleryApplication, 0)
				for _, application := range pointer.From(galleryApplications) {
					if !strings.EqualFold(id.GalleryApplicationVersionId.ID(), application.PackageReferenceId) {
						updatedApplications = append(updatedApplications, application)
					}
				}

				virtualMachineUpdate := &virtualmachines.VirtualMachineUpdate{
					Properties: &virtualmachines.VirtualMachineProperties{
						ApplicationProfile: &virtualmachines.ApplicationProfile{
							GalleryApplications: pointer.To(updatedApplications),
						},
					},
				}

				if err = client.UpdateThenPoll(ctx, id.VirtualMachineId, *virtualMachineUpdate, virtualmachines.DefaultUpdateOperationOptions()); err != nil {
					return fmt.Errorf("deleting Gallery Application Assignment %q: %+v", id, err)
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
