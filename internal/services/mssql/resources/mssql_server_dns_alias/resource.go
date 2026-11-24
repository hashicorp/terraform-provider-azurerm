// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_dns_alias

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverdnsaliases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = (*Resource)(nil)

func (m Resource) ResourceType() string {
	return "azurerm_mssql_server_dns_alias"
}

func (m Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return serverdnsaliases.ValidateDnsAliasID
}
