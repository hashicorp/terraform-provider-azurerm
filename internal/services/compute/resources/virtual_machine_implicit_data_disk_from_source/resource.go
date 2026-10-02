// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_implicit_data_disk_from_source

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
	_ sdk.ResourceWithCustomizeDiff  = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataDiskID
}

func (r Resource) ResourceType() string {
	return "azurerm_virtual_machine_implicit_data_disk_from_source"
}
