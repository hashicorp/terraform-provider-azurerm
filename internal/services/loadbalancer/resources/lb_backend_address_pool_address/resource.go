// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool_address

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const backendAddressPoolResourceName = "azurerm_lb_backend_address_pool"

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_lb_backend_address_pool_address"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.BackendAddressPoolAddressID
}
