// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_gallery_application_assignment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient
			var state VirtualMachineGalleryApplicationAssignmentResourceResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := parse.VirtualMachineGalleryApplicationAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.VirtualMachineId.ID())
			defer locks.UnlockByID(id.VirtualMachineId.ID())

			resp, err := client.Get(ctx, id.VirtualMachineId, virtualmachines.GetOperationOptions{Expand: pointer.To(virtualmachines.InstanceViewTypesUserData)})
			if err != nil {
				return fmt.Errorf("checking for presence of existing %q: %+v", id.VirtualMachineId, err)
			}

			virtualMachine := resp.Model
			if virtualMachine == nil {
				return fmt.Errorf("checking model of existing %q: %+v", id.VirtualMachineId, err)
			}

			if virtualMachine.Properties == nil || virtualMachine.Properties.ApplicationProfile == nil || virtualMachine.Properties.ApplicationProfile.GalleryApplications == nil {
				return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
			}

			for i, application := range pointer.From(virtualMachine.Properties.ApplicationProfile.GalleryApplications) {
				if strings.EqualFold(id.GalleryApplicationVersionId.ID(), application.PackageReferenceId) {
					updatedApplication := application
					if metadata.ResourceData.HasChange("order") {
						updatedApplication.Order = pointer.To(state.Order)
					}
					(*virtualMachine.Properties.ApplicationProfile.GalleryApplications)[i] = updatedApplication

					virtualMachineUpdate := &virtualmachines.VirtualMachineUpdate{
						Properties: &virtualmachines.VirtualMachineProperties{
							ApplicationProfile: &virtualmachines.ApplicationProfile{
								GalleryApplications: virtualMachine.Properties.ApplicationProfile.GalleryApplications,
							},
						},
					}

					if err = client.UpdateThenPoll(ctx, id.VirtualMachineId, *virtualMachineUpdate, virtualmachines.DefaultUpdateOperationOptions()); err != nil {
						return fmt.Errorf("updating Gallery Application Assignment %q: %+v", id, err)
					}

					return nil
				}
			}

			return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
		},
		Timeout: 30 * time.Minute,
	}
}
