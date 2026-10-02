// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application_version

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name gallery_application_version -properties "name" -service-package-name compute -compare-values "subscription_id:gallery_application_id,resource_group_name:gallery_application_id,gallery_name:gallery_application_id,application_name:gallery_application_id"

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_gallery_application_version"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return galleryapplicationversions.ValidateApplicationVersionID
}
