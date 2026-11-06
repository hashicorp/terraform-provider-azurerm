// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_app_attach_package

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/appattachpackage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var (
	_ sdk.Resource                  = Resource{}
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_virtual_desktop_app_attach_package"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appattachpackage.ValidateAppAttachPackageID
}

func (r Resource) Identity() resourceids.ResourceId {
	return &appattachpackage.AppAttachPackageId{}
}
