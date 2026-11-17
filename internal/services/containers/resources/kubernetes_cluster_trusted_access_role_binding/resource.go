// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_trusted_access_role_binding

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/trustedaccess"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return trustedaccess.ValidateTrustedAccessRoleBindingID
}

func (r Resource) ResourceType() string {
	return "azurerm_kubernetes_cluster_trusted_access_role_binding"
}
