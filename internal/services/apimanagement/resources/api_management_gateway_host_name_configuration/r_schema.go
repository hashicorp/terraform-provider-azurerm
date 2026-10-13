// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_host_name_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementGatewayHostNameConfigurationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
		},

		"gateway_name": schemaz.SchemaApiManagementChildName(),

		"certificate_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.AsGeneratedID(certificate.ParseCertificateIDInsensitively),
		},

		"host_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"request_client_certificate_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"http2_enabled": {
			Type:     pluginsdk.TypeBool,
			Default:  true,
			Optional: true,
		},

		"tls10_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tls11_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}
