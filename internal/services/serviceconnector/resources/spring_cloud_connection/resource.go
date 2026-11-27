// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_connection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_spring_cloud_connection"
}

func (r Resource) Description() string {
	return "Manages a service connector for Spring Cloud app."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return servicelinker.ValidateScopedLinkerID
}
