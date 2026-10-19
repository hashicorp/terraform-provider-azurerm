// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_management_policy

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.ValidateRoleManagementPolicyId
}

func (r Resource) ResourceType() string {
	return "azurerm_role_management_policy"
}
