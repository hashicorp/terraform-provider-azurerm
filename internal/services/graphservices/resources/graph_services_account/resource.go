// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package graph_services_account

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/graphservices/2023-04-13/graphservicesprods"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return graphservicesprods.ValidateAccountID
}

func (r Resource) ResourceType() string {
	return "azurerm_graph_services_account"
}
