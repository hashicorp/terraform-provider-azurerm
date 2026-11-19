// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_subnet_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/associationsinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (t Resource) ResourceType() string {
	return "azurerm_application_load_balancer_subnet_association"
}

func (t Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return associationsinterface.ValidateAssociationID
}
