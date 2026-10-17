// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_group_user

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementGroupUserSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"user_id": schemaz.SchemaApiManagementChildName(),

		"group_name": schemaz.SchemaApiManagementChildName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),
	}
}
