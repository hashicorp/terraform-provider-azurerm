// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_agreement

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountagreements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logicAppIntegrationAccountAgreementSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Logic App Integration Account Agreement.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountAgreementName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"integration_account_name": {
			Description:  "The name of the Logic App Integration Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountName(),
		},

		"agreement_type": {
			Description: "The type of the Logic App Integration Account Agreement.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(integrationaccountagreements.AgreementTypeASTwo),
				string(integrationaccountagreements.AgreementTypeXOneTwo),
				string(integrationaccountagreements.AgreementTypeEdifact),
			}, false),
		},

		"content": {
			Description:      "The content of the Logic App Integration Account Agreement.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"guest_identity": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"qualifier": {
						Description:  "The authenticating body that provides unique guest identities to organizations.",
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

		"guest_partner_name": {
			Description:  "The name of the guest Logic App Integration Account Partner.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.IntegrationAccountPartnerName(),
		},

		"host_identity": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"qualifier": {
						Description:  "The authenticating body that provides unique host identities to organizations.",
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

		"host_partner_name": {
			Description:  "The name of the host Logic App Integration Account Partner.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.IntegrationAccountPartnerName(),
		},

		"metadata": {
			Description: "The metadata of the Logic App Integration Account Agreement.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
