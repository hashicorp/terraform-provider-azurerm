// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (t Resource) ResourceType() string {
	return "azurerm_application_load_balancer"
}

func (t Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return trafficcontrollerinterface.ValidateTrafficControllerID
}
