// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_active_directory_administrator

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_mysql_flexible_server_active_directory_administrator"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.FlexibleServerAzureActiveDirectoryAdministratorID
}
