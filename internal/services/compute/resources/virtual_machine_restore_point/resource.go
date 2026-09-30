// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_restore_point

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return restorepoints.ValidateRestorePointID
}

func (r Resource) ResourceType() string {
	return "azurerm_virtual_machine_restore_point"
}
