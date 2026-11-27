// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_certificate_authority

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementGatewayCertificateAuthoritySchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"api_management_id": {
			Description:  "The ID of the API Management Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
		},

		"certificate_name": schemaz.SchemaApiManagementChildName("The name of the API Management Certificate."),

		"gateway_name": schemaz.SchemaApiManagementChildName("The name of the API Management Gateway."),

		"trusted": {
			Description: "Whether the API Management Gateway Certificate Authority is trusted.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},
	}

	if !features.SixPointOh() {
		s["is_trusted"] = &pluginsdk.Schema{ // azignore:AZS023 - the old name, kept until 6.0 removes it
			Description:   "Whether the API Management Gateway Certificate Authority is trusted.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"trusted"},
			Deprecated:    "`is_trusted` has been deprecated in favour of `trusted` and will be removed in v6.0 of the AzureRM Provider",
		}

		s["trusted"] = &pluginsdk.Schema{
			Description:   "Whether the API Management Gateway Certificate Authority is trusted.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"is_trusted"},
		}
	}

	return s
}
