// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_private_link_scope

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

// Resource is a Resource implementation for the Azure Arc Private Link Scope resource
type Resource struct{}

func (a Resource) ResourceType() string {
	return "azurerm_arc_private_link_scope"
}

func (a Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return privatelinkscopes.ValidateProviderPrivateLinkScopeID
}
