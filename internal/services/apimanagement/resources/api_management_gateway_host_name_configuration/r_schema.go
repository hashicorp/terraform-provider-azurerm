// Copyright IBM Corp. 2014, 2026
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
		"name": schemaz.SchemaApiManagementChildName("The name of the API Management Gateway Host Name Configuration."),

		"api_management_id": {
			Description:  "The ID of the API Management Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
		},

		"gateway_name": schemaz.SchemaApiManagementChildName("The name of the API Management Gateway."),

		"certificate_id": {
			Description:  "The certificate ID to be used for TLS connection establishment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.AsGeneratedID(certificate.ParseCertificateIDInsensitively),
		},

		"host_name": {
			Description:  "The host name to use for the API Management Gateway Host Name Configuration.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"request_client_certificate_enabled": {
			Description: "Whether the API Management Gateway requests a client certificate.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"http2_enabled": {
			Description: "Whether HTTP/2.0 is supported.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"tls10_enabled": {
			Description: "Whether TLS 1.0 is supported.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"tls11_enabled": {
			Description: "Whether TLS 1.1 is supported.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},
	}
}
