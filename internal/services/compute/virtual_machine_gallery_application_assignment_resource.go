// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type VirtualMachineGalleryApplicationAssignmentResource struct{}

var (
	_ sdk.ResourceWithUpdate = VirtualMachineGalleryApplicationAssignmentResource{}
)

type VirtualMachineGalleryApplicationAssignmentResourceResourceModel struct {
	GalleryApplicationVersionId string `tfschema:"gallery_application_version_id"`
	VirtualMachineId            string `tfschema:"virtual_machine_id"`
	ConfigurationBlobUri        string `tfschema:"configuration_blob_uri"`
	Order                       int64  `tfschema:"order"`
	Tag                         string `tfschema:"tag"`
}

func (r VirtualMachineGalleryApplicationAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"gallery_application_version_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: galleryapplicationversions.ValidateApplicationVersionID,
		},

		"virtual_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualmachines.ValidateVirtualMachineID,
		},

		"configuration_blob_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"order": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntBetween(0, 2147483647),
		},

		"tag": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r VirtualMachineGalleryApplicationAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r VirtualMachineGalleryApplicationAssignmentResource) ResourceType() string {
	return "azurerm_virtual_machine_gallery_application_assignment"
}

func (r VirtualMachineGalleryApplicationAssignmentResource) ModelObject() interface{} {
	return &VirtualMachineGalleryApplicationAssignmentResourceResourceModel{}
}

func (r VirtualMachineGalleryApplicationAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.VirtualMachineGalleryApplicationAssignmentIDValidation
}

func (r VirtualMachineGalleryApplicationAssignmentResource) Create() sdk.ResourceFunc {
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
			for _, application := range pointer.From(applications) {
				if strings.EqualFold(galleryApplicationVersionId.ID(), application.PackageReferenceId) {
					return tf.ImportAsExistsError(r.ResourceType(), virtualMachineID.ID())
				}
			}

			*applications = append(pointer.From(applications), virtualmachines.VMGalleryApplication{
				PackageReferenceId:     galleryApplicationVersionId.ID(),
				ConfigurationReference: pointer.To(state.ConfigurationBlobUri),
				Order:                  pointer.To(state.Order),
				Tags:                   pointer.To(state.Tag),
			})

			virtualMachineUpdate := &virtualmachines.VirtualMachineUpdate{
				Properties: &virtualmachines.VirtualMachineProperties{
					ApplicationProfile: &virtualmachines.ApplicationProfile{
						GalleryApplications: applications,
					},
				},
			}

			if err = client.UpdateThenPoll(ctx, *virtualMachineID, pointer.From(virtualMachineUpdate), virtualmachines.DefaultUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating Gallery Application Assignment %q: %+v", virtualMachineID, err)
			}

			metadata.SetID(parse.NewVirtualMachineGalleryApplicationAssignmentID(*virtualMachineID, *galleryApplicationVersionId))
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r VirtualMachineGalleryApplicationAssignmentResource) Read() sdk.ResourceFunc {
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

func (r VirtualMachineGalleryApplicationAssignmentResource) Update() sdk.ResourceFunc {
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

func (r VirtualMachineGalleryApplicationAssignmentResource) Delete() sdk.ResourceFunc {
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
