// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const (
	MarketplaceScope = "/providers/Microsoft.Marketplace"
)

var _ sdk.Resource = RoleAssignmentMarketplaceResource{}

type RoleAssignmentMarketplaceResource struct {
	base roleAssignmentBaseResource
}

func (r RoleAssignmentMarketplaceResource) Arguments() map[string]*pluginsdk.Schema {
	return r.base.arguments()
}

func (r RoleAssignmentMarketplaceResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r RoleAssignmentMarketplaceResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), MarketplaceScope)
}

func (r RoleAssignmentMarketplaceResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r RoleAssignmentMarketplaceResource) Read() sdk.ResourceFunc {
	return r.base.readFunc(MarketplaceScope, true)
}

func (r RoleAssignmentMarketplaceResource) ResourceType() string {
	return "azurerm_marketplace_role_assignment"
}

func (r RoleAssignmentMarketplaceResource) ModelObject() interface{} {
	return nil
}

func (r RoleAssignmentMarketplaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.ValidateScopedRoleAssignmentID
}
