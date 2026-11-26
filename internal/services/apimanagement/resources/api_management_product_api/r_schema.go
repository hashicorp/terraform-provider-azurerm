// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_product_api

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementProductApiSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_name": schemaz.SchemaApiManagementApiName("The Name of the API Management API within the API Management Service."),

		"product_id": schemaz.SchemaApiManagementChildName("The ID of the API Management Product within the API Management Service."),

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),
	}
}
