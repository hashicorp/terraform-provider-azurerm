// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_active_directory_administrator

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/administratormicrosoftentras"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func postgresqlFlexibleServerActiveDirectoryAdministratorSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"server_name": {
			Description: "The name of the PostgreSQL Flexible Server on which to set the administrator.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"object_id": {
			Description:  "The object ID of a user, Service Principal or security group in the Azure Active Directory tenant set as the Flexible Server Admin.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"principal_name": {
			Description:  "The name of Azure Active Directory principal.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"principal_type": {
			Description: "The type of Azure Active Directory principal.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(administratormicrosoftentras.PrincipalTypeGroup),
				string(administratormicrosoftentras.PrincipalTypeServicePrincipal),
				string(administratormicrosoftentras.PrincipalTypeUser),
			}, false),
		},

		"tenant_id": {
			Description:  "The Azure Tenant ID.",
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}
