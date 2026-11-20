// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/availabilitygrouplisteners"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_mssql_virtual_machine_group"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return availabilitygrouplisteners.ValidateSqlVirtualMachineGroupID
}
