// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_gallery_application_assignment

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_virtual_machine_gallery_application_assignment"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.VirtualMachineGalleryApplicationAssignmentIDValidation
}
