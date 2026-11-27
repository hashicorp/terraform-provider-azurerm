// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_gallery_application_assignment

import (
	"math"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"gallery_application_version_id": {
			Description:  "The ID of the Gallery Application Version.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: galleryapplicationversions.ValidateApplicationVersionID,
		},

		"virtual_machine_id": {
			Description:  "The ID of the Virtual Machine.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualmachines.ValidateVirtualMachineID,
		},

		"configuration_blob_uri": {
			Description:  "The URI to an Azure Blob that will replace the default configuration for the package if provided.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"order": {
			Description:  "The order in which the packages have to be installed.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntBetween(0, math.MaxInt32),
		},

		"tag": {
			Description:  "A passthrough value for more generic context. This field can be any valid `string` value.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
