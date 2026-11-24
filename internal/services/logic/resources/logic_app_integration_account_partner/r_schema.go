// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_partner

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logicAppIntegrationAccountPartnerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Logic App Integration Account Partner.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountPartnerName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"integration_account_name": {
			Description:  "The name of the Logic App Integration Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountName(),
		},

		"business_identity": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"qualifier": {
						Description:  "The authenticating body that provides unique business identities to organizations.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.IntegrationAccountPartnerBusinessIdentityQualifier(),
					},

					"value": {
						Description:  "The value that identifies the documents that your Logic Apps receive.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.IntegrationAccountPartnerBusinessIdentityValue(),
					},
				},
			},
		},

		"metadata": {
			Description:      "A JSON mapping of any Metadata for this Logic App Integration Account Partner.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},
	}
}
