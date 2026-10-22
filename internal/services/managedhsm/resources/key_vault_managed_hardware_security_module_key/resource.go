// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_key

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedHSMDataPlaneVersionlessKeyID
}

func (r Resource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_key"
}

var _ deleteAndPurgeNestedItem = deleteAndPurgeKey{}
