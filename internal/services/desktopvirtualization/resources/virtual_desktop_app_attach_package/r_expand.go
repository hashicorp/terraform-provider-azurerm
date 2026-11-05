// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_app_attach_package

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/appattachpackage"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/msiximage"
)

func (r Resource) expandVirtualDesktopAppAttachPackageImage(model VirtualDesktopAppAttachPackageModel, msixImageProperties *msiximage.ExpandMsixImageProperties) *appattachpackage.AppAttachPackageInfoProperties {
	// Convert from URL format (`StorageShareFileId`) to UNC format (`ImagePath`)
	imagePath := strings.ReplaceAll(model.StorageShareFileId, "/", "\\")
	imagePath = strings.TrimPrefix(imagePath, "https:")

	return &appattachpackage.AppAttachPackageInfoProperties{
		DisplayName:           pointer.To(model.DisplayName),
		PackageFullName:       pointer.To(model.MsixPackageName),
		ImagePath:             pointer.To(imagePath),
		IsRegularRegistration: pointer.To(model.RegisterAtLogOnEnabled),
		IsActive:              pointer.To(model.StateEnabled),
		LastUpdated:           msixImageProperties.LastUpdated,
		PackageApplications:   r.expandVirtualDesktopAppAttachPackageApplications(msixImageProperties.PackageApplications),
		PackageFamilyName:     msixImageProperties.PackageFamilyName,
		PackageName:           msixImageProperties.PackageName,
		PackageRelativePath:   msixImageProperties.PackageRelativePath,
		Version:               msixImageProperties.Version,
	}
}

func (r Resource) expandVirtualDesktopAppAttachPackageApplications(inputs *[]msiximage.MsixPackageApplications) *[]appattachpackage.MsixPackageApplications {
	outputs := make([]appattachpackage.MsixPackageApplications, 0)
	if inputs == nil {
		return pointer.To(outputs)
	}

	for _, input := range *inputs {
		outputs = append(outputs, appattachpackage.MsixPackageApplications{
			AppId:          input.AppId,
			AppUserModelID: input.AppUserModelID,
			Description:    input.Description,
			FriendlyName:   input.FriendlyName,
			IconImageName:  input.IconImageName,
			RawIcon:        input.RawIcon,
			RawPng:         input.RawPng,
		})
	}

	return pointer.To(outputs)
}
