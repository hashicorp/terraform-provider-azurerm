// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package marketplace_role_assignment

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.ValidateScopedRoleAssignmentID
}

func (r Resource) ResourceType() string {
	return "azurerm_marketplace_role_assignment"
}
