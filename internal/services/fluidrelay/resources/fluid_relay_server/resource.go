// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package fluid_relay_server

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = (*Resource)(nil)

func (s Resource) ResourceType() string {
	return "azurerm_fluid_relay_server"
}

func (s Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fluidrelayservers.ValidateFluidRelayServerID
}
