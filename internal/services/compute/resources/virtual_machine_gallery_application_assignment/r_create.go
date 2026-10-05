// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_gallery_application_assignment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient
			var state VirtualMachineGalleryApplicationAssignmentResourceResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			virtualMachineID, err := virtualmachines.ParseVirtualMachineID(state.VirtualMachineId)
			if err != nil {
				return fmt.Errorf("parsing `virtual_machine_id`, %+v", err)
			}

			locks.ByID(virtualMachineID.ID())
			defer locks.UnlockByID(virtualMachineID.ID())

			resp, err := client.Get(ctx, *virtualMachineID, virtualmachines.GetOperationOptions{Expand: pointer.To(virtualmachines.InstanceViewTypesUserData)})
			if err != nil {
				return fmt.Errorf("checking for presence of existing %q: %+v", *virtualMachineID, err)
			}

			virtualMachine := resp.Model
			if virtualMachine == nil {
				return fmt.Errorf("retrieving model of existing %q: %+v", *virtualMachineID, err)
			}

			if virtualMachine.Properties == nil {
				virtualMachine.Properties = pointer.To(virtualmachines.VirtualMachineProperties{})
			}
			if virtualMachine.Properties.ApplicationProfile == nil {
				virtualMachine.Properties.ApplicationProfile = pointer.To(virtualmachines.ApplicationProfile{})
			}
			if virtualMachine.Properties.ApplicationProfile.GalleryApplications == nil {
				virtualMachine.Properties.ApplicationProfile.GalleryApplications = pointer.To(make([]virtualmachines.VMGalleryApplication, 0))
			}

			galleryApplicationVersionId, err := galleryapplicationversions.ParseApplicationVersionID(state.GalleryApplicationVersionId)
			if err != nil {
				return fmt.Errorf("parsing `gallery_application_version_id`: %+v", err)
			}

			applications := virtualMachine.Properties.ApplicationProfile.GalleryApplications
			exists := false
			for _, application := range pointer.From(applications) {
				if strings.EqualFold(galleryApplicationVersionId.ID(), application.PackageReferenceId) {
					exists = true
					if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return tf.ImportAsExistsError(r.ResourceType(), virtualMachineID.ID())
					}
				}
			}

			if !exists {
				*applications = append(pointer.From(applications), virtualmachines.VMGalleryApplication{
					PackageReferenceId:     galleryApplicationVersionId.ID(),
					ConfigurationReference: pointer.To(state.ConfigurationBlobUri),
					Order:                  pointer.To(state.Order),
					Tags:                   pointer.To(state.Tag),
				})
			}

			virtualMachineUpdate := &virtualmachines.VirtualMachineUpdate{
				Properties: &virtualmachines.VirtualMachineProperties{
					ApplicationProfile: &virtualmachines.ApplicationProfile{
						GalleryApplications: applications,
					},
				},
			}

			// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
			if err = client.UpdateThenPoll(ctx, *virtualMachineID, pointer.From(virtualMachineUpdate), virtualmachines.DefaultUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating Gallery Application Assignment %q: %+v", virtualMachineID, err)
			}

			metadata.SetID(parse.NewVirtualMachineGalleryApplicationAssignmentID(*virtualMachineID, *galleryApplicationVersionId))
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
