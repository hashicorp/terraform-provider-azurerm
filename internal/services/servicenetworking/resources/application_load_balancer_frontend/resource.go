// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_frontend

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/frontendsinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (f Resource) ResourceType() string {
	return "azurerm_application_load_balancer_frontend"
}

func (f Resource) Description() string {
	return "Manages an Application Gateway for Containers Frontend."
}

func (f Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return frontendsinterface.ValidateFrontendID
}
