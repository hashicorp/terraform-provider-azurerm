// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_connection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.ResourceWithUpdate = Resource{}

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_app_service_connection"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return servicelinker.ValidateScopedLinkerID
}
