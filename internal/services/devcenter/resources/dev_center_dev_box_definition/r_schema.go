// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_dev_box_definition

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devboxdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/images"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevCenterDevBoxDefinitionName,
		},

		"location": commonschema.Location(),

		"dev_center_id": commonschema.ResourceIDReferenceRequiredForceNew(&devboxdefinitions.DevCenterId{}),

		"hibernate_support_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"image_reference_id": commonschema.ResourceIDReferenceRequired(&images.GalleryImageId{}),

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
