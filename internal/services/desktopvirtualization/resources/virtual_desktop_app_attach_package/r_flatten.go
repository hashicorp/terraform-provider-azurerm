// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_app_attach_package

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/appattachpackage"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *appattachpackage.AppAttachPackageId, model *appattachpackage.AppAttachPackage) error {
	state := VirtualDesktopAppAttachPackageModel{
		Name:              id.AppAttachPackageName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)

		props := model.Properties
		if image := props.Image; image != nil {
			state.DisplayName = pointer.From(image.DisplayName)
			state.MsixPackageName = pointer.From(image.PackageFullName)

			if image.ImagePath != nil {
				// Convert from UNC format (`ImagePath`) to URL format (`StorageShareFileId`)
				storageShareFileId := strings.ReplaceAll(pointer.From(image.ImagePath), "\\", "/")
				storageShareFileId = fmt.Sprintf("https:%s", storageShareFileId)
				state.StorageShareFileId = storageShareFileId
			}

			if image.IsRegularRegistration != nil {
				state.RegisterAtLogOnEnabled = pointer.From(image.IsRegularRegistration)
			}

			if image.IsActive != nil {
				state.StateEnabled = pointer.From(image.IsActive)
			}

			state.LastUpdated = pointer.From(image.LastUpdated)
			state.PackageApplications = r.flattenVirtualDesktopAppAttachPackageApplications(image.PackageApplications)
			state.PackageFamilyName = pointer.From(image.PackageFamilyName)
			state.PackageName = pointer.From(image.PackageName)
			state.PackageRelativePath = pointer.From(image.PackageRelativePath)
			state.Version = pointer.From(image.Version)
		}

		state.HostPoolIds = pointer.From(props.HostPoolReferences)
		for i, hostPoolId := range state.HostPoolIds {
			parsedHostPoolId, err := hostpool.ParseHostPoolIDInsensitively(hostPoolId)
			if err != nil {
				return fmt.Errorf("parsing host pool ID %s: %+v", hostPoolId, err)
			}
			state.HostPoolIds[i] = parsedHostPoolId.ID()
		}

		if props.FailHealthCheckOnStagingFailure != nil {
			state.HealthCheckStatusOnFailure = pointer.FromEnum(props.FailHealthCheckOnStagingFailure)
		}

		state.Tags = pointer.From(model.Tags)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func (r Resource) flattenVirtualDesktopAppAttachPackageApplications(inputs *[]appattachpackage.MsixPackageApplications) []MsixPackageApplicationModel {
	outputs := make([]MsixPackageApplicationModel, 0)
	if inputs == nil {
		return outputs
	}

	for _, input := range pointer.From(inputs) {
		outputs = append(outputs, MsixPackageApplicationModel{
			AppId:          pointer.From(input.AppId),
			AppUserModelID: pointer.From(input.AppUserModelID),
			Description:    pointer.From(input.Description),
			FriendlyName:   pointer.From(input.FriendlyName),
			IconImageName:  pointer.From(input.IconImageName),
			RawIcon:        pointer.From(input.RawIcon),
			RawPng:         pointer.From(input.RawPng),
		})
	}

	return outputs
}
