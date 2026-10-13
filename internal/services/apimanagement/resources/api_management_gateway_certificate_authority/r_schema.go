// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_certificate_authority

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementGatewayCertificateAuthoritySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
		},

		"certificate_name": schemaz.SchemaApiManagementChildName(),

		"gateway_name": schemaz.SchemaApiManagementChildName(),

		"is_trusted": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}
