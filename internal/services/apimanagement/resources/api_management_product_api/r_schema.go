// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_product_api

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementProductApiSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_name": schemaz.SchemaApiManagementApiName(),

		"product_id": schemaz.SchemaApiManagementChildName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),
	}
}
