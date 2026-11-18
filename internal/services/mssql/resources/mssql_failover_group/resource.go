// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_failover_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/failovergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_mssql_failover_group"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return failovergroups.ValidateFailoverGroupID
}
