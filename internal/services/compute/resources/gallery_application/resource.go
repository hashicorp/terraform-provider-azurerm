// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name gallery_application -properties "name" -service-package-name compute -compare-values "subscription_id:gallery_id,resource_group_name:gallery_id,gallery_name:gallery_id"

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_gallery_application"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return galleryapplications.ValidateApplicationID
}
