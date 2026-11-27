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
			Description:  "The name of this Dev Center Dev Box Definition.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevCenterDevBoxDefinitionName,
		},

		"location": commonschema.Location(),

		"dev_center_id": commonschema.ResourceIDReferenceRequiredForceNew(&devboxdefinitions.DevCenterId{}, "Dev Center"),

		"hibernate_support_enabled": {
			Description: "Whether the Dev Boxes created with this definition are capable of hibernation.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"image_reference_id": commonschema.ResourceIDReferenceRequired(&images.GalleryImageId{}, "Dev Center Gallery Image"),

		"sku_name": {
			Description:  "The name of the SKU for the Dev Center Dev Box Definition.",
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
